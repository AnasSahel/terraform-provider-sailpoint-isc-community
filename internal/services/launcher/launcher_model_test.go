// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package launcher

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustOwnerObject(t *testing.T, typ, id, name string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(common.ObjectRefObjectType.AttrTypes, map[string]attr.Value{
		"type": types.StringValue(typ),
		"id":   types.StringValue(id),
		"name": types.StringValue(name),
	})
	if diags.HasError() {
		t.Fatalf("building owner object: %v", diags)
	}
	return obj
}

func mustOwnerAttribute(t *testing.T) schema.SingleNestedAttribute {
	t.Helper()
	resp := &resource.SchemaResponse{}
	NewLauncherResource().Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	attr, ok := resp.Schema.Attributes["owner"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("owner is not a SingleNestedAttribute")
	}
	return attr
}

// ISC stamps owner with the calling identity on every write (#177), so it must
// not be Required and must not be pinned to the prior state.
func TestLauncherSchema_OwnerIsOptionalComputedAndNotPinned(t *testing.T) {
	t.Parallel()

	owner := mustOwnerAttribute(t)
	if owner.Required {
		t.Error("owner must not be Required")
	}
	if !owner.Optional || !owner.Computed {
		t.Errorf("owner must be Optional+Computed, got Optional=%v Computed=%v", owner.Optional, owner.Computed)
	}
	if len(owner.PlanModifiers) != 0 {
		t.Errorf("owner must not carry plan modifiers (UseStateForUnknown would pin a stale owner), got %d", len(owner.PlanModifiers))
	}
}

func TestLauncherFromAPI_TakesOwnerFromResponse(t *testing.T) {
	t.Parallel()

	var m launcherModel
	diags := m.FromAPI(context.Background(), client.LauncherAPI{
		ID:    "l-1",
		Owner: &client.ObjectRefAPI{Type: "USER", ID: "caller-b", Name: "Caller B"},
	})
	if diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	if want := mustOwnerObject(t, "USER", "caller-b", "Caller B"); !m.Owner.Equal(want) {
		t.Errorf("owner = %v, want %v", m.Owner, want)
	}
}

func TestLauncherFromAPI_NoOwnerIsNull(t *testing.T) {
	t.Parallel()

	var m launcherModel
	if diags := m.FromAPI(context.Background(), client.LauncherAPI{ID: "l-1"}); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	if !m.Owner.IsNull() {
		t.Errorf("owner = %v, want null", m.Owner)
	}
}

func TestLauncherToAPI_OmitsUnsetOrUnknownOwner(t *testing.T) {
	t.Parallel()

	for name, owner := range map[string]types.Object{
		"null":    types.ObjectNull(common.ObjectRefObjectType.AttrTypes),
		"unknown": types.ObjectUnknown(common.ObjectRefObjectType.AttrTypes),
	} {
		m := launcherModel{Name: types.StringValue("n"), Owner: owner}
		req, diags := m.ToAPI(context.Background())
		if diags.HasError() {
			t.Fatalf("%s: ToAPI: %v", name, diags)
		}
		if req.Owner != nil {
			t.Errorf("%s: owner = %+v, want nil", name, req.Owner)
		}
	}
}

func TestLauncherToAPI_SendsConfiguredOwnerWithoutName(t *testing.T) {
	t.Parallel()

	m := launcherModel{Name: types.StringValue("n"), Owner: mustOwnerObject(t, "USER", "id-1", "stale")}
	req, diags := m.ToAPI(context.Background())
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags)
	}
	if req.Owner == nil || req.Owner.ID != "id-1" || req.Owner.Type != "USER" || req.Owner.Name != "" {
		t.Errorf("owner = %+v, want {USER id-1 \"\"}", req.Owner)
	}
}
