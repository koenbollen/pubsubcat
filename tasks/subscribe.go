package tasks

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/koenbollen/pubsubcat/utils"
)

const temporarySubscriptionTemplate = "pubsubcat-%s-%x-%d"

var newline = []byte("\n")

// SubscribeParams allows config over the Subscribe task.
type SubscribeParams struct {
	Verbosity      int
	TopicID        string
	SubscriptionID string
	Count          int
	NoCleanup      bool
}

// Subscribe will connect to pubsub, created a temporary subscription on the
// given topic and listens for message to output to Stdout.
//
// Cancel the given context to stop.
func Subscribe(ctx context.Context, client *pubsub.Client, params SubscribeParams) error {
	topic := client.Topic(params.TopicID)

	var subscription *pubsub.Subscription
	if params.SubscriptionID != "" {
		params.NoCleanup = true
		subscription = client.Subscription(params.SubscriptionID)
	} else {
		if params.Verbosity >= 2 {
			log.Println("] creating temporary subscription on topic:", topic.ID())
		}

		var err error
		subscription, err = createTemporarySubscription(ctx, client, topic)
		if err != nil {
			return err
		}
	}

	if params.Verbosity >= 1 {
		log.Printf("] listening on topic %q using subscription %q",
			topic.ID(),
			subscription.ID())
	}

	// Receive messages on subscription and output them to Stdout:
	var mu sync.Mutex
	received := 0
	cctx, cancel := context.WithCancel(ctx)
	err := subscription.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock() // TODO: Maybe only lock when --counting or --blocking?
		defer mu.Unlock()

		received++
		if params.Count > 0 && received > params.Count {
			cancel()
			msg.Nack()
			return
		}

		if params.Verbosity >= 2 {
			log.Printf("] received %q at %v", msg.ID, msg.PublishTime)
		}
		if params.Verbosity >= 3 {
			log.Printf("]   attributes: %v", msg.Attributes)
		}
		os.Stdout.Write(msg.Data)
		os.Stdout.Write(newline)
		msg.Ack()
		if params.Verbosity >= 4 {
			log.Printf("] message %v acknowledged", msg.ID)
		}
	})
	if err != nil {
		return fmt.Errorf("error whilst receving messages from %s: %w", subscription.ID(), err)
	}

	if params.NoCleanup {
		return nil
	}

	if params.Verbosity >= 1 {
		log.Printf("] stopped receiving, cleaning up temporary subscription")
	}
	return cleanupTemporarySubscription(subscription)
}

func createTemporarySubscription(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	subscriptionID := fmt.Sprintf(temporarySubscriptionTemplate, topic.ID(), rand.Int(), time.Now().Unix())
	subscription, err := client.CreateSubscription(ctx, subscriptionID, pubsub.SubscriptionConfig{
		Topic:               topic,
		AckDeadline:         10 * time.Second,
		RetentionDuration:   10 * time.Minute,
		RetainAckedMessages: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary subscription: %w", err)
	}
	return subscription, nil
}

func cleanupTemporarySubscription(subscription *pubsub.Subscription) error {
	// Use a new context because the old might be cancelled already:
	deleteContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	utils.CancelOnSignal(deleteContext, cancel, os.Interrupt)
	defer cancel()
	err := subscription.Delete(deleteContext)
	if err != nil {
		return fmt.Errorf("failed to cleanup subscription %v: %w", subscription.ID(), err)
	}
	return nil
}
