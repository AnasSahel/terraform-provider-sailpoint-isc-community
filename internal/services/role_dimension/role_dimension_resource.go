// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = &roleDimensionResource{}
	_ resource.ResourceWithConfigure   = &roleDimensionResource{}
	_ resource.ResourceWithImportState = &roleDimensionResource{}
)

type roleDimensionResource struct {
	client *client.Client
}

// NewRoleDimensionResource returns the sailpoint_role_dimension resource.
func NewRoleDimensionResource() resource.Resource { return &roleDimensionResource{} }

func (r *roleDimensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_dimension"
}

func (r *roleDimensionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "role dimension resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Shared attribute helpers.

func objectRefSingle(desc string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: desc,
		Optional:            true,
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
	}
}

func objectRefSet(desc string) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: desc,
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
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
	}
}

// criteriaKeyAttr is the criteria key. Dimension keys have no source_id.
func criteriaKeyAttr() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"type":     schema.StringAttribute{Required: true, MarkdownDescription: "Only `IDENTITY`."},
			"property": schema.StringAttribute{Required: true, MarkdownDescription: "Identity attribute, for example `attribute.department`."},
		},
	}
}

// criteriaLeafAttrs is the level-3 (leaf) criteria attribute set.
func criteriaLeafAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operation":    schema.StringAttribute{Required: true},
		"key":          criteriaKeyAttr(),
		"string_value": schema.StringAttribute{Optional: true},
	}
}

func (r *roleDimensionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for a SailPoint role dimension.",
		MarkdownDescription: "Resource for a dimension of a SailPoint dynamic role. The role must be dimensional (`dimensional = true`)." +
			" A dimension grants its access profiles and entitlements to the identities that match its membership criteria." +
			" Only changed fields are patched. Import with `<role_id>/<dimension_id>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the dimension.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the dimensional role the dimension belongs to. Changing it forces a new dimension.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the dimension (max 128 chars).",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description (max 2000 chars).",
			},
			"owner": objectRefSingle("The owner of the dimension. Typically `type = IDENTITY`. SailPoint stores an owner name" +
				" only when one is sent, and the provider sends only `type` and `id`, so `owner.name` stays empty for dimensions" +
				" created through the API."),
			"access_profiles": objectRefSet("Access profiles the dimension grants. This is a set matched by `id`, so the order" +
				" SailPoint returns them in does not matter."),
			"entitlements": objectRefSet("Entitlements the dimension grants. This is a set matched by `id`, so the order" +
				" SailPoint returns them in does not matter."),
			"membership": schema.SingleNestedAttribute{
				MarkdownDescription: "Dimension membership rules. Only `type = STANDARD` with `criteria` is supported." +
					" Saving the dimension in the SailPoint UI can clear its membership; the next apply restores it.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Only `STANDARD`.",
						Required:            true,
					},
					"criteria": schema.SingleNestedAttribute{
						MarkdownDescription: "Criteria tree. Max 3 levels deep. `AND` and `OR` nodes carry no `string_value`.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"operation":    schema.StringAttribute{Required: true},
							"key":          criteriaKeyAttr(),
							"string_value": schema.StringAttribute{Optional: true},
							"children": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operation":    schema.StringAttribute{Required: true},
										"key":          criteriaKeyAttr(),
										"string_value": schema.StringAttribute{Optional: true},
										"children": schema.ListNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: criteriaLeafAttrs(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
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

func (r *roleDimensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleDimensionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := plan.RoleID.ValueString()
	apiResp, err := r.client.CreateRoleDimension(ctx, roleID, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Role Dimension",
			fmt.Sprintf("Could not create dimension %q on role %q: %s", plan.Name.ValueString(), roleID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Creating SailPoint Role Dimension", "Received nil response from SailPoint API")
		return
	}

	var state roleDimensionModel
	resp.Diagnostics.Append(state.FromAPI(ctx, roleID, apiResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created role dimension", map[string]any{
		"role_id": roleID,
		"id":      state.ID.ValueString(),
	})
}

func (r *roleDimensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleDimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.RoleID.ValueString()
	id := state.ID.ValueString()
	apiResp, err := r.client.GetRoleDimension(ctx, roleID, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Role dimension not found, removing from state", map[string]any{"role_id": roleID, "id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Role Dimension",
			fmt.Sprintf("Could not read dimension %q of role %q: %s", id, roleID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Role Dimension", "Received nil response from SailPoint API")
		return
	}

	var newState roleDimensionModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, roleID, apiResp, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *roleDimensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleDimensionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state roleDimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.RoleID.ValueString()
	id := state.ID.ValueString()
	ops, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.PatchRoleDimension(ctx, roleID, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Role Dimension",
			fmt.Sprintf("Could not update dimension %q of role %q: %s", id, roleID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Role Dimension", "Received nil response from SailPoint API")
		return
	}

	var newState roleDimensionModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, roleID, apiResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *roleDimensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleDimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := state.RoleID.ValueString()
	id := state.ID.ValueString()
	if err := r.client.DeleteRoleDimension(ctx, roleID, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Role Dimension",
			fmt.Sprintf("Could not delete dimension %q of role %q: %s", id, roleID, err.Error()),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
// Import format: role_id/dimension_id.
func (r *roleDimensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: role_id/dimension_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
