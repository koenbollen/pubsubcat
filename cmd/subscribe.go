package cmd

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/koenbollen/pubsubcat/tasks"
	"github.com/koenbollen/pubsubcat/utils"
	"github.com/spf13/cobra"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to a topic using a temporary subscription",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity := GetVerbosity(cmd.Flags())

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

		count, _ := cmd.Flags().GetInt("count")
		subscribeParams := tasks.SubscribeParams{
			TopicID:   topicID,
			Verbosity: verbosity,
			Count:     count,
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

	// TODO: Support --output=FILE, -o FILE
	// TODO: Support --unbuffered, -u
	// TODO: Support --no-cleanup
	// TODO: Support --subscription mycustomsubscription (or as a second positional argument)
	// TODO: Support https://godoc.org/cloud.google.com/go/pubsub#Subscription.SeekToTime
}
