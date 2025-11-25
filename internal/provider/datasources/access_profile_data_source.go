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
	_ datasource.DataSource              = &accessProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &accessProfileDataSource{}
)

type accessProfileDataSource struct {
	client *client.Client
}

func NewAccessProfileDataSource() datasource.DataSource {
	return &accessProfileDataSource{}
}

func (d *accessProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accessProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_profile"
}

func (d *accessProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.AccessProfileSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Access Profile.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Access Profile. Access Profiles are collections of entitlements from a source that can be requested by users. See [Access Profile Documentation](https://developer.sailpoint.com/docs/api/v2025/get-access-profile) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *accessProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Access Profile data source")

	var config models.AccessProfile
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the access profile via API
	fetchedAccessProfile, err := d.client.GetAccessProfile(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Access Profile",
			fmt.Sprintf("Could not read access profile ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.AccessProfile
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedAccessProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Access Profile Response",
			fmt.Sprintf("Could not convert access profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access Profile data source read successfully", map[string]interface{}{
		"access_profile_id": state.ID.ValueString(),
	})
}
