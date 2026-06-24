// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package jsonpath

import (
	"encoding/json"
	"fmt"
	"testing"
)

// mustArray asserts that v is a []interface{} and returns it. Fails the test otherwise.
func mustArray(t *testing.T, v interface{}, what string) []interface{} {
	t.Helper()
	arr, ok := v.([]interface{})
	if !ok {
		t.Fatalf("%s: expected []interface{}, got %T", what, v)
	}
	return arr
}

// mustObject asserts that v is a map[string]interface{} and returns it. Fails the test otherwise.
func mustObject(t *testing.T, v interface{}, what string) map[string]interface{} {
	t.Helper()
	obj, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("%s: expected map[string]interface{}, got %T", what, v)
	}
	return obj
}

func TestParse_valid(t *testing.T) {
	cases := []struct {
		path      string
		wantKinds []segKind
		wantKeys  []string
		wantIdxs  []int
	}{
		{
			path:      "$.foo",
			wantKinds: []segKind{segKey},
			wantKeys:  []string{"foo"},
			wantIdxs:  []int{0},
		},
		{
			path:      "$.foo.bar",
			wantKinds: []segKind{segKey, segKey},
			wantKeys:  []string{"foo", "bar"},
			wantIdxs:  []int{0, 0},
		},
		{
			path:      "$.foo[*].bar",
			wantKinds: []segKind{segKey, segWildcard, segKey},
			wantKeys:  []string{"foo", "", "bar"},
			wantIdxs:  []int{0, -1, 0},
		},
		{
			path:      "$.foo[0].bar",
			wantKinds: []segKind{segKey, segIndex, segKey},
			wantKeys:  []string{"foo", "", "bar"},
			wantIdxs:  []int{0, 0, 0},
		},
		{
			path:      "$.foo[2].bar",
			wantKinds: []segKind{segKey, segIndex, segKey},
			wantKeys:  []string{"foo", "", "bar"},
			wantIdxs:  []int{0, 2, 0},
		},
		{
			path:      "$.a.b.c",
			wantKinds: []segKind{segKey, segKey, segKey},
			wantKeys:  []string{"a", "b", "c"},
			wantIdxs:  []int{0, 0, 0},
		},
		// Bracket-quoted keys (single quotes).
		{
			path:      "$['foo']",
			wantKinds: []segKind{segKey},
			wantKeys:  []string{"foo"},
			wantIdxs:  []int{0},
		},
		{
			path:      "$['key with spaces']",
			wantKinds: []segKind{segKey},
			wantKeys:  []string{"key with spaces"},
			wantIdxs:  []int{0},
		},
		// Bracket-quoted keys (double quotes).
		{
			path:      "$[\"foo\"]",
			wantKinds: []segKind{segKey},
			wantKeys:  []string{"foo"},
			wantIdxs:  []int{0},
		},
		{
			path:      "$[\"key with spaces\"]",
			wantKinds: []segKind{segKey},
			wantKeys:  []string{"key with spaces"},
			wantIdxs:  []int{0},
		},
		// Mixed: dot-key prefix followed by bracket-quoted key.
		{
			path:      "$.steps['Create Ticket'].attr",
			wantKinds: []segKind{segKey, segKey, segKey},
			wantKeys:  []string{"steps", "Create Ticket", "attr"},
			wantIdxs:  []int{0, 0, 0},
		},
		// Chained bracket-quoted keys.
		{
			path:      "$['outer']['inner']",
			wantKinds: []segKind{segKey, segKey},
			wantKeys:  []string{"outer", "inner"},
			wantIdxs:  []int{0, 0},
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			segs, err := Parse(tc.path)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.path, err)
			}
			if len(segs) != len(tc.wantKinds) {
				t.Fatalf("Parse(%q) got %d segments, want %d", tc.path, len(segs), len(tc.wantKinds))
			}
			for i, seg := range segs {
				if seg.kind != tc.wantKinds[i] {
					t.Errorf("segment[%d].kind = %v, want %v", i, seg.kind, tc.wantKinds[i])
				}
				if tc.wantKinds[i] == segKey && seg.key != tc.wantKeys[i] {
					t.Errorf("segment[%d].key = %q, want %q", i, seg.key, tc.wantKeys[i])
				}
				if tc.wantKinds[i] == segIndex && seg.idx != tc.wantIdxs[i] {
					t.Errorf("segment[%d].idx = %d, want %d", i, seg.idx, tc.wantIdxs[i])
				}
			}
		})
	}
}

