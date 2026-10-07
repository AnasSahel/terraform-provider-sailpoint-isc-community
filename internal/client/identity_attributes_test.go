// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeIdentityAttributeServer mimics the ISC behaviour behind issue #179: each
// create is a non-atomic read-modify-write of one shared document, so
// overlapping creates are acknowledged with 2xx but only the last write survives.
type fakeIdentityAttributeServer struct {
	mu      sync.Mutex
	attrs   map[string]IdentityAttributeAPI
	dropAll bool // acknowledge creates but never persist them
}

func (f *fakeIdentityAttributeServer) handler(t *testing.T) http.Handler {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
	})
	mux.HandleFunc("POST /v2025/identity-attributes", func(w http.ResponseWriter, r *http.Request) {
		var attr IdentityAttributeAPI
		if err := json.NewDecoder(r.Body).Decode(&attr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Read the shared document...
		f.mu.Lock()
		snapshot := make(map[string]IdentityAttributeAPI, len(f.attrs))
		for k, v := range f.attrs {
			snapshot[k] = v
		}
		f.mu.Unlock()

		// ...let concurrent requests interleave...
		time.Sleep(20 * time.Millisecond)

		// ...and write it back, clobbering anything written meanwhile.
		if !f.dropAll {
			snapshot[attr.Name] = attr
			f.mu.Lock()
			f.attrs = snapshot
			f.mu.Unlock()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(attr)
	})
	mux.HandleFunc("GET /v2025/identity-attributes/{name}", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		attr, ok := f.attrs[r.PathValue("name")]
		f.mu.Unlock()
		if !ok {
			http.Error(w, `{"messages":[{"text":"NotFoundException"}]}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(attr)
	})
	return mux
}

func newTestClient(t *testing.T, f *fakeIdentityAttributeServer) *Client {
	t.Helper()
	srv := httptest.NewServer(f.handler(t))
	t.Cleanup(srv.Close)

	c, err := NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	c.identityAttributeReadBackDelays = []time.Duration{0, time.Millisecond, time.Millisecond}
	return c
}

// TestCreateIdentityAttribute_ParallelCreatesAllPersist covers the first
// acceptance criterion of issue #179: 10 creates fired at once all exist afterwards.
func TestCreateIdentityAttribute_ParallelCreatesAllPersist(t *testing.T) {
	f := &fakeIdentityAttributeServer{attrs: map[string]IdentityAttributeAPI{}}
	c := newTestClient(t, f)

	const n = 10
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("attr%d", i)
			got, err := c.CreateIdentityAttribute(context.Background(), &IdentityAttributeAPI{Name: name})
			if err != nil {
				errs <- err
				return
			}
			if got.Name != name {
				errs <- fmt.Errorf("got name %q, want %q", got.Name, name)
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("create failed: %v", err)
	}

	for i := range n {
		name := fmt.Sprintf("attr%d", i)
		if _, err := c.GetIdentityAttribute(context.Background(), name); err != nil {
			t.Errorf("attribute %q missing after parallel create: %v", name, err)
		}
	}
}

// TestCreateIdentityAttribute_FailsWhenReadBackNotFound covers the second
// acceptance criterion: a 2xx create that never becomes readable is an error.
func TestCreateIdentityAttribute_FailsWhenReadBackNotFound(t *testing.T) {
	f := &fakeIdentityAttributeServer{attrs: map[string]IdentityAttributeAPI{}, dropAll: true}
	c := newTestClient(t, f)

	_, err := c.CreateIdentityAttribute(context.Background(), &IdentityAttributeAPI{Name: "lost"})
	if err == nil {
		t.Fatal("expected an error when the created attribute is not found on read-back")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected error to wrap ErrNotFound, got: %v", err)
	}
	if !strings.Contains(err.Error(), "not found on read-back") {
		t.Errorf("expected an explicit read-back message, got: %v", err)
	}
}
