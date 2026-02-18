// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// lifecycleStateModel represents the Terraform state for a Lifecycle State resource.
type lifecycleStateModel struct {
	ID                        types.String `tfsdk:"id"`
	IdentityProfileID         types.String `tfsdk:"identity_profile_id"`
	Name                      types.String `tfsdk:"name"`
	TechnicalName             types.String `tfsdk:"technical_name"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
	IdentityState             types.String `tfsdk:"identity_state"`
	Priority                  types.Int32  `tfsdk:"priority"`
	Created                   types.String `tfsdk:"created"`
	Modified                  types.String `tfsdk:"modified"`
	EmailNotificationOption   types.Object `tfsdk:"email_notification_option"`
	AccountActions            types.List   `tfsdk:"account_actions"`
	AccessProfileIds          types.List   `tfsdk:"access_profile_ids"`
	AccessActionConfiguration types.Object `tfsdk:"access_action_configuration"`
}

// emailNotificationOptionModel represents the email notification configuration.
type emailNotificationOptionModel struct {
	NotifyManagers      types.Bool `tfsdk:"notify_managers"`
	NotifyAllAdmins     types.Bool `tfsdk:"notify_all_admins"`
	NotifySpecificUsers types.Bool `tfsdk:"notify_specific_users"`
	EmailAddressList    types.List `tfsdk:"email_address_list"`
}

// accountActionModel represents an account action configuration.
type accountActionModel struct {
	Action           types.String `tfsdk:"action"`
	SourceIds        types.List   `tfsdk:"source_ids"`
	ExcludeSourceIds types.List   `tfsdk:"exclude_source_ids"`
	AllSources       types.Bool   `tfsdk:"all_sources"`
}

// accessActionConfigurationModel represents the access action configuration.
type accessActionConfigurationModel struct {
	RemoveAllAccessEnabled types.Bool `tfsdk:"remove_all_access_enabled"`
}

// emailNotificationOptionAttrTypes defines the attribute types for email notification option.
var emailNotificationOptionAttrTypes = map[string]attr.Type{
	"notify_managers":       types.BoolType,
	"notify_all_admins":     types.BoolType,
	"notify_specific_users": types.BoolType,
	"email_address_list":    types.ListType{ElemType: types.StringType},
}

// accountActionAttrTypes defines the attribute types for account action.
var accountActionAttrTypes = map[string]attr.Type{
	"action":             types.StringType,
	"source_ids":         types.ListType{ElemType: types.StringType},
	"exclude_source_ids": types.ListType{ElemType: types.StringType},
	"all_sources":        types.BoolType,
}

// accessActionConfigurationAttrTypes defines the attribute types for access action configuration.
var accessActionConfigurationAttrTypes = map[string]attr.Type{
	"remove_all_access_enabled": types.BoolType,
}

// FromAPI maps fields from the API model to the Terraform model.
func (m *lifecycleStateModel) FromAPI(ctx context.Context, api client.LifecycleStateAPI, identityProfileID string) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.IdentityProfileID = types.StringValue(identityProfileID)
	m.Name = types.StringValue(api.Name)
	m.TechnicalName = types.StringValue(api.TechnicalName)
	m.Enabled = types.BoolValue(api.Enabled)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map nullable string fields
	m.Description = common.StringOrNull(api.Description)
	m.IdentityState = common.StringOrNull(api.IdentityState)

	// Map Priority (nullable)
	if api.Priority != nil {
		m.Priority = types.Int32Value(*api.Priority)
	} else {
		m.Priority = types.Int32Null()
	}

	// Map EmailNotificationOption
	emailList, d := types.ListValueFrom(ctx, types.StringType, api.EmailNotificationOption.EmailAddressList)
	diags.Append(d...)
	emailNotifObj, d := types.ObjectValue(emailNotificationOptionAttrTypes, map[string]attr.Value{
		"notify_managers":       types.BoolValue(api.EmailNotificationOption.NotifyManagers),
		"notify_all_admins":     types.BoolValue(api.EmailNotificationOption.NotifyAllAdmins),
		"notify_specific_users": types.BoolValue(api.EmailNotificationOption.NotifySpecificUsers),
		"email_address_list":    emailList,
	})
	diags.Append(d...)
	m.EmailNotificationOption = emailNotifObj

	// Map AccountActions
	if len(api.AccountActions) > 0 {
		accountActionObjects := make([]attr.Value, len(api.AccountActions))
		for i, action := range api.AccountActions {
			sourceIdsList, d := types.ListValueFrom(ctx, types.StringType, action.SourceIds)
			diags.Append(d...)
			excludeSourceIdsList, d := types.ListValueFrom(ctx, types.StringType, action.ExcludeSourceIds)
			diags.Append(d...)

			actionObj, d := types.ObjectValue(accountActionAttrTypes, map[string]attr.Value{
				"action":             types.StringValue(action.Action),
				"source_ids":         sourceIdsList,
				"exclude_source_ids": excludeSourceIdsList,
				"all_sources":        types.BoolValue(action.AllSources),
			})
			diags.Append(d...)
			accountActionObjects[i] = actionObj
		}
		accountActionsList, d := types.ListValue(types.ObjectType{AttrTypes: accountActionAttrTypes}, accountActionObjects)
		diags.Append(d...)
		m.AccountActions = accountActionsList
	} else {
		m.AccountActions = types.ListNull(types.ObjectType{AttrTypes: accountActionAttrTypes})
	}

	// Map AccessProfileIds
	if len(api.AccessProfileIds) > 0 {
		accessProfileIdsList, d := types.ListValueFrom(ctx, types.StringType, api.AccessProfileIds)
		diags.Append(d...)
		m.AccessProfileIds = accessProfileIdsList
	} else {
		m.AccessProfileIds = types.ListNull(types.StringType)
	}

	// Map AccessActionConfiguration
	accessActionConfigObj, d := types.ObjectValue(accessActionConfigurationAttrTypes, map[string]attr.Value{
		"remove_all_access_enabled": types.BoolValue(api.AccessActionConfiguration.RemoveAllAccessEnabled),
	})
	diags.Append(d...)
	m.AccessActionConfiguration = accessActionConfigObj

	return diags
}

// ToAPI maps fields from the Terraform model to the API create request.
func (m *lifecycleStateModel) ToAPI(ctx context.Context) (client.LifecycleStateCreateAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiRequest := client.LifecycleStateCreateAPI{
		Name:          m.Name.ValueString(),
		TechnicalName: m.TechnicalName.ValueString(),
		Enabled:       m.Enabled.ValueBool(),
	}

	// Map Description (optional, nullable)
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		desc := m.Description.ValueString()
		apiRequest.Description = &desc
	}

	// Map IdentityState (optional, nullable)
	if !m.IdentityState.IsNull() && !m.IdentityState.IsUnknown() {
		identityState := m.IdentityState.ValueString()
		apiRequest.IdentityState = &identityState
	}

	// Map Priority (optional, nullable)
	if !m.Priority.IsNull() && !m.Priority.IsUnknown() {
		priority := m.Priority.ValueInt32()
		apiRequest.Priority = &priority
	}

	// Map EmailNotificationOption
	if !m.EmailNotificationOption.IsNull() && !m.EmailNotificationOption.IsUnknown() {
		var emailNotifModel emailNotificationOptionModel
		d := m.EmailNotificationOption.As(ctx, &emailNotifModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			var emailList []string
			if !emailNotifModel.EmailAddressList.IsNull() && !emailNotifModel.EmailAddressList.IsUnknown() {
				d := emailNotifModel.EmailAddressList.ElementsAs(ctx, &emailList, false)
				diags.Append(d...)
			}
			apiRequest.EmailNotificationOption = client.EmailNotificationOptionAPI{
				NotifyManagers:      emailNotifModel.NotifyManagers.ValueBool(),
				NotifyAllAdmins:     emailNotifModel.NotifyAllAdmins.ValueBool(),
				NotifySpecificUsers: emailNotifModel.NotifySpecificUsers.ValueBool(),
				EmailAddressList:    emailList,
			}
		}
	}

	// Map AccountActions
	if !m.AccountActions.IsNull() && !m.AccountActions.IsUnknown() {
		var accountActionsModels []accountActionModel
		d := m.AccountActions.ElementsAs(ctx, &accountActionsModels, false)
		diags.Append(d...)
		if !diags.HasError() {
			apiRequest.AccountActions = make([]client.AccountActionAPI, len(accountActionsModels))
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
				apiRequest.AccountActions[i] = client.AccountActionAPI{
					Action:           actionModel.Action.ValueString(),
					SourceIds:        sourceIds,
					ExcludeSourceIds: excludeSourceIds,
					AllSources:       actionModel.AllSources.ValueBool(),
				}
			}
		}
	}

	// Map AccessProfileIds
	if !m.AccessProfileIds.IsNull() && !m.AccessProfileIds.IsUnknown() {
		d := m.AccessProfileIds.ElementsAs(ctx, &apiRequest.AccessProfileIds, false)
		diags.Append(d...)
	}

	// Map AccessActionConfiguration
	if !m.AccessActionConfiguration.IsNull() && !m.AccessActionConfiguration.IsUnknown() {
		var accessActionConfigModel accessActionConfigurationModel
		d := m.AccessActionConfiguration.As(ctx, &accessActionConfigModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			apiRequest.AccessActionConfiguration = client.AccessActionConfigurationAPI{
				RemoveAllAccessEnabled: accessActionConfigModel.RemoveAllAccessEnabled.ValueBool(),
			}
		}
	}

	return apiRequest, diags
}

// ToPatchOperations compares the plan (m) with the current state and generates JSON Patch operations
// for fields that have changed.
func (m *lifecycleStateModel) ToPatchOperations(ctx context.Context, state *lifecycleStateModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diags diag.Diagnostics
	var operations []client.JSONPatchOperation

	// Enabled
	if !m.Enabled.Equal(state.Enabled) {
		operations = append(operations, client.NewReplacePatch("/enabled", m.Enabled.ValueBool()))
	}

	// Description
	if !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
			operations = append(operations, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			operations = append(operations, client.NewRemovePatch("/description"))
		}
	}

	// Priority
	if !m.Priority.Equal(state.Priority) {
		if !m.Priority.IsNull() {
			operations = append(operations, client.NewReplacePatch("/priority", m.Priority.ValueInt32()))
		} else {
			operations = append(operations, client.NewRemovePatch("/priority"))
		}
	}

	// EmailNotificationOption
	if !m.EmailNotificationOption.Equal(state.EmailNotificationOption) {
		var emailNotifModel emailNotificationOptionModel
		d := m.EmailNotificationOption.As(ctx, &emailNotifModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			var emailList []string
			if !emailNotifModel.EmailAddressList.IsNull() && !emailNotifModel.EmailAddressList.IsUnknown() {
				d := emailNotifModel.EmailAddressList.ElementsAs(ctx, &emailList, false)
				diags.Append(d...)
			}
			operations = append(operations, client.NewReplacePatch("/emailNotificationOption", client.EmailNotificationOptionAPI{
				NotifyManagers:      emailNotifModel.NotifyManagers.ValueBool(),
				NotifyAllAdmins:     emailNotifModel.NotifyAllAdmins.ValueBool(),
				NotifySpecificUsers: emailNotifModel.NotifySpecificUsers.ValueBool(),
				EmailAddressList:    emailList,
			}))
		}
	}

	// AccountActions
	if !m.AccountActions.Equal(state.AccountActions) {
		if m.AccountActions.IsNull() {
			operations = append(operations, client.NewReplacePatch("/accountActions", []client.AccountActionAPI{}))
		} else {
			var accountActionsModels []accountActionModel
			d := m.AccountActions.ElementsAs(ctx, &accountActionsModels, false)
			diags.Append(d...)
			if !diags.HasError() {
				accountActionsAPI := make([]client.AccountActionAPI, len(accountActionsModels))
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
					accountActionsAPI[i] = client.AccountActionAPI{
						Action:           actionModel.Action.ValueString(),
						SourceIds:        sourceIds,
						ExcludeSourceIds: excludeSourceIds,
						AllSources:       actionModel.AllSources.ValueBool(),
					}
				}
				operations = append(operations, client.NewReplacePatch("/accountActions", accountActionsAPI))
			}
		}
	}

	// AccessProfileIds
	if !m.AccessProfileIds.Equal(state.AccessProfileIds) {
		if m.AccessProfileIds.IsNull() {
			operations = append(operations, client.NewReplacePatch("/accessProfileIds", []string{}))
		} else {
			var accessProfileIds []string
			d := m.AccessProfileIds.ElementsAs(ctx, &accessProfileIds, false)
			diags.Append(d...)
			if !diags.HasError() {
				operations = append(operations, client.NewReplacePatch("/accessProfileIds", accessProfileIds))
			}
		}
	}

	// AccessActionConfiguration
	if !m.AccessActionConfiguration.Equal(state.AccessActionConfiguration) {
		var accessActionConfigModel accessActionConfigurationModel
		d := m.AccessActionConfiguration.As(ctx, &accessActionConfigModel, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if !diags.HasError() {
			operations = append(operations, client.NewReplacePatch("/accessActionConfiguration", client.AccessActionConfigurationAPI{
				RemoveAllAccessEnabled: accessActionConfigModel.RemoveAllAccessEnabled.ValueBool(),
			}))
		}
	}

	return operations, diags
}

// lifecycleStateDataSourceModel embeds the resource model and adds server-managed read-only fields.
type lifecycleStateDataSourceModel struct {
	lifecycleStateModel
	IdentityCount types.Int32 `tfsdk:"identity_count"`
}

// FromAPI maps fields from the API response to the data source model.
func (m *lifecycleStateDataSourceModel) FromAPI(ctx context.Context, api client.LifecycleStateAPI, identityProfileID string) diag.Diagnostics {
	// Map shared fields via the embedded resource model
	diagnostics := m.lifecycleStateModel.FromAPI(ctx, api, identityProfileID)

	// Map server-managed read-only fields
	m.IdentityCount = types.Int32Value(api.IdentityCount)

	return diagnostics
}