func TestParse_invalid(t *testing.T) {
	cases := []struct {
		path string
	}{
		{path: "foo"},              // no leading $
		{path: "$"},                // no segments
		{path: "$."},               // empty key
		{path: "$.foo["},           // unclosed bracket
		{path: "$.foo[-1].bar"},    // negative index
		{path: "$.foo[abc].bar"},   // non-integer, non-quoted index
		{path: "$!foo"},            // unexpected char
		{path: "$['']"},            // empty single-quoted key
		{path: "$[\"\"]"},          // empty double-quoted key
		{path: "$['mismatched\"]"}, // mismatched quotes (opens ', closes ")
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			_, err := Parse(tc.path)
			if err == nil {
				t.Errorf("Parse(%q) expected error, got nil", tc.path)
			}
		})
	}
}

func TestPreservePaths_topLevelField(t *testing.T) {
	merged := map[string]interface{}{
		"host": "ldap.example.com",
	}
	server := map[string]interface{}{
		"host":     "ldap.example.com",
		"password": "secret",
	}
	if err := PreservePaths(merged, server, []string{"$.password"}); err != nil {
		t.Fatal(err)
	}
	if merged["password"] != "secret" {
		t.Errorf("merged[password] = %v, want %q", merged["password"], "secret")
	}
}

func TestPreservePaths_wildcardArrayField(t *testing.T) {
	merged := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{
				"forestName": "corp.example.com",
				"user":       "svc-account",
			},
			map[string]interface{}{
				"forestName": "other.example.com",
				"user":       "svc-other",
			},
		},
	}
	server := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{
				"forestName": "corp.example.com",
				"user":       "svc-account",
				"password":   "pass1",
			},
			map[string]interface{}{
				"forestName": "other.example.com",
				"user":       "svc-other",
				"password":   "pass2",
			},
		},
	}

	if err := PreservePaths(merged, server, []string{"$.domainSettings[*].password"}); err != nil {
		t.Fatal(err)
	}

	arr := mustArray(t, merged["domainSettings"], "merged[domainSettings]")
	if v := mustObject(t, arr[0], "arr[0]")["password"]; v != "pass1" {
		t.Errorf("element[0].password = %v, want %q", v, "pass1")
	}
	if v := mustObject(t, arr[1], "arr[1]")["password"]; v != "pass2" {
		t.Errorf("element[1].password = %v, want %q", v, "pass2")
	}
}

func TestPreservePaths_specificIndexField(t *testing.T) {
	merged := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{"forestName": "a"},
			map[string]interface{}{"forestName": "b"},
		},
	}
	server := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{"forestName": "a", "password": "secret0"},
			map[string]interface{}{"forestName": "b", "password": "secret1"},
		},
	}

	if err := PreservePaths(merged, server, []string{"$.domainSettings[1].password"}); err != nil {
		t.Fatal(err)
	}

	arr := mustArray(t, merged["domainSettings"], "merged[domainSettings]")
	if _, ok := mustObject(t, arr[0], "arr[0]")["password"]; ok {
		t.Error("element[0].password should not be set")
	}
	if v := mustObject(t, arr[1], "arr[1]")["password"]; v != "secret1" {
		t.Errorf("element[1].password = %v, want %q", v, "secret1")
	}
}

func TestPreservePaths_serverMissingKey(t *testing.T) {
	merged := map[string]interface{}{
		"foo": "bar",
	}
	server := map[string]interface{}{
		"foo": "bar",
		// "password" absent from server
	}
	// Should silently skip; merged unchanged.
	if err := PreservePaths(merged, server, []string{"$.password"}); err != nil {
		t.Fatal(err)
	}
	if _, ok := merged["password"]; ok {
		t.Error("password should not be injected when absent from server")
	}
}

