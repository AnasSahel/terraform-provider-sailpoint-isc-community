// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// recordedCall is one non-OAuth request the fake tenant received.
type recordedCall struct {
	Method      string
	Path        string
	ContentType string
	Body        []byte
}

// fakeTenant answers OAuth and delegates the rest to handler, recording each non-OAuth call.
type fakeTenant struct {
	mu    sync.Mutex
	calls []recordedCall
}

func newTestResource(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, body []byte)) (*roleResource, *fakeTenant) {
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
		f.calls = append(f.calls, recordedCall{Method: r.Method, Path: r.URL.Path, ContentType: r.Header.Get("Content-Type"), Body: body})
		f.mu.Unlock()
		handler(w, r, body)
	}))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &roleResource{client: c}, f
}

func emptyRaw(t *testing.T) tftypes.Value {
	t.Helper()
	objType, ok := testSchema(t).Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	return tftypes.NewValue(objType, nil)
}

func planFor(t *testing.T, m *roleModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: testSchema(t), Raw: emptyRaw(t)}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("set plan: %v", diags)
	}
	return plan
}

func stateFor(t *testing.T, m *roleModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: testSchema(t), Raw: emptyRaw(t)}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
	return state
}

// createPlan is the plan for a new role with only name, owner and dimensional configured.
func createPlan() *roleModel {
	return &roleModel{
		ID:                  types.StringUnknown(),
		Name:                types.StringValue("Dynamic Role"),
		Description:         types.StringNull(),
		Enabled:             types.BoolUnknown(),
		Requestable:         types.BoolUnknown(),
		Dimensional:         types.BoolValue(true),
		Owner:               &common.ObjectRefModel{Type: types.StringValue("IDENTITY"), ID: types.StringValue("identity-1"), Name: types.StringUnknown()},
		Segments:            types.SetValueMust(types.StringType, []attr.Value{}),
		AccessRequestConfig: types.ObjectUnknown(accessRequestConfigAttrTypes),
		DimensionRefs:       types.SetUnknown(common.ObjectRefObjectType),
		Created:             types.StringUnknown(),
		Modified:            types.StringUnknown(),
	}
}

