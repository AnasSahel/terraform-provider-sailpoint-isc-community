// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

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
	_ datasource.DataSource              = &roleDimensionDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDimensionDataSource{}
)

type roleDimensionDataSource struct {
	client *client.Client
}

// NewRoleDimensionDataSource returns the sailpoint_role_dimension data source.
func NewRoleDimensionDataSource() datasource.DataSource { return &roleDimensionDataSource{} }

func (d *roleDimensionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_dimension"
}

func (d *roleDimensionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "role dimension data source")
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

func computedCriteriaKey() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"type":     schema.StringAttribute{Computed: true},
			"property": schema.StringAttribute{Computed: true},
		},
	}
}

func computedCriteriaLeaf() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"operation":    schema.StringAttribute{Computed: true},
		"key":          computedCriteriaKey(),
		"string_value": schema.StringAttribute{Computed: true},
	}
}

func (d *roleDimensionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for a SailPoint role dimension.",
		MarkdownDescription: "Reads a dimension of a SailPoint dynamic role by `role_id` and `id`.",
		Attributes: map[string]schema.Attribute{
			"role_id":         schema.StringAttribute{Required: true, MarkdownDescription: "ID of the dimensional role."},
			"id":              schema.StringAttribute{Required: true, MarkdownDescription: "ID of the dimension."},
			"name":            schema.StringAttribute{Computed: true},
			"description":     schema.StringAttribute{Computed: true},
			"owner":           computedRef(),
			"access_profiles": computedRefSet(),
			"entitlements":    computedRefSet(),
			"membership": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"criteria": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"operation":    schema.StringAttribute{Computed: true},
							"key":          computedCriteriaKey(),
							"string_value": schema.StringAttribute{Computed: true},
							"children": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operation":    schema.StringAttribute{Computed: true},
										"key":          computedCriteriaKey(),
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
				},
			},
			"created":  schema.StringAttribute{Computed: true},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *roleDimensionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config roleDimensionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := config.RoleID.ValueString()
	id := config.ID.ValueString()
	tflog.Debug(ctx, "Reading role dimension data source", map[string]any{"role_id": roleID, "id": id})

	apiResp, err := d.client.GetRoleDimension(ctx, roleID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Role Dimension",
			fmt.Sprintf("Could not read dimension %q of role %q: %s", id, roleID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Role Dimension", "Received nil response from SailPoint API")
		return
	}

	var state roleDimensionModel
	resp.Diagnostics.Append(state.FromAPI(ctx, roleID, apiResp, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
