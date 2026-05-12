package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/slackernews/cli/pkg/auth"
	"github.com/slackernews/cli/pkg/config"
)

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

// NewClient creates an API client using saved configuration.
func NewClient(insecure bool) (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.InstanceURL == "" {
		return nil, fmt.Errorf("not configured: run 'slackernews configure --url <url>'")
	}

	token, err := auth.GetToken()
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:    cfg.InstanceURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		insecure:   insecure,
	}, nil
}

func (c *Client) request(method, path string, body interface{}) (*http.Response, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "https" && !c.insecure {
		return nil, fmt.Errorf("insecure URL detected: %s (use --insecure)", c.baseURL)
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("server unreachable: %w", err)
	}

	if resp.StatusCode >= 500 {
		resp.Body.Close()
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		return nil, fmt.Errorf("authentication failed: check your API token")
	}

	return resp, nil
}

// Get performs a GET request.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.request("GET", path, nil)
}

// Post performs a POST request.
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.request("POST", path, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.request("DELETE", path, nil)
}

// DecodeJSON reads and decodes a JSON response body.
func DecodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// GetLinks fetches top links for a duration and page.
func (c *Client) GetLinks(duration string, page int) ([]RenderableLink, error) {
	path := fmt.Sprintf("/api/v1/cli/links?duration=%s&page=%d", url.QueryEscape(duration), page)
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}
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

// SearchLinks searches links by query.
func (c *Client) SearchLinks(query string) ([]RenderableLink, error) {
	path := fmt.Sprintf("/api/v1/cli/links/search?q=%s", url.QueryEscape(query))
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}
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

// Upvote upvotes a link by its URL.
func (c *Client) Upvote(linkID string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/upvote", url.QueryEscape(linkID))
	resp, err := c.Post(path, nil)
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
func (c *Client) Unvote(linkID string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/upvote", url.QueryEscape(linkID))
	resp, err := c.Delete(path)
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
func (c *Client) Comment(linkID string, body string) error {
	path := fmt.Sprintf("/api/v1/cli/links/%s/comments", url.QueryEscape(linkID))
	resp, err := c.Post(path, map[string]string{"body": body})
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
