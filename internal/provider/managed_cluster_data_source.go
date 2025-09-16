// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

var (
	_ datasource.DataSource              = &ManagedClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &ManagedClusterDataSource{}
)

type ManagedClusterDataSource struct {
	client *api_v2025.APIClient
}

type ManagedClusterDataSourceModel struct {
	ID                   types.String             `tfsdk:"id"`
	Name                 types.String             `tfsdk:"name"`
	Pod                  types.String             `tfsdk:"pod"`
	Org                  types.String             `tfsdk:"org"`
	Type                 types.String             `tfsdk:"type"`
	Configuration        types.Map                `tfsdk:"configuration"`
	KeyPair              ManagedClusterKeyPair    `tfsdk:"key_pair"`
	Attributes           ManagedClusterAttributes `tfsdk:"attributes"`
	Description          types.String             `tfsdk:"description"`
	Redis                ManagedClusterRedis      `tfsdk:"redis"`
	ClientType           types.String             `tfsdk:"client_type"`
	CcgVersion           types.String             `tfsdk:"ccg_version"`
	PinnedConfig         types.Bool               `tfsdk:"pinned_config"`
	Operational          types.Bool               `tfsdk:"operational"`
	Status               types.String             `tfsdk:"status"`
	PublicKeyCertificate types.String             `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint  types.String             `tfsdk:"public_key_thumbprint"`
	PublicKey            types.String             `tfsdk:"public_key_type"`
	AlertKey             types.String             `tfsdk:"alert_key"`
	ClientIds            types.List               `tfsdk:"client_ids"`
	ServiceCount         types.Int32              `tfsdk:"service_count"`
	CcId                 types.String             `tfsdk:"cc_id"`
	CreatedAt            types.String             `tfsdk:"created_at"`
	UpdatedAt            types.String             `tfsdk:"updated_at"`
}

type ManagedClusterKeyPair struct {
	PublicKey            types.String `tfsdk:"public_key"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
}

type ManagedClusterAttributes struct {
	Queue    ManagedClusterQueue `tfsdk:"queue"`
	Keystore types.String        `tfsdk:"keystore"`
}

type ManagedClusterQueue struct {
	Name   types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
}

type ManagedClusterRedis struct {
	RedisHost types.String `tfsdk:"redis_host"`
	RedisPort types.Int32  `tfsdk:"redis_port"`
}

func NewManagedClusterDataSource() datasource.DataSource {
	return &ManagedClusterDataSource{}
}

func (d *ManagedClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ManagedClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_managed_cluster", req.ProviderTypeName)
}

func (d *ManagedClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Required: true},
			"name": schema.StringAttribute{Computed: true},
			"pod":  schema.StringAttribute{Computed: true},
			"org":  schema.StringAttribute{Computed: true},
			"type": schema.StringAttribute{Computed: true},
			"configuration": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"key_pair": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"public_key":             schema.StringAttribute{Computed: true},
					"public_key_thumbprint":  schema.StringAttribute{Computed: true},
					"public_key_certificate": schema.StringAttribute{Computed: true},
				},
			},
			"attributes": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"queue": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"name":   schema.StringAttribute{Computed: true},
							"region": schema.StringAttribute{Computed: true},
						},
					},
					"keystore": schema.StringAttribute{Computed: true},
				},
			},
			"description": schema.StringAttribute{Computed: true},
			"redis": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"redis_host": schema.StringAttribute{Computed: true},
					"redis_port": schema.Int32Attribute{Computed: true},
				},
			},
			"client_type":            schema.StringAttribute{Computed: true},
			"ccg_version":            schema.StringAttribute{Computed: true},
			"pinned_config":          schema.BoolAttribute{Computed: true},
			"operational":            schema.BoolAttribute{Computed: true},
			"status":                 schema.StringAttribute{Computed: true},
			"public_key_certificate": schema.StringAttribute{Computed: true},
			"public_key_thumbprint":  schema.StringAttribute{Computed: true},
			"public_key_type":        schema.StringAttribute{Computed: true},
			"alert_key":              schema.StringAttribute{Computed: true},
			"client_ids":             schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"service_count":          schema.Int32Attribute{Computed: true},
			"cc_id":                  schema.StringAttribute{Computed: true},
			"created_at":             schema.StringAttribute{Computed: true},
			"updated_at":             schema.StringAttribute{Computed: true},
		},
	}
}

