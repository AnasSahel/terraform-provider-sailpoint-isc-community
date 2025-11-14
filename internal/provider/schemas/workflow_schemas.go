// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
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
		"owner": resource_schema.StringAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Required:            true,
		},
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
		},
		"definition": resource_schema.StringAttribute{
			Description:         desc["definition"].description,
			MarkdownDescription: desc["definition"].markdown,
			Required:            true,
		},
		"trigger": resource_schema.StringAttribute{
			Description:         desc["trigger"].description,
			MarkdownDescription: desc["trigger"].markdown,
			Required:            true,
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
		"owner": datasource_schema.StringAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
		},
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"definition": datasource_schema.StringAttribute{
			Description:         desc["definition"].description,
			MarkdownDescription: desc["definition"].markdown,
			Computed:            true,
		},
		"trigger": datasource_schema.StringAttribute{
			Description:         desc["trigger"].description,
			MarkdownDescription: desc["trigger"].markdown,
			Computed:            true,
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
			description: "Owner of the workflow as a JSON string.",
			markdown:    "Owner of the workflow as a JSON string. Must be a valid identity reference with `type`, `id`, and `name` fields. Example: `{\"type\":\"IDENTITY\",\"id\":\"2c91808568c529c60168cca6f90c1313\",\"name\":\"William Wilson\"}`",
		},
		"description": {
			description: "Description of the workflow.",
			markdown:    "Description of the workflow's purpose and functionality.",
		},
		"definition": {
			description: "Workflow definition as a JSON string.",
			markdown:    "Workflow definition as a JSON string containing the workflow logic. Must include `start` (name of first step) and `steps` (object containing all workflow steps with their actions and configurations). See [Workflows Documentation](https://developer.sailpoint.com/docs/extensibility/workflows) for structure details.",
		},
		"trigger": {
			description: "Trigger configuration as a JSON string.",
			markdown:    "Trigger configuration as a JSON string defining what initiates the workflow. Must include `type` (e.g., `EVENT`) and `attributes` (trigger-specific configuration). See [Workflow Triggers](https://developer.sailpoint.com/docs/extensibility/workflows/triggers) for available trigger types.",
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