// echoCreate answers POST /roles/v1 with the sent body plus the fields the API fills in.
func echoCreate(t *testing.T, defaultARC map[string]any) func(w http.ResponseWriter, r *http.Request, body []byte) {
	return func(w http.ResponseWriter, r *http.Request, body []byte) {
		if r.Method != http.MethodPost || r.URL.Path != "/roles/v1" {
			t.Errorf("unexpected call %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		role := decodeMap(t, body)
		role["id"] = "role-1"
		role["enabled"] = false
		role["requestable"] = false
		role["segments"] = []string{}
		role["dimensionRefs"] = []any{}
		role["owner"] = map[string]any{"type": "IDENTITY", "id": "identity-1", "name": "Owner"}
		if arc, ok := role["accessRequestConfig"].(map[string]any); ok {
			for k, v := range defaultARC {
				if _, set := arc[k]; !set {
					arc[k] = v
				}
			}
		} else {
			role["accessRequestConfig"] = defaultARC
		}
		_ = json.NewEncoder(w).Encode(role)
	}
}

// apiDefaultARC is the access request config the API fills in when a role is created without one.
func apiDefaultARC() map[string]any {
	return map[string]any{
		"commentsRequired": nil, "denialCommentsRequired": nil, "approvalSchemes": []any{},
		"reauthorizationRequired": false, "requireEndDate": false, "formDefinitionId": nil,
	}
}

func runCreate(t *testing.T, r *roleResource, plan *roleModel) roleModel {
	t.Helper()
	ctx := context.Background()
	resp := resource.CreateResponse{State: tfsdk.State{Schema: testSchema(t), Raw: emptyRaw(t)}}
	r.Create(ctx, resource.CreateRequest{Plan: planFor(t, plan)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create: %v", resp.Diagnostics)
	}
	var got roleModel
	if diags := resp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("get state: %v", diags)
	}
	return got
}

func arcModel(t *testing.T, obj types.Object) accessRequestConfigModel {
	t.Helper()
	var m accessRequestConfigModel
	if diags := obj.As(context.Background(), &m, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("access_request_config: %v", diags)
	}
	return m
}

func dimensionAttrs(t *testing.T, obj types.Object) []dimensionAttributeModel {
	t.Helper()
	var m dimensionSchemaModel
	if diags := obj.As(context.Background(), &m, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("dimension_schema: %v", diags)
	}
	return m.DimensionAttributes
}

func TestCreateWithoutSegmentsOrAccessRequestConfig(t *testing.T) {
	r, f := newTestResource(t, echoCreate(t, apiDefaultARC()))
	got := runCreate(t, r, createPlan())

	if len(f.calls) != 1 {
		t.Fatalf("calls = %d, want 1", len(f.calls))
	}
	sent := decodeMap(t, f.calls[0].Body)
	for _, key := range []string{"accessRequestConfig", "dimensionRefs"} {
		if _, found := sent[key]; found {
			t.Errorf("create body = %s, want no %s", f.calls[0].Body, key)
		}
	}
	if sent["dimensional"] != true {
		t.Errorf("create body = %s, want dimensional true", f.calls[0].Body)
	}

	if got.Segments.IsNull() || got.Segments.IsUnknown() || len(got.Segments.Elements()) != 0 {
		t.Errorf("segments = %v, want a known empty set", got.Segments)
	}
	if got.AccessRequestConfig.IsNull() || got.AccessRequestConfig.IsUnknown() {
		t.Fatalf("access_request_config = %v, want the API defaults", got.AccessRequestConfig)
	}
	arc := arcModel(t, got.AccessRequestConfig)
	if !arc.ReauthorizationRequired.Equal(types.BoolValue(false)) || !arc.DimensionSchema.IsNull() || !arc.FormDefinitionID.IsNull() {
		t.Errorf("access_request_config = %+v, want reauthorization_required false and no schema or form", arc)
	}
}

func TestCreateWithDimensionSchemaAndForm(t *testing.T) {
	ctx := context.Background()
	r, f := newTestResource(t, echoCreate(t, apiDefaultARC()))

	ds, diags := types.ObjectValueFrom(ctx, dimensionSchemaAttrTypes, dimensionSchemaModel{DimensionAttributes: []dimensionAttributeModel{
		{Name: types.StringValue("department"), DisplayName: types.StringValue("Department"), Derived: types.BoolValue(true)},
		{Name: types.StringValue("location"), DisplayName: types.StringUnknown(), Derived: types.BoolValue(false)},
	}})
	if diags.HasError() {
		t.Fatalf("dimension_schema: %v", diags)
	}
	arc, diags := types.ObjectValue(accessRequestConfigAttrTypes, map[string]attr.Value{
		"comments_required":             types.BoolUnknown(),
		"denial_comments_required":      types.BoolUnknown(),
		"reauthorization_required":      types.BoolUnknown(),
		"require_end_date":              types.BoolUnknown(),
		"max_permitted_access_duration": types.ObjectNull(accessDurationAttrTypes),
		"approval_schemes":              types.ListNull(types.ObjectType{AttrTypes: approvalSchemeAttrTypes}),
		"dimension_schema":              ds,
		"form_definition_id":            types.StringValue("form-def-1"),
	})
	if diags.HasError() {
		t.Fatalf("access_request_config: %v", diags)
	}
	plan := createPlan()
	plan.AccessRequestConfig = arc
	got := runCreate(t, r, plan)

	var sent client.RoleAPI
	if err := json.Unmarshal(f.calls[0].Body, &sent); err != nil {
		t.Fatalf("unmarshal create body: %v", err)
	}
	sentARC := sent.AccessRequestConfig
	if sentARC == nil || sentARC.FormDefinitionID == nil || *sentARC.FormDefinitionID != "form-def-1" ||
		sentARC.DimensionSchema == nil || len(sentARC.DimensionSchema.DimensionAttributes) != 2 {
		t.Fatalf("create body = %s, want both dimension attributes and the form", f.calls[0].Body)
	}
	if a := sentARC.DimensionSchema.DimensionAttributes[1]; a.Name != "location" || a.Derived == nil || *a.Derived || a.DisplayName != nil {
		t.Errorf("second dimension attribute = %+v, want location, derived false, no display name", a)
	}

	gotARC := arcModel(t, got.AccessRequestConfig)
	if !gotARC.FormDefinitionID.Equal(types.StringValue("form-def-1")) {
		t.Errorf("form_definition_id = %v, want form-def-1", gotARC.FormDefinitionID)
	}
	attrs := dimensionAttrs(t, gotARC.DimensionSchema)
	if len(attrs) != 2 || !attrs[0].Name.Equal(types.StringValue("department")) ||
		!attrs[1].Derived.Equal(types.BoolValue(false)) {
		t.Errorf("dimension_attributes = %+v, want department then location (derived false)", attrs)
	}
}

func TestUpdateKeepsDimensionSchemaInPatchBody(t *testing.T) {
	ctx := context.Background()
	updated := map[string]any{}
	if err := json.Unmarshal([]byte(dynamicRoleJSON), &updated); err != nil {
		t.Fatal(err)
	}
	arcJSON, ok := updated["accessRequestConfig"].(map[string]any)
	if !ok {
		t.Fatal("fixture has no accessRequestConfig object")
	}
	arcJSON["commentsRequired"] = true

	r, f := newTestResource(t, func(w http.ResponseWriter, r *http.Request, _ []byte) {
		if r.Method != http.MethodPatch || r.URL.Path != "/roles/v1/role-1" {
			t.Errorf("unexpected call %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(updated)
	})

	state := modelFromJSON(t, dynamicRoleJSON)
	// The plan for a config that sets comments_required = true and leaves dimension_schema and
	// form_definition_id out: UseStateForUnknown has copied them from state.
	plan := modelFromJSON(t, withJSON(t, map[string]any{
		"dimensional":         false,
		"accessRequestConfig": arcJSON,
		"dimensionRefs":       []any{},
	}))

	resp := resource.UpdateResponse{State: stateFor(t, &state)}
	r.Update(ctx, resource.UpdateRequest{Plan: planFor(t, &plan), State: stateFor(t, &state)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update: %v", resp.Diagnostics)
	}

	if len(f.calls) != 1 {
		t.Fatalf("calls = %d, want 1 PATCH", len(f.calls))
	}
	call := f.calls[0]
	if call.ContentType != "application/json-patch+json" {
		t.Errorf("Content-Type = %q, want application/json-patch+json", call.ContentType)
	}
	var ops []struct {
		Op    string                          `json:"op"`
		Path  string                          `json:"path"`
		Value client.RequestabilityForRoleAPI `json:"value"`
	}
	if err := json.Unmarshal(call.Body, &ops); err != nil {
		t.Fatalf("unmarshal PATCH body %s: %v", call.Body, err)
	}
	if len(ops) != 1 || ops[0].Op != "replace" || ops[0].Path != "/accessRequestConfig" {
		t.Fatalf("PATCH body = %s, want only replace /accessRequestConfig", call.Body)
	}
	v := ops[0].Value
	if v.DimensionSchema == nil || len(v.DimensionSchema.DimensionAttributes) != 2 ||
		v.FormDefinitionID == nil || *v.FormDefinitionID != "form-def-1" {
		t.Errorf("PATCH body = %s, want the dimension schema and form kept", call.Body)
	}

	var got roleModel
	if diags := resp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("get state: %v", diags)
	}
	if attrs := dimensionAttrs(t, arcModel(t, got.AccessRequestConfig).DimensionSchema); len(attrs) != 2 {
		t.Errorf("dimension_attributes = %+v, want 2", attrs)
	}
}
