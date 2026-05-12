package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/slackernews/cli/pkg/api"
	"github.com/slackernews/cli/pkg/config"
)

func setupTestEnv(t *testing.T, handler http.HandlerFunc) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	cfg := config.Config{InstanceURL: ts.URL}
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SLACKERNEWS_TOKEN", "test-token")
}

func resetCmd() {
	allowInsecure = false
	topJSON = false
	searchJSON = false
	topDuration = "7d"
	rootCmd.SetArgs([]string{})
	rootCmd.SetOut(nil)
	rootCmd.SetErr(nil)
}

func TestTopJSON(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.RenderableLink{
			{
				Link:         api.Link{URL: "https://example.com", Title: "Example"},
				FirstShare:   api.FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Alice"},
				DisplayScore: 5,
				ReplyCount:   2,
			},
		})
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"top", "--json", "--insecure"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var links []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &links); err != nil {
		t.Fatalf("failed to unmarshal output: %v\noutput: %s", err, buf.String())
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if links[0]["title"] != "Example" {
		t.Errorf("expected title 'Example', got %v", links[0]["title"])
	}
	if links[0]["url"] != "https://example.com" {
		t.Errorf("expected url 'https://example.com', got %v", links[0]["url"])
	}
}

func TestTopNoLinks(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.RenderableLink{})
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"top", "--insecure"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("No links found")) {
		t.Errorf("expected output to contain 'No links found', got: %s", buf.String())
	}
}
