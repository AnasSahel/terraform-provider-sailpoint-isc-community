// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &lifecycleStateResource{}
	_ resource.ResourceWithConfigure   = &lifecycleStateResource{}
	_ resource.ResourceWithImportState = &lifecycleStateResource{}
)

type lifecycleStateResource struct {
	client *client.Client
}

func NewLifecycleStateResource() resource.Resource {
	return &lifecycleStateResource{}
}

func (r *lifecycleStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *lifecycleStateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

func (r *lifecycleStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.LifecycleStateSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Lifecycle State.",
		MarkdownDescription: "Manages a SailPoint Lifecycle State within an Identity Profile. Lifecycle states define different stages in an identity's lifecycle (e.g., onboarding, active, termination) and associated actions. See [Lifecycle State Documentation](https://developer.sailpoint.com/docs/api/v2025/lifecycle-states/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *lifecycleStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating LifecycleState resource")

	var plan models.LifecycleState
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity_profile_id from the plan
	identityProfileID := plan.IdentityProfileID.ValueString()

	// Convert Terraform model to API model
	apiLifecycleState, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Lifecycle State",
			fmt.Sprintf("Could not convert lifecycle state: %s", err.Error()),
		)
		return
	}

	// Create the lifecycle state via API
	createdLifecycleState, err := r.client.CreateLifecycleState(ctx, identityProfileID, apiLifecycleState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Lifecycle State",
			fmt.Sprintf("An error occurred while creating the lifecycle state: %s", err.Error()),
		)
		return
	}

	// Preserve the plan's optional nested objects before converting
	planHadEmailNotification := plan.EmailNotificationOption != nil
	planHadAccessActionConfig := plan.AccessActionConfiguration != nil

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdLifecycleState, identityProfileID); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Lifecycle State Response",
			fmt.Sprintf("Could not convert lifecycle state response: %s", err.Error()),
		)
		return
	}

	// If user didn't configure these optional nested objects, keep them as nil
	// to avoid "inconsistent result" errors from Terraform
	if !planHadEmailNotification {
		plan.EmailNotificationOption = nil
	}
	if !planHadAccessActionConfig {
		plan.AccessActionConfiguration = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Lifecycle State resource created successfully", map[string]interface{}{
		"lifecycle_state_id":  plan.ID.ValueString(),
		"identity_profile_id": identityProfileID,
	})
}

func (r *lifecycleStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.LifecycleState
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity_profile_id and lifecycle_state_id from state
	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	// Preserve the state's optional nested objects before fetching
	stateHadEmailNotification := state.EmailNotificationOption != nil
	stateHadAccessActionConfig := state.AccessActionConfiguration != nil

	// Get the lifecycle state via API
	fetchedLifecycleState, err := r.client.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lifecycle State",
			fmt.Sprintf("Could not read lifecycle state ID %s for identity profile %s: %s", lifecycleStateID, identityProfileID, err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedLifecycleState, identityProfileID); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Lifecycle State Response",
			fmt.Sprintf("Could not convert lifecycle state response: %s", err.Error()),
		)
		return
	}

	// Preserve nil for optional nested objects if they weren't in the original state
	// This prevents configuration drift when API returns defaults
	if !stateHadEmailNotification {
		state.EmailNotificationOption = nil
	}
	if !stateHadAccessActionConfig {
		state.AccessActionConfiguration = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *lifecycleStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating LifecycleState resource")

	var plan models.LifecycleState
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.LifecycleState
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity_profile_id and lifecycle_state_id from state
	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	// Generate patch operations
	patchOperations, err := state.GeneratePatchOperations(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Patch Operations",
			fmt.Sprintf("Could not generate patch operations for lifecycle state: %s", err.Error()),
		)
		return
	}

	// Debug: log the patch operations
	tflog.Debug(ctx, "Generated patch operations", map[string]interface{}{
		"operations": fmt.Sprintf("%+v", patchOperations),
	})

	// If no operations, skip the update
	if len(patchOperations) == 0 {
		tflog.Info(ctx, "No changes detected, skipping update", map[string]interface{}{
			"lifecycle_state_id":  lifecycleStateID,
			"identity_profile_id": identityProfileID,
		})
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Preserve the plan's optional nested objects before updating
	planHadEmailNotification := plan.EmailNotificationOption != nil
	planHadAccessActionConfig := plan.AccessActionConfiguration != nil

	// Update the lifecycle state via API using PATCH
	updatedLifecycleState, err := r.client.PatchLifecycleState(ctx, identityProfileID, lifecycleStateID, patchOperations)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Lifecycle State",
			fmt.Sprintf("An error occurred while updating the lifecycle state: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedLifecycleState, identityProfileID); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Lifecycle State Response",
			fmt.Sprintf("Could not convert lifecycle state response: %s", err.Error()),
		)
		return
	}

	// If user didn't configure these optional nested objects, keep them as nil
	if !planHadEmailNotification {
		plan.EmailNotificationOption = nil
	}
	if !planHadAccessActionConfig {
		plan.AccessActionConfiguration = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Lifecycle State resource updated successfully", map[string]interface{}{
		"lifecycle_state_id":  plan.ID.ValueString(),
		"identity_profile_id": identityProfileID,
		"operations_count":    len(patchOperations),
	})
}

func (r *lifecycleStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting LifecycleState resource")

	var state models.LifecycleState
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity_profile_id and lifecycle_state_id from state
	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	// Delete the lifecycle state via API
	err := r.client.DeleteLifecycleState(ctx, identityProfileID, lifecycleStateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Lifecycle State",
			fmt.Sprintf("Could not delete lifecycle state ID %s for identity profile %s: %s", lifecycleStateID, identityProfileID, err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Lifecycle State resource deleted successfully", map[string]interface{}{
		"lifecycle_state_id":  lifecycleStateID,
		"identity_profile_id": identityProfileID,
	})
}

func (r *lifecycleStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: identity_profile_id:lifecycle_state_id
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in the format 'identity_profile_id:lifecycle_state_id', got: %s", req.ID),
		)
		return
	}

	identityProfileID := parts[0]
	lifecycleStateID := parts[1]

	// Validate both parts are not empty
	if identityProfileID == "" || lifecycleStateID == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Both identity_profile_id and lifecycle_state_id must be non-empty. Got: %s", req.ID),
		)
		return
	}

	// Set the identity_profile_id and id in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identity_profile_id"), identityProfileID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lifecycleStateID)...)

	tflog.Info(ctx, "Lifecycle State resource imported successfully", map[string]interface{}{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
}
