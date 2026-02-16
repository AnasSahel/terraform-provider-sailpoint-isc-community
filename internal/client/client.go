// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"resty.dev/v3"
)

// ErrNotFound is returned when a resource is not found (HTTP 404).
// Use errors.Is(err, client.ErrNotFound) to check for this error.
var ErrNotFound = errors.New("resource not found")

type Client struct {
	BaseURL      string
	HTTPClient   *resty.Client
	ClientID     string
	ClientSecret string

	token       string
	tokenExpiry time.Time
	tokenMutex  sync.RWMutex
}

func NewClient(baseURL, clientID, clientSecret string) (*Client, error) {
	client := &Client{
		BaseURL:      baseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Configure Resty HTTP client with retry logic
	client.HTTPClient = resty.New().
		SetBaseURL(baseURL).
		SetTimeout(30 * time.Second).
		// Retry configuration
		SetRetryCount(5).                      // Retry up to 5 times
		SetRetryWaitTime(1 * time.Second).     // Wait 1 second between retries
		SetRetryMaxWaitTime(30 * time.Second). // Wait up to 30 seconds between retries
		// SetAllowNonIdempotentRetry(true).      // Retry POST/PATCH too (v3 only retries idempotent methods by default)
		AddRetryConditions(retryCondition).
		AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
			// Add authentication token to every request
			token, err := client.getToken(req.Context())
			if err != nil {
				return err
			}
			req.
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				SetHeader("Accept", "application/json")
			if req.Header.Get("Content-Type") == "" {
				req.SetHeader("Content-Type", "application/json")
			}

			return nil
		}).
		AddResponseMiddleware(func(c *resty.Client, resp *resty.Response) error {
			// Log rate limit headers for debugging
			if remaining := resp.Header().Get("X-RateLimit-Remaining"); remaining != "" {
				tflog.Debug(resp.Request.Context(), "Rate limit remaining", map[string]any{
					"X-RateLimit-Remaining": remaining,
				})
			}
			return nil
		})

	// Initial authentication
	if err := client.refreshToken(context.Background()); err != nil {
		return nil, fmt.Errorf("initial authentication failed: %w", err)
	}

	return client, nil
}

func retryCondition(r *resty.Response, err error) bool {
	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on 5xx server errors
	if r.StatusCode() >= 500 && r.StatusCode() <= 599 {
		return true
	}

	// Retry on 429 rate limit (SailPoint has 100 req/10s limit)
	if r.StatusCode() == http.StatusTooManyRequests {
		return true
	}

	// Retry on 408 request timeout
	if r.StatusCode() == http.StatusRequestTimeout {
		return true
	}

	return false
}

func (c *Client) doRequest(ctx context.Context, method string, url string, body any, result any) (*resty.Response, error) {
	return c.doRequestWithHeaders(ctx, method, url, body, result, nil)
}

func (c *Client) doRequestWithHeaders(ctx context.Context, method string, url string, body any, result any, headers map[string]string) (*resty.Response, error) {
	req := c.prepareRequest(ctx)

	// Add custom headers if provided
	for key, value := range headers {
		req.SetHeader(key, value)
	}

	if body != nil {
		req.SetBody(body)
	}

	if result != nil {
		req.SetResult(result)
	}

	switch method {
	case http.MethodGet:
		return req.Get(url)
	case http.MethodPost:
		return req.Post(url)
	case http.MethodPut:
		return req.Put(url)
	case http.MethodPatch:
		req.SetHeader("Content-Type", "application/json-patch+json")
		return req.Patch(url)
	case http.MethodDelete:
		return req.Delete(url)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}
}
