package tasks

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CleanParams allows config over the Clean task.
type CleanParams struct {
	Verbosity int
	TopicID   string
}

// CleanTopic will look at all subscriptions in the given topic and remove
// old pubsubcat temporary topics.
func CleanTopic(ctx context.Context, client *pubsub.Client, params CleanParams) error {
	topic := client.Topic(params.TopicID)
	subscriptions := topic.Subscriptions(ctx)
	if params.Verbosity >= 2 {
		log.Println("] checking for lingering pubsubcat subscriptions to clean")
	}
	for {
		s, err := subscriptions.Next() // nolint
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrapf(err, "failed to iterate subscriptions for topic: %s", params.TopicID)
		}
		if params.Verbosity >= 3 {
			log.Println("] checking subscription", s)
		}
		var subscriptionTopicID string
		var random int
		var timestamp int64
		if n, _ := fmt.Sscanf(s.ID(), "pubsubcat-%s-%x-%d", &subscriptionTopicID, &random, &timestamp); n == 3 {
			if subscriptionTopicID != params.TopicID { // sanity check, should not happen
				continue
			}
			createdAt := time.Unix(timestamp, 0)
			if createdAt.Before(time.Now().Add(-24 * time.Hour)) {
				if params.Verbosity >= 1 {
					log.Println("] cleaning lingering subscription:", s.ID())
				}
				err = s.Delete(ctx)
				if err != nil {
					return errors.Wrapf(err, "failed to delete old subscription for topic %s", params.TopicID)
				}
			}
		}
	}
	return nil
}
