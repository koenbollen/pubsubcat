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

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish input lines as messages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Fprintf(os.Stderr, "failed to create pubsub client: %v\n", err)
			os.Exit(1)
		}
		defer client.Close()

		err = tasks.Publish(ctx, client, topicID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to publish messages: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// TODO: Checkout these for defaults and --flags
	// topic.PublishSettings = pubsub.PublishSettings{
	// 	ByteThreshold:  5000,
	// 	CountThreshold: 10,
	// 	DelayThreshold: 100 * time.Millisecond,
	// }
	// TODO: Support --blocking, -b
}
