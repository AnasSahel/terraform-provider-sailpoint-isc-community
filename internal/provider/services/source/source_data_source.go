// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type SourceDataSource struct {
	client *api_v2025.APIClient
}

var (
	_ datasource.DataSource              = &SourceDataSource{}
	_ datasource.DataSourceWithConfigure = &SourceDataSource{}
)

func NewSourceDataSource() datasource.DataSource {
	return &SourceDataSource{}
}

func (d *SourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring SourceDataSource")

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

	tflog.Debug(ctx, "Configured SourceDataSource")
	d.client = client
}

func (d *SourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *SourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Trace(ctx, "Preparing SourceDataSource schema")
	resp.Schema = GetSourceDataSourceSchema()
}

func (d *SourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Source Data Source")

	var config SourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var source *api_v2025.Source
	var httpResponse interface{}
	var err error

	// Fetch source by ID if provided
	if !config.Id.IsNull() && !config.Id.IsUnknown() {
		sourceID := config.Id.ValueString()
		tflog.Debug(ctx, "Fetching source by ID", map[string]interface{}{
			"id": sourceID,
		})

		source, httpResponse, err = d.client.SourcesAPI.GetSource(
			context.Background(),
			sourceID,
		).Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Source",
				fmt.Sprintf("Could not read source with ID '%s': %s\n\nHTTP Response: %v",
					sourceID, err.Error(), httpResponse),
			)
			return
		}
	} else if !config.Name.IsNull() && !config.Name.IsUnknown() {
		// Fetch source by name
		sourceName := config.Name.ValueString()
		tflog.Debug(ctx, "Fetching source by name", map[string]interface{}{
			"name": sourceName,
		})

		sources, httpResp, err := d.client.SourcesAPI.ListSources(
			context.Background(),
		).Filters(fmt.Sprintf(`name eq "%s"`, sourceName)).Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing Sources",
				fmt.Sprintf("Could not list sources to find '%s': %s\n\nHTTP Response: %v",
					sourceName, err.Error(), httpResp),
			)
			return
		}

		if len(sources) == 0 {
			resp.Diagnostics.AddError(
				"Source Not Found",
				fmt.Sprintf("No source found with name '%s'", sourceName),
			)
			return
		}

		if len(sources) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Sources Found",
				fmt.Sprintf("Found %d sources with name '%s', expected exactly 1", len(sources), sourceName),
			)
			return
		}

		source = &sources[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a source",
		)
		return
	}

	// Convert API response to Terraform state
	var state SourceDataSourceModel
	conversionDiags := state.FromSailPointSource(ctx, source)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Successfully read source", map[string]interface{}{
		"id":     state.Id.ValueString(),
		"name":   state.Name.ValueString(),
		"status": state.Status.ValueString(),
	})

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
