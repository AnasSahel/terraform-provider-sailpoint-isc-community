// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// schemaAttr builds one `attributes` element. Pass nil for a flag to plan it
// as unknown (the value the framework gives an omitted Computed attribute).
func schemaAttr(t *testing.T, name string, isMulti, isEntitlement, isGroup *bool) attr.Value {
	t.Helper()
	flag := func(b *bool) types.Bool {
		if b == nil {
			return types.BoolUnknown()
		}
		return types.BoolValue(*b)
	}
	obj, diags := types.ObjectValue(sourceSchemaAttributeElementType().AttrTypes, map[string]attr.Value{
		"name":           types.StringValue(name),
		"native_name":    types.StringNull(),
		"type":           types.StringValue("STRING"),
		"description":    types.StringValue(name + " description"),
		"is_multi":       flag(isMulti),
		"is_entitlement": flag(isEntitlement),
		"is_group":       flag(isGroup),
		"schema":         types.ObjectNull(schemaAttributeSchemaRefAttrTypes),
	})
	if diags.HasError() {
		t.Fatalf("building attribute %q: %v", name, diags)
	}
	return obj
}

func schemaAttrSet(t *testing.T, elems ...attr.Value) types.Set {
	t.Helper()
	s, diags := types.SetValue(sourceSchemaAttributeElementType(), elems)
	if diags.HasError() {
		t.Fatalf("building set: %v", diags)
	}
	return s
}

func runAttributesModifier(t *testing.T, state, plan types.Set) types.Set {
	t.Helper()
	req := planmodifier.SetRequest{StateValue: state, PlanValue: plan, ConfigValue: plan}
	resp := &planmodifier.SetResponse{PlanValue: plan}
	attributesUseStateForUnknownByName().PlanModifySet(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("plan modifier returned diagnostics: %v", resp.Diagnostics)
	}
	return resp.PlanValue
}

func boolPtr(b bool) *bool { return &b }

// TestUnit_attributesUseStateForUnknownByName_sameSetDifferentOrder covers the
// #165 + #167 scenario: config omits the flags and lists attributes in a
// different order than the API returns them. The plan must equal prior state.
func TestUnit_attributesUseStateForUnknownByName_sameSetDifferentOrder(t *testing.T) {
	t.Parallel()

	state := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(true), boolPtr(true), boolPtr(true)),
		schemaAttr(t, "sAMAccountName", boolPtr(false), boolPtr(false), boolPtr(false)),
		schemaAttr(t, "mail", boolPtr(false), boolPtr(false), boolPtr(false)),
	)
	plan := schemaAttrSet(t,
		schemaAttr(t, "sAMAccountName", nil, nil, nil),
		schemaAttr(t, "mail", nil, nil, nil),
		schemaAttr(t, "memberOf", nil, nil, nil),
	)

	got := runAttributesModifier(t, state, plan)

	if !got.Equal(state) {
		t.Fatalf("plan should equal prior state (empty diff)\n got: %s\nwant: %s", got, state)
	}
}

// TestUnit_attributesUseStateForUnknownByName_addedAttribute checks that a new
// attribute keeps unknown flags while existing ones take their prior values.
func TestUnit_attributesUseStateForUnknownByName_addedAttribute(t *testing.T) {
	t.Parallel()

	state := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(true), boolPtr(true), boolPtr(true)),
	)
	plan := schemaAttrSet(t,
		schemaAttr(t, "memberOf", nil, nil, nil),
		schemaAttr(t, "department", nil, nil, nil),
	)

	got := runAttributesModifier(t, state, plan)

	want := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(true), boolPtr(true), boolPtr(true)),
		schemaAttr(t, "department", nil, nil, nil),
	)
	if !got.Equal(want) {
		t.Fatalf("unexpected plan\n got: %s\nwant: %s", got, want)
	}
}

// TestUnit_attributesUseStateForUnknownByName_configWins checks that a flag set
// in config is never overwritten by prior state.
func TestUnit_attributesUseStateForUnknownByName_configWins(t *testing.T) {
	t.Parallel()

	state := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(true), boolPtr(true), boolPtr(true)),
	)
	plan := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(false), nil, nil),
	)

	got := runAttributesModifier(t, state, plan)

	want := schemaAttrSet(t,
		schemaAttr(t, "memberOf", boolPtr(false), boolPtr(true), boolPtr(true)),
	)
	if !got.Equal(want) {
		t.Fatalf("unexpected plan\n got: %s\nwant: %s", got, want)
	}
}

// TestUnit_attributesUseStateForUnknownByName_create checks the modifier is a
// no-op when there is no prior state.
func TestUnit_attributesUseStateForUnknownByName_create(t *testing.T) {
	t.Parallel()

	plan := schemaAttrSet(t, schemaAttr(t, "memberOf", nil, nil, nil))

	got := runAttributesModifier(t, types.SetNull(sourceSchemaAttributeElementType()), plan)

	if !got.Equal(plan) {
		t.Fatalf("plan should be unchanged on create\n got: %s\nwant: %s", got, plan)
	}
}
