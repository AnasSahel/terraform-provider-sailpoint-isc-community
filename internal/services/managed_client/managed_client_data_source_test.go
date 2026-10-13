// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package managed_client

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

// fakeTenantBody is a raw /managed-clients/v1 list response, including the
// secret fields the data source must never put in state.
const fakeTenantBody = `[
  {"id":"mc-1","name":"VA-1","clusterId":"cl-1","type":"VA","status":"NORMAL","vaVersion":"va-1","secret":"s3cr3t","alertKey":"k","vaDownloadUrl":"https://example.invalid"},
  {"id":"mc-2","name":"VA-1","clusterId":"cl-2","type":"VA","status":"WARNING"},
  {"id":"mc-3","name":"VA-3","clusterId":"cl-1","type":"VA","status":"NORMAL"}
]`

func newTestDataSource(t *testing.T) *managedClientDataSource {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/oauth/token":
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
		case "/managed-clients/v1":
			_, _ = w.Write([]byte(fakeTenantBody))
		case "/managed-clients/v1/mc-1":
			_, _ = w.Write([]byte(`{"id":"mc-1","name":"VA-1","clusterId":"cl-1","status":"NORMAL","secret":"s3cr3t"}`))
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
	return &managedClientDataSource{client: c}
}

func readWithConfig(t *testing.T, ds *managedClientDataSource, args map[string]string) (datasource.ReadResponse, managedClientModel) {
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

	var state managedClientModel
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	}
	return resp, state
}

func TestManagedClientDataSourceRead(t *testing.T) {
	ds := newTestDataSource(t)

	cases := []struct {
		name    string
		args    map[string]string
		wantID  string
		wantErr string
	}{
		{name: "by id", args: map[string]string{"id": "mc-1"}, wantID: "mc-1"},
		{name: "unique name", args: map[string]string{"name": "VA-3"}, wantID: "mc-3"},
		{name: "name narrowed by cluster", args: map[string]string{"name": "VA-1", "cluster_id": "cl-2"}, wantID: "mc-2"},
		{name: "ambiguous name", args: map[string]string{"name": "VA-1"}, wantErr: "Multiple managed clients found"},
		{name: "unknown name", args: map[string]string{"name": "nope"}, wantErr: "No managed client found"},
		{name: "unknown id", args: map[string]string{"id": "missing"}, wantErr: "Error Reading SailPoint Managed Client"},
		{name: "no lookup key", args: map[string]string{}, wantErr: "Missing required argument"},
		{name: "id and name", args: map[string]string{"id": "mc-1", "name": "VA-1"}, wantErr: "Conflicting arguments"},
		{name: "id and cluster", args: map[string]string{"id": "mc-1", "cluster_id": "cl-1"}, wantErr: "Conflicting arguments"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, state := readWithConfig(t, ds, tc.args)
			if tc.wantErr != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), tc.wantErr) {
					t.Fatalf("want error %q, got %v", tc.wantErr, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if state.ID.ValueString() != tc.wantID {
				t.Errorf("id = %q, want %q", state.ID.ValueString(), tc.wantID)
			}
		})
	}
}
