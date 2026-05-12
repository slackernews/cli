package cmd

import (
	"github.com/spf13/cobra"
)

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search links by keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient(allowInsecure)
		if err != nil {
			return err
		}

		links, err := client.SearchLinks(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		return printLinks(cmd, links, searchJSON)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output as JSON")
}
