// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type managedClusterModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	Description                        types.String `tfsdk:"description"`
	Type                               types.String `tfsdk:"type"`
	ClientType                         types.String `tfsdk:"client_type"`
	Pod                                types.String `tfsdk:"pod"`
	Org                                types.String `tfsdk:"org"`
	CcgVersion                         types.String `tfsdk:"ccg_version"`
	Configuration                      types.Map    `tfsdk:"configuration"`
	PinnedConfig                       types.Bool   `tfsdk:"pinned_config"`
	Operational                        types.Bool   `tfsdk:"operational"`
	Status                             types.String `tfsdk:"status"`
	AlertKey                           types.String `tfsdk:"alert_key"`
	ClientIDs                          types.List   `tfsdk:"client_ids"`
	ServiceCount                       types.Int64  `tfsdk:"service_count"`
	CurrentInstalledReleaseVersion     types.String `tfsdk:"current_installed_release_version"`
	ConsolidatedHealthIndicatorsStatus types.String `tfsdk:"consolidated_health_indicators_status"`
	CreatedAt                          types.String `tfsdk:"created_at"`
	UpdatedAt                          types.String `tfsdk:"updated_at"`
}

func (m *managedClusterModel) FromAPI(ctx context.Context, api *client.ManagedClusterAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if api == nil {
		return diagnostics
	}

	m.ID = types.StringValue(api.ID)
	m.Name = common.StringOrNullIfEmpty(api.Name)
	m.Description = common.StringOrNull(api.Description)
	m.Type = common.StringOrNullIfEmpty(api.Type)
	m.ClientType = common.StringOrNull(api.ClientType)
	m.Pod = common.StringOrNullIfEmpty(api.Pod)
	m.Org = common.StringOrNullIfEmpty(api.Org)
	m.CcgVersion = common.StringOrNullIfEmpty(api.CcgVersion)
	m.Status = common.StringOrNull(api.Status)
	m.AlertKey = common.StringOrNull(api.AlertKey)
	m.CurrentInstalledReleaseVersion = common.StringOrNull(api.CurrentInstalledReleaseVersion)
	m.ConsolidatedHealthIndicatorsStatus = common.StringOrNull(api.ConsolidatedHealthIndicatorsStatus)
	m.CreatedAt = common.StringOrNull(api.CreatedAt)
	m.UpdatedAt = common.StringOrNull(api.UpdatedAt)

	if api.PinnedConfig != nil {
		m.PinnedConfig = types.BoolValue(*api.PinnedConfig)
	} else {
		m.PinnedConfig = types.BoolNull()
	}
	if api.Operational != nil {
		m.Operational = types.BoolValue(*api.Operational)
	} else {
		m.Operational = types.BoolNull()
	}
	if api.ServiceCount != nil {
		m.ServiceCount = types.Int64Value(int64(*api.ServiceCount))
	} else {
		m.ServiceCount = types.Int64Null()
	}

	configuration := make(map[string]types.String, len(api.Configuration))
	for k, v := range api.Configuration {
		configuration[k] = common.StringOrNull(v)
	}
	configValue, diags := types.MapValueFrom(ctx, types.StringType, configuration)
	diagnostics.Append(diags...)
	m.Configuration = configValue

	clientIDs := api.ClientIDs
	if clientIDs == nil {
		clientIDs = []string{}
	}
	clientIDsValue, diags := types.ListValueFrom(ctx, types.StringType, clientIDs)
	diagnostics.Append(diags...)
	m.ClientIDs = clientIDsValue

	return diagnostics
}
