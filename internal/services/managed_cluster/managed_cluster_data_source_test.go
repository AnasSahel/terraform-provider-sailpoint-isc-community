// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_cluster

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const testClustersJSON = `[
  {"id":"c1","name":"Prod","type":"idn","clientType":"CCG","ccgVersion":"v01","status":"NORMAL","operational":true,
   "clientIds":["va1"],"serviceCount":2,"configuration":{"gmtOffset":"1"},
   "keyPair":{"publicKey":"-----BEGIN PUBLIC KEY-----x"},"attributes":{"keystore":"secret"}},
  {"id":"c2","name":"Dup","clientType":"CCG","ccgVersion":"v01"},
  {"id":"c3","name":"Dup","clientType":"CCG","ccgVersion":"v01"}
]`

func newTestDataSource(t *testing.T) *managedClusterDataSource {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/oauth/token":
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
		case "/managed-clusters/v1":
			_, _ = w.Write([]byte(testClustersJSON))
		case "/managed-clusters/v1/c1":
			_, _ = w.Write([]byte(`{"id":"c1","name":"Prod","clientType":"CCG","ccgVersion":"v01"}`))
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
	return &managedClusterDataSource{client: c}
}

func readWithConfig(t *testing.T, ds *managedClusterDataSource, args map[string]string) datasource.ReadResponse {
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

func TestManagedClusterDataSourceReadByName(t *testing.T) {
	resp := readWithConfig(t, newTestDataSource(t), map[string]string{"name": "Prod"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var state managedClusterModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state: %v", diags)
	}
	if state.ID.ValueString() != "c1" || state.Status.ValueString() != "NORMAL" || !state.Operational.ValueBool() {
		t.Errorf("state = %+v", state)
	}
	if raw := resp.State.Raw.String(); strings.Contains(raw, "secret") || strings.Contains(raw, "PUBLIC KEY") {
		t.Errorf("key material leaked into state: %s", raw)
	}
}

func TestManagedClusterDataSourceReadByID(t *testing.T) {
	resp := readWithConfig(t, newTestDataSource(t), map[string]string{"id": "c1"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	var state managedClusterModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state: %v", diags)
	}
	if state.Name.ValueString() != "Prod" {
		t.Errorf("name = %v, want Prod", state.Name)
	}
}

func TestManagedClusterDataSourceReadErrors(t *testing.T) {
	cases := map[string]struct {
		args    map[string]string
		summary string
	}{
		"no match":   {map[string]string{"name": "Nope"}, "No managed cluster found"},
		"ambiguous":  {map[string]string{"name": "Dup"}, "Multiple managed clusters found"},
		"no input":   {map[string]string{}, "Missing required argument"},
		"both":       {map[string]string{"id": "c1", "name": "Prod"}, "Conflicting arguments"},
		"id missing": {map[string]string{"id": "zzz"}, "Error Reading SailPoint Managed Cluster"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := readWithConfig(t, newTestDataSource(t), tc.args)
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected an error")
			}
			if got := resp.Diagnostics.Errors()[0].Summary(); got != tc.summary {
				t.Errorf("summary = %q, want %q", got, tc.summary)
			}
		})
	}
}
