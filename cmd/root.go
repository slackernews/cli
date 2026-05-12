package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var allowInsecure bool

var rootCmd = &cobra.Command{
	Use:   "slackernews",
	Short: "SlackerNews CLI for terminal-based link browsing",
	Long: `Browse top links, vote, and comment on SlackerNews
directly from your terminal.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&allowInsecure, "insecure", false, "Allow non-HTTPS URLs (development only)")
}
