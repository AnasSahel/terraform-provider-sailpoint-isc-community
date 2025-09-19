// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type ManagedClusterResource struct {
	client *api_v2025.APIClient
}

var (
	_ resource.Resource                = &ManagedClusterResource{}
	_ resource.ResourceWithConfigure   = &ManagedClusterResource{}
	_ resource.ResourceWithImportState = &ManagedClusterResource{}
)

func NewManagedClusterResource() resource.Resource {
	return &ManagedClusterResource{}
}

func (r *ManagedClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring ManagedClusterResource")

	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	tflog.Debug(ctx, "Configured ManagedClusterResource")
	r.client = client
}

func (r *ManagedClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

func (r *ManagedClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Trace(ctx, "Preparing ManagedClusterResource schema")
	resp.Schema = GetManagedClusterResourceSchema()
}

func (r *ManagedClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Managed Cluster")

	var plan ManagedClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required fields before making API call
	validationDiags := ValidateRequiredFields(&plan)
	resp.Diagnostics.Append(validationDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to SailPoint API request
	managedClusterRequest, diags := plan.ToSailPointCreateManagedClusterRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Sending create request to SailPoint API", map[string]interface{}{
		"name": managedClusterRequest.GetName(),
		"type": managedClusterRequest.GetType(),
	})

	// Create the managed cluster via SailPoint API
	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.CreateManagedCluster(
		context.Background(),
	).ManagedClusterRequest(*managedClusterRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Managed Cluster",
			fmt.Sprintf("Could not create managed cluster '%s': %s\n\nHTTP Response: %v",
				plan.Name.ValueString(), err.Error(), httpResponse),
		)
		return
	}

	// Convert API response to Terraform state
	var state ManagedClusterResourceModel
	conversionDiags := state.FromSailPointManagedCluster(ctx, managedCluster)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the original planned configuration to avoid inconsistency errors
	if !plan.Configuration.IsNull() {
		state.Configuration = plan.Configuration
	}

	tflog.Info(ctx, "Successfully created managed cluster", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ManagedClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Managed Cluster")

	var state ManagedClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to get current state for managed cluster")
		return
	}

	clusterID := state.Id.ValueString()
	tflog.Debug(ctx, "Fetching managed cluster from SailPoint API", map[string]interface{}{
		"id": clusterID,
	})

	// Fetch the managed cluster from SailPoint API
	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.GetManagedCluster(
		context.Background(),
		clusterID,
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Managed Cluster",
			fmt.Sprintf("Could not read managed cluster ID '%s': %s\n\nHTTP Response: %v",
				clusterID, err.Error(), httpResponse),
		)
		return
	}

	// Convert API response to Terraform state
	conversionDiags := state.FromSailPointManagedCluster(ctx, managedCluster)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Successfully read managed cluster", map[string]interface{}{
		"id":     state.Id.ValueString(),
		"name":   state.Name.ValueString(),
		"status": state.Status.ValueString(),
	})

	// Set the refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ManagedClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Managed Cluster")

	// Get the current state
	var state ManagedClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the planned changes
	var plan ManagedClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build JSON Patch operations using the extracted function
	patchOps := BuildManagedClusterPatches(&state, &plan)

	// If no changes, return early
	if len(patchOps) == 0 {
		tflog.Debug(ctx, "No changes detected, skipping API call")
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Applying %d patch operations to managed cluster %s", len(patchOps), state.Id.ValueString()))

	// Call the SailPoint API to update the managed cluster
	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.UpdateManagedCluster(
		context.Background(),
		state.Id.ValueString(),
	).JsonPatchOperation(patchOps).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Managed Cluster",
			fmt.Sprintf("Could not update managed cluster %s: %s\nHTTP Response: %v", state.Id.ValueString(), err, httpResponse),
		)
		return
	}

	// Use the selective update method to preserve computed fields
	newState := state
	updateDiags := newState.UpdateSelectiveFields(ctx, managedCluster, &plan)
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ManagedClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Managed Cluster")

	var state ManagedClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.Id.ValueString()
	clusterName := state.Name.ValueString()

	tflog.Debug(ctx, "Sending delete request to SailPoint API", map[string]interface{}{
		"id":   clusterID,
		"name": clusterName,
	})

	// Delete the managed cluster via SailPoint API
	httpResponse, err := r.client.ManagedClustersAPI.DeleteManagedCluster(
		context.Background(),
		clusterID,
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Managed Cluster",
			fmt.Sprintf("Could not delete managed cluster '%s' (ID: %s): %s\n\nHTTP Response: %v",
				clusterName, clusterID, err.Error(), httpResponse),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted managed cluster", map[string]interface{}{
		"id":   clusterID,
		"name": clusterName,
	})
}

// ImportState enables importing existing managed clusters by ID.
func (r *ManagedClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
