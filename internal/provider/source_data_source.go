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
				Required:            true,
				Description:         "The unique identifier of the source.",
				MarkdownDescription: "The unique identifier (UUID) of the source to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "The display name of the source.",
				MarkdownDescription: "The human-readable name assigned to the source.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "A description of the source and its purpose.",
				MarkdownDescription: "A detailed description explaining the source and what system it represents.",
			},
			"owner": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The owner of the source.",
				MarkdownDescription: "The owner of the source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"cluster": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The cluster associated with the source.",
				MarkdownDescription: "The cluster to which this source belongs.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"account_correlation_config": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The account correlation configuration for the source.",
				MarkdownDescription: "The account correlation configuration associated with this source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"account_correlation_rule": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The account correlation rule for the source.",
				MarkdownDescription: "The account correlation rule associated with this source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"manager_correlation_rule": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The manager correlation rule for the source.",
				MarkdownDescription: "The manager correlation rule associated with this source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"before_provisioning_rule": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The before provisioning rule for the source.",
				MarkdownDescription: "The before provisioning rule associated with this source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"features": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "A list of features enabled for the source.",
				MarkdownDescription: "An array of features that are enabled or supported by this source.",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the source (e.g., 'Application', 'Database').",
				MarkdownDescription: "The category or type of the source within SailPoint ISC.",
			},
			"connector": schema.StringAttribute{
				Computed:            true,
				Description:         "The connector associated with the source.",
				MarkdownDescription: "The connector used to integrate with the source system.",
			},
			"connector_class": schema.StringAttribute{
				Computed:            true,
				Description:         "The class of the connector used by the source.",
				MarkdownDescription: "The specific class name of the connector implementation for this source.",
			},
			"delete_threshold": schema.Int32Attribute{
				Computed:            true,
				Description:         "The delete threshold for the source.",
				MarkdownDescription: "The threshold value that determines when accounts are deleted from the source.",
			},
			"authoritative": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if the source is authoritative.",
				MarkdownDescription: "A boolean flag indicating whether this source is considered authoritative for its data.",
			},
			"management_workgroup": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The management workgroup for the source.",
				MarkdownDescription: "The workgroup responsible for managing this source.",
				Attributes:          models.ObjectRefDataSourceSchema(),
			},
			"healthy": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if the source is healthy.",
				MarkdownDescription: "A boolean flag indicating the health status of the source.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The current status of the source.",
				MarkdownDescription: "The operational status of the source within SailPoint ISC.",
			},
			"since": schema.StringAttribute{
				Computed:            true,
				Description:         "The timestamp since when the source has been active.",
				MarkdownDescription: "The date and time indicating when the source became active in ISO 8601 format.",
			},
			"connector_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The unique identifier of the connector used by the source.",
				MarkdownDescription: "The UUID of the connector associated with this source.",
			},
			"connector_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the connector used by the source.",
				MarkdownDescription: "The human-readable name of the connector associated with this source.",
			},
			"connector_type": schema.StringAttribute{
				Computed:            true,
				Description:         "The type of the connector used by the source.",
				MarkdownDescription: "The category or type of connector used for this source.",
			},
			"connector_implementation_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The implementation ID of the connector used by the source.",
				MarkdownDescription: "The specific implementation identifier of the connector for this source.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "The timestamp when the source was created.",
				MarkdownDescription: "The date and time when the source was initially created in ISO 8601 format.",
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "The timestamp when the source was last modified.",
				MarkdownDescription: "The date and time when the source was last updated in ISO 8601 format.",
			},
			"credential_provider_enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if the credential provider is enabled for the source.",
				MarkdownDescription: "A boolean flag indicating whether the credential provider feature is enabled for this source.",
			},
			"category": schema.StringAttribute{
				Computed:            true,
				Description:         "The category of the source.",
				MarkdownDescription: "The classification or category assigned to this source.",
			},
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

	// Map simple string fields
	state.Name = types.StringValue(source.Name)
	state.Description = types.StringValue(source.Description)
	state.Type = types.StringValue(source.Type)
	state.Connector = types.StringValue(source.Connector)
	state.ConnectorClass = types.StringValue(source.ConnectorClass)
	state.DeleteThreshold = types.Int32Value(source.DeleteThreshold)
	state.Created = types.StringValue(source.Created)
	state.Modified = types.StringValue(source.Modified)
	state.Authoritative = types.BoolValue(source.Authoritative)
	state.Healthy = types.BoolValue(source.Healthy)
	state.Status = types.StringValue(source.Status)
	state.Since = types.StringValue(source.Since)
	state.ConnectorID = types.StringValue(source.ConnectorID)
	state.ConnectorName = types.StringValue(source.ConnectorName)
	state.ConnectorType = types.StringValue(source.ConnectorType)
	state.ConnectorImplementationID = types.StringValue(source.ConnectorImplementationID)
	state.CredentialProviderEnabled = types.BoolValue(source.CredentialProviderEnabled)
	state.Category = types.StringValue(source.Category)

	// Map simple list fields
	features, diags := types.ListValueFrom(ctx, types.StringType, source.Features)
	resp.Diagnostics.Append(diags...)
	state.Features = features

	// Map Object fields
	state.Owner = models.NewObjectRefFromAPI(source.Owner)
	state.Cluster = models.NewObjectRefFromAPI(source.Cluster)
	state.AccountCorrelationConfig = models.NewObjectRefFromAPI(source.AccountCorrelationConfig)
	state.AccountCorrelationRule = models.NewObjectRefFromAPI(source.AccountCorrelationRule)
	state.ManagerCorrelationRule = models.NewObjectRefFromAPI(source.ManagerCorrelationRule)
	state.BeforeProvisioningRule = models.NewObjectRefFromAPI(source.BeforeProvisioningRule)
	state.ManagementWorkgroup = models.NewObjectRefFromAPI(source.ManagementWorkgroup)

	tflog.Debug(ctx, "Successfully fetched source", map[string]any{"source_id": source.ID})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Helper functions
