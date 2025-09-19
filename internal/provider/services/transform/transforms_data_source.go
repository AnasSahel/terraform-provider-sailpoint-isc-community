// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &TransformsDataSource{}
var _ datasource.DataSourceWithConfigure = &TransformsDataSource{}

// TransformsDataSource is the data source implementation.
type TransformsDataSource struct {
	client *api_v2025.APIClient
}

// NewTransformsDataSource is a helper function to simplify the provider implementation.
func NewTransformsDataSource() datasource.DataSource {
	return &TransformsDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *TransformsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *TransformsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transforms"
}

// Schema defines the schema for the data source.
func (d *TransformsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetTransformsDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *TransformsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TransformsDataSourceModel

	transforms, response, err := d.client.TransformsAPI.ListTransforms(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SailPoint Transforms",
			fmt.Sprintf("SailPoint API error while reading transforms: %s\nHTTP Response: %v",
				err.Error(),
				response,
			),
		)
		return
	}

	// Map response to Terraform model
	diags := state.FromSailPointTransformsRead(ctx, transforms)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
