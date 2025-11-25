// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
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

func NewLauncherDataSource() datasource.DataSource {
	return &launcherDataSource{}
}

func (d *launcherDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = c
}

func (d *launcherDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_launcher"
}

func (d *launcherDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.LauncherSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Launcher.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Launcher. Launchers are interactive processes that can be triggered to execute workflows or other automation tasks. See [Launcher Documentation](https://developer.sailpoint.com/docs/api/v2025/get-launcher) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *launcherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Launcher data source")

	var config models.Launcher
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the launcher via API
	fetchedLauncher, err := d.client.GetLauncher(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Launcher",
			fmt.Sprintf("Could not read launcher ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.Launcher
	state.ConvertFromSailPointForDataSource(ctx, fetchedLauncher)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Launcher data source read successfully", map[string]interface{}{
		"launcher_id": state.ID.ValueString(),
	})
}
