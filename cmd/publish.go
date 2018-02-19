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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.CancelOnSignal(ctx, cancel, os.Interrupt)

		// TODO: Determine projectID/topicID by args[] and/or from default project.
		//       Support /project/MY_PROJECT_ID/topics/MY_TOPIC override

		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create pubsub client: %v", err)
		}
		defer client.Close()

		tasks.Publish(ctx, client, args[0])
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
