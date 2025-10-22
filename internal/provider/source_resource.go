package provider

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
				Description:         "The ID of the Source.",
				MarkdownDescription: "The ID of the Source.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the Source.",
				MarkdownDescription: "The name of the Source.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the source and its purpose.",
				MarkdownDescription: "A detailed description explaining the source and what system it represents.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.SingleNestedAttribute{
				Description:         "The owner of the Source.",
				MarkdownDescription: "The owner of the Source.",
				Required:            true,
				Attributes:          ObjectRefResourceSchema(),
			},
			"cluster": schema.SingleNestedAttribute{
				Description:         "The cluster associated with the source.",
				MarkdownDescription: "The cluster to which this source belongs.",
				Optional:            true,
				Attributes:          ObjectRefResourceSchema(),
			},
			"type": schema.StringAttribute{
				Description:         "The type of the source (e.g., 'Application', 'Database').",
				MarkdownDescription: "The category or type of the source within SailPoint ISC.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector": schema.StringAttribute{
				Description:         "The connector of the Source.",
				MarkdownDescription: "The connector of the Source.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_class": schema.StringAttribute{
				Description:         "The class of the connector used by the source.",
				MarkdownDescription: "The specific class name of the connector implementation for this source.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_threshold": schema.Int32Attribute{
				Description:         "The delete threshold for the source.",
				MarkdownDescription: "The threshold value that determines when accounts are deleted from the source.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Description:         "The creation timestamp of the Source.",
				MarkdownDescription: "The creation timestamp of the Source.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				Description:         "The last modified timestamp of the Source.",
				MarkdownDescription: "The last modified timestamp of the Source.",
				Computed:            true,
			},

			// "account_correlation_config": schema.SingleNestedAttribute{
			// 	Description:         "The account correlation configuration for the source.",
			// 	MarkdownDescription: "The account correlation configuration associated with this source.",
			// 	Optional:            true,
			// 	Attributes:          ObjectRefResourceSchema(),
			// },
			// "account_correlation_rule": schema.SingleNestedAttribute{
			// 	Description:         "The account correlation rule for the source.",
			// 	MarkdownDescription: "The account correlation rule associated with this source.",
			// 	Optional:            true,
			// 	Attributes:          ObjectRefResourceSchema(),
			// },
			// "manager_correlation_mapping": schema.SingleNestedAttribute{
			// 	Description:         "The manager correlation mapping for the source.",
			// 	MarkdownDescription: "The manager correlation mapping associated with this source.",
			// 	Optional:            true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"account_attribute_name": schema.StringAttribute{
			// 			Description:         "The account attribute name used for manager correlation.",
			// 			MarkdownDescription: "The name of the account attribute used in the manager correlation mapping.",
			// 			Optional:            true,
			// 		},
			// 		"identity_attribute_name": schema.StringAttribute{
			// 			Description:         "The identity attribute name used for manager correlation.",
			// 			MarkdownDescription: "The name of the identity attribute used in the manager correlation mapping.",
			// 			Optional:            true,
			// 		},
			// 	},
			// },
			// "manager_correlation_rule": schema.SingleNestedAttribute{
			// 	Description:         "The manager correlation rule for the source.",
			// 	MarkdownDescription: "The manager correlation rule associated with this source.",
			// 	Optional:            true,
			// 	Attributes:          ObjectRefResourceSchema(),
			// },
			// "before_provisioning_rule": schema.SingleNestedAttribute{
			// 	Description:         "The before provisioning rule for the source.",
			// 	MarkdownDescription: "The before provisioning rule associated with this source.",
			// 	Optional:            true,
			// 	Attributes:          ObjectRefResourceSchema(),
			// },
			// "schemas": schema.ListNestedAttribute{
			// 	Description:         "The schemas associated with the source.",
			// 	MarkdownDescription: "A list of schemas that define the structure of data for this source.",
			// 	Optional:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: ObjectRefResourceSchema(),
			// 	},
			// },
			// "password_policies": schema.ListNestedAttribute{
			// 	Description:         "The password policies associated with the source.",
			// 	MarkdownDescription: "A list of password policies that apply to this source.",
			// 	Optional:            true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: ObjectRefResourceSchema(),
			// 	},
			// },
			// "features": schema.ListAttribute{
			// 	Description:         "A list of features enabled for the source.",
			// 	MarkdownDescription: "An array of features that are enabled or supported by this source.",
			// 	Optional:            true,
			// 	ElementType:         types.StringType,
			// },

			// "connector_attributes": schema.StringAttribute{
			// 	Description:         "The attributes of the connector used by the source.",
			// 	MarkdownDescription: "A map of attributes and their values for the connector associated with this source.",
			// 	Optional:            true,
			// 	Computed:            true,
			// 	Sensitive:           true,
			// 	CustomType:          jsontypes.NormalizedType{},
			// },

			// "authoritative": schema.BoolAttribute{
			// 	Description:         "Indicates if the source is authoritative.",
			// 	MarkdownDescription: "A boolean flag indicating whether this source is considered authoritative for its data.",
			// 	Optional:            true,
			// },
			// "management_workgroup": schema.SingleNestedAttribute{
			// 	Description:         "The management workgroup for the source.",
			// 	MarkdownDescription: "The workgroup responsible for managing this source.",
			// 	Optional:            true,
			// 	Attributes:          ObjectRefResourceSchema(),
			// },
			// "healthy": schema.BoolAttribute{
			// 	Description:         "Indicates if the source is healthy.",
			// 	MarkdownDescription: "A boolean flag indicating the health status of the source.",
			// 	Optional:            true,
			// },
			// "status": schema.StringAttribute{
			// 	Description:         "The current status of the source.",
			// 	MarkdownDescription: "The operational status of the source within SailPoint ISC.",
			// 	Optional:            true,
			// },
			// "since": schema.StringAttribute{
			// 	Description:         "The timestamp since when the source has been active.",
			// 	MarkdownDescription: "The date and time indicating when the source became active in ISO 8601 format.",
			// 	Optional:            true,
			// },
			// "connector_id": schema.StringAttribute{
			// 	Description:         "The unique identifier of the connector used by the source.",
			// 	MarkdownDescription: "The UUID of the connector associated with this source.",
			// 	Optional:            true,
			// },
			// "connector_name": schema.StringAttribute{
			// 	Description:         "The name of the connector used by the source.",
			// 	MarkdownDescription: "The human-readable name of the connector associated with this source.",
			// 	Optional:            true,
			// },
			// "connector_type": schema.StringAttribute{
			// 	Description:         "The type of the connector used by the source.",
			// 	MarkdownDescription: "The category or type of connector used for this source.",
			// 	Optional:            true,
			// },
			// "connector_implementation_id": schema.StringAttribute{
			// 	Description:         "The implementation ID of the connector used by the source.",
			// 	MarkdownDescription: "The specific implementation identifier of the connector for this source.",
			// 	Optional:            true,
			// },
			// "credential_provider_enabled": schema.BoolAttribute{
			// 	Description:         "Indicates if the credential provider is enabled for the source.",
			// 	MarkdownDescription: "A boolean flag indicating whether the credential provider feature is enabled for this source.",
			// 	Optional:            true,
			// },
			// "category": schema.StringAttribute{
			// 	Description:         "The category of the source.",
			// 	MarkdownDescription: "The classification or category assigned to this source.",
			// 	Optional:            true,
			// },
		},
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.Source
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source := plan.ConvertToCreateRequestPtr(ctx)

	createdSource, err := r.client.CreateSource(ctx, source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			fmt.Sprintf("An error occurred while creating the source: %s", err.Error()),
		)
		return
	}

	plan.ConvertFromSailPointForResource(ctx, createdSource)

	if plan.Cluster != nil {
		plan.Cluster = models.NewObjectRefFromSailPoint(createdSource.Cluster)
	}

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

	state.ConvertFromSailPointForResource(ctx, fetchedSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// jsonPatches := []client.JSONPatchOperation{}
	var plan, state models.Source
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	patches := state.BuildPatchOptions(ctx, &plan)

	if len(patches) == 0 {
		return
	}

	updatedSource, err := r.client.PatchSource(ctx, state.ID.ValueString(), patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			fmt.Sprintf("An error occurred while updating the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Map Simple fields back to state
	// plan.Name = types.StringValue(updatedSource.Name)
	// plan.Description = types.StringValue(updatedSource.Description)
	// plan.Owner = models.NewObjectRefFromSailPoint(updatedSource.Owner)
	// plan.ConnectorClass = types.StringValue(updatedSource.ConnectorClass)
	// plan.DeleteThreshold = types.Int32Value(updatedSource.DeleteThreshold)
	// plan.Modified = types.StringValue(updatedSource.Modified)

	plan.ConvertFromSailPointForResource(ctx, updatedSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			fmt.Sprintf("An error occurred while deleting the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}
