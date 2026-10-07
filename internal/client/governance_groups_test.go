// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// newGovernanceGroupTestClient starts a fake tenant that answers the OAuth token
// endpoint and delegates every other request to handler.
func newGovernanceGroupTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		}
		handler(w, r)
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

func TestGovernanceGroupEndpointsUseWorkgroupsV1(t *testing.T) {
	var seen []string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		switch {
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/workgroups/v1":
			writeJSON(t, w, http.StatusOK, []GovernanceGroupAPI{{ID: "wg1", Name: "n"}})
		default:
			writeJSON(t, w, http.StatusOK, GovernanceGroupAPI{ID: "wg1", Name: "n"})
		}
	})
	ctx := context.Background()

	if _, err := c.ListGovernanceGroups(ctx, ""); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := c.CreateGovernanceGroup(ctx, &GovernanceGroupAPI{Name: "n"}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := c.GetGovernanceGroup(ctx, "wg1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := c.PatchGovernanceGroup(ctx, "wg1", []JSONPatchOperation{NewReplacePatch("/name", "m")}); err != nil {
		t.Fatalf("patch: %v", err)
	}
	if err := c.DeleteGovernanceGroup(ctx, "wg1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	want := []string{
		"GET /workgroups/v1",
		"POST /workgroups/v1",
		"GET /workgroups/v1/wg1",
		"PATCH /workgroups/v1/wg1",
		"DELETE /workgroups/v1/wg1",
	}
	if strings.Join(seen, "\n") != strings.Join(want, "\n") {
		t.Errorf("requests mismatch\n got: %v\nwant: %v", seen, want)
	}
}

func TestListGovernanceGroupMembersPaginates(t *testing.T) {
	const total = governanceGroupMembersPageSize + 7
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workgroups/v1/wg1/members" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		page := []GovernanceGroupMemberAPI{}
		for i := offset; i < total && i < offset+limit; i++ {
			page = append(page, GovernanceGroupMemberAPI{Type: "IDENTITY", ID: fmt.Sprintf("id%d", i)})
		}
		writeJSON(t, w, http.StatusOK, page)
	})

	members, err := c.ListGovernanceGroupMembers(context.Background(), "wg1")
	if err != nil {
		t.Fatalf("list members: %v", err)
	}
	if len(members) != total {
		t.Errorf("got %d members, want %d", len(members), total)
	}
}

func TestAddGovernanceGroupMembers(t *testing.T) {
	tests := []struct {
		name    string
		results []governanceGroupMemberBulkAddResult
		wantErr string
	}{
		{
			name: "added and already member",
			results: []governanceGroupMemberBulkAddResult{
				{ID: "a", Status: http.StatusCreated},
				{ID: "b", Status: http.StatusConflict},
			},
		},
		{
			name: "per-identity failure",
			results: []governanceGroupMemberBulkAddResult{
				{ID: "a", Status: http.StatusCreated},
				{ID: "c", Status: http.StatusNotFound, Description: "Identity not found."},
			},
			wantErr: "identity c: status 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/workgroups/v1/wg1/members/bulk-add" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				writeJSON(t, w, http.StatusMultiStatus, tt.results)
			})

			err := c.AddGovernanceGroupMembers(context.Background(), "wg1", []GovernanceGroupMemberAPI{
				{Type: "IDENTITY", ID: "a"},
			})
			switch {
			case tt.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)):
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}
