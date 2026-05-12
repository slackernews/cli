package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/slackernews/cli/pkg/formatters"
	"github.com/spf13/cobra"
)

var (
	topDuration string
	topJSON     bool
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "List top-ranked links",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient(globalInsecure)
		if err != nil {
			return err
		}

		links, err := client.GetLinks(topDuration, 1)
		if err != nil {
			return err
		}

		if len(links) == 0 {
			if topJSON {
				fmt.Println("[]")
				return nil
			}
			fmt.Println("No links found")
			return nil
		}

		if topJSON {
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
	rootCmd.AddCommand(topCmd)
	topCmd.Flags().StringVar(&topDuration, "duration", "7d", "Time window (e.g., 1d, 7d, 30d, all)")
	topCmd.Flags().BoolVar(&topJSON, "json", false, "Output as JSON")
}
