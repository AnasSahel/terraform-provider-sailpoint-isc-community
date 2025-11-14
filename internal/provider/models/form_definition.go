// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormDefinition represents the Terraform model for a SailPoint Form Definition.
type FormDefinition struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Owner          *ObjectRef   `tfsdk:"owner"`
	UsedBy         types.String `tfsdk:"used_by"`          // JSON string
	FormInput      types.String `tfsdk:"form_input"`       // JSON string
	FormElements   types.String `tfsdk:"form_elements"`    // JSON string
	FormConditions types.String `tfsdk:"form_conditions"`  // JSON string
	Created        types.String `tfsdk:"created"`
	Modified       types.String `tfsdk:"modified"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API FormDefinition.
func (f *FormDefinition) ConvertToSailPoint(ctx context.Context) (*client.FormDefinition, error) {
	if f == nil {
		return nil, nil
	}

	form := &client.FormDefinition{
		Name: f.Name.ValueString(),
	}

	// Description
	if !f.Description.IsNull() && !f.Description.IsUnknown() {
		form.Description = f.Description.ValueString()
	}

	// Owner
	if f.Owner != nil {
		ownerRef := f.Owner.ConvertToSailPoint(ctx)
		form.Owner = &ownerRef
	}

	// Parse FormElements JSON string to slice of maps
	if !f.FormElements.IsNull() && !f.FormElements.IsUnknown() {
		var elements []map[string]interface{}
		if err := json.Unmarshal([]byte(f.FormElements.ValueString()), &elements); err != nil {
			return nil, err
		}
		form.FormElements = elements
	}

	// Parse FormInput JSON string to slice of maps
	if !f.FormInput.IsNull() && !f.FormInput.IsUnknown() {
		var input []map[string]interface{}
		if err := json.Unmarshal([]byte(f.FormInput.ValueString()), &input); err != nil {
			return nil, err
		}
		form.FormInput = input
	}

	// Parse FormConditions JSON string to slice of maps
	if !f.FormConditions.IsNull() && !f.FormConditions.IsUnknown() {
		var conditions []map[string]interface{}
		if err := json.Unmarshal([]byte(f.FormConditions.ValueString()), &conditions); err != nil {
			return nil, err
		}
		form.FormConditions = conditions
	}

	return form, nil
}

// ConvertFromSailPoint converts a SailPoint API FormDefinition to the Terraform model.
// For resources, set includeNull to true. For data sources, set to false.
func (f *FormDefinition) ConvertFromSailPoint(ctx context.Context, form *client.FormDefinition, includeNull bool) error {
	if f == nil || form == nil {
		return nil
	}

	f.ID = types.StringValue(form.ID)
	f.Name = types.StringValue(form.Name)

	// Description
	if form.Description != "" {
		f.Description = types.StringValue(form.Description)
	} else if includeNull {
		f.Description = types.StringNull()
	}

	// Owner
	if form.Owner != nil {
		f.Owner = &ObjectRef{}
		f.Owner.ConvertFromSailPointForResource(ctx, form.Owner)
	} else if includeNull {
		f.Owner = nil
	}

	// UsedBy
	if form.UsedBy != nil && len(form.UsedBy) > 0 {
		usedByJSON, err := json.Marshal(form.UsedBy)
		if err != nil {
			return err
		}
		f.UsedBy = types.StringValue(string(usedByJSON))
	} else if includeNull {
		f.UsedBy = types.StringNull()
	}

	// FormInput
	if form.FormInput != nil && len(form.FormInput) > 0 {
		inputJSON, err := json.Marshal(form.FormInput)
		if err != nil {
			return err
		}
		f.FormInput = types.StringValue(string(inputJSON))
	} else if includeNull {
		f.FormInput = types.StringNull()
	}

	// FormElements
	if form.FormElements != nil && len(form.FormElements) > 0 {
		// Normalize form elements by removing empty validations arrays that the API adds
		normalizedElements := normalizeFormElements(form.FormElements)
		elementsJSON, err := json.Marshal(normalizedElements)
		if err != nil {
			return err
		}
		f.FormElements = types.StringValue(string(elementsJSON))
	} else if includeNull {
		f.FormElements = types.StringNull()
	}

	// FormConditions
	if form.FormConditions != nil && len(form.FormConditions) > 0 {
		conditionsJSON, err := json.Marshal(form.FormConditions)
		if err != nil {
			return err
		}
		f.FormConditions = types.StringValue(string(conditionsJSON))
	} else if includeNull {
		f.FormConditions = types.StringNull()
	}

	// Created and Modified timestamps
	if form.Created != "" {
		f.Created = types.StringValue(form.Created)
	} else if includeNull {
		f.Created = types.StringNull()
	}

	if form.Modified != "" {
		f.Modified = types.StringValue(form.Modified)
	} else if includeNull {
		f.Modified = types.StringNull()
	}

	return nil
}

// ConvertFromSailPointForResource converts for resource operations (includes all fields).
func (f *FormDefinition) ConvertFromSailPointForResource(ctx context.Context, form *client.FormDefinition) error {
	return f.ConvertFromSailPoint(ctx, form, true)
}

// ConvertFromSailPointForDataSource converts for data source operations (preserves state).
func (f *FormDefinition) ConvertFromSailPointForDataSource(ctx context.Context, form *client.FormDefinition) error {
	return f.ConvertFromSailPoint(ctx, form, false)
}

// GeneratePatchOperations generates JSON Patch operations for updating a form definition.
func (f *FormDefinition) GeneratePatchOperations(ctx context.Context, newForm FormDefinition) ([]map[string]interface{}, error) {
	var operations []map[string]interface{}

	// Name
	if !f.Name.Equal(newForm.Name) {
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/name",
			"value": newForm.Name.ValueString(),
		})
	}

	// Description
	if !f.Description.Equal(newForm.Description) {
		if newForm.Description.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/description",
			})
		} else {
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/description",
				"value": newForm.Description.ValueString(),
			})
		}
	}

	// Owner
	if (f.Owner == nil && newForm.Owner != nil) || (f.Owner != nil && newForm.Owner == nil) ||
	   (f.Owner != nil && newForm.Owner != nil && (!f.Owner.ID.Equal(newForm.Owner.ID) || !f.Owner.Type.Equal(newForm.Owner.Type))) {
		if newForm.Owner == nil {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/owner",
			})
		} else {
			ownerRef := newForm.Owner.ConvertToSailPoint(ctx)
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/owner",
				"value": &ownerRef,
			})
		}
	}

	// FormElements
	if !f.FormElements.Equal(newForm.FormElements) {
		if newForm.FormElements.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/formElements",
			})
		} else {
			var elements []map[string]interface{}
			if err := json.Unmarshal([]byte(newForm.FormElements.ValueString()), &elements); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/formElements",
				"value": elements,
			})
		}
	}

	// FormInput
	if !f.FormInput.Equal(newForm.FormInput) {
		if newForm.FormInput.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/formInput",
			})
		} else {
			var input []map[string]interface{}
			if err := json.Unmarshal([]byte(newForm.FormInput.ValueString()), &input); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/formInput",
				"value": input,
			})
		}
	}

	// FormConditions
	if !f.FormConditions.Equal(newForm.FormConditions) {
		if newForm.FormConditions.IsNull() {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/formConditions",
			})
		} else {
			var conditions []map[string]interface{}
			if err := json.Unmarshal([]byte(newForm.FormConditions.ValueString()), &conditions); err != nil {
				return nil, err
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/formConditions",
				"value": conditions,
			})
		}
	}

	return operations, nil
}

// normalizeFormElements removes empty arrays and API-added fields from form elements
// to prevent state inconsistency errors. The SailPoint API adds empty "validations"
// arrays to form elements even when not provided, which causes Terraform to detect
// a diff between the plan and the actual state.
func normalizeFormElements(elements []map[string]interface{}) []map[string]interface{} {
	normalized := make([]map[string]interface{}, len(elements))

	for i, element := range elements {
		normalizedElement := make(map[string]interface{})

		for key, value := range element {
			// Skip empty validations arrays
			if key == "validations" {
				if arr, ok := value.([]interface{}); ok && len(arr) == 0 {
					continue
				}
			}

			// Recursively normalize nested formElements (for sections)
			if key == "config" {
				if configMap, ok := value.(map[string]interface{}); ok {
					normalizedConfig := make(map[string]interface{})
					for configKey, configValue := range configMap {
						if configKey == "formElements" {
							if nestedElements, ok := configValue.([]interface{}); ok {
								// Convert to []map[string]interface{} for recursion
								nestedMaps := make([]map[string]interface{}, len(nestedElements))
								for j, ne := range nestedElements {
									if neMap, ok := ne.(map[string]interface{}); ok {
										nestedMaps[j] = neMap
									}
								}
								normalizedConfig[configKey] = normalizeFormElements(nestedMaps)
								continue
							}
						}
						normalizedConfig[configKey] = configValue
					}
					normalizedElement[key] = normalizedConfig
					continue
				}
			}

			normalizedElement[key] = value
		}

		normalized[i] = normalizedElement
	}

	return normalized
}
