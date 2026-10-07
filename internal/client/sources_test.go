// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"testing"
)

func TestListSourcesSendsFiltersAndLimit(t *testing.T) {
	var gotPath, gotFilters, gotLimit string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotFilters = r.URL.Query().Get("filters")
		gotLimit = r.URL.Query().Get("limit")
		writeJSON(t, w, http.StatusOK, []SourceAPI{{ID: "s1", Name: "Active Directory"}})
	})

	got, err := c.ListSources(context.Background(), `name eq "Active Directory"`, 2)
	if err != nil {
		t.Fatalf("ListSources: %v", err)
	}
	if gotPath != "/sources/v1" {
		t.Errorf("path = %q, want /sources/v1", gotPath)
	}
	if gotFilters != `name eq "Active Directory"` {
		t.Errorf("filters = %q", gotFilters)
	}
	if gotLimit != "2" {
		t.Errorf("limit = %q, want 2", gotLimit)
	}
	if len(got) != 1 || got[0].ID != "s1" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestListSourcesReturnsErrorOnFailure(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusBadRequest, map[string]string{"detailCode": "400.1 Bad request syntax"})
	})

	if _, err := c.ListSources(context.Background(), `name eq "x"`, 2); err == nil {
		t.Fatal("expected an error on a 400 response")
	}
}
