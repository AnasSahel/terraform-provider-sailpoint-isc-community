// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func makeTestAPI() *client.SourceAttributeSyncConfigAPI {
	return &client.SourceAttributeSyncConfigAPI{
		Source: client.ObjectRefAPI{Type: "SOURCE", ID: "src-1", Name: "AD"},
		Attributes: []client.SourceAttributeSyncAttributeAPI{
			{Name: "email", DisplayName: "Email", Target: "mail", Enabled: true},
			{Name: "department", DisplayName: "Department", Target: "dept", Enabled: false},
			{Name: "firstname", DisplayName: "First Name", Target: "givenName", Enabled: false},
		},
	}
}

func makeModel(attrs []sourceAttrSyncAttributeModel) sourceAttrSyncConfigModel {
	const sourceID = "src-1"
	return sourceAttrSyncConfigModel{
		ID:         types.StringValue(sourceID),
		SourceID:   types.StringValue(sourceID),
		Attributes: attrs,
	}
}

// TestFromAPI_NoFilter verifies that nil cfgNames returns all attributes.
func TestFromAPI_NoFilter(t *testing.T) {
	t.Parallel()
	api := makeTestAPI()
	m := makeModel(nil)
	diags := m.FromAPI(context.Background(), api, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(m.Attributes) != 3 {
		t.Errorf("expected 3 attributes, got %d", len(m.Attributes))
	}
}

// TestFromAPI_WithFilter verifies that cfgNames limits the returned attributes.
func TestFromAPI_WithFilter(t *testing.T) {
	t.Parallel()
	api := makeTestAPI()
	m := makeModel(nil)
	filter := map[string]struct{}{"email": {}, "department": {}}
	diags := m.FromAPI(context.Background(), api, filter)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(m.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(m.Attributes))
	}
	for _, a := range m.Attributes {
		name := a.Name.ValueString()
		if _, ok := filter[name]; !ok {
			t.Errorf("unexpected attribute in filtered result: %s", name)
		}
	}
}

// TestFromAPI_RoundTrip verifies source ID is mirrored to id.
func TestFromAPI_RoundTrip(t *testing.T) {
	t.Parallel()
	api := makeTestAPI()
	m := makeModel(nil)
	diags := m.FromAPI(context.Background(), api, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if m.ID.ValueString() != "src-1" {
		t.Errorf("expected id=src-1, got %s", m.ID.ValueString())
	}
	if m.Source == nil {
		t.Fatal("expected non-nil source")
	}
	if m.Source.ID.ValueString() != "src-1" {
		t.Errorf("expected source.id=src-1, got %s", m.Source.ID.ValueString())
	}
}

// TestToAPI_MergeLogic verifies that unlisted attributes keep their remote enabled value.
func TestToAPI_MergeLogic(t *testing.T) {
	t.Parallel()
	remote := makeTestAPI()

	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(false)},
	})
	m.Source = &common.ObjectRefModel{
		Type: types.StringValue("SOURCE"),
		ID:   types.StringValue("src-1"),
		Name: types.StringValue("AD"),
	}

	result, diags := m.ToAPI(context.Background(), remote.Attributes)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(result.Attributes) != 3 {
		t.Errorf("expected 3 attributes in PUT body, got %d", len(result.Attributes))
	}

	attrMap := make(map[string]client.SourceAttributeSyncAttributeAPI)
	for _, a := range result.Attributes {
		attrMap[a.Name] = a
	}

	// email was overridden to false
	if attrMap["email"].Enabled {
		t.Error("expected email.enabled=false after override")
	}
	// department was not listed — should keep remote value (false)
	if attrMap["department"].Enabled {
		t.Error("expected department.enabled=false (unchanged from remote)")
	}
}

// TestToAPI_InvalidAttributeName verifies a clear error for unknown attribute names.
func TestToAPI_InvalidAttributeName(t *testing.T) {
	t.Parallel()
	remote := makeTestAPI()

	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("nonexistent"), Enabled: types.BoolValue(true)},
	})
	m.Source = &common.ObjectRefModel{
		Type: types.StringValue("SOURCE"),
		ID:   types.StringValue("src-1"),
		Name: types.StringValue("AD"),
	}

	_, diags := m.ToAPI(context.Background(), remote.Attributes)
	if !diags.HasError() {
		t.Fatal("expected error for unknown attribute name, got none")
	}
}

// TestConfiguredAttributeNames verifies the helper extracts names correctly.
func TestConfiguredAttributeNames(t *testing.T) {
	t.Parallel()
	attrs := []sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email")},
		{Name: types.StringValue("department")},
	}
	names := configuredAttributeNames(attrs)
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
	if _, ok := names["email"]; !ok {
		t.Error("expected 'email' in names")
	}
	if _, ok := names["department"]; !ok {
		t.Error("expected 'department' in names")
	}
}

// TestConfiguredAttributeNames_Empty verifies an empty slice returns an empty map.
func TestConfiguredAttributeNames_Empty(t *testing.T) {
	t.Parallel()
	names := configuredAttributeNames(nil)
	if len(names) != 0 {
		t.Errorf("expected empty map, got %d entries", len(names))
	}
}
