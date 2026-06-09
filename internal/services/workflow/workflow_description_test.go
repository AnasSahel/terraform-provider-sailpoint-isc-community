// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestStringsEmptyEquivalent covers the "" ↔ null normalization guard used to
// keep workflow.description stable across the SailPoint API's empty→null
// rewrite (#136).
func TestStringsEmptyEquivalent(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		a, b types.String
		want bool
	}{
		{"both null", types.StringNull(), types.StringNull(), true},
		{"null and empty", types.StringNull(), types.StringValue(""), true},
		{"both empty", types.StringValue(""), types.StringValue(""), true},
		{"empty and null", types.StringValue(""), types.StringNull(), true},
		{"value and empty", types.StringValue("hello"), types.StringValue(""), false},
		{"value and null", types.StringValue("hello"), types.StringNull(), false},
		{"both value", types.StringValue("hello"), types.StringValue("hello"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := stringsEmptyEquivalent(tc.a, tc.b); got != tc.want {
				t.Errorf("stringsEmptyEquivalent(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
