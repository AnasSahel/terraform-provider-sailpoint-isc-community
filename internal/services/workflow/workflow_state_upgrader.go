// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// workflowModelV0 mirrors workflowModel for state shape v0. Field set is
// identical; the only on-the-wire difference is that definition.steps used
// jsontypes.NormalizedType. The framework decodes prior state into this
// struct via priorWorkflowSchemaV0.
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

// priorWorkflowSchemaV0 returns the workflow resource schema as it shipped
// in provider 2.4.3. Frozen — do not edit when v1 evolves; v0 is a
// historical artifact used only by the framework to decode existing state
// during upgrade. The only attribute that differs from v1 is
// definition.steps (jsontypes.NormalizedType vs. workflowStepsType).
func priorWorkflowSchemaV0() *schema.Schema {
	return &schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
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
		},
	}
}

// upgradeWorkflowStateV0ToV1 maps a v0 workflow state to v1. The only
// on-the-wire change between v0 and v1 is the framework type identity of
// definition.steps (jsontypes.NormalizedType → workflowStepsType). The
// underlying string is preserved verbatim — this is a no-op cast at the
// data level, but the framework requires an explicit upgrader because it
// rejects type-identity mismatches when reading prior state.
//
// See #114: provider 2.4.4 introduced workflowStepsType without this
// upgrader, breaking every plan/refresh on existing state until users
// pinned back to 2.4.3.
func upgradeWorkflowStateV0ToV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var prior workflowModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgraded := workflowModel{
		ID:             prior.ID,
		Name:           prior.Name,
		Owner:          prior.Owner,
		Description:    prior.Description,
		Trigger:        prior.Trigger,
		Enabled:        prior.Enabled,
		Created:        prior.Created,
		Modified:       prior.Modified,
		Creator:        prior.Creator,
		ModifiedBy:     prior.ModifiedBy,
		ExecutionCount: prior.ExecutionCount,
		FailureCount:   prior.FailureCount,
	}

	// Re-emit definition with the v1 attr types: steps wrapped in
	// workflowStepsValue. The underlying JSON string is unchanged.
	if prior.Definition.IsNull() {
		upgraded.Definition = types.ObjectNull(definitionAttrTypes)
	} else if prior.Definition.IsUnknown() {
		upgraded.Definition = types.ObjectUnknown(definitionAttrTypes)
	} else {
		priorAttrs := prior.Definition.Attributes()
		startAttr, _ := priorAttrs["start"].(types.String)
		stepsAttr, _ := priorAttrs["steps"].(jsontypes.Normalized)

		defObj, d := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
			"start": startAttr,
			"steps": workflowStepsValue{Normalized: stepsAttr},
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		upgraded.Definition = defObj
	}

	tflog.Info(ctx, "Upgraded workflow state from v0 to v1", map[string]any{
		"id":   upgraded.ID.ValueString(),
		"name": upgraded.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &upgraded)...)
}
