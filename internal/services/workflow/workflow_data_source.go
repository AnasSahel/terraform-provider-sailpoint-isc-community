// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &workflowDataSource{}
	_ datasource.DataSourceWithConfigure = &workflowDataSource{}
)

type workflowDataSource struct {
	client *client.Client
}

func NewWorkflowDataSource() datasource.DataSource {
	return &workflowDataSource{}
}

// Metadata implements datasource.DataSource.
func (d *workflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *workflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "workflow data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

// Schema implements datasource.DataSource.
func (d *workflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a SailPoint Workflow by ID.",
		MarkdownDescription: "Retrieves a SailPoint Workflow by ID. Workflows are custom automation scripts that respond to event triggers and perform a series of actions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the workflow.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workflow.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the workflow.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the workflow is enabled.",
				Computed:            true,
			},
			"trigger": schema.StringAttribute{
				MarkdownDescription: "The trigger configuration as JSON.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the workflow was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the workflow was last modified.",
				Computed:            true,
			},
			"creator": schema.SingleNestedAttribute{
				MarkdownDescription: "The identity who created the workflow.",
				Computed:            true,
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
			},
			"failure_count": schema.Int32Attribute{
				MarkdownDescription: "The number of times the workflow has failed.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the workflow.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner (e.g., `IDENTITY`).",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
					},
				},
			},
			"definition": schema.SingleNestedAttribute{
				MarkdownDescription: "The workflow definition containing the steps to execute.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"start": schema.StringAttribute{
						MarkdownDescription: "The name of the starting step.",
						Computed:            true,
					},
					"steps": schema.StringAttribute{
						MarkdownDescription: "JSON object containing the workflow steps.",
						Computed:            true,
						CustomType:          jsontypes.NormalizedType{},
					},
				},
			},
		},
	}
}

// Read implements datasource.DataSource.
func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config workflowModel
	tflog.Debug(ctx, "Getting config for workflow data source")
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the workflow from SailPoint
	tflog.Debug(ctx, "Fetching workflow from SailPoint", map[string]any{
		"id": config.ID.ValueString(),
	})
	workflowResponse, err := d.client.GetWorkflow(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Workflow",
			fmt.Sprintf("Could not read SailPoint Workflow %q: %s", config.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Workflow", map[string]any{
			"id":    config.ID.ValueString(),
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

	// Map the response to the data source model
	var state workflowModel
	tflog.Debug(ctx, "Mapping SailPoint Workflow API response to data source model", map[string]any{
		"id": config.ID.ValueString(),
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *workflowResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for workflow data source", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Workflow data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}
