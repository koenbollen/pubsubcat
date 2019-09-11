package tasks

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

// PublishParams allows config over the Publish task.
type PublishParams struct {
	Verbosity  int
	TopicID    string
	Blocking   bool
	Attributes []string
}

// Publish lines read from os.Stdin to the given topic using the given client.
func Publish(ctx context.Context, client *pubsub.Client, params PublishParams) error {
	topic := client.Topic(params.TopicID)

	if params.Verbosity >= 1 {
		log.Println("] publishing lines from stdin to topic", topic.ID())
	}

	var attributes map[string]string
	if params.Attributes != nil {
		attributes = make(map[string]string)
		for _, line := range params.Attributes {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				attributes[parts[0]] = parts[1]
			}
		}
	}

	if params.Verbosity >= 3 {
		log.Printf("]   attributes: %v", attributes)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// TODO: Can skip the copybuffer when --blocking is on
		buffer := make([]byte, len(scanner.Bytes()))
		copy(buffer, scanner.Bytes())
		result := topic.Publish(ctx, &pubsub.Message{
			Data:       buffer,
			Attributes: attributes,
		})
		if params.Blocking {
			id, err := result.Get(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to publish message")
			}
			if params.Verbosity >= 2 {
				log.Println("] published message:", id)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "failed to read from stdin")
	}

	if params.Verbosity >= 1 {
		log.Println("] eof, syncing")
	}
	topic.Stop()

	if params.Verbosity >= 2 {
		log.Println("] done")
	}
	return nil
}
