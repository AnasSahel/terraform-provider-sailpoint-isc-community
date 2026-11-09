// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"sync"
	"testing"
)

func TestRoleDimensionEndpoints(t *testing.T) {
	var mu sync.Mutex
	var seen []string
	var patchContentType string
	var createBody map[string]any
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		seen = append(seen, r.Method+" "+r.URL.Path)
		switch r.Method {
		case http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &createBody); err != nil {
				t.Errorf("decode create body: %v", err)
			}
			writeJSON(t, w, http.StatusCreated, RoleDimensionAPI{ID: "d1", Name: "n"})
		case http.MethodPatch:
			patchContentType = r.Header.Get("Content-Type")
			writeJSON(t, w, http.StatusOK, RoleDimensionAPI{ID: "d1", Name: "m"})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			writeJSON(t, w, http.StatusOK, RoleDimensionAPI{ID: "d1", Name: "n"})
		}
	})
	ctx := context.Background()

	created, err := c.CreateRoleDimension(ctx, "r1", &RoleDimensionAPI{
		Name:           "n",
		Owner:          &ObjectRefAPI{Type: "IDENTITY", ID: "i1"},
		AccessProfiles: []ObjectRefAPI{{Type: "ACCESS_PROFILE", ID: "ap1"}},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != "d1" {
		t.Errorf("created id = %q, want d1", created.ID)
	}
	if _, err := c.GetRoleDimension(ctx, "r1", "d1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := c.PatchRoleDimension(ctx, "r1", "d1", []JSONPatchOperation{NewReplacePatch("/name", "m")}); err != nil {
		t.Fatalf("patch: %v", err)
	}
	if err := c.DeleteRoleDimension(ctx, "r1", "d1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	want := []string{
		"POST /roles/v1/r1/dimensions",
		"GET /roles/v1/r1/dimensions/d1",
		"PATCH /roles/v1/r1/dimensions/d1",
		"DELETE /roles/v1/r1/dimensions/d1",
	}
	if !reflect.DeepEqual(seen, want) {
		t.Errorf("calls = %v\nwant    %v", seen, want)
	}
	if patchContentType != "application/json-patch+json" {
		t.Errorf("PATCH Content-Type = %q, want application/json-patch+json", patchContentType)
	}
	for _, k := range []string{"id", "parentId"} {
		if _, ok := createBody[k]; ok {
			t.Errorf("create body has %q: %v", k, createBody)
		}
	}
}

func TestRoleDimensionPatchWithoutOpsReads(t *testing.T) {
	var seen []string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		writeJSON(t, w, http.StatusOK, RoleDimensionAPI{ID: "d1", Name: "n"})
	})
	if _, err := c.PatchRoleDimension(context.Background(), "r1", "d1", nil); err != nil {
		t.Fatalf("patch: %v", err)
	}
	if want := []string{"GET /roles/v1/r1/dimensions/d1"}; !reflect.DeepEqual(seen, want) {
		t.Errorf("calls = %v, want %v", seen, want)
	}
}

func TestRoleDimensionNotFound(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]any{"detailCode": "404 Not found"})
	})
	ctx := context.Background()
	if _, err := c.GetRoleDimension(ctx, "r1", "d1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("get: err = %v, want ErrNotFound", err)
	}
	if err := c.DeleteRoleDimension(ctx, "r1", "d1"); err != nil {
		t.Errorf("delete: err = %v, want nil", err)
	}
}

// TestRoleDimensionDeleteRetriesServerError locks in that the client's retry
// config covers DELETE, so the 500 SailPoint returns when two dimensions of a
// role are deleted in parallel is retried.
func TestRoleDimensionDeleteRetriesServerError(t *testing.T) {
	var mu sync.Mutex
	deletes := 0
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		deletes++
		if deletes == 1 {
			writeJSON(t, w, http.StatusInternalServerError, map[string]any{"detailCode": "500"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if err := c.DeleteRoleDimension(context.Background(), "r1", "d1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if deletes != 2 {
		t.Errorf("DELETE calls = %d, want 2", deletes)
	}
}

func TestRoleDimensionEmptyArguments(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected HTTP call %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	ctx := context.Background()
	dim := &RoleDimensionAPI{Name: "n"}
	for name, call := range map[string]func() error{
		"get no role":    func() error { _, err := c.GetRoleDimension(ctx, "", "d1"); return err },
		"get no id":      func() error { _, err := c.GetRoleDimension(ctx, "r1", ""); return err },
		"create no role": func() error { _, err := c.CreateRoleDimension(ctx, "", dim); return err },
		"create nil":     func() error { _, err := c.CreateRoleDimension(ctx, "r1", nil); return err },
		"patch no role":  func() error { _, err := c.PatchRoleDimension(ctx, "", "d1", nil); return err },
		"patch no id":    func() error { _, err := c.PatchRoleDimension(ctx, "r1", "", nil); return err },
		"delete no role": func() error { return c.DeleteRoleDimension(ctx, "", "d1") },
		"delete no id":   func() error { return c.DeleteRoleDimension(ctx, "r1", "") },
	} {
		if err := call(); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}
