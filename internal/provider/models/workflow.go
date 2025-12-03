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

// Workflow represents the Terraform model for a SailPoint Workflow.
type Workflow struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Owner       *ObjectRef           `tfsdk:"owner"`
	Description types.String         `tfsdk:"description"`
	Definition  *WorkflowDefinition  `tfsdk:"definition"`
	Trigger     jsontypes.Normalized `tfsdk:"trigger"` // Computed field as JSON string to handle unknown values
	Enabled     types.Bool           `tfsdk:"enabled"`
	Created     types.String         `tfsdk:"created"`
	Modified    types.String         `tfsdk:"modified"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API Workflow.
func (w *Workflow) ConvertToSailPoint(ctx context.Context) (*client.Workflow, error) {
	if w == nil {
		return nil, nil
	}

	workflow := &client.Workflow{
		Name: w.Name.ValueString(),
	}

	// Convert owner ObjectRef
	if w.Owner != nil {
		ownerRef := w.Owner.ConvertToSailPoint(ctx)
		workflow.Owner = &ownerRef
	}

	// Convert definition WorkflowDefinition
	if w.Definition != nil {
		definition, err := w.Definition.ConvertToSailPoint(ctx)
		if err != nil {
			return nil, err
		}
		workflow.Definition = definition
	}

	// Note: Trigger is now computed and managed by sailpoint_workflow_trigger resource,
	// so it should not be included during Create/Update operations.
	// The trigger field is read-only and only populated on Read operations.

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

	// Convert owner ObjectRef
	if workflow.Owner != nil {
		w.Owner = &ObjectRef{}
		if includeNull {
			w.Owner.ConvertFromSailPointForResource(ctx, workflow.Owner)
		} else {
			w.Owner.ConvertFromSailPointForDataSource(ctx, workflow.Owner)
		}
	} else if includeNull {
		w.Owner = nil
	}

	// Convert definition WorkflowDefinition
	if workflow.Definition != nil {
		w.Definition = &WorkflowDefinition{}
		var err error
		if includeNull {
			err = w.Definition.ConvertFromSailPointForResource(ctx, workflow.Definition)
		} else {
			err = w.Definition.ConvertFromSailPointForDataSource(ctx, workflow.Definition)
		}
		if err != nil {
			return err
		}
	} else if includeNull {
		w.Definition = nil
	}

	// Convert trigger - now stored as JSON string
	// Trigger is computed and managed by sailpoint_workflow_trigger resource
	if workflow.Trigger != nil {
		triggerJSON, err := json.Marshal(workflow.Trigger)
		if err != nil {
			return err
		}
		w.Trigger = jsontypes.NewNormalizedValue(string(triggerJSON))
	} else if includeNull {
		w.Trigger = jsontypes.NewNormalizedNull()
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
