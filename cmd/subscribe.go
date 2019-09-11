package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/koenbollen/pubsubcat/tasks"
	"github.com/koenbollen/pubsubcat/utils"
	"github.com/spf13/cobra"
)

var subscribeCmd = &cobra.Command{
	Use:   "subscribe [flags] TOPIC [SUBSCRIPTION]",
	Short: "Subscribe to a topic using a temporary subscription",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no TOPIC given")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		verbosity := GetVerbosity(cmd.Flags())
		count, _ := cmd.Flags().GetInt("count")
		noCleanup, _ := cmd.Flags().GetBool("no-cleanup")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.CancelOnSignal(ctx, cancel, os.Interrupt)

		inProjectID, topicID, err := utils.DetermineProject(args[0], globalProjectID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		client, err := pubsub.NewClient(ctx, inProjectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to pubsub: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		cleanParams := tasks.CleanParams{
			Verbosity: verbosity,
			TopicID:   topicID,
		}
		err = tasks.CleanTopic(ctx, client, cleanParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to clean old subscriptions: %v\n", err)
			os.Exit(1)
		}

		subscriptionID := ""
		if len(args) >= 2 {
			subscriptionID = args[1]
		}
		subscribeParams := tasks.SubscribeParams{
			TopicID:        topicID,
			SubscriptionID: subscriptionID,
			Verbosity:      verbosity,
			Count:          count,
			NoCleanup:      noCleanup,
		}
		err = tasks.Subscribe(ctx, client, subscribeParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe to topic: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)

	subscribeCmd.Flags().IntP("count", "c", 0, "only read <int> messages, then exit")
	subscribeCmd.Flags().Bool("no-cleanup", false, "do not cleanup temporary subscription")

	// TODO: Support --output=FILE, -o FILE
	// TODO: Support --unbuffered, -u
	// TODO: Support --subscription mycustomsubscription (or as a second positional argument)
	// TODO: Support https://godoc.org/cloud.google.com/go/pubsub#Subscription.SeekToTime

	// This should be in pop.go but wanted to ensure the order of init:
	popCmd.Flags().AddFlagSet(subscribeCmd.Flags())
	popCmd.Flag("count").Hidden = true
}
