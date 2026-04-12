// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package segment

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
	_ datasource.DataSource              = &segmentDataSource{}
	_ datasource.DataSourceWithConfigure = &segmentDataSource{}
)

type segmentDataSource struct {
	client *client.Client
}

// NewSegmentDataSource creates a new data source for SailPoint Segment.
func NewSegmentDataSource() datasource.DataSource {
	return &segmentDataSource{}
}

func (d *segmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (d *segmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "segment data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *segmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Segment.",
		MarkdownDescription: "Data source for SailPoint Segment. Look up a segment by ID to retrieve its configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the segment.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the segment.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the segment.",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the segment is operational.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the segment.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"visibility_criteria": schema.SingleNestedAttribute{
				MarkdownDescription: "Visibility rules that determine which identities the segment applies to.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"expression": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"operator":  schema.StringAttribute{Computed: true},
							"attribute": schema.StringAttribute{Computed: true},
							"value": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"type":  schema.StringAttribute{Computed: true},
									"value": schema.StringAttribute{Computed: true},
								},
							},
							"children": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operator":  schema.StringAttribute{Computed: true},
										"attribute": schema.StringAttribute{Computed: true},
										"value": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"type":  schema.StringAttribute{Computed: true},
												"value": schema.StringAttribute{Computed: true},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the segment was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the segment was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *segmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state segmentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Debug(ctx, "Reading segment data source", map[string]any{"id": id})

	apiResp, err := d.client.GetSegment(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Segment",
			fmt.Sprintf("Could not read SailPoint Segment %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Segment", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
