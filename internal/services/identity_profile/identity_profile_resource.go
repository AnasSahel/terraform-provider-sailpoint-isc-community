// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_profile

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &identityProfileResource{}
	_ resource.ResourceWithConfigure   = &identityProfileResource{}
	_ resource.ResourceWithImportState = &identityProfileResource{}
)

type identityProfileResource struct {
	client *client.Client
}

// NewIdentityProfileResource creates a new resource for Identity Profile.
func NewIdentityProfileResource() resource.Resource {
	return &identityProfileResource{}
}

// Metadata implements resource.Resource.
func (r *identityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_profile"
}

// Configure implements resource.ResourceWithConfigure.
func (r *identityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "identity profile resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *identityProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Identity Profile.",
		MarkdownDescription: "Resource for SailPoint Identity Profile. Identity profiles define the source of identities and how identity attributes are mapped.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the identity profile.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity profile.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the identity profile.",
				Optional:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the identity profile.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner object. Must be `IDENTITY`.",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of the identity profile.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authoritative_source": schema.SingleNestedAttribute{
				MarkdownDescription: "The authoritative source for the identity profile. Changing this will recreate the resource.",
				Required:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the source object. Always `SOURCE`.",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the authoritative source.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the authoritative source.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"identity_attribute_config": schema.SingleNestedAttribute{
				MarkdownDescription: "The identity attribute configuration that defines how identity attributes are mapped.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the identity attribute configuration is enabled.",
						Optional:            true,
						Computed:            true,
					},
					"attribute_transforms": schema.ListNestedAttribute{
						MarkdownDescription: "List of identity attribute transforms.",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identity_attribute_name": schema.StringAttribute{
									MarkdownDescription: "The name of the identity attribute being mapped.",
									Required:            true,
								},
								"transform_definition": schema.SingleNestedAttribute{
									MarkdownDescription: "The transform definition for the identity attribute.",
									Required:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "The type of the transform definition (e.g., `accountAttribute`, `rule`).",
											Required:            true,
										},
										"attributes": schema.StringAttribute{
											MarkdownDescription: "The attributes of the transform definition as a JSON string.",
											Optional:            true,
											Computed:            true,
											CustomType:          jsontypes.NormalizedType{},
										},
									},
								},
							},
						},
					},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the identity profile was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the identity profile was last modified.",
				Computed:            true,
			},
		},
	}
}

// Create implements resource.Resource.
func (r *identityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identityProfileModel
	tflog.Debug(ctx, "Getting plan for identity profile resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping identity profile resource model to API create request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the identity profile via the API client
	tflog.Debug(ctx, "Creating identity profile via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiResponse, err := r.client.CreateIdentityProfile(ctx, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Identity Profile",
			fmt.Sprintf("Could not create SailPoint Identity Profile %q: %s",
				plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Identity Profile", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Identity Profile",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state identityProfileModel
	tflog.Debug(ctx, "Mapping SailPoint Identity Profile API response to resource model", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *apiResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Identity Profile resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *identityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identityProfileModel
	tflog.Debug(ctx, "Getting state for identity profile resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.ID.ValueString()

	// Read the identity profile from SailPoint
	tflog.Debug(ctx, "Fetching identity profile from SailPoint", map[string]any{
		"id": identityProfileID,
	})
	apiResponse, err := r.client.GetIdentityProfile(ctx, identityProfileID)
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Identity Profile not found, removing from state", map[string]any{
				"id": identityProfileID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Profile",
			fmt.Sprintf("Could not read SailPoint Identity Profile %q: %s",
				identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Identity Profile", map[string]any{
			"id":    identityProfileID,
			"error": err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Profile",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	resp.Diagnostics.Append(state.FromAPI(ctx, *apiResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Identity Profile resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *identityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identityProfileModel
	tflog.Debug(ctx, "Getting plan for identity profile resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the ID
	var state identityProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.ID.ValueString()

	// Build patch operations for changed fields
	tflog.Debug(ctx, "Building patch operations for identity profile update", map[string]any{
		"id": identityProfileID,
	})
	patchOperations, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(patchOperations) == 0 {
		tflog.Debug(ctx, "No changes detected, skipping update", map[string]any{
			"id": identityProfileID,
		})
		return
	}

	// Update the identity profile via the API client (PATCH)
	tflog.Debug(ctx, "Updating identity profile via SailPoint API", map[string]any{
		"id":               identityProfileID,
		"operations_count": len(patchOperations),
	})
	apiResponse, err := r.client.PatchIdentityProfile(ctx, identityProfileID, patchOperations)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Identity Profile",
			fmt.Sprintf("Could not update SailPoint Identity Profile %q: %s",
				identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Identity Profile", map[string]any{
			"id":    identityProfileID,
			"error": err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Identity Profile",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState identityProfileModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, *apiResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Identity Profile resource", map[string]any{
		"id":   identityProfileID,
		"name": newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *identityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identityProfileModel
	tflog.Debug(ctx, "Getting state for identity profile resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.ID.ValueString()

	tflog.Debug(ctx, "Deleting identity profile via SailPoint API", map[string]any{
		"id": identityProfileID,
	})
	_, err := r.client.DeleteIdentityProfile(ctx, identityProfileID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Identity Profile",
			fmt.Sprintf("Could not delete SailPoint Identity Profile %q: %s",
				identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Identity Profile", map[string]any{
			"id":    identityProfileID,
			"error": err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Identity Profile resource", map[string]any{
		"id":   identityProfileID,
		"name": state.Name.ValueString(),
	})
}

// ImportState implements resource.ResourceWithImportState.
func (r *identityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing identity profile resource", map[string]any{
		"import_id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Identity Profile resource", map[string]any{
		"id": req.ID,
	})
}
