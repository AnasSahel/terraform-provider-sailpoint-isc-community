package resources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/sailpoint_sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.Resource              = &formDefinitionResource{}
	_ resource.ResourceWithConfigure = &formDefinitionResource{}
)

type formDefinitionResource struct {
	client *sailpoint_sdk.Client
}

func NewFormDefinitionResource() resource.Resource {
	return &formDefinitionResource{}
}

func (r *formDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

func (r *formDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint_sdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sailpoint_sdk.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *formDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the form definition.",
				MarkdownDescription: "The ID of the form definition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the form definition.",
				MarkdownDescription: "The name of the form definition.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "The description of the form definition.",
				MarkdownDescription: "The description of the form definition.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation timestamp of the form definition.",
				MarkdownDescription: "The creation timestamp of the form definition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "The last modified timestamp of the form definition.",
				MarkdownDescription: "The last modified timestamp of the form definition.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.SingleNestedAttribute{
				Required:            true,
				Description:         "The owner of the form definition.",
				MarkdownDescription: "The owner of the form definition.",
				Attributes:          formOwnerSchema(),
			},
		},
	}
}

func (r *formDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.FormDefinitionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fd, err := r.client.FormDefinitionApi.CreateFormDefinition(ctx, plan.ToCreateApiModel())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Form Definition",
			fmt.Sprintf("Failed to create form definition: %v", err),
		)
		return
	}

	plan.FromApiModel(fd)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *formDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.FormDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fd, err := r.client.FormDefinitionApi.GetFormDefinitionById(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Form Definition",
			fmt.Sprintf("Failed to read form definition with ID %s: %v", state.Id.ValueString(), err),
		)
		return
	}

	state.FromApiModel(fd)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *formDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *formDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.FormDefinitionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.FormDefinitionApi.DeleteFormDefinition(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Form Definition",
			fmt.Sprintf("Failed to delete form definition with ID %s: %v", state.Id.ValueString(), err),
		)
		return
	}
}

func formOwnerSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Required:            true,
			Description:         "The type of the owner.",
			MarkdownDescription: "The type of the owner.",
		},
		"id": schema.StringAttribute{
			Required:            true,
			Description:         "The ID of the owner.",
			MarkdownDescription: "The ID of the owner.",
		},
		"name": schema.StringAttribute{
			Optional:            true,
			Description:         "The name of the owner.",
			MarkdownDescription: "The name of the owner.",
		},
	}
}
