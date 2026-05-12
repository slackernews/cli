package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/spf13/cobra"
)

var upvoteCmd = &cobra.Command{
	Use:   "upvote <id>",
	Short: "Upvote a link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(globalInsecure)
		if err != nil {
			return err
		}

		if err := client.Upvote(args[0]); err != nil {
			return err
		}

		fmt.Println("Upvoted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upvoteCmd)
}
