package identity_attribute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

var (
	_ resource.Resource              = &identityAttributeResource{}
	_ resource.ResourceWithConfigure = &identityAttributeResource{}
)

type identityAttributeResource struct {
	client *sailpoint.APIClient
}

func NewIdentityAttributeResource() resource.Resource {
	return &identityAttributeResource{}
}

func (r *identityAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)

	if !ok {
		resp.Diagnostics.AddError(ErrProviderDataTitle, fmt.Sprintf(ErrProviderDataMsg, req.ProviderData))
		return
	}

	r.client = client
}

func (r *identityAttributeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + ResourceTypeName
}

func (r *identityAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema()
}

func (r *identityAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IdentityAttributeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityAttribute, httpResponse, err := r.client.V2025.IdentityAttributesAPI.
		CreateIdentityAttribute(ctx).
		IdentityAttribute(MapTerraformResourceToIdentityAttribute(plan)).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			ErrResourceCreateTitle,
			fmt.Sprintf(ErrResourceCreateMsg, plan.Name.ValueString(), err, httpResponse),
		)
		return
	}

	state := MapIdentityAttributeToResourceModel(*identityAttribute)
	if plan.Sources != nil {
		state.Sources = make([]Source1, len(plan.Sources))
		for i, source := range identityAttribute.Sources {
			// sourceJson, _ := source.MarshalJSON()
			state.Sources[i] = Source1{
				Type:       types.StringValue(source.GetType()),
				Properties: types.StringValue("string(sourceJson)"),
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IdentityAttributeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityAttribute, httpResponse, err := r.client.V2025.IdentityAttributesAPI.
		GetIdentityAttribute(ctx, state.Name.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			ErrResourceReadTitle,
			fmt.Sprintf(ErrResourceReadMsg, state.Name.ValueString(), err, httpResponse),
		)
		return
	}

	updatedState := MapIdentityAttributeToResourceModel(*identityAttribute)
	if state.Sources != nil {
		updatedState.Sources = make([]Source1, len(state.Sources))
		for i, source := range identityAttribute.Sources {
			// sourceJson, _ := source.MarshalJSON()
			updatedState.Sources[i] = Source1{
				Type:       types.StringValue(source.GetType()),
				Properties: types.StringValue("string(sourceJson)"),
			}
		}
	}

	// updatedState := MapIdentityAttributeToResourceModel(*identityAttribute)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *identityAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *identityAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IdentityAttributeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.client.V2025.IdentityAttributesAPI.
		DeleteIdentityAttribute(ctx, state.Name.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			ErrResourceReadTitle,
			fmt.Sprintf(ErrResourceReadMsg, state.Name.ValueString(), err, httpResponse),
		)
		return
	}
}
