// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceDataSource{}
)

type sourceDataSource struct {
	client *client.Client
}

func NewSourceDataSource() datasource.DataSource {
	return &sourceDataSource{}
}

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Source.",
		MarkdownDescription: "Data source for SailPoint Source. Sources represent managed systems (e.g., Active Directory, Workday) in Identity Security Cloud.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the source.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the source.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the source.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
					},
				},
			},
			"cluster": schema.SingleNestedAttribute{
				MarkdownDescription: "The cluster associated with this source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the cluster.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the cluster.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the cluster.",
						Computed:            true,
					},
				},
			},
			"connector": schema.StringAttribute{
				MarkdownDescription: "The connector script name.",
				Computed:            true,
			},
			"connector_class": schema.StringAttribute{
				MarkdownDescription: "The fully qualified name of the Java class that implements the connector interface.",
				Computed:            true,
			},
			"connector_attributes": schema.StringAttribute{
				MarkdownDescription: "A JSON object containing connector-specific configuration.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"connection_type": schema.StringAttribute{
				MarkdownDescription: "The connection type.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of system being managed.",
				Computed:            true,
			},
			"delete_threshold": schema.Int64Attribute{
				MarkdownDescription: "The percentage threshold for skipping the delete phase (0-100).",
				Computed:            true,
			},
			"authoritative": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is referenced by an identity profile.",
				Computed:            true,
			},
			"healthy": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is healthy.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the source.",
				Computed:            true,
			},
			"features": schema.ListAttribute{
				MarkdownDescription: "The list of features enabled for the source.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"credential_provider_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether credential provider is enabled for the source.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The source category.",
				Computed:            true,
			},
			"provision_as_csv": schema.BoolAttribute{
				MarkdownDescription: "Whether the source was provisioned as a CSV source.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SailPoint Source data source")

	var config sourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching source from SailPoint", map[string]any{
		"id": config.ID.ValueString(),
	})
	sourceResponse, err := d.client.GetSource(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source",
			fmt.Sprintf("Could not read SailPoint Source %q: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	if sourceResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source",
			"Received nil response from SailPoint API",
		)
		return
	}

	var state sourceModel
	resp.Diagnostics.Append(state.FromAPI(ctx, *sourceResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Source data source", map[string]any{
		"id":   config.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}
