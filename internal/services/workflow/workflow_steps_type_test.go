// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// stepsValue is a small helper to keep the test cases compact.
func stepsValue(s string) workflowStepsValue {
	return workflowStepsValue{Normalized: jsontypes.NewNormalizedValue(s)}
}

func TestWorkflowStepsValue_StringSemanticEquals(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old, new string
		want     bool
		wantErr  bool
	}{
		"identical": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"url":"https://example.com"}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"url":"https://example.com"}}}`,
			want: true,
		},
		"whitespace and key order differences": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"url":"https://example.com"}}}`,
			new:  "{\n  \"step1\": {\n    \"attributes\": {\"url\": \"https://example.com\"},\n    \"actionId\": \"sp:http\"\n  }\n}",
			want: true,
		},
		"sp:http param_oauth.refID differs only": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"old-uuid","paramType":"1.4"}}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"new-uuid","paramType":"1.4"}}}}`,
			want: true,
		},
		"sp:http param_header.refID differs only": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"param_header":{"refID":"old","paramType":"1.3"}}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"param_header":{"refID":"new","paramType":"1.3"}}}}`,
			want: true,
		},
		"sp:http param_oauth_scopes.refID differs only": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth_scopes":{"refID":"old","paramType":"3.1"}}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth_scopes":{"refID":"new","paramType":"3.1"}}}}`,
			want: true,
		},
		"sp:http with all three minted refIDs differing": {
			old: `{"step1":{"actionId":"sp:http","attributes":{` +
				`"param_oauth":{"refID":"o1","paramType":"1.4"},` +
				`"param_header":{"refID":"h1","paramType":"1.3"},` +
				`"param_oauth_scopes":{"refID":"s1","paramType":"3.1"}` +
				`}}}`,
			new: `{"step1":{"actionId":"sp:http","attributes":{` +
				`"param_oauth":{"refID":"o2","paramType":"1.4"},` +
				`"param_header":{"refID":"h2","paramType":"1.3"},` +
				`"param_oauth_scopes":{"refID":"s2","paramType":"3.1"}` +
				`}}}`,
			want: true,
		},
		"sp:http refID differs AND another field differs": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"url":"https://a","param_oauth":{"refID":"old"}}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"url":"https://b","param_oauth":{"refID":"new"}}}}`,
			want: false,
		},
		"sp:http paramID differs (not in ignore list)": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"x","paramID":"old"}}}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"x","paramID":"new"}}}}`,
			want: false,
		},
		"unknown actionId — refID diff is NOT ignored": {
			old:  `{"step1":{"actionId":"sp:custom","attributes":{"param_oauth":{"refID":"old"}}}}`,
			new:  `{"step1":{"actionId":"sp:custom","attributes":{"param_oauth":{"refID":"new"}}}}`,
			want: false,
		},
		"missing actionId on step — refID diff is NOT ignored": {
			old:  `{"step1":{"attributes":{"param_oauth":{"refID":"old"}}}}`,
			new:  `{"step1":{"attributes":{"param_oauth":{"refID":"new"}}}}`,
			want: false,
		},
		"step renamed (key changes) — not equal": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{}}}`,
			new:  `{"renamed":{"actionId":"sp:http","attributes":{}}}`,
			want: false,
		},
		"sibling step (non sp:http) differs — not equal": {
			old:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"x"}}},"endOk":{"type":"success","displayName":"End A"}}`,
			new:  `{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"y"}}},"endOk":{"type":"success","displayName":"End B"}}`,
			want: false,
		},
		"step with no attributes block": {
			old:  `{"endOk":{"type":"success","displayName":"End"}}`,
			new:  `{"endOk":{"type":"success","displayName":"End"}}`,
			want: true,
		},
		"invalid JSON on old side": {
			old:     `{not json`,
			new:     `{"step1":{}}`,
			wantErr: true,
		},
		"invalid JSON on new side": {
			old:     `{"step1":{}}`,
			new:     `{not json`,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			old := stepsValue(tc.old)
			got, diags := old.StringSemanticEquals(context.Background(), stepsValue(tc.new))

			if tc.wantErr {
				if !diags.HasError() {
					t.Fatalf("expected diagnostics error, got none")
				}
				return
			}
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if got != tc.want {
				t.Errorf("StringSemanticEquals = %v, want %v\nold=%s\nnew=%s", got, tc.want, tc.old, tc.new)
			}
		})
	}
}

func TestWorkflowStepsValue_StringSemanticEquals_AcceptsBareNormalized(t *testing.T) {
	t.Parallel()
	// The framework may pass a plain jsontypes.Normalized (not wrapped in our
	// custom value) when comparing. The implementation must accept both.
	old := stepsValue(`{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"a"}}}}`)
	bareNew := jsontypes.NewNormalizedValue(`{"step1":{"actionId":"sp:http","attributes":{"param_oauth":{"refID":"b"}}}}`)

	got, diags := old.StringSemanticEquals(context.Background(), bareNew)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !got {
		t.Errorf("expected semantic equal across our value and bare jsontypes.Normalized, got false")
	}
}
