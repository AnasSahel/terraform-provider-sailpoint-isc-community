// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const syncConfigPath = "/sources/v1/src-1/attribute-sync-config"

// fakeTenant serves the sync config, the identity attributes, and records the PUT body.
// The GET returns a stale display name for email; the PUT drops the disabled "dropped" entry,
// like ISC does for entries the CREATE policy does not map.
type fakeTenant struct {
	mu      sync.Mutex
	putBody *client.SourceAttributeSyncConfigAPI
}

func (f *fakeTenant) serve(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/oauth/token":
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
		case r.URL.Path == "/identity-attributes/v1" && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`[{"name":"email","displayName":"Work Email"},{"name":"department","displayName":"Department"}]`))
		case r.URL.Path == syncConfigPath && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{"source":{"type":"SOURCE","id":"src-1","name":"AD"},"attributes":[` +
				`{"name":"email","displayName":"Old Email","target":"mail","enabled":false},` +
				`{"name":"department","displayName":"Department","target":"dept","enabled":true},` +
				`{"name":"dropped","displayName":"Dropped","target":"drop","enabled":true}]}`))
		case r.URL.Path == syncConfigPath && r.Method == http.MethodPut:
			var body client.SourceAttributeSyncConfigAPI
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode PUT body: %v", err)
			}
			f.mu.Lock()
			f.putBody = &body
			f.mu.Unlock()
			kept := body.Attributes[:0]
			for _, a := range body.Attributes {
				if a.Name == "dropped" && !a.Enabled {
					continue
				}
				kept = append(kept, a)
			}
			body.Attributes = kept
			out, _ := json.Marshal(body)
			_, _ = w.Write(out)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func newTestResource(t *testing.T) (*sourceAttrSyncConfigResource, *fakeTenant) {
	t.Helper()
	f := &fakeTenant{}
	srv := httptest.NewServer(f.serve(t))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &sourceAttrSyncConfigResource{client: c}, f
}

func mustSchema(t *testing.T) (schema.Schema, tftypes.Object) {
	t.Helper()
	var resp resource.SchemaResponse
	NewSourceAttributeSyncConfigResource().Schema(context.Background(), resource.SchemaRequest{}, &resp)
	objType, ok := resp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	return resp.Schema, objType
}

// TestCreate_UnknownSourceStaleDisplayNameDroppedEntry reproduces the three defects of #211:
// a plan with an unknown `source`, a stale displayName in the GET, and a disabled entry that
// ISC drops on PUT.
func TestCreate_UnknownSourceStaleDisplayNameDroppedEntry(t *testing.T) {
	r, f := newTestResource(t)
	sch, objType := mustSchema(t)

	attrsType, ok := objType.AttributeTypes["attributes"].(tftypes.Set)
	if !ok {
		t.Fatal("attributes is not a set")
	}
	elemType, ok := attrsType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatal("attributes element is not an object")
	}
	entry := func(name string, enabled bool) tftypes.Value {
		return tftypes.NewValue(elemType, map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, name),
			"enabled":      tftypes.NewValue(tftypes.Bool, enabled),
			"display_name": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"target":       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		})
	}
	planRaw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"source_id":  tftypes.NewValue(tftypes.String, "src-1"),
		"source":     tftypes.NewValue(objType.AttributeTypes["source"], tftypes.UnknownValue),
		"attributes": tftypes.NewValue(attrsType, []tftypes.Value{entry("email", true), entry("dropped", false)}),
	})
	plan := tfsdk.Plan{Schema: sch, Raw: planRaw}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create diagnostics: %v", resp.Diagnostics)
	}

	if f.putBody == nil {
		t.Fatal("no PUT was sent")
	}
	if f.putBody.Source.ID != "src-1" {
		t.Errorf("PUT source = %+v, want id src-1", f.putBody.Source)
	}
	for _, a := range f.putBody.Attributes {
		if a.Name == "email" && (a.DisplayName != "Work Email" || !a.Enabled) {
			t.Errorf("PUT email = %+v, want enabled with current display name", a)
		}
	}

	var state sourceAttrSyncConfigModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("state: %v", diags)
	}
	if state.Source.IsUnknown() || state.Source.IsNull() {
		t.Error("source must be known in state")
	}
	names := map[string]bool{}
	for _, a := range state.Attributes {
		names[a.Name.ValueString()] = a.Enabled.ValueBool()
	}
	if enabled, ok := names["dropped"]; !ok || enabled {
		t.Errorf("dropped entry should stay in state as disabled, got %v", names)
	}
	if !names["email"] {
		t.Errorf("email should be enabled in state, got %v", names)
	}
}
