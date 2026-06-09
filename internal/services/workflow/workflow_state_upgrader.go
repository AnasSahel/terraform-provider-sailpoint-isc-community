// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// workflowModelV0 mirrors the state shape stored by provider ≤2.4.3.
// definition.steps was jsontypes.NormalizedType (plain JSON string).
type workflowModelV0 struct {
	ID             types.String         `tfsdk:"id"`
	Name           types.String         `tfsdk:"name"`
	Owner          types.Object         `tfsdk:"owner"`
	Description    types.String         `tfsdk:"description"`
	Definition     types.Object         `tfsdk:"definition"`
	Trigger        jsontypes.Normalized `tfsdk:"trigger"`
	Enabled        types.Bool           `tfsdk:"enabled"`
	Created        types.String         `tfsdk:"created"`
	Modified       types.String         `tfsdk:"modified"`
	Creator        types.Object         `tfsdk:"creator"`
	ModifiedBy     types.Object         `tfsdk:"modified_by"`
	ExecutionCount types.Int32          `tfsdk:"execution_count"`
	FailureCount   types.Int32          `tfsdk:"failure_count"`
}

// workflowModelV1 mirrors the state shape stored by provider 2.4.4–2.x.
// definition.steps was workflowStepsType (a custom jsontypes.Normalized subtype).
// The wire format is identical to V0 (a JSON string), so the same Go struct
// and jsontypes.NormalizedType schema field decodes it correctly.
type workflowModelV1 = workflowModelV0

// priorWorkflowSchemaV0 returns the workflow resource schema as it shipped
// in provider ≤2.4.3. Frozen — do not modify.
func priorWorkflowSchemaV0() *schema.Schema {
	return &schema.Schema{
		Version:    0,
		Attributes: workflowSchemaAttrsWithStringSteps(),
	}
}

// priorWorkflowSchemaV1 returns the workflow resource schema as it shipped
// in provider 2.4.4–2.x. The only on-the-wire difference from v0 is the
// custom type identity of definition.steps (workflowStepsType); both are
// stored as a plain JSON string. Frozen — do not modify.
func priorWorkflowSchemaV1() *schema.Schema {
	return &schema.Schema{
		Version:    1,
		Attributes: workflowSchemaAttrsWithStringSteps(),
	}
}

// workflowSchemaAttrsWithStringSteps returns the attribute map shared by
// priorWorkflowSchemaV0 and priorWorkflowSchemaV1. In both prior versions
// definition.steps was stored as a JSON string.
func workflowSchemaAttrsWithStringSteps() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"description": schema.StringAttribute{
			Optional: true,
		},
		"enabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			Default:  booldefault.StaticBool(false),
		},
		"trigger": schema.StringAttribute{
			Computed:   true,
			CustomType: jsontypes.NormalizedType{},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"modified": schema.StringAttribute{
			Computed: true,
		},
		"creator": schema.SingleNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{Computed: true},
				"id":   schema.StringAttribute{Computed: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
		"modified_by": schema.SingleNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{Computed: true},
				"id":   schema.StringAttribute{Computed: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
		"execution_count": schema.Int32Attribute{
			Computed: true,
		},
		"failure_count": schema.Int32Attribute{
			Computed: true,
		},
		"owner": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{Required: true},
				"id":   schema.StringAttribute{Required: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
		"definition": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"start": schema.StringAttribute{
					Required: true,
				},
				"steps": schema.StringAttribute{
					Required:   true,
					CustomType: jsontypes.NormalizedType{},
				},
			},
		},
	}
}

// upgradeDefinitionToV2 converts a prior-state definition object (steps as a
// JSON string) into the v2 definition object (steps as a typed map). Used by
// both the v0→v2 and v1→v2 upgraders.
func upgradeDefinitionToV2(priorDef types.Object) (types.Object, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	if priorDef.IsNull() {
		return types.ObjectNull(definitionAttrTypes), diagnostics
	}
	if priorDef.IsUnknown() {
		return types.ObjectUnknown(definitionAttrTypes), diagnostics
	}

	priorAttrs := priorDef.Attributes()
	startAttr, _ := priorAttrs["start"].(types.String)
	stepsAttr, _ := priorAttrs["steps"].(jsontypes.Normalized)

	var rawAllSteps map[string]any
	if !stepsAttr.IsNull() && !stepsAttr.IsUnknown() && stepsAttr.ValueString() != "" {
		if err := json.Unmarshal([]byte(stepsAttr.ValueString()), &rawAllSteps); err != nil {
			diagnostics.AddError("Error parsing steps JSON during state upgrade", err.Error())
			return types.ObjectNull(definitionAttrTypes), diagnostics
		}
	}

	stepsMap, diags := stepsFromAPI(rawAllSteps)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return types.ObjectNull(definitionAttrTypes), diagnostics
	}

	defObj, diags := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
		"start": startAttr,
		"steps": stepsMap,
	})
	diagnostics.Append(diags...)
	return defObj, diagnostics
}

// upgradeWorkflowStateV0ToV2 maps a v0 workflow state to v2. In v0,
// definition.steps was a plain JSON string (jsontypes.NormalizedType).
func upgradeWorkflowStateV0ToV2(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var prior workflowModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := workflowModel{
		ID:                prior.ID,
		Name:              prior.Name,
		Owner:             prior.Owner,
		Description:       prior.Description,
		Trigger:           prior.Trigger,
		Enabled:           prior.Enabled,
		Created:           prior.Created,
		Modified:          prior.Modified,
		Creator:           prior.Creator,
		ModifiedBy:        prior.ModifiedBy,
		ExecutionCount:    prior.ExecutionCount,
		FailureCount:      prior.FailureCount,
		IgnoreJSONChanges: types.ListNull(types.StringType),
	}

	defObj, diags := upgradeDefinitionToV2(prior.Definition)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	upgraded.Definition = defObj

	tflog.Info(ctx, "Upgraded workflow state from v0 to v2", map[string]any{
		"id":   upgraded.ID.ValueString(),
		"name": upgraded.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &upgraded)...)
}

// upgradeWorkflowStateV1ToV2 maps a v1 workflow state to v2. In v1,
// definition.steps was workflowStepsType — a custom string subtype with the
// same wire format as jsontypes.NormalizedType.
func upgradeWorkflowStateV1ToV2(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var prior workflowModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := workflowModel{
		ID:                prior.ID,
		Name:              prior.Name,
		Owner:             prior.Owner,
		Description:       prior.Description,
		Trigger:           prior.Trigger,
		Enabled:           prior.Enabled,
		Created:           prior.Created,
		Modified:          prior.Modified,
		Creator:           prior.Creator,
		ModifiedBy:        prior.ModifiedBy,
		ExecutionCount:    prior.ExecutionCount,
		FailureCount:      prior.FailureCount,
		IgnoreJSONChanges: types.ListNull(types.StringType),
	}

	defObj, diags := upgradeDefinitionToV2(prior.Definition)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	upgraded.Definition = defObj

	tflog.Info(ctx, "Upgraded workflow state from v1 to v2", map[string]any{
		"id":   upgraded.ID.ValueString(),
		"name": upgraded.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &upgraded)...)
}
