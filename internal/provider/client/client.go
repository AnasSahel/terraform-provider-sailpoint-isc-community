package client

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"resty.dev/v3"
)

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
		AddRetryConditions(retryCondition).
		AddRequestMiddleware(func(c *resty.Client, req *resty.Request) error {
			// Add authentication token to every request
			token, err := client.getToken(req.Context())
			if err != nil {
				return err
			}
			req.
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept", "application/json")

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
	if err := client.refreshToken(client.HTTPClient.Context()); err != nil {
		return nil, fmt.Errorf("Initial authentication failed: %w", err)
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
