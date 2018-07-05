package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/koenbollen/pubsubcat/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the publish command
var createCmd = &cobra.Command{
	Use:   "create [flags] TOPIC [SUBSCRIPTION]",
	Short: "Create topics and subscriptions",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no TOPIC given")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.CancelOnSignal(ctx, cancel, os.Interrupt)

		projectID, topicID, err := utils.DetermineProject(args[0], globalProjectID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to pubsub: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		if len(args) == 1 {

			topic, err := client.CreateTopic(ctx, topicID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create topic: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(topic)

		} else {

			topic := client.TopicInProject(topicID, projectID)
			config := pubsub.SubscriptionConfig{
				Topic: topic,
			}
			subscription, err := client.CreateSubscription(ctx, args[1], config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create subscription: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(subscription)

		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// TODO: Support subscription creation config (eg RetainAckedMessages)
	// TODO: if tty, make resource gray except the topicID part.
}
