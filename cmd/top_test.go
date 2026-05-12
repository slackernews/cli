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
	globalInsecure = false
	topJSON = false
	searchJSON = false
	rootCmd.SetArgs(nil)
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

	if !bytes.Contains(buf.Bytes(), []byte(`"title": "Example"`)) {
		t.Errorf("expected output to contain title, got: %s", buf.String())
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
