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
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// TODO: Test if msg.Data needs to be copied if non-blocking
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
		_ = result
		// if blocking {
		// 	id, err := result.Get(ctx)
		// 	if err != nil {
		// 		return errors.Wrap(err, "failed to publish message")
		// 	}
		// 	if params.Verbosity >= 3 {
		// 		log.Println("] published message:", id)
		// 	}
		// }
		msg.Ack()
		if params.Verbosity >= 3 {
			log.Printf("] message %v acknowledged", msg.ID)
		}
	})
	if err != nil {
		return errors.Wrapf(err, "error whilst receving messages from: %s", subscription.ID())
	}
	outTopic.Stop()

	if params.Verbosity >= 1 {
		log.Printf("] pipe stopped, cleaning up temporary subscription")
	}
	return cleanupTemporarySubscription(subscription)
}
