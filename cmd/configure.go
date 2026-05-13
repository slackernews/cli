package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/slackernews/cli/pkg/auth"
	"github.com/slackernews/cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	configureURL   string
	configureToken string
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the SlackerNews CLI",
	Long:  `Set the instance URL and API token for the CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configureURL == "" {
			return fmt.Errorf("--url is required")
		}

		u, err := url.Parse(configureURL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		if u.Scheme != "https" && !allowInsecure {
			return fmt.Errorf("URL must use https:// (use --insecure to allow http:// for development)")
		}

		cfg := config.Config{
			InstanceURL: strings.TrimSuffix(configureURL, "/"),
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if configureToken != "" {
			if err := auth.SetToken(configureToken); err != nil {
				return fmt.Errorf("failed to store token: %w", err)
			}
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Configuration saved successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringVar(&configureURL, "url", "", "SlackerNews instance URL")
	configureCmd.Flags().StringVar(&configureToken, "token", "", "API token")
}
