package tasks

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// PipeParams allows config over the Pipe task.
type PipeParams struct {
	Verbosity  int
	InTopicID  string
	OutTopicID string
	Blocking   bool
	Count      int
	NoCleanup  bool
}

// Pipe will create a random subscription on inTopicID and forward all messages
// to the given outTopicID.
func Pipe(ctx context.Context, client *pubsub.Client, params PipeParams) error {
	inTopic := client.Topic(params.InTopicID)
	outTopic := client.Topic(params.OutTopicID)

	if params.Verbosity >= 2 {
		log.Println("] creating temporary subscription on topic:", inTopic.ID())
	}
	subscription, err := createTemporarySubscription(ctx, client, inTopic)
	if err != nil {
		return err
	}

	if params.Verbosity >= 1 {
		log.Printf("] listening on topic %q using subscription %q, forwarding to %q",
			inTopic.ID(),
			subscription.ID(),
			outTopic.ID())
	}

	// Receive messages on subscription and output them to Stdout:
	received := 0
	var midReceiveError error
	cctx, cancel := context.WithCancel(ctx)
	err = subscription.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		// TODO: Test if msg.Data needs to be copied if non-blocking
		// TODO: Test if we need a lock here, maybe only if --count?

		received++
		if params.Count > 0 && received > params.Count {
			cancel()
			msg.Nack()
			return
		}

		if params.Verbosity >= 3 {
			log.Printf("] received %q at %v", msg.ID, msg.PublishTime)
		}
		if params.Verbosity >= 4 {
			log.Printf("]   attributes: %v", msg.Attributes)
		}
		result := outTopic.Publish(ctx, &pubsub.Message{
			Data:        msg.Data,
			ID:          msg.ID,
			Attributes:  msg.Attributes,
			PublishTime: msg.PublishTime,
		})
		if params.Blocking {
			var id string
			id, err = result.Get(ctx)
			if err != nil {
				midReceiveError = errors.Wrap(err, "failed to publish message")
				cancel()
				msg.Nack()
				return
			}
			if params.Verbosity >= 3 {
				log.Println("] published message:", id)
			}
		}
		msg.Ack()
		if params.Verbosity >= 3 {
			log.Printf("] message %v acknowledged", msg.ID)
		}
	})
	outTopic.Stop()
	if err != nil {
		return errors.Wrapf(err, "error whilst receving messages from: %s", subscription.ID())
	}
	if midReceiveError != nil {
		return errors.Wrapf(err, "error whilst sending messages to: %s", outTopic.ID())
	}

	if params.NoCleanup {
		return nil
	}

	if params.Verbosity >= 1 {
		log.Printf("] pipe stopped, cleaning up temporary subscription")
	}
	return cleanupTemporarySubscription(subscription)
}
