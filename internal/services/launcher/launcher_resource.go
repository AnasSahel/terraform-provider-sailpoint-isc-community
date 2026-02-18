// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// NewLauncherResource creates a new resource for Launcher.
func NewLauncherResource() resource.Resource {
	return &launcherResource{}
}

// Metadata implements resource.Resource.
func (r *launcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launcher"
}

// Configure implements resource.ResourceWithConfigure.
func (r *launcherResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "launcher resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *launcherResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Launcher.",
		MarkdownDescription: "Resource for SailPoint Launcher. Launchers are used to trigger workflows through the SailPoint UI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the launcher.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the launcher, limited to 255 characters.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the launcher, limited to 2000 characters.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the launcher. Currently only `INTERACTIVE_PROCESS` is supported.",
				Required:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the launcher is disabled. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
			},
			"config": schema.StringAttribute{
				MarkdownDescription: "JSON configuration associated with this launcher, restricted to a max size of 4KB.",
				Required:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the launcher was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the launcher was last modified.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the launcher.",
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
			"reference": schema.SingleNestedAttribute{
				MarkdownDescription: "The reference to the resource this launcher triggers (e.g., a workflow).",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the reference (e.g., `WORKFLOW`).",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the referenced resource.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the referenced resource.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *launcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan launcherModel
	tflog.Debug(ctx, "Getting plan for launcher resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping launcher resource model to API create request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the launcher via the API client
	tflog.Debug(ctx, "Creating launcher via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	launcherAPIResponse, err := r.client.CreateLauncher(ctx, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Launcher",
			fmt.Sprintf("Could not create SailPoint Launcher %q: %s", plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Launcher", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if launcherAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Launcher",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state launcherModel
	tflog.Debug(ctx, "Mapping SailPoint Launcher API response to resource model", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *launcherAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for launcher resource", map[string]any{
		"name": plan.Name.ValueString(),
		"id":   state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Launcher resource", map[string]any{
		"name": plan.Name.ValueString(),
		"id":   state.ID.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *launcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state launcherModel
	tflog.Debug(ctx, "Getting state for launcher resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the launcher from SailPoint
	tflog.Debug(ctx, "Fetching launcher from SailPoint", map[string]any{
		"id": state.ID.ValueString(),
	})
	launcherResponse, err := r.client.GetLauncher(ctx, state.ID.ValueString())
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Launcher not found, removing from state", map[string]any{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Launcher",
			fmt.Sprintf("Could not read SailPoint Launcher %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Launcher", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if launcherResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Launcher",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Launcher API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *launcherResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for launcher resource", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Launcher resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *launcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan launcherModel
	tflog.Debug(ctx, "Getting plan for launcher resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the ID
	var state launcherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping launcher resource model to API update request", map[string]any{
		"id": state.ID.ValueString(),
	})
	apiUpdateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the launcher via the API client (PUT)
	tflog.Debug(ctx, "Updating launcher via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	launcherAPIResponse, err := r.client.UpdateLauncher(ctx, state.ID.ValueString(), &apiUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Launcher",
			fmt.Sprintf("Could not update SailPoint Launcher %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Launcher", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if launcherAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Launcher",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState launcherModel
	tflog.Debug(ctx, "Mapping SailPoint Launcher API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(newState.FromAPI(ctx, *launcherAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the created timestamp from the prior state (it should never change)
	newState.Created = state.Created

	// Set the state
	tflog.Debug(ctx, "Setting state for launcher resource", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Launcher resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *launcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state launcherModel
	tflog.Debug(ctx, "Getting state for launcher resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting launcher via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	err := r.client.DeleteLauncher(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Launcher",
			fmt.Sprintf("Could not delete SailPoint Launcher %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Launcher", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Launcher resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// ImportState implements resource.ResourceWithImportState.
func (r *launcherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing launcher resource", map[string]any{
		"id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Launcher resource", map[string]any{
		"id": req.ID,
	})
}
