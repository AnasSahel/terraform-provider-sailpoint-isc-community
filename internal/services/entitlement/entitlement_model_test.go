// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func mustObjectRef(t *testing.T, typ, id, name string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(common.ObjectRefObjectType.AttrTypes, map[string]attr.Value{
		"type": types.StringValue(typ),
		"id":   types.StringValue(id),
		"name": types.StringValue(name),
	})
	if diags.HasError() {
		t.Fatalf("building object ref: %v", diags)
	}
	return obj
}

func mustStringSet(t *testing.T, values ...string) types.Set {
	t.Helper()
	set, diags := types.SetValueFrom(context.Background(), types.StringType, values)
	if diags.HasError() {
		t.Fatalf("building set: %v", diags)
	}
	return set
}

func baseState(t *testing.T) entitlementModel {
	t.Helper()
	return entitlementModel{
		ID:          types.StringValue("ent-1"),
		Name:        types.StringValue("Original"),
		Description: types.StringValue("desc"),
		Privileged:  types.BoolValue(false),
		Requestable: types.BoolValue(true),
		Owner:       mustObjectRef(t, "IDENTITY", "owner-1", "Owner One"),
		Source:      mustObjectRef(t, "SOURCE", "src-1", "AD"),
		Segments:    mustStringSet(t, "seg-1"),
	}
}

func opPaths(ops []client.JSONPatchOperation) []string {
	paths := make([]string, 0, len(ops))
	for _, op := range ops {
		paths = append(paths, op.Op+" "+op.Path)
	}
	return paths
}

// TestToPatchOperations_UnknownValuesAreSkipped covers issue #181: values not set in config
// are unknown in the plan and must not produce remove patches.
func TestToPatchOperations_UnknownValuesAreSkipped(t *testing.T) {
	state := baseState(t)
	plan := state
	plan.Name = types.StringValue("Updated")
	plan.Description = types.StringUnknown()
	plan.Privileged = types.BoolUnknown()
	plan.Requestable = types.BoolUnknown()
	plan.Owner = types.ObjectUnknown(common.ObjectRefObjectType.AttrTypes)
	plan.Source = types.ObjectUnknown(common.ObjectRefObjectType.AttrTypes)
	plan.Segments = types.SetUnknown(types.StringType)

	ops, diags := plan.ToPatchOperations(context.Background(), &state)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	got := opPaths(ops)
	if len(got) != 1 || got[0] != "replace /name" {
		t.Fatalf("expected only [replace /name], got %v", got)
	}
}

func TestToPatchOperations_NoChanges(t *testing.T) {
	state := baseState(t)
	plan := state

	ops, diags := plan.ToPatchOperations(context.Background(), &state)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(ops) != 0 {
		t.Fatalf("expected no ops, got %v", opPaths(ops))
	}
}

func TestToPatchOperations_ExplicitNullRemoves(t *testing.T) {
	state := baseState(t)
	plan := state
	plan.Description = types.StringNull()
	plan.Owner = types.ObjectNull(common.ObjectRefObjectType.AttrTypes)
	plan.Segments = types.SetNull(types.StringType)

	ops, diags := plan.ToPatchOperations(context.Background(), &state)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	want := map[string]bool{"remove /description": true, "remove /owner": true, "remove /segments": true}
	got := opPaths(ops)
	if len(got) != len(want) {
		t.Fatalf("expected %d ops, got %v", len(want), got)
	}
	for _, p := range got {
		if !want[p] {
			t.Fatalf("unexpected op %q in %v", p, got)
		}
	}
}

func TestToPatchOperations_OwnerChange(t *testing.T) {
	state := baseState(t)
	plan := state

	// Owner name unknown (server-resolved) but same id: no patch.
	sameOwner, diags := types.ObjectValue(common.ObjectRefObjectType.AttrTypes, map[string]attr.Value{
		"type": types.StringValue("IDENTITY"),
		"id":   types.StringValue("owner-1"),
		"name": types.StringUnknown(),
	})
	if diags.HasError() {
		t.Fatalf("building owner: %v", diags)
	}
	plan.Owner = sameOwner
	ops, diags := plan.ToPatchOperations(context.Background(), &state)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(ops) != 0 {
		t.Fatalf("expected no ops for unchanged owner id, got %v", opPaths(ops))
	}

	// New owner id: replace patch without the stale name.
	newOwner, diags := types.ObjectValue(common.ObjectRefObjectType.AttrTypes, map[string]attr.Value{
		"type": types.StringValue("IDENTITY"),
		"id":   types.StringValue("owner-2"),
		"name": types.StringUnknown(),
	})
	if diags.HasError() {
		t.Fatalf("building owner: %v", diags)
	}
	plan.Owner = newOwner
	ops, diags = plan.ToPatchOperations(context.Background(), &state)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if got := opPaths(ops); len(got) != 1 || got[0] != "replace /owner" {
		t.Fatalf("expected [replace /owner], got %v", got)
	}
	ref, ok := ops[0].Value.(*client.ObjectRefAPI)
	if !ok {
		t.Fatalf("expected *client.ObjectRefAPI value, got %T", ops[0].Value)
	}
	if ref.ID != "owner-2" || ref.Name != "" {
		t.Fatalf("unexpected owner patch value: %+v", ref)
	}
}

// TestPlanGet_UnknownComputedObjects reproduces the decode step of issue #181: a plan with
// unknown owner/source must decode into the model without a Value Conversion Error.
func TestPlanGet_UnknownComputedObjects(t *testing.T) {
	ctx := context.Background()
	r := NewEntitlementResource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", schemaResp.Diagnostics)
	}
	s := schemaResp.Schema
	objType, ok := s.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not an object")
	}

	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, tftypes.UnknownValue)
	}
	vals["id"] = tftypes.NewValue(tftypes.String, "ent-1")
	vals["name"] = tftypes.NewValue(tftypes.String, "Updated")

	plan := tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(objType, vals)}
	var m entitlementModel
	if diags := plan.Get(ctx, &m); diags.HasError() {
		t.Fatalf("plan.Get failed: %v", diags)
	}
	if !m.Source.IsUnknown() || !m.Owner.IsUnknown() {
		t.Fatalf("expected unknown owner/source, got owner=%v source=%v", m.Owner, m.Source)
	}
}

func TestFromAPI_NilObjectRefsAreNull(t *testing.T) {
	var m entitlementModel
	diags := m.FromAPI(context.Background(), &client.EntitlementAPI{ID: "ent-1", Name: "n"})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !m.Owner.IsNull() || !m.Source.IsNull() {
		t.Fatalf("expected null owner/source, got owner=%v source=%v", m.Owner, m.Source)
	}
}
