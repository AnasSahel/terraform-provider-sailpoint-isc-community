// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &lifecycleStateResource{}
	_ resource.ResourceWithConfigure   = &lifecycleStateResource{}
	_ resource.ResourceWithImportState = &lifecycleStateResource{}
)

type lifecycleStateResource struct {
	client *client.Client
}

// NewLifecycleStateResource creates a new resource for Lifecycle State.
func NewLifecycleStateResource() resource.Resource {
	return &lifecycleStateResource{}
}

// Metadata implements resource.Resource.
func (r *lifecycleStateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

// Configure implements resource.ResourceWithConfigure.
func (r *lifecycleStateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "lifecycle_state resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *lifecycleStateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Lifecycle State.",
		MarkdownDescription: "Resource for SailPoint Lifecycle State. Lifecycle states define the different stages an identity can be in within an identity profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the lifecycle state.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the identity profile this lifecycle state belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the lifecycle state. Cannot be changed after creation.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"technical_name": schema.StringAttribute{
				MarkdownDescription: "The technical name of the lifecycle state. This is used for internal purposes and cannot be changed after creation.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the lifecycle state.",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the lifecycle state is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_state": schema.StringAttribute{
				MarkdownDescription: "The identity state associated with this lifecycle state. Possible values: `ACTIVE`, `INACTIVE_SHORT_TERM`, `INACTIVE_LONG_TERM`. Cannot be changed after creation.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "The priority of the lifecycle state. Lower numbers appear first when listing with `?sorters=priority`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the lifecycle state was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the lifecycle state was last modified.",
				Computed:            true,
			},
			"email_notification_option": schema.SingleNestedAttribute{
				MarkdownDescription: "Email notification configuration for the lifecycle state. " +
					"Defaults to all notifications disabled with an empty email list. " +
					"Remove this block from your configuration to reset to defaults.",
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					emailNotificationOptionAttrTypes,
					map[string]attr.Value{
						"notify_managers":       types.BoolValue(false),
						"notify_all_admins":     types.BoolValue(false),
						"notify_specific_users": types.BoolValue(false),
						"email_address_list":    types.ListValueMust(types.StringType, []attr.Value{}),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"notify_managers": schema.BoolAttribute{
						MarkdownDescription: "If true, managers are notified of lifecycle state changes. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"notify_all_admins": schema.BoolAttribute{
						MarkdownDescription: "If true, all admins are notified of lifecycle state changes. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"notify_specific_users": schema.BoolAttribute{
						MarkdownDescription: "If true, users specified in `email_address_list` are notified of lifecycle state changes. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"email_address_list": schema.ListAttribute{
						MarkdownDescription: "List of email addresses to notify when `notify_specific_users` is true. Defaults to an empty list.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"account_actions": schema.ListNestedAttribute{
				MarkdownDescription: "List of account actions to perform when an identity enters this lifecycle state.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							MarkdownDescription: "The action to perform. Possible values: `ENABLE`, `DISABLE`, `DELETE`.",
							Required:            true,
						},
						"source_ids": schema.ListAttribute{
							MarkdownDescription: "List of source IDs to apply the action to. Required if `all_sources` is false.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"exclude_source_ids": schema.ListAttribute{
							MarkdownDescription: "List of source IDs to exclude from the action. Cannot be used with `source_ids`.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"all_sources": schema.BoolAttribute{
							MarkdownDescription: "If true, the action applies to all sources. If true, `source_ids` must not be provided.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"access_profile_ids": schema.ListAttribute{
				MarkdownDescription: "List of access profile IDs associated with this lifecycle state.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"access_action_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "Access action configuration for the lifecycle state. " +
					"Defaults to `remove_all_access_enabled = false`. " +
					"Remove this block from your configuration to reset to defaults.",
				Optional: true,
				Computed: true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					accessActionConfigurationAttrTypes,
					map[string]attr.Value{
						"remove_all_access_enabled": types.BoolValue(false),
					},
				)),
				Attributes: map[string]schema.Attribute{
					"remove_all_access_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, all access is removed when an identity enters this lifecycle state. Defaults to `false`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *lifecycleStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan lifecycleStateModel
	tflog.Debug(ctx, "Getting plan for lifecycle state resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := plan.IdentityProfileID.ValueString()

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping lifecycle state resource model to API create request", map[string]any{
		"identity_profile_id": identityProfileID,
		"name":                plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the lifecycle state via the API client
	tflog.Debug(ctx, "Creating lifecycle state via SailPoint API", map[string]any{
		"identity_profile_id": identityProfileID,
		"name":                plan.Name.ValueString(),
	})
	apiResponse, err := r.client.CreateLifecycleState(ctx, identityProfileID, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Lifecycle State",
			fmt.Sprintf("Could not create SailPoint Lifecycle State %q in identity profile %q: %s",
				plan.Name.ValueString(), identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Lifecycle State", map[string]any{
			"identity_profile_id": identityProfileID,
			"name":                plan.Name.ValueString(),
			"error":               err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state lifecycleStateModel
	tflog.Debug(ctx, "Mapping SailPoint Lifecycle State API response to resource model", map[string]any{
		"identity_profile_id": identityProfileID,
		"name":                plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromAPI(ctx, *apiResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Lifecycle State resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"id":                  state.ID.ValueString(),
		"name":                state.Name.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *lifecycleStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lifecycleStateModel
	tflog.Debug(ctx, "Getting state for lifecycle state resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	// Read the lifecycle state from SailPoint
	tflog.Debug(ctx, "Fetching lifecycle state from SailPoint", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	apiResponse, err := r.client.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Lifecycle State not found, removing from state", map[string]any{
				"identity_profile_id": identityProfileID,
				"lifecycle_state_id":  lifecycleStateID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Lifecycle State",
			fmt.Sprintf("Could not read SailPoint Lifecycle State %q in identity profile %q: %s",
				lifecycleStateID, identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Lifecycle State", map[string]any{
			"identity_profile_id": identityProfileID,
			"lifecycle_state_id":  lifecycleStateID,
			"error":               err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	resp.Diagnostics.Append(state.FromAPI(ctx, *apiResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Lifecycle State resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"name":                state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *lifecycleStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan lifecycleStateModel
	tflog.Debug(ctx, "Getting plan for lifecycle state resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the ID
	var state lifecycleStateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	// Build patch operations for changed fields
	tflog.Debug(ctx, "Building patch operations for lifecycle state update", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	patchOperations, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(patchOperations) == 0 {
		tflog.Debug(ctx, "No changes detected, skipping update", map[string]any{
			"identity_profile_id": identityProfileID,
			"lifecycle_state_id":  lifecycleStateID,
		})
		return
	}

	// Update the lifecycle state via the API client (PATCH)
	tflog.Debug(ctx, "Updating lifecycle state via SailPoint API", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"operations_count":    len(patchOperations),
	})
	apiResponse, err := r.client.UpdateLifecycleState(ctx, identityProfileID, lifecycleStateID, patchOperations)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Lifecycle State",
			fmt.Sprintf("Could not update SailPoint Lifecycle State %q in identity profile %q: %s",
				lifecycleStateID, identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Lifecycle State", map[string]any{
			"identity_profile_id": identityProfileID,
			"lifecycle_state_id":  lifecycleStateID,
			"error":               err.Error(),
		})
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState lifecycleStateModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, *apiResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Lifecycle State resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"name":                newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *lifecycleStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lifecycleStateModel
	tflog.Debug(ctx, "Getting state for lifecycle state resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := state.IdentityProfileID.ValueString()
	lifecycleStateID := state.ID.ValueString()

	tflog.Debug(ctx, "Deleting lifecycle state via SailPoint API", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	err := r.client.DeleteLifecycleState(ctx, identityProfileID, lifecycleStateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Lifecycle State",
			fmt.Sprintf("Could not delete SailPoint Lifecycle State %q in identity profile %q: %s",
				lifecycleStateID, identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Lifecycle State", map[string]any{
			"identity_profile_id": identityProfileID,
			"lifecycle_state_id":  lifecycleStateID,
			"error":               err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Lifecycle State resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
		"name":                state.Name.ValueString(),
	})
}

// ImportState implements resource.ResourceWithImportState.
// Import format: identity_profile_id/lifecycle_state_id.
func (r *lifecycleStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing lifecycle state resource", map[string]any{
		"import_id": req.ID,
	})

	// Parse the import ID (format: identity_profile_id/lifecycle_state_id)
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: identity_profile_id/lifecycle_state_id, got: %s", req.ID),
		)
		return
	}

	identityProfileID := parts[0]
	lifecycleStateID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identity_profile_id"), identityProfileID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lifecycleStateID)...)

	tflog.Info(ctx, "Successfully imported SailPoint Lifecycle State resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
}