func TestPreservePaths_mergedHasMoreElementsThanServer(t *testing.T) {
	// User added a new element that doesn't exist on the server yet.
	// Wildcard should only process indices present in both.
	merged := map[string]interface{}{
		"settings": []interface{}{
			map[string]interface{}{"name": "existing"},
			map[string]interface{}{"name": "new"},
		},
	}
	server := map[string]interface{}{
		"settings": []interface{}{
			map[string]interface{}{"name": "existing", "token": "abc"},
			// server only has 1 element
		},
	}

	if err := PreservePaths(merged, server, []string{"$.settings[*].token"}); err != nil {
		t.Fatal(err)
	}

	arr := mustArray(t, merged["settings"], "merged[settings]")
	if v := mustObject(t, arr[0], "arr[0]")["token"]; v != "abc" {
		t.Errorf("element[0].token = %v, want %q", v, "abc")
	}
	if _, ok := mustObject(t, arr[1], "arr[1]")["token"]; ok {
		t.Error("element[1].token should not be set (new element not in server)")
	}
}

func TestPreservePaths_nestedObject(t *testing.T) {
	merged := map[string]interface{}{
		"connection": map[string]interface{}{
			"host": "ldap.example.com",
		},
	}
	server := map[string]interface{}{
		"connection": map[string]interface{}{
			"host":     "ldap.example.com",
			"password": "secret",
		},
	}

	if err := PreservePaths(merged, server, []string{"$.connection.password"}); err != nil {
		t.Fatal(err)
	}

	conn := mustObject(t, merged["connection"], "merged[connection]")
	if conn["password"] != "secret" {
		t.Errorf("connection.password = %v, want %q", conn["password"], "secret")
	}
}

