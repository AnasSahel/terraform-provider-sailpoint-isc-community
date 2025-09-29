// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource Schema Functions

func LifecycleStateListDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Lifecycle State List Data Source",
		MarkdownDescription: "Use this data source to retrieve information about all lifecycle states within an identity profile in SailPoint Identity Security Cloud (ISC). Lifecycle states define stages in an identity's lifecycle and control access provisioning and deprovisioning.",
		Attributes: map[string]schema.Attribute{
			"identity_profile_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the identity profile associated with the lifecycle state.",
				MarkdownDescription: "The unique identifier of the identity profile that contains this lifecycle state. Identity profiles define how identities are processed and managed in SailPoint ISC.",
			},
			"lifecycle_state_list": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of lifecycle states associated with the identity profile.",
				MarkdownDescription: "A list containing the retrieved lifecycle state information. When querying a specific lifecycle state by ID, this list will contain exactly one item with the detailed information about that lifecycle state.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: lifecycleStateAttributesForDataSource(),
				},
			},
		},
	}
}

func LifecycleStateDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Lifecycle State Data Source",
		MarkdownDescription: "Use this data source to retrieve information about a specific lifecycle state within an identity profile in SailPoint Identity Security Cloud (ISC). Lifecycle states define stages in an identity's lifecycle and control access provisioning and deprovisioning.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the lifecycle state to retrieve.",
				MarkdownDescription: "The unique identifier of the lifecycle state to retrieve. This is a system-generated UUID that uniquely identifies the lifecycle state within the identity profile.",
			},
			"identity_profile_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the identity profile associated with the lifecycle state.",
				MarkdownDescription: "The unique identifier of the identity profile that contains this lifecycle state. Identity profiles define how identities are processed and managed in SailPoint ISC.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the lifecycle state.",
				MarkdownDescription: "The human-readable display name of the lifecycle state, such as 'Active', 'Inactive', 'Terminated', etc.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "Creation date of the lifecycle state.",
				MarkdownDescription: "The timestamp when this lifecycle state was created in SailPoint ISC, in ISO 8601 format.",
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "Last modification date of the lifecycle state.",
				MarkdownDescription: "The timestamp when this lifecycle state was last modified in SailPoint ISC, in ISO 8601 format.",
			},
			"enabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates whether the lifecycle state is enabled or disabled.",
				MarkdownDescription: "Boolean flag indicating whether this lifecycle state is currently enabled (`true`) or disabled (`false`). Disabled lifecycle states are not active and won't be applied to identities.",
			},
			"technical_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The lifecycle state's technical name. This is for internal use.",
				MarkdownDescription: "The technical/internal name of the lifecycle state, used for system identification and API operations. This is typically a machine-readable identifier.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "Lifecycle state's description.",
				MarkdownDescription: "A detailed description of the lifecycle state, explaining its purpose and when it should be used in the identity lifecycle management process.",
			},
			"identity_count": schema.Int32Attribute{
				Computed:            true,
				Description:         "Number of identities that have the lifecycle state.",
				MarkdownDescription: "The current number of identities that are assigned to this lifecycle state. This provides insight into how many users are currently in this stage of their lifecycle.",
			},
			"access_profile_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				Description:         "List of unique access-profile IDs that are associated with the lifecycle state.",
				MarkdownDescription: "A list of access profile IDs that are automatically granted to identities when they enter this lifecycle state. Access profiles bundle together related entitlements and roles.",
			},
			"identity_state": schema.StringAttribute{
				Computed:            true,
				Description:         "The lifecycle state's associated identity state. This field is generally 'null'.",
				MarkdownDescription: "The identity state associated with this lifecycle state. This field controls the identity's overall status in the system and is typically `null` for most lifecycle states.",
			},
			"priority": schema.Int32Attribute{
				Computed:            true,
				Description:         "Priority level used to determine which profile to assign when a user exists in multiple profiles. Lower numeric values have higher priority. By default, new profiles are assigned the lowest priority. The assigned profile also controls access granted or removed during provisioning based on lifecycle state changes.",
				MarkdownDescription: "The priority level of this lifecycle state, used to determine precedence when an identity exists in multiple identity profiles. Lower numeric values indicate higher priority (e.g., priority `1` takes precedence over priority `10`). This affects which lifecycle state's access profiles and rules are applied during provisioning operations.",
			},
			"email_notification_option": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Options for email notifications when identities enter this lifecycle state.",
				MarkdownDescription: "Configuration options for sending email notifications when identities transition into this lifecycle state. This includes notifying managers, all admins, specific users, and providing a list of email addresses to notify.",
				Attributes: map[string]schema.Attribute{
					"notify_managers": schema.BoolAttribute{
						Computed:            true,
						Description:         "Whether to notify the identity's managers when they enter this lifecycle state.",
						MarkdownDescription: "Boolean flag indicating if the identity's direct managers should receive an email notification when the identity transitions into this lifecycle state.",
					},
					"notify_all_admins": schema.BoolAttribute{
						Computed:            true,
						Description:         "Whether to notify all administrators when identities enter this lifecycle state.",
						MarkdownDescription: "Boolean flag indicating if all system administrators should receive an email notification when any identity transitions into this lifecycle state.",
					},
					"notify_specific_users": schema.BoolAttribute{
						Computed:            true,
						Description:         "Whether to notify specific users when identities enter this lifecycle state.",
						MarkdownDescription: "Boolean flag indicating if a predefined list of specific users should receive an email notification when identities transition into this lifecycle state.",
					},
					"email_address_list": schema.ListAttribute{
						ElementType:         types.StringType,
						Computed:            true,
						Description:         "List of additional email addresses to notify when identities enter this lifecycle state.",
						MarkdownDescription: "A list of additional email addresses that should receive notifications when identities transition into this lifecycle state. This allows for custom notification recipients beyond managers and admins.",
					},
				},
			},
			"access_action_configuration": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Configuration for access actions to be performed during this lifecycle state transition.",
				MarkdownDescription: "Configuration that determines what access management actions should be taken when identities transition into this lifecycle state. This includes options for automatically removing access.",
				Attributes:          accessActionConfigurationAttributesForDataSource(),
			},
			"account_actions": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of account actions to be performed when identities enter this lifecycle state.",
				MarkdownDescription: "A list of account actions that define specific operations to be performed on identity accounts when they transition into this lifecycle state. Each action can enable, disable, or delete accounts on specific sources.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: accountActionAttributesForDataSource(),
				},
			},
		},
	}
}

