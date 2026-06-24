// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Package ignorejson implements the practitioner-driven `ignore_json_changes`
// feature for resources with one or more JSON-string attributes. Each entry in
// the list is a resource-rooted path of the form `<field>.<json-path>` or
// `<field>[i].<json-path>`, where `<field>` is a JSON attribute on the resource
// and the remainder is a minimal jsonpath (see internal/common/jsonpath) into
// that field's parsed JSON. At apply time the practitioner's value is preserved
// at those paths so server-side drift never surfaces.
package ignorejson

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common/jsonpath"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fieldRemainder reports whether path p targets the JSON field `field` and, if
// so, returns the jsonpath remainder (already prefixed with "$"). A path targets
// a field when it is the field name followed by "." (object descent) or "["
// (array index/wildcard).
func fieldRemainder(p, field string) (jsonPath string, ok bool) {
	if strings.HasPrefix(p, field+".") || strings.HasPrefix(p, field+"[") {
		return "$" + p[len(field):], true
	}
	return "", false
}

// elements extracts the path list from the Terraform value, treating
// null/unknown as empty.
func elements(ctx context.Context, ignoreList types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if ignoreList.IsNull() || ignoreList.IsUnknown() {
		return nil, diags
	}
	var paths []string
	diags.Append(ignoreList.ElementsAs(ctx, &paths, false)...)
	return paths, diags
}

// ApplyToField re-injects sourceField's values into newField at every
// ignore_json_changes path that targets fieldName, and returns the reconciled
// field. It is a no-op when either value is null/unknown or no path matches.
func ApplyToField(ctx context.Context, newField, sourceField jsontypes.Normalized, fieldName string, ignoreList types.List) (jsontypes.Normalized, diag.Diagnostics) {
	paths, diags := elements(ctx, ignoreList)
	if diags.HasError() || len(paths) == 0 {
		return newField, diags
	}
	if newField.IsNull() || newField.IsUnknown() || sourceField.IsNull() || sourceField.IsUnknown() {
		return newField, diags
	}

	var jsonPaths []string
	for _, p := range paths {
		if jp, ok := fieldRemainder(p, fieldName); ok {
			jsonPaths = append(jsonPaths, jp)
		}
	}
	if len(jsonPaths) == 0 {
		return newField, diags
	}

	merged, err := jsonpath.PreservePathsInJSON(newField.ValueString(), sourceField.ValueString(), jsonPaths)
	if err != nil {
		diags.AddError("Error applying ignore_json_changes", err.Error())
		return newField, diags
	}
	return jsontypes.NewNormalizedValue(merged), diags
}

// Remainders returns the jsonpath remainders (each rooted at "$") of every
// ignore_json_changes entry that targets fieldName. Entries that target a
// different field are skipped. Use this when a resource needs the raw jsonpaths
// to feed jsonpath.PreservePaths / jsonpath.RemovePaths directly (e.g. a field
// projected from the API that must be pruned, not re-injected).
func Remainders(ctx context.Context, ignoreList types.List, fieldName string) ([]string, diag.Diagnostics) {
	paths, diags := elements(ctx, ignoreList)
	if diags.HasError() || len(paths) == 0 {
		return nil, diags
	}
	var remainders []string
	for _, p := range paths {
		if jp, ok := fieldRemainder(p, fieldName); ok {
			remainders = append(remainders, jp)
		}
	}
	return remainders, diags
}

// ValidatePaths checks that every ignore_json_changes entry targets one of the
// resource's JSON fields and has a parseable jsonpath remainder.
func ValidatePaths(ctx context.Context, ignoreList types.List, jsonFields []string) diag.Diagnostics {
	paths, diags := elements(ctx, ignoreList)
	if diags.HasError() {
		return diags
	}
	for _, p := range paths {
		matched := false
		for _, field := range jsonFields {
			jp, ok := fieldRemainder(p, field)
			if !ok {
				continue
			}
			matched = true
			if err := jsonpath.Validate(jp); err != nil {
				diags.AddError(
					"Invalid ignore_json_changes path",
					fmt.Sprintf("path %q: %s", p, err.Error()),
				)
			}
			break
		}
		if !matched {
			diags.AddError(
				"Invalid ignore_json_changes path",
				fmt.Sprintf("path %q must target one of the JSON fields %v followed by '.' or '['", p, jsonFields),
			)
		}
	}
	return diags
}
