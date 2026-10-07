// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestApplicationEndpointsUseSourceAppsV1WithExperimentalHeader(t *testing.T) {
	var seen []string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		if got := r.Header.Get("X-SailPoint-Experimental"); got != "true" {
			t.Errorf("%s %s: X-SailPoint-Experimental = %q, want \"true\"", r.Method, r.URL.Path, got)
		}
		switch {
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/source-apps/v1/all":
			if got := r.URL.Query().Get("filters"); got != `name eq "App"` {
				t.Errorf("filters = %q", got)
			}
			writeJSON(t, w, http.StatusOK, []ApplicationAPI{{ID: "app1", Name: "App"}})
		case r.Method == http.MethodPatch:
			if got := r.Header.Get("Content-Type"); got != "application/json-patch+json" {
				t.Errorf("PATCH Content-Type = %q, want application/json-patch+json", got)
			}
			writeJSON(t, w, http.StatusOK, ApplicationAPI{ID: "app1", Name: "App"})
		default:
			writeJSON(t, w, http.StatusOK, ApplicationAPI{ID: "app1", Name: "App"})
		}
	})
	ctx := context.Background()

	if _, err := c.CreateApplication(ctx, &ApplicationCreateAPI{Name: "App", AccountSource: ApplicationCreateAccountSourceAPI{ID: "src1"}}); err != nil {
		t.Fatalf("CreateApplication: %v", err)
	}
	if _, err := c.GetApplication(ctx, "app1"); err != nil {
		t.Fatalf("GetApplication: %v", err)
	}
	if _, err := c.ListApplications(ctx, `name eq "App"`, 2); err != nil {
		t.Fatalf("ListApplications: %v", err)
	}
	if _, err := c.PatchApplication(ctx, "app1", []JSONPatchOperation{NewReplacePatch("/enabled", true)}); err != nil {
		t.Fatalf("PatchApplication: %v", err)
	}
	if err := c.DeleteApplication(ctx, "app1"); err != nil {
		t.Fatalf("DeleteApplication: %v", err)
	}

	want := []string{
		"POST /source-apps/v1",
		"GET /source-apps/v1/app1",
		"GET /source-apps/v1/all",
		"PATCH /source-apps/v1/app1",
		"DELETE /source-apps/v1/app1",
	}
	if len(seen) != len(want) {
		t.Fatalf("requests = %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("request %d = %q, want %q", i, seen[i], want[i])
		}
	}
}

func TestGetApplicationNotFound(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]string{"error": "not found"})
	})

	_, err := c.GetApplication(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
