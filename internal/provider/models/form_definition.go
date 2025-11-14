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
	UsedBy         []ObjectRef  `tfsdk:"used_by"`          // List of object references
	FormInput      []FormInput  `tfsdk:"form_input"`       // List of form inputs
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

	// UsedBy
	if f.UsedBy != nil && len(f.UsedBy) > 0 {
		usedByMaps := make([]map[string]interface{}, len(f.UsedBy))
		for i, ref := range f.UsedBy {
			apiRef := ref.ConvertToSailPoint(ctx)
			usedByMaps[i] = map[string]interface{}{
				"type": apiRef.Type,
				"id":   apiRef.ID,
			}
			if apiRef.Name != "" {
				usedByMaps[i]["name"] = apiRef.Name
			}
		}
		form.UsedBy = usedByMaps
	}

	// Parse FormElements JSON string to slice of maps
	if !f.FormElements.IsNull() && !f.FormElements.IsUnknown() {
		var elements []map[string]interface{}
		if err := json.Unmarshal([]byte(f.FormElements.ValueString()), &elements); err != nil {
			return nil, err
		}
		form.FormElements = elements
	}

	// Convert FormInput array to slice of maps
	if f.FormInput != nil && len(f.FormInput) > 0 {
		inputMaps := make([]map[string]interface{}, len(f.FormInput))
		for i, input := range f.FormInput {
			inputMaps[i] = input.ConvertToSailPoint(ctx)
		}
		form.FormInput = inputMaps
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
		usedByRefs := make([]ObjectRef, len(form.UsedBy))
		for i, usedByMap := range form.UsedBy {
			// Convert map to client.ObjectRef
			objRef := &client.ObjectRef{}
			if typeVal, ok := usedByMap["type"].(string); ok {
				objRef.Type = typeVal
			}
			if idVal, ok := usedByMap["id"].(string); ok {
				objRef.ID = idVal
			}
			if nameVal, ok := usedByMap["name"].(string); ok {
				objRef.Name = nameVal
			}
			usedByRefs[i].ConvertFromSailPointForResource(ctx, objRef)
		}
		f.UsedBy = usedByRefs
	} else {
		f.UsedBy = []ObjectRef{}
	}

	// FormInput
	if form.FormInput != nil && len(form.FormInput) > 0 {
		formInputs := make([]FormInput, len(form.FormInput))
		for i, inputMap := range form.FormInput {
			formInputs[i].ConvertFromSailPoint(ctx, inputMap)
		}
		f.FormInput = formInputs
	} else {
		f.FormInput = []FormInput{}
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

	// UsedBy
	if !usedByEqual(f.UsedBy, newForm.UsedBy) {
		usedByMaps := make([]map[string]interface{}, len(newForm.UsedBy))
		for i, ref := range newForm.UsedBy {
			apiRef := ref.ConvertToSailPoint(ctx)
			usedByMaps[i] = map[string]interface{}{
				"type": apiRef.Type,
				"id":   apiRef.ID,
			}
			if apiRef.Name != "" {
				usedByMaps[i]["name"] = apiRef.Name
			}
		}
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/usedBy",
			"value": usedByMaps,
		})
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
	if !formInputEqual(f.FormInput, newForm.FormInput) {
		inputMaps := make([]map[string]interface{}, len(newForm.FormInput))
		for i, input := range newForm.FormInput {
			inputMaps[i] = input.ConvertToSailPoint(ctx)
		}
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/formInput",
			"value": inputMaps,
		})
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

// usedByEqual compares two UsedBy slices to determine if they are equal.
// Two UsedBy slices are equal if they have the same length and all corresponding
// elements have matching type, id, and name values.
func usedByEqual(a, b []ObjectRef) bool {
	if len(a) != len(b) {
		return false
	}

	// Create a map to track matches
	// This handles cases where the order might be different
	matchedIndices := make(map[int]bool)

	for _, aRef := range a {
		found := false
		for j, bRef := range b {
			if matchedIndices[j] {
				continue // Already matched
			}
			if aRef.Type.Equal(bRef.Type) && aRef.ID.Equal(bRef.ID) && aRef.Name.Equal(bRef.Name) {
				matchedIndices[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// formInputEqual compares two FormInput slices to determine if they are equal.
// Two FormInput slices are equal if they have the same length and all corresponding
// elements have matching id, type, label, and description values.
func formInputEqual(a, b []FormInput) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].ID.Equal(b[i].ID) ||
			!a[i].Type.Equal(b[i].Type) ||
			!a[i].Label.Equal(b[i].Label) ||
			!a[i].Description.Equal(b[i].Description) {
			return false
		}
	}

	return true
}
