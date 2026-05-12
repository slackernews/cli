package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestClient(handler http.HandlerFunc) (*Client, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	return &Client{
		baseURL:    ts.URL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      "test-token",
		insecure:   true,
	}, ts
}

func TestGetLinks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/cli/links" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("duration") != "7d" {
			t.Errorf("unexpected duration: %s", r.URL.Query().Get("duration"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected authorization header: %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]RenderableLink{
			{
				Link:         Link{URL: "https://example.com", Title: "Example"},
				FirstShare:   FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Alice"},
				DisplayScore: 5,
				IsUpvoted:    false,
				ReplyCount:   2,
			},
		})
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	links, err := client.GetLinks("7d", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if links[0].Link.Title != "Example" {
		t.Errorf("unexpected title: %s", links[0].Link.Title)
	}
}

func TestSearchLinks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/cli/links/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "kubernetes" {
			t.Errorf("unexpected query: %s", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]RenderableLink{})
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	links, err := client.SearchLinks("kubernetes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 0 {
		t.Fatalf("expected 0 links, got %d", len(links))
	}
}

func TestUpvote(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/api/v1/cli/links/https://example.com/upvote" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	if err := client.Upvote("https://example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpvoteAlreadyUpvoted(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	err := client.Upvote("https://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "already upvoted" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestUnvote(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	if err := client.Unvote("https://example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnvoteNotFound(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	err := client.Unvote("https://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "no vote to remove" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestComment(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/api/v1/cli/links/https://example.com/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		if body["body"] != "Great read" {
			t.Errorf("unexpected body: %s", body["body"])
		}
		w.WriteHeader(http.StatusCreated)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	if err := client.Comment("https://example.com", "Great read"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	_, err := client.GetLinks("7d", 1)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "authentication failed: check your API token" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestServerError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	client, ts := newTestClient(handler)
	defer ts.Close()

	_, err := client.GetLinks("7d", 1)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "server error: 500 Internal Server Error" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNetworkError(t *testing.T) {
	client := &Client{
		baseURL:    "http://localhost:1",
		httpClient: &http.Client{Timeout: 100 * time.Millisecond},
		token:      "test-token",
		insecure:   true,
	}

	_, err := client.GetLinks("7d", 1)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "connection refused") && !strings.Contains(err.Error(), "server unreachable") {
		t.Errorf("unexpected error message: %v", err)
	}
}
