// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

type LifecycleStateResource struct {
	client       *sailpoint.APIClient
	validator    *Validator
	errorHandler *ErrorHandler
}

var (
	_ resource.Resource              = &LifecycleStateResource{}
	_ resource.ResourceWithConfigure = &LifecycleStateResource{}
)

func NewLifecycleStateResource() resource.Resource {
	return &LifecycleStateResource{
		validator:    NewValidator(),
		errorHandler: NewErrorHandler(),
	}
}

func (r *LifecycleStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)
	if !ok {
		resp.Diagnostics.Append(r.errorHandler.HandleConfigurationError(
			ErrUnexpectedConfigureType,
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this to the provider developers.", req.ProviderData),
		)...)
		return
	}

	r.client = client
}

func (r *LifecycleStateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + ResourceTypeName
}

func (r *LifecycleStateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = LifecycleStateResourceSchema()
}

func (r *LifecycleStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LifecycleStateResourceModel

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	validationDiags := r.validator.ValidateResourceModel(&plan)
	resp.Diagnostics.Append(validationDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate identity profile ID
	identityProfileID, validationDiags := r.validator.ValidateIdentityProfileID(plan.IdentityProfileId)
	resp.Diagnostics.Append(validationDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan to API object
	lifecycleStateRequest := FromCreatePlanToSailPointLifecycleState(ctx, &plan)

	// Call the API
	lifecycleState, httpResponse, err := r.client.V2025.LifecycleStatesAPI.
		CreateLifecycleState(ctx, identityProfileID).
		LifecycleState(lifecycleStateRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.Append(r.errorHandler.HandleAPIError(
			"Creating",
			err,
			httpResponse,
			fmt.Sprintf("identity profile ID: %s", identityProfileID),
		)...)
		return
	}

	// Convert API response to Terraform state
	newState := ToTerraformResource(ctx, &plan, lifecycleState)

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *LifecycleStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LifecycleStateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lifecycleState, httpResponse, err := r.client.V2025.LifecycleStatesAPI.
		GetLifecycleState(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lifecycle State",
			fmt.Sprintf("Could not read lifecycle state (profile ID: '%s', state ID: '%s'): %s\n\nHTTP Response: %v",
				state.IdentityProfileId.ValueString(),
				state.Id.ValueString(),
				err.Error(),
				httpResponse,
			),
		)
		return
	}

	newState := ToTerraformResource(ctx, &state, lifecycleState)

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *LifecycleStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state LifecycleStateResourceModel
	var plan LifecycleStateResourceModel

	var patches []api_v2025.JsonPatchOperation

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate patch operations for changed fields
	if patch := utils.CreateBoolPatch(plan.Enabled, state.Enabled, "/enabled"); patch != nil {
		patches = append(patches, *patch)
	}
	if patch := utils.CreateStringPatch(plan.Description, state.Description, "/description"); patch != nil {
		patches = append(patches, *patch)
	}
	if patch := utils.CreateInt32Patch(plan.Priority, state.Priority, "/priority"); patch != nil {
		patches = append(patches, *patch)
	}

	lifecycleState, httpResponse, err := r.client.V2025.LifecycleStatesAPI.
		UpdateLifecycleStates(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).
		JsonPatchOperation(patches).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Lifecycle State",
			fmt.Sprintf("Could not update lifecycle state: %s\n\n%v", err.Error(), httpResponse),
		)
		return
	}

	state.Enabled = types.BoolValue(lifecycleState.GetEnabled())
	state.Description = types.StringValue(lifecycleState.GetDescription())
	state.Priority = types.Int32Value(lifecycleState.GetPriority())

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *LifecycleStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LifecycleStateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResponse, err := r.client.V2025.LifecycleStatesAPI.
		DeleteLifecycleState(ctx, state.IdentityProfileId.ValueString(), state.Id.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Lifecycle State",
			fmt.Sprintf("Could not delete lifecycle state: %s\n\n%v", err.Error(), httpResponse),
		)
		return
	}
}

func (r *LifecycleStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected import format: "identity_profile_id:lifecycle_state_id"
	const importIDSeparator = ":"
	const expectedParts = 2

	idParts := strings.Split(req.ID, importIDSeparator)
	if len(idParts) != expectedParts || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in format 'identity_profile_id:lifecycle_state_id'",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identity_profile_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
