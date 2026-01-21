// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LifecycleStateSchemaBuilder struct{}

var (
	_ SchemaBuilder = &LifecycleStateSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for LifecycleState resource.
func (sb *LifecycleStateSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]resource_schema.Attribute{
		"id": resource_schema.StringAttribute{
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"identity_profile_id": resource_schema.StringAttribute{
			Description:         desc["identity_profile_id"].description,
			MarkdownDescription: desc["identity_profile_id"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": resource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"technical_name": resource_schema.StringAttribute{
			Description:         desc["technical_name"].description,
			MarkdownDescription: desc["technical_name"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
		},
		"enabled": resource_schema.BoolAttribute{
			Description:         desc["enabled"].description,
			MarkdownDescription: desc["enabled"].markdown,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"priority": resource_schema.Int64Attribute{
			Description:         desc["priority"].description,
			MarkdownDescription: desc["priority"].markdown,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"identity_state": resource_schema.StringAttribute{
			Description:         desc["identity_state"].description,
			MarkdownDescription: desc["identity_state"].markdown,
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"ACTIVE",
					"INACTIVE_SHORT_TERM",
					"INACTIVE_LONG_TERM",
				),
			},
		},
		"identity_count": resource_schema.Int64Attribute{
			Description:         desc["identity_count"].description,
			MarkdownDescription: desc["identity_count"].markdown,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"access_profile_ids": resource_schema.ListAttribute{
			Description:         desc["access_profile_ids"].description,
			MarkdownDescription: desc["access_profile_ids"].markdown,
			Optional:            true,
			ElementType:         types.StringType,
		},
		"email_notification_option": resource_schema.SingleNestedAttribute{
			Description:         desc["email_notification_option"].description,
			MarkdownDescription: desc["email_notification_option"].markdown,
			Optional:            true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]resource_schema.Attribute{
				"notify_managers": resource_schema.BoolAttribute{
					Description: "Whether to notify the identity's managers when entering this lifecycle state.",
					Optional:    true,
					// Computed:    true,
				},
				"notify_all_admins": resource_schema.BoolAttribute{
					Description: "Whether to notify all administrators when entering this lifecycle state.",
					Optional:    true,
					// Computed:    true,
				},
				"notify_specific_users": resource_schema.BoolAttribute{
					Description: "Whether to notify specific users when entering this lifecycle state.",
					Optional:    true,
					// Computed:    true,
				},
				"email_address_list": resource_schema.ListAttribute{
					Description: "List of email addresses to notify when entering this lifecycle state (requires notify_specific_users to be true).",
					Optional:    true,
					// Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		"account_actions": resource_schema.ListNestedAttribute{
			Description:         desc["account_actions"].description,
			MarkdownDescription: desc["account_actions"].markdown,
			Optional:            true,
			NestedObject: resource_schema.NestedAttributeObject{
				Attributes: map[string]resource_schema.Attribute{
					"action": resource_schema.StringAttribute{
						Description: "The action to perform on accounts (ENABLE, DISABLE, or DELETE).",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ENABLE",
								"DISABLE",
								"DELETE",
							),
						},
					},
					"source_ids": resource_schema.ListAttribute{
						Description: "List of source IDs on which to perform the action. Mutually exclusive with exclude_source_ids and all_sources.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"exclude_source_ids": resource_schema.ListAttribute{
						Description: "List of source IDs to exclude from the action. Requires all_sources to be true.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"all_sources": resource_schema.BoolAttribute{
						Description: "Whether to perform the action on all sources. If true, source_ids should not be specified.",
						Optional:    true,
					},
				},
			},
		},
		"access_action_configuration": resource_schema.SingleNestedAttribute{
			Description:         desc["access_action_configuration"].description,
			MarkdownDescription: desc["access_action_configuration"].markdown,
			Optional:            true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]resource_schema.Attribute{
				"remove_all_access_enabled": resource_schema.BoolAttribute{
					Description: "Whether to remove all access when entering this lifecycle state.",
					Optional:    true,
				},
			},
		},
		"created": resource_schema.StringAttribute{
			Description:         desc["created"].description,
			MarkdownDescription: desc["created"].markdown,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"modified": resource_schema.StringAttribute{
			Description:         desc["modified"].description,
			MarkdownDescription: desc["modified"].markdown,
			Computed:            true,
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for LifecycleState data source.
func (sb *LifecycleStateSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]datasource_schema.Attribute{
		"id": datasource_schema.StringAttribute{
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
			Required:            true,
		},
		"identity_profile_id": datasource_schema.StringAttribute{
			Description:         desc["identity_profile_id"].description,
			MarkdownDescription: desc["identity_profile_id"].markdown,
			Required:            true,
		},
		"name": datasource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Computed:            true,
		},
		"technical_name": datasource_schema.StringAttribute{
			Description:         desc["technical_name"].description,
			MarkdownDescription: desc["technical_name"].markdown,
			Computed:            true,
		},
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"enabled": datasource_schema.BoolAttribute{
			Description:         desc["enabled"].description,
			MarkdownDescription: desc["enabled"].markdown,
			Computed:            true,
		},
		"priority": datasource_schema.Int64Attribute{
			Description:         desc["priority"].description,
			MarkdownDescription: desc["priority"].markdown,
			Computed:            true,
		},
		"identity_state": datasource_schema.StringAttribute{
			Description:         desc["identity_state"].description,
			MarkdownDescription: desc["identity_state"].markdown,
			Computed:            true,
		},
		"identity_count": datasource_schema.Int64Attribute{
			Description:         desc["identity_count"].description,
			MarkdownDescription: desc["identity_count"].markdown,
			Computed:            true,
		},
		"access_profile_ids": datasource_schema.ListAttribute{
			Description:         desc["access_profile_ids"].description,
			MarkdownDescription: desc["access_profile_ids"].markdown,
			Computed:            true,
			ElementType:         types.StringType,
		},
		"email_notification_option": datasource_schema.SingleNestedAttribute{
			Description:         desc["email_notification_option"].description,
			MarkdownDescription: desc["email_notification_option"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"notify_managers": datasource_schema.BoolAttribute{
					Description: "Whether to notify the identity's managers when entering this lifecycle state.",
					Computed:    true,
				},
				"notify_all_admins": datasource_schema.BoolAttribute{
					Description: "Whether to notify all administrators when entering this lifecycle state.",
					Computed:    true,
				},
				"notify_specific_users": datasource_schema.BoolAttribute{
					Description: "Whether to notify specific users when entering this lifecycle state.",
					Computed:    true,
				},
				"email_address_list": datasource_schema.ListAttribute{
					Description: "List of email addresses to notify when entering this lifecycle state (requires notify_specific_users to be true).",
					Computed:    true,
					ElementType: types.StringType,
				},
			},
		},
		"account_actions": datasource_schema.ListNestedAttribute{
			Description:         desc["account_actions"].description,
			MarkdownDescription: desc["account_actions"].markdown,
			Computed:            true,
			NestedObject: datasource_schema.NestedAttributeObject{
				Attributes: map[string]datasource_schema.Attribute{
					"action": datasource_schema.StringAttribute{
						Description: "The action to perform on accounts (ENABLE, DISABLE, or DELETE).",
						Computed:    true,
					},
					"source_ids": datasource_schema.ListAttribute{
						Description: "List of source IDs on which to perform the action. Mutually exclusive with exclude_source_ids and all_sources.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"exclude_source_ids": datasource_schema.ListAttribute{
						Description: "List of source IDs to exclude from the action. Requires all_sources to be true.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"all_sources": datasource_schema.BoolAttribute{
						Description: "Whether to perform the action on all sources. If true, source_ids should not be specified.",
						Computed:    true,
					},
				},
			},
		},
		"access_action_configuration": datasource_schema.SingleNestedAttribute{
			Description:         desc["access_action_configuration"].description,
			MarkdownDescription: desc["access_action_configuration"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"remove_all_access_enabled": datasource_schema.BoolAttribute{
					Description: "Whether to remove all access when entering this lifecycle state.",
					Computed:    true,
				},
			},
		},
		"created": datasource_schema.StringAttribute{
			Description:         desc["created"].description,
			MarkdownDescription: desc["created"].markdown,
			Computed:            true,
		},
		"modified": datasource_schema.StringAttribute{
			Description:         desc["modified"].description,
			MarkdownDescription: desc["modified"].markdown,
			Computed:            true,
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *LifecycleStateSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the lifecycle state.",
			markdown:    "Unique identifier (UUID) of the lifecycle state.",
		},
		"identity_profile_id": {
			description: "ID of the parent identity profile.",
			markdown:    "The unique identifier (UUID) of the identity profile that this lifecycle state belongs to. **Changing this forces recreation of the resource.**",
		},
		"name": {
			description: "Display name of the lifecycle state.",
			markdown:    "The display name of the lifecycle state as it appears in the UI.",
		},
		"technical_name": {
			description: "Technical name used for referencing the lifecycle state.",
			markdown:    "The technical name used for referencing the lifecycle state programmatically. Must be unique within the identity profile. **Changing this forces recreation of the resource.**",
		},
		"description": {
			description: "Description of the lifecycle state.",
			markdown:    "Optional description explaining the purpose and behavior of this lifecycle state.",
		},
		"enabled": {
			description: "Whether the lifecycle state is enabled.",
			markdown:    "Whether the lifecycle state is enabled (true) or disabled (false). When disabled, identities cannot enter this lifecycle state. Defaults to `false` if not specified.",
		},
		"priority": {
			description: "Priority order for the lifecycle state.",
			markdown:    "Priority order for the lifecycle state. Lower numbers have higher priority. This determines the order in which lifecycle states are evaluated.",
		},
		"identity_state": {
			description: "Identity state classification.",
			markdown:    "The identity state classification. Valid values are `ACTIVE`, `INACTIVE_SHORT_TERM`, or `INACTIVE_LONG_TERM`. This determines how the identity's access is managed.",
		},
		"identity_count": {
			description: "Number of identities currently in this lifecycle state.",
			markdown:    "The number of identities currently in this lifecycle state (computed, read-only).",
		},
		"access_profile_ids": {
			description: "List of access profile IDs to grant when entering this lifecycle state.",
			markdown:    "List of access profile IDs (UUIDs) to automatically grant to identities when they enter this lifecycle state.",
		},
		"email_notification_option": {
			description: "Email notification configuration for this lifecycle state.",
			markdown:    "Configuration for email notifications to send when an identity enters this lifecycle state. Allows notifying managers, administrators, or specific email addresses.",
		},
		"account_actions": {
			description: "Account actions to perform when entering this lifecycle state.",
			markdown:    "List of account actions to perform on identity accounts when entering this lifecycle state. Actions include ENABLE, DISABLE, or DELETE, and can target specific sources or all sources.",
		},
		"access_action_configuration": {
			description: "Access action configuration for this lifecycle state.",
			markdown:    "Configuration for access-related actions to perform when entering this lifecycle state, such as removing all access.",
		},
		"created": {
			description: "ISO-8601 timestamp when the lifecycle state was created.",
			markdown:    "ISO-8601 timestamp when the lifecycle state was created (computed, read-only).",
		},
		"modified": {
			description: "ISO-8601 timestamp when the lifecycle state was last modified.",
			markdown:    "ISO-8601 timestamp when the lifecycle state was last modified (computed, read-only).",
		},
	}
}
