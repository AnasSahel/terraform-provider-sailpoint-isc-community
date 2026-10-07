// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"testing"
)

func TestListEntitlementsSendsFiltersAndLimit(t *testing.T) {
	var gotPath, gotFilters, gotLimit string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotFilters = r.URL.Query().Get("filters")
		gotLimit = r.URL.Query().Get("limit")
		writeJSON(t, w, http.StatusOK, []EntitlementAPI{{ID: "e1", Name: "Launcher"}})
	})

	got, err := c.ListEntitlements(context.Background(), `name eq "Launcher"`, 2)
	if err != nil {
		t.Fatalf("ListEntitlements: %v", err)
	}
	if gotPath != "/v2025/entitlements" {
		t.Errorf("path = %q, want /v2025/entitlements", gotPath)
	}
	if gotFilters != `name eq "Launcher"` {
		t.Errorf("filters = %q", gotFilters)
	}
	if gotLimit != "2" {
		t.Errorf("limit = %q, want 2", gotLimit)
	}
	if len(got) != 1 || got[0].ID != "e1" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestListEntitlementsOmitsEmptyFilters(t *testing.T) {
	var rawQuery string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		rawQuery = r.URL.RawQuery
		writeJSON(t, w, http.StatusOK, []EntitlementAPI{})
	})

	if _, err := c.ListEntitlements(context.Background(), "", 0); err != nil {
		t.Fatalf("ListEntitlements: %v", err)
	}
	if rawQuery != "" {
		t.Errorf("expected no query params, got %q", rawQuery)
	}
}
