package tasks

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Pipe will create a random subscription on inTopicID and forward all messages
// to the given outTopicID.
func Pipe(ctx context.Context, client *pubsub.Client, inTopicID, outTopicID string) error {
	inTopic := client.Topic(inTopicID)
	outTopic := client.Topic(outTopicID)

	subscription, err := createTemporarySubscription(ctx, client, inTopic)
	if err != nil {
		return err
	}

	// TODO: Log here that the subcription is created.

	// Receive messages on subscription and output them to Stdout:
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// TODO: Log msg.ID if -vvv
		// TODO: Test if msg.Data needs to be copied if non-blocking
		log.Println(msg)
		result := outTopic.Publish(ctx, &pubsub.Message{
			Data:        msg.Data,
			ID:          msg.ID,
			Attributes:  msg.Attributes,
			PublishTime: msg.PublishTime,
		})
		_ = result
		// if blocking {
		// 	<-result.Ready()
		// }
		msg.Ack()
	})
	if err != nil {
		return errors.Wrapf(err, "error whilst receving messages from %s", subscription.ID())
	}

	// TODO: log the cleanup action

	return cleanupTemporarySubscription(subscription)
}
