// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type IdentityAttributeSchemaBuilder struct{}

var (
	_ SchemaBuilder = &IdentityAttributeSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for IdentityAttribute resource.
func (sb *IdentityAttributeSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]resource_schema.Attribute{
		"name": resource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"display_name": resource_schema.StringAttribute{
			Description:         desc["display_name"].description,
			MarkdownDescription: desc["display_name"].markdown,
			Optional:            true,
		},
		"type": resource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Required:            true,
		},
		"system": resource_schema.BoolAttribute{
			Description:         desc["system"].description,
			MarkdownDescription: desc["system"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"standard": resource_schema.BoolAttribute{
			Description:         desc["standard"].description,
			MarkdownDescription: desc["standard"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"multi": resource_schema.BoolAttribute{
			Description:         desc["multi"].description,
			MarkdownDescription: desc["multi"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"searchable": resource_schema.BoolAttribute{
			Description:         desc["searchable"].description,
			MarkdownDescription: desc["searchable"].markdown,
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"sources": resource_schema.ListNestedAttribute{
			Description:         desc["sources"].description,
			MarkdownDescription: desc["sources"].markdown,
			Optional:            true,
			NestedObject: resource_schema.NestedAttributeObject{
				Attributes: map[string]resource_schema.Attribute{
					"type": resource_schema.StringAttribute{
						Description:         desc["sources.type"].description,
						MarkdownDescription: desc["sources.type"].markdown,
						Optional:            true,
					},
					"properties": resource_schema.StringAttribute{
						Description:         desc["sources.properties"].description,
						MarkdownDescription: desc["sources.properties"].markdown,
						Optional:            true,
					},
				},
			},
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for IdentityAttribute data source.
func (sb *IdentityAttributeSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]datasource_schema.Attribute{
		"name": datasource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Required:            true,
		},
		"display_name": datasource_schema.StringAttribute{
			Description:         desc["display_name"].description,
			MarkdownDescription: desc["display_name"].markdown,
			Computed:            true,
		},
		"type": datasource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Computed:            true,
		},
		"system": datasource_schema.BoolAttribute{
			Description:         desc["system"].description,
			MarkdownDescription: desc["system"].markdown,
			Computed:            true,
		},
		"standard": datasource_schema.BoolAttribute{
			Description:         desc["standard"].description,
			MarkdownDescription: desc["standard"].markdown,
			Computed:            true,
		},
		"multi": datasource_schema.BoolAttribute{
			Description:         desc["multi"].description,
			MarkdownDescription: desc["multi"].markdown,
			Computed:            true,
		},
		"searchable": datasource_schema.BoolAttribute{
			Description:         desc["searchable"].description,
			MarkdownDescription: desc["searchable"].markdown,
			Computed:            true,
		},
		"sources": datasource_schema.ListNestedAttribute{
			Description:         desc["sources"].description,
			MarkdownDescription: desc["sources"].markdown,
			Computed:            true,
			NestedObject: datasource_schema.NestedAttributeObject{
				Attributes: map[string]datasource_schema.Attribute{
					"type": datasource_schema.StringAttribute{
						Description:         desc["sources.type"].description,
						MarkdownDescription: desc["sources.type"].markdown,
						Computed:            true,
					},
					"properties": datasource_schema.StringAttribute{
						Description:         desc["sources.properties"].description,
						MarkdownDescription: desc["sources.properties"].markdown,
						Computed:            true,
					},
				},
			},
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *IdentityAttributeSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"name": {
			description: "The technical name of the identity attribute.",
			markdown:    "The technical name of the identity attribute. This is used as the identifier and is immutable after creation.",
		},
		"display_name": {
			description: "The display name of the identity attribute.",
			markdown:    "The user-friendly display name of the identity attribute shown in the UI.",
		},
		"type": {
			description: "The data type of the identity attribute.",
			markdown:    "The data type of the identity attribute (e.g., `string`, `int`, `date`).",
		},
		"system": {
			description: "Indicates whether this is a system-managed attribute.",
			markdown:    "Indicates whether this is a system-managed attribute. System attributes cannot be deleted. Must be set to `false` to make the attribute searchable or deletable. Defaults to `false`.",
		},
		"standard": {
			description: "Indicates whether this is a standard attribute.",
			markdown:    "Indicates whether this is a standard attribute. Standard attributes cannot be deleted. Must be set to `false` to make the attribute searchable or deletable. Defaults to `false`.",
		},
		"multi": {
			description: "Indicates whether this attribute can hold multiple values.",
			markdown:    "Indicates whether this attribute can hold multiple values. Must be set to `false` to make the attribute searchable. Defaults to `false`.",
		},
		"searchable": {
			description: "Indicates whether this attribute can be used in searches.",
			markdown:    "Indicates whether this attribute can be used in searches. Can only be `true` if `system`, `standard`, and `multi` are all `false`. Defaults to `false`.",
		},
		"sources": {
			description: "List of sources for this identity attribute.",
			markdown:    "List of sources that define how the identity attribute is populated.",
		},
		"sources.type": {
			description: "The type of the source.",
			markdown:    "The type of the source (e.g., `rule`, `static`).",
		},
		"sources.properties": {
			description: "Configuration properties for the source as a JSON string.",
			markdown:    "Configuration properties for the source as a JSON string. The structure varies by source type.",
		},
	}
}
