package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"
)

// Error types for distinct exit codes and programmatic handling.
type (
	AuthError     struct{ Message string }
	NetworkError  struct{ Message string }
	ServerError   struct{ Message string }
	RateLimitError struct {
		Message    string
		RetryAfter time.Duration
	}
)

func (e *AuthError) Error() string     { return e.Message }
func (e *NetworkError) Error() string  { return e.Message }
func (e *ServerError) Error() string   { return e.Message }
func (e *RateLimitError) Error() string { return e.Message }

// Link represents a SlackerNews link.
type Link struct {
	URL      string `json:"url"`
	Domain   string `json:"domain"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	IsHidden bool   `json:"isHidden"`
}

// FirstShare represents the first share metadata.
type FirstShare struct {
	SharedAt    int64  `json:"sharedAt"`
	MessageTs   string `json:"messageTs"`
	UserID      string `json:"userId"`
	FullName    string `json:"fullName"`
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	Permalink   string `json:"permalink"`
	ReplyCount  int    `json:"replyCount"`
}

// RenderableLink is the full link object returned by the API.
type RenderableLink struct {
	Link         Link       `json:"link"`
	FirstShare   FirstShare `json:"firstShare"`
	DisplayScore float64    `json:"displayScore"`
	IsUpvoted    bool       `json:"isUpvoted"`
	ReplyCount   int        `json:"replyCount"`
}

// Client communicates with the SlackerNews API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	insecure   bool
}

// NewClient creates an API client with the given parameters.
func NewClient(baseURL, token string, insecure bool) (*Client, error) {
	if baseURL == "" {
		return nil, &NetworkError{Message: "not configured: run 'slackernews configure --url <url>'"}
	}

	if token == "" {
		return nil, &AuthError{Message: "no API token found: run 'slackernews configure --token <token>' or set SLACKERNEWS_TOKEN"}
	}

	timeout := 30 * time.Second
	if d := os.Getenv("SLACKERNEWS_TIMEOUT"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil && parsed > 0 {
			timeout = parsed
		}
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		token:    token,
		insecure: insecure,
	}, nil
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	base := strings.TrimSuffix(c.baseURL, "/")
	u, err := url.Parse(base + path)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "https" && !c.insecure {
		return nil, &NetworkError{Message: fmt.Sprintf("insecure URL detected: %s (use --insecure)", c.baseURL)}
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			return nil, &NetworkError{Message: fmt.Sprintf("connection timed out: %v", err)}
		}
		if opErr, ok := err.(*net.OpError); ok {
			if dnsErr, ok := opErr.Err.(*net.DNSError); ok {
				return nil, &NetworkError{Message: fmt.Sprintf("DNS lookup failed (%s): %v", dnsErr.Name, err)}
			}
		}
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, &NetworkError{Message: fmt.Sprintf("connection refused: %v", err)}
		}
		return nil, &NetworkError{Message: fmt.Sprintf("server unreachable: %v", err)}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		retryAfter := 0
		if h := resp.Header.Get("Retry-After"); h != "" {
			// Try parsing as seconds first
			if d, err := time.ParseDuration(h + "s"); err == nil {
				retryAfter = int(d.Seconds())
			}
		}
		return nil, &RateLimitError{
			Message:    fmt.Sprintf("rate limited: retry after %ds", retryAfter),
			RetryAfter: time.Duration(retryAfter) * time.Second,
		}
	}

	if resp.StatusCode >= 500 {
		resp.Body.Close()
		return nil, &ServerError{Message: fmt.Sprintf("server error: %s", resp.Status)}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		return nil, &AuthError{Message: "authentication failed: check your API token"}
	}

	return resp, nil
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.withRetry(ctx, "GET", path, nil)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.withRetry(ctx, "POST", path, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.withRetry(ctx, "DELETE", path, nil)
}

func (c *Client) withRetry(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

	for attempt := 0; attempt <= len(backoffs); attempt++ {
		resp, err = c.request(ctx, method, path, body)
		if err == nil {
			return resp, nil
		}

		// Don't retry auth errors, client errors, or context cancellation
		if errors.As(err, new(*AuthError)) {
			return nil, err
		}
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, err
		}

		// Check if it's a retriable error
		var netErr *NetworkError
		var srvErr *ServerError
		var rateErr *RateLimitError
		isRetriable := errors.As(err, &netErr) || errors.As(err, &srvErr) || errors.As(err, &rateErr)
		if !isRetriable {
			return nil, err
		}

		// On last attempt, return the error
		if attempt == len(backoffs) {
			return nil, err
		}

		// Wait before retry, but respect context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoffs[attempt]):
			// continue to next attempt
		}
	}

	return nil, err
}

// DecodeJSON reads and decodes a JSON response body.
func DecodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response from %s %s: %w", resp.Request.Method, resp.Request.URL.Path, err)
	}
	return nil
}

func decodeLinks(resp *http.Response) ([]RenderableLink, error) {
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var links []RenderableLink
	if err := DecodeJSON(resp, &links); err != nil {
		return nil, err
	}
	return links, nil
}

// GetLinks fetches top links for a duration and page.
func (c *Client) GetLinks(ctx context.Context, duration string, page int) ([]RenderableLink, error) {
	path := fmt.Sprintf("/api/v1/cli/links?duration=%s&page=%d", url.QueryEscape(duration), page)
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	return decodeLinks(resp)
}

// SearchLinks searches links by query.
func (c *Client) SearchLinks(ctx context.Context, query string) ([]RenderableLink, error) {
	path := fmt.Sprintf("/api/v1/cli/links/search?q=%s", url.QueryEscape(query))
	resp, err := c.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	return decodeLinks(resp)
}

// Upvote upvotes a link by its URL.
func (c *Client) Upvote(ctx context.Context, linkID string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/upvote", url.PathEscape(linkID))
	resp, err := c.Post(ctx, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("already upvoted")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// Unvote removes an upvote from a link by its URL.
func (c *Client) Unvote(ctx context.Context, linkID string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/upvote", url.PathEscape(linkID))
	resp, err := c.Delete(ctx, path)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no vote to remove")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// Comment posts a comment on a link by its URL.
func (c *Client) Comment(ctx context.Context, linkID string, body string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/comments", url.PathEscape(linkID))
	resp, err := c.Post(ctx, path, map[string]string{"body": body})
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("link not found")
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
