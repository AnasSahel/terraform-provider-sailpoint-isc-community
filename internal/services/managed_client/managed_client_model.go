// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_client

import (
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type managedClientModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ClusterID         types.String `tfsdk:"cluster_id"`
	ClientID          types.String `tfsdk:"client_id"`
	Description       types.String `tfsdk:"description"`
	Type              types.String `tfsdk:"type"`
	ClusterType       types.String `tfsdk:"cluster_type"`
	Status            types.String `tfsdk:"status"`
	IPAddress         types.String `tfsdk:"ip_address"`
	VaVersion         types.String `tfsdk:"va_version"`
	LastSeen          types.String `tfsdk:"last_seen"`
	SinceLastSeen     types.String `tfsdk:"since_last_seen"`
	ProvisionStatus   types.String `tfsdk:"provision_status"`
	APIGatewayBaseURL types.String `tfsdk:"api_gateway_base_url"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

// FromAPI maps a managed client API response into the Terraform model.
func (m *managedClientModel) FromAPI(api *client.ManagedClientAPI) {
	if api == nil {
		return
	}

	m.ID = types.StringValue(api.ID)
	m.Name = common.StringOrNull(api.Name)
	m.ClusterID = common.StringOrNullIfEmpty(api.ClusterID)
	m.ClientID = common.StringOrNullIfEmpty(api.ClientID)
	m.Description = common.StringOrNull(api.Description)
	m.Type = common.StringOrNull(api.Type)
	m.ClusterType = common.StringOrNull(api.ClusterType)
	m.Status = common.StringOrNull(api.Status)
	m.IPAddress = common.StringOrNull(api.IPAddress)
	m.VaVersion = common.StringOrNull(api.VaVersion)
	m.LastSeen = common.StringOrNull(api.LastSeen)
	m.SinceLastSeen = common.StringOrNull(api.SinceLastSeen)
	m.ProvisionStatus = common.StringOrNull(api.ProvisionStatus)
	m.APIGatewayBaseURL = common.StringOrNull(api.APIGatewayBaseURL)
	m.CreatedAt = common.StringOrNull(api.CreatedAt)
	m.UpdatedAt = common.StringOrNull(api.UpdatedAt)
}

// matchManagedClients returns the managed clients whose name equals name and,
// when clusterID is not empty, that belong to that cluster.
func matchManagedClients(all []client.ManagedClientAPI, name, clusterID string) []client.ManagedClientAPI {
	var matches []client.ManagedClientAPI
	for _, mc := range all {
		if mc.Name == nil || *mc.Name != name {
			continue
		}
		if clusterID != "" && mc.ClusterID != clusterID {
			continue
		}
		matches = append(matches, mc)
	}
	return matches
}
