package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment <id> <body>",
	Short: "Comment on a link",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient(allowInsecure)
		if err != nil {
			return err
		}

		if err := client.Comment(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Comment posted successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
