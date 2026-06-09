// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"encoding/json"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stepAttrTypes defines the attribute types for a single typed step object.
var stepAttrTypes = map[string]attr.Type{
	"type":           types.StringType,
	"action_id":      types.StringType,
	"attributes":     jsontypes.NormalizedType{},
	"next_step":      types.StringType,
	"display_name":   types.StringType,
	"version_number": types.Int32Type,
	"catch":          jsontypes.NormalizedType{},
	"config":         jsontypes.NormalizedType{},
}

// definitionAttrTypes defines the attribute types for the workflow definition object.
var definitionAttrTypes = map[string]attr.Type{
	"start": types.StringType,
	"steps": types.MapType{ElemType: types.ObjectType{AttrTypes: stepAttrTypes}},
}

// promotedStepKeys are API-level keys extracted into named step attributes.
// All other keys from a flat step object are collected into the "config" catch-all.
var promotedStepKeys = map[string]bool{
	"type":          true,
	"actionId":      true,
	"attributes":    true,
	"nextStep":      true,
	"displayName":   true,
	"versionNumber": true,
	"catch":         true,
}

// workflowModel represents the Terraform model for a SailPoint Workflow.
type workflowModel struct {
	ID                types.String         `tfsdk:"id"`
	Name              types.String         `tfsdk:"name"`
	Owner             types.Object         `tfsdk:"owner"`
	Description       types.String         `tfsdk:"description"`
	Definition        types.Object         `tfsdk:"definition"`
	Trigger           jsontypes.Normalized `tfsdk:"trigger"`
	Enabled           types.Bool           `tfsdk:"enabled"`
	Created           types.String         `tfsdk:"created"`
	Modified          types.String         `tfsdk:"modified"`
	Creator           types.Object         `tfsdk:"creator"`
	ModifiedBy        types.Object         `tfsdk:"modified_by"`
	ExecutionCount    types.Int32          `tfsdk:"execution_count"`
	FailureCount      types.Int32          `tfsdk:"failure_count"`
	IgnoreJSONChanges types.List           `tfsdk:"ignore_json_changes"`
}

// objectRefAttrTypes returns the attribute types for ObjectRef-like nested objects.
func objectRefAttrTypes() map[string]attr.Type {
	return common.ObjectRefObjectType.AttrTypes
}

// objectRefFromAPI converts a *client.ObjectRefAPI to a types.Object.
func objectRefFromAPI(api *client.ObjectRefAPI) (types.Object, diag.Diagnostics) {
	if api == nil || api.Type == "" {
		return types.ObjectNull(objectRefAttrTypes()), nil
	}
	return types.ObjectValue(objectRefAttrTypes(), map[string]attr.Value{
		"type": types.StringValue(api.Type),
		"id":   types.StringValue(api.ID),
		"name": common.StringOrNullIfEmpty(api.Name),
	})
}

// splitAPIStep converts a flat API step map into a typed step attribute map.
// This is the single shared function used by both FromAPI and the state upgraders.
// Promoted keys become named fields; all other keys go into "config" (JSON catch-all).
func splitAPIStep(raw map[string]any) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	getString := func(key string) types.String {
		if v, ok := raw[key].(string); ok && v != "" {
			return types.StringValue(v)
		}
		return types.StringNull()
	}

	getJSON := func(key string) jsontypes.Normalized {
		v, ok := raw[key]
		if !ok || v == nil {
			return jsontypes.NewNormalizedNull()
		}
		b, err := json.Marshal(v)
		if err != nil {
			diags.AddError("Error marshaling JSON field "+key, err.Error())
			return jsontypes.NewNormalizedNull()
		}
		return jsontypes.NewNormalizedValue(string(b))
	}

	// JSON unmarshal always decodes numbers as float64 for map[string]any.
	var versionNum types.Int32
	if v, ok := raw["versionNumber"]; ok && v != nil {
		if n, ok := v.(float64); ok {
			versionNum = types.Int32Value(int32(n))
		} else {
			versionNum = types.Int32Null()
		}
	} else {
		versionNum = types.Int32Null()
	}

	configMap := map[string]any{}
	for k, v := range raw {
		if !promotedStepKeys[k] {
			configMap[k] = v
		}
	}

	var configVal jsontypes.Normalized
	if len(configMap) > 0 {
		b, err := json.Marshal(configMap)
		if err != nil {
			diags.AddError("Error marshaling step config", err.Error())
			return nil, diags
		}
		configVal = jsontypes.NewNormalizedValue(string(b))
	} else {
		configVal = jsontypes.NewNormalizedNull()
	}

	return map[string]attr.Value{
		"type":           getString("type"),
		"action_id":      getString("actionId"),
		"attributes":     getJSON("attributes"),
		"next_step":      getString("nextStep"),
		"display_name":   getString("displayName"),
		"version_number": versionNum,
		"catch":          getJSON("catch"),
		"config":         configVal,
	}, diags
}

// stepsFromAPI converts the API's flat steps map into a types.Map of typed step objects.
func stepsFromAPI(apiSteps map[string]any) (types.Map, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: stepAttrTypes}

	if len(apiSteps) == 0 {
		result, diags := types.MapValue(elemType, map[string]attr.Value{})
		diagnostics.Append(diags...)
		return result, diagnostics
	}

	stepObjs := make(map[string]attr.Value, len(apiSteps))
	for name, raw := range apiSteps {
		rawMap, ok := raw.(map[string]any)
		if !ok {
			diagnostics.AddError("Invalid step format", "step "+name+" is not an object")
			return types.MapNull(elemType), diagnostics
		}
		attrs, diags := splitAPIStep(rawMap)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return types.MapNull(elemType), diagnostics
		}
		obj, diags := types.ObjectValue(stepAttrTypes, attrs)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return types.MapNull(elemType), diagnostics
		}
		stepObjs[name] = obj
	}

	result, diags := types.MapValue(elemType, stepObjs)
	diagnostics.Append(diags...)
	return result, diagnostics
}

