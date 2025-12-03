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
	_ resource.Resource                = &workflowTriggerResource{}
	_ resource.ResourceWithConfigure   = &workflowTriggerResource{}
	_ resource.ResourceWithImportState = &workflowTriggerResource{}
)

type workflowTriggerResource struct {
	client *client.Client
}

func NewWorkflowTriggerResource() resource.Resource {
	return &workflowTriggerResource{}
}

func (r *workflowTriggerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workflowTriggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_trigger"
}

func (r *workflowTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.WorkflowTriggerSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource for managing a SailPoint Workflow trigger.",
		MarkdownDescription: "Manages a trigger for a SailPoint Workflow. This resource is used to attach triggers to workflows separately from the workflow definition, allowing for flexible trigger configuration. Triggers define what initiates a workflow execution.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *workflowTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Workflow Trigger resource")

	var plan models.WorkflowTriggerResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiTrigger, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Trigger",
			fmt.Sprintf("Could not convert workflow trigger: %s", err.Error()),
		)
		return
	}

	// Set the trigger on the workflow
	workflowID := plan.WorkflowID.ValueString()
	workflow, err := r.client.SetWorkflowTrigger(ctx, workflowID, apiTrigger)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Workflow Trigger",
			fmt.Sprintf("Could not set trigger on workflow %s: %s", workflowID, err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPoint(ctx, workflowID, workflow.Trigger); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Trigger Response",
			fmt.Sprintf("Could not convert workflow trigger response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Workflow Trigger resource created successfully", map[string]interface{}{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
}

func (r *workflowTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.WorkflowTriggerResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := state.WorkflowID.ValueString()

	// Get the workflow to retrieve the trigger
	workflow, err := r.client.GetWorkflow(ctx, workflowID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflow Trigger",
			fmt.Sprintf("Could not read workflow %s: %s", workflowID, err.Error()),
		)
		return
	}

	// Check if trigger exists
	if workflow.Trigger == nil {
		tflog.Warn(ctx, "Trigger not found on workflow, removing from state", map[string]interface{}{
			"workflow_id": workflowID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPoint(ctx, workflowID, workflow.Trigger); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Trigger Response",
			fmt.Sprintf("Could not convert workflow trigger response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *workflowTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Workflow Trigger resource")

	var plan models.WorkflowTriggerResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiTrigger, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Trigger",
			fmt.Sprintf("Could not convert workflow trigger: %s", err.Error()),
		)
		return
	}

	// Update the trigger on the workflow
	workflowID := plan.WorkflowID.ValueString()
	workflow, err := r.client.SetWorkflowTrigger(ctx, workflowID, apiTrigger)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workflow Trigger",
			fmt.Sprintf("Could not update trigger on workflow %s: %s", workflowID, err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPoint(ctx, workflowID, workflow.Trigger); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Trigger Response",
			fmt.Sprintf("Could not convert workflow trigger response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Workflow Trigger resource updated successfully", map[string]interface{}{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
}

func (r *workflowTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Workflow Trigger resource")

	var state models.WorkflowTriggerResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := state.WorkflowID.ValueString()

	// Remove the trigger from the workflow (set to null)
	_, err := r.client.RemoveWorkflowTrigger(ctx, workflowID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workflow Trigger",
			fmt.Sprintf("Could not remove trigger from workflow %s: %s", workflowID, err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Workflow Trigger resource deleted successfully", map[string]interface{}{
		"workflow_id": workflowID,
	})
}

func (r *workflowTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is the workflow_id
	resource.ImportStatePassthroughID(ctx, path.Root("workflow_id"), req, resp)
}
