package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upvoteCmd = &cobra.Command{
	Use:   "upvote <id>",
	Short: "Upvote a link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient(allowInsecure)
		if err != nil {
			return err
		}

		if err := client.Upvote(cmd.Context(), args[0]); err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Upvoted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upvoteCmd)
}
