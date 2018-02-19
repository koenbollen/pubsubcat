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

		err = tasks.CleanTopic(ctx, client, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to clean old subscriptions: %v", err)
		}

		err = tasks.Subscribe(ctx, client, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// TODO: Support --output=FILE, -o FILE
	// TODO: Support --count=N, -c N
	// TODO: Support --unbuffered, -u
	// TODO: Support --no-cleanup
}
