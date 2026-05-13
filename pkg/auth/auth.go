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

// GetToken retrieves the API token from the environment or OS keychain.
// Environment variable is checked first so CI/headless environments can
// override the keychain without needing a secrets service.
func GetToken() (string, error) {
	if token := os.Getenv("SLACKERNEWS_TOKEN"); token != "" {
		return token, nil
	}

	token, err := keyring.Get(serviceName, accountName)
	if err == nil {
		return token, nil
	}

	if err != keyring.ErrNotFound {
		return "", fmt.Errorf("failed to retrieve token from keychain: %w", err)
	}

	return "", fmt.Errorf("no API token found: run 'slackernews configure --token <token>' or set SLACKERNEWS_TOKEN")
}

// SetToken stores the API token in the OS keychain.
func SetToken(token string) error {
	return keyring.Set(serviceName, accountName, token)
}
