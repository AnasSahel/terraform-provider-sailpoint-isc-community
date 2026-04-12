// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package access_profile

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
	_ datasource.DataSource              = &accessProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &accessProfileDataSource{}
)

type accessProfileDataSource struct {
	client *client.Client
}

func NewAccessProfileDataSource() datasource.DataSource { return &accessProfileDataSource{} }

func (d *accessProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_profile"
}

func (d *accessProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "access profile data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func computedObjectRef() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{Computed: true},
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *accessProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Access Profile.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Required: true},
			"name":        schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"enabled":     schema.BoolAttribute{Computed: true},
			"requestable": schema.BoolAttribute{Computed: true},
			"owner":       computedObjectRef(),
			"source":      computedObjectRef(),
			"entitlements": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Computed: true},
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"segments": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"additional_owners": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Computed: true},
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"access_request_config": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"comments_required":        schema.BoolAttribute{Computed: true},
					"denial_comments_required": schema.BoolAttribute{Computed: true},
					"reauthorization_required": schema.BoolAttribute{Computed: true},
					"require_end_date":         schema.BoolAttribute{Computed: true},
					"max_permitted_access_duration": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"value":     schema.Int64Attribute{Computed: true},
							"time_unit": schema.StringAttribute{Computed: true},
						},
					},
					"approval_schemes": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"approver_type": schema.StringAttribute{Computed: true},
								"approver_id":   schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"revoke_request_config": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"approval_schemes": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"approver_type": schema.StringAttribute{Computed: true},
								"approver_id":   schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"provisioning_criteria": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"operation": schema.StringAttribute{Computed: true},
					"attribute": schema.StringAttribute{Computed: true},
					"value":     schema.StringAttribute{Computed: true},
					"children": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"operation": schema.StringAttribute{Computed: true},
								"attribute": schema.StringAttribute{Computed: true},
								"value":     schema.StringAttribute{Computed: true},
								"children": schema.ListNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"operation": schema.StringAttribute{Computed: true},
											"attribute": schema.StringAttribute{Computed: true},
											"value":     schema.StringAttribute{Computed: true},
										},
									},
								},
							},
						},
					},
				},
			},
			"created":  schema.StringAttribute{Computed: true},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *accessProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accessProfileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Debug(ctx, "Reading access profile data source", map[string]any{"id": id})

	apiResp, err := d.client.GetAccessProfile(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Access Profile",
			fmt.Sprintf("Could not read access profile %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Access Profile", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
