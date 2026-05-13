package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unvoteCmd = &cobra.Command{
	Use:   "unvote <id>",
	Short: "Remove your upvote from a link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient(allowInsecure)
		if err != nil {
			return err
		}

		if err := client.Unvote(cmd.Context(), args[0]); err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Vote removed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(unvoteCmd)
}
