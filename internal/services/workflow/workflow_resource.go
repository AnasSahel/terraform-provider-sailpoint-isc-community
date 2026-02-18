// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// Metadata implements resource.Resource.
func (r *workflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Configure implements resource.ResourceWithConfigure.
func (r *workflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "workflow resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *workflowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a SailPoint Workflow.",
		MarkdownDescription: "Manages a SailPoint Workflow. Workflows are custom automation scripts that respond to event triggers and perform a series of actions. The trigger is managed separately using the `sailpoint_workflow_trigger` resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the workflow.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workflow.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the workflow.",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is enabled. Workflows cannot be created in an enabled state. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"trigger": schema.StringAttribute{
				MarkdownDescription: "The trigger configuration as JSON. This is a computed field - use `sailpoint_workflow_trigger` resource to manage triggers.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the workflow was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the workflow was last modified.",
				Computed:            true,
			},
			"creator": schema.SingleNestedAttribute{
				MarkdownDescription: "The identity who created the workflow.",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the creator (e.g., `IDENTITY`).",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the creator.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the creator.",
						Computed:            true,
					},
				},
			},
			"modified_by": schema.SingleNestedAttribute{
				MarkdownDescription: "The identity who last modified the workflow.",
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the modifier (e.g., `IDENTITY`).",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the modifier.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the modifier.",
						Computed:            true,
					},
				},
			},
			"execution_count": schema.Int32Attribute{
				MarkdownDescription: "The number of times the workflow has been executed.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"failure_count": schema.Int32Attribute{
				MarkdownDescription: "The number of times the workflow has failed.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the workflow.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner (e.g., `IDENTITY`).",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"definition": schema.SingleNestedAttribute{
				MarkdownDescription: "The workflow definition containing the steps to execute. If not specified, the workflow will have no definition.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"start": schema.StringAttribute{
						MarkdownDescription: "The name of the starting step.",
						Required:            true,
					},
					"steps": schema.StringAttribute{
						MarkdownDescription: "JSON object containing the workflow steps. Each key is a step name and the value defines the step configuration including action type, attributes, and flow control.\n\n" +
							"~> **Note:** When configuring steps that use secrets (e.g., OAuth client secrets for `sp:http` actions), " +
							"set the secret value through the SailPoint UI first, then copy the resulting vault reference " +
							"(e.g., `$.secrets.<uid>`) into your Terraform configuration to prevent drift.",
						Required:   true,
						CustomType: jsontypes.NormalizedType{},
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workflowModel
	tflog.Debug(ctx, "Getting plan for workflow resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping workflow resource model to API create request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the workflow via the API client
	tflog.Debug(ctx, "Creating workflow via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	workflowAPIResponse, err := r.client.CreateWorkflow(ctx, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Workflow",
			fmt.Sprintf("Could not create SailPoint Workflow %q: %s", plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Workflow", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if workflowAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Workflow",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state workflowModel
	tflog.Debug(ctx, "Mapping SailPoint Workflow API response to resource model", map[string]any{
		"id":   workflowAPIResponse.ID,
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *workflowAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for workflow resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Workflow resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowModel
	tflog.Debug(ctx, "Getting state for workflow resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the workflow from SailPoint
	tflog.Debug(ctx, "Fetching workflow from SailPoint", map[string]any{
		"id": state.ID.ValueString(),
	})
	workflowResponse, err := r.client.GetWorkflow(ctx, state.ID.ValueString())
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Workflow not found, removing from state", map[string]any{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Workflow",
			fmt.Sprintf("Could not read SailPoint Workflow %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Workflow", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if workflowResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Workflow",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Workflow API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *workflowResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for workflow resource", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Workflow resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workflowModel
	tflog.Debug(ctx, "Getting plan for workflow resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to retrieve the ID
	var state workflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping workflow resource model to API update request", map[string]any{
		"id":   state.ID.ValueString(),
		"name": plan.Name.ValueString(),
	})
	apiUpdateRequest, diags := plan.ToAPIUpdate(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the workflow via the API client using PUT (full update)
	tflog.Debug(ctx, "Updating workflow via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	workflowAPIResponse, err := r.client.UpdateWorkflow(ctx, state.ID.ValueString(), &apiUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Workflow",
			fmt.Sprintf("Could not update SailPoint Workflow %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Workflow", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if workflowAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Workflow",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState workflowModel
	tflog.Debug(ctx, "Mapping SailPoint Workflow API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(newState.FromAPI(ctx, *workflowAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for workflow resource", map[string]any{
		"id": newState.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Workflow resource", map[string]any{
		"id":   newState.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workflowModel
	tflog.Debug(ctx, "Getting state for workflow resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Workflows must be disabled before deletion
	// Check if workflow is enabled and disable it first if needed
	if !state.Enabled.IsNull() && state.Enabled.ValueBool() {
		tflog.Info(ctx, "Workflow is enabled, disabling before deletion", map[string]any{
			"id": state.ID.ValueString(),
		})

		// Convert to API model and disable
		apiWorkflow, diags := state.ToAPIUpdate(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		apiWorkflow.Enabled = false

		// Update to disable the workflow
		_, err := r.client.UpdateWorkflow(ctx, state.ID.ValueString(), &apiWorkflow)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Disabling Workflow",
				fmt.Sprintf("Could not disable workflow %q before deletion: %s", state.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Info(ctx, "Workflow disabled successfully before deletion", map[string]any{
			"id": state.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, "Deleting workflow via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	err := r.client.DeleteWorkflow(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Workflow",
			fmt.Sprintf("Could not delete SailPoint Workflow %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Workflow", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Workflow resource", map[string]any{
		"id": state.ID.ValueString(),
	})
}

// ImportState implements resource.ResourceWithImportState.
func (r *workflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing workflow resource", map[string]any{
		"id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Workflow resource", map[string]any{
		"id": req.ID,
	})
}
