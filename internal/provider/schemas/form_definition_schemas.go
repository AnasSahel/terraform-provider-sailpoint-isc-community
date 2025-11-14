// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type FormDefinitionSchemaBuilder struct{}

var (
	_ SchemaBuilder = &FormDefinitionSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for FormDefinition resource.
func (sb *FormDefinitionSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
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
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description: "The type of the referenced object.",
					Required:    true,
				},
				"id": resource_schema.StringAttribute{
					Description: "The unique identifier of the referenced object.",
					Required:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "The name of the referenced object.",
					Optional:    true,
				},
			},
		},
		"used_by": resource_schema.ListNestedAttribute{
			Description:         desc["used_by"].description,
			MarkdownDescription: desc["used_by"].markdown,
			Computed:            true,
			NestedObject: resource_schema.NestedAttributeObject{
				Attributes: map[string]resource_schema.Attribute{
					"type": resource_schema.StringAttribute{
						Description: "The type of the referenced object.",
						Computed:    true,
					},
					"id": resource_schema.StringAttribute{
						Description: "The unique identifier of the referenced object.",
						Computed:    true,
					},
					"name": resource_schema.StringAttribute{
						Description: "The name of the referenced object.",
						Computed:    true,
					},
				},
			},
		},
		"form_input": resource_schema.StringAttribute{
			Description:         desc["form_input"].description,
			MarkdownDescription: desc["form_input"].markdown,
			Optional:            true,
		},
		"form_elements": resource_schema.StringAttribute{
			Description:         desc["form_elements"].description,
			MarkdownDescription: desc["form_elements"].markdown,
			Required:            true,
		},
		"form_conditions": resource_schema.StringAttribute{
			Description:         desc["form_conditions"].description,
			MarkdownDescription: desc["form_conditions"].markdown,
			Optional:            true,
		},
		"created": resource_schema.StringAttribute{
			Description:         desc["created"].description,
			MarkdownDescription: desc["created"].markdown,
			Computed:            true,
		},
		"modified": resource_schema.StringAttribute{
			Description:         desc["modified"].description,
			MarkdownDescription: desc["modified"].markdown,
			Computed:            true,
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for FormDefinition data source.
func (sb *FormDefinitionSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
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
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of the referenced object.",
					Computed:    true,
				},
				"id": datasource_schema.StringAttribute{
					Description: "The unique identifier of the referenced object.",
					Computed:    true,
				},
				"name": datasource_schema.StringAttribute{
					Description: "The name of the referenced object.",
					Computed:    true,
				},
			},
		},
		"used_by": datasource_schema.ListNestedAttribute{
			Description:         desc["used_by"].description,
			MarkdownDescription: desc["used_by"].markdown,
			Computed:            true,
			NestedObject: datasource_schema.NestedAttributeObject{
				Attributes: map[string]datasource_schema.Attribute{
					"type": datasource_schema.StringAttribute{
						Description: "The type of the referenced object.",
						Computed:    true,
					},
					"id": datasource_schema.StringAttribute{
						Description: "The unique identifier of the referenced object.",
						Computed:    true,
					},
					"name": datasource_schema.StringAttribute{
						Description: "The name of the referenced object.",
						Computed:    true,
					},
				},
			},
		},
		"form_input": datasource_schema.StringAttribute{
			Description:         desc["form_input"].description,
			MarkdownDescription: desc["form_input"].markdown,
			Computed:            true,
		},
		"form_elements": datasource_schema.StringAttribute{
			Description:         desc["form_elements"].description,
			MarkdownDescription: desc["form_elements"].markdown,
			Computed:            true,
		},
		"form_conditions": datasource_schema.StringAttribute{
			Description:         desc["form_conditions"].description,
			MarkdownDescription: desc["form_conditions"].markdown,
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
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *FormDefinitionSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the form definition.",
			markdown:    "Unique identifier (UUID) of the form definition.",
		},
		"name": {
			description: "Name of the form.",
			markdown:    "Name of the form as it appears in the UI.",
		},
		"description": {
			description: "Description of the form.",
			markdown:    "Description text that explains the purpose of this form.",
		},
		"owner": {
			description: "Owner of the form definition. Required - must specify type and id.",
			markdown:    "**Required.** Owner reference containing the identity who owns this form. Must include type (e.g., 'IDENTITY') and id fields.",
		},
		"used_by": {
			description: "List of objects using this form definition.",
			markdown:    "List of object references showing which systems are using this form definition. Automatically tracked by the API when systems use the form. Each reference includes type, id, and name.",
		},
		"form_input": {
			description: "Form input configuration as a JSON string.",
			markdown:    "Form input configuration defining the data sources and inputs for the form, represented as a JSON string.",
		},
		"form_elements": {
			description: "Form elements configuration as a JSON string. Required - forms must have at least one section with fields.",
			markdown:    "**Required.** Form elements configuration defining sections and fields for data collection, represented as a JSON string. Forms are composed of sections that split the form into logical groups, and fields that are the data collection points. At minimum, a form must contain one section with at least one field.",
		},
		"form_conditions": {
			description: "Form conditions configuration as a JSON string.",
			markdown:    "Form conditions configuration defining conditional logic that modifies the form dynamically, represented as a JSON string.",
		},
		"created": {
			description: "Timestamp when the form was created.",
			markdown:    "ISO 8601 timestamp indicating when the form definition was created.",
		},
		"modified": {
			description: "Timestamp when the form was last modified.",
			markdown:    "ISO 8601 timestamp indicating when the form definition was last modified.",
		},
	}
}
