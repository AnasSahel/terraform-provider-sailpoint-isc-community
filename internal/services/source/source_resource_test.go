// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustStringList(t *testing.T, vals ...string) types.List {
	t.Helper()
	l, diags := types.ListValueFrom(context.Background(), types.StringType, vals)
	if diags.HasError() {
		t.Fatalf("building list: %v", diags)
	}
	return l
}

// TestUnit_projectConnectorAttributes verifies the projection helper in
// isolation — no SailPoint credentials needed. All values are invented; no
// real IDs, tenant names, or client data appear in this file.
func TestUnit_projectConnectorAttributes(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		managed  string // JSON: keys the user declared
		full     string // JSON: full API response (superset)
		wantKeys []string
		wantVals map[string]interface{}
	}{
		"subset of server keys": {
			// User declares two keys; server injects status, healthy, since, etc.
			// Projection must retain only the two declared keys.
			managed: `{"cloudDisplayName":"Example Connector","delimiter":"|"}`,
			full: `{
				"cloudDisplayName": "Example Connector",
				"delimiter":        "|",
				"status":           "SOURCE_STATE_HEALTHY",
				"healthy":          true,
				"since":            "2026-01-01T00:00:00Z",
				"cloudEnabled":     true
			}`,
			wantKeys: []string{"cloudDisplayName", "delimiter"},
			wantVals: map[string]interface{}{
				"cloudDisplayName": "Example Connector",
				"delimiter":        "|",
			},
		},
		"server normalised a declared value": {
			// Server may canonicalise a value (e.g., fix casing). The projected
			// value must reflect the server's version so the plan surfaces drift.
			managed:  `{"cloudDisplayName":"example connector"}`,
			full:     `{"cloudDisplayName":"Example Connector","status":"healthy"}`,
			wantKeys: []string{"cloudDisplayName"},
			wantVals: map[string]interface{}{
				"cloudDisplayName": "Example Connector",
			},
		},
		"server deleted a declared key": {
			// If a key the user declared is no longer in the API response,
			// keep the managed value so the plan surfaces the divergence.
			managed:  `{"requiredKey":"expected-value","other":"val"}`,
			full:     `{"other":"val","serverOnly":"ignored"}`,
			wantKeys: []string{"requiredKey", "other"},
			wantVals: map[string]interface{}{
				"requiredKey": "expected-value",
				"other":       "val",
			},
		},
		"all declared keys present no server extras": {
			managed:  `{"a":"1","b":"2"}`,
			full:     `{"a":"1","b":"2"}`,
			wantKeys: []string{"a", "b"},
			wantVals: map[string]interface{}{"a": "1", "b": "2"},
		},
		"empty managed set": {
			// Empty managed set: projection must produce an empty object,
			// not propagate any server keys.
			managed:  `{}`,
			full:     `{"status":"healthy","since":"2026-01-01"}`,
			wantKeys: []string{},
			wantVals: map[string]interface{}{},
		},
		"boolean and numeric values round-trip": {
			managed:  `{"hasHeader":true,"maxRows":1000}`,
			full:     `{"hasHeader":true,"maxRows":1000,"serverKey":"serverval"}`,
			wantKeys: []string{"hasHeader", "maxRows"},
			wantVals: map[string]interface{}{
				"hasHeader": true,
				"maxRows":   float64(1000), // JSON unmarshal produces float64
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			managed := jsontypes.NewNormalizedValue(tc.managed)
			full := jsontypes.NewNormalizedValue(tc.full)

			got, diags := projectConnectorAttributes(managed, full)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			var gotMap map[string]interface{}
			if err := json.Unmarshal([]byte(got.ValueString()), &gotMap); err != nil {
				t.Fatalf("projected value is not valid JSON: %v", err)
			}

			// Assert exact key set (no extras).
			if len(gotMap) != len(tc.wantKeys) {
				t.Errorf("key count: got %d (%v), want %d (%v)",
					len(gotMap), mapKeys(gotMap), len(tc.wantKeys), tc.wantKeys)
			}
			for _, wk := range tc.wantKeys {
				if _, ok := gotMap[wk]; !ok {
					t.Errorf("expected key %q missing from projected output (got: %v)", wk, mapKeys(gotMap))
				}
			}

			// Assert values for declared keys.
			for k, wantV := range tc.wantVals {
				gotV, ok := gotMap[k]
				if !ok {
					t.Errorf("missing expected key %q in projected output", k)
					continue
				}
				// Compare via re-marshalled JSON to avoid float64/interface{} mismatches.
				gotJSON, _ := json.Marshal(gotV)
				wantJSON, _ := json.Marshal(wantV)
				if string(gotJSON) != string(wantJSON) {
					t.Errorf("key %q: got %s, want %s", k, gotJSON, wantJSON)
				}
			}

			// Assert no extra keys beyond the declared set.
			wantKeySet := make(map[string]bool, len(tc.wantKeys))
			for _, k := range tc.wantKeys {
				wantKeySet[k] = true
			}
			for k := range gotMap {
				if !wantKeySet[k] {
					t.Errorf("unexpected extra key %q in projected output", k)
				}
			}
		})
	}
}

