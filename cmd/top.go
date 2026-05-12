package cmd

import (
	"github.com/slackernews/cli/pkg/api"
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
		client, err := api.NewClient(allowInsecure)
		if err != nil {
			return err
		}

		links, err := client.GetLinks(topDuration, 1)
		if err != nil {
			return err
		}

		return printLinks(cmd, links, topJSON)
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.Flags().StringVar(&topDuration, "duration", "7d", "Time window (e.g., 1d, 7d, 30d, all)")
	topCmd.Flags().BoolVar(&topJSON, "json", false, "Output as JSON")
}
