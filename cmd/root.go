package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/slackernews/cli/pkg/api"
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
		os.Exit(exitCodeFor(err))
	}
}

func exitCodeFor(err error) int {
	var authErr *api.AuthError
	if errors.As(err, &authErr) {
		return 2
	}

	var netErr *api.NetworkError
	if errors.As(err, &netErr) {
		return 3
	}

	var srvErr *api.ServerError
	if errors.As(err, &srvErr) {
		return 4
	}

	var rateErr *api.RateLimitError
	if errors.As(err, &rateErr) {
		return 5
	}

	return 1
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&allowInsecure, "insecure", false, "Allow non-HTTPS URLs (development only)")
}
