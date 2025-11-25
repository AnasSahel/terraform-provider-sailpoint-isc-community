// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type LauncherSchemaBuilder struct{}

var (
	_ SchemaBuilder = &LauncherSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for Launcher resource.
func (sb *LauncherSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
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
			Required:            true,
		},
		"type": resource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Required:            true,
		},
		"disabled": resource_schema.BoolAttribute{
			Description:         desc["disabled"].description,
			MarkdownDescription: desc["disabled"].markdown,
			Required:            true,
		},
		"reference": resource_schema.SingleNestedAttribute{
			Description:         desc["reference"].description,
			MarkdownDescription: desc["reference"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description: "The type of the referenced resource (typically 'WORKFLOW').",
					Required:    true,
				},
				"id": resource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the referenced resource.",
					Required:    true,
				},
			},
		},
		"config": resource_schema.StringAttribute{
			Description:         desc["config"].description,
			MarkdownDescription: desc["config"].markdown,
			Required:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description: "The type of the referenced object (e.g., IDENTITY).",
					Computed:    true,
				},
				"id": resource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the owner identity.",
					Computed:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "The name of the owner identity.",
					Computed:    true,
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

// GetDataSourceSchema implements SchemaBuilder for Launcher data source.
func (sb *LauncherSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
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
		"type": datasource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Computed:            true,
		},
		"disabled": datasource_schema.BoolAttribute{
			Description:         desc["disabled"].description,
			MarkdownDescription: desc["disabled"].markdown,
			Computed:            true,
		},
		"reference": datasource_schema.SingleNestedAttribute{
			Description:         desc["reference"].description,
			MarkdownDescription: desc["reference"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of the referenced resource (typically 'WORKFLOW').",
					Computed:    true,
				},
				"id": datasource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the referenced resource.",
					Computed:    true,
				},
			},
		},
		"config": datasource_schema.StringAttribute{
			Description:         desc["config"].description,
			MarkdownDescription: desc["config"].markdown,
			Computed:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of the referenced object (e.g., IDENTITY).",
					Computed:    true,
				},
				"id": datasource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the owner identity.",
					Computed:    true,
				},
				"name": datasource_schema.StringAttribute{
					Description: "The name of the owner identity.",
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
func (sb *LauncherSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the launcher.",
			markdown:    "Unique identifier (UUID) of the launcher.",
		},
		"name": {
			description: "Name of the launcher.",
			markdown:    "Name of the launcher as it appears in the UI. Limited to 255 characters.",
		},
		"description": {
			description: "Description of the launcher.",
			markdown:    "Description of the launcher's purpose. Limited to 2000 characters.",
		},
		"type": {
			description: "Type of the launcher.",
			markdown:    "Type of the launcher. Currently only 'INTERACTIVE_PROCESS' is supported.",
		},
		"disabled": {
			description: "Whether the launcher is disabled.",
			markdown:    "Whether the launcher is disabled (true) or enabled (false).",
		},
		"reference": {
			description: "Reference to the workflow or resource.",
			markdown:    "Reference to the workflow or other resource that this launcher will execute. Must include `type` (typically 'WORKFLOW') and `id` (UUID of the referenced resource).",
		},
		"config": {
			description: "JSON configuration associated with the launcher.",
			markdown:    "JSON configuration associated with this launcher. Maximum size of 4KB.",
		},
		"owner": {
			description: "Owner of the launcher.",
			markdown:    "Owner of the launcher (computed). Contains `type` (typically 'IDENTITY'), `id` (UUID), and optionally `name`.",
		},
		"created": {
			description: "ISO-8601 timestamp when the launcher was created.",
			markdown:    "ISO-8601 timestamp when the launcher was created (computed).",
		},
		"modified": {
			description: "ISO-8601 timestamp when the launcher was last modified.",
			markdown:    "ISO-8601 timestamp when the launcher was last modified (computed).",
		},
	}
}
