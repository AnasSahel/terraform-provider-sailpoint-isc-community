// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ignoredFieldsByActionID lists the JSON paths inside each step's `attributes`
// block that SailPoint mints/normalizes server-side and that the provider must
// not treat as a divergence between plan and state. Without this list, a
// `tofu apply` that creates or updates a workflow with one of these step types
// fails with `Provider produced inconsistent result after apply` on the very
// first apply (the second apply silently resyncs because the state then
// matches the API).
//
// Paths use dotted notation, relative to the step's value (i.e. *not* including
// the step name). The walker traverses nested objects only — no array support
// for now since none of the known minted fields live inside arrays.
//
// Add a new entry whenever a SailPoint action surface starts minting another
// field. Keep entries grouped by SailPoint action id so the source of truth
// stays readable.
var ignoredFieldsByActionID = map[string][]string{
	// `sp:http` Storage Parameter Service refs: SailPoint mints a fresh refID
	// at workflow POST time regardless of what the client sends. paramID and
	// other auth-related metadata are preserved as submitted.
	"sp:http": {
		"attributes.param_oauth.refID",
		"attributes.param_header.refID",
		"attributes.param_oauth_scopes.refID",
	},
}

// workflowStepsType extends jsontypes.NormalizedType with a SemanticEquals
// implementation that ignores server-minted fields per action id. See
// `ignoredFieldsByActionID` for the maintained allow-list and #90 for the
// underlying SailPoint behavior.
type workflowStepsType struct {
	jsontypes.NormalizedType
}

func (t workflowStepsType) String() string {
	return "workflow.workflowStepsType"
}

func (t workflowStepsType) ValueType(_ context.Context) attr.Value {
	return workflowStepsValue{}
}

func (t workflowStepsType) Equal(o attr.Type) bool {
	other, ok := o.(workflowStepsType)
	if !ok {
		return false
	}
	return t.NormalizedType.Equal(other.NormalizedType)
}

func (t workflowStepsType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return workflowStepsValue{Normalized: jsontypes.Normalized{StringValue: in}}, nil
}

func (t workflowStepsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.NormalizedType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	normalized, ok := attrValue.(jsontypes.Normalized)
	if !ok {
		return nil, fmt.Errorf("unexpected value type from NormalizedType.ValueFromTerraform: %T", attrValue)
	}
	return workflowStepsValue{Normalized: normalized}, nil
}

// workflowStepsValue wraps jsontypes.Normalized with the extended
// SemanticEquals logic.
type workflowStepsValue struct {
	jsontypes.Normalized
}

func (v workflowStepsValue) Type(_ context.Context) attr.Type {
	return workflowStepsType{}
}

func (v workflowStepsValue) Equal(o attr.Value) bool {
	other, ok := o.(workflowStepsValue)
	if !ok {
		return false
	}
	return v.Normalized.Equal(other.Normalized)
}

// StringSemanticEquals returns true when the two JSON strings are equal after
// stripping the server-minted fields documented in `ignoredFieldsByActionID`.
// A debug-level log line is emitted on every ignored divergence so a user
// running with `TF_LOG=debug` can audit what the provider is masking.
func (v workflowStepsValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	var newStr string
	switch nv := newValuable.(type) {
	case workflowStepsValue:
		newStr = nv.ValueString()
	case jsontypes.Normalized:
		newStr = nv.ValueString()
	default:
		diags.AddError(
			"Semantic Equality Check Error",
			fmt.Sprintf("workflowStepsValue.StringSemanticEquals received an unexpected value type %T; please report this to the provider developers.", newValuable),
		)
		return false, diags
	}

	oldStripped, oldOriginal, err := stripIgnoredFields(v.ValueString())
	if err != nil {
		diags.AddError("Semantic Equality Check Error", "Failed to parse prior JSON: "+err.Error())
		return false, diags
	}
	newStripped, newOriginal, err := stripIgnoredFields(newStr)
	if err != nil {
		diags.AddError("Semantic Equality Check Error", "Failed to parse new JSON: "+err.Error())
		return false, diags
	}

	if !reflect.DeepEqual(oldStripped, newStripped) {
		return false, diags
	}

	// Equal post-strip. If the original JSONs differ, we masked something —
	// log the masked paths for traceability.
	if !reflect.DeepEqual(oldOriginal, newOriginal) {
		tflog.Debug(ctx, "workflow_steps: semantic-equal after stripping server-minted fields", map[string]any{
			"ignored_paths_per_action_id": ignoredFieldsByActionID,
		})
	}

	return true, diags
}

// stripIgnoredFields parses the JSON and removes, in-place, the fields listed
// in `ignoredFieldsByActionID` for every step that matches a known action id.
// Returns the stripped representation (used for comparison) AND the original
// parsed representation (used to detect whether anything was actually masked).
func stripIgnoredFields(jsonStr string) (stripped, original map[string]any, err error) {
	if err := json.Unmarshal([]byte(jsonStr), &stripped); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal([]byte(jsonStr), &original); err != nil {
		return nil, nil, err
	}

	for stepName, raw := range stripped {
		step, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		actionID, _ := step["actionId"].(string)
		paths, has := ignoredFieldsByActionID[actionID]
		if !has {
			continue
		}
		for _, p := range paths {
			deletePath(step, strings.Split(p, "."))
		}
		_ = stepName
	}

	return stripped, original, nil
}

// deletePath removes the leaf key reached by walking `parts` through nested
// objects in `node`. Missing keys along the way are silently ignored.
func deletePath(node map[string]any, parts []string) {
	if len(parts) == 0 {
		return
	}
	if len(parts) == 1 {
		delete(node, parts[0])
		return
	}
	next, ok := node[parts[0]].(map[string]any)
	if !ok {
		return
	}
	deletePath(next, parts[1:])
}
