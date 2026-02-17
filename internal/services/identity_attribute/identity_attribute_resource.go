// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_attribute

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	_ resource.Resource                = &identityAttributeResource{}
	_ resource.ResourceWithConfigure   = &identityAttributeResource{}
	_ resource.ResourceWithImportState = &identityAttributeResource{}
)

type identityAttributeResource struct {
	client *client.Client
}

func NewIdentityAttributeResource() resource.Resource {
	return &identityAttributeResource{}
}

// Metadata implements resource.Resource.
func (r *identityAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_attribute"
}

// Configure implements resource.ResourceWithConfigure.
func (r *identityAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "identity attribute resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *identityAttributeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Identity Attribute.",
		MarkdownDescription: "Resource for SailPoint Identity Attribute.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity attribute.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the identity attribute.",
				Optional:            true,
				Computed:            true,
			},
			"standard": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is a standard attribute.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the identity attribute.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multi": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute supports multiple values.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"searchable": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is searchable.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is a system attribute.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sources": schema.ListNestedAttribute{
				MarkdownDescription: "The sources associated with the identity attribute.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Attribute mapping type. Mostly `rule`.",
							Optional:            true,
							Computed:            true,
						},
						"properties": schema.StringAttribute{
							MarkdownDescription: "Attribute mapping properties.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *identityAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the plan
	var plan identityAttributeModel
	tflog.Debug(ctx, "Getting plan for identity attribute resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping identity attribute resource model to API create request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the identity attribute via the API client
	tflog.Debug(ctx, "Creating identity attribute via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	identityAttributeAPIResponse, err := r.client.CreateIdentityAttribute(ctx, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Identity Attribute",
			fmt.Sprintf("Could not create SailPoint Identity Attribute %q: %s", plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Identity Attribute", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if identityAttributeAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Identity Attribute",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state identityAttributeModel
	tflog.Debug(ctx, "Mapping SailPoint Identity Attribute API response to resource model", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *identityAttributeAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for identity attribute resource", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Identity Attribute resource", map[string]any{
		"name": plan.Name.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *identityAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state identityAttributeModel
	tflog.Debug(ctx, "Getting state for identity attribute resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the identity attribute from SailPoint
	tflog.Debug(ctx, "Fetching identity attribute from SailPoint", map[string]any{
		"name": state.Name.ValueString(),
	})
	identityAttributeResponse, err := r.client.GetIdentityAttribute(ctx, state.Name.ValueString())
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Identity Attribute not found, removing from state", map[string]any{
				"name": state.Name.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Attribute",
			fmt.Sprintf("Could not read SailPoint Identity Attribute %q: %s", state.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Identity Attribute", map[string]any{
			"name":  state.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if identityAttributeResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Attribute",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Identity Attribute API response to resource model", map[string]any{
		"name": state.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *identityAttributeResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for identity attribute resource", map[string]any{
		"name": state.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Identity Attribute resource", map[string]any{
		"name": state.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *identityAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identityAttributeModel
	tflog.Debug(ctx, "Getting state for identity attribute resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting identity attribute via SailPoint API", map[string]any{
		"name": state.Name.ValueString(),
	})
	err := r.client.DeleteIdentityAttribute(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Identity Attribute",
			fmt.Sprintf("Could not delete SailPoint Identity Attribute %q: %s", state.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Identity Attribute", map[string]any{
			"name":  state.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Identity Attribute resource", map[string]any{
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *identityAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identityAttributeModel
	tflog.Debug(ctx, "Getting plan for identity attribute resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping identity attribute resource model to API update request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiUpdateRequest, diags := plan.ToAPIUpdateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the identity attribute via the API client
	tflog.Debug(ctx, "Updating identity attribute via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	identityAttributeAPIResponse, err := r.client.UpdateIdentityAttribute(ctx, plan.Name.ValueString(), &apiUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Identity Attribute",
			fmt.Sprintf("Could not update SailPoint Identity Attribute %q: %s", plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Identity Attribute", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if identityAttributeAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Identity Attribute",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state identityAttributeModel
	tflog.Debug(ctx, "Mapping SailPoint Identity Attribute API response to resource model", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *identityAttributeAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for identity attribute resource", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Identity Attribute resource", map[string]any{
		"name": plan.Name.ValueString(),
	})
}

func (r *identityAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by setting the "name" attribute
	tflog.Debug(ctx, "Importing identity attribute resource", map[string]any{
		"name": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Identity Attribute resource", map[string]any{
		"name": req.ID,
	})
}