// TestUnit_ignoreJSONChanges_readProjectionPrune exercises the #162 fix path at
// the model level: the Read projection copies the full server value of a declared
// top-level key (masked nested secret included), then the ignored path is pruned
// so the managed connector_attributes matches the config (which omits it) and
// reaches No changes — while non-ignored siblings still surface drift. No
// SailPoint credentials needed; all values are invented.
func TestUnit_ignoreJSONChanges_readProjectionPrune(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Config / prior managed state: declares domainSettings WITHOUT password.
	const managed = `{"domainSettings":[{"domainDN":"DC=example,DC=com","user":"svc@example.com"}]}`
	// Full API set: server returns the password masked, plus undeclared keys.
	const full = `{"domainSettings":[{"domainDN":"DC=example,DC=com","user":"svc@example.com","password":"********"}],"status":"SOURCE_STATE_HEALTHY"}`

	cases := map[string]struct {
		model sourceModel
	}{
		"via ignore_json_changes (new, field-rooted)": {
			model: sourceModel{IgnoreJSONChanges: mustStringList(t, "connector_attributes.domainSettings[*].password")},
		},
		"via ignore_attributes_paths (deprecated, $-rooted)": {
			model: sourceModel{IgnoreAttributesPaths: mustStringList(t, "$.domainSettings[*].password")},
		},
		"both set, union applied": {
			model: sourceModel{
				IgnoreJSONChanges:     mustStringList(t, "connector_attributes.domainSettings[*].password"),
				IgnoreAttributesPaths: mustStringList(t, "$.status"),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			paths, diags := tc.model.resolvedIgnorePaths(ctx)
			if diags.HasError() {
				t.Fatalf("resolvedIgnorePaths: %v", diags)
			}

			projected, diags := projectConnectorAttributes(jsontypes.NewNormalizedValue(managed), jsontypes.NewNormalizedValue(full))
			if diags.HasError() {
				t.Fatalf("projectConnectorAttributes: %v", diags)
			}
			pruned, diags := pruneIgnoredPaths(projected, paths)
			if diags.HasError() {
				t.Fatalf("pruneIgnoredPaths: %v", diags)
			}

			var got map[string]interface{}
			if err := json.Unmarshal([]byte(pruned.ValueString()), &got); err != nil {
				t.Fatalf("pruned value not valid JSON: %v", err)
			}

			arr, ok := got["domainSettings"].([]interface{})
			if !ok || len(arr) != 1 {
				t.Fatalf("domainSettings shape unexpected: %v", got["domainSettings"])
			}
			el, ok := arr[0].(map[string]interface{})
			if !ok {
				t.Fatalf("domainSettings[0] not an object: %v", arr[0])
			}
			if _, ok := el["password"]; ok {
				t.Errorf("ignored password must be pruned, got %v", el["password"])
			}
			if el["domainDN"] != "DC=example,DC=com" || el["user"] != "svc@example.com" {
				t.Errorf("non-ignored siblings must survive, got %v", el)
			}
			// status is a server-managed undeclared key — projection drops it
			// regardless of the ignore set.
			if _, ok := got["status"]; ok {
				t.Errorf("undeclared server key must not be in managed attributes, got %v", got)
			}
		})
	}
}

// TestUnit_ignoreJSONChanges_siblingDriftSurfaces confirms a real change on a
// NON-ignored field inside a declared key is still reflected after prune.
func TestUnit_ignoreJSONChanges_siblingDriftSurfaces(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	const managed = `{"domainSettings":[{"domainDN":"DC=old,DC=com","password":"********"}]}`
	// Server changed domainDN (real drift) and still masks the password.
	const full = `{"domainSettings":[{"domainDN":"DC=new,DC=com","password":"********"}]}`

	m := sourceModel{IgnoreJSONChanges: mustStringList(t, "connector_attributes.domainSettings[*].password")}
	paths, _ := m.resolvedIgnorePaths(ctx)
	projected, _ := projectConnectorAttributes(jsontypes.NewNormalizedValue(managed), jsontypes.NewNormalizedValue(full))
	pruned, _ := pruneIgnoredPaths(projected, paths)

	var got map[string]interface{}
	if err := json.Unmarshal([]byte(pruned.ValueString()), &got); err != nil {
		t.Fatal(err)
	}
	arr, ok := got["domainSettings"].([]interface{})
	if !ok || len(arr) != 1 {
		t.Fatalf("domainSettings shape unexpected: %v", got["domainSettings"])
	}
	el, ok := arr[0].(map[string]interface{})
	if !ok {
		t.Fatalf("domainSettings[0] not an object: %v", arr[0])
	}
	if el["domainDN"] != "DC=new,DC=com" {
		t.Errorf("sibling drift must surface: domainDN = %v, want DC=new,DC=com", el["domainDN"])
	}
	if _, ok := el["password"]; ok {
		t.Error("ignored password must still be pruned")
	}
}

// mapKeys returns the key slice of a map, for use in error messages only.
func mapKeys(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
