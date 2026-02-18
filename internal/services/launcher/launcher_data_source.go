// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &launcherDataSource{}
	_ datasource.DataSourceWithConfigure = &launcherDataSource{}
)

type launcherDataSource struct {
	client *client.Client
}

// NewLauncherDataSource creates a new data source for Launcher.
func NewLauncherDataSource() datasource.DataSource {
	return &launcherDataSource{}
}

// Metadata implements datasource.DataSource.
func (d *launcherDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launcher"
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *launcherDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "launcher data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

// Schema implements datasource.DataSource.
func (d *launcherDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a SailPoint Launcher by ID.",
		MarkdownDescription: "Retrieves a SailPoint Launcher by ID. Launchers are used to trigger workflows through the SailPoint UI.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the launcher.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the launcher, limited to 255 characters.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the launcher, limited to 2000 characters.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the launcher. Currently only `INTERACTIVE_PROCESS` is supported.",
				Computed:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the launcher is disabled.",
				Computed:            true,
			},
			"config": schema.StringAttribute{
				MarkdownDescription: "JSON configuration associated with this launcher, restricted to a max size of 4KB.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the launcher was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the launcher was last modified.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the launcher.",
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
			"reference": schema.SingleNestedAttribute{
				MarkdownDescription: "The reference to the resource this launcher triggers (e.g., a workflow).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the reference (e.g., `WORKFLOW`).",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the referenced resource.",
						Computed:            true,
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

// Read implements datasource.DataSource.
func (d *launcherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config launcherModel
	tflog.Debug(ctx, "Getting config for launcher data source")
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the launcher from SailPoint
	tflog.Debug(ctx, "Fetching launcher from SailPoint", map[string]any{
		"id": config.ID.ValueString(),
	})
	launcherResponse, err := d.client.GetLauncher(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Launcher",
			fmt.Sprintf("Could not read SailPoint Launcher %q: %s", config.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Launcher", map[string]any{
			"id":    config.ID.ValueString(),
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

	// Map the response to the data source model
	var state launcherModel
	tflog.Debug(ctx, "Mapping SailPoint Launcher API response to data source model", map[string]any{
		"id": config.ID.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *launcherResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for launcher data source", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Launcher data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}
