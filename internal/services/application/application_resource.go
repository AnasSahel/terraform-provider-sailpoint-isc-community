// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package application

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

type applicationResource struct {
	client *client.Client
}

func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "application resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func accountSourceSchema(required bool) schema.SingleNestedAttribute {
	idAttr := schema.StringAttribute{
		MarkdownDescription: "The ID of the backing SailPoint source.",
	}
	if required {
		idAttr.Required = true
		idAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
	} else {
		idAttr.Computed = true
	}

	return schema.SingleNestedAttribute{
		MarkdownDescription: "The SailPoint source that owns this application.",
		Required:            required,
		Computed:            !required,
		Attributes: map[string]schema.Attribute{
			"id": idAttr,
			"type": schema.StringAttribute{
				MarkdownDescription: "The source reference type. Typically `SOURCE`.",
				Optional:            required,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the backing source.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownUnlessSiblingChanges("id"),
				},
			},
			"use_for_password_management": schema.BoolAttribute{
				MarkdownDescription: "Whether the backing source is used for password management.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"password_policies": schema.ListNestedAttribute{
				MarkdownDescription: "Password policies associated with the backing source.",
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Computed: true},
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for SailPoint Application (source app).",
		MarkdownDescription: "Resource for SailPoint Application (source app). Applications represent " +
			"individual apps within a SailPoint source for password management and access request workflows. " +
			"This resource uses the experimental `/source-apps/v1` API (`X-SailPoint-Experimental: true`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the application.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_app_id": schema.StringAttribute{
				MarkdownDescription: "The deprecated cloud application identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the application.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"provision_request_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether access requests are enabled for this application.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"match_all_accounts": schema.BoolAttribute{
				MarkdownDescription: "Whether the application matches all accounts on the backing source.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"app_center_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is visible in the request center.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_source": accountSourceSchema(true),
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the application.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The owner type. Must be `IDENTITY`.",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The owner ID.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The owner name. Resolved by the server.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							planmodifiers.UseStateForUnknownUnlessSiblingChanges("id"),
						},
					},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the application was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "When the application was last modified.",
				Computed:            true,
			},
		},
	}
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan applicationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToCreateAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating application", map[string]any{"name": plan.Name.ValueString()})
	apiResp, err := r.client.CreateApplication(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Application",
			fmt.Sprintf("Could not create application %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}
	if apiResp == nil || apiResp.ID == "" {
		resp.Diagnostics.AddError("Error Creating SailPoint Application", "Received empty response from SailPoint API")
		return
	}

	var state applicationModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ops, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(ops) > 0 {
		tflog.Debug(ctx, "Applying post-create application configuration", map[string]any{
			"id":        state.ID.ValueString(),
			"patch_ops": len(ops),
		})
		apiResp, err = r.client.PatchApplication(ctx, state.ID.ValueString(), ops)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating SailPoint Application",
				fmt.Sprintf("Application %q was created with ID %q, but configuring additional attributes failed: %s",
					plan.Name.ValueString(), state.ID.ValueString(), err.Error()),
			)
			return
		}
		if apiResp == nil {
			resp.Diagnostics.AddError("Error Creating SailPoint Application", "Received nil response while configuring application after create")
			return
		}
		resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created application", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	apiResp, err := r.client.GetApplication(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Application not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Application",
			fmt.Sprintf("Could not read application %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Application", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan applicationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state applicationModel
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

	apiResp, err := r.client.PatchApplication(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Application",
			fmt.Sprintf("Could not update application %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Application", "Received nil response from SailPoint API")
		return
	}

	var newState applicationModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated application", map[string]any{
		"id":   newState.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if err := r.client.DeleteApplication(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Application",
			fmt.Sprintf("Could not delete application %q: %s", id, err.Error()),
		)
		return
	}
	tflog.Info(ctx, "Successfully deleted application", map[string]any{"id": id})
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
