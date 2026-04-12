// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &entitlementDataSource{}
	_ datasource.DataSourceWithConfigure = &entitlementDataSource{}
)

type entitlementDataSource struct {
	client *client.Client
}

// NewEntitlementDataSource creates a new data source for SailPoint Entitlement.
func NewEntitlementDataSource() datasource.DataSource {
	return &entitlementDataSource{}
}

func (d *entitlementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (d *entitlementDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "entitlement data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *entitlementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Entitlement.",
		MarkdownDescription: "Data source for SailPoint Entitlement. Look up an entitlement by ID to retrieve its attributes.",
		Attributes: map[string]schema.Attribute{
			"id":                        schema.StringAttribute{Required: true, MarkdownDescription: "The unique identifier of the entitlement."},
			"name":                      schema.StringAttribute{Computed: true},
			"description":               schema.StringAttribute{Computed: true},
			"attribute":                 schema.StringAttribute{Computed: true},
			"value":                     schema.StringAttribute{Computed: true},
			"source_schema_object_type": schema.StringAttribute{Computed: true},
			"privileged":                schema.BoolAttribute{Computed: true},
			"cloud_governed":            schema.BoolAttribute{Computed: true},
			"requestable":               schema.BoolAttribute{Computed: true},
			"owner": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"source": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"segments": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"manually_updated_fields": schema.MapAttribute{
				Computed:    true,
				ElementType: types.BoolType,
			},
			"created":  schema.StringAttribute{Computed: true},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *entitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state entitlementModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Debug(ctx, "Reading entitlement data source", map[string]any{"id": id})

	apiResp, err := d.client.GetEntitlement(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Entitlement",
			fmt.Sprintf("Could not read entitlement %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Entitlement", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
