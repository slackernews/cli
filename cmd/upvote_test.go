package cmd

import (
	"bytes"
	"net/http"
	"testing"
)

func TestUpvoteSuccess(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"upvote", "--insecure", "https://example.com"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Upvoted successfully")) {
		t.Errorf("expected output to contain confirmation, got: %s", buf.String())
	}
}

func TestUpvoteAlreadyUpvoted(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"upvote", "--insecure", "https://example.com"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("already upvoted")) {
		t.Errorf("expected error to contain 'already upvoted', got: %v", err)
	}
}

func TestUnvoteSuccess(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"unvote", "--insecure", "https://example.com"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Vote removed successfully")) {
		t.Errorf("expected output to contain confirmation, got: %s", buf.String())
	}
}

func TestUnvoteNotFound(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"unvote", "--insecure", "https://example.com"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("no vote to remove")) {
		t.Errorf("expected error to contain 'no vote to remove', got: %v", err)
	}
}

func TestCommentSuccess(t *testing.T) {
	t.Cleanup(resetCmd)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}
	setupTestEnv(t, handler)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"comment", "--insecure", "https://example.com", "Great read"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Comment posted successfully")) {
		t.Errorf("expected output to contain confirmation, got: %s", buf.String())
	}
}
