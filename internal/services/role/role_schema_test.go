// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func testSchema(t *testing.T) schema.Schema {
	t.Helper()
	var resp resource.SchemaResponse
	(&roleResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func testDataSourceSchema(t *testing.T) datasource.SchemaResponse {
	t.Helper()
	var resp datasource.SchemaResponse
	(&roleDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp
}

func TestSchemaValid(t *testing.T) {
	ctx := context.Background()
	if diags := testSchema(t).ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("resource schema invalid: %v", diags)
	}
	if diags := testDataSourceSchema(t).Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("data source schema invalid: %v", diags)
	}
}

func TestAttrTypesMatchSchema(t *testing.T) {
	arcWant := types.ObjectType{AttrTypes: accessRequestConfigAttrTypes}
	refsWant := types.SetType{ElemType: common.ObjectRefObjectType}
	rs, ds := testSchema(t), testDataSourceSchema(t).Schema
	for name, attrs := range map[string]map[string]attr.Type{
		"resource": {
			"access_request_config": rs.Attributes["access_request_config"].GetType(),
			"dimension_refs":        rs.Attributes["dimension_refs"].GetType(),
		},
		"data source": {
			"access_request_config": ds.Attributes["access_request_config"].GetType(),
			"dimension_refs":        ds.Attributes["dimension_refs"].GetType(),
		},
	} {
		if got := attrs["access_request_config"]; !got.Equal(arcWant) {
			t.Errorf("%s access_request_config: schema type %v != attr types %v", name, got, arcWant)
		}
		if got := attrs["dimension_refs"]; !got.Equal(refsWant) {
			t.Errorf("%s dimension_refs: schema type %v != %v", name, got, refsWant)
		}
	}
}

// runDimensionalModifiers runs the dimensional plan modifiers in schema order.
func runDimensionalModifiers(t *testing.T, state bool, plan, config types.Bool) planmodifier.BoolResponse {
	t.Helper()
	ctx := context.Background()
	sch := testSchema(t)
	a, ok := sch.Attributes["dimensional"].(schema.BoolAttribute)
	if !ok {
		t.Fatal("dimensional is not a BoolAttribute")
	}
	objType, ok := sch.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	raw := tftypes.NewValue(objType, vals)

	req := planmodifier.BoolRequest{
		State:       tfsdk.State{Schema: sch, Raw: raw},
		Plan:        tfsdk.Plan{Schema: sch, Raw: raw},
		Config:      tfsdk.Config{Schema: sch, Raw: raw},
		StateValue:  types.BoolValue(state),
		PlanValue:   plan,
		ConfigValue: config,
	}
	resp := planmodifier.BoolResponse{PlanValue: req.PlanValue}
	for _, m := range a.PlanModifiers {
		m.PlanModifyBool(ctx, req, &resp)
		req.PlanValue = resp.PlanValue
	}
	return resp
}

func TestDimensionalPlanModifierOrder(t *testing.T) {
	// Unconfigured: the Computed value plans as unknown, which must not replace the role.
	resp := runDimensionalModifiers(t, true, types.BoolUnknown(), types.BoolNull())
	if resp.RequiresReplace || !resp.PlanValue.Equal(types.BoolValue(true)) {
		t.Errorf("unconfigured: RequiresReplace=%v plan=%v, want false/true", resp.RequiresReplace, resp.PlanValue)
	}
	// Changed in config: the role must be replaced.
	resp = runDimensionalModifiers(t, true, types.BoolValue(false), types.BoolValue(false))
	if !resp.RequiresReplace {
		t.Error("changed: RequiresReplace=false, want true")
	}
}
