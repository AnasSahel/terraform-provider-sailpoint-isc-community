// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &lifecycleStateDataSource{}
	_ datasource.DataSourceWithConfigure = &lifecycleStateDataSource{}
)

type lifecycleStateDataSource struct {
	client *client.Client
}

func NewLifecycleStateDataSource() datasource.DataSource {
	return &lifecycleStateDataSource{}
}

func (d *lifecycleStateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *lifecycleStateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

func (d *lifecycleStateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.LifecycleStateSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Lifecycle State.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Lifecycle State within an Identity Profile. Lifecycle States define the stages an identity goes through, such as joiner, mover, or leaver states. See [Lifecycle State Documentation](https://developer.sailpoint.com/docs/api/v2024/get-lifecycle-state) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *lifecycleStateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading LifecycleState data source")

	var config models.LifecycleState
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the lifecycle state via API using both identity_profile_id and id
	fetchedLifecycleState, err := d.client.GetLifecycleState(ctx, config.IdentityProfileID.ValueString(), config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lifecycle State",
			fmt.Sprintf("Could not read lifecycle state ID %s from identity profile %s: %s",
				config.ID.ValueString(), config.IdentityProfileID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model (includeNull: false for data sources)
	var state models.LifecycleState
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedLifecycleState, config.IdentityProfileID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Lifecycle State Response",
			fmt.Sprintf("Could not convert lifecycle state response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Lifecycle State data source read successfully", map[string]interface{}{
		"lifecycle_state_id":  state.ID.ValueString(),
		"identity_profile_id": state.IdentityProfileID.ValueString(),
	})
}
