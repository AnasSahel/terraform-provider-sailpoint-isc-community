package provider

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &sourceResource{}
	_ resource.ResourceWithConfigure = &sourceResource{}
)

type sourceResource struct {
	client *client.Client
}

func NewSourceResource() resource.Resource {
	return new(sourceResource)
}

func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Source.",
		MarkdownDescription: "Resource representing a SailPoint Source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the Source.",
				MarkdownDescription: "The ID of the Source.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the Source.",
				MarkdownDescription: "The name of the Source.",
			},
			"owner": schema.SingleNestedAttribute{
				Required:            true,
				Description:         "The owner of the Source.",
				MarkdownDescription: "The owner of the Source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"connector": schema.StringAttribute{
				Required:            true,
				Description:         "The connector of the Source.",
				MarkdownDescription: "The connector of the Source.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation timestamp of the Source.",
				MarkdownDescription: "The creation timestamp of the Source.",
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "The last modified timestamp of the Source.",
				MarkdownDescription: "The last modified timestamp of the Source.",
			},
			// Additional attributes can be added here following the same pattern
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "A description of the source and its purpose.",
				MarkdownDescription: "A detailed description explaining the source and what system it represents.",
			},
			"cluster": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The cluster associated with the source.",
				MarkdownDescription: "The cluster to which this source belongs.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"account_correlation_config": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The account correlation configuration for the source.",
				MarkdownDescription: "The account correlation configuration associated with this source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"account_correlation_rule": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The account correlation rule for the source.",
				MarkdownDescription: "The account correlation rule associated with this source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"manager_correlation_mapping": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The manager correlation mapping for the source.",
				MarkdownDescription: "The manager correlation mapping associated with this source.",
				Attributes: map[string]schema.Attribute{
					"account_attribute_name": schema.StringAttribute{
						Optional:            true,
						Description:         "The account attribute name used for manager correlation.",
						MarkdownDescription: "The name of the account attribute used in the manager correlation mapping.",
					},
					"identity_attribute_name": schema.StringAttribute{
						Optional:            true,
						Description:         "The identity attribute name used for manager correlation.",
						MarkdownDescription: "The name of the identity attribute used in the manager correlation mapping.",
					},
				},
			},
			"manager_correlation_rule": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The manager correlation rule for the source.",
				MarkdownDescription: "The manager correlation rule associated with this source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"before_provisioning_rule": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The before provisioning rule for the source.",
				MarkdownDescription: "The before provisioning rule associated with this source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"schemas": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "The schemas associated with the source.",
				MarkdownDescription: "A list of schemas that define the structure of data for this source.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: ObjectRefResourceSchema(),
				},
			},
			"password_policies": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "The password policies associated with the source.",
				MarkdownDescription: "A list of password policies that apply to this source.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: ObjectRefResourceSchema(),
				},
			},
			"features": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "A list of features enabled for the source.",
				MarkdownDescription: "An array of features that are enabled or supported by this source.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Description:         "The type of the source (e.g., 'Application', 'Database').",
				MarkdownDescription: "The category or type of the source within SailPoint ISC.",
			},
			"connector_class": schema.StringAttribute{
				Optional:            true,
				Description:         "The class of the connector used by the source.",
				MarkdownDescription: "The specific class name of the connector implementation for this source.",
			},
			"connector_attributes": schema.DynamicAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The attributes of the connector used by the source.",
				MarkdownDescription: "A map of attributes and their values for the connector associated with this source.",
			},
			"delete_threshold": schema.Int32Attribute{
				Optional:            true,
				Description:         "The delete threshold for the source.",
				MarkdownDescription: "The threshold value that determines when accounts are deleted from the source.",
			},
			"authoritative": schema.BoolAttribute{
				Optional:            true,
				Description:         "Indicates if the source is authoritative.",
				MarkdownDescription: "A boolean flag indicating whether this source is considered authoritative for its data.",
			},
			"management_workgroup": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "The management workgroup for the source.",
				MarkdownDescription: "The workgroup responsible for managing this source.",
				Attributes:          ObjectRefResourceSchema(),
			},
			"healthy": schema.BoolAttribute{
				Optional:            true,
				Description:         "Indicates if the source is healthy.",
				MarkdownDescription: "A boolean flag indicating the health status of the source.",
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Description:         "The current status of the source.",
				MarkdownDescription: "The operational status of the source within SailPoint ISC.",
			},
			"since": schema.StringAttribute{
				Optional:            true,
				Description:         "The timestamp since when the source has been active.",
				MarkdownDescription: "The date and time indicating when the source became active in ISO 8601 format.",
			},
			"connector_id": schema.StringAttribute{
				Optional:            true,
				Description:         "The unique identifier of the connector used by the source.",
				MarkdownDescription: "The UUID of the connector associated with this source.",
			},
			"connector_name": schema.StringAttribute{
				Optional:            true,
				Description:         "The name of the connector used by the source.",
				MarkdownDescription: "The human-readable name of the connector associated with this source.",
			},
			"connector_type": schema.StringAttribute{
				Optional:            true,
				Description:         "The type of the connector used by the source.",
				MarkdownDescription: "The category or type of connector used for this source.",
			},
			"connector_implementation_id": schema.StringAttribute{
				Optional:            true,
				Description:         "The implementation ID of the connector used by the source.",
				MarkdownDescription: "The specific implementation identifier of the connector for this source.",
			},
			"credential_provider_enabled": schema.BoolAttribute{
				Optional:            true,
				Description:         "Indicates if the credential provider is enabled for the source.",
				MarkdownDescription: "A boolean flag indicating whether the credential provider feature is enabled for this source.",
			},
			"category": schema.StringAttribute{
				Optional:            true,
				Description:         "The category of the source.",
				MarkdownDescription: "The classification or category assigned to this source.",
			},
		},
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Source
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source := &client.Source{
		Name:      plan.Name.ValueString(),
		Connector: plan.Connector.ValueString(),
		Owner:     &client.ObjectRef{ID: plan.Owner.ID.ValueString(), Type: plan.Owner.Type.ValueString(), Name: plan.Owner.Name.ValueString()},
	}

	createdSource, err := r.client.CreateSource(ctx, source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			fmt.Sprintf("An error occurred while creating the source: %s", err.Error()),
		)
		return
	}

	// Map Simple fields back to state
	plan.ID = types.StringValue(createdSource.ID)
	plan.Created = types.StringValue(createdSource.Created)
	plan.Modified = types.StringValue(createdSource.Modified)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fetchedSource, err := r.client.GetSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			fmt.Sprintf("An error occurred while reading the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map Simple fields back to state
	state.Name = types.StringValue(fetchedSource.Name)
	state.Connector = types.StringValue(fetchedSource.Connector)
	state.Created = types.StringValue(fetchedSource.Created)
	state.Modified = types.StringValue(fetchedSource.Modified)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
