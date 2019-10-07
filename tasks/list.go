package tasks

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/fatih/color"
	"google.golang.org/api/iterator"
)

var faint = color.New(color.Faint)
var bold = color.New(color.Bold)

// ListTopics outputs all the topics fetched from the pubsub.Client.
func ListTopics(ctx context.Context, client *pubsub.Client, recursive bool) error {
	iter := client.Topics(ctx)
	for {
		t, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if recursive {
			if err := ListSubscriptions(ctx, client, t.ID()); err != nil {
				return err
			}
		} else {
			fmt.Println(colorizeResource(t.String()))
		}
	}
	return nil
}

// ListSubscriptions outputs the request topic and all it's subscriptions
func ListSubscriptions(ctx context.Context, client *pubsub.Client, topicID string) error {
	topic := client.Topic(topicID)
	fmt.Println(colorizeResource(topic.String()))
	iter := topic.Subscriptions(ctx)
	for {
		s, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(" ", colorizeResource(s.String()))
	}
	return nil
}

func colorizeResource(resource string) string {
	if color.NoColor {
		return resource
	}
	parts := strings.Split(resource, "/")
	if parts[0] != "projects" {
		return resource
	}
	project := parts[1]
	resourceType := parts[2]
	resourceID := parts[3]
	return faint.Sprintf("projects/%s/%s/", project, resourceType) + bold.Sprint(resourceID)
}
