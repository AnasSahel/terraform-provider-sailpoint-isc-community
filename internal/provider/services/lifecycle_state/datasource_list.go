// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

type LifecycleStateListDataSource struct {
	client *sailpoint.APIClient
}

var (
	_ datasource.DataSource              = &LifecycleStateListDataSource{}
	_ datasource.DataSourceWithConfigure = &LifecycleStateListDataSource{}
)

func NewLifecycleStateListDataSource() datasource.DataSource {
	return &LifecycleStateListDataSource{}
}

func (d *LifecycleStateListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *sailpoint.APIClient. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

func (d *LifecycleStateListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state_list"
}

func (d *LifecycleStateListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = LifecycleStateListDataSourceSchema()
}

func (d *LifecycleStateListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config LifecycleStateListDataSourceModel

	// Read the existing config
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate required fields
	if config.IdentityProfileId.IsNull() || config.IdentityProfileId.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing Identity Profile ID",
			"The identity profile ID must be specified.",
		)
		return
	}

	// Call the API
	lifecycleStateList, httpResponse, err := d.client.V2025.LifecycleStatesAPI.GetLifecycleStates(
		ctx,
		config.IdentityProfileId.ValueString(),
	).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lifecycle States",
			fmt.Sprintf("Could not read lifecycle states for identity profile ID '%s': %s\n\nHTTP Response: %v",
				config.IdentityProfileId.ValueString(),
				err.Error(),
				httpResponse,
			),
		)
		return
	}

	// Transform API response to Terraform models
	data := LifecycleStateListDataSourceModel{
		IdentityProfileId:  config.IdentityProfileId,
		LifecycleStateList: make([]LifecycleStateModel, 0, len(lifecycleStateList)),
	}

	for _, lifecycleStateItem := range lifecycleStateList {
		data.LifecycleStateList = append(data.LifecycleStateList, ToTerraformDataSource(ctx, lifecycleStateItem))
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
