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
	_ resource.Resource                = &formDefinitionResource{}
	_ resource.ResourceWithConfigure   = &formDefinitionResource{}
	_ resource.ResourceWithImportState = &formDefinitionResource{}
)

type formDefinitionResource struct {
	client *client.Client
}

func NewFormDefinitionResource() resource.Resource {
	return &formDefinitionResource{}
}

func (r *formDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *formDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

func (r *formDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.FormDefinitionSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Form Definition.",
		MarkdownDescription: "Manages a SailPoint Form Definition. Forms are composed of sections and fields for data collection in workflows. See [Custom Forms Documentation](https://developer.sailpoint.com/docs/api/v2025/custom-forms/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *formDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating FormDefinition resource")

	var plan models.FormDefinition
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiForm, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Form Definition",
			fmt.Sprintf("Could not convert form definition: %s", err.Error()),
		)
		return
	}

	// Create the form definition via API
	createdForm, err := r.client.CreateFormDefinition(ctx, apiForm)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Form Definition",
			fmt.Sprintf("An error occurred while creating the form definition: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdForm); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Form Definition Response",
			fmt.Sprintf("Could not convert form definition response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Form Definition resource created successfully", map[string]interface{}{
		"form_id": plan.ID.ValueString(),
	})
}

func (r *formDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.FormDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the form definition via API
	fetchedForm, err := r.client.GetFormDefinition(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Form Definition",
			fmt.Sprintf("Could not read form definition ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedForm); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Form Definition Response",
			fmt.Sprintf("Could not convert form definition response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *formDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Form Definition resource")

	var plan models.FormDefinition
	var state models.FormDefinition

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

	// If no operations, nothing to update
	if len(operations) == 0 {
		tflog.Info(ctx, "No changes detected, skipping update")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Update the form definition via API
	_, err = r.client.PatchFormDefinition(ctx, plan.ID.ValueString(), operations)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Form Definition",
			fmt.Sprintf("An error occurred while updating the form definition: %s", err.Error()),
		)
		return
	}

	// Read back the updated form to get the latest state including timestamps
	updatedForm, err := r.client.GetFormDefinition(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Form Definition",
			fmt.Sprintf("Could not read form definition after update: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedForm); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Form Definition Response",
			fmt.Sprintf("Could not convert form definition response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Form Definition resource updated successfully", map[string]interface{}{
		"form_id": plan.ID.ValueString(),
	})
}

func (r *formDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Form Definition resource")

	var state models.FormDefinition
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the form definition via API
	err := r.client.DeleteFormDefinition(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Form Definition",
			fmt.Sprintf("Could not delete form definition ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Form Definition resource deleted successfully", map[string]interface{}{
		"form_id": state.ID.ValueString(),
	})
}

func (r *formDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