func TestPreservePaths_invalidPath(t *testing.T) {
	merged := map[string]interface{}{}
	server := map[string]interface{}{}
	err := PreservePaths(merged, server, []string{"invalid-path"})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestPreservePathsInJSON(t *testing.T) {
	cases := []struct {
		name       string
		mergedJSON string
		priorJSON  string
		paths      []string
		wantJSON   string // expected result; use checkKey/checkMissing instead when order varies
		wantErr    bool
	}{
		{
			name:       "top-level key",
			mergedJSON: `{"host":"ldap.example.com"}`,
			priorJSON:  `{"host":"ldap.example.com","password":"secret"}`,
			paths:      []string{"$.password"},
			wantJSON:   `{"host":"ldap.example.com","password":"secret"}`,
		},
		{
			name:       "nested key",
			mergedJSON: `{"conn":{"host":"ldap.example.com"}}`,
			priorJSON:  `{"conn":{"host":"ldap.example.com","password":"secret"}}`,
			paths:      []string{"$.conn.password"},
			wantJSON:   `{"conn":{"host":"ldap.example.com","password":"secret"}}`,
		},
		{
			name:       "wildcard array",
			mergedJSON: `{"items":[{"name":"a"},{"name":"b"}]}`,
			priorJSON:  `{"items":[{"name":"a","token":"t1"},{"name":"b","token":"t2"}]}`,
			paths:      []string{"$.items[*].token"},
		},
		{
			name:       "bracket-quoted key with spaces",
			mergedJSON: `{"steps":{"Create Ticket":{"type":"action"}}}`,
			priorJSON:  `{"steps":{"Create Ticket":{"type":"action","refID":"abc"}}}`,
			paths:      []string{"$.steps['Create Ticket'].refID"},
			wantJSON:   `{"steps":{"Create Ticket":{"refID":"abc","type":"action"}}}`,
		},
		{
			name:       "prior-missing path is no-op",
			mergedJSON: `{"foo":"bar"}`,
			priorJSON:  `{"foo":"bar"}`,
			paths:      []string{"$.password"},
			wantJSON:   `{"foo":"bar"}`,
		},
		{
			name:       "invalid path error",
			mergedJSON: `{"foo":"bar"}`,
			priorJSON:  `{"foo":"bar"}`,
			paths:      []string{"invalid-path"},
			wantErr:    true,
		},
		{
			name:       "invalid merged JSON error",
			mergedJSON: `not-json`,
			priorJSON:  `{}`,
			paths:      []string{"$.foo"},
			wantErr:    true,
		},
		{
			name:       "invalid prior JSON error",
			mergedJSON: `{}`,
			priorJSON:  `not-json`,
			paths:      []string{"$.foo"},
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := PreservePathsInJSON(tc.mergedJSON, tc.priorJSON, tc.paths)
			if tc.wantErr {
				if err == nil {
					t.Errorf("PreservePathsInJSON(%q, %q, %v) expected error, got nil", tc.mergedJSON, tc.priorJSON, tc.paths)
				}
				return
			}
			if err != nil {
				t.Fatalf("PreservePathsInJSON unexpected error: %v", err)
			}
			if tc.wantJSON != "" {
				// Normalise both through JSON round-trip to avoid key-ordering sensitivity.
				var gotObj, wantObj interface{}
				if e := json.Unmarshal([]byte(got), &gotObj); e != nil {
					t.Fatalf("result is not valid JSON: %v", e)
				}
				if e := json.Unmarshal([]byte(tc.wantJSON), &wantObj); e != nil {
					t.Fatalf("wantJSON is not valid JSON: %v", e)
				}
				gotNorm, gotMarshalErr := json.Marshal(gotObj)
				if gotMarshalErr != nil {
					t.Fatalf("marshal result: %v", gotMarshalErr)
				}
				wantNorm, wantMarshalErr := json.Marshal(wantObj)
				if wantMarshalErr != nil {
					t.Fatalf("marshal want: %v", wantMarshalErr)
				}
				if string(gotNorm) != string(wantNorm) {
					t.Errorf("got  %s\nwant %s", gotNorm, wantNorm)
				}
			}
			// For the wildcard case, verify element-level values directly.
			if tc.name == "wildcard array" {
				var obj map[string]interface{}
				if e := json.Unmarshal([]byte(got), &obj); e != nil {
					t.Fatalf("result is not valid JSON: %v", e)
				}
				items := mustArray(t, obj["items"], "items")
				for i, wantToken := range []string{"t1", "t2"} {
					elem := mustObject(t, items[i], fmt.Sprintf("items[%d]", i))
					tokenAny, exists := elem["token"]
					if !exists {
						t.Errorf("items[%d].token missing", i)
						continue
					}
					token, ok := tokenAny.(string)
					if !ok {
						t.Errorf("items[%d].token: expected string, got %T", i, tokenAny)
						continue
					}
					if token != wantToken {
						t.Errorf("items[%d].token = %q, want %q", i, token, wantToken)
					}
				}
			}
		})
	}
}

func TestPreservePaths_multiplePaths(t *testing.T) {
	merged := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{"forestName": "corp"},
		},
		"forestSettings": []interface{}{
			map[string]interface{}{"name": "forest1"},
		},
	}
	server := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{"forestName": "corp", "password": "pass1"},
		},
		"forestSettings": []interface{}{
			map[string]interface{}{"name": "forest1", "servicePassword": "pass2"},
		},
	}

	paths := []string{
		"$.domainSettings[*].password",
		"$.forestSettings[*].servicePassword",
	}
	if err := PreservePaths(merged, server, paths); err != nil {
		t.Fatal(err)
	}

	domArr := mustArray(t, merged["domainSettings"], "merged[domainSettings]")
	if v := mustObject(t, domArr[0], "domArr[0]")["password"]; v != "pass1" {
		t.Errorf("domainSettings[0].password = %v, want %q", v, "pass1")
	}
	forArr := mustArray(t, merged["forestSettings"], "merged[forestSettings]")
	if v := mustObject(t, forArr[0], "forArr[0]")["servicePassword"]; v != "pass2" {
		t.Errorf("forestSettings[0].servicePassword = %v, want %q", v, "pass2")
	}
}