func (d *ManagedClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManagedClusterDataSourceModel
	diags := req.Config.GetAttribute(ctx, path.Root("id"), &data.ID)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Configuration Error",
			"Failed to read configuration data. Please ensure the 'id' attribute is set correctly.",
		)
		return
	}

	managedCluster, httpResponse, err := d.client.ManagedClustersAPI.GetManagedCluster(context.Background(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Managed Cluster",
			fmt.Sprintf("An error occurred while reading the managed cluster: %s. Response: %v", err.Error(), httpResponse),
		)
		return
	}

	// Set data
	data.ID = types.StringValue(managedCluster.GetId())
	data.Name = types.StringValue(managedCluster.GetName())
	data.Pod = types.StringValue(managedCluster.GetPod())
	data.Org = types.StringValue(managedCluster.GetOrg())
	data.Type = types.StringValue(string(managedCluster.GetType()))

	data.Configuration, diags = formatConfigurationToTerraform(managedCluster.GetConfiguration())
	resp.Diagnostics.Append(diags...)

	data.KeyPair = ManagedClusterKeyPair{
		PublicKey:            types.StringValue(*managedCluster.GetKeyPair().PublicKey.Get()),
		PublicKeyThumbprint:  types.StringValue(*managedCluster.GetKeyPair().PublicKeyThumbprint.Get()),
		PublicKeyCertificate: types.StringValue(*managedCluster.GetKeyPair().PublicKeyCertificate.Get()),
	}

	data.Attributes = ManagedClusterAttributes{
		Queue: ManagedClusterQueue{
			Name:   types.StringValue(managedCluster.GetAttributes().Queue.GetName()),
			Region: types.StringValue(managedCluster.GetAttributes().Queue.GetRegion()),
		},
	}
	if managedCluster.GetAttributes().Keystore.IsSet() {
		data.Attributes.Keystore = types.StringValue(*managedCluster.GetAttributes().Keystore.Get())
	}

	data.Description = types.StringValue(managedCluster.GetDescription())

	data.Redis = ManagedClusterRedis{
		RedisHost: types.StringValue(*managedCluster.GetRedis().RedisHost),
		RedisPort: types.Int32Value(*managedCluster.GetRedis().RedisPort),
	}

	data.ClientType = types.StringValue(string(managedCluster.GetClientType()))
	data.CcgVersion = types.StringValue(managedCluster.GetCcgVersion())
	data.PinnedConfig = types.BoolValue(managedCluster.GetPinnedConfig())

	data.Operational = types.BoolValue(managedCluster.GetOperational())
	data.Status = types.StringValue(string(managedCluster.GetStatus()))
	data.PublicKeyCertificate = types.StringValue(*managedCluster.GetKeyPair().PublicKeyCertificate.Get())
	data.PublicKeyThumbprint = types.StringValue(*managedCluster.GetKeyPair().PublicKeyThumbprint.Get())
	data.PublicKey = types.StringValue(*managedCluster.GetKeyPair().PublicKey.Get())
	data.AlertKey = types.StringValue(managedCluster.GetAlertKey())
	if managedCluster.GetClientIds() != nil {
		clientIds := make([]attr.Value, len(managedCluster.GetClientIds()))
		for i, id := range managedCluster.GetClientIds() {
			clientIds[i] = types.StringValue(id)
		}
		data.ClientIds = types.ListValueMust(types.StringType, clientIds)
	} else {
		data.ClientIds = types.ListNull(types.StringType)
	}
	data.ServiceCount = types.Int32Value(int32(managedCluster.GetServiceCount()))
	data.CcId = types.StringValue(managedCluster.GetCcId())
	data.CreatedAt = types.StringValue(managedCluster.GetCreatedAt().String())
	data.UpdatedAt = types.StringValue(managedCluster.GetUpdatedAt().String())

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"State Error",
			"Failed to set state data after reading managed cluster.",
		)
		return
	}

}

func formatConfigurationToTerraform(config map[string]string) (types.Map, diag.Diagnostics) {
	if config == nil {
		return types.MapNull(types.StringType), nil
	}

	var elements = make(map[string]attr.Value, len(config))

	for key, value := range config {
		elements[key] = types.StringValue(value)
	}

	mapValue, diags := types.MapValue(types.StringType, elements)
	return mapValue, diags
}
