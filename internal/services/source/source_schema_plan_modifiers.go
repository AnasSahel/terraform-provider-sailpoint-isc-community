// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sourceSchemaAttributeComputedBools lists the Optional+Computed bool flags of
// a schema attribute. When omitted from config they are planned as unknown on
// every run unless a prior value is carried over.
var sourceSchemaAttributeComputedBools = []string{"is_multi", "is_entitlement", "is_group"}

// attributesUseStateForUnknownByName returns a set plan modifier for
// `attributes` that fills each element's unknown computed bools from the prior
// state element with the same `name`.
//
// A per-attribute `boolplanmodifier.UseStateForUnknown()` cannot be used here:
// for nested attributes inside a set, the framework pairs plan and state
// elements by iteration index, not by identity, so it would copy flags from an
// unrelated attribute whenever the two sets iterate in a different order.
// Matching on `name` (the key ISC uses for schema attributes) is stable.
func attributesUseStateForUnknownByName() planmodifier.Set {
	return attributesUseStateForUnknownByNameModifier{}
}

type attributesUseStateForUnknownByNameModifier struct{}

func (m attributesUseStateForUnknownByNameModifier) Description(_ context.Context) string {
	return "Carries over prior-state values of unset computed attribute flags, matching attributes by name."
}

func (m attributesUseStateForUnknownByNameModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m attributesUseStateForUnknownByNameModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// On Create there is no prior state; on destroy or an unknown set there is
	// nothing to correlate.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	priorByName := make(map[string]map[string]attr.Value, len(req.StateValue.Elements()))
	for _, elem := range req.StateValue.Elements() {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}
		attrs := obj.Attributes()
		name, ok := attrs["name"].(types.String)
		if !ok || name.IsNull() || name.IsUnknown() {
			continue
		}
		priorByName[name.ValueString()] = attrs
	}

	changed := false
	planElems := req.PlanValue.Elements()
	newElems := make([]attr.Value, 0, len(planElems))
	for _, elem := range planElems {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			newElems = append(newElems, elem)
			continue
		}
		attrs := obj.Attributes()
		name, ok := attrs["name"].(types.String)
		if !ok || name.IsNull() || name.IsUnknown() {
			newElems = append(newElems, elem)
			continue
		}
		prior, found := priorByName[name.ValueString()]
		if !found {
			// New attribute: leave its flags unknown so the API resolves them.
			newElems = append(newElems, elem)
			continue
		}

		updated := make(map[string]attr.Value, len(attrs))
		for k, v := range attrs {
			updated[k] = v
		}
		elemChanged := false
		for _, key := range sourceSchemaAttributeComputedBools {
			planned, ok := attrs[key].(types.Bool)
			if !ok || !planned.IsUnknown() {
				continue
			}
			priorVal, ok := prior[key].(types.Bool)
			if !ok || priorVal.IsUnknown() {
				continue
			}
			updated[key] = priorVal
			elemChanged = true
		}
		if !elemChanged {
			newElems = append(newElems, elem)
			continue
		}

		newObj, diags := types.ObjectValue(obj.AttributeTypes(ctx), updated)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		newElems = append(newElems, newObj)
		changed = true
	}

	if !changed {
		return
	}

	newSet, diags := types.SetValue(req.PlanValue.ElementType(ctx), newElems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = newSet
}
