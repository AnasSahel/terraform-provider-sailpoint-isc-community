// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LifecycleStateListDataSourceModel struct {
	IdentityProfileId  types.String          `tfsdk:"identity_profile_id"`
	LifecycleStateList []LifecycleStateModel `tfsdk:"lifecycle_state_list"`
}

type LifecycleStateDataSourceModel struct {
	IdentityProfileId types.String `tfsdk:"identity_profile_id"`
	LifecycleStateModel
}

type LifecycleStateResourceModel struct {
	IdentityProfileId types.String `tfsdk:"identity_profile_id"`
	LifecycleStateModel
}

type LifecycleStateModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	TechnicalName    types.String `tfsdk:"technical_name"`
	Description      types.String `tfsdk:"description"`
	IdentityCount    types.Int32  `tfsdk:"identity_count"`
	AccessProfileIds types.List   `tfsdk:"access_profile_ids"`
	IdentityState    types.String `tfsdk:"identity_state"`
	Priority         types.Int32  `tfsdk:"priority"`

	EmailNotificationOption   *EmailNotificationOptionModel   `tfsdk:"email_notification_option"`
	AccessActionConfiguration *AccessActionConfigurationModel `tfsdk:"access_action_configuration"`
	AccountActions            []AccountActionModel            `tfsdk:"account_actions"`

	Created  types.String `tfsdk:"created"`
	Modified types.String `tfsdk:"modified"`
}

type EmailNotificationOptionModel struct {
	NotifyManagers      types.Bool `tfsdk:"notify_managers"`
	NotifyAllAdmins     types.Bool `tfsdk:"notify_all_admins"`
	NotifySpecificUsers types.Bool `tfsdk:"notify_specific_users"`
	EmailAddressList    types.List `tfsdk:"email_address_list"`
}

type AccessActionConfigurationModel struct {
	RemoveAllAccessEnabled types.Bool `tfsdk:"remove_all_access_enabled"`
}

type AccountActionModel struct {
	Action           types.String `tfsdk:"action"`
	SourceIds        types.List   `tfsdk:"source_ids"`
	ExcludeSourceIds types.List   `tfsdk:"exclude_source_ids"`
	AllSources       types.Bool   `tfsdk:"all_sources"`
}
