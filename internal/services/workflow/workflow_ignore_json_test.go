// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mustListValue builds a types.List of strings for tests.
func mustListValue(t *testing.T, elems []string) types.List {
	t.Helper()
	vals := make([]attr.Value, len(elems))
	for i, s := range elems {
		vals[i] = types.StringValue(s)
	}
	list, diags := types.ListValue(types.StringType, vals)
	if diags.HasError() {
		t.Fatalf("mustListValue: %v", diags)
	}
	return list
}

func TestResolveIgnoreJSONPath_valid(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path           string
		wantStep       string
		wantField      string
		wantInJSONPath string
	}{
		{
			path:           "definition.steps['Create Ticket'].attributes.param_oauth.refID",
			wantStep:       "Create Ticket",
			wantField:      "attributes",
			wantInJSONPath: "param_oauth.refID",
		},
		{
			path:           "definition.steps['step1'].config.someKey",
			wantStep:       "step1",
			wantField:      "config",
			wantInJSONPath: "someKey",
		},
		{
			path:           "definition.steps[\"my step\"].catch.error.code",
			wantStep:       "my step",
			wantField:      "catch",
			wantInJSONPath: "error.code",
		},
		{
			path:           "definition.steps['s'].attributes.items[*].refID",
			wantStep:       "s",
			wantField:      "attributes",
			wantInJSONPath: "items[*].refID",
		},
		{
			path:           "definition.steps['s'].attributes.nested.deep.key",
			wantStep:       "s",
			wantField:      "attributes",
			wantInJSONPath: "nested.deep.key",
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			stepName, field, inJSONPath, err := resolveIgnoreJSONPath(tc.path)
			if err != nil {
				t.Fatalf("resolveIgnoreJSONPath(%q) unexpected error: %v", tc.path, err)
			}
			if stepName != tc.wantStep {
				t.Errorf("stepName = %q, want %q", stepName, tc.wantStep)
			}
			if field != tc.wantField {
				t.Errorf("field = %q, want %q", field, tc.wantField)
			}
			if inJSONPath != tc.wantInJSONPath {
				t.Errorf("inJSONPath = %q, want %q", inJSONPath, tc.wantInJSONPath)
			}
		})
	}
}

func TestResolveIgnoreJSONPath_invalid(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path string
	}{
		{path: "steps['step'].attributes.key"},           // missing "definition." prefix
		{path: "definition.steps.step.attributes.key"},   // step not bracket-quoted
		{path: "definition.steps[''].attributes.key"},    // empty step name
		{path: "definition.steps['step'].attributes"},    // no in-JSON path (missing dot + path)
		{path: "definition.steps['step'].name.key"},      // "name" is not a JSON step field
		{path: "definition.steps['step'].attributes."},   // trailing dot in in-JSON path
		{path: "definition.steps['step'].owner.somekey"}, // "owner" is not a JSON step field
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			_, _, _, err := resolveIgnoreJSONPath(tc.path)
			if err == nil {
				t.Errorf("resolveIgnoreJSONPath(%q) expected error, got nil", tc.path)
			}
		})
	}
}

