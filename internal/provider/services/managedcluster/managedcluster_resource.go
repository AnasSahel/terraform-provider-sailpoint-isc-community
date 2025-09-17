package managedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"type": schema.StringAttribute{
				Required:            true,
				Description:         "The type of the managed cluster.",
				MarkdownDescription: "The type of the managed cluster.",
			},
			"configuration": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "The configuration of the managed cluster.",
				MarkdownDescription: "The configuration of the managed cluster.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				Description:         "The description of the managed cluster.",
				MarkdownDescription: "The description of the managed cluster.",
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

	resp.Diagnostics.Append(plan.FromSailPointManagedCluster(ctx, managedCluster)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ManagedClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
