package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/slackernews/cli/pkg/auth"
	"github.com/slackernews/cli/pkg/config"
)

func createClient(insecure bool) (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.InstanceURL == "" {
		return nil, fmt.Errorf("not configured: run 'slackernews configure --url <url>'")
	}

	token, err := auth.GetToken()
	if err != nil {
		return nil, err
	}

	return api.NewClient(cfg.InstanceURL, token, insecure)
}
