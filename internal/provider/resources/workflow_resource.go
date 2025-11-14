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
	_ resource.Resource                = &workflowResource{}
	_ resource.ResourceWithConfigure   = &workflowResource{}
	_ resource.ResourceWithImportState = &workflowResource{}
)

type workflowResource struct {
	client *client.Client
}

func NewWorkflowResource() resource.Resource {
	return &workflowResource{}
}

func (r *workflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *workflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.WorkflowSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Workflow.",
		MarkdownDescription: "Manages a SailPoint Workflow. Workflows are custom automation scripts that respond to event triggers and perform a series of actions. See [Workflow Documentation](https://developer.sailpoint.com/docs/extensibility/workflows/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Workflow resource")

	var plan models.Workflow
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiWorkflow, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow",
			fmt.Sprintf("Could not convert workflow: %s", err.Error()),
		)
		return
	}

	// Create the workflow via API
	createdWorkflow, err := r.client.CreateWorkflow(ctx, apiWorkflow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Workflow",
			fmt.Sprintf("An error occurred while creating the workflow: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdWorkflow); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Response",
			fmt.Sprintf("Could not convert workflow response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Workflow resource created successfully", map[string]interface{}{
		"workflow_id": plan.ID.ValueString(),
	})
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Workflow
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the workflow via API
	fetchedWorkflow, err := r.client.GetWorkflow(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedWorkflow); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Response",
			fmt.Sprintf("Could not convert workflow response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Workflow resource")

	var plan models.Workflow
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiWorkflow, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow",
			fmt.Sprintf("Could not convert workflow: %s", err.Error()),
		)
		return
	}

	// Update the workflow via API using PUT (full update)
	updatedWorkflow, err := r.client.UpdateWorkflow(ctx, plan.ID.ValueString(), apiWorkflow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workflow",
			fmt.Sprintf("An error occurred while updating the workflow: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedWorkflow); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Response",
			fmt.Sprintf("Could not convert workflow response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Workflow resource updated successfully", map[string]interface{}{
		"workflow_id": plan.ID.ValueString(),
	})
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Workflow resource")

	var state models.Workflow
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: Workflows must be disabled before deletion
	// Check if workflow is enabled and disable it first if needed
	if !state.Enabled.IsNull() && state.Enabled.ValueBool() {
		tflog.Info(ctx, "Workflow is enabled, disabling before deletion", map[string]interface{}{
			"workflow_id": state.ID.ValueString(),
		})

		// Convert to API model and disable
		apiWorkflow, err := state.ConvertToSailPoint(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Converting Workflow",
				fmt.Sprintf("Could not convert workflow for disabling: %s", err.Error()),
			)
			return
		}

		falseValue := false
		apiWorkflow.Enabled = &falseValue

		// Update to disable the workflow
		_, err = r.client.UpdateWorkflow(ctx, state.ID.ValueString(), apiWorkflow)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Disabling Workflow",
				fmt.Sprintf("Could not disable workflow ID %s before deletion: %s", state.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Info(ctx, "Workflow disabled successfully before deletion", map[string]interface{}{
			"workflow_id": state.ID.ValueString(),
		})
	}

	// Delete the workflow via API
	err := r.client.DeleteWorkflow(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workflow",
			fmt.Sprintf("Could not delete workflow ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Workflow resource deleted successfully", map[string]interface{}{
		"workflow_id": state.ID.ValueString(),
	})
}

func (r *workflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
