// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EmailNotificationOption represents the Terraform model for email notification configuration.
type EmailNotificationOption struct {
	NotifyManagers      types.Bool `tfsdk:"notify_managers"`
	NotifyAllAdmins     types.Bool `tfsdk:"notify_all_admins"`
	NotifySpecificUsers types.Bool `tfsdk:"notify_specific_users"`
	EmailAddressList    types.List `tfsdk:"email_address_list"` // List of strings
}

// AccountAction represents the Terraform model for an account action configuration.
type AccountAction struct {
	Action           types.String `tfsdk:"action"`             // ENABLE, DISABLE, DELETE
	SourceIds        types.List   `tfsdk:"source_ids"`         // List of strings
	ExcludeSourceIds types.List   `tfsdk:"exclude_source_ids"` // List of strings
	AllSources       types.Bool   `tfsdk:"all_sources"`
}

// AccessActionConfiguration represents the Terraform model for access action configuration.
type AccessActionConfiguration struct {
	RemoveAllAccessEnabled types.Bool `tfsdk:"remove_all_access_enabled"`
}

// LifecycleState represents the Terraform model for a SailPoint Lifecycle State.
type LifecycleState struct {
	ID                        types.String               `tfsdk:"id"`
	IdentityProfileID         types.String               `tfsdk:"identity_profile_id"` // Parent identity profile ID (path parameter)
	Name                      types.String               `tfsdk:"name"`
	TechnicalName             types.String               `tfsdk:"technical_name"`
	Enabled                   types.Bool                 `tfsdk:"enabled"`
	Description               types.String               `tfsdk:"description"`
	IdentityCount             types.Int64                `tfsdk:"identity_count"`
	EmailNotificationOption   *EmailNotificationOption   `tfsdk:"email_notification_option"`
	AccountActions            []AccountAction            `tfsdk:"account_actions"`
	AccessProfileIds          types.List                 `tfsdk:"access_profile_ids"` // List of strings
	IdentityState             types.String               `tfsdk:"identity_state"`
	AccessActionConfiguration *AccessActionConfiguration `tfsdk:"access_action_configuration"`
	Priority                  types.Int64                `tfsdk:"priority"`
	Created                   types.String               `tfsdk:"created"`
	Modified                  types.String               `tfsdk:"modified"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API LifecycleState.
func (ls *LifecycleState) ConvertToSailPoint(ctx context.Context) (*client.LifecycleState, error) {
	if ls == nil {
		return nil, nil
	}

	state := &client.LifecycleState{
		Name:          ls.Name.ValueString(),
		TechnicalName: ls.TechnicalName.ValueString(),
	}

	// Optional fields
	if !ls.Description.IsNull() && !ls.Description.IsUnknown() {
		description := ls.Description.ValueString()
		state.Description = &description
	}

	if !ls.Enabled.IsNull() && !ls.Enabled.IsUnknown() {
		enabled := ls.Enabled.ValueBool()
		state.Enabled = &enabled
	}

	if !ls.Priority.IsNull() && !ls.Priority.IsUnknown() {
		priority := int32(ls.Priority.ValueInt64())
		state.Priority = &priority
	}

	if !ls.IdentityState.IsNull() && !ls.IdentityState.IsUnknown() {
		identityState := ls.IdentityState.ValueString()
		state.IdentityState = &identityState
	}

	// EmailNotificationOption
	if ls.EmailNotificationOption != nil {
		emailOpt := &client.EmailNotificationOption{}

		if !ls.EmailNotificationOption.NotifyManagers.IsNull() && !ls.EmailNotificationOption.NotifyManagers.IsUnknown() {
			notifyManagers := ls.EmailNotificationOption.NotifyManagers.ValueBool()
			emailOpt.NotifyManagers = &notifyManagers
		}

		if !ls.EmailNotificationOption.NotifyAllAdmins.IsNull() && !ls.EmailNotificationOption.NotifyAllAdmins.IsUnknown() {
			notifyAllAdmins := ls.EmailNotificationOption.NotifyAllAdmins.ValueBool()
			emailOpt.NotifyAllAdmins = &notifyAllAdmins
		}

		if !ls.EmailNotificationOption.NotifySpecificUsers.IsNull() && !ls.EmailNotificationOption.NotifySpecificUsers.IsUnknown() {
			notifySpecificUsers := ls.EmailNotificationOption.NotifySpecificUsers.ValueBool()
			emailOpt.NotifySpecificUsers = &notifySpecificUsers
		}

		if !ls.EmailNotificationOption.EmailAddressList.IsNull() && !ls.EmailNotificationOption.EmailAddressList.IsUnknown() {
			var emailAddresses []string
			diags := ls.EmailNotificationOption.EmailAddressList.ElementsAs(ctx, &emailAddresses, false)
			if !diags.HasError() && len(emailAddresses) > 0 {
				emailOpt.EmailAddressList = emailAddresses
			}
		}

		state.EmailNotificationOption = emailOpt
	}

	// AccountActions
	if len(ls.AccountActions) > 0 {
		accountActions := make([]client.AccountAction, len(ls.AccountActions))
		for i, action := range ls.AccountActions {
			accountActions[i] = client.AccountAction{
				Action: action.Action.ValueString(),
			}

			if !action.AllSources.IsNull() && !action.AllSources.IsUnknown() {
				allSources := action.AllSources.ValueBool()
				accountActions[i].AllSources = &allSources
			}

			if !action.SourceIds.IsNull() && !action.SourceIds.IsUnknown() {
				var sourceIds []string
				diags := action.SourceIds.ElementsAs(ctx, &sourceIds, false)
				if !diags.HasError() && len(sourceIds) > 0 {
					accountActions[i].SourceIds = sourceIds
				}
			}

			if !action.ExcludeSourceIds.IsNull() && !action.ExcludeSourceIds.IsUnknown() {
				var excludeSourceIds []string
				diags := action.ExcludeSourceIds.ElementsAs(ctx, &excludeSourceIds, false)
				if !diags.HasError() && len(excludeSourceIds) > 0 {
					accountActions[i].ExcludeSourceIds = excludeSourceIds
				}
			}
		}
		state.AccountActions = accountActions
	}

	// AccessProfileIds
	if !ls.AccessProfileIds.IsNull() && !ls.AccessProfileIds.IsUnknown() {
		var accessProfileIds []string
		diags := ls.AccessProfileIds.ElementsAs(ctx, &accessProfileIds, false)
		if !diags.HasError() && len(accessProfileIds) > 0 {
			state.AccessProfileIds = accessProfileIds
		}
	}

	// AccessActionConfiguration
	if ls.AccessActionConfiguration != nil {
		accessActionConfig := &client.AccessActionConfiguration{}

		if !ls.AccessActionConfiguration.RemoveAllAccessEnabled.IsNull() &&
			!ls.AccessActionConfiguration.RemoveAllAccessEnabled.IsUnknown() {
			removeAllAccessEnabled := ls.AccessActionConfiguration.RemoveAllAccessEnabled.ValueBool()
			accessActionConfig.RemoveAllAccessEnabled = &removeAllAccessEnabled
		}

		state.AccessActionConfiguration = accessActionConfig
	}

	return state, nil
}

// ConvertFromSailPoint converts a SailPoint API LifecycleState to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (ls *LifecycleState) ConvertFromSailPoint(ctx context.Context, state *client.LifecycleState, identityProfileID string, includeNull bool) error {
	if ls == nil || state == nil {
		return nil
	}

	// Set IDs
	ls.ID = types.StringValue(state.ID)
	ls.IdentityProfileID = types.StringValue(identityProfileID)

	// Required fields
	ls.Name = types.StringValue(state.Name)
	ls.TechnicalName = types.StringValue(state.TechnicalName)

	// Optional fields with null handling
	if state.Description != nil {
		ls.Description = types.StringValue(*state.Description)
	} else if includeNull {
		ls.Description = types.StringNull()
	}

	if state.Enabled != nil {
		ls.Enabled = types.BoolValue(*state.Enabled)
	} else if includeNull {
		ls.Enabled = types.BoolNull()
	}

	if state.Priority != nil {
		ls.Priority = types.Int64Value(int64(*state.Priority))
	} else if includeNull {
		ls.Priority = types.Int64Null()
	}

	if state.IdentityState != nil {
		ls.IdentityState = types.StringValue(*state.IdentityState)
	} else if includeNull {
		ls.IdentityState = types.StringNull()
	}

	// Computed field
	if state.IdentityCount != nil {
		ls.IdentityCount = types.Int64Value(int64(*state.IdentityCount))
	} else if includeNull {
		ls.IdentityCount = types.Int64Null()
	}

	// Timestamps
	if state.Created != nil {
		ls.Created = types.StringValue(*state.Created)
	} else if includeNull {
		ls.Created = types.StringNull()
	}

	if state.Modified != nil {
		ls.Modified = types.StringValue(*state.Modified)
	} else if includeNull {
		ls.Modified = types.StringNull()
	}

	// EmailNotificationOption - populate if API returned it (even with default values)
	if state.EmailNotificationOption != nil {
		ls.EmailNotificationOption = &EmailNotificationOption{}

		if state.EmailNotificationOption.NotifyManagers != nil {
			ls.EmailNotificationOption.NotifyManagers = types.BoolValue(*state.EmailNotificationOption.NotifyManagers)
		} else if includeNull {
			ls.EmailNotificationOption.NotifyManagers = types.BoolNull()
		}

		if state.EmailNotificationOption.NotifyAllAdmins != nil {
			ls.EmailNotificationOption.NotifyAllAdmins = types.BoolValue(*state.EmailNotificationOption.NotifyAllAdmins)
		} else if includeNull {
			ls.EmailNotificationOption.NotifyAllAdmins = types.BoolNull()
		}

		if state.EmailNotificationOption.NotifySpecificUsers != nil {
			ls.EmailNotificationOption.NotifySpecificUsers = types.BoolValue(*state.EmailNotificationOption.NotifySpecificUsers)
		} else if includeNull {
			ls.EmailNotificationOption.NotifySpecificUsers = types.BoolNull()
		}

		if len(state.EmailNotificationOption.EmailAddressList) > 0 {
			emailElements := make([]attr.Value, len(state.EmailNotificationOption.EmailAddressList))
			for i, email := range state.EmailNotificationOption.EmailAddressList {
				emailElements[i] = types.StringValue(email)
			}
			listValue, diags := types.ListValue(types.StringType, emailElements)
			if diags.HasError() {
				return fmt.Errorf("error creating email_address_list: %v", diags)
			}
			ls.EmailNotificationOption.EmailAddressList = listValue
		} else if includeNull {
			ls.EmailNotificationOption.EmailAddressList = types.ListNull(types.StringType)
		}
	} else if includeNull {
		ls.EmailNotificationOption = nil
	}

	// AccountActions
	if len(state.AccountActions) > 0 {
		ls.AccountActions = make([]AccountAction, len(state.AccountActions))
		for i, action := range state.AccountActions {
			ls.AccountActions[i].Action = types.StringValue(action.Action)

			if action.AllSources != nil {
				ls.AccountActions[i].AllSources = types.BoolValue(*action.AllSources)
			} else if includeNull {
				ls.AccountActions[i].AllSources = types.BoolNull()
			}

			if len(action.SourceIds) > 0 {
				sourceElements := make([]attr.Value, len(action.SourceIds))
				for j, sourceID := range action.SourceIds {
					sourceElements[j] = types.StringValue(sourceID)
				}
				listValue, diags := types.ListValue(types.StringType, sourceElements)
				if diags.HasError() {
					return fmt.Errorf("error creating source_ids list: %v", diags)
				}
				ls.AccountActions[i].SourceIds = listValue
			} else if includeNull {
				ls.AccountActions[i].SourceIds = types.ListNull(types.StringType)
			}

			if len(action.ExcludeSourceIds) > 0 {
				excludeElements := make([]attr.Value, len(action.ExcludeSourceIds))
				for j, excludeID := range action.ExcludeSourceIds {
					excludeElements[j] = types.StringValue(excludeID)
				}
				listValue, diags := types.ListValue(types.StringType, excludeElements)
				if diags.HasError() {
					return fmt.Errorf("error creating exclude_source_ids list: %v", diags)
				}
				ls.AccountActions[i].ExcludeSourceIds = listValue
			} else if includeNull {
				ls.AccountActions[i].ExcludeSourceIds = types.ListNull(types.StringType)
			}
		}
	} else if includeNull {
		ls.AccountActions = nil
	}

	// AccessProfileIds
	if len(state.AccessProfileIds) > 0 {
		accessProfileElements := make([]attr.Value, len(state.AccessProfileIds))
		for i, profileID := range state.AccessProfileIds {
			accessProfileElements[i] = types.StringValue(profileID)
		}
		listValue, diags := types.ListValue(types.StringType, accessProfileElements)
		if diags.HasError() {
			return fmt.Errorf("error creating access_profile_ids list: %v", diags)
		}
		ls.AccessProfileIds = listValue
	} else if includeNull {
		ls.AccessProfileIds = types.ListNull(types.StringType)
	}

	// AccessActionConfiguration - populate whatever the API returns
	if state.AccessActionConfiguration != nil {
		ls.AccessActionConfiguration = &AccessActionConfiguration{}

		if state.AccessActionConfiguration.RemoveAllAccessEnabled != nil {
			ls.AccessActionConfiguration.RemoveAllAccessEnabled = types.BoolValue(*state.AccessActionConfiguration.RemoveAllAccessEnabled)
		} else if includeNull {
			ls.AccessActionConfiguration.RemoveAllAccessEnabled = types.BoolNull()
		}
	} else if includeNull {
		ls.AccessActionConfiguration = nil
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (ls *LifecycleState) ConvertFromSailPointForResource(ctx context.Context, state *client.LifecycleState, identityProfileID string) error {
	return ls.ConvertFromSailPoint(ctx, state, identityProfileID, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (ls *LifecycleState) ConvertFromSailPointForDataSource(ctx context.Context, state *client.LifecycleState, identityProfileID string) error {
	return ls.ConvertFromSailPoint(ctx, state, identityProfileID, false)
}

// GeneratePatchOperations generates JSON Patch operations for updating a lifecycle state.
// Only updatable fields are included: name, description, enabled, email_notification_option,
// account_actions, access_profile_ids, identity_state, access_action_configuration, priority.
func (ls *LifecycleState) GeneratePatchOperations(ctx context.Context, newState *LifecycleState) ([]map[string]interface{}, error) {
	operations := make([]map[string]interface{}, 0)

	// Name
	if !ls.Name.Equal(newState.Name) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/name",
			"value": newState.Name.ValueString(),
		})
	}

	// Description
	if !ls.Description.Equal(newState.Description) {
		if newState.Description.IsNull() {
			// Remove description if it's being set to null
			if !ls.Description.IsNull() {
				operations = append(operations, map[string]interface{}{
					"op":   "remove",
					"path": "/description",
				})
			}
		} else {
			// Use 'add' if old state didn't have description, 'replace' if it did
			op := "replace"
			if ls.Description.IsNull() {
				op = "add"
			}
			operations = append(operations, map[string]interface{}{
				"op":    op,
				"path":  "/description",
				"value": newState.Description.ValueString(),
			})
		}
	}

	// Enabled
	// Skip if the new value is unknown (will be computed by API)
	if !newState.Enabled.IsUnknown() && !ls.Enabled.Equal(newState.Enabled) {
		if newState.Enabled.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/enabled",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/enabled",
				"value": newState.Enabled.ValueBool(),
			})
		}
	}

	// Priority
	// Skip if the new value is unknown (will be computed by API)
	if !newState.Priority.IsUnknown() && !ls.Priority.Equal(newState.Priority) {
		if newState.Priority.IsNull() {
			if !ls.Priority.IsNull() {
				operations = append(operations, map[string]interface{}{
					"op":   "remove",
					"path": "/priority",
				})
			}
		} else {
			op := "replace"
			if ls.Priority.IsNull() {
				op = "add"
			}
			operations = append(operations, map[string]interface{}{
				"op":    op,
				"path":  "/priority",
				"value": int32(newState.Priority.ValueInt64()),
			})
		}
	}

	// IdentityState
	// Skip if the new value is unknown (will be computed by API)
	if !newState.IdentityState.IsUnknown() && !ls.IdentityState.Equal(newState.IdentityState) {
		if newState.IdentityState.IsNull() {
			if !ls.IdentityState.IsNull() {
				operations = append(operations, map[string]interface{}{
					"op":   "remove",
					"path": "/identityState",
				})
			}
		} else {
			op := "replace"
			if ls.IdentityState.IsNull() {
				op = "add"
			}
			operations = append(operations, map[string]interface{}{
				"op":    op,
				"path":  "/identityState",
				"value": newState.IdentityState.ValueString(),
			})
		}
	}

	// EmailNotificationOption
	if !emailNotificationOptionEqual(ls.EmailNotificationOption, newState.EmailNotificationOption) {
		if newState.EmailNotificationOption == nil {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/emailNotificationOption",
			})
		} else {
			emailOpt := map[string]interface{}{}

			if !newState.EmailNotificationOption.NotifyManagers.IsNull() && !newState.EmailNotificationOption.NotifyManagers.IsUnknown() {
				emailOpt["notifyManagers"] = newState.EmailNotificationOption.NotifyManagers.ValueBool()
			}

			if !newState.EmailNotificationOption.NotifyAllAdmins.IsNull() && !newState.EmailNotificationOption.NotifyAllAdmins.IsUnknown() {
				emailOpt["notifyAllAdmins"] = newState.EmailNotificationOption.NotifyAllAdmins.ValueBool()
			}

			if !newState.EmailNotificationOption.NotifySpecificUsers.IsNull() && !newState.EmailNotificationOption.NotifySpecificUsers.IsUnknown() {
				emailOpt["notifySpecificUsers"] = newState.EmailNotificationOption.NotifySpecificUsers.ValueBool()
			}

			if !newState.EmailNotificationOption.EmailAddressList.IsNull() && !newState.EmailNotificationOption.EmailAddressList.IsUnknown() {
				var emailAddresses []string
				diags := newState.EmailNotificationOption.EmailAddressList.ElementsAs(ctx, &emailAddresses, false)
				if !diags.HasError() {
					emailOpt["emailAddressList"] = emailAddresses
				}
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/emailNotificationOption",
				"value": emailOpt,
			})
		}
	}

	// AccountActions
	if !accountActionsEqual(ls.AccountActions, newState.AccountActions) {
		if len(newState.AccountActions) == 0 {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/accountActions",
			})
		} else {
			accountActions := make([]map[string]interface{}, len(newState.AccountActions))
			for i, action := range newState.AccountActions {
				accountAction := map[string]interface{}{
					"action": action.Action.ValueString(),
				}

				if !action.AllSources.IsNull() && !action.AllSources.IsUnknown() {
					accountAction["allSources"] = action.AllSources.ValueBool()
				}

				if !action.SourceIds.IsNull() && !action.SourceIds.IsUnknown() {
					var sourceIds []string
					diags := action.SourceIds.ElementsAs(ctx, &sourceIds, false)
					if !diags.HasError() {
						accountAction["sourceIds"] = sourceIds
					}
				}

				if !action.ExcludeSourceIds.IsNull() && !action.ExcludeSourceIds.IsUnknown() {
					var excludeSourceIds []string
					diags := action.ExcludeSourceIds.ElementsAs(ctx, &excludeSourceIds, false)
					if !diags.HasError() {
						accountAction["excludeSourceIds"] = excludeSourceIds
					}
				}

				accountActions[i] = accountAction
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/accountActions",
				"value": accountActions,
			})
		}
	}

	// AccessProfileIds
	if !ls.AccessProfileIds.Equal(newState.AccessProfileIds) {
		if newState.AccessProfileIds.IsNull() || len(newState.AccessProfileIds.Elements()) == 0 {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/accessProfileIds",
			})
		} else {
			var accessProfileIds []string
			diags := newState.AccessProfileIds.ElementsAs(ctx, &accessProfileIds, false)
			if !diags.HasError() {
				operations = append(operations, map[string]interface{}{
					"op":    "replace",
					"path":  "/accessProfileIds",
					"value": accessProfileIds,
				})
			}
		}
	}

	// AccessActionConfiguration
	if !accessActionConfigurationEqual(ls.AccessActionConfiguration, newState.AccessActionConfiguration) {
		if newState.AccessActionConfiguration == nil {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/accessActionConfiguration",
			})
		} else {
			accessActionConfig := map[string]interface{}{}

			if !newState.AccessActionConfiguration.RemoveAllAccessEnabled.IsNull() && !newState.AccessActionConfiguration.RemoveAllAccessEnabled.IsUnknown() {
				accessActionConfig["removeAllAccessEnabled"] = newState.AccessActionConfiguration.RemoveAllAccessEnabled.ValueBool()
			}

			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/accessActionConfiguration",
				"value": accessActionConfig,
			})
		}
	}

	return operations, nil
}

// emailNotificationOptionEqual compares two EmailNotificationOption structs for equality.
func emailNotificationOptionEqual(a, b *EmailNotificationOption) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}

	if !a.NotifyManagers.Equal(b.NotifyManagers) ||
		!a.NotifyAllAdmins.Equal(b.NotifyAllAdmins) ||
		!a.NotifySpecificUsers.Equal(b.NotifySpecificUsers) ||
		!a.EmailAddressList.Equal(b.EmailAddressList) {
		return false
	}

	return true
}

// accountActionsEqual compares two AccountAction slices for equality.
func accountActionsEqual(a, b []AccountAction) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Action.Equal(b[i].Action) ||
			!a[i].AllSources.Equal(b[i].AllSources) ||
			!a[i].SourceIds.Equal(b[i].SourceIds) ||
			!a[i].ExcludeSourceIds.Equal(b[i].ExcludeSourceIds) {
			return false
		}
	}

	return true
}

// accessActionConfigurationEqual compares two AccessActionConfiguration structs for equality.
func accessActionConfigurationEqual(a, b *AccessActionConfiguration) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}

	return a.RemoveAllAccessEnabled.Equal(b.RemoveAllAccessEnabled)
}
