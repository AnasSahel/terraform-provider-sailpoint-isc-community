// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &TransformsDataSource{}
	_ datasource.DataSourceWithConfigure = &TransformsDataSource{}
)

// TransformsDataSource is the data source implementation.
type TransformsDataSource struct {
	client *sailpoint.APIClient
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

	client, ok := req.ProviderData.(*sailpoint.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set up API request with optional filtering
	apiReq := d.client.V2025.TransformsAPI.ListTransforms(context.Background())

	// Apply filters if specified
	if !state.Filters.IsNull() && !state.Filters.IsUnknown() && state.Filters.ValueString() != "" {
		apiReq = apiReq.Filters(state.Filters.ValueString())
	}

	// Execute API request
	transforms, response, err := apiReq.Execute()
	if err != nil {
		errorMsg := "Failed to retrieve transforms"
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if response != nil {
			switch response.StatusCode {
			case 400:
				detailMsg = fmt.Sprintf("Bad Request - Invalid filter syntax. Please check the 'filters' parameter. API error: %s", err.Error())
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to read transforms. Please check your user permissions in SailPoint."
			case 429:
				detailMsg = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
			default:
				detailMsg = fmt.Sprintf("HTTP %d - %s", response.StatusCode, err.Error())
			}
			detailMsg += fmt.Sprintf("\nHTTP Response: %v", response)
		}

		resp.Diagnostics.AddError(errorMsg, detailMsg)
		return
	}

	// Map response to Terraform model
	diags := state.FromSailPointTransformsRead(ctx, transforms)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set ID for the data source
	state.Id = types.StringValue(fmt.Sprintf("transforms-%d", time.Now().Unix()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
