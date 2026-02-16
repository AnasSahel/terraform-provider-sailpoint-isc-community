// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"time"

	"resty.dev/v3"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (c *Client) refreshToken(ctx context.Context) error {
	var tokenResp tokenResponse

	resp, err := resty.New().R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     c.ClientID,
			"client_secret": c.ClientSecret,
		}).
		SetResult(&tokenResp).
		Post(fmt.Sprintf("%s/oauth/token", c.BaseURL))

	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("token request returned %s: %s", resp.Status(), resp.String())
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	c.token = tokenResp.AccessToken
	// Refresh 5 minutes before actual expiry for safety
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	return nil
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	// Fast path: check if token is still valid (read lock only)
	c.tokenMutex.RLock()
	if time.Now().Before(c.tokenExpiry) {
		token := c.token
		c.tokenMutex.RUnlock()
		return token, nil
	}
	c.tokenMutex.RUnlock()

	// Slow path: need to refresh token (write lock)
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Double-check after acquiring write lock (another goroutine may have refreshed)
	if time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	// Actually refresh the token
	if err := c.refreshToken(ctx); err != nil {
		return "", err
	}

	return c.token, nil
}
