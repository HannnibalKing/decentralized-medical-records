package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// IPFSClient defines the minimal interface for fetching and pinning.
type IPFSClient interface {
	Get(ctx context.Context, cid string) ([]byte, error)
	Pin(ctx context.Context, cid string) error
}

// HTTPIPFSClient is a simple Kubo-compatible client using the HTTP API.
type HTTPIPFSClient struct {
	baseURL string
	http    *http.Client
}

// NewIPFS builds a client pointing at a Kubo API (e.g., http://localhost:5001).
func NewIPFS(baseURL string) *HTTPIPFSClient {
	return &HTTPIPFSClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Get fetches a CID via /api/v0/cat.
func (c *HTTPIPFSClient) Get(ctx context.Context, cid string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v0/cat?arg=%s", c.baseURL, cid)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipfs cat status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// Pin adds a CID to the pin set via /api/v0/pin/add.
func (c *HTTPIPFSClient) Pin(ctx context.Context, cid string) error {
	url := fmt.Sprintf("%s/api/v0/pin/add?arg=%s", c.baseURL, cid)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ipfs pin status %d", resp.StatusCode)
	}
	io.Copy(io.Discard, resp.Body)
	return nil
}

// Arweave-backed fallback can be used when IPFS fails.
