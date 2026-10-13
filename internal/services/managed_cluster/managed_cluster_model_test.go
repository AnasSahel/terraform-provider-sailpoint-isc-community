// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

func strPtr(s string) *string { return &s }

func TestFromAPIMapsFields(t *testing.T) {
	ctx := context.Background()
	operational := true
	count := int32(4)
	api := &client.ManagedClusterAPI{
		ID:           "c1",
		Name:         "Cluster",
		Type:         "idn",
		ClientType:   strPtr("CCG"),
		Operational:  &operational,
		Status:       strPtr("NORMAL"),
		ClientIDs:    []string{"a", "b"},
		ServiceCount: &count,
	}

	var m managedClusterModel
	if diags := m.FromAPI(ctx, api); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}

	if m.ID.ValueString() != "c1" || m.Name.ValueString() != "Cluster" || m.Type.ValueString() != "idn" {
		t.Errorf("id/name/type = %v/%v/%v", m.ID, m.Name, m.Type)
	}
	if !m.Operational.ValueBool() || m.ServiceCount.ValueInt64() != 4 {
		t.Errorf("operational/service_count = %v/%v", m.Operational, m.ServiceCount)
	}
	if !m.Description.IsNull() {
		t.Errorf("absent fields should be null: description=%v", m.Description)
	}
	if len(m.ClientIDs.Elements()) != 2 {
		t.Errorf("client_ids = %v", m.ClientIDs)
	}
}

// Absent collections must be empty, not null, so a second plan stays clean.
func TestFromAPIEmptyCollections(t *testing.T) {
	var m managedClusterModel
	if diags := m.FromAPI(context.Background(), &client.ManagedClusterAPI{ID: "c1"}); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	if m.ClientIDs.IsNull() || len(m.ClientIDs.Elements()) != 0 {
		t.Errorf("client_ids = %v, want empty list", m.ClientIDs)
	}
}

func TestMatchClustersByName(t *testing.T) {
	clusters := []client.ManagedClusterAPI{
		{ID: "1", Name: "Prod"},
		{ID: "2", Name: "prod"},
		{ID: "3", Name: "Dup"},
		{ID: "4", Name: "Dup"},
	}
	if got := matchClustersByName(clusters, "Prod"); len(got) != 1 || got[0].ID != "1" {
		t.Errorf("Prod: got %v", got)
	}
	if got := matchClustersByName(clusters, "Dup"); len(got) != 2 {
		t.Errorf("Dup: got %d matches, want 2", len(got))
	}
	if got := matchClustersByName(clusters, "missing"); len(got) != 0 {
		t.Errorf("missing: got %v", got)
	}
}
