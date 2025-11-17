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
	_ datasource.DataSource              = &identityProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &identityProfileDataSource{}
)

type identityProfileDataSource struct {
	client *client.Client
}

func NewIdentityProfileDataSource() datasource.DataSource {
	return &identityProfileDataSource{}
}

func (d *identityProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identityProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_profile"
}

func (d *identityProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.IdentityProfileSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Identity Profile.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Identity Profile. Identity profiles define configurations for identities including authoritative sources and attribute mappings. See [Identity Profiles API](https://developer.sailpoint.com/docs/api/v2025/list-identity-profiles/) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *identityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Identity Profile data source")

	var config models.IdentityProfile
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity profile via API
	fetchedProfile, err := d.client.GetIdentityProfile(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Profile",
			fmt.Sprintf("Could not read identity profile %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.IdentityProfile
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedProfile); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Profile Response",
			fmt.Sprintf("Could not convert identity profile response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Profile data source read successfully", map[string]interface{}{
		"profile_id":   state.ID.ValueString(),
		"profile_name": state.Name.ValueString(),
	})
}
