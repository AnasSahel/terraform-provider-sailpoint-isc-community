// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package access_profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &accessProfileResource{}
	_ resource.ResourceWithConfigure   = &accessProfileResource{}
	_ resource.ResourceWithImportState = &accessProfileResource{}
)

type accessProfileResource struct {
	client *client.Client
}

func NewAccessProfileResource() resource.Resource { return &accessProfileResource{} }

func (r *accessProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_profile"
}

func (r *accessProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "access profile resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// objectRefNestedAttr builds a standard SingleNestedAttribute for {type, id, name} refs.
func objectRefNestedAttr(desc string, required bool) schema.SingleNestedAttribute {
	attrs := map[string]schema.Attribute{
		"type": schema.StringAttribute{Required: true},
		"id":   schema.StringAttribute{Required: true},
		"name": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
	return schema.SingleNestedAttribute{
		MarkdownDescription: desc,
		Required:            required,
		Optional:            !required,
		Attributes:          attrs,
	}
}

func approvalSchemesAttr(desc string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: desc,
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"approver_type": schema.StringAttribute{
					MarkdownDescription: "One of `APP_OWNER`, `OWNER`, `SOURCE_OWNER`, `MANAGER`, `GOVERNANCE_GROUP`, `WORKFLOW`.",
					Required:            true,
				},
				"approver_id": schema.StringAttribute{
					MarkdownDescription: "ID of the approver. Required when `approver_type` is `GOVERNANCE_GROUP` or `WORKFLOW`.",
					Optional:            true,
				},
			},
		},
	}
}

// level3Attrs is the leaf level of the provisioning criteria tree — no further children.
func level3Attrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operation": schema.StringAttribute{Required: true},
		"attribute": schema.StringAttribute{Optional: true},
		"value":     schema.StringAttribute{Optional: true},
	}
}

func (r *accessProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Access Profile.",
		MarkdownDescription: "Resource for SailPoint Access Profile. Access profiles bundle entitlements from a single source into a reusable unit that can be assigned to roles or requested directly.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the access profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the access profile.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description (max 2000 characters).",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the access profile is enabled. When `true`, at least one entitlement must be provided.",
			},
			"requestable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the access profile can be requested. Defaults to `true`.",
			},
			"owner":  objectRefNestedAttr("The owner of the access profile. Typically `type = IDENTITY`.", true),
			"source": objectRefNestedAttr("The source the access profile draws entitlements from. `type = SOURCE`.", true),
			"entitlements": schema.SetNestedAttribute{
				MarkdownDescription: "Entitlements bundled into this access profile.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Required: true},
						"id":   schema.StringAttribute{Required: true},
						"name": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"segments": schema.SetAttribute{
				MarkdownDescription: "Segment UUIDs this access profile is visible in.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"additional_owners": schema.SetNestedAttribute{
				MarkdownDescription: "Additional owners. Each may be `IDENTITY` or `GOVERNANCE_GROUP`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Required: true},
						"id":   schema.StringAttribute{Required: true},
						"name": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"access_request_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Access request configuration.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"comments_required":        schema.BoolAttribute{Optional: true},
					"denial_comments_required": schema.BoolAttribute{Optional: true},
					"reauthorization_required": schema.BoolAttribute{Optional: true},
					"require_end_date":         schema.BoolAttribute{Optional: true},
					"max_permitted_access_duration": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"value": schema.Int64Attribute{Optional: true},
							"time_unit": schema.StringAttribute{
								MarkdownDescription: "One of `HOURS`, `DAYS`, `WEEKS`, `MONTHS`.",
								Optional:            true,
							},
						},
					},
					"approval_schemes": approvalSchemesAttr("Ordered approval chain for access requests."),
				},
			},
			"revoke_request_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Revoke request configuration.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"approval_schemes": approvalSchemesAttr("Ordered approval chain for revoke requests."),
				},
			},
			"provisioning_criteria": schema.SingleNestedAttribute{
				MarkdownDescription: "Provisioning criteria tree. Max 3 levels deep.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"operation": schema.StringAttribute{
						MarkdownDescription: "Root operator: `AND`, `OR`, `EQUALS`, `NOT_EQUALS`, `CONTAINS`, `HAS`.",
						Required:            true,
					},
					"attribute": schema.StringAttribute{Optional: true},
					"value":     schema.StringAttribute{Optional: true},
					"children": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"operation": schema.StringAttribute{Required: true},
								"attribute": schema.StringAttribute{Optional: true},
								"value":     schema.StringAttribute{Optional: true},
								"children": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: level3Attrs(),
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

func (r *accessProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accessProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.CreateAccessProfile(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Access Profile",
			fmt.Sprintf("Could not create access profile %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Creating SailPoint Access Profile", "Received nil response from SailPoint API")
		return
	}

	var state accessProfileModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created access profile", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *accessProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accessProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	apiResp, err := r.client.GetAccessProfile(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Access profile not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Access Profile",
			fmt.Sprintf("Could not read access profile %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Access Profile", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accessProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accessProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state accessProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	ops, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.PatchAccessProfile(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Access Profile",
			fmt.Sprintf("Could not update access profile %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Access Profile", "Received nil response from SailPoint API")
		return
	}

	var newState accessProfileModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *accessProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accessProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if err := r.client.DeleteAccessProfile(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Access Profile",
			fmt.Sprintf("Could not delete access profile %q: %s", id, err.Error()),
		)
		return
	}
}

func (r *accessProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
