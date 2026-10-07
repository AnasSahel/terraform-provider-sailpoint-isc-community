// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestBuildEntitlementFilters(t *testing.T) {
	cases := []struct {
		name  string
		model entitlementDSModel
		want  string
	}{
		{name: "none", want: ""},
		{
			name: "name only",
			model: entitlementDSModel{
				entitlementModel: entitlementModel{Name: types.StringValue("My Launcher")},
			},
			want: `name eq "My Launcher"`,
		},
		{
			name: "all filters in stable order",
			model: entitlementDSModel{
				entitlementModel: entitlementModel{
					Name:      types.StringValue("n"),
					Value:     types.StringValue("v"),
					Attribute: types.StringValue("a"),
				},
				SourceID: types.StringValue("src"),
			},
			want: `name eq "n" and value eq "v" and attribute eq "a" and source.id eq "src"`,
		},
		{
			name: "quotes and backslashes are escaped",
			model: entitlementDSModel{
				entitlementModel: entitlementModel{Value: types.StringValue(`CN="x",OU=a\b`)},
			},
			want: `value eq "CN=\"x\",OU=a\\b"`,
		},
		{
			name: "null and unknown are ignored",
			model: entitlementDSModel{
				entitlementModel: entitlementModel{Name: types.StringNull(), Value: types.StringUnknown()},
				SourceID:         types.StringValue(""),
			},
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := buildEntitlementFilters(tc.model); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// newTestDataSource wires the data source to a fake tenant whose
// /v2025/entitlements list endpoint returns the given entitlements.
func newTestDataSource(t *testing.T, ents []client.EntitlementAPI, gotFilters *string) *entitlementDataSource {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/oauth/token" {
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		}
		if r.URL.Path != "/v2025/entitlements" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		*gotFilters = r.URL.Query().Get("filters")
		if err := json.NewEncoder(w).Encode(ents); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	t.Cleanup(srv.Close)

	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &entitlementDataSource{client: c}
}

func readWithConfig(t *testing.T, ds *entitlementDataSource, args map[string]string) datasource.ReadResponse {
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

func TestEntitlementDataSourceReadByFilters(t *testing.T) {
	var gotFilters string
	ds := newTestDataSource(t, []client.EntitlementAPI{{
		ID:        "ent-1",
		Name:      "My Launcher",
		Attribute: "launcher",
		Value:     "launcher-123",
		Source:    &client.ObjectRefAPI{Type: "SOURCE", ID: "src-1", Name: "IdentityNow"},
	}}, &gotFilters)

	resp := readWithConfig(t, ds, map[string]string{"name": "My Launcher"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if gotFilters != `name eq "My Launcher"` {
		t.Errorf("filters = %q", gotFilters)
	}

	var state entitlementDSModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("reading state: %v", diags)
	}
	if state.ID.ValueString() != "ent-1" {
		t.Errorf("id = %q, want ent-1", state.ID.ValueString())
	}
	if state.SourceID.ValueString() != "src-1" {
		t.Errorf("source_id = %q, want src-1", state.SourceID.ValueString())
	}
	if state.Value.ValueString() != "launcher-123" {
		t.Errorf("value = %q", state.Value.ValueString())
	}
}

func TestEntitlementDataSourceReadErrors(t *testing.T) {
	cases := []struct {
		name    string
		ents    []client.EntitlementAPI
		args    map[string]string
		summary string
		detail  string
	}{
		{
			name:    "no lookup argument",
			summary: "Missing required argument",
		},
		{
			name:    "id combined with filters",
			args:    map[string]string{"id": "x", "name": "y"},
			summary: "Conflicting arguments",
		},
		{
			name:    "not found mentions generation delay",
			ents:    []client.EntitlementAPI{},
			args:    map[string]string{"name": "My Launcher"},
			summary: "No entitlement found",
			detail:  "up to an hour",
		},
		{
			name:    "ambiguous",
			ents:    []client.EntitlementAPI{{ID: "a"}, {ID: "b"}},
			args:    map[string]string{"name": "dup"},
			summary: "Multiple entitlements found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotFilters string
			ds := newTestDataSource(t, tc.ents, &gotFilters)
			resp := readWithConfig(t, ds, tc.args)
			if !resp.Diagnostics.HasError() {
				t.Fatalf("expected an error")
			}
			d := resp.Diagnostics.Errors()[0]
			if d.Summary() != tc.summary {
				t.Errorf("summary = %q, want %q", d.Summary(), tc.summary)
			}
			if tc.detail != "" && !strings.Contains(d.Detail(), tc.detail) {
				t.Errorf("detail %q does not mention %q", d.Detail(), tc.detail)
			}
		})
	}
}
