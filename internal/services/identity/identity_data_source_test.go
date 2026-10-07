// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustString(t *testing.T, m types.Map, key string) types.String {
	t.Helper()
	v, ok := m.Elements()[key]
	if !ok {
		t.Fatalf("missing attribute %q", key)
	}
	s, ok := v.(types.String)
	if !ok {
		t.Fatalf("attribute %q is %T, want types.String", key, v)
	}
	return s
}

// Regression test for #171: multi-valued attributes used to fail unmarshalling.
func TestPopulateIdentityDSModel_MixedAttributes(t *testing.T) {
	raw := `{
		"id": "2448605646d0486fb6eb3b4c8bada9a8",
		"name": "Amanda Ross",
		"attributes": {
			"department": "Engineering",
			"groups": ["Admins", "Users"],
			"active": true,
			"employeeCount": 12345678901234567,
			"ratio": 0.5,
			"address": {"city": "Austin", "zip": "78701"},
			"middleName": null
		}
	}`

	var api client.IdentityAPI
	if err := json.Unmarshal([]byte(raw), &api); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	var m identityDSModel
	if diags := populateIdentityDSModel(context.Background(), &m, &api); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	want := map[string]string{
		"department":    "Engineering",
		"groups":        `["Admins","Users"]`,
		"active":        "true",
		"employeeCount": "12345678901234567",
		"ratio":         "0.5",
		"address":       `{"city":"Austin","zip":"78701"}`,
	}
	for key, w := range want {
		if got := mustString(t, m.Attributes, key).ValueString(); got != w {
			t.Errorf("attribute %q: got %q, want %q", key, got, w)
		}
	}
	if !mustString(t, m.Attributes, "middleName").IsNull() {
		t.Errorf("attribute %q: want null", "middleName")
	}
}

func TestPopulateIdentityDSModel_NoAttributes(t *testing.T) {
	api := client.IdentityAPI{ID: "id", Name: "name"}

	var m identityDSModel
	if diags := populateIdentityDSModel(context.Background(), &m, &api); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if !m.Attributes.IsNull() {
		t.Errorf("attributes: want null map, got %v", m.Attributes)
	}
}
