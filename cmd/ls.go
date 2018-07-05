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

// lsCmd represents the publish command
var lsCmd = &cobra.Command{
	Use:   "ls [flags] [TOPIC]",
	Short: "List topics and subscriptions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		recursive, _ := cmd.Flags().GetBool("recursive")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.CancelOnSignal(ctx, cancel, os.Interrupt)

		resource := ""
		if len(args) > 0 {
			resource = args[0]
		}
		projectID, topicID, err := utils.DetermineProject(resource, globalProjectID)
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

		if topicID != "" {

			err := tasks.ListSubscriptions(ctx, client, topicID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to list subscriptions: %v\n", err)
				os.Exit(1)
			}

		} else {

			err := tasks.ListTopics(ctx, client, recursive)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to list topics: %v\n", err)
				os.Exit(1)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("recursive", "r", false, "also output all subscriptions")

	// TODO: Support -s, --short that doesn't output the entire project/resource/topic syntax
}
