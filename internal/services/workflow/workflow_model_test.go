// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mustStepObject builds a typed step types.Object for tests.
func mustStepObject(t *testing.T, attrs map[string]attr.Value) types.Object {
	t.Helper()
	full := map[string]attr.Value{
		"type":           types.StringNull(),
		"action_id":      types.StringNull(),
		"attributes":     jsontypes.NewNormalizedNull(),
		"next_step":      types.StringNull(),
		"display_name":   types.StringNull(),
		"version_number": types.Int32Null(),
		"catch":          jsontypes.NewNormalizedNull(),
		"config":         jsontypes.NewNormalizedNull(),
	}
	for k, v := range attrs {
		full[k] = v
	}
	obj, diags := types.ObjectValue(stepAttrTypes, full)
	if diags.HasError() {
		t.Fatalf("mustStepObject: %v", diags)
	}
	return obj
}

// mustStepsMap builds a types.Map of step objects for tests.
func mustStepsMap(t *testing.T, steps map[string]types.Object) types.Map {
	t.Helper()
	elems := make(map[string]attr.Value, len(steps))
	for k, v := range steps {
		elems[k] = v
	}
	m, diags := types.MapValue(types.ObjectType{AttrTypes: stepAttrTypes}, elems)
	if diags.HasError() {
		t.Fatalf("mustStepsMap: %v", diags)
	}
	return m
}

// mustDefinitionObject builds a typed definition types.Object for tests.
func mustDefinitionObject(t *testing.T, start string, steps types.Map) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
		"start": types.StringValue(start),
		"steps": steps,
	})
	if diags.HasError() {
		t.Fatalf("mustDefinitionObject: %v", diags)
	}
	return obj
}

// TestWorkflowModel_ToAPI_FromAPI_RoundTrip_ActionStep verifies a single action step survives
// a full ToAPI → FromAPI round-trip.
func TestWorkflowModel_ToAPI_FromAPI_RoundTrip_ActionStep(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	actionStep := mustStepObject(t, map[string]attr.Value{
		"type":       types.StringValue("action"),
		"action_id":  types.StringValue("sp:http"),
		"attributes": jsontypes.NewNormalizedValue(`{"url":"https://example.com"}`),
		"next_step":  types.StringValue("end"),
	})
	stepsMap := mustStepsMap(t, map[string]types.Object{"httpStep": actionStep})
	original := workflowModel{
		Name:           types.StringValue("test-workflow"),
		Description:    types.StringValue("test description"),
		Enabled:        types.BoolValue(false),
		Definition:     mustDefinitionObject(t, "httpStep", stepsMap),
		Owner:          types.ObjectNull(objectRefAttrTypes()),
		Trigger:        jsontypes.NewNormalizedNull(),
		Created:        types.StringNull(),
		Modified:       types.StringNull(),
		Creator:        types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:     types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount: types.Int32Value(0),
		FailureCount:   types.Int32Value(0),
	}

	api, toAPIDiags := original.ToAPI(ctx)
	if toAPIDiags.HasError() {
		t.Fatalf("ToAPI produced diagnostics: %v", toAPIDiags)
	}
	if api.Definition == nil {
		t.Fatal("ToAPI: definition is nil")
	}
	if api.Definition.Start != "httpStep" {
		t.Errorf("ToAPI: start = %q, want %q", api.Definition.Start, "httpStep")
	}
	step, ok := api.Definition.Steps["httpStep"].(map[string]any)
	if !ok {
		t.Fatalf("ToAPI: httpStep not a map[string]any")
	}
	if step["actionId"] != "sp:http" {
		t.Errorf("ToAPI: actionId = %v, want sp:http", step["actionId"])
	}
	if step["nextStep"] != "end" {
		t.Errorf("ToAPI: nextStep = %v, want end", step["nextStep"])
	}

	apiResponse := client.WorkflowAPI{
		ID:         "test-id",
		Name:       api.Name,
		Definition: api.Definition,
		Enabled:    api.Enabled,
	}

	var roundTripped workflowModel
	fromAPIDiags := roundTripped.FromAPI(ctx, apiResponse)
	if fromAPIDiags.HasError() {
		t.Fatalf("FromAPI produced diagnostics: %v", fromAPIDiags)
	}

	if roundTripped.Definition.IsNull() || roundTripped.Definition.IsUnknown() {
		t.Fatal("FromAPI: definition is null/unknown after round-trip")
	}

	defAttrs := roundTripped.Definition.Attributes()
	rtSteps, ok := defAttrs["steps"].(types.Map)
	if !ok {
		t.Fatalf("FromAPI: steps has type %T, want types.Map", defAttrs["steps"])
	}
	if rtSteps.IsNull() || rtSteps.IsUnknown() {
		t.Fatal("FromAPI: steps is null/unknown")
	}

	rtStep, ok := rtSteps.Elements()["httpStep"].(types.Object)
	if !ok {
		t.Fatalf("FromAPI: httpStep missing from steps map")
	}
	rtAttrs := rtStep.Attributes()

	if v, ok := rtAttrs["action_id"].(types.String); !ok || v.ValueString() != "sp:http" {
		t.Errorf("FromAPI: action_id = %v, want sp:http", rtAttrs["action_id"])
	}
	if v, ok := rtAttrs["next_step"].(types.String); !ok || v.ValueString() != "end" {
		t.Errorf("FromAPI: next_step = %v, want end", rtAttrs["next_step"])
	}
}

