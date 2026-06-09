// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// TestStripTopLevelNullKeys verifies that server-injected null-valued keys are
// removed from trigger attributes while real values (including nested nulls)
// are preserved (#136).
func TestStripTopLevelNullKeys(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		in    jsontypes.Normalized
		want  string // expected ValueString(); ignored when wantNull
		isNul bool
	}{
		{
			name: "strips server-added null key",
			in:   jsontypes.NewNormalizedValue(`{"id":"x","integrationId":null}`),
			want: `{"id":"x"}`,
		},
		{
			name: "no null keys is unchanged",
			in:   jsontypes.NewNormalizedValue(`{"id":"x"}`),
			want: `{"id":"x"}`,
		},
		{
			name: "all null keys yields empty object",
			in:   jsontypes.NewNormalizedValue(`{"a":null,"b":null}`),
			want: `{}`,
		},
		{
			name: "nested null is preserved (only top level stripped)",
			in:   jsontypes.NewNormalizedValue(`{"a":{"b":null}}`),
			want: `{"a":{"b":null}}`,
		},
		{
			name:  "null value passes through",
			in:    jsontypes.NewNormalizedNull(),
			isNul: true,
		},
		{
			name: "non-object is left untouched",
			in:   jsontypes.NewNormalizedValue(`["x"]`),
			want: `["x"]`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := stripTopLevelNullKeys(tc.in)
			if tc.isNul {
				if !got.IsNull() {
					t.Fatalf("expected null, got %q", got.ValueString())
				}
				return
			}
			if got.ValueString() != tc.want {
				t.Errorf("stripTopLevelNullKeys(%q) = %q, want %q", tc.in.ValueString(), got.ValueString(), tc.want)
			}
		})
	}
}
