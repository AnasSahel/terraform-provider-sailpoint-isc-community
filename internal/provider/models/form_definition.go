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

// FormElementValidation represents a validation rule for a form element.
type FormElementValidation struct {
	ValidationType types.String `tfsdk:"validation_type"`
	Config         types.Object `tfsdk:"config"` // Optional config for the validation
}

// FormElement represents a single form element (field, section, etc).
type FormElement struct {
	ID           types.String            `tfsdk:"id"`
	ElementType  types.String            `tfsdk:"element_type"`
	Key          types.String            `tfsdk:"key"`
	Config       jsontypes.Normalized    `tfsdk:"config"` // Complex config as JSON
	Validations  []FormElementValidation `tfsdk:"validations"`
	FormElements []FormElement           `tfsdk:"form_elements"` // Nested elements for sections
}

// FormDefinition represents the Terraform model for a SailPoint Form Definition.
type FormDefinition struct {
	ID             types.String    `tfsdk:"id"`
	Name           types.String    `tfsdk:"name"`
	Description    types.String    `tfsdk:"description"`
	Owner          *ObjectRef      `tfsdk:"owner"`
	UsedBy         []ObjectRef     `tfsdk:"used_by"`         // List of object references
	FormInput      []FormInput     `tfsdk:"form_input"`      // List of form inputs
	FormElements   []FormElement   `tfsdk:"form_elements"`   // Structured list of form elements
	FormConditions []FormCondition `tfsdk:"form_conditions"` // List of form conditions
	Created        types.String    `tfsdk:"created"`
	Modified       types.String    `tfsdk:"modified"`
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
	if len(f.UsedBy) > 0 {
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

	// Convert FormElements array to slice of maps for API
	if len(f.FormElements) > 0 {
		elementMaps := make([]map[string]interface{}, len(f.FormElements))
		for i, elem := range f.FormElements {
			elementMaps[i] = elem.ConvertToSailPoint(ctx)
		}
		form.FormElements = elementMaps
	}

	// Convert FormInput array to slice of maps
	if len(f.FormInput) > 0 {
		inputMaps := make([]map[string]interface{}, len(f.FormInput))
		for i, input := range f.FormInput {
			inputMaps[i] = input.ConvertToSailPoint(ctx)
		}
		form.FormInput = inputMaps
	}

	// Convert FormConditions array to slice of maps
	if len(f.FormConditions) > 0 {
		conditionMaps := make([]map[string]interface{}, len(f.FormConditions))
		for i, condition := range f.FormConditions {
			conditionMaps[i] = condition.ConvertToSailPoint(ctx)
		}
		form.FormConditions = conditionMaps
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
	if len(form.UsedBy) > 0 {
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
	}
	// If nil or empty, leave UsedBy as nil to preserve null vs [] distinction

	// FormInput
	if len(form.FormInput) > 0 {
		formInputs := make([]FormInput, len(form.FormInput))
		for i, inputMap := range form.FormInput {
			formInputs[i].ConvertFromSailPoint(ctx, inputMap)
		}
		f.FormInput = formInputs
	}
	// If nil or empty, leave FormInput as nil to preserve null vs [] distinction

	// FormElements - convert to structured format
	if len(form.FormElements) > 0 {
		formElements := make([]FormElement, len(form.FormElements))
		for i, elemMap := range form.FormElements {
			formElements[i].ConvertFromSailPoint(ctx, elemMap)
		}
		f.FormElements = formElements
	}

	// FormConditions
	if len(form.FormConditions) > 0 {
		formConditions := make([]FormCondition, len(form.FormConditions))
		for i, conditionMap := range form.FormConditions {
			formConditions[i].ConvertFromSailPoint(ctx, conditionMap)
		}
		f.FormConditions = formConditions
	}
	// If nil or empty, leave FormConditions as nil to preserve null vs [] distinction

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
	if !formElementsEqual(f.FormElements, newForm.FormElements) {
		if len(newForm.FormElements) == 0 {
			operations = append(operations, map[string]interface{}{
				"op":   "remove",
				"path": "/formElements",
			})
		} else {
			elementMaps := make([]map[string]interface{}, len(newForm.FormElements))
			for i, elem := range newForm.FormElements {
				elementMaps[i] = elem.ConvertToSailPoint(ctx)
			}
			operations = append(operations, map[string]interface{}{
				"op":    "replace",
				"path":  "/formElements",
				"value": elementMaps,
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
	if !formConditionsEqual(f.FormConditions, newForm.FormConditions) {
		conditionMaps := make([]map[string]interface{}, len(newForm.FormConditions))
		for i, condition := range newForm.FormConditions {
			conditionMaps[i] = condition.ConvertToSailPoint(ctx)
		}
		operations = append(operations, map[string]interface{}{
			"op":    "replace",
			"path":  "/formConditions",
			"value": conditionMaps,
		})
	}

	return operations, nil
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

// formConditionsEqual compares two FormCondition slices to determine if they are equal.
func formConditionsEqual(a, b []FormCondition) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		// Compare rule operator
		if !a[i].RuleOperator.Equal(b[i].RuleOperator) {
			return false
		}

		// Compare rules
		if len(a[i].Rules) != len(b[i].Rules) {
			return false
		}
		for j := range a[i].Rules {
			if !a[i].Rules[j].SourceType.Equal(b[i].Rules[j].SourceType) ||
				!a[i].Rules[j].Source.Equal(b[i].Rules[j].Source) ||
				!a[i].Rules[j].Operator.Equal(b[i].Rules[j].Operator) ||
				!a[i].Rules[j].ValueType.Equal(b[i].Rules[j].ValueType) ||
				!a[i].Rules[j].Value.Equal(b[i].Rules[j].Value) {
				return false
			}
		}

		// Compare effects
		if len(a[i].Effects) != len(b[i].Effects) {
			return false
		}
		for j := range a[i].Effects {
			if !a[i].Effects[j].EffectType.Equal(b[i].Effects[j].EffectType) {
				return false
			}
			// Compare configs
			aConfig := a[i].Effects[j].Config
			bConfig := b[i].Effects[j].Config
			if (aConfig == nil) != (bConfig == nil) {
				return false
			}
			if aConfig != nil && bConfig != nil {
				if !aConfig.DefaultValueLabel.Equal(bConfig.DefaultValueLabel) ||
					!aConfig.Element.Equal(bConfig.Element) {
					return false
				}
			}
		}
	}

	return true
}

// ConvertToSailPoint converts a FormElement to a map[string]interface{} for the API.
func (fe *FormElement) ConvertToSailPoint(ctx context.Context) map[string]interface{} {
	element := make(map[string]interface{})

	if !fe.ID.IsNull() && !fe.ID.IsUnknown() {
		element["id"] = fe.ID.ValueString()
	}

	if !fe.ElementType.IsNull() && !fe.ElementType.IsUnknown() {
		element["elementType"] = fe.ElementType.ValueString()
	}

	if !fe.Key.IsNull() && !fe.Key.IsUnknown() {
		element["key"] = fe.Key.ValueString()
	}

	// Config is JSON, parse it to get the actual config object
	if !fe.Config.IsNull() && !fe.Config.IsUnknown() {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(fe.Config.ValueString()), &configMap); err == nil {
			element["config"] = configMap
		}
	}

	// Validations
	if len(fe.Validations) > 0 {
		validationMaps := make([]map[string]interface{}, len(fe.Validations))
		for i, val := range fe.Validations {
			validationMaps[i] = map[string]interface{}{
				"validationType": val.ValidationType.ValueString(),
			}
			// TODO: Handle validation config if needed
		}
		element["validations"] = validationMaps
	}

	// Nested form elements (for sections)
	if len(fe.FormElements) > 0 {
		nestedElementMaps := make([]map[string]interface{}, len(fe.FormElements))
		for i, nestedElem := range fe.FormElements {
			nestedElementMaps[i] = nestedElem.ConvertToSailPoint(ctx)
		}
		// Put nested elements in config.formElements
		if configMap, ok := element["config"].(map[string]interface{}); ok {
			configMap["formElements"] = nestedElementMaps
		}
	}

	return element
}

// formElementsEqual compares two FormElement slices for equality.
func formElementsEqual(a, b []FormElement) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].ID.Equal(b[i].ID) ||
			!a[i].ElementType.Equal(b[i].ElementType) ||
			!a[i].Key.Equal(b[i].Key) ||
			!a[i].Config.Equal(b[i].Config) {
			return false
		}
		// Compare validations
		if len(a[i].Validations) != len(b[i].Validations) {
			return false
		}
		for j := range a[i].Validations {
			if !a[i].Validations[j].ValidationType.Equal(b[i].Validations[j].ValidationType) {
				return false
			}
		}
		// Recursively compare nested form elements
		if !formElementsEqual(a[i].FormElements, b[i].FormElements) {
			return false
		}
	}
	return true
}

// ConvertFromSailPoint converts an API form element map to a Terraform FormElement.
func (fe *FormElement) ConvertFromSailPoint(ctx context.Context, elem map[string]interface{}) {
	if idVal, ok := elem["id"].(string); ok {
		fe.ID = types.StringValue(idVal)
	}

	if typeVal, ok := elem["elementType"].(string); ok {
		fe.ElementType = types.StringValue(typeVal)
	}

	if keyVal, ok := elem["key"].(string); ok {
		fe.Key = types.StringValue(keyVal)
	}

	// Handle config
	if configVal, ok := elem["config"].(map[string]interface{}); ok {
		// Remove formElements from config before storing as JSON
		configMap := make(map[string]interface{})
		for k, v := range configVal {
			if k != "formElements" {
				configMap[k] = v
			}
		}
		configJSON, _ := json.Marshal(configMap)
		fe.Config = jsontypes.NewNormalizedValue(string(configJSON))

		// Handle nested formElements if present
		if nestedElements, ok := configVal["formElements"].([]interface{}); ok {
			fe.FormElements = make([]FormElement, len(nestedElements))
			for i, nestedElem := range nestedElements {
				if nestedMap, ok := nestedElem.(map[string]interface{}); ok {
					fe.FormElements[i].ConvertFromSailPoint(ctx, nestedMap)
				}
			}
		}
	}

	// Handle validations
	if validationsVal, ok := elem["validations"].([]interface{}); ok {
		fe.Validations = make([]FormElementValidation, len(validationsVal))
		for i, valItem := range validationsVal {
			if valMap, ok := valItem.(map[string]interface{}); ok {
				if typeVal, ok := valMap["validationType"].(string); ok {
					fe.Validations[i].ValidationType = types.StringValue(typeVal)
				}
				// TODO: Handle validation config if needed
			}
		}
	}
}
