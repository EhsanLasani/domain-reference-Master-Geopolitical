package httpx

import (
	"context"
	"net/http"
	"time"
)

// Client implements guideline 12 - Shared HTTP client with retries
type Client struct {
	httpClient *http.Client
	retries    int
	timeout    time.Duration
}

func NewClient(timeout time.Duration, retries int) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		retries: retries,
		timeout: timeout,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	
	var resp *http.Response
	var err error
	
	for i := 0; i <= c.retries; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		
		if i < c.retries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(i+1) * time.Second):
				// Exponential backoff
			}
		}
	}
	
	return resp, err
}