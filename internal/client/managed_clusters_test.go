// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestListManagedClustersFollowsPagination(t *testing.T) {
	var offsets []string
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/managed-clusters/v1" {
			t.Errorf("path = %q, want /managed-clusters/v1", r.URL.Path)
		}
		offsets = append(offsets, r.URL.Query().Get("offset"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		n := managedClustersPageSize
		if offset > 0 {
			n = 3
		}
		page := make([]ManagedClusterAPI, n)
		for i := range page {
			page[i] = ManagedClusterAPI{ID: fmt.Sprintf("c%d", offset+i)}
		}
		writeJSON(t, w, http.StatusOK, page)
	})

	clusters, err := c.ListManagedClusters(context.Background())
	if err != nil {
		t.Fatalf("ListManagedClusters: %v", err)
	}
	if len(clusters) != managedClustersPageSize+3 {
		t.Errorf("len = %d, want %d", len(clusters), managedClustersPageSize+3)
	}
	if want := []string{"0", strconv.Itoa(managedClustersPageSize)}; fmt.Sprint(offsets) != fmt.Sprint(want) {
		t.Errorf("offsets = %v, want %v", offsets, want)
	}
}

func TestGetManagedCluster(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/managed-clusters/v1/c1" {
			t.Errorf("path = %q, want /managed-clusters/v1/c1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"c1","name":"Cluster","keyPair":{"publicKey":"secret"},"attributes":{"keystore":"secret"}}`))
	})

	got, err := c.GetManagedCluster(context.Background(), "c1")
	if err != nil {
		t.Fatalf("GetManagedCluster: %v", err)
	}
	if got.ID != "c1" || got.Name != "Cluster" {
		t.Errorf("got %+v", got)
	}
}

func TestGetManagedClusterNotFound(t *testing.T) {
	c := newGovernanceGroupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusNotFound, map[string]string{"error": "not found"})
	})

	_, err := c.GetManagedCluster(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
