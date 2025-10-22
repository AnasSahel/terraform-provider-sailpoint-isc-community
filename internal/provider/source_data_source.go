package provider

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &sourceDataSource{}
)

type sourceDataSource struct {
	client *client.Client
}

func NewSourceDataSource() datasource.DataSource {
	return new(sourceDataSource)
}

func (d *sourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Identity Security Cloud (ISC) source.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific source in SailPoint ISC. Sources represent systems or applications from which identity data is aggregated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the source.",
				MarkdownDescription: "The unique identifier (UUID) of the source to retrieve.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The display name of the source.",
				MarkdownDescription: "The human-readable name assigned to the source.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the source and its purpose.",
				MarkdownDescription: "A detailed description explaining the source and what system it represents.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				Description:         "The owner of the source.",
				MarkdownDescription: "The owner of the source.",
				Computed:            true,
				Attributes:          ObjectRefDataSourceSchema(),
			},
			"cluster": schema.SingleNestedAttribute{
				Description:         "The cluster associated with the source.",
				MarkdownDescription: "The cluster to which this source belongs.",
				Computed:            true,
				Attributes:          ObjectRefDataSourceSchema(),
			},
			"features": schema.SetAttribute{
				Description:         "A list of features enabled for the source.",
				MarkdownDescription: "An array of features that are enabled or supported by this source.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"type": schema.StringAttribute{
				Description:         "The type of the source (e.g., 'Application', 'Database').",
				MarkdownDescription: "The category or type of the source within SailPoint ISC.",
				Computed:            true,
			},
			"connector": schema.StringAttribute{
				Description:         "The connector associated with the source.",
				MarkdownDescription: "The connector used to integrate with the source system.",
				Computed:            true,
			},
			"connector_class": schema.StringAttribute{
				Description:         "The class of the connector used by the source.",
				MarkdownDescription: "The specific class name of the connector implementation for this source.",
				Computed:            true,
			},
			"delete_threshold": schema.Int32Attribute{
				Description:         "The delete threshold for the source.",
				MarkdownDescription: "The threshold value that determines when accounts are deleted from the source.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				Description:         "The timestamp when the source was created.",
				MarkdownDescription: "The date and time when the source was initially created in ISO 8601 format.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				Description:         "The timestamp when the source was last modified.",
				MarkdownDescription: "The date and time when the source was last updated in ISO 8601 format.",
				Computed:            true,
			},

			// "account_correlation_config": schema.SingleNestedAttribute{
			// 	Description:         "The account correlation configuration for the source.",
			// 	MarkdownDescription: "The account correlation configuration associated with this source.",
			// 	Computed:            true,
			// 	Attributes:          ObjectRefDataSourceSchema(),
			// },
			// "account_correlation_rule": schema.SingleNestedAttribute{
			// 	Description:         "The account correlation rule for the source.",
			// 	MarkdownDescription: "The account correlation rule associated with this source.",
			// 	Computed:            true,
			// 	Attributes:          ObjectRefDataSourceSchema(),
			// },
			// "manager_correlation_mapping": schema.SingleNestedAttribute{
			// 	Description:         "The manager correlation mapping for the source.",
			// 	MarkdownDescription: "The manager correlation mapping associated with this source.",
			// 	Computed:            true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"account_attribute_name": schema.StringAttribute{
			// 			Description:         "The account attribute name used for manager correlation.",
			// 			MarkdownDescription: "The name of the account attribute used in the manager correlation mapping.",
			// 			Computed:            true,
			// 		},
			// 		"identity_attribute_name": schema.StringAttribute{
			// 			Description:         "The identity attribute name used for manager correlation.",
			// 			MarkdownDescription: "The name of the identity attribute used in the manager correlation mapping.",
			// 			Computed:            true,
			// 		},
			// 	},
			// },
			// "manager_correlation_rule": schema.SingleNestedAttribute{
			// 	Description:         "The manager correlation rule for the source.",
			// 	MarkdownDescription: "The manager correlation rule associated with this source.",
			// 	Computed:            true,
			// 	Attributes:          ObjectRefDataSourceSchema(),
			// },
			// "before_provisioning_rule": schema.SingleNestedAttribute{
			// 	Description:         "The before provisioning rule for the source.",
			// 	MarkdownDescription: "The before provisioning rule associated with this source.",
			// 	Computed:            true,
			// 	Attributes:          ObjectRefDataSourceSchema(),
			// },
			// "schemas": schema.ListNestedAttribute{
			// 	Description:         "The schemas associated with the source.",
			// 	MarkdownDescription: "A list of schemas that define the structure of data for this source.",
			// 	Computed:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: ObjectRefDataSourceSchema(),
			// 	},
			// },
			// "password_policies": schema.ListNestedAttribute{
			// 	Description:         "The password policies associated with the source.",
			// 	MarkdownDescription: "A list of password policies that apply to this source.",
			// 	Computed:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: ObjectRefDataSourceSchema(),
			// 	},
			// },

			// "connector_attributes": schema.StringAttribute{
			// 	Description:         "The attributes of the connector used by the source.",
			// 	MarkdownDescription: "A map of attributes and their values for the connector associated with this source.",
			// 	Computed:            true,
			// 	Sensitive:           true,
			// 	CustomType:          jsontypes.NormalizedType{},
			// },

			// "authoritative": schema.BoolAttribute{
			// 	Description:         "Indicates if the source is authoritative.",
			// 	MarkdownDescription: "A boolean flag indicating whether this source is considered authoritative for its data.",
			// 	Computed:            true,
			// },
			// "management_workgroup": schema.SingleNestedAttribute{
			// 	Description:         "The management workgroup for the source.",
			// 	MarkdownDescription: "The workgroup responsible for managing this source.",
			// 	Computed:            true,
			// 	Attributes:          ObjectRefDataSourceSchema(),
			// },
			// "healthy": schema.BoolAttribute{
			// 	Description:         "Indicates if the source is healthy.",
			// 	MarkdownDescription: "A boolean flag indicating the health status of the source.",
			// 	Computed:            true,
			// },
			// "status": schema.StringAttribute{
			// 	Description:         "The current status of the source.",
			// 	MarkdownDescription: "The operational status of the source within SailPoint ISC.",
			// 	Computed:            true,
			// },
			// "since": schema.StringAttribute{
			// 	Description:         "The timestamp since when the source has been active.",
			// 	MarkdownDescription: "The date and time indicating when the source became active in ISO 8601 format.",
			// 	Computed:            true,
			// },
			// "connector_id": schema.StringAttribute{
			// 	Description:         "The unique identifier of the connector used by the source.",
			// 	MarkdownDescription: "The UUID of the connector associated with this source.",
			// 	Computed:            true,
			// },
			// "connector_name": schema.StringAttribute{
			// 	Description:         "The name of the connector used by the source.",
			// 	MarkdownDescription: "The human-readable name of the connector associated with this source.",
			// 	Computed:            true,
			// },
			// "connector_type": schema.StringAttribute{
			// 	Description:         "The type of the connector used by the source.",
			// 	MarkdownDescription: "The category or type of connector used for this source.",
			// 	Computed:            true,
			// },
			// "connector_implementation_id": schema.StringAttribute{
			// 	Description:         "The implementation ID of the connector used by the source.",
			// 	MarkdownDescription: "The specific implementation identifier of the connector for this source.",
			// 	Computed:            true,
			// },
			// "credential_provider_enabled": schema.BoolAttribute{
			// 	Description:         "Indicates if the credential provider is enabled for the source.",
			// 	MarkdownDescription: "A boolean flag indicating whether the credential provider feature is enabled for this source.",
			// 	Computed:            true,
			// },
			// "category": schema.StringAttribute{
			// 	Description:         "The category of the source.",
			// 	MarkdownDescription: "The classification or category assigned to this source.",
			// 	Computed:            true,
			// },
		},
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching source by ID", map[string]any{"id": state.ID.ValueString()})

	source, err := d.client.GetSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Fetch Source",
			fmt.Sprintf("Failed to fetch source with ID %q: %v", state.ID.ValueString(), err),
		)
		return
	}

	state.ConvertFromSailPointForDataSource(ctx, source)
	tflog.Debug(ctx, "Successfully fetched source", map[string]any{"source_id": source.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Helper functions
