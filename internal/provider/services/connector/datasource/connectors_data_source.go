// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_datasource

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

const (
	DefaultMaxResults = 10000 // Default max results for pagination
	DefaultPageSize   = 250   // Default page size for pagination
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

// applyCommonFilters applies filters and locale to the API request.
func (d *ConnectorsDataSource) applyCommonFilters(req api_v2025.ApiGetConnectorListRequest, data *ConnectorsDataSourceModel) api_v2025.ApiGetConnectorListRequest {
	// Apply filters if specified
	if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
		req = req.Filters(data.Filters.ValueString())
	}

	// Apply locale if specified
	if !data.Locale.IsNull() && !data.Locale.IsUnknown() {
		req = req.Locale(data.Locale.ValueString())
	}

	return req
}

// applyStandardFilters applies additional filters for non-paginated requests.
func (d *ConnectorsDataSource) applyStandardFilters(req api_v2025.ApiGetConnectorListRequest, data *ConnectorsDataSourceModel) api_v2025.ApiGetConnectorListRequest {
	// Apply limit if specified
	if !data.Limit.IsNull() && !data.Limit.IsUnknown() {
		req = req.Limit(data.Limit.ValueInt32())
	}

	// Apply offset if specified
	if !data.Offset.IsNull() && !data.Offset.IsUnknown() {
		req = req.Offset(data.Offset.ValueInt32())
	}

	// Apply include_count if specified
	if !data.IncludeCount.IsNull() && !data.IncludeCount.IsUnknown() {
		req = req.Count(data.IncludeCount.ValueBool())
	}

	return req
}

// fetchConnectorsWithPagination retrieves connectors using SailPoint SDK pagination.
func (d *ConnectorsDataSource) fetchConnectorsWithPagination(ctx context.Context, apiReq api_v2025.ApiGetConnectorListRequest, data *ConnectorsDataSourceModel) ([]api_v2025.V3ConnectorDto, error) {
	tflog.Debug(ctx, "Using SailPoint pagination to fetch all connectors")

	// Determine pagination parameters
	var maxResults int32 = DefaultMaxResults
	if !data.MaxResults.IsNull() && !data.MaxResults.IsUnknown() {
		maxResults = data.MaxResults.ValueInt32()
	}

	var pageSize int32 = DefaultPageSize
	if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
		pageSize = data.PageSize.ValueInt32()
	}

	// Use SailPoint SDK pagination
	connectors, _, err := sailpoint.Paginate[api_v2025.V3ConnectorDto](apiReq, 0, pageSize, maxResults)
	if err != nil {
		return nil, fmt.Errorf("could not paginate connectors: %v", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d connectors using pagination", len(connectors)))
	return connectors, nil
}

// fetchConnectorsStandard retrieves connectors using standard API request.
func (d *ConnectorsDataSource) fetchConnectorsStandard(ctx context.Context, apiReq api_v2025.ApiGetConnectorListRequest, data *ConnectorsDataSourceModel) ([]api_v2025.V3ConnectorDto, error) {
	tflog.Debug(ctx, "Using standard API request for connectors")

	apiReq = d.applyStandardFilters(apiReq, data)

	// Execute API request
	connectors, httpResp, err := apiReq.Execute()
	if err != nil {
		return nil, fmt.Errorf("could not read connectors: %v\nHTTP Response: %+v", err, httpResp)
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d connectors using standard API request", len(connectors)))
	return connectors, nil
}

// convertConnectorsToModel converts API response to Terraform model.
func (d *ConnectorsDataSource) convertConnectorsToModel(ctx context.Context, connectors []api_v2025.V3ConnectorDto) ([]ConnectorSummaryModel, error) {
	models := make([]ConnectorSummaryModel, len(connectors))
	for i, conn := range connectors {
		connectorModel := ConnectorSummaryModel{}
		if err := connectorModel.FromSailPointV3ConnectorDto(ctx, &conn); err != nil {
			return nil, fmt.Errorf("could not convert connector %s: %v", conn.GetName(), err)
		}
		models[i] = connectorModel
	}
	return models, nil
}

func (d *ConnectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectorsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set up base API request with common filters
	apiReq := d.client.ConnectorsAPI.GetConnectorList(ctx)
	apiReq = d.applyCommonFilters(apiReq, &data)

	// Fetch connectors based on pagination preference
	var connectors []api_v2025.V3ConnectorDto
	var err error

	if !data.PaginateAll.IsNull() && data.PaginateAll.ValueBool() {
		connectors, err = d.fetchConnectorsWithPagination(ctx, apiReq, &data)
		if err != nil {
			resp.Diagnostics.AddError("Error Reading Connectors with Pagination", err.Error())
			return
		}
	} else {
		connectors, err = d.fetchConnectorsStandard(ctx, apiReq, &data)
		if err != nil {
			resp.Diagnostics.AddError("Error Reading Connectors", err.Error())
			return
		}
	}

	// Convert API response to Terraform model
	data.Connectors, err = d.convertConnectorsToModel(ctx, connectors)
	if err != nil {
		resp.Diagnostics.AddError("Error Converting Connector Data", err.Error())
		return
	}

	// Set ID for the data source
	data.ID = types.StringValue(fmt.Sprintf("connectors-%d", time.Now().Unix()))

	tflog.Trace(ctx, fmt.Sprintf("Successfully processed %d connectors", len(data.Connectors)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
