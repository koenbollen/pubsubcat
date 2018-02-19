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

// CleanTopic will look at all subscriptions in the given topic and remove
// old pubsubcat temporary topics.
func CleanTopic(ctx context.Context, client *pubsub.Client, topicID string) error {
	topic := client.Topic(topicID)
	subscriptions := topic.Subscriptions(ctx)
	log.Println("Checking for old subscriptions to clean")
	for {
		s, err := subscriptions.Next() // nolint
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrapf(err, "failed to iterate subscriptions for topic %s", topicID)
		}
		log.Println("checking subscription", s)
		var subscriptionTopicID string
		var random int
		var timestamp int64
		if n, _ := fmt.Sscanf(s.ID(), "pubsubcat-%s-%x-%d", &subscriptionTopicID, &random, &timestamp); n == 3 {
			if subscriptionTopicID != topicID { // sanity check, should not happen
				continue
			}
			createdAt := time.Unix(timestamp, 0)
			if createdAt.Before(time.Now().Add(-24 * time.Hour)) {
				log.Println("Cleaning old subscription", s.ID(), random, timestamp, n)
				err = s.Delete(ctx)
				if err != nil {
					return errors.Wrapf(err, "failed to delete old subscription for topic %s", topicID)
				}
			}
		}
	}
	return nil
}
