// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_client

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &managedClientDataSource{}
	_ datasource.DataSourceWithConfigure = &managedClientDataSource{}
)

type managedClientDataSource struct {
	client *client.Client
}

func NewManagedClientDataSource() datasource.DataSource {
	return &managedClientDataSource{}
}

func (d *managedClientDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_client"
}

func (d *managedClientDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "managed client data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *managedClientDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	computed := func(desc string) schema.StringAttribute {
		return schema.StringAttribute{MarkdownDescription: desc, Computed: true}
	}

	resp.Schema = schema.Schema{
		Description: "Data source for a SailPoint managed client (virtual appliance).",
		MarkdownDescription: "Data source for a SailPoint managed client, such as a virtual appliance (VA). " +
			"Look up a managed client by `id`, or by `name` optionally narrowed with `cluster_id`. " +
			"Exactly one of `id` or `name` must be provided.\n\n" +
			"The API cannot filter managed clients by name, so a lookup by `name` lists every managed client " +
			"and matches in the provider; it fails if zero or several managed clients match.\n\n" +
			"The client's API secret, alert key and VA download URL are never read into state. " +
			"`last_seen` and `since_last_seen` change on every read.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the managed client. Mutually exclusive with `name`.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The managed client name. Use as a lookup key (mutually exclusive with `id`) or read the resolved name after lookup by `id`.",
			},
			"cluster_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the cluster the managed client belongs to. When set together with `name`, only managed clients of that cluster match. Cannot be set together with `id`.",
			},
			"client_id":            computed("The client ID used in API management."),
			"description":          computed("The managed client description."),
			"type":                 computed("The managed client type, such as `VA` or `CCG`."),
			"cluster_type":         computed("The type of the cluster the managed client belongs to, such as `idn`."),
			"status":               computed("The managed client status, such as `NORMAL`, `WARNING`, `ERROR` or `NOT_CONFIGURED`."),
			"ip_address":           computed("The public IP address of the managed client."),
			"va_version":           computed("The version the virtual appliance is running."),
			"last_seen":            computed("When the managed client was last seen by the server."),
			"since_last_seen":      computed("Milliseconds since the managed client last polled the server."),
			"provision_status":     computed("The provisioning status of the managed client, such as `PROVISIONED` or `DRAFT`."),
			"api_gateway_base_url": computed("The API gateway base URL of the managed client."),
			"created_at":           computed("When the managed client was created."),
			"updated_at":           computed("When the managed client was last updated."),
		},
	}
}

func (d *managedClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config managedClientModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := isSet(config.ID)
	hasName := isSet(config.Name)
	hasCluster := isSet(config.ClusterID)

	switch {
	case !hasID && !hasName:
		resp.Diagnostics.AddError("Missing required argument", "One of `id` or `name` must be provided.")
		return
	case hasID && hasName:
		resp.Diagnostics.AddError("Conflicting arguments", "`id` and `name` are mutually exclusive. Provide only one.")
		return
	case hasID && hasCluster:
		resp.Diagnostics.AddError("Conflicting arguments", "`cluster_id` can only be used together with `name`.")
		return
	}

	var apiResp *client.ManagedClientAPI

	if hasID {
		id := config.ID.ValueString()
		tflog.Debug(ctx, "Reading managed client data source by ID", map[string]any{"id": id})
		mc, err := d.client.GetManagedClient(ctx, id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Managed Client",
				fmt.Sprintf("Could not read managed client %q: %s", id, err.Error()),
			)
			return
		}
		apiResp = mc
	} else {
		name := config.Name.ValueString()
		clusterID := config.ClusterID.ValueString()
		tflog.Debug(ctx, "Reading managed client data source by name", map[string]any{"name": name, "cluster_id": clusterID})

		all, err := d.client.ListManagedClients(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Managed Clients",
				fmt.Sprintf("Could not list managed clients: %s", err.Error()),
			)
			return
		}

		where := fmt.Sprintf("name %q", name)
		if clusterID != "" {
			where += fmt.Sprintf(" in cluster %q", clusterID)
		}

		matches := matchManagedClients(all, name, clusterID)
		switch len(matches) {
		case 0:
			resp.Diagnostics.AddError(
				"No managed client found",
				fmt.Sprintf("No managed client matched %s. Verify the name and cluster.", where),
			)
			return
		case 1:
			apiResp = &matches[0]
		default:
			resp.Diagnostics.AddError(
				"Multiple managed clients found",
				fmt.Sprintf("%d managed clients matched %s. Use `id`, or set `cluster_id` to narrow the match.", len(matches), where),
			)
			return
		}
	}

	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Managed Client", "Received nil response from SailPoint API")
		return
	}

	var state managedClientModel
	state.FromAPI(apiResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully read managed client data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func isSet(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown() && v.ValueString() != ""
}
