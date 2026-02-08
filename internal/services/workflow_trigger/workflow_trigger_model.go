// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// workflowTriggerModel represents the Terraform model for managing a workflow trigger.
// This is a separate resource from the workflow itself to allow flexible trigger management.
type workflowTriggerModel struct {
	WorkflowID  types.String         `tfsdk:"workflow_id"` // Required: The workflow to attach this trigger to
	Type        types.String         `tfsdk:"type"`        // Required: EVENT, EXTERNAL, SCHEDULED
	DisplayName types.String         `tfsdk:"display_name"`
	Attributes  jsontypes.Normalized `tfsdk:"attributes"` // Trigger-specific attributes as JSON
}

// ToAPIRequest converts the Terraform model to a SailPoint API WorkflowTrigger.
func (m *workflowTriggerModel) ToAPIRequest(ctx context.Context) (*client.WorkflowTriggerAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	trigger := &client.WorkflowTriggerAPI{
		Type: m.Type.ValueString(),
	}

	// Set optional displayName
	if !m.DisplayName.IsNull() && !m.DisplayName.IsUnknown() {
		trigger.DisplayName = m.DisplayName.ValueString()
	}

	// Parse attributes JSON string to map
	if !m.Attributes.IsNull() && !m.Attributes.IsUnknown() {
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(m.Attributes.ValueString()), &attributes); err != nil {
			diags.AddError("Invalid Attributes JSON", err.Error())
			return nil, diags
		}
		trigger.Attributes = attributes
	}

	return trigger, diags
}

// FromSailPointAPI populates the Terraform model from a SailPoint API response.
// The workflowID parameter is needed since the API doesn't return it in the trigger object.
func (m *workflowTriggerModel) FromSailPointAPI(ctx context.Context, workflowID string, trigger *client.WorkflowTriggerAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	if trigger == nil {
		return diags
	}

	m.WorkflowID = types.StringValue(workflowID)
	m.Type = types.StringValue(trigger.Type)

	// Handle optional displayName
	if trigger.DisplayName != "" {
		m.DisplayName = types.StringValue(trigger.DisplayName)
	} else {
		m.DisplayName = types.StringNull()
	}

	// Convert attributes map to JSON string with normalization
	if len(trigger.Attributes) > 0 {
		attributesJSON, err := json.Marshal(trigger.Attributes)
		if err != nil {
			diags.AddError("Failed to marshal attributes", err.Error())
			return diags
		}
		m.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	} else {
		m.Attributes = jsontypes.NewNormalizedNull()
	}

	return diags
}
