// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"
)

func TestGetManagedClientUsesV1AndDropsSecret(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/managed-clients/v1/mc1" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"id":            "mc1",
			"name":          "VA-1",
			"clusterId":     "cl1",
			"status":        "NORMAL",
			"secret":        "top-secret",
			"alertKey":      "CLIENT_STATUS_NOT_CONFIGURED",
			"vaDownloadUrl": "https://example.invalid/va",
		})
	})

	mc, err := c.GetManagedClient(context.Background(), "mc1")
	if err != nil {
		t.Fatalf("GetManagedClient: %v", err)
	}
	if mc.ID != "mc1" || mc.ClusterID != "cl1" || mc.Name == nil || *mc.Name != "VA-1" {
		t.Errorf("unexpected managed client: %+v", mc)
	}
}

func TestGetManagedClientNotFound(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]any{})
	})

	_, err := c.GetManagedClient(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestListManagedClientsFollowsPagination(t *testing.T) {
	var offsets []string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/managed-clients/v1" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		offset := r.URL.Query().Get("offset")
		offsets = append(offsets, offset)
		if got := r.URL.Query().Get("limit"); got != strconv.Itoa(managedClientsPageSize) {
			t.Errorf("limit = %q", got)
		}
		size := managedClientsPageSize
		if offset != "0" {
			size = 3
		}
		page := make([]ManagedClientAPI, size)
		for i := range page {
			page[i] = ManagedClientAPI{ID: offset + "-" + strconv.Itoa(i)}
		}
		writeJSON(t, w, http.StatusOK, page)
	})

	all, err := c.ListManagedClients(context.Background())
	if err != nil {
		t.Fatalf("ListManagedClients: %v", err)
	}
	if len(all) != managedClientsPageSize+3 {
		t.Errorf("len = %d, want %d", len(all), managedClientsPageSize+3)
	}
	if len(offsets) != 2 || offsets[0] != "0" || offsets[1] != strconv.Itoa(managedClientsPageSize) {
		t.Errorf("offsets = %v", offsets)
	}
}
