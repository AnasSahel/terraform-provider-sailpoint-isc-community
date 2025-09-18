// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type ManagedClusterResource struct {
	client *api_v2025.APIClient
}

var (
	_ resource.Resource              = &ManagedClusterResource{}
	_ resource.ResourceWithConfigure = &ManagedClusterResource{}
)

func NewManagedClusterResource() resource.Resource {
	return &ManagedClusterResource{}
}

func (r *ManagedClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Managed Cluster Resource client",
		map[string]any{"provider_data_type": fmt.Sprintf("%T", req.ProviderData)})

	if req.ProviderData == nil {
		resp.Diagnostics.AddError(
			"Unable to Configure Managed Cluster Resource",
			"Provider data is nil. This usually means that the provider hasn't been configured correctly.",
		)
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unable to Configure Managed Cluster Resource",
			"Provider data is of unexpected type. Expected '*api_v2025.APIClient'.",
		)
		return
	}

	r.client = client
	tflog.Info(ctx, "Managed Cluster Resource client configured successfully")
}

func (r *ManagedClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

func (r *ManagedClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ManagedClusterResourceModel

	tflog.Info(ctx, "Reading Terraform plan data into the model")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read plan data into the model")

	// Set data
	tflog.Info(ctx, "Mapping Terraform plan data to API request model")
	managedClusterRequest, diags := data.toSailPointCreateManagedClusterRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully mapped Terraform plan data to API request model")

	tflog.Info(ctx, "Creating Managed Cluster via API")
	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.CreateManagedCluster(context.Background()).ManagedClusterRequest(*managedClusterRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Managed Cluster",
			fmt.Sprintf("An error was encountered creating the Managed Cluster: %s. HTTP Response: %v", err.Error(), httpResponse),
		)
		return
	}
	tflog.Info(ctx, "Successfully created Managed Cluster via API")
	resp.Diagnostics.Append(data.fromSailPointManagedCluster(ctx, managedCluster)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error mapping API response to Terraform model")
		return
	}
	tflog.Info(ctx, "Successfully mapped API response to Terraform model")

	// Set state
	tflog.Info(ctx, "Setting state with the new Managed Cluster data")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error setting state with the new Managed Cluster data")
		return
	}
	tflog.Info(ctx, "Successfully set state with the new Managed Cluster data")
}

func (r *ManagedClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ManagedClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ManagedClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
