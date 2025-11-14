// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Workflow represents the Terraform model for a SailPoint Workflow.
type Workflow struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"` // JSON string
	Description types.String `tfsdk:"description"`
	Definition  types.String `tfsdk:"definition"` // JSON string
	Trigger     types.String `tfsdk:"trigger"`    // JSON string
	Enabled     types.Bool   `tfsdk:"enabled"`
	Created     types.String `tfsdk:"created"`
	Modified    types.String `tfsdk:"modified"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API Workflow.
func (w *Workflow) ConvertToSailPoint(ctx context.Context) (*client.Workflow, error) {
	if w == nil {
		return nil, nil
	}

	workflow := &client.Workflow{
		Name: w.Name.ValueString(),
	}

	// Parse owner JSON string to map
	if !w.Owner.IsNull() && !w.Owner.IsUnknown() {
		var owner map[string]interface{}
		if err := json.Unmarshal([]byte(w.Owner.ValueString()), &owner); err != nil {
			return nil, err
		}
		workflow.Owner = owner
	}

	// Parse definition JSON string to map
	if !w.Definition.IsNull() && !w.Definition.IsUnknown() {
		var definition map[string]interface{}
		if err := json.Unmarshal([]byte(w.Definition.ValueString()), &definition); err != nil {
			return nil, err
		}
		workflow.Definition = definition
	}

	// Parse trigger JSON string to map
	if !w.Trigger.IsNull() && !w.Trigger.IsUnknown() {
		var trigger map[string]interface{}
		if err := json.Unmarshal([]byte(w.Trigger.ValueString()), &trigger); err != nil {
			return nil, err
		}
		workflow.Trigger = trigger
	}

	// Set optional fields
	if !w.Description.IsNull() && !w.Description.IsUnknown() {
		desc := w.Description.ValueString()
		workflow.Description = &desc
	}

	if !w.Enabled.IsNull() && !w.Enabled.IsUnknown() {
		enabled := w.Enabled.ValueBool()
		workflow.Enabled = &enabled
	}

	return workflow, nil
}

// ConvertFromSailPoint converts a SailPoint API Workflow to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (w *Workflow) ConvertFromSailPoint(ctx context.Context, workflow *client.Workflow, includeNull bool) error {
	if w == nil || workflow == nil {
		return nil
	}

	w.ID = types.StringValue(workflow.ID)
	w.Name = types.StringValue(workflow.Name)

	// Convert owner map to JSON string
	if workflow.Owner != nil {
		ownerJSON, err := json.Marshal(workflow.Owner)
		if err != nil {
			return err
		}
		w.Owner = types.StringValue(string(ownerJSON))
	} else if includeNull {
		w.Owner = types.StringNull()
	}

	// Convert definition map to JSON string
	if workflow.Definition != nil {
		definitionJSON, err := json.Marshal(workflow.Definition)
		if err != nil {
			return err
		}
		w.Definition = types.StringValue(string(definitionJSON))
	} else if includeNull {
		w.Definition = types.StringNull()
	}

	// Convert trigger map to JSON string
	if workflow.Trigger != nil {
		triggerJSON, err := json.Marshal(workflow.Trigger)
		if err != nil {
			return err
		}
		w.Trigger = types.StringValue(string(triggerJSON))
	} else if includeNull {
		w.Trigger = types.StringNull()
	}

	// Handle optional fields
	if workflow.Description != nil {
		w.Description = types.StringValue(*workflow.Description)
	} else if includeNull {
		w.Description = types.StringNull()
	}

	if workflow.Enabled != nil {
		w.Enabled = types.BoolValue(*workflow.Enabled)
	} else if includeNull {
		w.Enabled = types.BoolNull()
	}

	// Handle computed fields
	if workflow.Created != nil {
		w.Created = types.StringValue(*workflow.Created)
	} else if includeNull {
		w.Created = types.StringNull()
	}

	if workflow.Modified != nil {
		w.Modified = types.StringValue(*workflow.Modified)
	} else if includeNull {
		w.Modified = types.StringNull()
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (w *Workflow) ConvertFromSailPointForResource(ctx context.Context, workflow *client.Workflow) error {
	return w.ConvertFromSailPoint(ctx, workflow, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (w *Workflow) ConvertFromSailPointForDataSource(ctx context.Context, workflow *client.Workflow) error {
	return w.ConvertFromSailPoint(ctx, workflow, false)
}
