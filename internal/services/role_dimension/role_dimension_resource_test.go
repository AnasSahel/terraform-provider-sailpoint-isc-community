// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func testSchema(t *testing.T) resource.SchemaResponse {
	t.Helper()
	var resp resource.SchemaResponse
	(&roleDimensionResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp
}

func TestSchemaValid(t *testing.T) {
	resp := testSchema(t)
	if diags := resp.Schema.ValidateImplementation(context.Background()); diags.HasError() {
		t.Fatalf("schema invalid: %v", diags)
	}
}

func TestRefSetTypesMatchSchema(t *testing.T) {
	sch := testSchema(t).Schema
	want := types.SetType{ElemType: common.ObjectRefObjectType}
	for _, name := range []string{"access_profiles", "entitlements"} {
		if got := sch.Attributes[name].GetType(); !got.Equal(want) {
			t.Errorf("%s: schema type %v != %v", name, got, want)
		}
	}
}

func TestImportState(t *testing.T) {
	sch := testSchema(t).Schema
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	r := &roleDimensionResource{}
	do := func(id string) resource.ImportStateResponse {
		resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		return resp
	}

	resp := do("r1/d1")
	if resp.Diagnostics.HasError() {
		t.Fatalf("r1/d1: %v", resp.Diagnostics)
	}
	var roleID, id types.String
	resp.State.GetAttribute(context.Background(), path.Root("role_id"), &roleID)
	resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	if roleID.ValueString() != "r1" || id.ValueString() != "d1" {
		t.Errorf("state = %s/%s", roleID, id)
	}

	for _, bad := range []string{"d1", "r1/", "/d1", "r1/d1/x"} {
		resp := do(bad)
		if !resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != "Invalid Import ID" {
			t.Errorf("%q: diagnostics = %v, want Invalid Import ID", bad, resp.Diagnostics)
		}
	}
}

// fakeTenant answers OAuth and delegates the rest to handler, recording "METHOD path" for non-OAuth calls.
type fakeTenant struct {
	mu          sync.Mutex
	calls       []string
	contentType string
	body        string
}

func newTestResource(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*roleDimensionResource, *fakeTenant) {
	t.Helper()
	f := &fakeTenant{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/oauth/token" {
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		}
		body, _ := io.ReadAll(r.Body)
		f.mu.Lock()
		f.calls = append(f.calls, r.Method+" "+r.URL.Path)
		f.contentType = r.Header.Get("Content-Type")
		f.body = strings.TrimSpace(string(body))
		f.mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &roleDimensionResource{client: c}, f
}

func objType(t *testing.T) tftypes.Object {
	t.Helper()
	ot, ok := testSchema(t).Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	return ot
}

func stateFor(t *testing.T, m roleDimensionModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: testSchema(t).Schema, Raw: tftypes.NewValue(objType(t), nil)}
	if diags := state.Set(context.Background(), &m); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
	return state
}

func planFor(t *testing.T, m roleDimensionModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: testSchema(t).Schema, Raw: tftypes.NewValue(objType(t), nil)}
	if diags := plan.Set(context.Background(), &m); diags.HasError() {
		t.Fatalf("set plan: %v", diags)
	}
	return plan
}

func getState(t *testing.T, s tfsdk.State) roleDimensionModel {
	t.Helper()
	var m roleDimensionModel
	if diags := s.Get(context.Background(), &m); diags.HasError() {
		t.Fatalf("get state: %v", diags)
	}
	return m
}

