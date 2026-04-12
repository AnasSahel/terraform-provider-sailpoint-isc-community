// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

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
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

type roleResource struct {
	client *client.Client
}

func NewRoleResource() resource.Resource { return &roleResource{} }

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "role resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Shared attribute helpers.

func objectRefSingle(desc string, required bool) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: desc,
		Required:            required,
		Optional:            !required,
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
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}
}

func approvalSchemesAttr(desc string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: desc,
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"approver_type": schema.StringAttribute{
					MarkdownDescription: "One of `OWNER`, `MANAGER`, `GOVERNANCE_GROUP`, `WORKFLOW`.",
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

// criteriaLeafAttrs is the level-3 (leaf) criteria attribute set.
func criteriaLeafAttrs() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operation": schema.StringAttribute{Required: true},
		"key": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"type":      schema.StringAttribute{Required: true},
				"property":  schema.StringAttribute{Required: true},
				"source_id": schema.StringAttribute{Optional: true},
			},
		},
		"string_value": schema.StringAttribute{Optional: true},
	}
}

func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for SailPoint Role.",
		MarkdownDescription: "Resource for SailPoint Role. Roles are the top of the access hierarchy — they bundle access profiles and" +
			" entitlements together and can define dynamic membership via criteria-based rules or explicit identity lists.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the role.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the role (max 128 chars).",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description (max 2000 chars).",
			},
			"enabled":         schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether the role is enabled."},
			"requestable":     schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether the role can be requested. Defaults to `false`."},
			"dimensional":     schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether this is a dimensional role."},
			"owner":           objectRefSingle("The owner of the role. Typically `type = IDENTITY`.", true),
			"access_profiles": objectRefSet("Access profiles bundled into this role."),
			"entitlements":    objectRefSet("Entitlements bundled directly into this role."),
			"segments": schema.SetAttribute{
				MarkdownDescription: "Segment UUIDs the role is visible in.",
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
			"membership": schema.SingleNestedAttribute{
				MarkdownDescription: "Role membership rules. `type = STANDARD` uses `criteria`; `type = IDENTITY_LIST` uses `identities`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "`STANDARD` or `IDENTITY_LIST`.",
						Required:            true,
					},
					"criteria": schema.SingleNestedAttribute{
						MarkdownDescription: "Criteria tree for STANDARD membership. Max 3 levels deep.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"operation": schema.StringAttribute{Required: true},
							"key": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"type":      schema.StringAttribute{Required: true},
									"property":  schema.StringAttribute{Required: true},
									"source_id": schema.StringAttribute{Optional: true},
								},
							},
							"string_value": schema.StringAttribute{Optional: true},
							"children": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operation": schema.StringAttribute{Required: true},
										"key": schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{
												"type":      schema.StringAttribute{Required: true},
												"property":  schema.StringAttribute{Required: true},
												"source_id": schema.StringAttribute{Optional: true},
											},
										},
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
					"identities": schema.ListNestedAttribute{
						MarkdownDescription: "Explicit identity list for IDENTITY_LIST membership. Max 500 modifications per PATCH.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id":   schema.StringAttribute{Required: true},
								"type": schema.StringAttribute{Optional: true, Computed: true},
								"name": schema.StringAttribute{
									Computed: true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"alias_name": schema.StringAttribute{
									Computed: true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"access_request_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"comments_required":        schema.BoolAttribute{Optional: true},
					"denial_comments_required": schema.BoolAttribute{Optional: true},
					"reauthorization_required": schema.BoolAttribute{Optional: true},
					"require_end_date":         schema.BoolAttribute{Optional: true},
					"max_permitted_access_duration": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"value":     schema.Int64Attribute{Optional: true},
							"time_unit": schema.StringAttribute{Optional: true},
						},
					},
					"approval_schemes": approvalSchemesAttr("Ordered approval chain for access requests."),
				},
			},
			"revoke_request_config": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"comments_required":        schema.BoolAttribute{Optional: true},
					"denial_comments_required": schema.BoolAttribute{Optional: true},
					"approval_schemes":         approvalSchemesAttr("Ordered approval chain for revoke requests."),
				},
			},
			"dimension_refs": objectRefSet("Dimensions referenced by this role (when dimensional=true)."),
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

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan roleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.CreateRole(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Role",
			fmt.Sprintf("Could not create role %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Creating SailPoint Role", "Received nil response from SailPoint API")
		return
	}

	var state roleModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created role", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state roleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	apiResp, err := r.client.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Role not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Role",
			fmt.Sprintf("Could not read role %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Role", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state roleModel
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

	apiResp, err := r.client.PatchRole(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Role",
			fmt.Sprintf("Could not update role %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Role", "Received nil response from SailPoint API")
		return
	}

	var newState roleModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state roleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if err := r.client.DeleteRole(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Role",
			fmt.Sprintf("Could not delete role %q: %s", id, err.Error()),
		)
		return
	}
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
