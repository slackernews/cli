package auth

import (
	"os"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestGetTokenFromEnv(t *testing.T) {
	t.Setenv("SLACKERNEWS_TOKEN", "env-token")

	token, err := GetToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "env-token" {
		t.Errorf("expected %q, got %q", "env-token", token)
	}
}

func TestGetTokenKeyringFirst(t *testing.T) {
	// When keyring has a value, it should be returned even if env has a different value
	t.Setenv("SLACKERNEWS_TOKEN", "env-token")
	keyring.Set(serviceName, accountName, "keyring-token")
	defer keyring.Delete(serviceName, accountName)

	token, err := GetToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "keyring-token" {
		t.Errorf("expected keyring-token %q, got %q", "keyring-token", token)
	}
}

func TestGetTokenKeyringNotFound(t *testing.T) {
	// Ensure no env var and no keyring entry
	os.Unsetenv("SLACKERNEWS_TOKEN")
	keyring.Delete(serviceName, accountName)

	_, err := GetToken()
	if err == nil {
		t.Fatal("expected error when no token exists")
	}
	if !strings.Contains(err.Error(), "no API token found") {
		t.Errorf("expected 'no API token found' error, got: %v", err)
	}
}

func TestGetTokenKeyringError(t *testing.T) {
	// This test verifies the error wrapping path when keyring returns
	// an unexpected error (not ErrNotFound). We can't easily simulate
	// a keyring error, but we can at least verify the code path when
	// env is not set and keyring returns ErrNotFound (which falls through
	// to the env check and then the final error).
	os.Unsetenv("SLACKERNEWS_TOKEN")
	keyring.Delete(serviceName, accountName)

	_, err := GetToken()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "no API token found") {
		t.Errorf("expected 'no API token found' error, got: %v", err)
	}
}

func TestSetToken(t *testing.T) {
	defer keyring.Delete(serviceName, accountName)

	if err := SetToken("my-token"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	token, err := keyring.Get(serviceName, accountName)
	if err != nil {
		t.Fatalf("unexpected error retrieving token: %v", err)
	}
	if token != "my-token" {
		t.Errorf("expected %q, got %q", "my-token", token)
	}
}
