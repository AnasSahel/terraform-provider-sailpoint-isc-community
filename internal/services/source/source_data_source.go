// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceDataSource{}
)

type sourceDataSource struct {
	client *client.Client
}

// sourceDataSourceModel is the data source's view of a source. It mirrors
// sourceModel minus the resource-only ignore_* arguments, which the data
// source schema does not declare.
type sourceDataSourceModel struct {
	ID                        types.String           `tfsdk:"id"`
	Name                      types.String           `tfsdk:"name"`
	Description               types.String           `tfsdk:"description"`
	Owner                     *common.ObjectRefModel `tfsdk:"owner"`
	Cluster                   *common.ObjectRefModel `tfsdk:"cluster"`
	Connector                 types.String           `tfsdk:"connector"`
	ConnectorClass            types.String           `tfsdk:"connector_class"`
	ConnectorAttributes       jsontypes.Normalized   `tfsdk:"connector_attributes"`
	ConnectorAttributesAll    jsontypes.Normalized   `tfsdk:"connector_attributes_all"`
	ConnectionType            types.String           `tfsdk:"connection_type"`
	Type                      types.String           `tfsdk:"type"`
	DeleteThreshold           types.Int64            `tfsdk:"delete_threshold"`
	Authoritative             types.Bool             `tfsdk:"authoritative"`
	Healthy                   types.Bool             `tfsdk:"healthy"`
	Status                    types.String           `tfsdk:"status"`
	Features                  types.Set              `tfsdk:"features"`
	CredentialProviderEnabled types.Bool             `tfsdk:"credential_provider_enabled"`
	Category                  types.String           `tfsdk:"category"`
	ProvisionAsCsv            types.Bool             `tfsdk:"provision_as_csv"`
	Created                   types.String           `tfsdk:"created"`
	Modified                  types.String           `tfsdk:"modified"`
}

func newSourceDataSourceModel(m sourceModel) sourceDataSourceModel {
	return sourceDataSourceModel{
		ID:                        m.ID,
		Name:                      m.Name,
		Description:               m.Description,
		Owner:                     m.Owner,
		Cluster:                   m.Cluster,
		Connector:                 m.Connector,
		ConnectorClass:            m.ConnectorClass,
		ConnectorAttributes:       m.ConnectorAttributes,
		ConnectorAttributesAll:    m.ConnectorAttributesAll,
		ConnectionType:            m.ConnectionType,
		Type:                      m.Type,
		DeleteThreshold:           m.DeleteThreshold,
		Authoritative:             m.Authoritative,
		Healthy:                   m.Healthy,
		Status:                    m.Status,
		Features:                  m.Features,
		CredentialProviderEnabled: m.CredentialProviderEnabled,
		Category:                  m.Category,
		ProvisionAsCsv:            m.ProvisionAsCsv,
		Created:                   m.Created,
		Modified:                  m.Modified,
	}
}

func NewSourceDataSource() datasource.DataSource {
	return &sourceDataSource{}
}

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Source.",
		MarkdownDescription: "Data source for SailPoint Source. Sources represent managed systems (e.g., Active Directory, Workday) in Identity Security Cloud.\n\n" +
			"Look up a source either by `id` or by `name`; exactly one of the two must be set. A name lookup must match exactly one " +
			"source, which lets a configuration reference a source by its name instead of a tenant-specific ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the source. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the source. When set instead of `id`, the source is looked up by exact name. Exactly one of `id` or `name` must be set.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the source.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
					},
				},
			},
			"cluster": schema.SingleNestedAttribute{
				MarkdownDescription: "The cluster associated with this source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the cluster.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the cluster.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the cluster.",
						Computed:            true,
					},
				},
			},
			"connector": schema.StringAttribute{
				MarkdownDescription: "The connector script name.",
				Computed:            true,
			},
			"connector_class": schema.StringAttribute{
				MarkdownDescription: "The fully qualified name of the Java class that implements the connector interface.",
				Computed:            true,
			},
			"connector_attributes": schema.StringAttribute{
				MarkdownDescription: "A JSON object containing connector-specific configuration.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"connector_attributes_all": schema.StringAttribute{
				MarkdownDescription: "The full connector attributes as returned by the API, including both user-configured and server-managed keys.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"connection_type": schema.StringAttribute{
				MarkdownDescription: "The connection type.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of system being managed.",
				Computed:            true,
			},
			"delete_threshold": schema.Int64Attribute{
				MarkdownDescription: "The percentage threshold for skipping the delete phase (0-100).",
				Computed:            true,
			},
			"authoritative": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is referenced by an identity profile.",
				Computed:            true,
			},
			"healthy": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is healthy.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the source.",
				Computed:            true,
			},
			"features": schema.SetAttribute{
				MarkdownDescription: "The list of features enabled for the source.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"credential_provider_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether credential provider is enabled for the source.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The source category.",
				Computed:            true,
			},
			"provision_as_csv": schema.BoolAttribute{
				MarkdownDescription: "Whether the source was provisioned as a CSV source.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SailPoint Source data source")

	var config sourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := isSet(config.ID)
	hasName := isSet(config.Name)
	if !hasID && !hasName {
		resp.Diagnostics.AddError("Missing required argument", "Provide either `id` or `name` to look up a source.")
		return
	}
	if hasID && hasName {
		resp.Diagnostics.AddError("Conflicting arguments", "`id` and `name` cannot both be set. Provide only one of them.")
		return
	}

	var sourceResponse *client.SourceAPI
	if hasID {
		id := config.ID.ValueString()
		tflog.Debug(ctx, "Fetching source from SailPoint by ID", map[string]any{"id": id})

		src, err := d.client.GetSource(ctx, id)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Source",
				fmt.Sprintf("Could not read SailPoint Source %q: %s", id, err.Error()),
			)
			return
		}
		sourceResponse = src
	} else {
		name := config.Name.ValueString()
		filters := "name eq " + common.QuoteFilterValue(name)
		tflog.Debug(ctx, "Fetching source from SailPoint by name", map[string]any{"filters": filters})

		// Two results are enough to tell "exactly one" from "ambiguous".
		sources, err := d.client.ListSources(ctx, filters, 2)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Sources",
				fmt.Sprintf("Could not list sources named %q: %s", name, err.Error()),
			)
			return
		}
		switch len(sources) {
		case 0:
			resp.Diagnostics.AddError("No source found", fmt.Sprintf("No source is named %q.", name))
			return
		case 1:
			sourceResponse = &sources[0]
		default:
			resp.Diagnostics.AddError(
				"Multiple sources found",
				fmt.Sprintf("More than one source is named %q. Supply `id` instead.", name),
			)
			return
		}
	}

	if sourceResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source",
			"Received nil response from SailPoint API",
		)
		return
	}

	var model sourceModel
	resp.Diagnostics.Append(model.FromAPI(ctx, *sourceResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state := newSourceDataSourceModel(model)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Source data source", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func isSet(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown() && v.ValueString() != ""
}
