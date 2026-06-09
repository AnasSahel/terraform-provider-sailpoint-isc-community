// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// stringer is satisfied by types.String and other string-backed attr.Value types.
type stringer interface {
	ValueString() string
}

// creatorTftype returns the tftypes.Object shape used by owner/creator/modified_by.
func creatorTftype() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type": tftypes.String,
			"id":   tftypes.String,
			"name": tftypes.String,
		},
	}
}

// definitionV0V1Tftype returns the tftypes.Object shape for v0/v1 definition
// (steps is a plain string at the wire level in both versions).
func definitionV0V1Tftype() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"start": tftypes.String,
			"steps": tftypes.String,
		},
	}
}

// getWorkflowResourceSchema invokes the resource Schema method and returns the
// current (v2) schema. Used to initialise resp.State in upgrade tests.
func getWorkflowResourceSchema(ctx context.Context) schema.Schema {
	r := &workflowResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resp)
	return resp.Schema
}

// mustMapStepAttr is a test helper that extracts a named step from a types.Map
// returned in definition.steps and returns its attributes.
func mustMapStepAttr(t *testing.T, stepsVal attr.Value, stepName string) map[string]attr.Value {
	t.Helper()
	stepsMap, ok := stepsVal.(types.Map)
	if !ok {
		t.Fatalf("definition.steps: want types.Map, got %T", stepsVal)
	}
	stepElem, ok := stepsMap.Elements()[stepName]
	if !ok {
		t.Fatalf("step %q not found in steps map", stepName)
	}
	stepObj, ok := stepElem.(types.Object)
	if !ok {
		t.Fatalf("step %q is not a types.Object: %T", stepName, stepElem)
	}
	return stepObj.Attributes()
}

// runUpgraderTest is the shared test driver for both the v0→v2 and v1→v2 upgraders.
func runUpgraderTest(
	t *testing.T,
	priorSchema schema.Schema,
	raw tftypes.Value,
	upgraderFn func(context.Context, resource.UpgradeStateRequest, *resource.UpgradeStateResponse),
	wantID, wantName string,
	wantDefinitionNull bool,
	wantStart, wantStepName, wantStepType, wantActionID, wantConfigKey string,
) {
	t.Helper()
	ctx := context.Background()
	v2Schema := getWorkflowResourceSchema(ctx)

	req := resource.UpgradeStateRequest{
		State: &tfsdk.State{
			Schema: priorSchema,
			Raw:    raw,
		},
	}
	resp := &resource.UpgradeStateResponse{
		State: tfsdk.State{
			Schema: v2Schema,
		},
	}

	upgraderFn(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("upgrader returned diagnostics: %v", resp.Diagnostics)
	}

	var got workflowModel
	if diags := resp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("failed to decode upgraded state: %v", diags)
	}

	if got.ID.ValueString() != wantID {
		t.Errorf("ID: want %q, got %q", wantID, got.ID.ValueString())
	}
	if got.Name.ValueString() != wantName {
		t.Errorf("Name: want %q, got %q", wantName, got.Name.ValueString())
	}

	if wantDefinitionNull {
		if !got.Definition.IsNull() {
			t.Error("expected Definition to be null, got non-null")
		}
		return
	}

	defAttrs := got.Definition.Attributes()
	if defAttrs == nil {
		t.Fatal("Definition attributes nil")
	}

	startVal, ok := defAttrs["start"].(stringer)
	if !ok {
		t.Fatalf("Definition.start: unexpected type %T", defAttrs["start"])
	}
	if startVal.ValueString() != wantStart {
		t.Errorf("Definition.start: want %q, got %q", wantStart, startVal.ValueString())
	}

	// Verify steps is a types.Map (v2), not a string (v0/v1).
	stepsVal, exists := defAttrs["steps"]
	if !exists {
		t.Fatal("Definition.steps: key missing")
	}
	if _, isStr := stepsVal.(types.String); isStr {
		t.Fatal("Definition.steps is still types.String — upgrader did not convert to typed map")
	}
	if _, isStr := stepsVal.(jsontypes.Normalized); isStr {
		t.Fatal("Definition.steps is still jsontypes.Normalized — upgrader did not convert to typed map")
	}

	if wantStepName == "" {
		return
	}

	stepAttrs := mustMapStepAttr(t, stepsVal, wantStepName)

	if wantStepType != "" {
		tv, ok := stepAttrs["type"].(stringer)
		if !ok {
			t.Fatalf("step[%q].type: unexpected type %T", wantStepName, stepAttrs["type"])
		}
		if tv.ValueString() != wantStepType {
			t.Errorf("step[%q].type: want %q, got %q", wantStepName, wantStepType, tv.ValueString())
		}
	}

	if wantActionID != "" {
		av, ok := stepAttrs["action_id"].(stringer)
		if !ok {
			t.Fatalf("step[%q].action_id: unexpected type %T", wantStepName, stepAttrs["action_id"])
		}
		if av.ValueString() != wantActionID {
			t.Errorf("step[%q].action_id: want %q, got %q", wantStepName, wantActionID, av.ValueString())
		}
	}

	if wantConfigKey != "" {
		cv, ok := stepAttrs["config"].(jsontypes.Normalized)
		if !ok {
			t.Fatalf("step[%q].config: unexpected type %T", wantStepName, stepAttrs["config"])
		}
		if cv.IsNull() || cv.ValueString() == "" {
			t.Errorf("step[%q].config: expected non-null JSON containing %q", wantStepName, wantConfigKey)
		}
	}
}

