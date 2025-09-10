// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type transformModel struct {
	Id             types.String `tfsdk:"id"`
	Internal       types.Bool   `tfsdk:"internal"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	AttributesJson types.String `tfsdk:"attributes_json"`
}

type transformsDataSourceModel struct {
	Transforms []transformModel `tfsdk:"transforms"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &transformsDataSource{}
	_ datasource.DataSourceWithConfigure = &transformsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *transformsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// NewTransformsDataSource is a helper function to simplify the provider implementation.
func NewTransformsDataSource() datasource.DataSource {
	return &transformsDataSource{}
}

// transformsDataSource is the data source implementation.
type transformsDataSource struct {
	client *api_v2025.APIClient
}

// Metadata returns the data source type name.
func (d *transformsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transforms"
}

func (d *transformsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"transforms": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true},
						"internal":        schema.BoolAttribute{Computed: true},
						"name":            schema.StringAttribute{Computed: true},
						"type":            schema.StringAttribute{Computed: true},
						"attributes_json": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transformsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state transformsDataSourceModel

	transforms, _, err := d.client.TransformsAPI.ListTransforms(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read SailPoint Transforms",
			err.Error(),
		)
		return
	}

	for _, transform := range transforms {
		attrs, err := json.Marshal(transform.GetAttributes())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to parse Transform attributes",
				err.Error(),
			)
			return
		}

		transformState := transformModel{
			Id:             types.StringValue(transform.GetId()),
			Internal:       types.BoolValue(transform.GetInternal()),
			Name:           types.StringValue(transform.GetName()),
			Type:           types.StringValue(transform.GetType()),
			AttributesJson: types.StringValue(string(attrs)),
		}

		state.Transforms = append(state.Transforms, transformState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
