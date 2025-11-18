// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type WorkflowSchemaBuilder struct{}

var (
	_ SchemaBuilder = &WorkflowSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for Workflow resource.
func (sb *WorkflowSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
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
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description: "The type of the referenced object (e.g., IDENTITY).",
					Required:    true,
				},
				"id": resource_schema.StringAttribute{
					Description: "The unique identifier (UUID) of the owner identity.",
					Required:    true,
				},
				"name": resource_schema.StringAttribute{
					Description: "The name of the owner identity.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
		},
		"definition": resource_schema.SingleNestedAttribute{
			Description:         desc["definition"].description,
			MarkdownDescription: desc["definition"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"start": resource_schema.StringAttribute{
					Description: "The name of the first step to execute in the workflow.",
					Required:    true,
				},
				"steps": resource_schema.StringAttribute{
					Description: "Workflow steps as a JSON string. Each step defines an action or operator with its configuration.",
					Required:    true,
					CustomType:  jsontypes.NormalizedType{},
				},
			},
		},
		"trigger": resource_schema.SingleNestedAttribute{
			Description:         desc["trigger"].description,
			MarkdownDescription: desc["trigger"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description: "The type of trigger (e.g., EVENT, SCHEDULED, REQUEST_RESPONSE).",
					Required:    true,
				},
				"display_name": resource_schema.StringAttribute{
					Description: "Display name for the trigger.",
					Optional:    true,
				},
				"attributes": resource_schema.StringAttribute{
					Description: "Trigger-specific attributes as a JSON string. Structure varies by trigger type.",
					Optional:    true,
					CustomType:  jsontypes.NormalizedType{},
				},
			},
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

// GetDataSourceSchema implements SchemaBuilder for Workflow data source.
func (sb *WorkflowSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
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
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"definition": datasource_schema.SingleNestedAttribute{
			Description:         desc["definition"].description,
			MarkdownDescription: desc["definition"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"start": datasource_schema.StringAttribute{
					Description: "The name of the first step to execute in the workflow.",
					Computed:    true,
				},
				"steps": datasource_schema.StringAttribute{
					Description: "Workflow steps as a JSON string. Each step defines an action or operator with its configuration.",
					Computed:    true,
					CustomType:  jsontypes.NormalizedType{},
				},
			},
		},
		"trigger": datasource_schema.SingleNestedAttribute{
			Description:         desc["trigger"].description,
			MarkdownDescription: desc["trigger"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description: "The type of trigger (e.g., EVENT, SCHEDULED, REQUEST_RESPONSE).",
					Computed:    true,
				},
				"display_name": datasource_schema.StringAttribute{
					Description: "Display name for the trigger.",
					Computed:    true,
				},
				"attributes": datasource_schema.StringAttribute{
					Description: "Trigger-specific attributes as a JSON string. Structure varies by trigger type.",
					Computed:    true,
					CustomType:  jsontypes.NormalizedType{},
				},
			},
		},
		"enabled": datasource_schema.BoolAttribute{
			Description:         desc["enabled"].description,
			MarkdownDescription: desc["enabled"].markdown,
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
func (sb *WorkflowSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the workflow.",
			markdown:    "Unique identifier (UUID) of the workflow.",
		},
		"name": {
			description: "Name of the workflow.",
			markdown:    "Name of the workflow as it appears in the UI.",
		},
		"owner": {
			description: "Owner of the workflow.",
			markdown:    "Owner of the workflow. Must be a valid identity reference with `type` (typically 'IDENTITY'), `id` (UUID), and optionally `name`.",
		},
		"description": {
			description: "Description of the workflow.",
			markdown:    "Description of the workflow's purpose and functionality.",
		},
		"definition": {
			description: "Workflow definition.",
			markdown:    "Workflow definition containing the workflow logic. Must include `start` (name of first step to execute) and `steps` (JSON string with all workflow steps and their configurations).",
		},
		"trigger": {
			description: "Trigger configuration.",
			markdown:    "Trigger configuration defining what initiates the workflow. Must include `type` (e.g., EVENT, SCHEDULED, REQUEST_RESPONSE) and optional `attributes` (trigger-specific configuration as JSON string).",
		},
		"enabled": {
			description: "Whether the workflow is enabled.",
			markdown:    "Whether the workflow is enabled (true) or disabled (false). Disabled workflows do not execute when triggered. Note: Workflows must be disabled before they can be deleted.",
		},
		"created": {
			description: "ISO-8601 timestamp when the workflow was created.",
			markdown:    "ISO-8601 timestamp when the workflow was created (computed).",
		},
		"modified": {
			description: "ISO-8601 timestamp when the workflow was last modified.",
			markdown:    "ISO-8601 timestamp when the workflow was last modified (computed).",
		},
	}
}
