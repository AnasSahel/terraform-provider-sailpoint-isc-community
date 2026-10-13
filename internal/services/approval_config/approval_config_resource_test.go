// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package approval_config

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func testSchema(t *testing.T) resource.SchemaResponse {
	t.Helper()
	var resp resource.SchemaResponse
	(&approvalConfigResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	return resp
}

func TestSchemaValid(t *testing.T) {
	resp := testSchema(t)
	if diags := resp.Schema.ValidateImplementation(context.Background()); diags.HasError() {
		t.Fatalf("schema invalid: %v", diags)
	}
}

func TestAttrTypesMatchSchema(t *testing.T) {
	sch := testSchema(t).Schema
	checks := map[string]types.ObjectType{
		"reminder_config":   {AttrTypes: reminderConfigAttrTypes},
		"escalation_config": {AttrTypes: escalationConfigAttrTypes},
		"timeout_config":    {AttrTypes: timeoutConfigAttrTypes},
		"cron_timezone":     {AttrTypes: cronTimezoneAttrTypes},
		"fallback_approver": {AttrTypes: fallbackApproverAttrTypes},
	}
	for name, want := range checks {
		got := sch.Attributes[name].GetType()
		if !got.Equal(want) {
			t.Errorf("%s: schema type %v != attr types %v", name, got, want)
		}
	}
	chainWant := types.ListType{ElemType: chainItemObjectType}
	if got := sch.Attributes["serial_chain"].GetType(); !got.Equal(chainWant) {
		t.Errorf("serial_chain: schema type %v != %v", got, chainWant)
	}
}

func configWithScope(t *testing.T, scope *string, unknown bool) resource.ValidateConfigRequest {
	t.Helper()
	sch := testSchema(t).Schema
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	switch {
	case unknown:
		vals["scope"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	case scope != nil:
		vals["scope"] = tftypes.NewValue(tftypes.String, *scope)
	}
	return resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: sch, Raw: tftypes.NewValue(objType, vals)}}
}

func TestValidateConfigScope(t *testing.T) {
	r := &approvalConfigResource{}
	role, bad := "ROLE", "DOMAIN_OBJECT"

	var resp resource.ValidateConfigResponse
	r.ValidateConfig(context.Background(), configWithScope(t, &role, false), &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("ROLE: unexpected diagnostics %v", resp.Diagnostics)
	}

	resp = resource.ValidateConfigResponse{}
	r.ValidateConfig(context.Background(), configWithScope(t, &bad, false), &resp)
	if resp.Diagnostics.ErrorsCount() != 1 {
		t.Fatalf("DOMAIN_OBJECT: diagnostics = %v, want one error", resp.Diagnostics)
	}
	if attrDiag, ok := resp.Diagnostics.Errors()[0].(diag.DiagnosticWithPath); !ok || !attrDiag.Path().Equal(path.Root("scope")) {
		t.Errorf("DOMAIN_OBJECT: error is not on path scope: %v", resp.Diagnostics.Errors()[0])
	}

	resp = resource.ValidateConfigResponse{}
	r.ValidateConfig(context.Background(), configWithScope(t, nil, true), &resp)
	if len(resp.Diagnostics) != 0 {
		t.Errorf("unknown scope: diagnostics = %v, want none", resp.Diagnostics)
	}
}

func TestImportState(t *testing.T) {
	sch := testSchema(t).Schema
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	r := &approvalConfigResource{}

	do := func(id string) resource.ImportStateResponse {
		resp := resource.ImportStateResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: id}, &resp)
		return resp
	}

	resp := do("ROLE/r1")
	if resp.Diagnostics.HasError() {
		t.Fatalf("ROLE/r1: %v", resp.Diagnostics)
	}
	var scope, objectID, id types.String
	resp.State.GetAttribute(context.Background(), path.Root("scope"), &scope)
	resp.State.GetAttribute(context.Background(), path.Root("object_id"), &objectID)
	resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	if scope.ValueString() != "ROLE" || objectID.ValueString() != "r1" || id.ValueString() != "ROLE/r1" {
		t.Errorf("state = %s/%s/%s", scope, objectID, id)
	}

	for _, bad := range []string{"r1", "ROLE/", "DOMAIN_OBJECT/r1"} {
		resp := do(bad)
		if !resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != "Invalid Import ID" {
			t.Errorf("%q: diagnostics = %v, want Invalid Import ID", bad, resp.Diagnostics)
		}
	}
}

// fakeTenant answers OAuth and delegates the rest to handler, recording "METHOD path" for non-OAuth calls.
type fakeTenant struct {
	mu    sync.Mutex
	calls []string
}

