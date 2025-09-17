package managedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type ManagedClusterResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	Configuration types.Map    `tfsdk:"configuration"`
}

func (r ManagedClusterResourceModel) ToSailPointCreateManagedClusterRequest(ctx context.Context) (*api_v2025.ManagedClusterRequest, diag.Diagnostics) {
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

	if !r.Configuration.IsNull() && !r.Configuration.IsUnknown() {
		configuration := make(map[string]string)
		diags := r.Configuration.ElementsAs(ctx, &configuration, false)
		if diags.HasError() {
			return nil, diags
		}
		managedClusterRequest.SetConfiguration(configuration)
	}

	return managedClusterRequest, nil
}

func (r *ManagedClusterResourceModel) FromSailPointManagedCluster(ctx context.Context, apiModel *api_v2025.ManagedCluster) diag.Diagnostics {
	conf, diagsConf := types.MapValueFrom(ctx, types.StringType, apiModel.GetConfiguration())
	// clientIds, diagsClientIds := types.ListValueFrom(ctx, types.StringType, apiModel.GetClientIds())

	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Type = types.StringValue(string(apiModel.GetType()))
	r.Configuration = conf
	r.Description = types.StringValue(apiModel.GetDescription())

	// return append(diagsConf, diagsClientIds...)
	return diagsConf
}
