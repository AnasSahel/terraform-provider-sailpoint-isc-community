// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &managedClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &managedClusterDataSource{}
)

type managedClusterDataSource struct {
	client *client.Client
}

func NewManagedClusterDataSource() datasource.DataSource {
	return &managedClusterDataSource{}
}

func (d *managedClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

func (d *managedClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "managed cluster data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *managedClusterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for a SailPoint managed cluster (virtual appliance cluster).",
		MarkdownDescription: "Data source for a SailPoint managed cluster (virtual appliance cluster). Look up a cluster by `id` " +
			"or by `name`. Exactly one of `id` or `name` must be provided. Typical use: feed `id` into " +
			"`sailpoint_source.cluster` instead of hardcoding it. Key material (key pair, keystore, public keys) " +
			"is deliberately not exposed, so it never lands in state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the managed cluster. Mutually exclusive with `name`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The managed cluster name (exact, case-sensitive match). Use as a lookup key (mutually exclusive with `id`) or read the resolved name after lookup by `id`. Returns an error if zero or several clusters match.",
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the managed cluster.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The cluster type (for example `idn` or `iai`).",
				Computed:            true,
			},
			"client_type": schema.StringAttribute{
				MarkdownDescription: "The type of client in the cluster (for example `CCG` or `VA`).",
				Computed:            true,
			},
			"pod": schema.StringAttribute{
				MarkdownDescription: "The SailPoint pod hosting the cluster.",
				Computed:            true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The tenant (org) the cluster belongs to.",
				Computed:            true,
			},
			"ccg_version": schema.StringAttribute{
				MarkdownDescription: "The CCG version used by the cluster.",
				Computed:            true,
			},
			"operational": schema.BoolAttribute{
				MarkdownDescription: "Whether the cluster is operational.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The cluster status: `CONFIGURING`, `FAILED`, `NO_CLIENTS`, `NORMAL` or `WARNING`.",
				Computed:            true,
			},
			"client_ids": schema.ListAttribute{
				MarkdownDescription: "IDs of the clients (virtual appliances) in the cluster.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"service_count": schema.Int64Attribute{
				MarkdownDescription: "Number of services bound to the cluster.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the cluster was created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the cluster was last updated.",
				Computed:            true,
			},
		},
	}
}

func (d *managedClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config managedClusterModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueString() != ""
	hasName := !config.Name.IsNull() && !config.Name.IsUnknown() && config.Name.ValueString() != ""

	if !hasID && !hasName {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"One of `id` or `name` must be provided.",
		)
		return
	}
	if hasID && hasName {
		resp.Diagnostics.AddError(
			"Conflicting arguments",
			"`id` and `name` are mutually exclusive. Provide only one.",
		)
		return
	}

	var apiResp *client.ManagedClusterAPI

	if hasID {
		tflog.Debug(ctx, "Reading managed cluster data source by ID", map[string]any{"id": config.ID.ValueString()})
		cluster, err := d.client.GetManagedCluster(ctx, config.ID.ValueString())
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError("No managed cluster found", fmt.Sprintf("No managed cluster has id %q.", config.ID.ValueString()))
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Managed Cluster",
				fmt.Sprintf("Could not read managed cluster %q: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}
		apiResp = cluster
	} else {
		clusterName := config.Name.ValueString()
		tflog.Debug(ctx, "Reading managed cluster data source by name", map[string]any{"name": clusterName})

		// The name filter of /managed-clusters/v1 is not reliable, so list every
		// cluster (a tenant has a handful) and match the name here.
		clusters, err := d.client.ListManagedClusters(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Managed Clusters",
				fmt.Sprintf("Could not list managed clusters: %s", err.Error()),
			)
			return
		}

		matches := matchClustersByName(clusters, clusterName)
		switch len(matches) {
		case 0:
			resp.Diagnostics.AddError(
				"No managed cluster found",
				fmt.Sprintf("No managed cluster matched name %q. Verify the cluster name (the match is exact and case-sensitive).", clusterName),
			)
			return
		case 1:
			apiResp = &matches[0]
		default:
			resp.Diagnostics.AddError(
				"Multiple managed clusters found",
				fmt.Sprintf("name %q matched %d managed clusters. Use `id` instead.", clusterName, len(matches)),
			)
			return
		}
	}

	var state managedClusterModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully read managed cluster data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func matchClustersByName(clusters []client.ManagedClusterAPI, name string) []client.ManagedClusterAPI {
	var matches []client.ManagedClusterAPI
	for _, c := range clusters {
		if c.Name == name {
			matches = append(matches, c)
		}
	}
	return matches
}
