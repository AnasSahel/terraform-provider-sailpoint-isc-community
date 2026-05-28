// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package governance_group

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &governanceGroupResource{}
	_ resource.ResourceWithConfigure   = &governanceGroupResource{}
	_ resource.ResourceWithImportState = &governanceGroupResource{}
)

type governanceGroupResource struct {
	client *client.Client
}

// NewGovernanceGroupResource creates a new Governance Group resource.
func NewGovernanceGroupResource() resource.Resource {
	return &governanceGroupResource{}
}

func (r *governanceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_group"
}

func (r *governanceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "governance group resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func memberNestedAttr(required bool) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Identities that are members of this governance group. " +
			"Member changes are applied incrementally (add/remove) rather than replacing the full list.",
		Optional: !required,
		Required: required,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					MarkdownDescription: "Object type. Must be `IDENTITY`.",
					Required:            true,
				},
				"id": schema.StringAttribute{
					MarkdownDescription: "Identity ID of the member.",
					Required:            true,
				},
				"name": schema.StringAttribute{
					MarkdownDescription: "Name of the member identity. Server-resolved.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						planmodifiers.UseStateForUnknownUnlessSiblingChanges("id"),
					},
				},
			},
		},
	}
}

func (r *governanceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for SailPoint Governance Group.",
		MarkdownDescription: "Resource for SailPoint Governance Group. " +
			"Governance groups are collections of identities used as approvers in access request and " +
			"certification workflows. They can be referenced in role `approval_schemes` and certification " +
			"campaign reviewer configurations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the governance group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the governance group.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the governance group.",
			},
			"owner": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The owner of the governance group. Typically `type = IDENTITY`.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Required: true},
					"id":   schema.StringAttribute{Required: true},
					"name": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							planmodifiers.UseStateForUnknownUnlessSiblingChanges("id"),
						},
					},
				},
			},
			"members":      memberNestedAttr(false),
			"member_count": schema.Int64Attribute{Computed: true, MarkdownDescription: "Total number of members. Server-resolved."},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *governanceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan governanceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.CreateGovernanceGroup(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Governance Group",
			fmt.Sprintf("Could not create governance group %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Creating SailPoint Governance Group", "Received nil response from SailPoint API")
		return
	}

	id := apiResp.ID

	// Add any declared members after creation.
	if len(plan.Members) > 0 {
		membersToAdd := make([]client.GovernanceGroupMemberAPI, 0, len(plan.Members))
		for _, m := range plan.Members {
			membersToAdd = append(membersToAdd, client.GovernanceGroupMemberAPI{
				Type: m.Type.ValueString(),
				ID:   m.ID.ValueString(),
			})
		}
		if err := r.client.AddGovernanceGroupMembers(ctx, id, membersToAdd); err != nil {
			resp.Diagnostics.AddError(
				"Error Adding Members to SailPoint Governance Group",
				fmt.Sprintf("Governance group %q was created but members could not be added: %s", id, err.Error()),
			)
			return
		}
	}

	members, err := r.client.ListGovernanceGroupMembers(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Governance Group Members",
			fmt.Sprintf("Could not read members for governance group %q: %s", id, err.Error()),
		)
		return
	}

	var state governanceGroupModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp, members)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created governance group", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *governanceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state governanceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	apiResp, err := r.client.GetGovernanceGroup(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Governance group not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Governance Group",
			fmt.Sprintf("Could not read governance group %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Governance Group", "Received nil response from SailPoint API")
		return
	}

	members, err := r.client.ListGovernanceGroupMembers(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Governance Group Members",
			fmt.Sprintf("Could not read members for governance group %q: %s", id, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp, members)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *governanceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan governanceGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state governanceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	// Apply PATCH for scalar fields.
	ops, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.PatchGovernanceGroup(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Governance Group",
			fmt.Sprintf("Could not update governance group %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Governance Group", "Received nil response from SailPoint API")
		return
	}

	// Reconcile member additions and removals.
	toAdd, toRemoveIDs := MemberDiff(plan.Members, state.Members)

	if len(toRemoveIDs) > 0 {
		if err := r.client.RemoveGovernanceGroupMembers(ctx, id, toRemoveIDs); err != nil {
			resp.Diagnostics.AddError(
				"Error Removing Members from SailPoint Governance Group",
				fmt.Sprintf("Could not remove members from governance group %q: %s", id, err.Error()),
			)
			return
		}
	}

	if len(toAdd) > 0 {
		if err := r.client.AddGovernanceGroupMembers(ctx, id, toAdd); err != nil {
			resp.Diagnostics.AddError(
				"Error Adding Members to SailPoint Governance Group",
				fmt.Sprintf("Could not add members to governance group %q: %s", id, err.Error()),
			)
			return
		}
	}

	// Re-read final state to pick up server-resolved names on members.
	members, err := r.client.ListGovernanceGroupMembers(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Governance Group Members",
			fmt.Sprintf("Could not read members for governance group %q after update: %s", id, err.Error()),
		)
		return
	}

	var newState governanceGroupModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp, members)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated governance group", map[string]any{
		"id":      newState.ID.ValueString(),
		"patches": len(ops),
	})
}

func (r *governanceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state governanceGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if err := r.client.DeleteGovernanceGroup(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Governance Group",
			fmt.Sprintf("Could not delete governance group %q: %s", id, err.Error()),
		)
		return
	}
}

func (r *governanceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