// TestApplyIgnoreJSONChanges_DriftSuppression_Create verifies that after a Create,
// the server-minted refID is replaced by the practitioner's placeholder value.
func TestApplyIgnoreJSONChanges_DriftSuppression_Create(t *testing.T) {
	t.Parallel()

	// Plan: param_oauth.refID = "placeholder" (practitioner's input)
	planStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com","param_oauth":{"refID":"placeholder"}}`),
	})
	planModel := workflowModel{
		Name:           types.StringValue("test"),
		Definition:     mustDefinitionObject(t, "Create Ticket", mustStepsMap(t, map[string]types.Object{"Create Ticket": planStep})),
		Owner:          types.ObjectNull(objectRefAttrTypes()),
		Trigger:        jsontypes.NewNormalizedNull(),
		Created:        types.StringNull(),
		Modified:       types.StringNull(),
		Creator:        types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:     types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount: types.Int32Value(0),
		FailureCount:   types.Int32Value(0),
		IgnoreJSONChanges: mustListValue(t, []string{
			"definition.steps['Create Ticket'].attributes.param_oauth.refID",
		}),
	}

	// Simulate API response: server minted refID to a different value
	apiStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com","param_oauth":{"refID":"server-minted-id"}}`),
	})
	newModel := workflowModel{
		Name:              types.StringValue("test"),
		Definition:        mustDefinitionObject(t, "Create Ticket", mustStepsMap(t, map[string]types.Object{"Create Ticket": apiStep})),
		Owner:             types.ObjectNull(objectRefAttrTypes()),
		Trigger:           jsontypes.NewNormalizedNull(),
		Created:           types.StringNull(),
		Modified:          types.StringNull(),
		Creator:           types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:        types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount:    types.Int32Value(0),
		FailureCount:      types.Int32Value(0),
		IgnoreJSONChanges: planModel.IgnoreJSONChanges,
	}

	// Apply (Create flow: source = plan)
	diags := applyIgnoreJSONChanges(&newModel, &planModel, planModel.IgnoreJSONChanges)
	if diags.HasError() {
		t.Fatalf("applyIgnoreJSONChanges: %v", diags)
	}

	// Verify refID kept the practitioner's value, not the server's
	refID := mustExtractRefID(t, newModel, "Create Ticket")
	if refID != "placeholder" {
		t.Errorf("refID = %q, want %q (server drift not suppressed)", refID, "placeholder")
	}
}

// TestApplyIgnoreJSONChanges_DriftSuppression_Read verifies that on subsequent reads
// the prior-state value is preserved at ignored paths.
func TestApplyIgnoreJSONChanges_DriftSuppression_Read(t *testing.T) {
	t.Parallel()

	// Prior state: param_oauth.refID = "placeholder"
	priorStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com","param_oauth":{"refID":"placeholder"}}`),
	})
	priorState := workflowModel{
		Name:           types.StringValue("test"),
		Definition:     mustDefinitionObject(t, "Create Ticket", mustStepsMap(t, map[string]types.Object{"Create Ticket": priorStep})),
		Owner:          types.ObjectNull(objectRefAttrTypes()),
		Trigger:        jsontypes.NewNormalizedNull(),
		Created:        types.StringNull(),
		Modified:       types.StringNull(),
		Creator:        types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:     types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount: types.Int32Value(0),
		FailureCount:   types.Int32Value(0),
		IgnoreJSONChanges: mustListValue(t, []string{
			"definition.steps['Create Ticket'].attributes.param_oauth.refID",
		}),
	}

	// API returns the server-minted value again
	apiStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com","param_oauth":{"refID":"server-minted-id"}}`),
	})
	newModel := workflowModel{
		Name:              types.StringValue("test"),
		Definition:        mustDefinitionObject(t, "Create Ticket", mustStepsMap(t, map[string]types.Object{"Create Ticket": apiStep})),
		Owner:             types.ObjectNull(objectRefAttrTypes()),
		Trigger:           jsontypes.NewNormalizedNull(),
		Created:           types.StringNull(),
		Modified:          types.StringNull(),
		Creator:           types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:        types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount:    types.Int32Value(0),
		FailureCount:      types.Int32Value(0),
		IgnoreJSONChanges: priorState.IgnoreJSONChanges,
	}

	// Apply (Read flow: source = prior state)
	diags := applyIgnoreJSONChanges(&newModel, &priorState, priorState.IgnoreJSONChanges)
	if diags.HasError() {
		t.Fatalf("applyIgnoreJSONChanges: %v", diags)
	}

	refID := mustExtractRefID(t, newModel, "Create Ticket")
	if refID != "placeholder" {
		t.Errorf("refID = %q, want %q (server drift not suppressed on read)", refID, "placeholder")
	}
}

// TestApplyIgnoreJSONChanges_NoIgnorePaths verifies that a nil/null ignore list is a no-op.
func TestApplyIgnoreJSONChanges_NoIgnorePaths(t *testing.T) {
	t.Parallel()

	apiStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com","param_oauth":{"refID":"server-minted-id"}}`),
	})
	model := workflowModel{
		Name:              types.StringValue("test"),
		Definition:        mustDefinitionObject(t, "s", mustStepsMap(t, map[string]types.Object{"s": apiStep})),
		Owner:             types.ObjectNull(objectRefAttrTypes()),
		Trigger:           jsontypes.NewNormalizedNull(),
		Created:           types.StringNull(),
		Modified:          types.StringNull(),
		Creator:           types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:        types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount:    types.Int32Value(0),
		FailureCount:      types.Int32Value(0),
		IgnoreJSONChanges: types.ListNull(types.StringType),
	}
	source := model

	diags := applyIgnoreJSONChanges(&model, &source, types.ListNull(types.StringType))
	if diags.HasError() {
		t.Fatalf("applyIgnoreJSONChanges: %v", diags)
	}
	// No change expected — server-minted value stays
	refID := mustExtractRefID(t, model, "s")
	if refID != "server-minted-id" {
		t.Errorf("refID = %q, want %q", refID, "server-minted-id")
	}
}

