// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &sailpointProvider{}
)

// sailpointProvider is the provider implementation.
type sailpointProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

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
	tflog.Info(ctx, "Configuring SailPoint client")

	var config sailpointProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...) // Get the config data
	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("base_url"), "Invalid Base URL", "Base URL must be configured.")
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("client_id"), "Invalid Client ID", "Client ID must be configured.")
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("client_secret"), "Invalid Client Secret", "Client Secret must be configured.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Use environment variables as fallback if config not set
	var (
		baseUrl      = os.Getenv("SAILPOINT_BASE_URL")
		clientId     = os.Getenv("SAILPOINT_CLIENT_ID")
		clientSecret = os.Getenv("SAILPOINT_CLIENT_SECRET")
	)

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
		resp.Diagnostics.AddAttributeError(path.Root("baseUrl"), "Missing Base URL", "Set base_url in config or SAILPOINT_BASE_URL environment variable.")
	}

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(path.Root("clientId"), "Missing Client ID", "Set client_id in config or SAILPOINT_CLIENT_ID environment variable.")
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(path.Root("clientSecret"), "Missing Client Secret", "Set client_secret in config or SAILPOINT_CLIENT_SECRET environment variable.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sailpoint_base_url", baseUrl)
	ctx = tflog.SetField(ctx, "sailpoint_client_id", clientId)
	ctx = tflog.SetField(ctx, "sailpoint_client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sailpoint_client_secret")

	tflog.Debug(ctx, "Creating SailPoint client")

	apiClient, err := client.NewClient(baseUrl, clientId, clientSecret)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SailPoint Client", fmt.Sprintf("An error occurred creating the SailPoint client: %s", err.Error()))
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient

	tflog.Info(ctx, "Configured SailPoint client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *sailpointProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *sailpointProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
