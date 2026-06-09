// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// Metadata implements resource.Resource.
func (r *workflowTriggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_trigger"
}

// Configure implements resource.ResourceWithConfigure.
func (r *workflowTriggerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "workflow trigger resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *workflowTriggerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a trigger for a SailPoint Workflow.",
		MarkdownDescription: "Manages a trigger for a SailPoint Workflow. This resource is used to attach triggers to workflows separately from the workflow definition, allowing for flexible trigger configuration. Triggers define what initiates a workflow execution.",
		Attributes: map[string]schema.Attribute{
			"workflow_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the workflow to attach this trigger to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of trigger. Valid values are `EVENT`, `EXTERNAL`, or `SCHEDULED`.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the trigger.",
				Optional:            true,
			},
			"attributes": schema.StringAttribute{
				MarkdownDescription: "JSON object containing trigger-specific attributes. For EVENT triggers, this includes the event type (`id`, `filter`). For SCHEDULED triggers, this includes `cronString` and `frequency`.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *workflowTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workflowTriggerModel
	tflog.Debug(ctx, "Getting plan for workflow trigger resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := plan.WorkflowID.ValueString()
	tflog.Debug(ctx, "Mapping workflow trigger resource model to API request", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
	apiTrigger, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The ISC API rejects PATCH on enabled workflows; disable → patch → re-enable.
	var workflow *client.WorkflowAPI
	fnErr, reEnableErr := withWorkflowDisabled(ctx, r.client, workflowID, func() error {
		var err error
		tflog.Debug(ctx, "Setting workflow trigger via SailPoint API", map[string]any{
			"workflow_id":  workflowID,
			"trigger_type": plan.Type.ValueString(),
		})
		workflow, err = r.client.SetWorkflowTrigger(ctx, workflowID, apiTrigger)
		return err
	})

	if reEnableErr != nil {
		resp.Diagnostics.AddError(
			"Error Re-enabling Workflow After Trigger Create",
			fmt.Sprintf("Could not re-enable workflow %q after setting trigger: %s", workflowID, reEnableErr.Error()),
		)
	}
	if fnErr != nil {
		resp.Diagnostics.AddError(
			"Error Creating Workflow Trigger",
			fmt.Sprintf("Could not set trigger on workflow %q: %s", workflowID, fnErr.Error()),
		)
		tflog.Error(ctx, "Failed to create workflow trigger", map[string]any{
			"workflow_id": workflowID,
			"error":       fnErr.Error(),
		})
		return
	}

	if workflow == nil || workflow.Trigger == nil {
		resp.Diagnostics.AddError(
			"Error Creating Workflow Trigger",
			"Received nil response or trigger from SailPoint API",
		)
		return
	}

	var state workflowTriggerModel
	tflog.Debug(ctx, "Mapping SailPoint API response to resource model", map[string]any{
		"workflow_id": workflowID,
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, workflowID, workflow.Trigger)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// attributes is user-owned; the API may inject null-valued keys (e.g.
	// integrationId). Preserve the configured value so server-added nulls don't
	// surface as "inconsistent result after apply". (#136)
	if !plan.Attributes.IsNull() && !plan.Attributes.IsUnknown() {
		state.Attributes = plan.Attributes
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created workflow trigger resource", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *workflowTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowTriggerModel
	tflog.Debug(ctx, "Getting state for workflow trigger resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := state.WorkflowID.ValueString()

	tflog.Debug(ctx, "Fetching workflow from SailPoint", map[string]any{
		"workflow_id": workflowID,
	})
	workflow, err := r.client.GetWorkflow(ctx, workflowID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflow Trigger",
			fmt.Sprintf("Could not read workflow %q: %s", workflowID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read workflow", map[string]any{
			"workflow_id": workflowID,
			"error":       err.Error(),
		})
		return
	}

	if workflow.Trigger == nil || workflow.Trigger.Type == "" {
		tflog.Warn(ctx, "Trigger not found on workflow, removing from state", map[string]any{
			"workflow_id": workflowID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Debug(ctx, "Mapping SailPoint API response to resource model", map[string]any{
		"workflow_id": workflowID,
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, workflowID, workflow.Trigger)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Drop server-injected top-level null keys (e.g. integrationId) so a refresh
	// of an unchanged trigger doesn't report drift. (#136)
	state.Attributes = stripTopLevelNullKeys(state.Attributes)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read workflow trigger resource", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": state.Type.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *workflowTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workflowTriggerModel
	tflog.Debug(ctx, "Getting plan for workflow trigger resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := plan.WorkflowID.ValueString()
	tflog.Debug(ctx, "Mapping workflow trigger resource model to API request", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
	apiTrigger, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The ISC API rejects PATCH on enabled workflows; disable → patch → re-enable.
	var workflow *client.WorkflowAPI
	fnErr, reEnableErr := withWorkflowDisabled(ctx, r.client, workflowID, func() error {
		var err error
		tflog.Debug(ctx, "Updating workflow trigger via SailPoint API", map[string]any{
			"workflow_id":  workflowID,
			"trigger_type": plan.Type.ValueString(),
		})
		workflow, err = r.client.SetWorkflowTrigger(ctx, workflowID, apiTrigger)
		return err
	})

	if reEnableErr != nil {
		resp.Diagnostics.AddError(
			"Error Re-enabling Workflow After Trigger Update",
			fmt.Sprintf("Could not re-enable workflow %q after updating trigger: %s", workflowID, reEnableErr.Error()),
		)
	}
	if fnErr != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workflow Trigger",
			fmt.Sprintf("Could not update trigger on workflow %q: %s", workflowID, fnErr.Error()),
		)
		tflog.Error(ctx, "Failed to update workflow trigger", map[string]any{
			"workflow_id": workflowID,
			"error":       fnErr.Error(),
		})
		return
	}

	if workflow == nil || workflow.Trigger == nil {
		resp.Diagnostics.AddError(
			"Error Updating Workflow Trigger",
			"Received nil response or trigger from SailPoint API",
		)
		return
	}

	var state workflowTriggerModel
	tflog.Debug(ctx, "Mapping SailPoint API response to resource model", map[string]any{
		"workflow_id": workflowID,
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, workflowID, workflow.Trigger)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the configured attributes (see #136 note in Create).
	if !plan.Attributes.IsNull() && !plan.Attributes.IsUnknown() {
		state.Attributes = plan.Attributes
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated workflow trigger resource", map[string]any{
		"workflow_id":  workflowID,
		"trigger_type": plan.Type.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *workflowTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workflowTriggerModel
	tflog.Debug(ctx, "Getting state for workflow trigger resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowID := state.WorkflowID.ValueString()

	// The ISC API rejects PATCH on enabled workflows; disable → patch → re-enable.
	fnErr, reEnableErr := withWorkflowDisabled(ctx, r.client, workflowID, func() error {
		tflog.Debug(ctx, "Removing workflow trigger via SailPoint API", map[string]any{
			"workflow_id":  workflowID,
			"trigger_type": state.Type.ValueString(),
		})
		_, err := r.client.RemoveWorkflowTrigger(ctx, workflowID)
		return err
	})

	if reEnableErr != nil {
		resp.Diagnostics.AddError(
			"Error Re-enabling Workflow After Trigger Delete",
			fmt.Sprintf("Could not re-enable workflow %q after removing trigger: %s", workflowID, reEnableErr.Error()),
		)
	}
	if fnErr != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workflow Trigger",
			fmt.Sprintf("Could not remove trigger from workflow %q: %s", workflowID, fnErr.Error()),
		)
		tflog.Error(ctx, "Failed to remove workflow trigger", map[string]any{
			"workflow_id": workflowID,
			"error":       fnErr.Error(),
		})
		return
	}

	tflog.Info(ctx, "Successfully deleted workflow trigger resource", map[string]any{
		"workflow_id": workflowID,
	})
}

// ImportState implements resource.ResourceWithImportState.
func (r *workflowTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing workflow trigger resource", map[string]any{
		"workflow_id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("workflow_id"), req, resp)

	tflog.Info(ctx, "Successfully imported workflow trigger resource", map[string]any{
		"workflow_id": req.ID,
	})
}
