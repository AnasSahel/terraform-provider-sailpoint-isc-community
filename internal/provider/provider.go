// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	connector_datasource "github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/services/connector/datasource"
	connector_resource "github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/services/connector/resource"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/services/managedcluster"
	transform_datasource "github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/services/transform/datasource"
	transform_resource "github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/services/transform/resource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &sailpointProvider{}
)

// sailpointProviderModel maps provider schema data to a Go type.
type sailpointProviderModel struct {
	BaseUrl      types.String `tfsdk:"base_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sailpointProvider{
			version: version,
		}
	}
}

// sailpointProvider is the provider implementation.
type sailpointProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *sailpointProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sailpoint"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *sailpointProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a SailPoint API client for data sources and resources.
func (p *sailpointProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config sailpointProviderModel

	tflog.Info(ctx, "Configuring SailPoint client")

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing or invalid Base URL",
			"Ensure the Base URL attribute is set in the configuration.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing or invalid Client ID",
			"Ensure the Client ID attribute is set in the configuration.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing or invalid Client Secret",
			"Ensure the Client Secret attribute is set in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl := os.Getenv("SAILPOINT_BASE_URL")
	clientId := os.Getenv("SAILPOINT_CLIENT_ID")
	clientSecret := os.Getenv("SAILPOINT_CLIENT_SECRET")

	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	}
	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}
	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if baseUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseUrl"),
			"Missing Base URL",
			"The provider cannot create the API client as there is a missing value for the Base URL. "+
				"Set the Base URL attribute in the configuration or use the SAILPOINT_BASE_URL environment variable.",
		)
	}

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clientId"),
			"Missing Client ID",
			"The provider cannot create the API client as there is a missing value for the Client ID. "+
				"Set the Client ID attribute in the configuration or use the SAILPOINT_CLIENT_ID environment variable.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clientSecret"),
			"Missing Client Secret",
			"The provider cannot create the API client as there is a missing value for the Client Secret. "+
				"Set the Client Secret attribute in the configuration or use the SAILPOINT_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sailpoint_base_url", baseUrl)
	ctx = tflog.SetField(ctx, "sailpoint_client_id", clientId)
	ctx = tflog.SetField(ctx, "sailpoint_client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sailpoint_client_secret")

	os.Setenv("SAIL_BASE_URL", baseUrl)
	os.Setenv("SAIL_CLIENT_ID", clientId)
	os.Setenv("SAIL_CLIENT_SECRET", clientSecret)
	sailpointConfiguration := sailpoint.NewDefaultConfiguration()

	sailpointClient := sailpoint.NewAPIClient(sailpointConfiguration).V2025

	resp.DataSourceData = sailpointClient
	resp.ResourceData = sailpointClient

	tflog.Info(ctx, "SailPoint client configured", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *sailpointProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		transform_datasource.NewTransformsDataSource,
		transform_datasource.NewTransformDataSource,
		managedcluster.NewManagedClusterDataSource,
		connector_datasource.NewConnectorsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *sailpointProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		transform_resource.NewTransformResource,
		managedcluster.NewManagedClusterResource,
		connector_resource.NewConnectorResource,
	}
}
