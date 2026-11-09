// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	if m.Source.IsNull() || m.Source.IsUnknown() {
		t.Fatal("expected a known source")
	}
	var src common.ObjectRefModel
	if diags := m.Source.As(context.Background(), &src, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("source: %v", diags)
	}
	if src.ID.ValueString() != "src-1" {
		t.Errorf("expected source.id=src-1, got %s", src.ID.ValueString())
	}
}

// TestToAPI_MergeLogic verifies that unlisted attributes keep their remote enabled value.
func TestToAPI_MergeLogic(t *testing.T) {
	t.Parallel()
	remote := makeTestAPI()

	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(false)},
	})
	result, diags := m.ToAPI(remote, nil)
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
	_, diags := m.ToAPI(remote, nil)
	if !diags.HasError() {
		t.Fatal("expected error for unknown attribute name, got none")
	}
}

// TestToAPI_SourceFromRemote verifies the PUT body carries the source returned by the GET,
// so an unknown planned `source` never reaches the request (#211).
func TestToAPI_SourceFromRemote(t *testing.T) {
	t.Parallel()
	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(true)},
	})
	m.Source = types.ObjectUnknown(common.ObjectRefObjectType.AttrTypes)

	result, diags := m.ToAPI(makeTestAPI(), nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if result.Source.ID != "src-1" || result.Source.Type != "SOURCE" {
		t.Errorf("source = %+v, want SOURCE/src-1", result.Source)
	}
}

// TestToAPI_CurrentDisplayName verifies a stale display name from the GET is replaced by the
// identity attribute's current one, and unknown names keep the remote value (#211).
func TestToAPI_CurrentDisplayName(t *testing.T) {
	t.Parallel()
	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(true)},
	})
	result, diags := m.ToAPI(makeTestAPI(), map[string]string{"email": "Work Email", "department": ""})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	got := map[string]string{}
	for _, a := range result.Attributes {
		got[a.Name] = a.DisplayName
	}
	if got["email"] != "Work Email" {
		t.Errorf("email displayName = %q, want %q", got["email"], "Work Email")
	}
	if got["department"] != "Department" {
		t.Errorf("department displayName = %q, want remote value %q", got["department"], "Department")
	}
	if got["firstname"] != "First Name" {
		t.Errorf("firstname displayName = %q, want remote value %q", got["firstname"], "First Name")
	}
}

// TestToAPI_AbsentDisabledIsConverged verifies a declared disabled entry that ISC dropped is
// skipped without error, while an absent enabled entry still fails (#211).
func TestToAPI_AbsentDisabledIsConverged(t *testing.T) {
	t.Parallel()
	m := makeModel([]sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(true)},
		{Name: types.StringValue("dropped"), Enabled: types.BoolValue(false)},
	})
	result, diags := m.ToAPI(makeTestAPI(), nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	for _, a := range result.Attributes {
		if a.Name == "dropped" {
			t.Error("absent entry must not be added to the PUT body")
		}
	}

	m.Attributes[1].Enabled = types.BoolValue(true)
	if _, diags := m.ToAPI(makeTestAPI(), nil); !diags.HasError() {
		t.Error("expected an error for an absent entry declared enabled")
	}
}

// TestKeepAbsentDisabled verifies the dropped disabled entry stays in state (no perpetual diff),
// and that absent enabled entries are not kept (real drift must show).
func TestKeepAbsentDisabled(t *testing.T) {
	t.Parallel()
	api := makeTestAPI()
	declared := []sourceAttrSyncAttributeModel{
		{Name: types.StringValue("email"), Enabled: types.BoolValue(true), DisplayName: types.StringUnknown(), Target: types.StringUnknown()},
		{Name: types.StringValue("dropped"), Enabled: types.BoolValue(false), DisplayName: types.StringUnknown(), Target: types.StringValue("x")},
		{Name: types.StringValue("gone"), Enabled: types.BoolValue(true), DisplayName: types.StringUnknown(), Target: types.StringUnknown()},
	}
	m := makeModel(nil)
	if diags := m.FromAPI(context.Background(), api, configuredAttributeNames(declared)); diags.HasError() {
		t.Fatal(diags)
	}
	m.keepAbsentDisabled(declared, api)

	byName := map[string]sourceAttrSyncAttributeModel{}
	for _, a := range m.Attributes {
		byName[a.Name.ValueString()] = a
	}
	if len(byName) != 2 {
		t.Fatalf("expected email + dropped in state, got %v", byName)
	}
	d, ok := byName["dropped"]
	if !ok {
		t.Fatal("dropped entry missing from state")
	}
	if d.Enabled.ValueBool() || !d.DisplayName.IsNull() || d.Target.ValueString() != "x" {
		t.Errorf("dropped entry = %+v, want enabled=false, display_name=null, target=x", d)
	}
	if _, ok := byName["gone"]; ok {
		t.Error("absent enabled entry must not be kept in state")
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
