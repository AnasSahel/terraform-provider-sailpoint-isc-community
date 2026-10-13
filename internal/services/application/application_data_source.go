// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &applicationDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationDataSource{}
)

type applicationDataSource struct {
	client *client.Client
}

type applicationDSModel struct {
	applicationModel
}

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "application data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *applicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Application (source app).",
		MarkdownDescription: "Data source for SailPoint Application (source app). Look up an application by `id` " +
			"or by `name`. Exactly one of `id` or `name` must be provided. " +
			"This data source uses the experimental `/source-apps/v1` API (`X-SailPoint-Experimental: true`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the application. Mutually exclusive with `name`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The application name. Use as a lookup key (mutually exclusive with `id`) or read the resolved name after lookup by `id`. Returns an error if more than one application matches.",
			},
			"cloud_app_id": schema.StringAttribute{
				MarkdownDescription: "The deprecated cloud application identifier.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the application.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is enabled.",
				Computed:            true,
			},
			"provision_request_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether access requests are enabled for this application.",
				Computed:            true,
			},
			"match_all_accounts": schema.BoolAttribute{
				MarkdownDescription: "Whether the application matches all accounts on the backing source.",
				Computed:            true,
			},
			"app_center_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the application is visible in the request center.",
				Computed:            true,
			},
			"account_source": schema.SingleNestedAttribute{
				MarkdownDescription: "The SailPoint source that owns this application.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the backing SailPoint source.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The source reference type. Typically `SOURCE`.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the backing source.",
						Computed:            true,
					},
					"use_for_password_management": schema.BoolAttribute{
						MarkdownDescription: "Whether the backing source is used for password management.",
						Computed:            true,
					},
					"password_policies": schema.ListNestedAttribute{
						MarkdownDescription: "Password policies associated with the backing source.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{Computed: true},
								"id":   schema.StringAttribute{Computed: true},
								"name": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the application.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the application was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "When the application was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config applicationDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueString() != ""
	hasName := !config.Name.IsNull() && !config.Name.IsUnknown() && config.Name.ValueString() != ""

	if !hasID && !hasName {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"One of `id` or `name` must be provided.",
		)
		return
	}
	if hasID && hasName {
		resp.Diagnostics.AddError(
			"Conflicting arguments",
			"`id` and `name` are mutually exclusive. Provide only one.",
		)
		return
	}

	var apiResp *client.ApplicationAPI
	var err error

	if hasID {
		tflog.Debug(ctx, "Reading application data source by ID", map[string]any{"id": config.ID.ValueString()})
		apiResp, err = d.client.GetApplication(ctx, config.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Application",
				fmt.Sprintf("Could not read application %q: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}
	} else {
		appName := config.Name.ValueString()
		tflog.Debug(ctx, "Reading application data source by name", map[string]any{"name": appName})
		apps, listErr := d.client.ListApplications(ctx, applicationNameFilter(appName), 2)
		if listErr != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Applications",
				fmt.Sprintf("Could not list applications with name %q: %s", appName, listErr.Error()),
			)
			return
		}

		switch len(apps) {
		case 0:
			resp.Diagnostics.AddError(
				"No application found",
				fmt.Sprintf("No application matched name %q. Verify the application name.", appName),
			)
			return
		case 1:
			apiResp, err = d.client.GetApplication(ctx, apps[0].ID)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Reading SailPoint Application",
					fmt.Sprintf("Could not read application %q: %s", apps[0].ID, err.Error()),
				)
				return
			}
		default:
			resp.Diagnostics.AddError(
				"Multiple applications found",
				fmt.Sprintf("name %q matched more than one application. Use `id` directly or choose a more specific name.", appName),
			)
			return
		}
	}

	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Application", "Received nil response from SailPoint API")
		return
	}

	var state applicationDSModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully read application data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func applicationNameFilter(name string) string {
	return fmt.Sprintf(`name eq "%s"`, strings.ReplaceAll(name, `"`, `\"`))
}
