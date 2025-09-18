// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ManagedClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &ManagedClusterDataSource{}
)

// ManagedClusterDataSource defines the data source implementation.
type ManagedClusterDataSource struct {
	client *api_v2025.APIClient
}

// NewManagedClusterDataSource creates a new managed cluster data source.
func NewManagedClusterDataSource() datasource.DataSource {
	return &ManagedClusterDataSource{}
}

// Metadata returns the data source type name.
func (d *ManagedClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

// Schema defines the schema for the data source.
func (d *ManagedClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Trace(ctx, "Preparing ManagedClusterDataSource schema")
	resp.Schema = GetManagedClusterDataSourceSchema()
}

// Configure adds the provider configured client to the data source.
func (d *ManagedClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring ManagedClusterDataSource")

	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	tflog.Debug(ctx, "Configured ManagedClusterDataSource")
	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *ManagedClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Managed Cluster Data Source")

	var config ManagedClusterDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that we have either ID or name for lookup
	if config.Id.IsNull() && config.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a managed cluster.",
		)
		return
	}

	var managedCluster *api_v2025.ManagedCluster
	var err error
	var httpResponse interface{}

	if !config.Id.IsNull() {
		// Look up by ID
		clusterID := config.Id.ValueString()
		tflog.Debug(ctx, "Looking up managed cluster by ID", map[string]interface{}{
			"id": clusterID,
		})

		managedCluster, httpResponse, err = d.client.ManagedClustersAPI.GetManagedCluster(
			context.Background(),
			clusterID,
		).Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Managed Cluster",
				fmt.Sprintf("Could not read managed cluster with ID '%s': %s\n\nHTTP Response: %v",
					clusterID, err.Error(), httpResponse),
			)
			return
		}
	} else {
		// Look up by name - need to list all clusters and find by name
		clusterName := config.Name.ValueString()
		tflog.Debug(ctx, "Looking up managed cluster by name", map[string]interface{}{
			"name": clusterName,
		})

		clusters, httpResponse, err := d.client.ManagedClustersAPI.GetManagedClusters(
			context.Background(),
		).Execute()

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing Managed Clusters",
				fmt.Sprintf("Could not list managed clusters to find '%s': %s\n\nHTTP Response: %v",
					clusterName, err.Error(), httpResponse),
			)
			return
		}

		// Find the cluster with matching name
		var foundCluster *api_v2025.ManagedCluster
		for i := range clusters {
			if clusters[i].GetName() == clusterName {
				foundCluster = &clusters[i]
				break
			}
		}

		if foundCluster == nil {
			resp.Diagnostics.AddError(
				"Managed Cluster Not Found",
				fmt.Sprintf("No managed cluster found with name '%s'", clusterName),
			)
			return
		}

		managedCluster = foundCluster
	}

	// Convert API response to Terraform state
	var state ManagedClusterDataSourceModel
	conversionDiags := state.FromSailPointManagedClusterDataSource(ctx, managedCluster)
	resp.Diagnostics.Append(conversionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Successfully read managed cluster data", map[string]interface{}{
		"id":     state.Id.ValueString(),
		"name":   state.Name.ValueString(),
		"status": state.Status.ValueString(),
	})

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
