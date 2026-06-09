// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common/jsonpath"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resolveIgnoreJSONPath parses a resource-rooted ignore_json_changes path of the form
// "definition.steps['<step>'].<field>.<json-path>" into (stepName, field, inJSONPath).
// field must be one of "attributes", "config", or "catch"; inJSONPath is used as "$.<inJSONPath>".
func resolveIgnoreJSONPath(p string) (stepName, field, inJSONPath string, err error) {
	const prefix = "definition.steps["
	if !strings.HasPrefix(p, prefix) {
		return "", "", "", fmt.Errorf("path %q must start with %q followed by a bracket-quoted step name", p, prefix)
	}
	rest := p[len(prefix):]

	bracketEnd := strings.Index(rest, "]")
	if bracketEnd < 0 {
		return "", "", "", fmt.Errorf("path %q: unclosed '[' in step name", p)
	}
	inner := rest[:bracketEnd]
	rest = rest[bracketEnd+1:]

	if len(inner) < 2 {
		return "", "", "", fmt.Errorf("path %q: step name must be bracket-quoted (e.g. ['My Step'])", p)
	}
	quote := inner[0]
	if quote != '\'' && quote != '"' {
		return "", "", "", fmt.Errorf("path %q: step name must be bracket-quoted with single or double quotes", p)
	}
	if inner[len(inner)-1] != quote {
		return "", "", "", fmt.Errorf("path %q: mismatched quotes in step name", p)
	}
	stepName = inner[1 : len(inner)-1]
	if stepName == "" {
		return "", "", "", fmt.Errorf("path %q: step name must not be empty", p)
	}

	if !strings.HasPrefix(rest, ".") {
		return "", "", "", fmt.Errorf("path %q: expected '.<field>' after step name", p)
	}
	rest = rest[1:]

	dotIdx := strings.Index(rest, ".")
	if dotIdx < 0 {
		return "", "", "", fmt.Errorf("path %q: must specify an in-JSON path after the field (e.g. .attributes.param_oauth.refID)", p)
	}
	field = rest[:dotIdx]
	rest = rest[dotIdx+1:]

	switch field {
	case "attributes", "config", "catch":
	default:
		return "", "", "", fmt.Errorf("path %q: field %q is not a JSON step field; must be one of: attributes, config, catch", p, field)
	}

	if rest == "" {
		return "", "", "", fmt.Errorf("path %q: in-JSON path must not be empty", p)
	}

	if err := jsonpath.Validate("$." + rest); err != nil {
		return "", "", "", fmt.Errorf("path %q: invalid in-JSON path %q: %w", p, "$."+rest, err)
	}

	return stepName, field, rest, nil
}

// applyIgnoreJSONChanges re-injects sourceModel values at each ignored path to suppress server-minted drift.
// Pass plan as sourceModel on Create, prior state on Read/Update.
// Groups by (stepName, field); steps or fields absent from sourceModel are silently skipped.
func applyIgnoreJSONChanges(newModel *workflowModel, sourceModel *workflowModel, ignorePaths types.List) diag.Diagnostics {
	if ignorePaths.IsNull() || ignorePaths.IsUnknown() || len(ignorePaths.Elements()) == 0 {
		return nil
	}

	type groupKey struct{ step, field string }
	grouped := make(map[groupKey][]string)

	for _, elem := range ignorePaths.Elements() {
		pathStr, ok := elem.(types.String)
		if !ok || pathStr.IsNull() || pathStr.IsUnknown() {
			continue
		}
		stepName, field, inJSONPath, err := resolveIgnoreJSONPath(pathStr.ValueString())
		if err != nil {
			continue // ValidateConfig already rejected malformed paths
		}
		k := groupKey{stepName, field}
		grouped[k] = append(grouped[k], "$."+inJSONPath)
	}

	if len(grouped) == 0 {
		return nil
	}

	if newModel.Definition.IsNull() || newModel.Definition.IsUnknown() {
		return nil
	}
	if sourceModel.Definition.IsNull() || sourceModel.Definition.IsUnknown() {
		return nil
	}

	newDefAttrs := newModel.Definition.Attributes()
	newStepsMap, ok := newDefAttrs["steps"].(types.Map)
	if !ok || newStepsMap.IsNull() || newStepsMap.IsUnknown() {
		return nil
	}

	sourceDefAttrs := sourceModel.Definition.Attributes()
	sourceStepsMap, ok := sourceDefAttrs["steps"].(types.Map)
	if !ok || sourceStepsMap.IsNull() || sourceStepsMap.IsUnknown() {
		return nil
	}

	newStepElems := newStepsMap.Elements()
	sourceStepElems := sourceStepsMap.Elements()

	modified := false

	for k, paths := range grouped {
		sourceStepVal, ok := sourceStepElems[k.step]
		if !ok {
			continue // step absent from source (e.g. brand-new step on create)
		}
		newStepVal, ok := newStepElems[k.step]
		if !ok {
			continue
		}

		sourceStep, ok := sourceStepVal.(types.Object)
		if !ok {
			continue
		}
		newStep, ok := newStepVal.(types.Object)
		if !ok {
			continue
		}

		sourceStepAttrs := sourceStep.Attributes()
		newStepAttrs := newStep.Attributes()

		sourceFieldVal, ok := sourceStepAttrs[k.field].(jsontypes.Normalized)
		if !ok || sourceFieldVal.IsNull() || sourceFieldVal.IsUnknown() {
			continue
		}

		newFieldVal, ok := newStepAttrs[k.field].(jsontypes.Normalized)
		if !ok || newFieldVal.IsNull() || newFieldVal.IsUnknown() {
			continue
		}

		merged, mergeErr := jsonpath.PreservePathsInJSON(newFieldVal.ValueString(), sourceFieldVal.ValueString(), paths)
		if mergeErr != nil || merged == newFieldVal.ValueString() {
			continue
		}

		newStepAttrs[k.field] = jsontypes.NewNormalizedValue(merged)

		newStepObj, diags := types.ObjectValue(stepAttrTypes, newStepAttrs)
		if diags.HasError() {
			return diags
		}
		newStepElems[k.step] = newStepObj
		modified = true
	}

	if !modified {
		return nil
	}

	newStepsUpdated, diags := types.MapValue(types.ObjectType{AttrTypes: stepAttrTypes}, newStepElems)
	if diags.HasError() {
		return diags
	}

	newDefAttrs["steps"] = newStepsUpdated
	newDefObj, diags := types.ObjectValue(definitionAttrTypes, newDefAttrs)
	if diags.HasError() {
		return diags
	}

	newModel.Definition = newDefObj
	return nil
}
