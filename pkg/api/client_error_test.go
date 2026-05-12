package api

import (
	"errors"
	"testing"
)

func TestNewClientEmptyBaseURL(t *testing.T) {
	_, err := NewClient("", "token", false)
	if err == nil {
		t.Fatal("expected error when baseURL is empty")
	}
	var netErr *NetworkError
	if !errors.As(err, &netErr) {
		t.Errorf("expected NetworkError, got: %T %v", err, err)
	}
	if err.Error() != "not configured: run 'slackernews configure --url <url>'" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewClientEmptyToken(t *testing.T) {
	_, err := NewClient("https://example.com", "", false)
	if err == nil {
		t.Fatal("expected error when token is empty")
	}
	var authErr *AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got: %T %v", err, err)
	}
	if err.Error() != "no API token found: run 'slackernews configure --token <token>' or set SLACKERNEWS_TOKEN" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewClientSuccess(t *testing.T) {
	client, err := NewClient("https://example.com", "token", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.baseURL != "https://example.com" {
		t.Errorf("expected baseURL 'https://example.com', got %q", client.baseURL)
	}
	if client.token != "token" {
		t.Errorf("expected token 'token', got %q", client.token)
	}
	if client.insecure != false {
		t.Errorf("expected insecure false, got %v", client.insecure)
	}
}
