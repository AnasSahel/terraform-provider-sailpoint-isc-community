// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transform

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &transformDataSource{}
	_ datasource.DataSourceWithConfigure = &transformDataSource{}
)

type transformDataSource struct {
	client *client.Client
}

func NewTransformDataSource() datasource.DataSource {
	return &transformDataSource{}
}

func (d *transformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

func (d *transformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "transform data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *transformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Transform.",
		MarkdownDescription: "Data source for SailPoint Transform. Transforms are used to manipulate attribute values during identity processing.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the transform.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the transform.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the transform (e.g., 'lower', 'upper', 'concat', 'substring').",
				Computed:            true,
			},
			"attributes": schema.StringAttribute{
				MarkdownDescription: "A JSON object containing the transform-specific configuration attributes.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
		},
	}
}

func (d *transformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	tflog.Debug(ctx, "Getting config for transform data source")
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &config.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := config.ID.ValueString()

	tflog.Debug(ctx, "Fetching transform from SailPoint", map[string]any{
		"id": id,
	})
	transformResponse, err := d.client.GetTransform(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Transform",
			fmt.Sprintf("Could not read SailPoint Transform %q: %s", id, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Transform", map[string]any{
			"id":    id,
			"error": err.Error(),
		})
		return
	}

	if transformResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Transform",
			"Received nil response from SailPoint API",
		)
		return
	}

	var state transformModel
	resp.Diagnostics.Append(state.FromAPI(ctx, *transformResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Transform data source", map[string]any{
		"id":   id,
		"name": state.Name.ValueString(),
	})
}
