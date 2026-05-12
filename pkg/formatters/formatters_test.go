package formatters

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/slackernews/cli/pkg/api"
)

func TestFormatJSON(t *testing.T) {
	links := []api.RenderableLink{
		{
			Link:         api.Link{URL: "https://example.com", Title: "Example"},
			FirstShare:   api.FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Alice"},
			DisplayScore: 5,
			IsUpvoted:    false,
			ReplyCount:   2,
		},
	}

	out, err := FormatJSON(links)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, string(out))
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 link, got %d", len(result))
	}
	if result[0]["title"] != "Example" {
		t.Errorf("expected title 'Example', got %v", result[0]["title"])
	}
	if result[0]["url"] != "https://example.com" {
		t.Errorf("expected url 'https://example.com', got %v", result[0]["url"])
	}
	if result[0]["score"] != float64(5) {
		t.Errorf("expected score 5, got %v", result[0]["score"])
	}
	if result[0]["replyCount"] != float64(2) {
		t.Errorf("expected replyCount 2, got %v", result[0]["replyCount"])
	}
	if result[0]["isUpvoted"] != false {
		t.Errorf("expected isUpvoted false, got %v", result[0]["isUpvoted"])
	}
	if result[0]["firstSharedBy"] != "Alice" {
		t.Errorf("expected firstSharedBy 'Alice', got %v", result[0]["firstSharedBy"])
	}
}

func TestFormatJSONEmptyTitle(t *testing.T) {
	links := []api.RenderableLink{
		{
			Link:         api.Link{URL: "https://example.com", Title: ""},
			FirstShare:   api.FirstShare{SharedAt: time.Now().UnixMilli(), FullName: "Bob"},
			DisplayScore: 1,
			ReplyCount:   0,
		},
	}

	out, err := FormatJSON(links)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	json.Unmarshal(out, &result)

	if result[0]["title"] != "https://example.com" {
		t.Errorf("expected title fallback to URL, got %v", result[0]["title"])
	}
}

func TestFormatJSONEmptyLinks(t *testing.T) {
	out, err := FormatJSON([]api.RenderableLink{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "[]" {
		t.Errorf("expected '[]', got %s", string(out))
	}
}

func TestFormatTable(t *testing.T) {
	links := []api.RenderableLink{
		{
			Link:         api.Link{URL: "https://example.com", Title: "Example"},
			FirstShare:   api.FirstShare{SharedAt: time.Now().Add(-1 * time.Hour).UnixMilli(), FullName: "Alice"},
			DisplayScore: 5,
			IsUpvoted:    false,
			ReplyCount:   2,
		},
	}

	buf := new(bytes.Buffer)
	FormatTable(buf, links)
	output := buf.String()

	if !strings.Contains(output, "Example") {
		t.Errorf("expected output to contain 'Example', got: %s", output)
	}
	if !strings.Contains(output, "5") {
		t.Errorf("expected output to contain score '5', got: %s", output)
	}
	if !strings.Contains(output, "2") {
		t.Errorf("expected output to contain reply count '2', got: %s", output)
	}
}

func TestFormatTableEmpty(t *testing.T) {
	buf := new(bytes.Buffer)
	FormatTable(buf, []api.RenderableLink{})
	output := buf.String()

	if !strings.Contains(output, "No links found") {
		t.Errorf("expected output to contain 'No links found', got: %s", output)
	}
}

func TestFormatAgeMinutes(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(30*time.Minute/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "m") {
		t.Errorf("expected minutes suffix, got %q", result)
	}
}

func TestFormatAgeHours(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(5*time.Hour/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "h") {
		t.Errorf("expected hours suffix, got %q", result)
	}
}

func TestFormatAgeDays(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(3*24*time.Hour/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "d") {
		t.Errorf("expected days suffix, got %q", result)
	}
}

func TestFormatAgeMonths(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(90*24*time.Hour/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "mo") {
		t.Errorf("expected months suffix, got %q", result)
	}
}

func TestFormatAgeFuture(t *testing.T) {
	now := time.Now().UnixMilli()
	future := now + int64(5*time.Minute/time.Millisecond)
	result := formatAge(future)
	if !strings.HasSuffix(result, "m") {
		t.Errorf("expected positive minutes for future timestamp, got %q", result)
	}
	if strings.HasPrefix(result, "-") {
		t.Errorf("expected no negative sign for future timestamp, got %q", result)
	}
}

func TestFormatAgeExactlyOneHour(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(1*time.Hour/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "h") {
		t.Errorf("expected hours suffix at exactly 1 hour, got %q", result)
	}
}

func TestFormatAgeExactly24Hours(t *testing.T) {
	now := time.Now().UnixMilli()
	past := now - int64(24*time.Hour/time.Millisecond)
	result := formatAge(past)
	if !strings.HasSuffix(result, "d") {
		t.Errorf("expected days suffix at exactly 24 hours, got %q", result)
	}
}