// TestUpgradeWorkflowStateV0ToV2 verifies the v0→v2 state upgrader.
func TestUpgradeWorkflowStateV0ToV2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	v0Schema := *priorWorkflowSchemaV0()

	stepsJSON := `{"step1":{"actionId":"sp:operator-success","type":"success","nextStep":"end"}}`

	t.Run("workflow with definition and action step", func(t *testing.T) {
		t.Parallel()
		raw := tftypes.NewValue(v0Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":              tftypes.NewValue(tftypes.String, "wf-v0"),
			"name":            tftypes.NewValue(tftypes.String, "v0-workflow"),
			"description":     tftypes.NewValue(tftypes.String, nil),
			"enabled":         tftypes.NewValue(tftypes.Bool, false),
			"trigger":         tftypes.NewValue(tftypes.String, nil),
			"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"execution_count": tftypes.NewValue(tftypes.Number, 0),
			"failure_count":   tftypes.NewValue(tftypes.Number, 0),
			"creator":         tftypes.NewValue(creatorTftype(), nil),
			"modified_by":     tftypes.NewValue(creatorTftype(), nil),
			"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
				"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
				"id":   tftypes.NewValue(tftypes.String, "owner-v0"),
				"name": tftypes.NewValue(tftypes.String, nil),
			}),
			"definition": tftypes.NewValue(definitionV0V1Tftype(), map[string]tftypes.Value{
				"start": tftypes.NewValue(tftypes.String, "step1"),
				"steps": tftypes.NewValue(tftypes.String, stepsJSON),
			}),
		})

		runUpgraderTest(t, v0Schema, raw, upgradeWorkflowStateV0ToV2,
			"wf-v0", "v0-workflow",
			false, "step1", "step1", "success", "sp:operator-success", "",
		)
	})

	t.Run("workflow without definition", func(t *testing.T) {
		t.Parallel()
		raw := tftypes.NewValue(v0Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":              tftypes.NewValue(tftypes.String, "wf-v0-nodef"),
			"name":            tftypes.NewValue(tftypes.String, "no-def"),
			"description":     tftypes.NewValue(tftypes.String, nil),
			"enabled":         tftypes.NewValue(tftypes.Bool, false),
			"trigger":         tftypes.NewValue(tftypes.String, nil),
			"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"execution_count": tftypes.NewValue(tftypes.Number, 0),
			"failure_count":   tftypes.NewValue(tftypes.Number, 0),
			"creator":         tftypes.NewValue(creatorTftype(), nil),
			"modified_by":     tftypes.NewValue(creatorTftype(), nil),
			"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
				"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
				"id":   tftypes.NewValue(tftypes.String, "owner-nodef"),
				"name": tftypes.NewValue(tftypes.String, nil),
			}),
			"definition": tftypes.NewValue(definitionV0V1Tftype(), nil),
		})

		runUpgraderTest(t, v0Schema, raw, upgradeWorkflowStateV0ToV2,
			"wf-v0-nodef", "no-def",
			true, "", "", "", "", "",
		)
	})
}

