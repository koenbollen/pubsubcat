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
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity := GetVerbosity(cmd.Flags())
		blocking, _ := cmd.Flags().GetBool("blocking")
		count, _ := cmd.Flags().GetInt("count")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.CancelOnSignal(ctx, cancel, os.Interrupt)

		inProjectID, inTopicID, err := utils.DetermineProject(args[0], globalProjectID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		outProjectID, outTopicID, err := utils.DetermineProject(args[1], inProjectID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if inProjectID != outProjectID {
			fmt.Fprintln(os.Stderr, "unable to pipe between different Google Cloud Projects")
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
			TopicID:   inTopicID,
		}
		err = tasks.CleanTopic(ctx, client, cleanParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to clean old subscriptions: %v\n", err)
			os.Exit(1)
		}

		pipeParams := tasks.PipeParams{
			Verbosity:  verbosity,
			InTopicID:  inTopicID,
			OutTopicID: outTopicID,
			Blocking:   blocking,
			Count:      count,
		}
		err = tasks.Pipe(ctx, client, pipeParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to subscribe: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pipeCmd)

	pipeCmd.Flags().BoolP("blocking", "b", false, "wait for server on each message")
	pipeCmd.Flags().IntP("count", "c", 0, "only read <int> messages, then exit")

	// TODO: Support --no-cleanup
}
