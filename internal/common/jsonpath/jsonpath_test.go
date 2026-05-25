// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package jsonpath

import (
	"testing"
)

func TestParse_valid(t *testing.T) {
	cases := []struct {
		path     string
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
		{path: "foo"},            // no leading $
		{path: "$"},              // no segments
		{path: "$."},             // empty key
		{path: "$.foo["},         // unclosed bracket
		{path: "$.foo[-1].bar"},  // negative index
		{path: "$.foo[abc].bar"}, // non-integer index
		{path: "$!foo"},          // unexpected char
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

	arr := merged["domainSettings"].([]interface{})
	if v := arr[0].(map[string]interface{})["password"]; v != "pass1" {
		t.Errorf("element[0].password = %v, want %q", v, "pass1")
	}
	if v := arr[1].(map[string]interface{})["password"]; v != "pass2" {
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

	arr := merged["domainSettings"].([]interface{})
	if _, ok := arr[0].(map[string]interface{})["password"]; ok {
		t.Error("element[0].password should not be set")
	}
	if v := arr[1].(map[string]interface{})["password"]; v != "secret1" {
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
	// Should silently skip; merged unchanged
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

	arr := merged["settings"].([]interface{})
	if v := arr[0].(map[string]interface{})["token"]; v != "abc" {
		t.Errorf("element[0].token = %v, want %q", v, "abc")
	}
	if _, ok := arr[1].(map[string]interface{})["token"]; ok {
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

	conn := merged["connection"].(map[string]interface{})
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

	domArr := merged["domainSettings"].([]interface{})
	if v := domArr[0].(map[string]interface{})["password"]; v != "pass1" {
		t.Errorf("domainSettings[0].password = %v, want %q", v, "pass1")
	}
	forArr := merged["forestSettings"].([]interface{})
	if v := forArr[0].(map[string]interface{})["servicePassword"]; v != "pass2" {
		t.Errorf("forestSettings[0].servicePassword = %v, want %q", v, "pass2")
	}
}
