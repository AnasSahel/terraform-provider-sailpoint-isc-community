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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the lifecycle state is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_count": schema.Int32Attribute{
				MarkdownDescription: "The number of identities that have this lifecycle state.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"identity_state": schema.StringAttribute{
				MarkdownDescription: "The identity state associated with this lifecycle state. Possible values: `ACTIVE`, `INACTIVE_SHORT_TERM`, `INACTIVE_LONG_TERM`. Can only be set during creation.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
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
				MarkdownDescription: "Email notification configuration for the lifecycle state.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"notify_managers": schema.BoolAttribute{
						MarkdownDescription: "If true, managers are notified of lifecycle state changes.",
						Optional:            true,
						Computed:            true,
					},
					"notify_all_admins": schema.BoolAttribute{
						MarkdownDescription: "If true, all admins are notified of lifecycle state changes.",
						Optional:            true,
						Computed:            true,
					},
					"notify_specific_users": schema.BoolAttribute{
						MarkdownDescription: "If true, users specified in `email_address_list` are notified of lifecycle state changes.",
						Optional:            true,
						Computed:            true,
					},
					"email_address_list": schema.ListAttribute{
						MarkdownDescription: "List of email addresses to notify when `notify_specific_users` is true.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"account_actions": schema.ListNestedAttribute{
				MarkdownDescription: "List of account actions to perform when an identity enters this lifecycle state.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"access_action_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "Access action configuration for the lifecycle state.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"remove_all_access_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, all access is removed when an identity enters this lifecycle state.",
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
	apiCreateRequest, diags := plan.ToAPICreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the lifecycle state via the API client
	tflog.Debug(ctx, "Creating lifecycle state via SailPoint API", map[string]any{
		"identity_profile_id": identityProfileID,
		"name":                plan.Name.ValueString(),
	})
	lifecycleStateAPIResponse, err := r.client.CreateLifecycleState(ctx, identityProfileID, &apiCreateRequest)
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

	if lifecycleStateAPIResponse == nil {
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
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *lifecycleStateAPIResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for lifecycle state resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"id":                  state.ID.ValueString(),
	})
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
	lifecycleStateResponse, err := r.client.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
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

	if lifecycleStateResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Lifecycle State API response to resource model", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *lifecycleStateResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for lifecycle state resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
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
	patchOperations, diags := r.buildPatchOperations(ctx, &state, &plan)
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
	lifecycleStateAPIResponse, err := r.client.UpdateLifecycleState(ctx, identityProfileID, lifecycleStateID, patchOperations)
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

	if lifecycleStateAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState lifecycleStateModel
	tflog.Debug(ctx, "Mapping SailPoint Lifecycle State API response to resource model", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	resp.Diagnostics.Append(newState.FromSailPointAPI(ctx, *lifecycleStateAPIResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for lifecycle state resource", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
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

// buildPatchOperations creates JSON Patch operations for changes between state and plan.
func (r *lifecycleStateResource) buildPatchOperations(ctx context.Context, state, plan *lifecycleStateModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var operations []client.JSONPatchOperation

	// Check for enabled change
	if !plan.Enabled.Equal(state.Enabled) {
		operations = append(operations, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/enabled",
			Value: plan.Enabled.ValueBool(),
		})
	}

	// Check for description change
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/description",
				Value: nil,
			})
		} else {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/description",
				Value: plan.Description.ValueString(),
			})
		}
	}

	// Check for priority change
	if !plan.Priority.Equal(state.Priority) {
		if plan.Priority.IsNull() {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/priority",
				Value: nil,
			})
		} else {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/priority",
				Value: plan.Priority.ValueInt32(),
			})
		}
	}

	// Check for email_notification_option change
	if !plan.EmailNotificationOption.Equal(state.EmailNotificationOption) {
		var emailNotifModel emailNotificationOptionModel
		d := plan.EmailNotificationOption.As(ctx, &emailNotifModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			var emailList []string
			if !emailNotifModel.EmailAddressList.IsNull() && !emailNotifModel.EmailAddressList.IsUnknown() {
				d := emailNotifModel.EmailAddressList.ElementsAs(ctx, &emailList, false)
				diags.Append(d...)
			}
			operations = append(operations, client.JSONPatchOperation{
				Op:   "replace",
				Path: "/emailNotificationOption",
				Value: map[string]interface{}{
					"notifyManagers":      emailNotifModel.NotifyManagers.ValueBool(),
					"notifyAllAdmins":     emailNotifModel.NotifyAllAdmins.ValueBool(),
					"notifySpecificUsers": emailNotifModel.NotifySpecificUsers.ValueBool(),
					"emailAddressList":    emailList,
				},
			})
		}
	}

	// Check for account_actions change
	if !plan.AccountActions.Equal(state.AccountActions) {
		if plan.AccountActions.IsNull() {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/accountActions",
				Value: []interface{}{},
			})
		} else {
			var accountActionsModels []accountActionModel
			d := plan.AccountActions.ElementsAs(ctx, &accountActionsModels, false)
			diags.Append(d...)
			if !diags.HasError() {
				accountActionsAPI := make([]map[string]interface{}, len(accountActionsModels))
				for i, actionModel := range accountActionsModels {
					var sourceIds []string
					if !actionModel.SourceIds.IsNull() && !actionModel.SourceIds.IsUnknown() {
						d := actionModel.SourceIds.ElementsAs(ctx, &sourceIds, false)
						diags.Append(d...)
					}
					var excludeSourceIds []string
					if !actionModel.ExcludeSourceIds.IsNull() && !actionModel.ExcludeSourceIds.IsUnknown() {
						d := actionModel.ExcludeSourceIds.ElementsAs(ctx, &excludeSourceIds, false)
						diags.Append(d...)
					}
					accountActionsAPI[i] = map[string]interface{}{
						"action":           actionModel.Action.ValueString(),
						"sourceIds":        sourceIds,
						"excludeSourceIds": excludeSourceIds,
						"allSources":       actionModel.AllSources.ValueBool(),
					}
				}
				operations = append(operations, client.JSONPatchOperation{
					Op:    "replace",
					Path:  "/accountActions",
					Value: accountActionsAPI,
				})
			}
		}
	}

	// Check for access_profile_ids change
	if !plan.AccessProfileIds.Equal(state.AccessProfileIds) {
		if plan.AccessProfileIds.IsNull() {
			operations = append(operations, client.JSONPatchOperation{
				Op:    "replace",
				Path:  "/accessProfileIds",
				Value: []string{},
			})
		} else {
			var accessProfileIds []string
			d := plan.AccessProfileIds.ElementsAs(ctx, &accessProfileIds, false)
			diags.Append(d...)
			if !diags.HasError() {
				operations = append(operations, client.JSONPatchOperation{
					Op:    "replace",
					Path:  "/accessProfileIds",
					Value: accessProfileIds,
				})
			}
		}
	}

	// Check for access_action_configuration change
	if !plan.AccessActionConfiguration.Equal(state.AccessActionConfiguration) {
		var accessActionConfigModel accessActionConfigurationModel
		d := plan.AccessActionConfiguration.As(ctx, &accessActionConfigModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			operations = append(operations, client.JSONPatchOperation{
				Op:   "replace",
				Path: "/accessActionConfiguration",
				Value: map[string]interface{}{
					"removeAllAccessEnabled": accessActionConfigModel.RemoveAllAccessEnabled.ValueBool(),
				},
			})
		}
	}

	return operations, diags
}