// Helper functions for shared attribute definitions

func lifecycleStateAttributesForDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "System-generated unique ID of the lifecycle state.",
			MarkdownDescription: "The system-generated unique identifier for this lifecycle state. This is a UUID that uniquely identifies the lifecycle state within SailPoint ISC.",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			Description:         "Name of the lifecycle state.",
			MarkdownDescription: "The human-readable display name of the lifecycle state, such as 'Active', 'Inactive', 'Terminated', etc.",
		},
		"created": schema.StringAttribute{
			Computed:            true,
			Description:         "Creation date of the lifecycle state.",
			MarkdownDescription: "The timestamp when this lifecycle state was created in SailPoint ISC, in ISO 8601 format.",
		},
		"modified": schema.StringAttribute{
			Computed:            true,
			Description:         "Last modification date of the lifecycle state.",
			MarkdownDescription: "The timestamp when this lifecycle state was last modified in SailPoint ISC, in ISO 8601 format.",
		},
		"enabled": schema.BoolAttribute{
			Computed:            true,
			Description:         "Indicates whether the lifecycle state is enabled or disabled.",
			MarkdownDescription: "Boolean flag indicating whether this lifecycle state is currently enabled (`true`) or disabled (`false`). Disabled lifecycle states are not active and won't be applied to identities.",
		},
		"technical_name": schema.StringAttribute{
			Computed:            true,
			Description:         "The lifecycle state's technical name. This is for internal use.",
			MarkdownDescription: "The technical/internal name of the lifecycle state, used for system identification and API operations. This is typically a machine-readable identifier.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			Description:         "Lifecycle state's description.",
			MarkdownDescription: "A detailed description of the lifecycle state, explaining its purpose and when it should be used in the identity lifecycle management process.",
		},
		"identity_count": schema.Int32Attribute{
			Computed:            true,
			Description:         "Number of identities that have the lifecycle state.",
			MarkdownDescription: "The current number of identities that are assigned to this lifecycle state. This provides insight into how many users are currently in this stage of their lifecycle.",
		},
		"access_profile_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			Description:         "List of unique access-profile IDs that are associated with the lifecycle state.",
			MarkdownDescription: "A list of access profile IDs that are automatically granted to identities when they enter this lifecycle state. Access profiles bundle together related entitlements and roles.",
		},
		"identity_state": schema.StringAttribute{
			Computed:            true,
			Description:         "The lifecycle state's associated identity state. This field is generally 'null'.",
			MarkdownDescription: "The identity state associated with this lifecycle state. This field controls the identity's overall status in the system and is typically `null` for most lifecycle states.",
		},
		"priority": schema.Int32Attribute{
			Computed:            true,
			Description:         "Priority level used to determine which profile to assign when a user exists in multiple profiles. Lower numeric values have higher priority. By default, new profiles are assigned the lowest priority. The assigned profile also controls access granted or removed during provisioning based on lifecycle state changes.",
			MarkdownDescription: "The priority level of this lifecycle state, used to determine precedence when an identity exists in multiple identity profiles. Lower numeric values indicate higher priority (e.g., priority `1` takes precedence over priority `10`). This affects which lifecycle state's access profiles and rules are applied during provisioning operations.",
		},
		"email_notification_option": schema.SingleNestedAttribute{
			Computed:            true,
			Description:         "Options for email notifications when identities enter this lifecycle state.",
			MarkdownDescription: "Configuration options for sending email notifications when identities transition into this lifecycle state. This includes notifying managers, all admins, specific users, and providing a list of email addresses to notify.",
			Attributes: map[string]schema.Attribute{
				"notify_managers": schema.BoolAttribute{
					Computed:            true,
					Description:         "Whether to notify the identity's managers when they enter this lifecycle state.",
					MarkdownDescription: "Boolean flag indicating if the identity's direct managers should receive an email notification when the identity transitions into this lifecycle state.",
				},
				"notify_all_admins": schema.BoolAttribute{
					Computed:            true,
					Description:         "Whether to notify all administrators when identities enter this lifecycle state.",
					MarkdownDescription: "Boolean flag indicating if all system administrators should receive an email notification when any identity transitions into this lifecycle state.",
				},
				"notify_specific_users": schema.BoolAttribute{
					Computed:            true,
					Description:         "Whether to notify specific users when identities enter this lifecycle state.",
					MarkdownDescription: "Boolean flag indicating if a predefined list of specific users should receive an email notification when identities transition into this lifecycle state.",
				},
				"email_address_list": schema.ListAttribute{
					ElementType:         types.StringType,
					Computed:            true,
					Description:         "List of additional email addresses to notify when identities enter this lifecycle state.",
					MarkdownDescription: "A list of additional email addresses that should receive notifications when identities transition into this lifecycle state. This allows for custom notification recipients beyond managers and admins.",
				},
			},
		},
		"access_action_configuration": schema.SingleNestedAttribute{
			Computed:            true,
			Description:         "Configuration for access actions to be performed during this lifecycle state transition.",
			MarkdownDescription: "Configuration that determines what access management actions should be taken when identities transition into this lifecycle state. This includes options for automatically removing access.",
			Attributes:          accessActionConfigurationAttributesForDataSource(),
		},
		"account_actions": schema.ListNestedAttribute{
			Computed:            true,
			Description:         "List of account actions to be performed when identities enter this lifecycle state.",
			MarkdownDescription: "A list of account actions that define specific operations to be performed on identity accounts when they transition into this lifecycle state. Each action can enable, disable, or delete accounts on specific sources.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: accountActionAttributesForDataSource(),
			},
		},
	}
}

func accessActionConfigurationAttributesForDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"remove_all_access_enabled": schema.BoolAttribute{
			Computed:            true,
			Description:         "Whether to automatically remove all access when identities enter this lifecycle state.",
			MarkdownDescription: "Boolean flag indicating if all access profiles, entitlements, and roles should be automatically removed when identities transition into this lifecycle state. This is typically used for termination or suspension lifecycle states.",
		},
	}
}

func accountActionAttributesForDataSource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Computed:            true,
			Description:         "The action to be performed on the account.",
			MarkdownDescription: "The action to be performed on identity accounts. Values can be 'ENABLE', 'DISABLE', or 'DELETE'. This determines what happens to accounts on the specified sources when identities transition into this lifecycle state.",
		},
		"source_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			Description:         "List of source IDs to apply the action to.",
			MarkdownDescription: "A list of unique source identifiers where the account action is applied. The sources must have the ENABLE feature or be flat file sources.",
		},
		"exclude_source_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			Description:         "List of source IDs to exclude from the action.",
			MarkdownDescription: "A list of source identifiers that are excluded from the account action. This allows the action to be applied to most sources while excluding specific ones.",
		},
		"all_sources": schema.BoolAttribute{
			Computed:            true,
			Description:         "Whether to apply the action to all available sources.",
			MarkdownDescription: "Boolean flag indicating if the action is applied to all available sources. When true, source_ids is not provided. When false, source_ids specifies the target sources.",
		},
	}
}
