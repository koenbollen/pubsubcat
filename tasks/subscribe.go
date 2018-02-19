package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/koenbollen/pubsubcat/utils"
	"github.com/pkg/errors"
)

const temporarySubscriptionTemplate = "pubsubcat-%s-%x-%d"

var newline = []byte("\n")

// Subscribe will connect to pubsub, created a temporary subscription on the
// given topic and listens for message to output to Stdout.
//
// Cancel the given context to stop.
func Subscribe(ctx context.Context, client *pubsub.Client, topicID string) error {
	topic := client.Topic(topicID)

	subscription, err := createTemporarySubscription(ctx, client, topic)
	if err != nil {
		return err
	}

	// TODO: Log here that the subcription is created.

	// Receive messages on subscription and output them to Stdout:
	var mu sync.Mutex
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		os.Stdout.Write(msg.Data)
		os.Stdout.Write(newline)
		msg.Ack()
	})
	if err != nil {
		return errors.Wrapf(err, "error whilst receving messages from %s", subscription.ID())
	}

	// TODO: log the cleanup action

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
		return nil, errors.Wrap(err, "failed to create temporary subscription")
	}
	return subscription, nil
}

func cleanupTemporarySubscription(subscription *pubsub.Subscription) error {
	// Use a new context because the old might be cancelled already:
	deleteContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	utils.CancelOnSignal(deleteContext, cancel, os.Interrupt)
	defer cancel()
	err := subscription.Delete(deleteContext)
	if err != nil {
		return errors.Wrapf(err, "failed to cleanup subscription %v", subscription.ID())
	}
	return nil
}