// TestUpgradeWorkflowStateV1ToV2 verifies the v1→v2 state upgrader.
// In v1, definition.steps was workflowStepsType — wire-compatible with jsontypes.NormalizedType.
func TestUpgradeWorkflowStateV1ToV2(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	v1Schema := *priorWorkflowSchemaV1()

	t.Run("action step round-trips cleanly", func(t *testing.T) {
		t.Parallel()
		stepsJSON := `{"httpStep":{"type":"action","actionId":"sp:http","attributes":{"url":"https://api.example.com"},"nextStep":"end"}}`
		raw := tftypes.NewValue(v1Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":              tftypes.NewValue(tftypes.String, "wf-v1"),
			"name":            tftypes.NewValue(tftypes.String, "v1-workflow"),
			"description":     tftypes.NewValue(tftypes.String, nil),
			"enabled":         tftypes.NewValue(tftypes.Bool, false),
			"trigger":         tftypes.NewValue(tftypes.String, nil),
			"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"execution_count": tftypes.NewValue(tftypes.Number, 0),
			"failure_count":   tftypes.NewValue(tftypes.Number, 0),
			"creator":         tftypes.NewValue(creatorTftype(), nil),
			"modified_by":     tftypes.NewValue(creatorTftype(), nil),
			"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
				"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
				"id":   tftypes.NewValue(tftypes.String, "owner-v1"),
				"name": tftypes.NewValue(tftypes.String, nil),
			}),
			"definition": tftypes.NewValue(definitionV0V1Tftype(), map[string]tftypes.Value{
				"start": tftypes.NewValue(tftypes.String, "httpStep"),
				"steps": tftypes.NewValue(tftypes.String, stepsJSON),
			}),
		})

		runUpgraderTest(t, v1Schema, raw, upgradeWorkflowStateV1ToV2,
			"wf-v1", "v1-workflow",
			false, "httpStep", "httpStep", "action", "sp:http", "",
		)
	})

	t.Run("choice step config catch-all", func(t *testing.T) {
		t.Parallel()
		stepsJSON := `{"decide":{"type":"choice","choiceList":[{"comparator":"StringEquals","nextStep":"yes","variableA.$":"$.x","variableB":"ok"}],"defaultStep":"no"}}`
		raw := tftypes.NewValue(v1Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":              tftypes.NewValue(tftypes.String, "wf-v1-choice"),
			"name":            tftypes.NewValue(tftypes.String, "choice-workflow"),
			"description":     tftypes.NewValue(tftypes.String, nil),
			"enabled":         tftypes.NewValue(tftypes.Bool, false),
			"trigger":         tftypes.NewValue(tftypes.String, nil),
			"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"execution_count": tftypes.NewValue(tftypes.Number, 0),
			"failure_count":   tftypes.NewValue(tftypes.Number, 0),
			"creator":         tftypes.NewValue(creatorTftype(), nil),
			"modified_by":     tftypes.NewValue(creatorTftype(), nil),
			"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
				"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
				"id":   tftypes.NewValue(tftypes.String, "owner-choice"),
				"name": tftypes.NewValue(tftypes.String, nil),
			}),
			"definition": tftypes.NewValue(definitionV0V1Tftype(), map[string]tftypes.Value{
				"start": tftypes.NewValue(tftypes.String, "decide"),
				"steps": tftypes.NewValue(tftypes.String, stepsJSON),
			}),
		})

		runUpgraderTest(t, v1Schema, raw, upgradeWorkflowStateV1ToV2,
			"wf-v1-choice", "choice-workflow",
			false, "decide", "decide", "choice", "", "choiceList",
		)
	})

	t.Run("workflow without definition", func(t *testing.T) {
		t.Parallel()
		raw := tftypes.NewValue(v1Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
			"id":              tftypes.NewValue(tftypes.String, "wf-v1-nodef"),
			"name":            tftypes.NewValue(tftypes.String, "no-def"),
			"description":     tftypes.NewValue(tftypes.String, nil),
			"enabled":         tftypes.NewValue(tftypes.Bool, false),
			"trigger":         tftypes.NewValue(tftypes.String, nil),
			"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
			"execution_count": tftypes.NewValue(tftypes.Number, 0),
			"failure_count":   tftypes.NewValue(tftypes.Number, 0),
			"creator":         tftypes.NewValue(creatorTftype(), nil),
			"modified_by":     tftypes.NewValue(creatorTftype(), nil),
			"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
				"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
				"id":   tftypes.NewValue(tftypes.String, "owner-nodef"),
				"name": tftypes.NewValue(tftypes.String, nil),
			}),
			"definition": tftypes.NewValue(definitionV0V1Tftype(), nil),
		})

		runUpgraderTest(t, v1Schema, raw, upgradeWorkflowStateV1ToV2,
			"wf-v1-nodef", "no-def",
			true, "", "", "", "", "",
		)
	})
}