// TestApplyIgnoreJSONChanges_MissingSourceStep verifies that when the step is absent
// from the source (e.g. a brand-new step on create), the path is silently skipped.
func TestApplyIgnoreJSONChanges_MissingSourceStep(t *testing.T) {
	t.Parallel()

	// Source has no "New Step"
	sourceModel := workflowModel{
		Name:              types.StringValue("test"),
		Definition:        mustDefinitionObject(t, "Old Step", mustStepsMap(t, map[string]types.Object{})),
		Owner:             types.ObjectNull(objectRefAttrTypes()),
		Trigger:           jsontypes.NewNormalizedNull(),
		Created:           types.StringNull(),
		Modified:          types.StringNull(),
		Creator:           types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:        types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount:    types.Int32Value(0),
		FailureCount:      types.Int32Value(0),
		IgnoreJSONChanges: types.ListNull(types.StringType),
	}

	apiStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"attributes": jsontypes.NewNormalizedValue(`{"param_oauth":{"refID":"server-minted-id"}}`),
	})
	newModel := workflowModel{
		Name:              types.StringValue("test"),
		Definition:        mustDefinitionObject(t, "New Step", mustStepsMap(t, map[string]types.Object{"New Step": apiStep})),
		Owner:             types.ObjectNull(objectRefAttrTypes()),
		Trigger:           jsontypes.NewNormalizedNull(),
		Created:           types.StringNull(),
		Modified:          types.StringNull(),
		Creator:           types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:        types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount:    types.Int32Value(0),
		FailureCount:      types.Int32Value(0),
		IgnoreJSONChanges: types.ListNull(types.StringType),
	}

	paths := mustListValue(t, []string{"definition.steps['New Step'].attributes.param_oauth.refID"})
	diags := applyIgnoreJSONChanges(&newModel, &sourceModel, paths)
	if diags.HasError() {
		t.Fatalf("applyIgnoreJSONChanges: %v", diags)
	}
	// Server-minted value should remain unchanged (step absent from source → skip)
	refID := mustExtractRefID(t, newModel, "New Step")
	if refID != "server-minted-id" {
		t.Errorf("refID = %q, want %q (should be unchanged when step absent from source)", refID, "server-minted-id")
	}
}

// mustExtractRefID is a test helper that navigates newModel.definition.steps[stepName].attributes
// and returns param_oauth.refID as a string. Fails the test if any step is missing.
func mustExtractRefID(t *testing.T, m workflowModel, stepName string) string {
	t.Helper()
	defAttrs := m.Definition.Attributes()
	stepsMap, ok := defAttrs["steps"].(types.Map)
	if !ok {
		t.Fatalf("mustExtractRefID: steps is not types.Map")
	}
	stepVal, ok := stepsMap.Elements()[stepName].(types.Object)
	if !ok {
		t.Fatalf("mustExtractRefID: step %q missing", stepName)
	}
	attrVal, ok := stepVal.Attributes()["attributes"].(jsontypes.Normalized)
	if !ok || attrVal.IsNull() {
		t.Fatalf("mustExtractRefID: attributes is null or wrong type")
	}
	var attrMap map[string]any
	if err := json.Unmarshal([]byte(attrVal.ValueString()), &attrMap); err != nil {
		t.Fatalf("mustExtractRefID: unmarshal: %v", err)
	}
	paramOAuth, ok := attrMap["param_oauth"].(map[string]any)
	if !ok {
		t.Fatalf("mustExtractRefID: param_oauth not a map")
	}
	refID, ok := paramOAuth["refID"].(string)
	if !ok {
		t.Fatalf("mustExtractRefID: refID not a string")
	}
	return refID
}
