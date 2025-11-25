// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EntitlementSchemaBuilder struct{}

// GetDataSourceSchema implements SchemaBuilder for Entitlement data source.
func (sb *EntitlementSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
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
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"attribute": datasource_schema.StringAttribute{
			Description:         desc["attribute"].description,
			MarkdownDescription: desc["attribute"].markdown,
			Computed:            true,
		},
		"value": datasource_schema.StringAttribute{
			Description:         desc["value"].description,
			MarkdownDescription: desc["value"].markdown,
			Computed:            true,
		},
		"source_schema_object_type": datasource_schema.StringAttribute{
			Description:         desc["source_schema_object_type"].description,
			MarkdownDescription: desc["source_schema_object_type"].markdown,
			Computed:            true,
		},
		"privileged": datasource_schema.BoolAttribute{
			Description:         desc["privileged"].description,
			MarkdownDescription: desc["privileged"].markdown,
			Computed:            true,
		},
		"requestable": datasource_schema.BoolAttribute{
			Description:         desc["requestable"].description,
			MarkdownDescription: desc["requestable"].markdown,
			Computed:            true,
		},
		"cloud_governed": datasource_schema.BoolAttribute{
			Description:         desc["cloud_governed"].description,
			MarkdownDescription: desc["cloud_governed"].markdown,
			Computed:            true,
		},
		"source": datasource_schema.SingleNestedAttribute{
			Description:         desc["source"].description,
			MarkdownDescription: desc["source"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of the source (e.g., SOURCE).",
					Computed:    true,
				},
				"id": datasource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the source.",
					Computed:    true,
				},
				"name": datasource_schema.StringAttribute{
					Description: "The name of the source.",
					Computed:    true,
				},
			},
		},
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of the owner (e.g., IDENTITY).",
					Computed:    true,
				},
				"id": datasource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the owner.",
					Computed:    true,
				},
				"name": datasource_schema.StringAttribute{
					Description: "The name of the owner.",
					Computed:    true,
				},
			},
		},
		"attributes": datasource_schema.StringAttribute{
			Description:         desc["attributes"].description,
			MarkdownDescription: desc["attributes"].markdown,
			Computed:            true,
		},
		"direct_permissions": datasource_schema.ListNestedAttribute{
			Description:         desc["direct_permissions"].description,
			MarkdownDescription: desc["direct_permissions"].markdown,
			Computed:            true,
			NestedObject: datasource_schema.NestedAttributeObject{
				Attributes: map[string]datasource_schema.Attribute{
					"rights": datasource_schema.ListAttribute{
						Description: "Array of permission rights (e.g., SELECT, INSERT, UPDATE).",
						Computed:    true,
						ElementType: types.StringType,
					},
					"target": datasource_schema.StringAttribute{
						Description: "Target resource for the permission.",
						Computed:    true,
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
		"manually_updated_fields": datasource_schema.SingleNestedAttribute{
			Description:         desc["manually_updated_fields"].description,
			MarkdownDescription: desc["manually_updated_fields"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"display_name": datasource_schema.BoolAttribute{
					Description: "Indicates if the display name was manually updated.",
					Computed:    true,
				},
				"description": datasource_schema.BoolAttribute{
					Description: "Indicates if the description was manually updated.",
					Computed:    true,
				},
			},
		},
		"access_model_metadata": datasource_schema.SingleNestedAttribute{
			Description:         desc["access_model_metadata"].description,
			MarkdownDescription: desc["access_model_metadata"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"attributes": datasource_schema.ListNestedAttribute{
					Description: "Array of metadata attributes to classify the entitlement.",
					Computed:    true,
					NestedObject: datasource_schema.NestedAttributeObject{
						Attributes: map[string]datasource_schema.Attribute{
							"key": datasource_schema.StringAttribute{
								Description: "Unique identifier for the metadata type (e.g., 'iscCsp').",
								Computed:    true,
							},
							"name": datasource_schema.StringAttribute{
								Description: "Human readable name of the metadata type (e.g., 'CSP').",
								Computed:    true,
							},
							"multiselect": datasource_schema.BoolAttribute{
								Description: "Allows selecting multiple values (default: false).",
								Computed:    true,
							},
							"status": datasource_schema.StringAttribute{
								Description: "The state of the metadata item (e.g., 'active').",
								Computed:    true,
							},
							"type": datasource_schema.StringAttribute{
								Description: "The type of the metadata item (e.g., 'governance').",
								Computed:    true,
							},
							"object_types": datasource_schema.ListAttribute{
								Description: "The types of objects (e.g., ['general']).",
								Computed:    true,
								ElementType: types.StringType,
							},
							"description": datasource_schema.StringAttribute{
								Description: "Describes the metadata item.",
								Computed:    true,
							},
							"values": datasource_schema.ListNestedAttribute{
								Description: "The values to assign to the metadata item.",
								Computed:    true,
								NestedObject: datasource_schema.NestedAttributeObject{
									Attributes: map[string]datasource_schema.Attribute{
										"value": datasource_schema.StringAttribute{
											Description: "The value to assign to the metadata item (e.g., 'development').",
											Computed:    true,
										},
										"name": datasource_schema.StringAttribute{
											Description: "Display name of the value (e.g., 'Development').",
											Computed:    true,
										},
										"status": datasource_schema.StringAttribute{
											Description: "The status of the individual value (e.g., 'active').",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *EntitlementSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the entitlement.",
			markdown:    "Unique identifier (UUID) of the entitlement.",
		},
		"name": {
			description: "Name of the entitlement.",
			markdown:    "Display name of the entitlement.",
		},
		"created": {
			description: "ISO-8601 timestamp when the entitlement was created.",
			markdown:    "ISO-8601 timestamp when the entitlement was created (computed).",
		},
		"modified": {
			description: "ISO-8601 timestamp when the entitlement was last modified.",
			markdown:    "ISO-8601 timestamp when the entitlement was last modified (computed).",
		},
		"description": {
			description: "Description of the entitlement.",
			markdown:    "Description of the entitlement's purpose and scope.",
		},
		"attribute": {
			description: "The account schema attribute this entitlement represents.",
			markdown:    "The account schema attribute this entitlement represents (e.g., 'memberOf').",
		},
		"value": {
			description: "The actual value of the entitlement attribute.",
			markdown:    "The actual value of the entitlement attribute.",
		},
		"source_schema_object_type": {
			description: "The type of object in the source schema.",
			markdown:    "The type of object in the source schema (e.g., 'group').",
		},
		"privileged": {
			description: "Whether the entitlement is marked as privileged.",
			markdown:    "Boolean flag indicating if the entitlement is marked as privileged.",
		},
		"requestable": {
			description: "Whether the entitlement can be requested by users.",
			markdown:    "Boolean flag indicating if the entitlement can be requested by users through access requests.",
		},
		"cloud_governed": {
			description: "Whether the entitlement is cloud-governed.",
			markdown:    "Boolean flag indicating if the entitlement is cloud-governed.",
		},
		"source": {
			description: "The source system this entitlement comes from.",
			markdown:    "The source system this entitlement comes from. Contains `type` (e.g., 'SOURCE'), `id` (UUID), and `name`.",
		},
		"owner": {
			description: "Owner of the entitlement.",
			markdown:    "Owner of the entitlement. Contains `type` (e.g., 'IDENTITY'), `id` (UUID), and `name`.",
		},
		"attributes": {
			description: "Custom attributes as JSON.",
			markdown:    "Custom field name-value pairs as a JSON string.",
		},
		"direct_permissions": {
			description: "Direct permissions associated with this entitlement.",
			markdown:    "Array of direct permission objects with `rights` (array of permission rights) and `target` (target resource).",
		},
		"segments": {
			description: "Segments this entitlement belongs to.",
			markdown:    "Array of segment IDs that this entitlement belongs to.",
		},
		"manually_updated_fields": {
			description: "Fields that have been manually updated.",
			markdown:    "Object tracking which fields have been manually updated. Contains `display_name` and `description` boolean flags.",
		},
		"access_model_metadata": {
			description: "Access model metadata for this entitlement.",
			markdown:    "Object containing additional data to classify the entitlement. Contains an `attributes` array with metadata items including key, name, type, description, and values.",
		},
	}
}
