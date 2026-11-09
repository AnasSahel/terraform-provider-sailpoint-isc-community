// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

// dynamicRoleJSON is a dynamic role as the roles API returns it, with a dimension schema and a form.
const dynamicRoleJSON = `{
  "id": "role-1",
  "name": "Dynamic Role",
  "enabled": true,
  "requestable": true,
  "dimensional": true,
  "owner": {"type": "IDENTITY", "id": "identity-1", "name": "Owner"},
  "segments": [],
  "dimensionRefs": [{"type": "DIMENSION", "id": "dimension-1", "name": "Dimension One"}],
  "accessRequestConfig": {
    "commentsRequired": false,
    "denialCommentsRequired": false,
    "reauthorizationRequired": false,
    "requireEndDate": false,
    "approvalSchemes": [{"approverType": "MANAGER"}],
    "formDefinitionId": "form-def-1",
    "dimensionSchema": {"dimensionAttributes": [
      {"name": "department", "displayName": "Department", "derived": true},
      {"name": "location", "displayName": "Location", "derived": false}
    ]}
  }
}`

func modelFromJSON(t *testing.T, raw string) roleModel {
	t.Helper()
	var api client.RoleAPI
	if err := json.Unmarshal([]byte(raw), &api); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	var m roleModel
	if diags := m.FromAPI(context.Background(), &api); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	return m
}

// withJSON returns the fixture with top-level keys overridden.
func withJSON(t *testing.T, overrides map[string]any) string {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal([]byte(dynamicRoleJSON), &doc); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	for k, v := range overrides {
		doc[k] = v
	}
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return string(out)
}

func decodeMap(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal %s: %v", raw, err)
	}
	return out
}

func TestAccessRequestConfigRoundTrip(t *testing.T) {
	m := modelFromJSON(t, dynamicRoleJSON)
	api, diags := m.ToAPI(context.Background())
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags)
	}
	got, err := json.Marshal(api.AccessRequestConfig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want, ok := decodeMap(t, []byte(dynamicRoleJSON))["accessRequestConfig"].(map[string]any)
	if !ok {
		t.Fatal("fixture has no accessRequestConfig object")
	}
	if gotMap := decodeMap(t, got); !reflect.DeepEqual(gotMap, want) {
		t.Errorf("accessRequestConfig round trip:\n got  %v\n want %v", gotMap, want)
	}
}

func TestFromAPISegmentsNeverNull(t *testing.T) {
	for name, raw := range map[string]string{
		"absent": `{"id": "role-1", "name": "r", "owner": {"type": "IDENTITY", "id": "identity-1"}}`,
		"empty":  `{"id": "role-1", "name": "r", "owner": {"type": "IDENTITY", "id": "identity-1"}, "segments": []}`,
	} {
		m := modelFromJSON(t, raw)
		if m.Segments.IsNull() || m.Segments.IsUnknown() || len(m.Segments.Elements()) != 0 {
			t.Errorf("%s: segments = %v, want a known empty set", name, m.Segments)
		}
	}
}

func TestToPatchOperationsSegmentsClearedSendsEmptyArray(t *testing.T) {
	ctx := context.Background()
	state := modelFromJSON(t, withJSON(t, map[string]any{"segments": []string{"seg-1"}}))
	plan := modelFromJSON(t, dynamicRoleJSON)
	ops, diags := plan.ToPatchOperations(ctx, &state)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	if len(ops) != 1 || ops[0].Op != "replace" || ops[0].Path != "/segments" {
		t.Fatalf("ops = %+v, want one replace /segments", ops)
	}
	body, err := json.Marshal(ops[0])
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(body), `"value":[]`) {
		t.Errorf("op = %s, want value []", body)
	}
}

func TestToPatchOperationsKeepsDimensionSchema(t *testing.T) {
	ctx := context.Background()
	state := modelFromJSON(t, dynamicRoleJSON)

	var doc map[string]any
	if err := json.Unmarshal([]byte(dynamicRoleJSON), &doc); err != nil {
		t.Fatal(err)
	}
	arc, ok := doc["accessRequestConfig"].(map[string]any)
	if !ok {
		t.Fatal("fixture has no accessRequestConfig object")
	}
	arc["commentsRequired"] = true
	plan := modelFromJSON(t, withJSON(t, map[string]any{
		"name":                "Renamed Role",
		"dimensional":         false,
		"accessRequestConfig": arc,
		"dimensionRefs": []map[string]string{
			{"type": "DIMENSION", "id": "dimension-1", "name": "Dimension One"},
			{"type": "DIMENSION", "id": "dimension-2", "name": "Dimension Two"},
		},
	}))

	ops, diags := plan.ToPatchOperations(ctx, &state)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	paths := make([]string, 0, len(ops))
	for _, op := range ops {
		paths = append(paths, op.Op+" "+op.Path)
	}
	if want := []string{"replace /name", "replace /accessRequestConfig"}; !reflect.DeepEqual(paths, want) {
		t.Fatalf("ops = %v, want %v", paths, want)
	}

	body, err := json.Marshal(ops[1].Value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var sent client.RequestabilityForRoleAPI
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatalf("unmarshal %s: %v", body, err)
	}
	if sent.DimensionSchema == nil || len(sent.DimensionSchema.DimensionAttributes) != 2 ||
		sent.DimensionSchema.DimensionAttributes[0].Name != "department" ||
		sent.DimensionSchema.DimensionAttributes[1].Name != "location" {
		t.Errorf("accessRequestConfig = %s, want both dimension attributes", body)
	}
	if sent.FormDefinitionID == nil || *sent.FormDefinitionID != "form-def-1" {
		t.Errorf("accessRequestConfig = %s, want formDefinitionId form-def-1", body)
	}
	if sent.CommentsRequired == nil || !*sent.CommentsRequired {
		t.Errorf("accessRequestConfig = %s, want commentsRequired true", body)
	}
}

func TestToAPIOmitsDimensionRefs(t *testing.T) {
	m := modelFromJSON(t, dynamicRoleJSON)
	if len(m.DimensionRefs.Elements()) != 1 {
		t.Fatalf("dimension_refs = %v, want one element", m.DimensionRefs)
	}
	api, diags := m.ToAPI(context.Background())
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags)
	}
	body, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, found := decodeMap(t, body)["dimensionRefs"]; found {
		t.Errorf("create body = %s, want no dimensionRefs", body)
	}
}