// TestWorkflowModel_ToAPI_FromAPI_RoundTrip_ChoiceStep verifies a choice step with a
// config catch-all survives a full ToAPI → FromAPI round-trip.
func TestWorkflowModel_ToAPI_FromAPI_RoundTrip_ChoiceStep(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	choiceStep := mustStepObject(t, map[string]attr.Value{
		"type": types.StringValue("choice"),
		"config": jsontypes.NewNormalizedValue(
			`{"choiceList":[{"comparator":"StringEquals","nextStep":"yes","variableA.$":"$.x","variableB":"ok"}],"defaultStep":"no"}`,
		),
	})
	stepsMap := mustStepsMap(t, map[string]types.Object{"decide": choiceStep})
	original := workflowModel{
		Name:           types.StringValue("choice-wf"),
		Enabled:        types.BoolValue(false),
		Definition:     mustDefinitionObject(t, "decide", stepsMap),
		Owner:          types.ObjectNull(objectRefAttrTypes()),
		Trigger:        jsontypes.NewNormalizedNull(),
		Created:        types.StringNull(),
		Modified:       types.StringNull(),
		Creator:        types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:     types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount: types.Int32Value(0),
		FailureCount:   types.Int32Value(0),
	}

	api, toAPIDiags := original.ToAPI(ctx)
	if toAPIDiags.HasError() {
		t.Fatalf("ToAPI produced diagnostics: %v", toAPIDiags)
	}
	step, ok := api.Definition.Steps["decide"].(map[string]any)
	if !ok {
		t.Fatalf("ToAPI: decide step not a map[string]any")
	}
	if _, hasChoiceList := step["choiceList"]; !hasChoiceList {
		t.Error("ToAPI: choiceList missing from flattened step")
	}
	if _, hasDefault := step["defaultStep"]; !hasDefault {
		t.Error("ToAPI: defaultStep missing from flattened step")
	}

	// Round-trip back via FromAPI
	apiResponse := client.WorkflowAPI{
		ID:         "choice-id",
		Name:       api.Name,
		Definition: api.Definition,
	}

	var roundTripped workflowModel
	if diags := roundTripped.FromAPI(ctx, apiResponse); diags.HasError() {
		t.Fatalf("FromAPI produced diagnostics: %v", diags)
	}

	defAttrs := roundTripped.Definition.Attributes()
	rtSteps, _ := defAttrs["steps"].(types.Map)
	rtStep, ok := rtSteps.Elements()["decide"].(types.Object)
	if !ok {
		t.Fatal("FromAPI: decide step missing")
	}
	rtAttrs := rtStep.Attributes()

	// choiceList and defaultStep must be in config, not in named fields
	if v, ok := rtAttrs["config"].(jsontypes.Normalized); !ok || v.IsNull() {
		t.Error("FromAPI: config should be non-null for choice step")
	}
	if v, ok := rtAttrs["action_id"].(types.String); ok && !v.IsNull() {
		t.Errorf("FromAPI: action_id should be null for choice step, got %v", v.ValueString())
	}
}
