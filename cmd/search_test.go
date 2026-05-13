package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/slackernews/cli/pkg/api"
)

func TestSearchWithResults(t *testing.T) {
	t.Cleanup(resetCmd)

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]api.RenderableLink{
				{
					Link:         api.Link{URL: "https://kubernetes.io", Title: "Kubernetes"},
					FirstShare:   api.FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Bob"},
					DisplayScore: 10,
					ReplyCount:   3,
				},
			}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
		setupTestEnv(t, handler)

		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"search", "--insecure", "kubernetes"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bytes.Contains(buf.Bytes(), []byte("Kubernetes")) {
			t.Errorf("expected output to contain 'Kubernetes', got: %s", buf.String())
		}
	}

	func TestSearchNoResults(t *testing.T) {
		t.Cleanup(resetCmd)

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]api.RenderableLink{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
		setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"search", "--insecure", "nonexistent"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("No links found")) {
		t.Errorf("expected output to contain 'No links found', got: %s", buf.String())
	}
}

func TestSearchJSON(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]api.RenderableLink{
				{
					Link:         api.Link{URL: "https://kubernetes.io", Title: "Kubernetes"},
					FirstShare:   api.FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Bob"},
					DisplayScore: 10,
					ReplyCount:   3,
				},
			}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
		setupTestEnv(t, handler)

		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		rootCmd.SetArgs([]string{"search", "--json", "--insecure", "kubernetes"})

		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var links []map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &links); err != nil {
			t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
		}
		if len(links) != 1 {
			t.Fatalf("expected 1 link, got %d", len(links))
		}
		if links[0]["title"] != "Kubernetes" {
			t.Errorf("expected title 'Kubernetes', got %v", links[0]["title"])
		}
	}

	func TestSearchJSONEmpty(t *testing.T) {
		t.Cleanup(resetCmd)

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]api.RenderableLink{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
		setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"search", "--json", "--insecure", "nonexistent"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.String() != "[]\n" {
		t.Errorf("expected '[]\\n', got: %q", buf.String())
	}
}
