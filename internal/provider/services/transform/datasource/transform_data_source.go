// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &TransformDataSource{}
var _ datasource.DataSourceWithConfigure = &TransformDataSource{}

// TransformDataSource is the data source implementation.
type TransformDataSource struct {
	client *api_v2025.APIClient
}

// NewTransformDataSource is a helper function to simplify the provider implementation.
func NewTransformDataSource() datasource.DataSource {
	return &TransformDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *TransformDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *TransformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

// Schema defines the schema for the data source.
func (d *TransformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetTransformDataSourceSchema()
}

// Read refreshes the Terraform state with the latest data.
func (d *TransformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TransformDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either ID or Name is provided
	if (state.Id.IsNull() || state.Id.ValueString() == "") && (state.Name.IsNull() || state.Name.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Transform Identifier",
			"Either 'id' or 'name' must be specified to retrieve a transform.",
		)
		return
	}

	var transformId string

	// If ID is provided, use it directly
	if !state.Id.IsNull() && state.Id.ValueString() != "" {
		transformId = state.Id.ValueString()
	} else {
		// If only name is provided, search for the transform by name
		filterQuery := fmt.Sprintf("name eq \"%s\"", state.Name.ValueString())
		transforms, response, err := d.client.TransformsAPI.ListTransforms(context.Background()).
			Filters(filterQuery).Execute()
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to search for transform by name '%s'", state.Name.ValueString())
			detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

			// Add specific handling for common error scenarios
			if response != nil {
				switch response.StatusCode {
				case 400:
					detailMsg = fmt.Sprintf("Bad Request - Invalid search query. Please check the transform name '%s' is valid. API error: %s", state.Name.ValueString(), err.Error())
				case 401:
					detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
				case 403:
					detailMsg = "Forbidden - Insufficient permissions to search transforms. Please check your user permissions in SailPoint."
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

		if len(transforms) == 0 {
			resp.Diagnostics.AddError(
				"Transform Not Found",
				fmt.Sprintf("No transform found with name '%s'", state.Name.ValueString()),
			)
			return
		}

		if len(transforms) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Transforms Found",
				fmt.Sprintf("Multiple transforms found with name '%s'. Use 'id' instead to specify a unique transform.", state.Name.ValueString()),
			)
			return
		}

		transformId = transforms[0].GetId()
	}

	// Fetch the transform by ID
	transform, response, err := d.client.TransformsAPI.GetTransform(context.Background(), transformId).Execute()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to retrieve transform with ID '%s'", transformId)
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if response != nil {
			switch response.StatusCode {
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to read transforms. Please check your user permissions in SailPoint."
			case 404:
				detailMsg = fmt.Sprintf("Transform with ID '%s' not found. Please verify the ID is correct.", transformId)
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
	diags := state.FromSailPointTransformRead(ctx, *transform)
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
