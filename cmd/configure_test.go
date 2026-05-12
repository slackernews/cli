package cmd

import (
	"bytes"
	"testing"

	"github.com/slackernews/cli/pkg/config"
)

func TestConfigureMissingURL(t *testing.T) {
	t.Cleanup(resetCmd)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"configure"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --url is missing")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("--url is required")) {
		t.Errorf("expected '--url is required' error, got: %v", err)
	}
}

func TestConfigureInvalidURL(t *testing.T) {
	t.Cleanup(resetCmd)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"configure", "--url", "://invalid"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("invalid URL")) {
		t.Errorf("expected 'invalid URL' error, got: %v", err)
	}
}

func TestConfigureHTTPWithoutInsecure(t *testing.T) {
	t.Cleanup(resetCmd)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"configure", "--url", "http://example.com"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for HTTP URL without --insecure")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("https://")) {
		t.Errorf("expected HTTPS error, got: %v", err)
	}
}

func TestConfigureHTTPS(t *testing.T) {
	t.Cleanup(resetCmd)

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"configure", "--url", "https://slackernews.example.com"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Configuration saved")) {
		t.Errorf("expected 'Configuration saved' message, got: %s", buf.String())
	}

	// Verify config was saved
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.InstanceURL != "https://slackernews.example.com" {
		t.Errorf("expected URL %q, got %q", "https://slackernews.example.com", cfg.InstanceURL)
	}
}

func TestConfigureTrimsTrailingSlash(t *testing.T) {
	t.Cleanup(resetCmd)

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	rootCmd.SetArgs([]string{"configure", "--url", "https://slackernews.example.com/"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.InstanceURL != "https://slackernews.example.com" {
		t.Errorf("expected URL without trailing slash, got %q", cfg.InstanceURL)
	}
}

func TestConfigureInsecureHTTP(t *testing.T) {
	t.Cleanup(resetCmd)

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	rootCmd.SetArgs([]string{"configure", "--insecure", "--url", "http://example.com"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.InstanceURL != "http://example.com" {
		t.Errorf("expected URL %q, got %q", "http://example.com", cfg.InstanceURL)
	}
}
