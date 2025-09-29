// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource Schema Function

func LifecycleStateResourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Lifecycle State resource schema",
		MarkdownDescription: "Lifecycle State resource schema",
		Attributes: map[string]schema.Attribute{
			"identity_profile_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the identity profile associated with the lifecycle state.",
				MarkdownDescription: "The unique identifier of the identity profile that contains this lifecycle state. Identity profiles define how identities are processed and managed in SailPoint ISC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the lifecycle state.",
				MarkdownDescription: "The unique identifier of the lifecycle state. This is a system-generated UUID that uniquely identifies the lifecycle state within SailPoint Identity Security Cloud (ISC).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the lifecycle state.",
				MarkdownDescription: "The human-readable display name of the lifecycle state, such as 'Active', 'Inactive', 'Terminated', etc.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"technical_name": schema.StringAttribute{
				Required:            true,
				Description:         "The lifecycle state's technical name. This is for internal use.",
				MarkdownDescription: "The technical/internal name of the lifecycle state, used for system identification and API operations. This is typically a machine-readable identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "Description of the lifecycle state.",
				MarkdownDescription: "A detailed description of the lifecycle state, providing context about its purpose and usage within identity management processes.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Indicates whether the lifecycle state is enabled or disabled.",
				MarkdownDescription: "Boolean flag indicating whether this lifecycle state is currently enabled (`true`) or disabled (`false`). Disabled lifecycle states are not active and won't be applied to identities.",
			},
			"identity_count": schema.Int32Attribute{
				Computed:            true,
				Description:         "Number of identities currently in this lifecycle state.",
				MarkdownDescription: "The count of identities that are currently assigned to this lifecycle state within the associated identity profile.",
			},
			"access_profile_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Description:         "List of access profile IDs associated with this lifecycle state.",
				MarkdownDescription: "A list of unique identifiers for access profiles that are linked to this lifecycle state. Access profiles define sets of entitlements and permissions granted to identities in this state.",
			},
			"identity_state": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE_SHORT_TERM", "INACTIVE_LONG_TERM"),
				},
				Description:         "The identity state associated with this lifecycle state.",
				MarkdownDescription: "The broader identity state (e.g., 'Active', 'Inactive', 'Terminated') that this lifecycle state maps to within the identity management framework.",
			},
			"priority": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The priority of the lifecycle state within the identity profile.",
				MarkdownDescription: "An integer value representing the priority of this lifecycle state relative to other states in the same identity profile. Lower numbers indicate higher priority.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "Creation date of the lifecycle state.",
				MarkdownDescription: "The timestamp when this lifecycle state was created in SailPoint ISC, in ISO 8601 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "Last modification date of the lifecycle state.",
				MarkdownDescription: "The timestamp when this lifecycle state was last modified in SailPoint ISC, in ISO 8601 format.",
			},
			"email_notification_option": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Options for email notifications when identities enter this lifecycle state.",
				MarkdownDescription: "Configuration options for sending email notifications when identities transition into this lifecycle state. This includes notifying managers, all admins, specific users, and providing a list of email addresses to notify.",
				Attributes:          emailNotificationOptionAttributesForResource(),
			},
			"access_action_configuration": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Configuration for access actions to be performed during this lifecycle state transition.",
				MarkdownDescription: "Configuration that determines what access management actions should be taken when identities transition into this lifecycle state. This includes options for automatically removing access.",
				Attributes:          accessActionConfigurationAttributesForResource(),
			},
			"account_actions": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "List of account actions to be performed when identities enter this lifecycle state.",
				MarkdownDescription: "A list of account actions that define specific operations to be performed on identity accounts when they transition into this lifecycle state. Each action can enable, disable, or delete accounts on specific sources.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: accountActionAttributesForResource(),
				},
			},
		},
	}
}

func emailNotificationOptionAttributesForResource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"notify_managers": schema.BoolAttribute{
			Optional:            true,
			Description:         "Whether to notify the identity's managers when they enter this lifecycle state.",
			MarkdownDescription: "Boolean flag indicating if the identity's direct managers should receive an email notification when the identity transitions into this lifecycle state.",
		},
		"notify_all_admins": schema.BoolAttribute{
			Optional:            true,
			Description:         "Whether to notify all administrators when identities enter this lifecycle state.",
			MarkdownDescription: "Boolean flag indicating if all system administrators should receive an email notification when any identity transitions into this lifecycle state.",
		},
		"notify_specific_users": schema.BoolAttribute{
			Optional:            true,
			Description:         "Whether to notify specific users when identities enter this lifecycle state.",
			MarkdownDescription: "Boolean flag indicating if a predefined list of specific users should receive an email notification when identities transition into this lifecycle state.",
		},
		"email_address_list": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Description:         "List of additional email addresses to notify when identities enter this lifecycle state.",
			MarkdownDescription: "A list of additional email addresses that should receive notifications when identities transition into this lifecycle state. This allows for custom notification recipients beyond managers and admins.",
		},
	}
}

func accessActionConfigurationAttributesForResource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"remove_all_access_enabled": schema.BoolAttribute{
			Optional:            true,
			Description:         "Whether to automatically remove all access when identities enter this lifecycle state.",
			MarkdownDescription: "Boolean flag indicating if all access profiles, entitlements, and roles should be automatically removed when identities transition into this lifecycle state. This is typically used for termination or suspension lifecycle states.",
		},
	}
}

func accountActionAttributesForResource() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"action": schema.StringAttribute{
			Required:            true,
			Description:         "The action to be performed on the account.",
			MarkdownDescription: "The action to be performed on identity accounts. Valid values are 'ENABLE', 'DISABLE', or 'DELETE'. This determines what happens to accounts on the specified sources when identities transition into this lifecycle state.",
			Validators: []validator.String{
				stringvalidator.OneOf("ENABLE", "DISABLE", "DELETE"),
			},
		},
		"source_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Description:         "List of source IDs to apply the action to.",
			MarkdownDescription: "A list of unique source identifiers where the account action should be applied. The sources must have the ENABLE feature or be flat file sources. Required if `all_sources` is not true. Cannot be used together with `exclude_source_ids`.",
		},
		"exclude_source_ids": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Description:         "List of source IDs to exclude from the action.",
			MarkdownDescription: "A list of source identifiers to exclude from the account action. This allows you to apply the action to most sources while excluding specific ones. Cannot be used together with `source_ids`.",
		},
		"all_sources": schema.BoolAttribute{
			Optional:            true,
			Description:         "Whether to apply the action to all available sources.",
			MarkdownDescription: "Boolean flag indicating if the action should be applied to all available sources. When true, `source_ids` must not be provided. When false or not set, `source_ids` is required. Default is false.",
		},
	}
}
