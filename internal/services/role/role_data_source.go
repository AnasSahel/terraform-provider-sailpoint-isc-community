// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

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
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

type roleDataSource struct {
	client *client.Client
}

func NewRoleDataSource() datasource.DataSource { return &roleDataSource{} }

func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *roleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "role data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func computedRef() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{Computed: true},
			"id":   schema.StringAttribute{Computed: true},
			"name": schema.StringAttribute{Computed: true},
		},
	}
}

func computedRefSet() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{Computed: true},
				"id":   schema.StringAttribute{Computed: true},
				"name": schema.StringAttribute{Computed: true},
			},
		},
	}
}

func computedApprovalSchemes() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"approver_type": schema.StringAttribute{Computed: true},
				"approver_id":   schema.StringAttribute{Computed: true},
			},
		},
	}
}

func computedCriteriaLeaf() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operation": schema.StringAttribute{Computed: true},
		"key": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"type":      schema.StringAttribute{Computed: true},
				"property":  schema.StringAttribute{Computed: true},
				"source_id": schema.StringAttribute{Computed: true},
			},
		},
		"string_value": schema.StringAttribute{Computed: true},
	}
}

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Role.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Required: true},
			"name":            schema.StringAttribute{Computed: true},
			"description":     schema.StringAttribute{Computed: true},
			"enabled":         schema.BoolAttribute{Computed: true},
			"requestable":     schema.BoolAttribute{Computed: true},
			"dimensional":     schema.BoolAttribute{Computed: true},
			"owner":           computedRef(),
			"access_profiles": computedRefSet(),
			"entitlements":    computedRefSet(),
			"dimension_refs":  computedRefSet(),
			"segments": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"additional_owners": computedRefSet(),
			"membership": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"criteria": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"operation": schema.StringAttribute{Computed: true},
							"key": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"type":      schema.StringAttribute{Computed: true},
									"property":  schema.StringAttribute{Computed: true},
									"source_id": schema.StringAttribute{Computed: true},
								},
							},
							"string_value": schema.StringAttribute{Computed: true},
							"children": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operation": schema.StringAttribute{Computed: true},
										"key": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"type":      schema.StringAttribute{Computed: true},
												"property":  schema.StringAttribute{Computed: true},
												"source_id": schema.StringAttribute{Computed: true},
											},
										},
										"string_value": schema.StringAttribute{Computed: true},
										"children": schema.ListNestedAttribute{
											Computed: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: computedCriteriaLeaf(),
											},
										},
									},
								},
							},
						},
					},
					"identities": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id":         schema.StringAttribute{Computed: true},
								"type":       schema.StringAttribute{Computed: true},
								"name":       schema.StringAttribute{Computed: true},
								"alias_name": schema.StringAttribute{Computed: true},
							},
						},
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
					"approval_schemes": computedApprovalSchemes(),
				},
			},
			"revoke_request_config": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"comments_required":        schema.BoolAttribute{Computed: true},
					"denial_comments_required": schema.BoolAttribute{Computed: true},
					"approval_schemes":         computedApprovalSchemes(),
				},
			},
			"created":  schema.StringAttribute{Computed: true},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state roleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Debug(ctx, "Reading role data source", map[string]any{"id": id})

	apiResp, err := d.client.GetRole(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Role",
			fmt.Sprintf("Could not read role %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Role", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
