package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ArweaveClient defines minimal interface to fetch by transaction id.
type ArweaveClient interface {
	Get(ctx context.Context, tx string) ([]byte, error)
}

// HTTPArweaveClient fetches from a public Arweave gateway.
type HTTPArweaveClient struct {
	baseURL string
	http    *http.Client
}

// NewArweave builds a client pointing at a gateway (e.g., https://arweave.net).
func NewArweave(baseURL string) *HTTPArweaveClient {
	return &HTTPArweaveClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Get fetches a transaction payload.
func (c *HTTPArweaveClient) Get(ctx context.Context, tx string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, tx)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arweave status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
