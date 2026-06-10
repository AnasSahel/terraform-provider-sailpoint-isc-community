// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceAttrSyncConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceAttrSyncConfigDataSource{}
)

type sourceAttrSyncConfigDataSource struct {
	client *client.Client
}

// NewSourceAttributeSyncConfigDataSource creates a new data source for source attribute sync config.
func NewSourceAttributeSyncConfigDataSource() datasource.DataSource {
	return &sourceAttrSyncConfigDataSource{}
}

func (d *sourceAttrSyncConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_attribute_sync_config"
}

func (d *sourceAttrSyncConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source attribute sync config data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *sourceAttrSyncConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Source Attribute Sync Config.",
		MarkdownDescription: "Data source for SailPoint Source Attribute Sync Config. Returns the full attribute sync configuration for a source, including all syncable attributes and their current enabled state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier — mirrors `source_id`.",
			},
			"source_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the source to look up.",
			},
			"source": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Reference to the source.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"attributes": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "All syncable attributes for the source and their current enabled state.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":         schema.StringAttribute{Computed: true, MarkdownDescription: "The identity attribute name."},
						"enabled":      schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether this attribute is currently synced to the source."},
						"display_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The human-readable display name."},
						"target":       schema.StringAttribute{Computed: true, MarkdownDescription: "The account attribute name on the target source."},
					},
				},
			},
		},
	}
}

func (d *sourceAttrSyncConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourceAttrSyncConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	tflog.Debug(ctx, "Reading source attribute sync config data source", map[string]any{"source_id": sourceID})

	apiResp, err := d.client.GetSourceAttributeSyncConfig(ctx, sourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source Attribute Sync Config",
			fmt.Sprintf("Could not read attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	// Data source returns all attributes (no filtering).
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
