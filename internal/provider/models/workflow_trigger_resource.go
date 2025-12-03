// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowTriggerResource represents the Terraform resource for managing a workflow trigger.
// This is distinct from WorkflowTrigger which is the nested block within Workflow.
type WorkflowTriggerResource struct {
	ID          types.String         `tfsdk:"id"`           // Composite ID: workflow_id
	WorkflowID  types.String         `tfsdk:"workflow_id"`  // Required: The workflow to attach this trigger to
	Type        types.String         `tfsdk:"type"`         // Required: The type of trigger (e.g., EVENT, SCHEDULED, REQUEST_RESPONSE)
	DisplayName types.String         `tfsdk:"display_name"` // Optional: Display name for the trigger
	Attributes  jsontypes.Normalized `tfsdk:"attributes"`   // Optional: Trigger-specific attributes as JSON
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API WorkflowTrigger.
func (wtr *WorkflowTriggerResource) ConvertToSailPoint(ctx context.Context) (*client.WorkflowTrigger, error) {
	if wtr == nil {
		return nil, nil
	}

	trigger := &client.WorkflowTrigger{
		Type: wtr.Type.ValueString(),
	}

	// Set optional displayName
	if !wtr.DisplayName.IsNull() && !wtr.DisplayName.IsUnknown() {
		trigger.DisplayName = wtr.DisplayName.ValueString()
	}

	// Parse attributes JSON string to map
	if !wtr.Attributes.IsNull() && !wtr.Attributes.IsUnknown() {
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(wtr.Attributes.ValueString()), &attributes); err != nil {
			return nil, err
		}
		trigger.Attributes = attributes
	}

	return trigger, nil
}

// ConvertFromSailPoint converts from SailPoint API to Terraform model.
// The workflowID parameter is needed since the API doesn't return it in the trigger object.
func (wtr *WorkflowTriggerResource) ConvertFromSailPoint(ctx context.Context, workflowID string, trigger *client.WorkflowTrigger) error {
	if wtr == nil || trigger == nil {
		return nil
	}

	// Set composite ID (using workflow_id as the resource ID since trigger is nested)
	wtr.ID = types.StringValue(workflowID)
	wtr.WorkflowID = types.StringValue(workflowID)
	wtr.Type = types.StringValue(trigger.Type)

	// Handle optional displayName
	if trigger.DisplayName != "" {
		wtr.DisplayName = types.StringValue(trigger.DisplayName)
	} else {
		wtr.DisplayName = types.StringNull()
	}

	// Convert attributes map to JSON string with normalization
	if len(trigger.Attributes) > 0 {
		attributesJSON, err := json.Marshal(trigger.Attributes)
		if err != nil {
			return err
		}
		wtr.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	} else {
		wtr.Attributes = jsontypes.NewNormalizedNull()
	}

	return nil
}
