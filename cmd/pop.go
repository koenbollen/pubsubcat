package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop [flags] TOPIC",
	Short: "alias for: `pubsubcat subscribe --count 1 TOPIC`",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("no TOPIC given")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flags().Set("count", "1")
		subscribeCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(popCmd)

	// Flags copied from subscribeCmd in subscribe.go
}
