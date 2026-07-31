// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"encoding/json"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
)

func TestIdentityAPI_UnmarshalMixedAttributes(t *testing.T) {
	t.Helper()

	raw := `{
		"id": "abc123",
		"name": "Test User",
		"attributes": {
			"department": "Engineering",
			"roles": ["Admin", "User"],
			"active": true,
			"count": 42
		}
	}`

	var identity client.IdentityAPI
	if err := json.Unmarshal([]byte(raw), &identity); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if identity.ID != "abc123" {
		t.Fatalf("expected id abc123, got %q", identity.ID)
	}
	if len(identity.Attributes) != 4 {
		t.Fatalf("expected 4 attributes, got %d", len(identity.Attributes))
	}

	roles, ok := identity.Attributes["roles"].([]interface{})
	if !ok {
		t.Fatalf("roles attribute not an array, got %T", identity.Attributes["roles"])
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
}

func TestIdentityAttributesToStringMap(t *testing.T) {
	t.Helper()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name: "scalar and collection values",
			input: map[string]interface{}{
				"department": "Engineering",
				"roles":      []interface{}{"Admin", "User"},
				"active":     true,
				"count":      float64(42),
			},
			expected: map[string]string{
				"department": "Engineering",
				"roles":      `["Admin","User"]`,
				"active":     "true",
				"count":      "42",
			},
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, diags := identityAttributesToStringMap(tt.input)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			if tt.expected == nil {
				if result != nil {
					t.Fatalf("expected nil map, got %#v", result)
				}
				return
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d attributes, got %d: %#v", len(tt.expected), len(result), result)
			}

			for key, want := range tt.expected {
				got, ok := result[key]
				if !ok {
					t.Fatalf("missing key %q", key)
				}
				if got != want {
					t.Fatalf("key %q: expected %q, got %q", key, want, got)
				}
			}
		})
	}
}
