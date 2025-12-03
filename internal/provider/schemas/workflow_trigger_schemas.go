// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type WorkflowTriggerSchemaBuilder struct{}

// GetResourceSchema implements SchemaBuilder for WorkflowTrigger resource.
func (sb *WorkflowTriggerSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	return map[string]resource_schema.Attribute{
		"id": resource_schema.StringAttribute{
			Description:         "The composite ID of the workflow trigger (same as workflow_id).",
			MarkdownDescription: "The composite ID of the workflow trigger (same as workflow_id). Used by Terraform for resource identification.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"workflow_id": resource_schema.StringAttribute{
			Description:         "The ID of the workflow to attach this trigger to.",
			MarkdownDescription: "The ID of the workflow to attach this trigger to. This references the `sailpoint_workflow` resource.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"type": resource_schema.StringAttribute{
			Description:         "The type of trigger (e.g., EVENT, SCHEDULED, REQUEST_RESPONSE).",
			MarkdownDescription: "The type of trigger that initiates the workflow. Common types include EVENT, SCHEDULED, and REQUEST_RESPONSE.",
			Required:            true,
		},
		"display_name": resource_schema.StringAttribute{
			Description:         "Display name for the trigger.",
			MarkdownDescription: "An optional display name for the trigger that appears in the SailPoint UI.",
			Optional:            true,
		},
		"attributes": resource_schema.StringAttribute{
			Description:         "Trigger-specific attributes as a JSON string. Structure varies by trigger type.",
			MarkdownDescription: "Trigger-specific configuration attributes as a JSON string. The structure and required fields depend on the trigger `type`. For example, EVENT triggers may require event IDs, while SCHEDULED triggers require cron expressions.",
			Optional:            true,
			CustomType:          jsontypes.NormalizedType{},
		},
	}
}

// WorkflowTriggerSchemaBuilder is a standalone schema builder and doesn't implement data source schema.
// Triggers are typically managed via the resource and cannot be independently queried as data sources.
