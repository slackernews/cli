package api

import (
	"os"
	"testing"

	"github.com/slackernews/cli/pkg/config"
)

func TestNewClientNotConfigured(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	os.Unsetenv("SLACKERNEWS_TOKEN")

	_, err := NewClient(false)
	if err == nil {
		t.Fatal("expected error when not configured")
	}
	if !containsStr(err.Error(), "not configured") {
		t.Errorf("expected 'not configured' error, got: %v", err)
	}
}

func TestNewClientMissingToken(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	os.Unsetenv("SLACKERNEWS_TOKEN")

	// Save config with URL but no token
	cfg := config.Config{InstanceURL: "https://example.com"}
	config.Save(cfg)

	_, err := NewClient(false)
	if err == nil {
		t.Fatal("expected error when token is missing")
	}
	if !containsStr(err.Error(), "no API token found") && !containsStr(err.Error(), "token") {
		t.Errorf("expected token-related error, got: %v", err)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubStr(s, substr))
}

func containsSubStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
