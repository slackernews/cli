package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/slackernews/cli/pkg/formatters"
	"github.com/spf13/cobra"
)

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search links by keyword",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(globalInsecure)
		if err != nil {
			return err
		}

		links, err := client.SearchLinks(args[0])
		if err != nil {
			return err
		}

		if len(links) == 0 {
			if searchJSON {
				fmt.Println("[]")
				return nil
			}
			fmt.Println("No links found")
			return nil
		}

		if searchJSON {
			out, err := formatters.FormatJSON(links)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		}

		formatters.FormatTable(links)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output as JSON")
}
