// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func newTestDataSource(t *testing.T) *roleDimensionDataSource {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/oauth/token":
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/roles/v1/r1/dimensions/d1":
			_, _ = w.Write([]byte(`{"id":"d1","name":"Engineering - London","parentId":"r1",
				"owner":{"type":"IDENTITY","id":"i1","name":null},
				"accessProfiles":[{"type":"ACCESS_PROFILE","id":"ap2","name":"AP 2"},{"type":"ACCESS_PROFILE","id":"ap1","name":"AP 1"}],
				"entitlements":[],
				"membership":{"type":"STANDARD","criteria":{"operation":"AND","stringValue":"","children":[
					{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"stringValue":"London"}]}},
				"created":"2026-01-01T00:00:00Z","modified":"2026-01-01T00:00:00Z"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &roleDimensionDataSource{client: c}
}

func readWithConfig(t *testing.T, ds *roleDimensionDataSource, args map[string]string) datasource.ReadResponse {
	t.Helper()
	ctx := context.Background()
	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	sch := schemaResp.Schema
	objType, ok := sch.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not an object")
	}
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		if v, set := args[name]; set {
			vals[name] = tftypes.NewValue(typ, v)
		} else {
			vals[name] = tftypes.NewValue(typ, nil)
		}
	}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	ds.Read(ctx, datasource.ReadRequest{
		Config: tfsdk.Config{Schema: sch, Raw: tftypes.NewValue(objType, vals)},
	}, &resp)
	return resp
}

func TestDataSourceSchemaValid(t *testing.T) {
	var resp datasource.SchemaResponse
	(&roleDimensionDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	if diags := resp.Schema.ValidateImplementation(context.Background()); diags.HasError() {
		t.Fatalf("schema invalid: %v", diags)
	}
}

func TestDataSourceRead(t *testing.T) {
	resp := readWithConfig(t, newTestDataSource(t), map[string]string{"role_id": "r1", "id": "d1"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var state roleDimensionModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state: %v", diags)
	}
	if state.Name.ValueString() != "Engineering - London" || state.RoleID.ValueString() != "r1" {
		t.Errorf("state = %+v", state)
	}
	if state.Owner == nil || state.Owner.ID.ValueString() != "i1" {
		t.Errorf("owner = %+v", state.Owner)
	}
	if n := len(state.AccessProfiles.Elements()); n != 2 {
		t.Errorf("access_profiles = %d, want 2", n)
	}
	if state.Membership == nil || state.Membership.Criteria == nil || len(state.Membership.Criteria.Children) != 1 {
		t.Fatalf("membership = %+v", state.Membership)
	}
	if !state.Membership.Criteria.StringValue.IsNull() {
		t.Errorf("AND string_value = %v, want null", state.Membership.Criteria.StringValue)
	}
}

func TestDataSourceReadNotFound(t *testing.T) {
	resp := readWithConfig(t, newTestDataSource(t), map[string]string{"role_id": "r1", "id": "nope"})
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected an error")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error Reading SailPoint Role Dimension" {
		t.Errorf("summary = %q", got)
	}
}
