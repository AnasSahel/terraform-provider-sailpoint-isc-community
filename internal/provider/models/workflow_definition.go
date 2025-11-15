// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowDefinition represents the Terraform model for a workflow definition.
type WorkflowDefinition struct {
	Start types.String `tfsdk:"start"`
	Steps types.String `tfsdk:"steps"` // JSON string for flexibility
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API WorkflowDefinition.
func (d *WorkflowDefinition) ConvertToSailPoint(ctx context.Context) (*client.WorkflowDefinition, error) {
	if d == nil {
		return nil, nil
	}

	definition := &client.WorkflowDefinition{
		Start: d.Start.ValueString(),
	}

	// Parse steps JSON string to map
	if !d.Steps.IsNull() && !d.Steps.IsUnknown() && d.Steps.ValueString() != "" {
		var steps map[string]interface{}
		if err := json.Unmarshal([]byte(d.Steps.ValueString()), &steps); err != nil {
			return nil, err
		}
		definition.Steps = steps
	}

	return definition, nil
}

// ConvertFromSailPointForResource converts from SailPoint API to Terraform model for resources.
func (d *WorkflowDefinition) ConvertFromSailPointForResource(ctx context.Context, definition *client.WorkflowDefinition) error {
	if d == nil || definition == nil {
		return nil
	}

	d.Start = types.StringValue(definition.Start)

	// Convert steps map to JSON string
	if definition.Steps != nil && len(definition.Steps) > 0 {
		stepsJSON, err := json.Marshal(definition.Steps)
		if err != nil {
			return err
		}
		d.Steps = types.StringValue(string(stepsJSON))
	} else {
		d.Steps = types.StringNull()
	}

	return nil
}

// ConvertFromSailPointForDataSource converts from SailPoint API to Terraform model for data sources.
func (d *WorkflowDefinition) ConvertFromSailPointForDataSource(ctx context.Context, definition *client.WorkflowDefinition) error {
	if d == nil || definition == nil {
		return nil
	}

	d.Start = types.StringValue(definition.Start)

	// Convert steps map to JSON string
	if definition.Steps != nil && len(definition.Steps) > 0 {
		stepsJSON, err := json.Marshal(definition.Steps)
		if err != nil {
			return err
		}
		d.Steps = types.StringValue(string(stepsJSON))
	}
	// Don't set to null for data sources to preserve state

	return nil
}
