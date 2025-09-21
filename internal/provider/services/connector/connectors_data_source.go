// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ConnectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConnectorsDataSource{}
)

func NewConnectorsDataSource() datasource.DataSource {
	return &ConnectorsDataSource{}
}

// ConnectorsDataSource defines the data source implementation.
type ConnectorsDataSource struct {
	client *api_v2025.APIClient
}

func (d *ConnectorsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connectors"
}

func (d *ConnectorsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetConnectorsDataSourceSchema()
}

func (d *ConnectorsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *ConnectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectorsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var connectors []api_v2025.V3ConnectorDto
	var err error

	// Check if pagination is requested
	if !data.PaginateAll.IsNull() && data.PaginateAll.ValueBool() {
		tflog.Debug(ctx, "Using SailPoint pagination to fetch all connectors")

		// Set up base API request for pagination
		baseReq := d.client.ConnectorsAPI.GetConnectorList(ctx)

		// Apply filters if specified
		if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
			baseReq = baseReq.Filters(data.Filters.ValueString())
		}

		// Apply locale if specified
		if !data.Locale.IsNull() && !data.Locale.IsUnknown() {
			baseReq = baseReq.Locale(data.Locale.ValueString())
		}

		// Determine pagination parameters
		var maxResults int32 = 10000 // Default max results
		if !data.MaxResults.IsNull() && !data.MaxResults.IsUnknown() {
			maxResults = data.MaxResults.ValueInt32()
		}

		var pageSize int32 = 250 // Default page size
		if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
			pageSize = data.PageSize.ValueInt32()
		}

		// Use SailPoint SDK pagination
		connectors, _, err = sailpoint.Paginate[api_v2025.V3ConnectorDto](baseReq, 0, pageSize, maxResults)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Connectors with Pagination",
				fmt.Sprintf("Could not paginate connectors: %v", err),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Fetched %d connectors using pagination", len(connectors)))

	} else {
		// Use standard API request without pagination
		tflog.Debug(ctx, "Using standard API request for connectors")

		apiReq := d.client.ConnectorsAPI.GetConnectorList(ctx)

		// Apply filters if specified
		if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
			apiReq = apiReq.Filters(data.Filters.ValueString())
		}

		// Apply limit if specified
		if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
			apiReq = apiReq.Limit(data.Limit.ValueInt32())
		}

		// Apply offset if specified
		if !data.Offset.IsNull() && !data.Offset.IsUnknown() {
			apiReq = apiReq.Offset(data.Offset.ValueInt32())
		}

		// Apply include_count if specified
		if !data.IncludeCount.IsNull() && !data.IncludeCount.IsUnknown() {
			apiReq = apiReq.Count(data.IncludeCount.ValueBool())
		}

		// Apply locale if specified
		if !data.Locale.IsNull() && !data.Locale.IsUnknown() {
			apiReq = apiReq.Locale(data.Locale.ValueString())
		}

		// Execute API request
		var httpResp interface{}
		connectors, httpResp, err = apiReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Connectors",
				fmt.Sprintf("Could not read connectors: %v\nHTTP Response: %+v", err, httpResp),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Fetched %d connectors using standard API request", len(connectors)))
	}

	// Convert API response to Terraform model
	data.Connectors = make([]ConnectorSummaryModel, len(connectors))
	for i, conn := range connectors {
		connectorModel := ConnectorSummaryModel{}
		if err := connectorModel.FromSailPointV3ConnectorDto(ctx, &conn); err != nil {
			resp.Diagnostics.AddError(
				"Error Converting Connector Data",
				fmt.Sprintf("Could not convert connector %s: %v", conn.GetName(), err),
			)
			return
		}
		data.Connectors[i] = connectorModel
	}

	// Set ID for the data source
	data.ID = types.StringValue(fmt.Sprintf("connectors-%d", time.Now().Unix()))

	tflog.Trace(ctx, fmt.Sprintf("Successfully processed %d connectors", len(data.Connectors)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