func (f *fakeTenant) record(r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
}

func newTestResource(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, f *fakeTenant)) (*approvalConfigResource, *fakeTenant) {
	t.Helper()
	f := &fakeTenant{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/oauth/token" {
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		}
		f.record(r)
		handler(w, r, f)
	}))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &approvalConfigResource{client: c}, f
}

func stateFor(t *testing.T, scope string) tfsdk.State {
	t.Helper()
	sch := testSchema(t).Schema
	objType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	state := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}
	m := approvalConfigModel{
		ID:                               types.StringValue(scope + "/r1"),
		Scope:                            types.StringValue(scope),
		ObjectID:                         types.StringValue("r1"),
		ReminderConfig:                   types.ObjectNull(reminderConfigAttrTypes),
		EscalationConfig:                 types.ObjectNull(escalationConfigAttrTypes),
		TimeoutConfig:                    types.ObjectNull(timeoutConfigAttrTypes),
		CronTimezone:                     types.ObjectNull(cronTimezoneAttrTypes),
		SerialChain:                      types.ListNull(chainItemObjectType),
		RequiresComment:                  types.StringNull(),
		FallbackApprover:                 types.ObjectNull(fallbackApproverAttrTypes),
		MachineIdentityManagerAssignment: types.StringNull(),
		CircumventApprovalProcess:        types.BoolNull(),
		AutoApprove:                      types.StringNull(),
	}
	if diags := state.Set(context.Background(), &m); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
	return state
}

func TestReadInheritedRemovesResource(t *testing.T) {
	r, _ := newTestResource(t, func(w http.ResponseWriter, _ *http.Request, _ *fakeTenant) {
		_, _ = w.Write([]byte(`{"id":"tenant-1","scope":"TENANT"}`))
	})
	state := stateFor(t, "ROLE")
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("state should be removed, got %v", resp.State.Raw)
	}
}

func TestReadOwnConfig(t *testing.T) {
	r, _ := newTestResource(t, func(w http.ResponseWriter, _ *http.Request, _ *fakeTenant) {
		_, _ = w.Write([]byte(`{"id":"r1","scope":"DOMAIN_OBJECT","timeoutConfig":{"enabled":true,"daysUntilTimeout":90}}`))
	})
	state := stateFor(t, "ROLE")
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var got approvalConfigModel
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("state: %v", diags)
	}
	if got.Scope.ValueString() != "ROLE" {
		t.Errorf("scope = %s, want ROLE", got.Scope)
	}
	var tc timeoutConfigModel
	if diags := got.TimeoutConfig.As(context.Background(), &tc, objectAsOptions); diags.HasError() {
		t.Fatal(diags)
	}
	if tc.DaysUntilTimeout.ValueInt64() != 90 {
		t.Errorf("days_until_timeout = %v, want 90", tc.DaysUntilTimeout)
	}
}

func TestDeleteRetriesWithDomainObject(t *testing.T) {
	deletedDomainObject := false
	r, f := newTestResource(t, func(w http.ResponseWriter, req *http.Request, _ *fakeTenant) {
		switch {
		case req.Method == http.MethodDelete && strings.HasSuffix(req.URL.Path, "/DOMAIN_OBJECT"):
			deletedDomainObject = true
			w.WriteHeader(http.StatusNoContent)
		case req.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case deletedDomainObject:
			_, _ = w.Write([]byte(`{"id":"tenant-1","scope":"TENANT"}`))
		default:
			_, _ = w.Write([]byte(`{"id":"r1","scope":"DOMAIN_OBJECT"}`))
		}
	})
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: stateFor(t, "ROLE")}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	want := []string{
		"DELETE /generic-approvals/v1/config/r1/ROLE",
		"GET /generic-approvals/v1/config/r1/ROLE",
		"DELETE /generic-approvals/v1/config/r1/DOMAIN_OBJECT",
		"GET /generic-approvals/v1/config/r1/ROLE",
	}
	if !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls = %v\nwant    %v", f.calls, want)
	}
}

func TestDeleteTenantIsNoop(t *testing.T) {
	r, f := newTestResource(t, func(w http.ResponseWriter, req *http.Request, _ *fakeTenant) {
		t.Errorf("unexpected HTTP call %s %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: stateFor(t, "TENANT")}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics)
	}
	if resp.Diagnostics.WarningsCount() != 1 {
		t.Errorf("warnings = %d, want 1", resp.Diagnostics.WarningsCount())
	}
	if len(f.calls) != 0 {
		t.Errorf("calls = %v, want none", f.calls)
	}
}
