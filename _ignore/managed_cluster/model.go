// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type ManagedClusterModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Pod                  types.String `tfsdk:"pod"`
	Org                  types.String `tfsdk:"org"`
	Type                 types.String `tfsdk:"type"`
	Configuration        types.Map    `tfsdk:"configuration"`
	Description          types.String `tfsdk:"description"`
	ClientType           types.String `tfsdk:"client_type"`
	CcgVersion           types.String `tfsdk:"ccg_version"`
	PinnedConfig         types.Bool   `tfsdk:"pinned_config"`
	Operational          types.Bool   `tfsdk:"operational"`
	Status               types.String `tfsdk:"status"`
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKey            types.String `tfsdk:"public_key_type"`
	AlertKey             types.String `tfsdk:"alert_key"`
	ClientIds            types.List   `tfsdk:"client_ids"`
	ServiceCount         types.Int32  `tfsdk:"service_count"`
	CcId                 types.String `tfsdk:"cc_id"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

// type ManagedClusterDataSourceModel struct {
// 	ManagedClusterModel
// }

type ManagedClusterResourceModel struct {
	ManagedClusterModel
}

func (r ManagedClusterResourceModel) toSailPointCreateManagedClusterRequest(ctx context.Context) (*api_v2025.ManagedClusterRequest, diag.Diagnostics) {
	managedClusterRequest := api_v2025.NewManagedClusterRequestWithDefaults()

	if !r.Name.IsNull() {
		managedClusterRequest.SetName(r.Name.ValueString())
	}

	if !r.Type.IsNull() {
		managedClusterRequest.SetType(api_v2025.ManagedClusterTypes(r.Type.ValueString()))
	}

	if !r.Description.IsNull() {
		managedClusterRequest.SetDescription(r.Description.ValueString())
	}

	if !r.Configuration.IsNull() {
		configuration := make(map[string]string)
		diags := r.Configuration.ElementsAs(ctx, &configuration, false)
		if diags.HasError() {
			return nil, diags
		}
		managedClusterRequest.SetConfiguration(configuration)
	}

	return managedClusterRequest, nil
}

func (r *ManagedClusterResourceModel) fromSailPointManagedCluster(ctx context.Context, apiModel *api_v2025.ManagedCluster) diag.Diagnostics {
	conf, diagsConf := types.MapValueFrom(ctx, types.StringType, apiModel.GetConfiguration())
	clientIds, diagsClientIds := types.ListValueFrom(ctx, types.StringType, apiModel.GetClientIds())

	r.ID = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Pod = types.StringValue(apiModel.GetPod())
	r.Org = types.StringValue(apiModel.GetOrg())
	r.Type = types.StringValue(string(apiModel.GetType()))
	r.Configuration = conf
	r.Description = types.StringValue(apiModel.GetDescription())
	r.ClientType = types.StringValue(string(apiModel.GetClientType()))
	r.CcgVersion = types.StringValue(apiModel.GetCcgVersion())
	r.PinnedConfig = types.BoolValue(apiModel.GetPinnedConfig())
	r.Operational = types.BoolValue(apiModel.GetOperational())
	r.Status = types.StringValue(apiModel.GetStatus())
	r.PublicKeyCertificate = types.StringValue(apiModel.GetPublicKeyCertificate())
	r.PublicKeyThumbprint = types.StringValue(apiModel.GetPublicKeyThumbprint())
	r.PublicKey = types.StringValue(apiModel.GetPublicKey())
	r.AlertKey = types.StringValue(apiModel.GetAlertKey())
	r.ClientIds = clientIds
	r.ServiceCount = types.Int32Value(apiModel.GetServiceCount())
	r.CcId = types.StringValue(apiModel.GetCcId())
	r.CreatedAt = types.StringValue(apiModel.GetCreatedAt().String())
	r.UpdatedAt = types.StringValue(apiModel.GetUpdatedAt().String())

	return append(diagsConf, diagsClientIds...)
}
