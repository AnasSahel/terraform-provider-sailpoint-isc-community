// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package ignorejson

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustList(t *testing.T, vals ...string) types.List {
	t.Helper()
	l, diags := types.ListValueFrom(context.Background(), types.StringType, vals)
	if diags.HasError() {
		t.Fatalf("building list: %v", diags)
	}
	return l
}

func TestApplyToField(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("preserves source value at matching path", func(t *testing.T) {
		newF := jsontypes.NewNormalizedValue(`{"refID":"server-minted","x":1}`)
		srcF := jsontypes.NewNormalizedValue(`{"refID":"user","x":1}`)
		got, diags := ApplyToField(ctx, newF, srcF, "attributes", mustList(t, "attributes.refID"))
		if diags.HasError() {
			t.Fatalf("diags: %v", diags)
		}
		if got.ValueString() != `{"refID":"user","x":1}` {
			t.Errorf("got %q", got.ValueString())
		}
	})

	t.Run("no matching field is a no-op", func(t *testing.T) {
		newF := jsontypes.NewNormalizedValue(`{"refID":"server"}`)
		srcF := jsontypes.NewNormalizedValue(`{"refID":"user"}`)
		got, _ := ApplyToField(ctx, newF, srcF, "attributes", mustList(t, "properties.refID"))
		if got.ValueString() != `{"refID":"server"}` {
			t.Errorf("expected unchanged, got %q", got.ValueString())
		}
	})

	t.Run("null list is a no-op", func(t *testing.T) {
		newF := jsontypes.NewNormalizedValue(`{"a":1}`)
		got, _ := ApplyToField(ctx, newF, newF, "attributes", types.ListNull(types.StringType))
		if got.ValueString() != `{"a":1}` {
			t.Errorf("expected unchanged, got %q", got.ValueString())
		}
	})

	t.Run("null field is a no-op", func(t *testing.T) {
		newF := jsontypes.NewNormalizedNull()
		srcF := jsontypes.NewNormalizedValue(`{"a":1}`)
		got, _ := ApplyToField(ctx, newF, srcF, "attributes", mustList(t, "attributes.a"))
		if !got.IsNull() {
			t.Errorf("expected null, got %q", got.ValueString())
		}
	})
}

func TestValidatePaths(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name    string
		paths   []string
		fields  []string
		wantErr bool
	}{
		{"object path ok", []string{"attributes.param_oauth.refID"}, []string{"attributes"}, false},
		{"array path ok", []string{"form_elements[*].id"}, []string{"form_elements"}, false},
		{"wrong field", []string{"other.foo"}, []string{"attributes"}, true},
		{"trailing dot is invalid jsonpath", []string{"attributes."}, []string{"attributes"}, true},
		{"bare field with no subpath", []string{"attributes"}, []string{"attributes"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diags := ValidatePaths(ctx, mustList(t, tc.paths...), tc.fields)
			if diags.HasError() != tc.wantErr {
				t.Errorf("ValidatePaths(%v) hasErr=%v want %v: %v", tc.paths, diags.HasError(), tc.wantErr, diags)
			}
		})
	}
}
