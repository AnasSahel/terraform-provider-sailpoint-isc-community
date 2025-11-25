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
	_ datasource.DataSource              = &entitlementDataSource{}
	_ datasource.DataSourceWithConfigure = &entitlementDataSource{}
)

type entitlementDataSource struct {
	client *client.Client
}

func NewEntitlementDataSource() datasource.DataSource {
	return &entitlementDataSource{}
}

func (d *entitlementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *entitlementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (d *entitlementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.EntitlementSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Entitlement.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Entitlement. Entitlements represent permissions or access rights from source systems (e.g., AD groups, application roles). See [Entitlement Documentation](https://developer.sailpoint.com/docs/api/v2025/get-entitlement) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *entitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Entitlement data source")

	var config models.Entitlement
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the entitlement via API
	fetchedEntitlement, err := d.client.GetEntitlement(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Entitlement",
			fmt.Sprintf("Could not read entitlement ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.Entitlement
	state.ConvertFromSailPointForDataSource(ctx, fetchedEntitlement)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Entitlement data source read successfully", map[string]interface{}{
		"entitlement_id": state.ID.ValueString(),
	})
}
