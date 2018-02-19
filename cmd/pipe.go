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

// pipeCmd represents the pipe command
var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "A brief description of your command",
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
		// TODO: Fail if two projects.

		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create pubsub client: %v", err)
		}
		defer client.Close()

		err = tasks.CleanTopic(ctx, client, args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to clean old subscriptions: %v", err)
		}

		err = tasks.Pipe(ctx, client, args[0], args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pipeCmd)

	// TODO: Support --count=N, -c N
	// TODO: Support --no-cleanup
	// TODO: Support --blocking, -b
}
