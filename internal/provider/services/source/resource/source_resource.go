// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type SourceResource struct {
	client *api_v2025.APIClient
}

var (
	_ resource.Resource                = &SourceResource{}
	_ resource.ResourceWithConfigure   = &SourceResource{}
	_ resource.ResourceWithImportState = &SourceResource{}
)

func NewSourceResource() resource.Resource {
	return &SourceResource{}
}

func (r *SourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring SourceResource")

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

	tflog.Debug(ctx, "Configured SourceResource")
	r.client = client
}

func (r *SourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *SourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Trace(ctx, "Preparing SourceResource schema")
	resp.Schema = GetSourceResourceSchema()
}

func (r *SourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Source")

	var plan SourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to SailPoint API request
	sourceRequest, diags := plan.ToSailPointCreateSourceRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Sending create request to SailPoint API", map[string]interface{}{
		"name":      sourceRequest.GetName(),
		"connector": sourceRequest.GetConnector(),
	})

	// Create the source via SailPoint API
	source, httpResponse, err := r.client.SourcesAPI.CreateSource(
		context.Background(),
	).Source(*sourceRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			fmt.Sprintf("Could not create source '%s': %s\n\nHTTP Response: %v",
				plan.Name.ValueString(), err.Error(), httpResponse),
		)
		return
	}

	// Convert API response to Terraform state
	var state SourceResourceModel
	conversionDiags := state.FromSailPointSource(ctx, source)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully created source", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": state.Name.ValueString(),
	})

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Source")

	var state SourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to get current state for source")
		return
	}

	sourceID := state.Id.ValueString()
	tflog.Debug(ctx, "Fetching source from SailPoint API", map[string]interface{}{
		"id": sourceID,
	})

	// Fetch the source from SailPoint API
	source, httpResponse, err := r.client.SourcesAPI.GetSource(
		context.Background(),
		sourceID,
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			fmt.Sprintf("Could not read source ID '%s': %s\n\nHTTP Response: %v",
				sourceID, err.Error(), httpResponse),
		)
		return
	}

	// Convert API response to Terraform state
	conversionDiags := state.FromSailPointSource(ctx, source)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Successfully read source", map[string]interface{}{
		"id":     state.Id.ValueString(),
		"name":   state.Name.ValueString(),
		"status": state.Status.ValueString(),
	})

	// Set the refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Source")

	// Get the current state
	var state SourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the planned changes
	var plan SourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan to source update request
	sourceUpdateRequest, diags := plan.ToSailPointCreateSourceRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID for update
	sourceUpdateRequest.SetId(state.Id.ValueString())

	tflog.Debug(ctx, "Sending update request to SailPoint API", map[string]interface{}{
		"id":   state.Id.ValueString(),
		"name": sourceUpdateRequest.GetName(),
	})

	// Update the source via SailPoint API
	source, httpResponse, err := r.client.SourcesAPI.PutSource(
		context.Background(),
		state.Id.ValueString(),
	).Source(*sourceUpdateRequest).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			fmt.Sprintf("Could not update source %s: %s\nHTTP Response: %v",
				state.Id.ValueString(), err, httpResponse),
		)
		return
	}

	// Convert API response to Terraform state
	var newState SourceResourceModel
	conversionDiags := newState.FromSailPointSource(ctx, source)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully updated source", map[string]interface{}{
		"id":   newState.Id.ValueString(),
		"name": newState.Name.ValueString(),
	})

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *SourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Source")

	var state SourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.Id.ValueString()
	sourceName := state.Name.ValueString()

	tflog.Debug(ctx, "Sending delete request to SailPoint API", map[string]interface{}{
		"id":   sourceID,
		"name": sourceName,
	})

	// Delete the source via SailPoint API
	deleteResp, httpResponse, err := r.client.SourcesAPI.DeleteSource(
		context.Background(),
		sourceID,
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			fmt.Sprintf("Could not delete source '%s' (ID: %s): %s\n\nHTTP Response: %v",
				sourceName, sourceID, err.Error(), httpResponse),
		)
		return
	}

	// Log delete response if available
	if deleteResp != nil {
		tflog.Debug(ctx, "Delete response received", map[string]interface{}{
			"sourceId": sourceID,
		})
	}

	tflog.Info(ctx, "Successfully deleted source", map[string]interface{}{
		"id":   sourceID,
		"name": sourceName,
	})
}

// ImportState enables importing existing sources by ID.
func (r *SourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
