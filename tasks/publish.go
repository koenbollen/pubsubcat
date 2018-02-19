package tasks

import (
	"bufio"
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// Publish lines read from os.Stdin to the given topic using the given client.
func Publish(ctx context.Context, client *pubsub.Client, topicID string) error {
	topic := client.Topic(topicID)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// TODO: Can skip the copybuffer when --blocking is on
		buffer := make([]byte, len(scanner.Bytes()))
		copy(buffer, scanner.Bytes())
		result := topic.Publish(ctx, &pubsub.Message{
			Data: buffer,
		})
		_ = result
		// if blocking {
		// 	<-result.Ready()
		// }
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "failed to read from stdin")
	}
	topic.Stop()

	return nil
}
