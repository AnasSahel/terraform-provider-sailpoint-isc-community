// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AccessProfileSchemaBuilder struct{}

var (
	_ SchemaBuilder = &AccessProfileSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for AccessProfile resource.
func (sb *AccessProfileSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
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
		"name": resource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Required:            true,
		},
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
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
		"enabled": resource_schema.BoolAttribute{
			Description:         desc["enabled"].description,
			MarkdownDescription: desc["enabled"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"requestable": resource_schema.BoolAttribute{
			Description:         desc["requestable"].description,
			MarkdownDescription: desc["requestable"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description:         "The type of the referenced object (IDENTITY).",
					MarkdownDescription: "The type of the referenced object (IDENTITY).",
					Required:            true,
				},
				"id": resource_schema.StringAttribute{
					Description:         "The unique identifier of the identity.",
					MarkdownDescription: "The unique identifier (UUID) of the identity.",
					Required:            true,
				},
				"name": resource_schema.StringAttribute{
					Description:         "The name of the identity.",
					MarkdownDescription: "The name of the identity.",
					Computed:            true,
				},
			},
		},
		"source": resource_schema.SingleNestedAttribute{
			Description:         desc["source"].description,
			MarkdownDescription: desc["source"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description:         "The type of the referenced object (SOURCE).",
					MarkdownDescription: "The type of the referenced object (SOURCE).",
					Required:            true,
				},
				"id": resource_schema.StringAttribute{
					Description:         "The unique identifier of the source.",
					MarkdownDescription: "The unique identifier (UUID) of the source.",
					Required:            true,
				},
				"name": resource_schema.StringAttribute{
					Description:         "The name of the source.",
					MarkdownDescription: "The name of the source.",
					Computed:            true,
				},
			},
		},
		"entitlements": resource_schema.ListNestedAttribute{
			Description:         desc["entitlements"].description,
			MarkdownDescription: desc["entitlements"].markdown,
			Optional:            true,
			NestedObject: resource_schema.NestedAttributeObject{
				Attributes: map[string]resource_schema.Attribute{
					"type": resource_schema.StringAttribute{
						Description:         "The type of the referenced object (ENTITLEMENT).",
						MarkdownDescription: "The type of the referenced object (ENTITLEMENT).",
						Required:            true,
					},
					"id": resource_schema.StringAttribute{
						Description:         "The unique identifier of the entitlement.",
						MarkdownDescription: "The unique identifier (UUID) of the entitlement.",
						Required:            true,
					},
					"name": resource_schema.StringAttribute{
						Description:         "The name of the entitlement.",
						MarkdownDescription: "The name of the entitlement.",
						Computed:            true,
					},
				},
			},
		},
		"segments": resource_schema.ListAttribute{
			Description:         desc["segments"].description,
			MarkdownDescription: desc["segments"].markdown,
			Optional:            true,
			ElementType:         types.StringType,
		},
		"access_request_config": resource_schema.StringAttribute{
			Description:         desc["access_request_config"].description,
			MarkdownDescription: desc["access_request_config"].markdown,
			Optional:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"revoke_request_config": resource_schema.StringAttribute{
			Description:         desc["revoke_request_config"].description,
			MarkdownDescription: desc["revoke_request_config"].markdown,
			Optional:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"provisioning_criteria": resource_schema.StringAttribute{
			Description:         desc["provisioning_criteria"].description,
			MarkdownDescription: desc["provisioning_criteria"].markdown,
			Optional:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for AccessProfile data source.
func (sb *AccessProfileSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]datasource_schema.Attribute{
		"id": datasource_schema.StringAttribute{
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
			Required:            true,
		},
		"name": datasource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Computed:            true,
		},
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
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
		"enabled": datasource_schema.BoolAttribute{
			Description:         desc["enabled"].description,
			MarkdownDescription: desc["enabled"].markdown,
			Computed:            true,
		},
		"requestable": datasource_schema.BoolAttribute{
			Description:         desc["requestable"].description,
			MarkdownDescription: desc["requestable"].markdown,
			Computed:            true,
		},
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description:         "The type of the referenced object (IDENTITY).",
					MarkdownDescription: "The type of the referenced object (IDENTITY).",
					Computed:            true,
				},
				"id": datasource_schema.StringAttribute{
					Description:         "The unique identifier of the identity.",
					MarkdownDescription: "The unique identifier (UUID) of the identity.",
					Computed:            true,
				},
				"name": datasource_schema.StringAttribute{
					Description:         "The name of the identity.",
					MarkdownDescription: "The name of the identity.",
					Computed:            true,
				},
			},
		},
		"source": datasource_schema.SingleNestedAttribute{
			Description:         desc["source"].description,
			MarkdownDescription: desc["source"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description:         "The type of the referenced object (SOURCE).",
					MarkdownDescription: "The type of the referenced object (SOURCE).",
					Computed:            true,
				},
				"id": datasource_schema.StringAttribute{
					Description:         "The unique identifier of the source.",
					MarkdownDescription: "The unique identifier (UUID) of the source.",
					Computed:            true,
				},
				"name": datasource_schema.StringAttribute{
					Description:         "The name of the source.",
					MarkdownDescription: "The name of the source.",
					Computed:            true,
				},
			},
		},
		"entitlements": datasource_schema.ListNestedAttribute{
			Description:         desc["entitlements"].description,
			MarkdownDescription: desc["entitlements"].markdown,
			Computed:            true,
			NestedObject: datasource_schema.NestedAttributeObject{
				Attributes: map[string]datasource_schema.Attribute{
					"type": datasource_schema.StringAttribute{
						Description:         "The type of the referenced object (ENTITLEMENT).",
						MarkdownDescription: "The type of the referenced object (ENTITLEMENT).",
						Computed:            true,
					},
					"id": datasource_schema.StringAttribute{
						Description:         "The unique identifier of the entitlement.",
						MarkdownDescription: "The unique identifier (UUID) of the entitlement.",
						Computed:            true,
					},
					"name": datasource_schema.StringAttribute{
						Description:         "The name of the entitlement.",
						MarkdownDescription: "The name of the entitlement.",
						Computed:            true,
					},
				},
			},
		},
		"segments": datasource_schema.ListAttribute{
			Description:         desc["segments"].description,
			MarkdownDescription: desc["segments"].markdown,
			Computed:            true,
			ElementType:         types.StringType,
		},
		"access_request_config": datasource_schema.StringAttribute{
			Description:         desc["access_request_config"].description,
			MarkdownDescription: desc["access_request_config"].markdown,
			Computed:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"revoke_request_config": datasource_schema.StringAttribute{
			Description:         desc["revoke_request_config"].description,
			MarkdownDescription: desc["revoke_request_config"].markdown,
			Computed:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"provisioning_criteria": datasource_schema.StringAttribute{
			Description:         desc["provisioning_criteria"].description,
			MarkdownDescription: desc["provisioning_criteria"].markdown,
			Computed:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *AccessProfileSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the access profile.",
			markdown:    "Unique identifier (UUID) of the access profile.",
		},
		"name": {
			description: "Name of the access profile.",
			markdown:    "Name of the access profile as it appears in the UI.",
		},
		"description": {
			description: "Description of the access profile.",
			markdown:    "Description of the access profile (maximum 2000 characters).",
		},
		"created": {
			description: "Timestamp when the access profile was created.",
			markdown:    "ISO-8601 timestamp when the access profile was created.",
		},
		"modified": {
			description: "Timestamp when the access profile was last modified.",
			markdown:    "ISO-8601 timestamp when the access profile was last modified.",
		},
		"enabled": {
			description: "Whether the access profile is enabled.",
			markdown:    "Whether the access profile is enabled. Defaults to true.",
		},
		"requestable": {
			description: "Whether the access profile can be requested by users.",
			markdown:    "Whether users can request this access profile. Defaults to true.",
		},
		"owner": {
			description: "Reference to the identity that owns this access profile.",
			markdown:    "Reference to the identity that owns this access profile. The user must have ROLE_SUBADMIN or SOURCE_SUBADMIN authority.",
		},
		"source": {
			description: "Reference to the source that this access profile is attached to.",
			markdown:    "Reference to the source that this access profile is attached to. The source determines which entitlements are available.",
		},
		"entitlements": {
			description: "List of entitlements included in this access profile.",
			markdown:    "List of entitlement references included in this access profile. Entitlements must exist on the access profile's source. Use the [list entitlements endpoint](https://developer.sailpoint.com/docs/api/v2025/list-entitlements) with filters to find available entitlements.",
		},
		"segments": {
			description: "List of segment IDs associated with this access profile.",
			markdown:    "List of segment identifiers (UUIDs) associated with this access profile for governance segmentation.",
		},
		"access_request_config": {
			description: "Access request approval configuration as a JSON string.",
			markdown:    "Configuration for access request approval workflows as a JSON string. Defines how requests for this access profile are approved.",
		},
		"revoke_request_config": {
			description: "Revoke request approval configuration as a JSON string.",
			markdown:    "Configuration for revocation approval workflows as a JSON string. Defines how revocations of this access profile are processed.",
		},
		"provisioning_criteria": {
			description: "Provisioning criteria configuration as a JSON string.",
			markdown:    "Provisioning criteria for multi-account selection as a JSON string. Defines logic for selecting which account to provision when multiple exist.",
		},
	}
}
