// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
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

// ToAPI converts the Terraform model to a SailPoint API WorkflowTrigger.
func (m *workflowTriggerModel) ToAPI(ctx context.Context) (*client.WorkflowTriggerAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	trigger := &client.WorkflowTriggerAPI{
		Type: m.Type.ValueString(),
	}

	// Set optional displayName
	if !m.DisplayName.IsNull() && !m.DisplayName.IsUnknown() {
		trigger.DisplayName = m.DisplayName.ValueString()
	}

	// Parse attributes JSON string to map using common helper
	if attributes, diags := common.UnmarshalJSONField[map[string]interface{}](m.Attributes); attributes != nil {
		trigger.Attributes = *attributes
		diagnostics.Append(diags...)
	}

	return trigger, diagnostics
}

// FromAPI populates the Terraform model from a SailPoint API response.
// The workflowID parameter is needed since the API doesn't return it in the trigger object.
func (m *workflowTriggerModel) FromAPI(ctx context.Context, workflowID string, trigger *client.WorkflowTriggerAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	if trigger == nil {
		return diagnostics
	}

	m.WorkflowID = types.StringValue(workflowID)
	m.Type = types.StringValue(trigger.Type)

	// Handle optional displayName
	m.DisplayName = common.StringOrNullIfEmpty(trigger.DisplayName)

	// Convert attributes map to JSON string (nil → null, empty {} → "{}")
	if trigger.Attributes != nil {
		m.Attributes, diags = common.MarshalJSONOrDefault(trigger.Attributes, "{}")
		diagnostics.Append(diags...)
	} else {
		m.Attributes = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}
