package managedcluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/iancoleman/strcase"
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
	tflog.Debug(ctx, "Configuring ManagedClusterResource")

	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	tflog.Debug(ctx, "Configured ManagedClusterResource")
	r.client = client
}

func (r *ManagedClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_cluster"
}

func (r *ManagedClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Trace(ctx, "Preparing ManagedClusterResource schema")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the managed cluster.",
				MarkdownDescription: "The unique identifier for the managed cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the managed cluster.",
				MarkdownDescription: "The name of the managed cluster.",
			},
			"pod": schema.StringAttribute{
				Computed:            true,
				Description:         "The pod of the managed cluster.",
				MarkdownDescription: "The pod of the managed cluster.",
			},
			"org": schema.StringAttribute{
				Computed:            true,
				Description:         "The organization of the managed cluster.",
				MarkdownDescription: "The organization of the managed cluster.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				Description:         "The type of the managed cluster.",
				MarkdownDescription: "The type of the managed cluster.",
			},
			"configuration": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
				Description:         "The configuration of the managed cluster as key-value pairs.",
				MarkdownDescription: "The configuration of the managed cluster as key-value pairs.",
			},
			"description": schema.StringAttribute{
				Required:            true,
				Description:         "The description of the managed cluster.",
				MarkdownDescription: "The description of the managed cluster.",
			},
			"client_type": schema.StringAttribute{
				Computed:            true,
				Description:         "The client type of the managed cluster.",
				MarkdownDescription: "The client type of the managed cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ccg_version": schema.StringAttribute{
				Computed:            true,
				Description:         "The CCG version of the managed cluster.",
				MarkdownDescription: "The CCG version of the managed cluster.",
			},
			"pinned_config": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if the configuration is pinned.",
				MarkdownDescription: "Indicates if the configuration is pinned.",
			},
			"operational": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if the managed cluster is operational.",
				MarkdownDescription: "Indicates if the managed cluster is operational.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "The status of the managed cluster.",
				MarkdownDescription: "The status of the managed cluster.",
			},
			"public_key_certificate": schema.StringAttribute{
				Computed:            true,
				Description:         "The public key certificate of the managed cluster.",
				MarkdownDescription: "The public key certificate of the managed cluster.",
			},
			"public_key_thumbprint": schema.StringAttribute{
				Computed:            true,
				Description:         "The public key thumbprint of the managed cluster.",
				MarkdownDescription: "The public key thumbprint of the managed cluster.",
			},
			"public_key": schema.StringAttribute{
				Computed:            true,
				Description:         "The public key of the managed cluster.",
				MarkdownDescription: "The public key of the managed cluster.",
			},
			"alert_key": schema.StringAttribute{
				Computed:            true,
				Description:         "The alert key of the managed cluster.",
				MarkdownDescription: "The alert key of the managed cluster.",
			},
			"client_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "The client IDs associated with the managed cluster.",
				MarkdownDescription: "The client IDs associated with the managed cluster.",
			},
			"service_count": schema.Int32Attribute{
				Computed:            true,
				Description:         "The number of services associated with the managed cluster.",
				MarkdownDescription: "The number of services associated with the managed cluster.",
			},
			"cc_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The connected cloud ID associated with the managed cluster.",
				MarkdownDescription: "The connected cloud ID associated with the managed cluster.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation timestamp of the managed cluster.",
				MarkdownDescription: "The creation timestamp of the managed cluster.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				Description:         "The last update timestamp of the managed cluster.",
				MarkdownDescription: "The last update timestamp of the managed cluster.",
			},
		},
	}
}

