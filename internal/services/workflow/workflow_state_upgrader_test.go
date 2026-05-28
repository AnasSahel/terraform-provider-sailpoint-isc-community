// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestUpgradeWorkflowStateV0ToV1 verifies the v0 → v1 state upgrade for a
// realistic workflow state: same wire data, only the framework type
// identity of definition.steps changes. Closes #114.
func TestUpgradeWorkflowStateV0ToV1(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v0Schema := *priorWorkflowSchemaV0()
	v1Schema := getWorkflowResourceSchema(ctx)

	stepsJSON := `{"step1":{"actionId":"sp:operator-success","type":"success"}}`

	tests := map[string]struct {
		raw                tftypes.Value
		wantID             string
		wantName           string
		wantDefinitionNull bool
		wantStepsString    string
		wantStartString    string
		wantEnabled        bool
		wantHasOwner       bool
		wantOwnerType      string
		wantOwnerID        string
	}{
		"workflow with definition and owner": {
			raw: tftypes.NewValue(v0Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"id":              tftypes.NewValue(tftypes.String, "wf-123"),
				"name":            tftypes.NewValue(tftypes.String, "example-workflow"),
				"description":     tftypes.NewValue(tftypes.String, "test"),
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
					"id":   tftypes.NewValue(tftypes.String, "owner-id"),
					"name": tftypes.NewValue(tftypes.String, "Owner Name"),
				}),
				"definition": tftypes.NewValue(definitionV0Tftype(), map[string]tftypes.Value{
					"start": tftypes.NewValue(tftypes.String, "step1"),
					"steps": tftypes.NewValue(tftypes.String, stepsJSON),
				}),
			}),
			wantID:          "wf-123",
			wantName:        "example-workflow",
			wantStepsString: stepsJSON,
			wantStartString: "step1",
			wantEnabled:     false,
			wantHasOwner:    true,
			wantOwnerType:   "IDENTITY",
			wantOwnerID:     "owner-id",
		},
		"workflow without definition": {
			raw: tftypes.NewValue(v0Schema.Type().TerraformType(ctx), map[string]tftypes.Value{
				"id":              tftypes.NewValue(tftypes.String, "wf-456"),
				"name":            tftypes.NewValue(tftypes.String, "no-def"),
				"description":     tftypes.NewValue(tftypes.String, nil),
				"enabled":         tftypes.NewValue(tftypes.Bool, true),
				"trigger":         tftypes.NewValue(tftypes.String, nil),
				"created":         tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
				"modified":        tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z"),
				"execution_count": tftypes.NewValue(tftypes.Number, 0),
				"failure_count":   tftypes.NewValue(tftypes.Number, 0),
				"creator":         tftypes.NewValue(creatorTftype(), nil),
				"modified_by":     tftypes.NewValue(creatorTftype(), nil),
				"owner": tftypes.NewValue(creatorTftype(), map[string]tftypes.Value{
					"type": tftypes.NewValue(tftypes.String, "IDENTITY"),
					"id":   tftypes.NewValue(tftypes.String, "owner-2"),
					"name": tftypes.NewValue(tftypes.String, nil),
				}),
				"definition": tftypes.NewValue(definitionV0Tftype(), nil),
			}),
			wantID:             "wf-456",
			wantName:           "no-def",
			wantDefinitionNull: true,
			wantEnabled:        true,
			wantHasOwner:       true,
			wantOwnerType:      "IDENTITY",
			wantOwnerID:        "owner-2",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := resource.UpgradeStateRequest{
				State: &tfsdk.State{
					Schema: v0Schema,
					Raw:    tc.raw,
				},
			}
			resp := &resource.UpgradeStateResponse{
				State: tfsdk.State{
					Schema: v1Schema,
				},
			}

			upgradeWorkflowStateV0ToV1(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("upgrader returned diagnostics: %v", resp.Diagnostics)
			}

			var got workflowModel
			diags := resp.State.Get(ctx, &got)
			if diags.HasError() {
				t.Fatalf("failed to decode upgraded state: %v", diags)
			}

			if got.ID.ValueString() != tc.wantID {
				t.Errorf("ID: want %q, got %q", tc.wantID, got.ID.ValueString())
			}
			if got.Name.ValueString() != tc.wantName {
				t.Errorf("Name: want %q, got %q", tc.wantName, got.Name.ValueString())
			}
			if got.Enabled.ValueBool() != tc.wantEnabled {
				t.Errorf("Enabled: want %v, got %v", tc.wantEnabled, got.Enabled.ValueBool())
			}

			if tc.wantHasOwner {
				ownerAttrs := got.Owner.Attributes()
				if ownerAttrs == nil {
					t.Fatalf("Owner attributes nil")
				}
				if ownerType, ok := ownerAttrs["type"].(stringer); !ok {
					t.Errorf("Owner.type: not a string-like value, got %T", ownerAttrs["type"])
				} else if ownerType.ValueString() != tc.wantOwnerType {
					t.Errorf("Owner.type: want %q, got %q", tc.wantOwnerType, ownerType.ValueString())
				}
				if ownerID, ok := ownerAttrs["id"].(stringer); !ok {
					t.Errorf("Owner.id: not a string-like value")
				} else if ownerID.ValueString() != tc.wantOwnerID {
					t.Errorf("Owner.id: want %q, got %q", tc.wantOwnerID, ownerID.ValueString())
				}
			}

			if tc.wantDefinitionNull {
				if !got.Definition.IsNull() {
					t.Errorf("expected Definition to be null")
				}
				return
			}

			defAttrs := got.Definition.Attributes()
			if defAttrs == nil {
				t.Fatalf("Definition attributes nil")
			}
			startVal, ok := defAttrs["start"].(stringer)
			if !ok {
				t.Fatalf("Definition.start: not a string-like value")
			}
			if startVal.ValueString() != tc.wantStartString {
				t.Errorf("Definition.start: want %q, got %q", tc.wantStartString, startVal.ValueString())
			}

			// Critical assertion: steps must be wrapped in workflowStepsValue
			// (not jsontypes.Normalized) — that is the whole point of v1.
			stepsVal, ok := defAttrs["steps"].(workflowStepsValue)
			if !ok {
				t.Fatalf("Definition.steps: want workflowStepsValue, got %T", defAttrs["steps"])
			}
			if stepsVal.ValueString() != tc.wantStepsString {
				t.Errorf("Definition.steps content: want %q, got %q", tc.wantStepsString, stepsVal.ValueString())
			}
		})
	}
}

// stringer is satisfied by both types.String and any custom string type
// the framework wraps under attr.Value, so the test stays type-agnostic for
// fields that may evolve.
type stringer interface {
	ValueString() string
}

// creatorTftype returns the tftypes.Object shape used by owner/creator/
// modified_by in the v0 schema.
func creatorTftype() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type": tftypes.String,
			"id":   tftypes.String,
			"name": tftypes.String,
		},
	}
}

// definitionV0Tftype returns the tftypes.Object shape for definition in v0.
// Steps is a string at the wire level (jsontypes.NormalizedType is just a
// string with semantic equality).
func definitionV0Tftype() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"start": tftypes.String,
			"steps": tftypes.String,
		},
	}
}

// getWorkflowResourceSchema invokes the resource's Schema method and
// returns the resulting schema (v1). Used by the upgrade test to set up
// resp.State with the v1 schema as the framework would.
func getWorkflowResourceSchema(ctx context.Context) schema.Schema {
	r := &workflowResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resp)
	return resp.Schema
}
