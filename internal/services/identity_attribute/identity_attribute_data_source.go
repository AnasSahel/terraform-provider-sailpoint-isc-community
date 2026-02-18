// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_attribute

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &identityAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &identityAttributeDataSource{}
)

type identityAttributeDataSource struct {
	client *client.Client
}

func NewIdentityAttributeDataSource() datasource.DataSource {
	return &identityAttributeDataSource{}
}

func (d *identityAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_attribute"
}

func (d *identityAttributeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "identity attribute data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *identityAttributeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Identity Attribute.",
		MarkdownDescription: "Data source for SailPoint Identity Attribute.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity attribute.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the identity attribute.",
				Computed:            true,
			},
			"standard": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is a standard attribute.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the identity attribute.",
				Computed:            true,
			},
			"multi": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute supports multiple values.",
				Computed:            true,
			},
			"searchable": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is searchable.",
				Computed:            true,
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the identity attribute is a system attribute.",
				Computed:            true,
			},
			"sources": schema.ListNestedAttribute{
				MarkdownDescription: "The sources associated with the identity attribute.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Attribute mapping type. Mostly `rule`.",
							Computed:            true,
						},
						"properties": schema.StringAttribute{
							MarkdownDescription: "Attribute mapping properties.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

func (d *identityAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SailPoint Identity Attribute data source")

	var config identityAttributeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity attribute from SailPoint
	tflog.Debug(ctx, "Fetching identity attribute from SailPoint", map[string]any{
		"name": config.Name.ValueString(),
	})
	identityAttributeResponse, err := d.client.GetIdentityAttribute(ctx, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Attribute",
			fmt.Sprintf("Could not read SailPoint Identity Attribute %q: %s", config.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Identity Attribute", map[string]any{
			"name":  config.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if identityAttributeResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Attribute",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the data source model
	var state identityAttributeModel
	resp.Diagnostics.Append(state.FromAPI(ctx, *identityAttributeResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Identity Attribute data source", map[string]any{
		"name": config.Name.ValueString(),
	})
}
