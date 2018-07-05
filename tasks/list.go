package tasks

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

// ListTopics outputs all the topics fetched from the pubsub.Client.
func ListTopics(ctx context.Context, client *pubsub.Client) error {
	iter := client.Topics(ctx)
	for {
		t, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(t)
	}
	return nil
}

// ListSubscriptions outputs the request topic and all it's subscriptions
func ListSubscriptions(ctx context.Context, client *pubsub.Client, topicID, projectID string) error {
	topic := client.TopicInProject(topicID, projectID)
	fmt.Println(topic)
	iter := topic.Subscriptions(ctx)
	for {
		t, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(" ", t)
	}
	return nil
}
