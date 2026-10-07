// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newTokenServer serves /oauth/token with a new token on every call and
// /v2025/ping as an authenticated endpoint that echoes the bearer token.
func newTokenServer(t *testing.T, tokenStatus int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(tokenStatus)
		if tokenStatus != http.StatusOK {
			_, _ = fmt.Fprint(w, `{"error":"invalid_client"}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"access_token":"token-%d","expires_in":3600,"token_type":"bearer"}`, n)
	})
	mux.HandleFunc("/v2025/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"auth":%q}`, r.Header.Get("Authorization"))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, &calls
}

func TestGetToken_RefreshesExpiredTokenWithoutDeadlock(t *testing.T) {
	srv, calls := newTokenServer(t, http.StatusOK)

	c, err := NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Force the slow path: the cached token has expired.
	c.tokenMutex.Lock()
	c.tokenExpiry = time.Now().Add(-time.Minute)
	c.tokenMutex.Unlock()

	type result struct {
		body string
		err  error
	}
	done := make(chan result, 1)
	go func() {
		var out struct {
			Auth string `json:"auth"`
		}
		_, err := c.prepareRequest(context.Background()).SetResult(&out).Get("/v2025/ping")
		done <- result{out.Auth, err}
	}()

	select {
	case res := <-done:
		if res.err != nil {
			t.Fatalf("request failed: %v", res.err)
		}
		if res.body != "Bearer token-2" {
			t.Errorf("Authorization = %q, want %q", res.body, "Bearer token-2")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request did not complete: token refresh deadlocked")
	}

	if got := calls.Load(); got != 2 {
		t.Errorf("token endpoint called %d times, want 2", got)
	}
}

func TestNewClient_TokenErrorStatus(t *testing.T) {
	srv, _ := newTokenServer(t, http.StatusUnauthorized)

	if _, err := NewClient(srv.URL, "id", "bad"); err == nil {
		t.Fatal("NewClient succeeded, want an error for a 401 token response")
	}
}
