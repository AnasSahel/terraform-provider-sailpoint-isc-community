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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceProvisioningPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceProvisioningPolicyDataSource{}
)

type sourceProvisioningPolicyDataSource struct {
	client *client.Client
}

func NewSourceProvisioningPolicyDataSource() datasource.DataSource {
	return &sourceProvisioningPolicyDataSource{}
}

func (d *sourceProvisioningPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_provisioning_policy"
}

func (d *sourceProvisioningPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source provisioning policy data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *sourceProvisioningPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single SailPoint source provisioning policy.",
		MarkdownDescription: "Retrieves a single SailPoint source provisioning policy. " +
			"A provisioning policy defines the fields and transformations required for a specific operation type.",
		Attributes: map[string]schema.Attribute{
			// Input parameters
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the source to retrieve the provisioning policy from.",
				Required:            true,
			},
			"usage_type": schema.StringAttribute{
				MarkdownDescription: "The usage type of the provisioning policy (e.g., `CREATE`, `UPDATE`, `DELETE`, `ENABLE`, `DISABLE`, `ASSIGN`, `UNASSIGN`, `CREATE_GROUP`, `UPDATE_GROUP`, `DELETE_GROUP`, `REGISTER`, `CREATE_IDENTITY`, `UPDATE_IDENTITY`, `EDIT_GROUP`, `UNLOCK`, `CHANGE_PASSWORD`).",
				Required:            true,
			},

			// Output attributes
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the provisioning policy.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the provisioning policy.",
				Computed:            true,
			},
			"fields": schema.ListNestedAttribute{
				MarkdownDescription: "The list of fields defined by the provisioning policy.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the field.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the field.",
							Computed:            true,
						},
						"is_required": schema.BoolAttribute{
							MarkdownDescription: "Whether the field is required.",
							Computed:            true,
						},
						"is_multi_valued": schema.BoolAttribute{
							MarkdownDescription: "Whether the field supports multiple values.",
							Computed:            true,
						},
						"transform": schema.StringAttribute{
							MarkdownDescription: "The transformation applied to the field as a JSON object.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
						"attributes": schema.StringAttribute{
							MarkdownDescription: "Additional attributes for the field as a JSON object.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

func (d *sourceProvisioningPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		SourceID  types.String `tfsdk:"source_id"`
		UsageType types.String `tfsdk:"usage_type"`
	}
	tflog.Debug(ctx, "Getting config for source provisioning policy data source")
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("source_id"), &config.SourceID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("usage_type"), &config.UsageType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := config.SourceID.ValueString()
	usageType := config.UsageType.ValueString()

	tflog.Debug(ctx, "Fetching source provisioning policy from SailPoint", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})

	policy, err := d.client.GetProvisioningPolicy(ctx, sourceID, usageType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy",
			fmt.Sprintf("Could not read provisioning policy for source %q with usage type %q: %s", sourceID, usageType, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Source Provisioning Policy", map[string]any{
			"source_id":  sourceID,
			"usage_type": usageType,
			"error":      err.Error(),
		})
		return
	}

	if policy == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy",
			"Received nil response from SailPoint API",
		)
		return
	}

	var state sourceProvisioningPolicyModel
	resp.Diagnostics.Append(state.FromAPI(ctx, policy, sourceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully read SailPoint Source Provisioning Policy data source", map[string]any{
		"source_id":   sourceID,
		"usage_type":  usageType,
		"policy_name": state.Name.ValueString(),
	})
}
