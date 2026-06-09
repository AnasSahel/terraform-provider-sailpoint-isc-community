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

const testStepsJSON = `{"httpStep":{"actionId":"sp:http","attributes":{"url":"https://example.com"}}}`

func mustDefinitionObject(t *testing.T, start, stepsJSON string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(definitionAttrTypes, map[string]attr.Value{
		"start": types.StringValue(start),
		"steps": workflowStepsValue{Normalized: jsontypes.NewNormalizedValue(stepsJSON)},
	})
	if diags.HasError() {
		t.Fatalf("failed to build definition object: %v", diags)
	}
	return obj
}

// TestWorkflowModel_ToAPI_FromAPI_RoundTrip_WithDefinition verifies that a
// workflow with a definition (steps included) survives a full ToAPI → FromAPI
// round-trip without any diagnostics errors and with the steps preserved.
func TestWorkflowModel_ToAPI_FromAPI_RoundTrip_WithDefinition(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	original := workflowModel{
		Name:           types.StringValue("test-workflow"),
		Description:    types.StringValue("test description"),
		Enabled:        types.BoolValue(true),
		Definition:     mustDefinitionObject(t, "httpStep", testStepsJSON),
		Owner:          types.ObjectNull(objectRefAttrTypes()),
		Trigger:        jsontypes.NewNormalizedNull(),
		Created:        types.StringNull(),
		Modified:       types.StringNull(),
		Creator:        types.ObjectNull(objectRefAttrTypes()),
		ModifiedBy:     types.ObjectNull(objectRefAttrTypes()),
		ExecutionCount: types.Int32Value(0),
		FailureCount:   types.Int32Value(0),
	}

	// ToAPI
	api, toAPIDiags := original.ToAPI(ctx)
	if toAPIDiags.HasError() {
		t.Fatalf("ToAPI produced unexpected diagnostics: %v", toAPIDiags)
	}
	if api.Definition == nil {
		t.Fatal("ToAPI: definition is nil — steps were dropped")
	}
	if api.Definition.Start != "httpStep" {
		t.Errorf("ToAPI: definition.start = %q, want %q", api.Definition.Start, "httpStep")
	}
	if len(api.Definition.Steps) == 0 {
		t.Fatal("ToAPI: definition.steps is empty — steps were silently dropped")
	}

	// Simulate a read-back from the SailPoint API.
	apiResponse := client.WorkflowAPI{
		ID:         "test-id",
		Name:       api.Name,
		Definition: api.Definition,
		Enabled:    api.Enabled,
	}

	// FromAPI
	var roundTripped workflowModel
	fromAPIDiags := roundTripped.FromAPI(ctx, apiResponse)
	if fromAPIDiags.HasError() {
		t.Fatalf("FromAPI produced unexpected diagnostics: %v", fromAPIDiags)
	}

	// Verify the definition object was reconstructed.
	if roundTripped.Definition.IsNull() || roundTripped.Definition.IsUnknown() {
		t.Fatal("FromAPI: definition is null/unknown after round-trip")
	}

	// Verify the steps attribute is the correct custom value type.
	rtAttrs := roundTripped.Definition.Attributes()
	rtSteps, ok := rtAttrs["steps"].(workflowStepsValue)
	if !ok {
		t.Fatalf("FromAPI: steps attribute has type %T, want workflowStepsValue", rtAttrs["steps"])
	}

	// Verify the steps JSON content survived the round-trip.
	originalSteps, _ := original.Definition.Attributes()["steps"].(workflowStepsValue)
	eq, eqDiags := originalSteps.StringSemanticEquals(ctx, rtSteps)
	if eqDiags.HasError() {
		t.Fatalf("StringSemanticEquals diagnostics: %v", eqDiags)
	}
	if !eq {
		t.Errorf("steps not semantically equal after round-trip\noriginal : %s\nround-trip: %s",
			originalSteps.ValueString(), rtSteps.ValueString())
	}
}