func TestRemovePaths_wildcardLeaf(t *testing.T) {
	merged := map[string]interface{}{
		"domainSettings": []interface{}{
			map[string]interface{}{"domainDN": "DC=example,DC=com", "password": "********"},
			map[string]interface{}{"domainDN": "DC=other,DC=com", "password": "********"},
		},
	}
	if err := RemovePaths(merged, []string{"$.domainSettings[*].password"}); err != nil {
		t.Fatal(err)
	}
	arr := mustArray(t, merged["domainSettings"], "merged[domainSettings]")
	for i := range arr {
		o := mustObject(t, arr[i], "arr[i]")
		if _, ok := o["password"]; ok {
			t.Errorf("element %d still has password", i)
		}
		if o["domainDN"] == nil {
			t.Errorf("element %d lost domainDN (sibling must survive)", i)
		}
	}
}

func TestRemovePaths_topLevelAndIndexAndNested(t *testing.T) {
	merged := map[string]interface{}{
		"secret":   "x",
		"keep":     "y",
		"arr":      []interface{}{map[string]interface{}{"p": 1, "q": 2}, map[string]interface{}{"p": 3}},
		"nestedOb": map[string]interface{}{"inner": map[string]interface{}{"z": 9, "keepZ": 10}},
	}
	paths := []string{"$.secret", "$.arr[0].p", "$.nestedOb.inner.z"}
	if err := RemovePaths(merged, paths); err != nil {
		t.Fatal(err)
	}
	if _, ok := merged["secret"]; ok {
		t.Error("top-level secret not removed")
	}
	if merged["keep"] != "y" {
		t.Error("unrelated top-level key altered")
	}
	arr := mustArray(t, merged["arr"], "merged[arr]")
	if _, ok := mustObject(t, arr[0], "arr[0]")["p"]; ok {
		t.Error("arr[0].p not removed")
	}
	if mustObject(t, arr[1], "arr[1]")["p"] != 3 {
		t.Error("arr[1].p must survive (only index 0 targeted)")
	}
	inner := mustObject(t, mustObject(t, merged["nestedOb"], "nestedOb")["inner"], "inner")
	if _, ok := inner["z"]; ok {
		t.Error("nestedOb.inner.z not removed")
	}
	if inner["keepZ"] != 10 {
		t.Error("nestedOb.inner.keepZ must survive")
	}
}

func TestRemovePaths_absentAndTypeMismatchAreNoOps(t *testing.T) {
	merged := map[string]interface{}{"a": map[string]interface{}{"b": 1}, "scalar": "s"}
	// absent path, and a path that walks through a scalar as if it were a map/array
	paths := []string{"$.missing.deep", "$.scalar[*].x", "$.a.b.c"}
	if err := RemovePaths(merged, paths); err != nil {
		t.Fatal(err)
	}
	if mustObject(t, merged["a"], "a")["b"] != 1 {
		t.Error("a.b should be untouched by no-op paths")
	}
	if merged["scalar"] != "s" {
		t.Error("scalar should be untouched")
	}
}

// TestRemovePathsInJSON_maskedSecret is the #162 core case: a declared top-level
// key carries a masked nested secret the practitioner ignores. After pruning, the
// projected attributes equal the config (which omits the secret) → No changes,
// while non-ignored siblings remain.
func TestRemovePathsInJSON_maskedSecret(t *testing.T) {
	projected := `{"domainSettings":[{"domainDN":"DC=example,DC=com","servers":["192.0.2.10"],"user":"svc-bind@example.com","password":"********"}]}`
	out, err := RemovePathsInJSON(projected, []string{"$.domainSettings[*].password"})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	el := mustObject(t, mustArray(t, got["domainSettings"], "domainSettings")[0], "el")
	if _, ok := el["password"]; ok {
		t.Error("password must be pruned")
	}
	if el["domainDN"] != "DC=example,DC=com" || el["user"] != "svc-bind@example.com" {
		t.Error("non-ignored siblings must survive")
	}
}

func TestRemovePathsInJSON_invalidPath(t *testing.T) {
	if _, err := RemovePathsInJSON(`{"a":1}`, []string{"no-dollar"}); err == nil {
		t.Error("expected error for unparseable path")
	}
}
