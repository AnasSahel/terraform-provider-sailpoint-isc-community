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
	_ resource.Resource                = &identityProfileResource{}
	_ resource.ResourceWithConfigure   = &identityProfileResource{}
	_ resource.ResourceWithImportState = &identityProfileResource{}
)

type identityProfileResource struct {
	client *client.Client
}

func NewIdentityProfileResource() resource.Resource {
	return &identityProfileResource{}
}

func (r *identityProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *identityProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_profile"
}

func (r *identityProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.IdentityProfileSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Identity Profile.",
		MarkdownDescription: "Manages a SailPoint Identity Profile. Identity profiles define configurations for identities including authoritative sources and attribute mappings. See [Identity Profiles API](https://developer.sailpoint.com/docs/api/v2025/list-identity-profiles/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *identityProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating IdentityProfile resource")

	var plan models.IdentityProfile
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiProfile, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Profile",
			fmt.Sprintf("Could not convert identity profile: %s", err.Error()),
		)
		return
	}

	// Create the identity profile via API
	createdProfile, err := r.client.CreateIdentityProfile(ctx, apiProfile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Identity Profile",
			fmt.Sprintf("An error occurred while creating the identity profile: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Profile Response",
			fmt.Sprintf("Could not convert identity profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Profile resource created successfully", map[string]interface{}{
		"profile_id":   plan.ID.ValueString(),
		"profile_name": plan.Name.ValueString(),
	})
}

func (r *identityProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.IdentityProfile
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity profile via API
	fetchedProfile, err := r.client.GetIdentityProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Profile",
			fmt.Sprintf("Could not read identity profile %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Profile Response",
			fmt.Sprintf("Could not convert identity profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identityProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Identity Profile resource")

	var plan, state models.IdentityProfile
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate JSON Patch operations
	operations := state.GeneratePatchOperations(ctx, &plan)

	if len(operations) == 0 {
		tflog.Info(ctx, "No changes detected, skipping update")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	tflog.Debug(ctx, "Generated patch operations", map[string]interface{}{
		"operations": operations,
	})

	// Update the identity profile via API using PATCH
	updatedProfile, err := r.client.UpdateIdentityProfile(ctx, state.ID.ValueString(), operations)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Identity Profile",
			fmt.Sprintf("An error occurred while updating the identity profile: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Profile Response",
			fmt.Sprintf("Could not convert identity profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Profile resource updated successfully", map[string]interface{}{
		"profile_id":   plan.ID.ValueString(),
		"profile_name": plan.Name.ValueString(),
	})
}

func (r *identityProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Identity Profile resource")

	var state models.IdentityProfile
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the identity profile via API
	err := r.client.DeleteIdentityProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Profile",
			fmt.Sprintf("Could not delete identity profile %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Identity Profile resource deleted successfully", map[string]interface{}{
		"profile_id":   state.ID.ValueString(),
		"profile_name": state.Name.ValueString(),
	})
}

func (r *identityProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Identity profiles use "id" as the identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
