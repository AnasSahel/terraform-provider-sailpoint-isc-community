// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func NewAccessProfileResource() resource.Resource {
	return &accessProfileResource{}
}

func (r *accessProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *accessProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_profile"
}

func (r *accessProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.AccessProfileSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Access Profile.",
		MarkdownDescription: "Manages a SailPoint Access Profile. Access Profiles are collections of entitlements from a source that can be requested by users. See [Access Profile Documentation](https://developer.sailpoint.com/docs/api/v2025/create-access-profile) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *accessProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Access Profile resource")

	var plan models.AccessProfile
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiAccessProfile, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Access Profile",
			fmt.Sprintf("Could not convert access profile: %s", err.Error()),
		)
		return
	}

	// Create the access profile via API
	createdAccessProfile, err := r.client.CreateAccessProfile(ctx, apiAccessProfile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Access Profile",
			fmt.Sprintf("An error occurred while creating the access profile: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdAccessProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Access Profile Response",
			fmt.Sprintf("Could not convert access profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access Profile resource created successfully", map[string]interface{}{
		"access_profile_id": plan.ID.ValueString(),
	})
}

func (r *accessProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccessProfile
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the access profile via API
	fetchedAccessProfile, err := r.client.GetAccessProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Access Profile",
			fmt.Sprintf("Could not read access profile ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedAccessProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Access Profile Response",
			fmt.Sprintf("Could not convert access profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *accessProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Access Profile resource")

	var plan models.AccessProfile
	var state models.AccessProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate patch operations
	operations, err := state.GeneratePatchOperations(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Patch Operations",
			fmt.Sprintf("Could not generate patch operations: %s", err.Error()),
		)
		return
	}

	// Update the access profile via API if there are operations
	if len(operations) > 0 {
		_, err = r.client.PatchAccessProfile(ctx, plan.ID.ValueString(), operations)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Access Profile",
				fmt.Sprintf("An error occurred while updating the access profile: %s", err.Error()),
			)
			return
		}
	} else {
		tflog.Info(ctx, "No changes detected, skipping PATCH operation")
	}

	// Always read back the access profile to get the latest state including timestamps
	updatedAccessProfile, err := r.client.GetAccessProfile(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Access Profile",
			fmt.Sprintf("Could not read access profile after update: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedAccessProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Access Profile Response",
			fmt.Sprintf("Could not convert access profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access Profile resource updated successfully", map[string]interface{}{
		"access_profile_id": plan.ID.ValueString(),
	})
}

func (r *accessProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Access Profile resource")

	var state models.AccessProfile
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the access profile via API
	err := r.client.DeleteAccessProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Access Profile",
			fmt.Sprintf("Could not delete access profile ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Access Profile resource deleted successfully", map[string]interface{}{
		"access_profile_id": state.ID.ValueString(),
	})
}

func (r *accessProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
