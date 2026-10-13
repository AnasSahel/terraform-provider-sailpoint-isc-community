// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_client

import (
	"encoding/json"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

func strPtr(s string) *string { return &s }

func TestMatchManagedClients(t *testing.T) {
	all := []client.ManagedClientAPI{
		{ID: "a", Name: strPtr("VA-1"), ClusterID: "c1"},
		{ID: "b", Name: strPtr("VA-1"), ClusterID: "c2"},
		{ID: "c", Name: strPtr("VA-2"), ClusterID: "c1"},
		{ID: "d", Name: nil, ClusterID: "c1"},
	}

	tests := []struct {
		name, lookup, cluster string
		want                  []string
	}{
		{"name only, ambiguous", "VA-1", "", []string{"a", "b"}},
		{"name narrowed by cluster", "VA-1", "c2", []string{"b"}},
		{"unique name", "VA-2", "", []string{"c"}},
		{"no match", "VA-3", "", nil},
		{"name in wrong cluster", "VA-2", "c2", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchManagedClients(all, tt.lookup, tt.cluster)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d matches, want %d", len(got), len(tt.want))
			}
			for i, mc := range got {
				if mc.ID != tt.want[i] {
					t.Errorf("match %d = %q, want %q", i, mc.ID, tt.want[i])
				}
			}
		})
	}
}

func TestManagedClientFromAPI(t *testing.T) {
	var m managedClientModel
	m.FromAPI(&client.ManagedClientAPI{
		ID:        "a",
		Name:      strPtr("VA-1"),
		ClusterID: "c1",
		Status:    strPtr("NORMAL"),
	})

	if m.ID.ValueString() != "a" || m.Name.ValueString() != "VA-1" || m.ClusterID.ValueString() != "c1" || m.Status.ValueString() != "NORMAL" {
		t.Errorf("unexpected model: %+v", m)
	}
	if !m.VaVersion.IsNull() {
		t.Errorf("absent fields should be null: %+v", m)
	}
}

func TestRawToString(t *testing.T) {
	tests := []struct {
		name   string
		raw    json.RawMessage
		want   string
		isNull bool
	}{
		{"nil", nil, "", true},
		{"null", json.RawMessage(`null`), "", true},
		{"string", json.RawMessage(`"20000"`), "20000", false},
		{"number", json.RawMessage(`15000`), "15000", false},
		{"bool", json.RawMessage(`true`), "true", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rawToString(tt.raw)
			if got.IsNull() != tt.isNull || got.ValueString() != tt.want {
				t.Errorf("got %v, want %q (null=%v)", got, tt.want, tt.isNull)
			}
		})
	}
}