func (r *ManagedClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Managed Cluster")
	var plan ManagedClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	managedClusterRequest, diags := plan.ToSailPointCreateManagedClusterRequest(ctx)
	tflog.Debug(ctx, "Create Managed Cluster request prepared", map[string]interface{}{
		"managedClusterRequest": managedClusterRequest,
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.CreateManagedCluster(context.Background()).ManagedClusterRequest(*managedClusterRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Managed Cluster",
			fmt.Sprintf("Could not create managed cluster, unexpected error: %s\n\nFull HTTP response: %v", err.Error(), httpResponse),
		)
		return
	}

	// Create a new state model to populate from API response
	var state ManagedClusterResourceModel
	resp.Diagnostics.Append(state.FromSailPointManagedCluster(ctx, managedCluster)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the original planned configuration to avoid inconsistency errors
	// Only keep user-specified configuration values, let API-added values be computed
	if !plan.Configuration.IsNull() {
		state.Configuration = plan.Configuration
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ManagedClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManagedClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting state for Managed Cluster")
		return
	}

	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.GetManagedCluster(context.Background(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Managed Cluster",
			fmt.Sprintf("Could not read managed cluster ID %s: %s\n\nFull HTTP response: %v", state.Id.ValueString(), err.Error(), httpResponse),
		)
		return
	}

	resp.Diagnostics.Append(state.FromSailPointManagedCluster(ctx, managedCluster)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ManagedClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Managed Cluster")

	// Get the current state
	var state ManagedClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the planned changes
	var plan ManagedClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create JSON Patch operations for the changes
	var patchOps []api_v2025.JsonPatchOperation

	// Check if name changed
	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		value := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&name)
		patchOp := api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/name",
			Value: &value,
		}
		patchOps = append(patchOps, patchOp)
	}

	// Check if description changed
	if !plan.Description.Equal(state.Description) {
		patchOp := api_v2025.JsonPatchOperation{
			Op:   "replace",
			Path: "/description",
		}
		if plan.Description.IsNull() {
			// For null values, we don't set the Value field
		} else {
			desc := plan.Description.ValueString()
			value := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&desc)
			patchOp.Value = &value
		}
		patchOps = append(patchOps, patchOp)
	}

	// Check if type changed
	if !plan.Type.Equal(state.Type) {
		typeVal := plan.Type.ValueString()
		value := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&typeVal)
		patchOp := api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/type",
			Value: &value,
		}
		patchOps = append(patchOps, patchOp)
	}

	// Check if configuration changed
	if !plan.Configuration.Equal(state.Configuration) {
		// Convert the configuration map to the format expected by the API
		configMap := make(map[string]interface{})
		for k, v := range plan.Configuration.Elements() {
			// Remove quotes from the string value and convert key to camelCase
			stringVal := strings.Trim(v.String(), `"`)
			configMap[strcase.ToLowerCamel(k)] = stringVal
		}

		value := api_v2025.MapmapOfStringAnyAsUpdateMultiHostSourcesRequestInnerValue(&configMap)
		patchOp := api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/configuration",
			Value: &value,
		}
		patchOps = append(patchOps, patchOp)
	}

	// If no changes, return early
	if len(patchOps) == 0 {
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	// Call the SailPoint API to update the managed cluster
	managedCluster, httpResponse, err := r.client.ManagedClustersAPI.UpdateManagedCluster(
		context.Background(),
		state.Id.ValueString(),
	).JsonPatchOperation(patchOps).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating managed cluster",
			fmt.Sprintf("Could not update managed cluster %s: %s\nHTTP Response: %v", state.Id.ValueString(), err, httpResponse),
		)
		return
	}

	// Start with the current state to preserve computed fields
	newState := state

	// Always update the ID (required field)
	newState.Id = types.StringValue(managedCluster.GetId())

	// Update fields that were changed and returned by the API
	if managedCluster.HasName() {
		newState.Name = types.StringValue(managedCluster.GetName())
	}

	if managedCluster.HasDescription() {
		newState.Description = types.StringValue(managedCluster.GetDescription())
	}

	if managedCluster.HasType() {
		newState.Type = types.StringValue(string(managedCluster.GetType()))
	}

	// For configuration, preserve the planned configuration to avoid inconsistency errors
	if !plan.Configuration.IsNull() {
		newState.Configuration = plan.Configuration
	}

	// Update computed fields only if they have meaningful values in the response
	if managedCluster.HasPod() && managedCluster.GetPod() != "" {
		newState.Pod = types.StringValue(managedCluster.GetPod())
	}

	if managedCluster.HasOrg() && managedCluster.GetOrg() != "" {
		newState.Org = types.StringValue(managedCluster.GetOrg())
	}

	// ClientType is always present but may be nullable - check if it's valid
	clientType, ok := managedCluster.GetClientTypeOk()
	if ok && clientType != nil {
		newState.ClientType = types.StringValue(string(*clientType))
	}

	// CcgVersion is required but check if it has meaningful value
	if managedCluster.GetCcgVersion() != "" && managedCluster.GetCcgVersion() != "Undefined" {
		newState.CcgVersion = types.StringValue(managedCluster.GetCcgVersion())
	}

	if managedCluster.HasUpdatedAt() {
		newState.UpdatedAt = types.StringValue(managedCluster.GetUpdatedAt().String())
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *ManagedClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ManagedClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.client.ManagedClustersAPI.DeleteManagedCluster(context.Background(), state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Managed Cluster",
			fmt.Sprintf("Could not delete managed cluster ID %s: %s\n\nFull HTTP response: %v", state.Id.ValueString(), err.Error(), httpResponse),
		)
		return
	}
}
