// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &SourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &SourcesDataSource{}
)

func NewSourcesDataSource() datasource.DataSource {
	return &SourcesDataSource{}
}

// SourcesDataSource defines the data source implementation.
type SourcesDataSource struct {
	client *api_v2025.APIClient
}

func (d *SourcesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sources"
}

func (d *SourcesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetSourcesDataSourceSchema()
}

func (d *SourcesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SourcesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sources []api_v2025.Source
	var err error

	// Check if pagination is requested
	if !data.PaginateAll.IsNull() && data.PaginateAll.ValueBool() {
		tflog.Debug(ctx, "Using SailPoint pagination to fetch all sources")

		// Set up base API request for pagination
		baseReq := d.client.SourcesAPI.ListSources(ctx)

		// Apply filters if specified
		if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
			baseReq = baseReq.Filters(data.Filters.ValueString())
		}

		// Apply sorters if specified
		if !data.Sorters.IsNull() && !data.Sorters.IsUnknown() {
			baseReq = baseReq.Sorters(data.Sorters.ValueString())
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
		sources, _, err = sailpoint.Paginate[api_v2025.Source](baseReq, 0, pageSize, maxResults)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Sources with Pagination",
				fmt.Sprintf("Could not paginate sources: %v", err),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Fetched %d sources using pagination", len(sources)))

	} else {
		// Use standard API request without pagination
		tflog.Debug(ctx, "Using standard API request for sources")

		apiReq := d.client.SourcesAPI.ListSources(ctx)

		// Apply filters if specified
		if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
			apiReq = apiReq.Filters(data.Filters.ValueString())
		}

		// Apply sorters if specified
		if !data.Sorters.IsNull() && !data.Sorters.IsUnknown() {
			apiReq = apiReq.Sorters(data.Sorters.ValueString())
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

		// Execute API request
		var httpResp interface{}
		sources, httpResp, err = apiReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Sources",
				fmt.Sprintf("Could not read sources: %v\nHTTP Response: %+v", err, httpResp),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Fetched %d sources using standard API request", len(sources)))
	}

	// Convert API response to Terraform model
	data.Sources = make([]SourceSummaryModel, len(sources))
	for i, source := range sources {
		sourceModel := SourceSummaryModel{}
		if err := sourceModel.FromSailPointSource(ctx, &source); err != nil {
			resp.Diagnostics.AddError(
				"Error Converting Source Data",
				fmt.Sprintf("Could not convert source %s: %v", source.GetName(), err),
			)
			return
		}
		data.Sources[i] = sourceModel
	}

	// Set ID for the data source
	data.ID = types.StringValue(fmt.Sprintf("sources-%d", time.Now().Unix()))

	tflog.Trace(ctx, fmt.Sprintf("Successfully processed %d sources", len(data.Sources)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
