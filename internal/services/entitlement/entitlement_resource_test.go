// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func mustResourceSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var resp resource.SchemaResponse
	NewEntitlementResource().Schema(ctx, resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	objType, ok := resp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not an object")
	}
	return resp.Schema, objType
}

// objectValue builds a value of the resource type where every attribute not in
// vals takes the given fill value (nil for null, tftypes.UnknownValue for unknown).
func objectValue(objType tftypes.Object, vals map[string]tftypes.Value, fill any) tftypes.Value {
	all := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		if v, ok := vals[name]; ok {
			all[name] = v
		} else {
			all[name] = tftypes.NewValue(typ, fill)
		}
	}
	return tftypes.NewValue(objType, all)
}

func runModifyPlan(t *testing.T, config, state, plan map[string]tftypes.Value) resource.ModifyPlanResponse {
	t.Helper()
	sch, objType := mustResourceSchema(t)
	r, ok := NewEntitlementResource().(resource.ResourceWithModifyPlan)
	if !ok {
		t.Fatalf("resource does not implement ModifyPlan")
	}
	stateRaw := tftypes.NewValue(objType, nil)
	if state != nil {
		stateRaw = objectValue(objType, state, nil)
	}
	planTF := tfsdk.Plan{Schema: sch, Raw: objectValue(objType, plan, tftypes.UnknownValue)}
	resp := resource.ModifyPlanResponse{Plan: planTF}
	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{
		Config: tfsdk.Config{Schema: sch, Raw: objectValue(objType, config, nil)},
		State:  tfsdk.State{Schema: sch, Raw: stateRaw},
		Plan:   planTF,
	}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ModifyPlan diagnostics: %v", resp.Diagnostics)
	}
	return resp
}

func TestModifyPlan_PrivilegeLevelDerivesPrivileged(t *testing.T) {
	id := tftypes.NewValue(tftypes.String, "ent-1")
	resp := runModifyPlan(t,
		map[string]tftypes.Value{"id": id, "privilege_level": tftypes.NewValue(tftypes.String, "LOW")},
		map[string]tftypes.Value{"id": id, "privilege_level": tftypes.NewValue(tftypes.String, "NONE"), "privileged": tftypes.NewValue(tftypes.Bool, false)},
		map[string]tftypes.Value{"id": id, "privilege_level": tftypes.NewValue(tftypes.String, "LOW"), "privileged": tftypes.NewValue(tftypes.Bool, false)},
	)
	var got types.Bool
	resp.Diagnostics.Append(resp.Plan.GetAttribute(context.Background(), path.Root("privileged"), &got)...)
	if !got.Equal(types.BoolValue(true)) {
		t.Fatalf("expected privileged=true in plan, got %v", got)
	}
}

func TestModifyPlan_PrivilegedChangeUnknownsLevel(t *testing.T) {
	id := tftypes.NewValue(tftypes.String, "ent-1")
	resp := runModifyPlan(t,
		map[string]tftypes.Value{"id": id, "privileged": tftypes.NewValue(tftypes.Bool, true)},
		map[string]tftypes.Value{"id": id, "privilege_level": tftypes.NewValue(tftypes.String, "NONE"), "privileged": tftypes.NewValue(tftypes.Bool, false)},
		map[string]tftypes.Value{"id": id, "privilege_level": tftypes.NewValue(tftypes.String, "NONE"), "privileged": tftypes.NewValue(tftypes.Bool, true)},
	)
	var got types.String
	resp.Diagnostics.Append(resp.Plan.GetAttribute(context.Background(), path.Root("privilege_level"), &got)...)
	if !got.IsUnknown() {
		t.Fatalf("expected unknown privilege_level in plan, got %v", got)
	}
}

func TestValidateConfig_Privilege(t *testing.T) {
	sch, objType := mustResourceSchema(t)
	r, ok := NewEntitlementResource().(resource.ResourceWithValidateConfig)
	if !ok {
		t.Fatalf("resource does not implement ValidateConfig")
	}
	cases := []struct {
		name    string
		vals    map[string]tftypes.Value
		wantErr bool
	}{
		{"valid level", map[string]tftypes.Value{"privilege_level": tftypes.NewValue(tftypes.String, "HIGH")}, false},
		{"invalid level", map[string]tftypes.Value{"privilege_level": tftypes.NewValue(tftypes.String, "high")}, true},
		{"both set", map[string]tftypes.Value{
			"privilege_level": tftypes.NewValue(tftypes.String, "HIGH"),
			"privileged":      tftypes.NewValue(tftypes.Bool, true),
		}, true},
		{"deprecated only", map[string]tftypes.Value{"privileged": tftypes.NewValue(tftypes.Bool, true)}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var resp resource.ValidateConfigResponse
			r.ValidateConfig(context.Background(), resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: sch, Raw: objectValue(objType, tc.vals, nil)},
			}, &resp)
			if resp.Diagnostics.HasError() != tc.wantErr {
				t.Fatalf("wantErr=%v, got %v", tc.wantErr, resp.Diagnostics)
			}
		})
	}
}