func TestCreate(t *testing.T) {
	r, f := newTestResource(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"d1","name":"n","parentId":"r1","owner":{"type":"IDENTITY","id":"i1","name":null},
			"accessProfiles":[{"type":"ACCESS_PROFILE","id":"ap2","name":"AP 2"},{"type":"ACCESS_PROFILE","id":"ap1","name":"AP 1"}],
			"entitlements":[],"membership":null,"created":"2026-01-01T00:00:00Z","modified":"2026-01-01T00:00:00Z"}`))
	})
	plan := roleDimensionModel{
		ID:             types.StringUnknown(),
		RoleID:         types.StringValue("r1"),
		Name:           types.StringValue("n"),
		Description:    types.StringNull(),
		Owner:          &common.ObjectRefModel{Type: types.StringValue("IDENTITY"), ID: types.StringValue("i1"), Name: types.StringUnknown()},
		AccessProfiles: refSet(t, [3]string{"ACCESS_PROFILE", "ap1", ""}, [3]string{"ACCESS_PROFILE", "ap2", ""}),
		Entitlements:   types.SetNull(common.ObjectRefObjectType),
		Created:        types.StringUnknown(),
		Modified:       types.StringUnknown(),
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: testSchema(t).Schema, Raw: tftypes.NewValue(objType(t), nil)}}
	r.Create(context.Background(), resource.CreateRequest{Plan: planFor(t, plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if want := []string{"POST /roles/v1/r1/dimensions"}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
	if want := `{"name":"n","owner":{"type":"IDENTITY","id":"i1"},` +
		`"accessProfiles":[{"type":"ACCESS_PROFILE","id":"ap1"},{"type":"ACCESS_PROFILE","id":"ap2"}]}`; f.body != want {
		t.Errorf("body = %s\nwant   %s", f.body, want)
	}
	got := getState(t, resp.State)
	if got.ID.ValueString() != "d1" || got.RoleID.ValueString() != "r1" || got.Owner == nil || got.Owner.Name.ValueString() != "" {
		t.Errorf("state = %+v", got)
	}
	want := refSet(t, [3]string{"ACCESS_PROFILE", "ap1", "AP 1"}, [3]string{"ACCESS_PROFILE", "ap2", "AP 2"})
	if !got.AccessProfiles.Equal(want) {
		t.Errorf("access_profiles = %v, want %v", got.AccessProfiles, want)
	}
	if !got.Entitlements.IsNull() {
		t.Errorf("entitlements = %v, want null", got.Entitlements)
	}
}

func existingState(t *testing.T) roleDimensionModel {
	t.Helper()
	m := baseModel(t)
	m.Owner.Name = types.StringValue("")
	return m
}

func TestReadNotFoundRemovesResource(t *testing.T) {
	r, _ := newTestResource(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	})
	state := stateFor(t, existingState(t))
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("state should be removed, got %v", resp.State.Raw)
	}
}

func TestReadNormalisesUIValues(t *testing.T) {
	r, f := newTestResource(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"d1","name":"n","owner":null,"accessProfiles":[],"entitlements":[],
			"membership":{"type":"STANDARD","criteria":{"operation":"AND","stringValue":"","children":[
				{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"stringValue":"Engineering"}]}}}`))
	})
	state := stateFor(t, existingState(t))
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if want := []string{"GET /roles/v1/r1/dimensions/d1"}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
	got := getState(t, resp.State)
	if got.Owner != nil {
		t.Errorf("owner = %+v, want nil", got.Owner)
	}
	if !got.Membership.Criteria.StringValue.IsNull() {
		t.Errorf("AND string_value = %v, want null", got.Membership.Criteria.StringValue)
	}
	if !got.AccessProfiles.IsNull() {
		t.Errorf("access_profiles = %v, want null", got.AccessProfiles)
	}
}

func TestUpdateNameOnly(t *testing.T) {
	r, f := newTestResource(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"d1","name":"m","description":"desc","owner":{"type":"IDENTITY","id":"i1"},
			"accessProfiles":[{"type":"ACCESS_PROFILE","id":"ap2","name":"AP 2"},{"type":"ACCESS_PROFILE","id":"ap1","name":"AP 1"}],
			"entitlements":[{"type":"ENTITLEMENT","id":"e1","name":"E 1"}],
			"membership":{"type":"STANDARD","criteria":{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"stringValue":"Engineering"}},
			"created":"2026-01-01T00:00:00Z","modified":"2026-01-02T00:00:00Z"}`))
	})
	state := existingState(t)
	plan := existingState(t)
	plan.Name = types.StringValue("m")
	plan.Modified = types.StringUnknown()
	resp := resource.UpdateResponse{State: stateFor(t, state)}
	r.Update(context.Background(), resource.UpdateRequest{Plan: planFor(t, plan), State: stateFor(t, state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if want := []string{"PATCH /roles/v1/r1/dimensions/d1"}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
	if f.contentType != "application/json-patch+json" {
		t.Errorf("Content-Type = %q", f.contentType)
	}
	if want := `[{"op":"replace","path":"/name","value":"m"}]`; f.body != want {
		t.Errorf("body = %s, want %s", f.body, want)
	}
	got := getState(t, resp.State)
	if got.Name.ValueString() != "m" || got.Modified.ValueString() != "2026-01-02T00:00:00Z" {
		t.Errorf("state = %+v", got)
	}
	if !got.AccessProfiles.Equal(state.AccessProfiles) {
		t.Errorf("access_profiles = %v, want %v", got.AccessProfiles, state.AccessProfiles)
	}
}

func TestDeleteNotFoundIsOK(t *testing.T) {
	r, f := newTestResource(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{}`))
	})
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: stateFor(t, existingState(t))}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if want := []string{"DELETE /roles/v1/r1/dimensions/d1"}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls = %v, want %v", f.calls, want)
	}
}
