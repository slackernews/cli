package auth

import (
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "slackernews"
	accountName = "api-token"
)

// GetToken retrieves the API token from the OS keychain or environment.
func GetToken() (string, error) {
	// Fallback to environment variable for CI/headless environments
	if token := os.Getenv("SLACKERNEWS_TOKEN"); token != "" {
		return token, nil
	}

	token, err := keyring.Get(serviceName, accountName)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("no API token found: run 'slackernews configure --token <token>' or set SLACKERNEWS_TOKEN")
		}
		return "", fmt.Errorf("failed to retrieve token from keychain: %w", err)
	}
	return token, nil
}

// SetToken stores the API token in the OS keychain.
func SetToken(token string) error {
	return keyring.Set(serviceName, accountName, token)
}
