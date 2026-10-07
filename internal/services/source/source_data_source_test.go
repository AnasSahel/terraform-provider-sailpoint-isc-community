// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// newTestSourceDataSource wires the data source to a fake tenant that serves
// the given sources from both /v2025/sources and /v2025/sources/{id}.
func newTestSourceDataSource(t *testing.T, sources []client.SourceAPI, gotFilters *string) *sourceDataSource {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/oauth/token":
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		case "/v2025/sources":
			*gotFilters = r.URL.Query().Get("filters")
			if err := json.NewEncoder(w).Encode(sources); err != nil {
				t.Errorf("encode: %v", err)
			}
			return
		}
		for _, s := range sources {
			if r.URL.Path == "/v2025/sources/"+s.ID {
				if err := json.NewEncoder(w).Encode(s); err != nil {
					t.Errorf("encode: %v", err)
				}
				return
			}
		}
		t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &sourceDataSource{client: c}
}

func readSourceWithConfig(t *testing.T, ds *sourceDataSource, args map[string]string) datasource.ReadResponse {
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

var testADSource = client.SourceAPI{
	ID:        "src-1",
	Name:      "Active Directory",
	Connector: "active-directory",
	Owner:     &client.ObjectRefAPI{Type: "IDENTITY", ID: "owner-1", Name: "Jane"},
}

func TestSourceDataSourceRead(t *testing.T) {
	cases := []struct {
		name        string
		args        map[string]string
		wantFilters string
	}{
		{name: "by name", args: map[string]string{"name": "Active Directory"}, wantFilters: `name eq "Active Directory"`},
		{name: "by id", args: map[string]string{"id": "src-1"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotFilters string
			ds := newTestSourceDataSource(t, []client.SourceAPI{testADSource}, &gotFilters)

			resp := readSourceWithConfig(t, ds, tc.args)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if gotFilters != tc.wantFilters {
				t.Errorf("filters = %q, want %q", gotFilters, tc.wantFilters)
			}

			var state sourceDataSourceModel
			if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
				t.Fatalf("reading state: %v", diags)
			}
			if state.ID.ValueString() != "src-1" {
				t.Errorf("id = %q, want src-1", state.ID.ValueString())
			}
			if state.Name.ValueString() != "Active Directory" {
				t.Errorf("name = %q, want Active Directory", state.Name.ValueString())
			}
			if state.Connector.ValueString() != "active-directory" {
				t.Errorf("connector = %q", state.Connector.ValueString())
			}
		})
	}
}

func TestSourceDataSourceReadErrors(t *testing.T) {
	cases := []struct {
		name    string
		sources []client.SourceAPI
		args    map[string]string
		summary string
	}{
		{
			name:    "no lookup argument",
			summary: "Missing required argument",
		},
		{
			name:    "id and name together",
			args:    map[string]string{"id": "src-1", "name": "Active Directory"},
			summary: "Conflicting arguments",
		},
		{
			name:    "name not found",
			sources: []client.SourceAPI{},
			args:    map[string]string{"name": "Missing"},
			summary: "No source found",
		},
		{
			name:    "ambiguous name",
			sources: []client.SourceAPI{{ID: "a", Name: "dup"}, {ID: "b", Name: "dup"}},
			args:    map[string]string{"name": "dup"},
			summary: "Multiple sources found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotFilters string
			ds := newTestSourceDataSource(t, tc.sources, &gotFilters)
			resp := readSourceWithConfig(t, ds, tc.args)
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected an error")
			}
			if got := resp.Diagnostics.Errors()[0].Summary(); got != tc.summary {
				t.Errorf("summary = %q, want %q", got, tc.summary)
			}
		})
	}
}
