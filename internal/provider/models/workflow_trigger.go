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

// WorkflowTrigger represents the Terraform model for a workflow trigger.
type WorkflowTrigger struct {
	Type        types.String         `tfsdk:"type"`
	DisplayName types.String         `tfsdk:"display_name"`
	Attributes  jsontypes.Normalized `tfsdk:"attributes"` // JSON string with normalization
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API WorkflowTrigger.
func (t *WorkflowTrigger) ConvertToSailPoint(ctx context.Context) (*client.WorkflowTrigger, error) {
	if t == nil {
		return nil, nil
	}

	trigger := &client.WorkflowTrigger{
		Type: t.Type.ValueString(),
	}

	// Set optional displayName
	if !t.DisplayName.IsNull() && !t.DisplayName.IsUnknown() {
		trigger.DisplayName = t.DisplayName.ValueString()
	}

	// Parse attributes JSON string to map
	if !t.Attributes.IsNull() && !t.Attributes.IsUnknown() {
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(t.Attributes.ValueString()), &attributes); err != nil {
			return nil, err
		}
		trigger.Attributes = attributes
	}

	return trigger, nil
}

// ConvertFromSailPointForResource converts from SailPoint API to Terraform model for resources.
func (t *WorkflowTrigger) ConvertFromSailPointForResource(ctx context.Context, trigger *client.WorkflowTrigger) error {
	if t == nil || trigger == nil {
		return nil
	}

	t.Type = types.StringValue(trigger.Type)

	// Handle optional displayName
	if trigger.DisplayName != "" {
		t.DisplayName = types.StringValue(trigger.DisplayName)
	} else {
		t.DisplayName = types.StringNull()
	}

	// Convert attributes map to JSON string with normalization
	if len(trigger.Attributes) > 0 {
		attributesJSON, err := json.Marshal(trigger.Attributes)
		if err != nil {
			return err
		}
		t.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	} else {
		t.Attributes = jsontypes.NewNormalizedNull()
	}

	return nil
}

// ConvertFromSailPointForDataSource converts from SailPoint API to Terraform model for data sources.
func (t *WorkflowTrigger) ConvertFromSailPointForDataSource(ctx context.Context, trigger *client.WorkflowTrigger) error {
	if t == nil || trigger == nil {
		return nil
	}

	t.Type = types.StringValue(trigger.Type)

	// Handle optional displayName
	if trigger.DisplayName != "" {
		t.DisplayName = types.StringValue(trigger.DisplayName)
	}
	// Don't set to null for data sources to preserve state

	// Convert attributes map to JSON string with normalization
	if len(trigger.Attributes) > 0 {
		attributesJSON, err := json.Marshal(trigger.Attributes)
		if err != nil {
			return err
		}
		t.Attributes = jsontypes.NewNormalizedValue(string(attributesJSON))
	}
	// Don't set to null for data sources to preserve state

	return nil
}
