// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sourceProvisioningPolicyResource{}
	_ resource.ResourceWithConfigure   = &sourceProvisioningPolicyResource{}
	_ resource.ResourceWithImportState = &sourceProvisioningPolicyResource{}
)

type sourceProvisioningPolicyResource struct {
	client *client.Client
}

// NewSourceProvisioningPolicyResource creates a new resource for Source Provisioning Policy.
func NewSourceProvisioningPolicyResource() resource.Resource {
	return &sourceProvisioningPolicyResource{}
}

// Metadata implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_provisioning_policy"
}

// Configure implements resource.ResourceWithConfigure.
func (r *sourceProvisioningPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source provisioning policy resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a SailPoint Source Provisioning Policy.",
		MarkdownDescription: "Manages a SailPoint Source Provisioning Policy. " +
			"A provisioning policy defines the fields and transformations required for a specific provisioning operation type. " +
			"Use this resource to create, update, and delete provisioning policies for a source.",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the source this provisioning policy belongs to. Changing this forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"usage_type": schema.StringAttribute{
				MarkdownDescription: "The usage type of the provisioning policy (e.g., `CREATE`, `UPDATE`, `DELETE`, `ENABLE`, `DISABLE`, `ASSIGN`, `UNASSIGN`, `CREATE_GROUP`, `UPDATE_GROUP`, `DELETE_GROUP`, `REGISTER`, `CREATE_IDENTITY`, `UPDATE_IDENTITY`, `EDIT_GROUP`, `UNLOCK`, `CHANGE_PASSWORD`). This value is the unique identifier for the policy within a source. Changing this forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the provisioning policy.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the provisioning policy.",
				Optional:            true,
			},
			"fields": schema.ListNestedAttribute{
				MarkdownDescription: "The list of fields defined by the provisioning policy.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the field.",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the field. Can be null.",
							Optional:            true,
							Computed:            true,
						},
						"is_required": schema.BoolAttribute{
							MarkdownDescription: "Whether the field is required.",
							Optional:            true,
							Computed:            true,
						},
						"is_multi_valued": schema.BoolAttribute{
							MarkdownDescription: "Whether the field supports multiple values.",
							Optional:            true,
							Computed:            true,
						},
						"transform": schema.StringAttribute{
							MarkdownDescription: "The transformation applied to the field as a JSON object.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
						"attributes": schema.StringAttribute{
							MarkdownDescription: "Additional attributes for the field as a JSON object.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
						},
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Getting plan for source provisioning policy resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := plan.SourceID.ValueString()
	usageType := plan.UsageType.ValueString()

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping source provisioning policy resource model to API create request", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
		"name":       plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPICreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the provisioning policy via the API client
	tflog.Debug(ctx, "Creating source provisioning policy via SailPoint API", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
		"name":       plan.Name.ValueString(),
	})
	createResponse, err := r.client.CreateProvisioningPolicy(ctx, sourceID, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Source Provisioning Policy",
			fmt.Sprintf("Could not create SailPoint Source Provisioning Policy with usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Source Provisioning Policy", map[string]any{
			"source_id":  sourceID,
			"usage_type": usageType,
			"name":       plan.Name.ValueString(),
			"error":      err.Error(),
		})
		return
	}

	if createResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Source Provisioning Policy",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Read back the authoritative state from the API to avoid inconsistencies
	// (the POST response may not include all fields accurately)
	tflog.Debug(ctx, "Reading back source provisioning policy after create for authoritative state", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	policyResponse, err := r.client.GetProvisioningPolicy(ctx, sourceID, usageType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy After Create",
			fmt.Sprintf("Provisioning policy was created successfully but could not be read back. Usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		return
	}

	if policyResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy After Create",
			"Received nil response from SailPoint API when reading back created policy",
		)
		return
	}

	// Map the authoritative GET response to the resource model
	var state sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Mapping SailPoint Source Provisioning Policy API response to resource model", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, policyResponse, sourceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for source provisioning policy resource", map[string]any{
		"source_id":  sourceID,
		"usage_type": state.UsageType.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Source Provisioning Policy resource", map[string]any{
		"source_id":   sourceID,
		"usage_type":  state.UsageType.ValueString(),
		"policy_name": state.Name.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Getting state for source provisioning policy resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	usageType := state.UsageType.ValueString()

	// Read the provisioning policy from SailPoint
	tflog.Debug(ctx, "Fetching source provisioning policy from SailPoint", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	policyResponse, err := r.client.GetProvisioningPolicy(ctx, sourceID, usageType)
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Source Provisioning Policy not found, removing from state", map[string]any{
				"source_id":  sourceID,
				"usage_type": usageType,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy",
			fmt.Sprintf("Could not read SailPoint Source Provisioning Policy with usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Source Provisioning Policy", map[string]any{
			"source_id":  sourceID,
			"usage_type": usageType,
			"error":      err.Error(),
		})
		return
	}

	if policyResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Source Provisioning Policy API response to resource model", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, policyResponse, sourceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for source provisioning policy resource", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Source Provisioning Policy resource", map[string]any{
		"source_id":   sourceID,
		"usage_type":  usageType,
		"policy_name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Getting plan for source provisioning policy resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the identifiers
	var state sourceProvisioningPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	usageType := state.UsageType.ValueString()

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping source provisioning policy resource model to API update request", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	apiUpdateRequest, diags := plan.ToAPIUpdateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the provisioning policy via the API client (PUT)
	tflog.Debug(ctx, "Updating source provisioning policy via SailPoint API", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	_, err := r.client.UpdateProvisioningPolicy(ctx, sourceID, usageType, &apiUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Source Provisioning Policy",
			fmt.Sprintf("Could not update SailPoint Source Provisioning Policy with usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Source Provisioning Policy", map[string]any{
			"source_id":  sourceID,
			"usage_type": usageType,
			"error":      err.Error(),
		})
		return
	}

	// Read back the authoritative state from the API to avoid inconsistencies
	// (the PUT response may not include all fields accurately)
	tflog.Debug(ctx, "Reading back source provisioning policy after update for authoritative state", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	policyResponse, err := r.client.GetProvisioningPolicy(ctx, sourceID, usageType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy After Update",
			fmt.Sprintf("Provisioning policy was updated successfully but could not be read back. Usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		return
	}

	if policyResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Provisioning Policy After Update",
			"Received nil response from SailPoint API when reading back updated policy",
		)
		return
	}

	// Map the authoritative GET response to the resource model
	var newState sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Mapping SailPoint Source Provisioning Policy API response to resource model", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	resp.Diagnostics.Append(newState.FromSailPointAPI(ctx, policyResponse, sourceID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for source provisioning policy resource", map[string]any{
		"source_id":  sourceID,
		"usage_type": newState.UsageType.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Source Provisioning Policy resource", map[string]any{
		"source_id":   sourceID,
		"usage_type":  newState.UsageType.ValueString(),
		"policy_name": newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *sourceProvisioningPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceProvisioningPolicyResourceModel
	tflog.Debug(ctx, "Getting state for source provisioning policy resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	usageType := state.UsageType.ValueString()

	tflog.Debug(ctx, "Deleting source provisioning policy via SailPoint API", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
	err := r.client.DeleteProvisioningPolicy(ctx, sourceID, usageType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Source Provisioning Policy",
			fmt.Sprintf("Could not delete SailPoint Source Provisioning Policy with usage type %q for source %q: %s",
				usageType, sourceID, err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Source Provisioning Policy", map[string]any{
			"source_id":  sourceID,
			"usage_type": usageType,
			"error":      err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Source Provisioning Policy resource", map[string]any{
		"source_id":   sourceID,
		"usage_type":  usageType,
		"policy_name": state.Name.ValueString(),
	})
}

// ImportState implements resource.ResourceWithImportState.
// Import format: source_id/usage_type.
func (r *sourceProvisioningPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing source provisioning policy resource", map[string]any{
		"import_id": req.ID,
	})

	// Parse the import ID (format: source_id/usage_type)
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: source_id/usage_type, got: %s", req.ID),
		)
		return
	}

	sourceID := parts[0]
	usageType := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), sourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("usage_type"), usageType)...)

	tflog.Info(ctx, "Successfully imported SailPoint Source Provisioning Policy resource", map[string]any{
		"source_id":  sourceID,
		"usage_type": usageType,
	})
}
