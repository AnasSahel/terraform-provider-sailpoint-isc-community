// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_profile

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
	_ datasource.DataSource              = &identityProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &identityProfileDataSource{}
)

type identityProfileDataSource struct {
	client *client.Client
}

// NewIdentityProfileDataSource creates a new data source for Identity Profile.
func NewIdentityProfileDataSource() datasource.DataSource {
	return &identityProfileDataSource{}
}

// Metadata implements datasource.DataSource.
func (d *identityProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_profile"
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *identityProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "identity profile data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

// Schema implements datasource.DataSource.
func (d *identityProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a SailPoint Identity Profile by ID.",
		MarkdownDescription: "Retrieves a SailPoint Identity Profile by ID. Identity profiles define the source of identities and how identity attributes are mapped.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the identity profile.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity profile.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the identity profile was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the identity profile was last modified.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the identity profile.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the identity profile.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner object. Always `IDENTITY`.",
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
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of the identity profile.",
				Computed:            true,
			},
			"authoritative_source": schema.SingleNestedAttribute{
				MarkdownDescription: "The authoritative source for the identity profile.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the source object. Always `SOURCE`.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the authoritative source.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the authoritative source.",
						Computed:            true,
					},
				},
			},
			"identity_refresh_required": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether an identity refresh is required.",
				Computed:            true,
			},
			"identity_count": schema.Int32Attribute{
				MarkdownDescription: "The number of identities belonging to this identity profile.",
				Computed:            true,
			},
			"identity_attribute_config": schema.SingleNestedAttribute{
				MarkdownDescription: "The identity attribute configuration that defines how identity attributes are mapped.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the identity attribute configuration is enabled.",
						Computed:            true,
					},
					"attribute_transforms": schema.ListNestedAttribute{
						MarkdownDescription: "List of identity attribute transforms.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identity_attribute_name": schema.StringAttribute{
									MarkdownDescription: "The name of the identity attribute being mapped.",
									Computed:            true,
								},
								"transform_definition": schema.SingleNestedAttribute{
									MarkdownDescription: "The transform definition for the identity attribute.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "The type of the transform definition (e.g., `accountAttribute`, `rule`).",
											Computed:            true,
										},
										"attributes": schema.StringAttribute{
											MarkdownDescription: "The attributes of the transform definition as a JSON string.",
											Computed:            true,
											CustomType:          jsontypes.NormalizedType{},
										},
									},
								},
							},
						},
					},
				},
			},
			"identity_exception_report_reference": schema.SingleNestedAttribute{
				MarkdownDescription: "Reference to the identity exception report.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"task_result_id": schema.StringAttribute{
						MarkdownDescription: "The task result ID of the identity exception report.",
						Computed:            true,
					},
					"report_name": schema.StringAttribute{
						MarkdownDescription: "The name of the identity exception report.",
						Computed:            true,
					},
				},
			},
			"has_time_based_attr": schema.BoolAttribute{
				MarkdownDescription: "Indicates the value of `requiresPeriodicRefresh` attribute for the identity profile.",
				Computed:            true,
			},
		},
	}
}

// Read implements datasource.DataSource.
func (d *identityProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	tflog.Debug(ctx, "Getting config for identity profile data source")
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &config.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := config.ID.ValueString()

	// Read the identity profile from SailPoint
	tflog.Debug(ctx, "Fetching identity profile from SailPoint", map[string]any{
		"id": identityProfileID,
	})
	identityProfileResponse, err := d.client.GetIdentityProfile(ctx, identityProfileID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Profile",
			fmt.Sprintf("Could not read SailPoint Identity Profile %q: %s", identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Identity Profile", map[string]any{
			"id":    identityProfileID,
			"error": err.Error(),
		})
		return
	}

	if identityProfileResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Identity Profile",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the data source model
	var state identityProfileDataSourceModel
	tflog.Debug(ctx, "Mapping SailPoint Identity Profile API response to data source model", map[string]any{
		"id": identityProfileID,
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *identityProfileResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Identity Profile data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}
