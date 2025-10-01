// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

type LifecycleStateDataSource struct {
	client       *sailpoint.APIClient
	validator    *Validator
	errorHandler *ErrorHandler
}

var (
	_ datasource.DataSource              = &LifecycleStateDataSource{}
	_ datasource.DataSourceWithConfigure = &LifecycleStateDataSource{}
)

func NewLifecycleStateDataSource() datasource.DataSource {
	return &LifecycleStateDataSource{
		validator:    NewValidator(),
		errorHandler: NewErrorHandler(),
	}
}

func (d *LifecycleStateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)

	if !ok {
		resp.Diagnostics.Append(d.errorHandler.HandleConfigurationError(
			ErrUnexpectedConfigureType,
			MsgExpectedAPIClient,
		)...)
		return
	}

	d.client = client
}

func (d *LifecycleStateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

func (d *LifecycleStateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = LifecycleStateDataSourceSchema()
}

func (d *LifecycleStateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config LifecycleStateDataSourceModel

	// Read the existing config
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	lifecycleState, httpResponse, err := d.client.V2025.LifecycleStatesAPI.
		GetLifecycleState(ctx, config.IdentityProfileId.ValueString(), config.Id.ValueString()).
		Execute()

	if err != nil {
		resp.Diagnostics.Append(d.errorHandler.HandleAPIError(
			"Reading",
			err,
			httpResponse,
			fmt.Sprintf("profile ID: %s, state ID: %s",
				config.IdentityProfileId.ValueString(),
				config.Id.ValueString(),
			),
		)...)
		return
	}

	// Create the response data model
	data := LifecycleStateDataSourceModel{
		IdentityProfileId:   config.IdentityProfileId,
		LifecycleStateModel: ToTerraformDataSource(ctx, *lifecycleState),
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
