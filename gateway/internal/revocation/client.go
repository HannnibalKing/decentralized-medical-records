package revocation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client performs revocation lookups against an RPC/HTTP endpoint.
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient builds a revocation client pointing at baseURL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Lookup checks if a handle is revoked. The endpoint is expected to return JSON {"revoked":bool}.
func (c *Client) Lookup(ctx context.Context, handle string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/revocations/%s", c.baseURL, handle), nil)
	if err != nil {
		return false, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("revocation lookup status %d", resp.StatusCode)
	}
	var body struct {
		Revoked bool `json:"revoked"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, err
	}
	return body.Revoked, nil
}

// Revoke marks a handle revoked. Endpoint expected to accept POST {handle, reason}.
func (c *Client) Revoke(ctx context.Context, handle, reason string) error {
	payload := map[string]string{"handle": handle, "reason": reason}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/revocations", c.baseURL), strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("revoke status %d", resp.StatusCode)
	}
	return nil
}