// mergeStepToAPI converts a typed step object back to a flat API map.
// Promoted fields take precedence over any conflicting keys from "config".
func mergeStepToAPI(stepAttrs map[string]attr.Value) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics
	flat := map[string]any{}

	setString := func(attrKey, apiKey string) {
		if v, ok := stepAttrs[attrKey].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			flat[apiKey] = v.ValueString()
		}
	}

	setString("type", "type")
	setString("action_id", "actionId")
	setString("next_step", "nextStep")
	setString("display_name", "displayName")

	if v, ok := stepAttrs["version_number"].(types.Int32); ok && !v.IsNull() && !v.IsUnknown() {
		flat["versionNumber"] = v.ValueInt32()
	}

	unmarshalInto := func(attrKey, apiKey string) {
		v, ok := stepAttrs[attrKey].(jsontypes.Normalized)
		if !ok || v.IsNull() || v.IsUnknown() {
			return
		}
		var parsed any
		if err := json.Unmarshal([]byte(v.ValueString()), &parsed); err != nil {
			diags.AddError("Error parsing JSON field "+attrKey, err.Error())
			return
		}
		flat[apiKey] = parsed
	}

	unmarshalInto("attributes", "attributes")
	unmarshalInto("catch", "catch")

	// Merge config catch-all; promoted fields already set above win on collision.
	if v, ok := stepAttrs["config"].(jsontypes.Normalized); ok && !v.IsNull() && !v.IsUnknown() {
		var configMap map[string]any
		if err := json.Unmarshal([]byte(v.ValueString()), &configMap); err != nil {
			diags.AddError("Error parsing step config JSON", err.Error())
		} else {
			for k, val := range configMap {
				if _, exists := flat[k]; !exists {
					flat[k] = val
				}
			}
		}
	}

	return flat, diags
}

// ToAPI converts the Terraform model to a SailPoint API create request.
func (m *workflowModel) ToAPI(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.WorkflowAPI{
		Name: m.Name.ValueString(),
	}

	if !m.Owner.IsNull() && !m.Owner.IsUnknown() {
		attrs := m.Owner.Attributes()
		if attrs != nil {
			if t, ok := attrs["type"].(types.String); ok && !t.IsNull() && !t.IsUnknown() {
				if id, ok := attrs["id"].(types.String); ok && !id.IsNull() && !id.IsUnknown() {
					api.Owner = &client.ObjectRefAPI{
						Type: t.ValueString(),
						ID:   id.ValueString(),
					}
				}
			}
		}
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		api.Description = m.Description.ValueString()
	}

	if !m.Definition.IsNull() && !m.Definition.IsUnknown() {
		attrs := m.Definition.Attributes()
		if attrs != nil {
			if start, ok := attrs["start"].(types.String); ok && !start.IsNull() && !start.IsUnknown() {
				def := &client.WorkflowDefinitionAPI{
					Start: start.ValueString(),
					Steps: map[string]any{},
				}

				if stepsMap, ok := attrs["steps"].(types.Map); ok && !stepsMap.IsNull() && !stepsMap.IsUnknown() {
					for name, stepVal := range stepsMap.Elements() {
						stepObj, ok := stepVal.(types.Object)
						if !ok {
							continue
						}
						flat, diags := mergeStepToAPI(stepObj.Attributes())
						diagnostics.Append(diags...)
						if diagnostics.HasError() {
							return api, diagnostics
						}
						def.Steps[name] = flat
					}
				}

				api.Definition = def
			}
		}
	}

	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		api.Enabled = m.Enabled.ValueBool()
	}

	return api, diagnostics
}

// ToAPIUpdate converts the Terraform model to a SailPoint API update request (PUT).
func (m *workflowModel) ToAPIUpdate(ctx context.Context) (client.WorkflowAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}

// FromAPI populates the Terraform model from a SailPoint API response.
func (m *workflowModel) FromAPI(ctx context.Context, api client.WorkflowAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)

	m.Owner, diags = objectRefFromAPI(api.Owner)
	diagnostics.Append(diags...)

	m.Description = common.StringOrNullIfEmpty(api.Description)

	if api.Definition != nil && api.Definition.Start != "" {
		stepsMap, d := stepsFromAPI(api.Definition.Steps)
		diagnostics.Append(d...)
		if diagnostics.HasError() {
			return diagnostics
		}

		defObj, d := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
			"start": types.StringValue(api.Definition.Start),
			"steps": stepsMap,
		})
		diagnostics.Append(d...)
		m.Definition = defObj
	} else {
		m.Definition = types.ObjectNull(definitionAttrTypes)
	}

	if api.Trigger != nil {
		m.Trigger, diags = common.MarshalJSONOrDefault(api.Trigger, "{}")
		diagnostics.Append(diags...)
	} else {
		m.Trigger = jsontypes.NewNormalizedNull()
	}

	m.Enabled = types.BoolValue(api.Enabled)
	m.Created = common.StringOrNullIfEmpty(api.Created)
	m.Modified = common.StringOrNullIfEmpty(api.Modified)

	m.Creator, diags = objectRefFromAPI(api.Creator)
	diagnostics.Append(diags...)

	m.ModifiedBy, diags = objectRefFromAPI(api.ModifiedBy)
	diagnostics.Append(diags...)

	m.ExecutionCount = types.Int32Value(api.ExecutionCount)
	m.FailureCount = types.Int32Value(api.FailureCount)

	return diagnostics
}
