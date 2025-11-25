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
	_ resource.Resource                = &launcherResource{}
	_ resource.ResourceWithConfigure   = &launcherResource{}
	_ resource.ResourceWithImportState = &launcherResource{}
)

type launcherResource struct {
	client *client.Client
}

func NewLauncherResource() resource.Resource {
	return &launcherResource{}
}

func (r *launcherResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *launcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launcher"
}

func (r *launcherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.LauncherSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Launcher.",
		MarkdownDescription: "Manages a SailPoint Launcher. Launchers are interactive processes that can be triggered to execute workflows or other automation tasks. See [Launcher Documentation](https://developer.sailpoint.com/docs/api/v2025/create-launcher) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *launcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Launcher resource")

	var plan models.Launcher
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiLauncher := plan.ConvertToSailPoint(ctx)

	// Create the launcher via API
	createdLauncher, err := r.client.CreateLauncher(ctx, apiLauncher)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Launcher",
			fmt.Sprintf("An error occurred while creating the launcher: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	plan.ConvertFromSailPointForResource(ctx, createdLauncher)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launcher resource created successfully", map[string]interface{}{
		"launcher_id": plan.ID.ValueString(),
	})
}

func (r *launcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Launcher
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the launcher via API
	fetchedLauncher, err := r.client.GetLauncher(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Launcher",
			fmt.Sprintf("Could not read launcher ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	state.ConvertFromSailPointForResource(ctx, fetchedLauncher)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *launcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Launcher resource")

	var plan models.Launcher
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiLauncher := plan.ConvertToSailPoint(ctx)

	// Update the launcher via API using PUT (full update)
	updatedLauncher, err := r.client.UpdateLauncher(ctx, plan.ID.ValueString(), apiLauncher)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Launcher",
			fmt.Sprintf("An error occurred while updating the launcher: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	plan.ConvertFromSailPointForResource(ctx, updatedLauncher)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launcher resource updated successfully", map[string]interface{}{
		"launcher_id": plan.ID.ValueString(),
	})
}

func (r *launcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Launcher resource")

	var state models.Launcher
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the launcher via API
	err := r.client.DeleteLauncher(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Launcher",
			fmt.Sprintf("Could not delete launcher ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Launcher resource deleted successfully", map[string]interface{}{
		"launcher_id": state.ID.ValueString(),
	})
}

func (r *launcherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
