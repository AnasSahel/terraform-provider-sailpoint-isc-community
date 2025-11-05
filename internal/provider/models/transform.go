package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Transform represents the Terraform model for a SailPoint Transform.
type Transform struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Attributes types.String `tfsdk:"attributes"` // JSON string
	Internal   types.Bool   `tfsdk:"internal"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API Transform.
func (t *Transform) ConvertToSailPoint(ctx context.Context) (*client.Transform, error) {
	if t == nil {
		return nil, nil
	}

	transform := &client.Transform{
		Name: t.Name.ValueString(),
		Type: t.Type.ValueString(),
	}

	// Parse attributes JSON string to map
	if !t.Attributes.IsNull() && !t.Attributes.IsUnknown() {
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(t.Attributes.ValueString()), &attributes); err != nil {
			return nil, err
		}
		transform.Attributes = attributes
	} else {
		// Default to empty map if not provided
		transform.Attributes = make(map[string]interface{})
	}

	return transform, nil
}

// ConvertFromSailPoint converts a SailPoint API Transform to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (t *Transform) ConvertFromSailPoint(ctx context.Context, transform *client.Transform, includeNull bool) error {
	if t == nil || transform == nil {
		return nil
	}

	t.ID = types.StringValue(transform.ID)
	t.Name = types.StringValue(transform.Name)
	t.Type = types.StringValue(transform.Type)
	t.Internal = types.BoolValue(transform.Internal)

	// Convert attributes map to JSON string
	if transform.Attributes != nil {
		attributesJSON, err := json.Marshal(transform.Attributes)
		if err != nil {
			return err
		}
		t.Attributes = types.StringValue(string(attributesJSON))
	} else if includeNull {
		t.Attributes = types.StringNull()
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (t *Transform) ConvertFromSailPointForResource(ctx context.Context, transform *client.Transform) error {
	return t.ConvertFromSailPoint(ctx, transform, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (t *Transform) ConvertFromSailPointForDataSource(ctx context.Context, transform *client.Transform) error {
	return t.ConvertFromSailPoint(ctx, transform, false)
}
