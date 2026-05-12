package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/spf13/cobra"
)

var unvoteCmd = &cobra.Command{
	Use:   "unvote <id>",
	Short: "Remove your upvote from a link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(globalInsecure)
		if err != nil {
			return err
		}

		if err := client.Unvote(args[0]); err != nil {
			return err
		}

		fmt.Println("Vote removed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(unvoteCmd)
}
