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
	_ datasource.DataSource              = &identityAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &identityAttributeDataSource{}
)

type identityAttributeDataSource struct {
	client *client.Client
}

func NewIdentityAttributeDataSource() datasource.DataSource {
	return &identityAttributeDataSource{}
}

func (d *identityAttributeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identityAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_attribute"
}

func (d *identityAttributeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.IdentityAttributeSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Identity Attribute.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Identity Attribute. Identity attributes are configurable fields on identity objects. See [Identity Attributes API](https://developer.sailpoint.com/docs/api/v2025/list-identity-attributes/) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *identityAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Identity Attribute data source")

	var config models.IdentityAttribute
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity attribute via API
	fetchedAttribute, err := d.client.GetIdentityAttribute(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Attribute",
			fmt.Sprintf("Could not read identity attribute %s: %s", config.Name.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.IdentityAttribute
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedAttribute); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute Response",
			fmt.Sprintf("Could not convert identity attribute response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Attribute data source read successfully", map[string]interface{}{
		"attribute_name": state.Name.ValueString(),
	})
}
