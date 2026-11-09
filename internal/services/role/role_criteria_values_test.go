// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// cvUIRoleJSON is a role as the SailPoint UI saves it: every leaf uses values
// with stringValue null, and AND/OR nodes have values null.
const cvUIRoleJSON = `{
  "id": "role-1",
  "name": "UI Role",
  "owner": {"type": "IDENTITY", "id": "identity-1", "name": "Owner"},
  "enabled": false,
  "requestable": false,
  "dimensional": false,
  "segments": [],
  "accessRequestConfig": {"commentsRequired": false, "denialCommentsRequired": false,
    "reauthorizationRequired": false, "requireEndDate": false, "approvalSchemes": []},
  "membership": {"type": "STANDARD", "criteria": {
    "operation": "AND", "key": null, "stringValue": null, "values": null, "children": [
      {"operation": "EQUALS", "key": {"type": "IDENTITY", "property": "attribute.department", "sourceId": null},
       "stringValue": null, "values": ["ZZ-a", "ZZ-b"], "children": null},
      {"operation": "NOT_EQUALS", "key": {"type": "IDENTITY", "property": "attribute.location", "sourceId": null},
       "stringValue": null, "values": ["ZZ-c"], "children": null},
      {"operation": "OR", "key": null, "stringValue": null, "values": null, "children": [
        {"operation": "EQUALS", "key": {"type": "IDENTITY", "property": "attribute.x"},
         "stringValue": null, "values": ["ZZ-d"]},
        {"operation": "EQUALS", "key": {"type": "IDENTITY", "property": "attribute.y"},
         "stringValue": null, "values": ["ZZ-e", "ZZ-f"]}
      ]}
    ]}},
  "created": "2026-01-01T00:00:00Z",
  "modified": "2026-01-01T00:00:00Z"
}`

func cvResourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	var resp resource.SchemaResponse
	NewRoleResource().Schema(t.Context(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("resource schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func cvDataSourceSchema(t *testing.T) dsschema.Schema {
	t.Helper()
	var resp datasource.SchemaResponse
	NewRoleDataSource().Schema(t.Context(), datasource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// cvNullObject is a value of objType with every attribute null except vals.
func cvNullObject(objType tftypes.Object, vals map[string]tftypes.Value) tftypes.Value {
	all := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		if v, ok := vals[name]; ok {
			all[name] = v
		} else {
			all[name] = tftypes.NewValue(typ, nil)
		}
	}
	return tftypes.NewValue(objType, all)
}

func cvObjectType(t *testing.T, typ attr.Type) tftypes.Object {
	t.Helper()
	objType, ok := typ.TerraformType(t.Context()).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}
	return objType
}

func cvModelFromJSON(t *testing.T, raw string) roleModel {
	t.Helper()
	var api client.RoleAPI
	if err := json.Unmarshal([]byte(raw), &api); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var m roleModel
	if diags := m.FromAPI(t.Context(), &api); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	return m
}

func cvList(vals ...string) types.List {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, elems)
}

func cvListValues(t *testing.T, l types.List) []string {
	t.Helper()
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	return criteriaValuesToAPI(l)
}

// cvCall is one non-OAuth request the fake server received.
type cvCall struct {
	Method      string
	Path        string
	ContentType string
	Body        []byte
}

type cvRecorder struct {
	mu    sync.Mutex
	calls []cvCall
}

// cvNewServer answers OAuth and passes every other request to handler.
func cvNewServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request, body []byte)) (*client.Client, *cvRecorder) {
	t.Helper()
	rec := &cvRecorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/oauth/token" {
			_, _ = w.Write([]byte(`{"access_token":"test","expires_in":3600,"token_type":"bearer"}`))
			return
		}
		body, _ := io.ReadAll(r.Body)
		rec.mu.Lock()
		rec.calls = append(rec.calls, cvCall{Method: r.Method, Path: r.URL.Path, ContentType: r.Header.Get("Content-Type"), Body: body})
		rec.mu.Unlock()
		handler(w, r, body)
	}))
	t.Cleanup(srv.Close)
	c, err := client.NewClient(srv.URL, "id", "secret")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c, rec
}

func cvAssertNullList(t *testing.T, name string, l types.List) {
	t.Helper()
	if !l.Equal(types.ListNull(types.StringType)) {
		t.Errorf("%s: want ListNull(String), got %#v", name, l)
	}
}

func TestCriteriaValues_UIRoleRoundTrip(t *testing.T) {
	m := cvModelFromJSON(t, cvUIRoleJSON)
	c := m.Membership.Criteria

	cvAssertNullList(t, "AND values", c.Values)
	if !c.StringValue.IsNull() {
		t.Errorf("AND string_value: want null, got %v", c.StringValue)
	}
	leaves := []struct {
		name string
		sv   types.String
		vals types.List
		want []string
	}{
		{"children[0]", c.Children[0].StringValue, c.Children[0].Values, []string{"ZZ-a", "ZZ-b"}},
		{"children[1]", c.Children[1].StringValue, c.Children[1].Values, []string{"ZZ-c"}},
		{"children[2].children[0]", c.Children[2].Children[0].StringValue, c.Children[2].Children[0].Values, []string{"ZZ-d"}},
		{"children[2].children[1]", c.Children[2].Children[1].StringValue, c.Children[2].Children[1].Values, []string{"ZZ-e", "ZZ-f"}},
	}
	for _, l := range leaves {
		if !l.sv.IsNull() {
			t.Errorf("%s string_value: want null, got %v", l.name, l.sv)
		}
		if got := cvListValues(t, l.vals); !reflect.DeepEqual(got, l.want) {
			t.Errorf("%s values: want %v, got %v", l.name, l.want, got)
		}
	}
	cvAssertNullList(t, "OR values", c.Children[2].Values)
	if !c.Children[2].StringValue.IsNull() {
		t.Errorf("OR string_value: want null, got %v", c.Children[2].StringValue)
	}

	// ToAPI sends values on leaves, no stringValue anywhere and no values on AND/OR.
	api, diags := m.ToAPI(t.Context())
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags)
	}
	body, err := json.Marshal(api.Membership.Criteria)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(body), "stringValue") {
		t.Errorf("criteria body has stringValue: %s", body)
	}
	var tree map[string]any
	if err := json.Unmarshal(body, &tree); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := tree["values"]; ok {
		t.Errorf("AND node has values: %s", body)
	}
	children, ok := tree["children"].([]any)
	if !ok || len(children) != 3 {
		t.Fatalf("want 3 children: %s", body)
	}
	orNode, ok := children[2].(map[string]any)
	if !ok {
		t.Fatalf("OR node is not an object: %s", body)
	}
	if _, ok := orNode["values"]; ok {
		t.Errorf("OR node has values: %s", body)
	}
	leaf0, ok := children[0].(map[string]any)
	if !ok || !reflect.DeepEqual(leaf0["values"], []any{"ZZ-a", "ZZ-b"}) {
		t.Errorf("leaf 0 values: %s", body)
	}

	// State round trip keeps membership identical and plans no membership change.
	sch := cvResourceSchema(t)
	state := tfsdk.State{Schema: sch, Raw: tftypes.NewValue(cvObjectType(t, sch.Type()), nil)}
	if diags := state.Set(t.Context(), &m); diags.HasError() {
		t.Fatalf("State.Set: %v", diags)
	}
	var got roleModel
	if diags := state.Get(t.Context(), &got); diags.HasError() {
		t.Fatalf("State.Get: %v", diags)
	}
	if !reflect.DeepEqual(got.Membership, m.Membership) {
		t.Errorf("membership changed through state:\nwant %#v\ngot  %#v", m.Membership, got.Membership)
	}
	ops, diags := got.ToPatchOperations(t.Context(), &m)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	for _, op := range ops {
		if op.Path == "/membership" {
			t.Errorf("unexpected membership patch: %+v", op)
		}
	}
}

func TestCriteriaValues_MixedTree(t *testing.T) {
	m := cvModelFromJSON(t, `{"id":"role-1","name":"Mixed","owner":{"type":"IDENTITY","id":"identity-1"},
	  "membership":{"type":"STANDARD","criteria":{"operation":"OR","children":[
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.a"},"stringValue":"ZZ-x","values":null},
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.b"},"stringValue":null,"values":["ZZ-b","ZZ-a"]},
	    {"operation":"AND","children":[
	      {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.c"},"stringValue":"ZZ-y"},
	      {"operation":"NOT_EQUALS","key":{"type":"IDENTITY","property":"attribute.d"},"values":["ZZ-d2","ZZ-d1"]}
	    ]}
	  ]}}}`)
	c := m.Membership.Criteria
	if c.Children[0].StringValue.ValueString() != "ZZ-x" {
		t.Errorf("level-2 string_value: got %v", c.Children[0].StringValue)
	}
	cvAssertNullList(t, "level-2 string leaf values", c.Children[0].Values)
	if got := cvListValues(t, c.Children[1].Values); !reflect.DeepEqual(got, []string{"ZZ-b", "ZZ-a"}) {
		t.Errorf("level-2 values order: got %v", got)
	}
	if c.Children[2].Children[0].StringValue.ValueString() != "ZZ-y" {
		t.Errorf("level-3 string_value: got %v", c.Children[2].Children[0].StringValue)
	}
	if got := cvListValues(t, c.Children[2].Children[1].Values); !reflect.DeepEqual(got, []string{"ZZ-d2", "ZZ-d1"}) {
		t.Errorf("level-3 values order: got %v", got)
	}

	api := criteriaToAPI(c)
	if api.Children[0].StringValue == nil || *api.Children[0].StringValue != "ZZ-x" || api.Children[0].Values != nil {
		t.Errorf("level-2 string leaf to API: %+v", api.Children[0])
	}
	if api.Children[1].StringValue != nil || !reflect.DeepEqual(api.Children[1].Values, []string{"ZZ-b", "ZZ-a"}) {
		t.Errorf("level-2 values leaf to API: %+v", api.Children[1])
	}
	if api.Children[2].Children[0].StringValue == nil || api.Children[2].Children[0].Values != nil {
		t.Errorf("level-3 string leaf to API: %+v", api.Children[2].Children[0])
	}
	if api.Children[2].Children[1].StringValue != nil ||
		!reflect.DeepEqual(api.Children[2].Children[1].Values, []string{"ZZ-d2", "ZZ-d1"}) {
		t.Errorf("level-3 values leaf to API: %+v", api.Children[2].Children[1])
	}
}

func TestCriteriaValues_EmptyListReadsNull(t *testing.T) {
	m := cvModelFromJSON(t, `{"id":"role-1","name":"Empty","owner":{"type":"IDENTITY","id":"identity-1"},
	  "membership":{"type":"STANDARD","criteria":{"operation":"AND","values":[],"children":[
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.a"},"stringValue":"ZZ-x","values":[]}
	  ]}}}`)
	c := m.Membership.Criteria
	cvAssertNullList(t, "AND values", c.Values)
	if c.Children[0].StringValue.ValueString() != "ZZ-x" {
		t.Errorf("leaf string_value: got %v", c.Children[0].StringValue)
	}
	cvAssertNullList(t, "leaf values", c.Children[0].Values)
}

func TestCriteriaValues_Create(t *testing.T) {
	c, rec := cvNewServer(t, func(w http.ResponseWriter, r *http.Request, body []byte) {
		if r.Method != http.MethodPost || r.URL.Path != "/roles/v1" {
			t.Errorf("unexpected call %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var role map[string]any
		if err := json.Unmarshal(body, &role); err != nil {
			t.Errorf("decode: %v", err)
		}
		role["id"] = "role-1"
		role["created"] = "2026-01-01T00:00:00Z"
		role["modified"] = "2026-01-01T00:00:00Z"
		role["owner"] = map[string]any{"type": "IDENTITY", "id": "identity-1", "name": "Owner"}
		_ = json.NewEncoder(w).Encode(role)
	})
	plan := cvModelFromJSON(t, `{"name":"Create","owner":{"type":"IDENTITY","id":"identity-1"},
	  "enabled":false,"requestable":false,"dimensional":false,"segments":[],
	  "membership":{"type":"STANDARD","criteria":{"operation":"AND","children":[
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"values":["ZZ-b","ZZ-a"]},
	    {"operation":"OR","children":[
	      {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"values":["ZZ-c"]}
	    ]}
	  ]}}}`)

	r := &roleResource{client: c}
	sch := cvResourceSchema(t)
	nullRaw := tftypes.NewValue(cvObjectType(t, sch.Type()), nil)
	tfPlan := tfsdk.Plan{Schema: sch, Raw: nullRaw}
	if diags := tfPlan.Set(t.Context(), &plan); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sch, Raw: nullRaw}}
	r.Create(t.Context(), resource.CreateRequest{Plan: tfPlan}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create: %v", resp.Diagnostics)
	}

	if len(rec.calls) != 1 {
		t.Fatalf("want 1 call, got %d", len(rec.calls))
	}
	var sent client.RoleAPI
	if err := json.Unmarshal(rec.calls[0].Body, &sent); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}
	if strings.Contains(string(rec.calls[0].Body), "stringValue") {
		t.Errorf("POST body has stringValue: %s", rec.calls[0].Body)
	}
	crit := sent.Membership.Criteria
	if !reflect.DeepEqual(crit.Children[0].Values, []string{"ZZ-b", "ZZ-a"}) {
		t.Errorf("POST level-2 values: got %v", crit.Children[0].Values)
	}
	if !reflect.DeepEqual(crit.Children[1].Children[0].Values, []string{"ZZ-c"}) {
		t.Errorf("POST level-3 values: got %v", crit.Children[1].Children[0].Values)
	}

	var got types.List
	if diags := resp.State.GetAttribute(t.Context(), path.Root("membership").AtName("criteria").AtName("children").AtListIndex(0).AtName("values"), &got); diags.HasError() {
		t.Fatalf("GetAttribute: %v", diags)
	}
	if v := cvListValues(t, got); !reflect.DeepEqual(v, []string{"ZZ-b", "ZZ-a"}) {
		t.Errorf("state level-2 values: got %v", v)
	}
	if diags := resp.State.GetAttribute(t.Context(), path.Root("membership").AtName("criteria").AtName("children").AtListIndex(1).AtName("children").AtListIndex(0).AtName("values"), &got); diags.HasError() {
		t.Fatalf("GetAttribute: %v", diags)
	}
	if v := cvListValues(t, got); !reflect.DeepEqual(v, []string{"ZZ-c"}) {
		t.Errorf("state level-3 values: got %v", v)
	}
}

func TestCriteriaValues_PatchSendsValues(t *testing.T) {
	const tmpl = `{"id":"role-1","name":"Patch","owner":{"type":"IDENTITY","id":"identity-1","name":"Owner"},
	  "membership":{"type":"STANDARD","criteria":{"operation":"AND","children":[
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},LEAF}
	  ]}}}`
	state := cvModelFromJSON(t, strings.Replace(tmpl, "LEAF", `"stringValue":"ZZ-a"`, 1))
	plan := cvModelFromJSON(t, strings.Replace(tmpl, "LEAF", `"values":["ZZ-a","ZZ-b"]`, 1))

	ops, diags := plan.ToPatchOperations(t.Context(), &state)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	var found bool
	for _, op := range ops {
		if op.Path != "/membership" {
			continue
		}
		found = true
		if op.Op != "replace" {
			t.Errorf("membership op: want replace, got %s", op.Op)
		}
		body, err := json.Marshal(op.Value)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !strings.Contains(string(body), `"values":["ZZ-a","ZZ-b"]`) || strings.Contains(string(body), "stringValue") {
			t.Errorf("membership patch body: %s", body)
		}
	}
	if !found {
		t.Fatalf("no /membership op in %+v", ops)
	}
}

// TestCriteriaValues_UpdateRequest pins the membership update request that
// SailPoint accepts: PATCH /roles/v1/{id}, json-patch content type, one
// replace of /membership with values leaves and no null fields.
func TestCriteriaValues_UpdateRequest(t *testing.T) {
	const tmpl = `{"id":"role-1","name":"Update","owner":{"type":"IDENTITY","id":"identity-1","name":"Owner"},
	  "enabled":false,"requestable":false,"dimensional":false,"segments":[],
	  "created":"2026-01-01T00:00:00Z","modified":"2026-01-01T00:00:00Z",
	  "membership":{"type":"STANDARD","criteria":{"operation":"AND","children":[
	    {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"values":LEAF1},
	    {"operation":"NOT_EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"values":LEAF2}
	  ]}}}`
	stateJSON := strings.NewReplacer("LEAF1", `["ZZ-a","ZZ-b"]`, "LEAF2", `["ZZ-c"]`).Replace(tmpl)
	planJSON := strings.NewReplacer("LEAF1", `["ZZ-b","ZZ-d","ZZ-a"]`, "LEAF2", `["ZZ-e"]`).Replace(tmpl)

	c, rec := cvNewServer(t, func(w http.ResponseWriter, r *http.Request, _ []byte) {
		if r.Method != http.MethodPatch || r.URL.Path != "/roles/v1/role-1" {
			t.Errorf("unexpected call %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(planJSON))
	})

	stateModel := cvModelFromJSON(t, stateJSON)
	planModel := cvModelFromJSON(t, planJSON)
	sch := cvResourceSchema(t)
	nullRaw := tftypes.NewValue(cvObjectType(t, sch.Type()), nil)
	tfState := tfsdk.State{Schema: sch, Raw: nullRaw}
	if diags := tfState.Set(t.Context(), &stateModel); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	tfPlan := tfsdk.Plan{Schema: sch, Raw: nullRaw}
	if diags := tfPlan.Set(t.Context(), &planModel); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: nullRaw}}
	(&roleResource{client: c}).Update(t.Context(), resource.UpdateRequest{Plan: tfPlan, State: tfState}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update: %v", resp.Diagnostics)
	}

	if len(rec.calls) != 1 {
		t.Fatalf("want 1 call, got %d", len(rec.calls))
	}
	call := rec.calls[0]
	if call.ContentType != "application/json-patch+json" {
		t.Errorf("Content-Type: want application/json-patch+json, got %q", call.ContentType)
	}
	if strings.Contains(string(call.Body), "null") || strings.Contains(string(call.Body), "stringValue") {
		t.Errorf("PATCH body has null or stringValue: %s", call.Body)
	}
	const wantBody = `[{"op":"replace","path":"/membership","value":{"type":"STANDARD","criteria":{"operation":"AND","children":[
	  {"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"values":["ZZ-b","ZZ-d","ZZ-a"]},
	  {"operation":"NOT_EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"values":["ZZ-e"]}
	]}}}]`
	var got, want any
	if err := json.Unmarshal(call.Body, &got); err != nil {
		t.Fatalf("decode sent body: %v", err)
	}
	if err := json.Unmarshal([]byte(wantBody), &want); err != nil {
		t.Fatalf("decode want body: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PATCH body:\nwant %s\ngot  %s", wantBody, call.Body)
	}

	var vals types.List
	if diags := resp.State.GetAttribute(t.Context(), path.Root("membership").AtName("criteria").AtName("children").AtListIndex(0).AtName("values"), &vals); diags.HasError() {
		t.Fatalf("GetAttribute: %v", diags)
	}
	if v := cvListValues(t, vals); !reflect.DeepEqual(v, []string{"ZZ-b", "ZZ-d", "ZZ-a"}) {
		t.Errorf("state values after update: got %v", v)
	}
}

// cvNode builds a criteria node for ValidateConfig tests.
type cvNode struct {
	op   string
	sv   types.String
	vals types.List
}

func cvLevel3(n cvNode) criteriaLevel3Model {
	return criteriaLevel3Model{Operation: types.StringValue(n.op), StringValue: n.sv, Values: n.vals}
}

func cvLevel2(n cvNode, children ...criteriaLevel3Model) criteriaLevel2Model {
	return criteriaLevel2Model{Operation: types.StringValue(n.op), StringValue: n.sv, Values: n.vals, Children: children}
}

func cvRoot(n cvNode, children ...criteriaLevel2Model) *roleCriteriaModel {
	return &roleCriteriaModel{Operation: types.StringValue(n.op), StringValue: n.sv, Values: n.vals, Children: children}
}

func TestCriteriaValues_ValidateConfig(t *testing.T) {
	nullSV := types.StringNull()
	nullVals := types.ListNull(types.StringType)
	and := cvNode{"AND", nullSV, nullVals}
	or := cvNode{"OR", nullSV, nullVals}
	leafVals := cvNode{"EQUALS", nullSV, cvList("ZZ-a", "ZZ-b")}
	leafSV := cvNode{"EQUALS", types.StringValue("ZZ-a"), nullVals}
	leafBoth := cvNode{"EQUALS", types.StringValue("ZZ-a"), cvList("ZZ-a")}

	cases := []struct {
		name     string
		criteria *roleCriteriaModel
		wantErr  string
		wantPath string
	}{
		{name: "leaf values", criteria: cvRoot(and, cvLevel2(leafVals))},
		{name: "leaf string_value", criteria: cvRoot(and, cvLevel2(leafSV))},
		{name: "mixed three levels", criteria: cvRoot(and, cvLevel2(leafSV), cvLevel2(or, cvLevel3(leafVals), cvLevel3(leafSV)))},
		{name: "leaf both", criteria: cvRoot(and, cvLevel2(leafBoth)), wantErr: "not both", wantPath: "membership.criteria.children[0]"},
		{name: "leaf neither", criteria: cvRoot(and, cvLevel2(cvNode{"EQUALS", nullSV, nullVals})), wantErr: "Set `string_value` or `values`.", wantPath: "membership.criteria.children[0]"},
		{name: "empty values", criteria: cvRoot(and, cvLevel2(cvNode{"EQUALS", nullSV, cvList()})), wantErr: "at least one value", wantPath: "membership.criteria.children[0].values"},
		{name: "blank entry", criteria: cvRoot(and, cvLevel2(cvNode{"EQUALS", nullSV, cvList("ZZ-a", "  ")})), wantErr: "must not be blank", wantPath: "membership.criteria.children[0].values[1]"},
		{name: "values on AND", criteria: cvRoot(cvNode{"AND", nullSV, cvList("ZZ-a")}, cvLevel2(leafSV)), wantErr: "`values` is not allowed on an AND/OR node.", wantPath: "membership.criteria.values"},
		{name: "string_value on OR", criteria: cvRoot(and, cvLevel2(cvNode{"OR", types.StringValue("ZZ-a"), nullVals}, cvLevel3(leafSV))), wantErr: "`string_value` is not allowed on an AND/OR node.", wantPath: "membership.criteria.children[0].string_value"},
		{name: "level-3 both", criteria: cvRoot(and, cvLevel2(leafSV), cvLevel2(or, cvLevel3(leafBoth))), wantErr: "not both", wantPath: "membership.criteria.children[1].children[0]"},
		{name: "unknown values", criteria: cvRoot(and, cvLevel2(cvNode{"EQUALS", nullSV, types.ListUnknown(types.StringType)}))},
		{name: "unknown element", criteria: cvRoot(and, cvLevel2(cvNode{"EQUALS", nullSV, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("ZZ-a"), types.StringUnknown()})}))},
	}

	sch := cvResourceSchema(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := tfsdk.State{Schema: sch, Raw: cvNullObject(cvObjectType(t, sch.Type()), nil)}
			m := membershipModel{Type: types.StringValue("STANDARD"), Criteria: tc.criteria}
			if diags := state.SetAttribute(t.Context(), path.Root("membership"), m); diags.HasError() {
				t.Fatalf("SetAttribute: %v", diags)
			}
			resp := resource.ValidateConfigResponse{}
			(&roleResource{}).ValidateConfig(t.Context(), resource.ValidateConfigRequest{
				Config: tfsdk.Config{Schema: sch, Raw: state.Raw},
			}, &resp)

			if tc.wantErr == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected errors: %v", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.ErrorsCount() != 1 {
				t.Fatalf("want 1 error, got %v", resp.Diagnostics)
			}
			d := resp.Diagnostics.Errors()[0]
			if !strings.Contains(d.Detail(), tc.wantErr) {
				t.Errorf("detail: want %q, got %q", tc.wantErr, d.Detail())
			}
			withPath, ok := d.(interface{ Path() path.Path })
			if !ok {
				t.Fatalf("diagnostic has no path: %v", d)
			}
			if got := withPath.Path().String(); got != tc.wantPath {
				t.Errorf("path: want %s, got %s", tc.wantPath, got)
			}
		})
	}
}

func TestCriteriaValues_DataSource(t *testing.T) {
	c, _ := cvNewServer(t, func(w http.ResponseWriter, r *http.Request, _ []byte) {
		if r.Method != http.MethodGet || r.URL.Path != "/roles/v1/role-1" {
			t.Errorf("unexpected call %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(cvUIRoleJSON))
	})
	sch := cvDataSourceSchema(t)
	objType := cvObjectType(t, sch.Type())
	config := tfsdk.Config{Schema: sch, Raw: cvNullObject(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "role-1"),
	})}
	resp := datasource.ReadResponse{State: tfsdk.State{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	(&roleDataSource{client: c}).Read(t.Context(), datasource.ReadRequest{Config: config}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read: %v", resp.Diagnostics)
	}

	leaf := path.Root("membership").AtName("criteria").AtName("children").AtListIndex(0)
	var vals types.List
	if diags := resp.State.GetAttribute(t.Context(), leaf.AtName("values"), &vals); diags.HasError() {
		t.Fatalf("GetAttribute values: %v", diags)
	}
	if got := cvListValues(t, vals); !reflect.DeepEqual(got, []string{"ZZ-a", "ZZ-b"}) {
		t.Errorf("values: got %v", got)
	}
	var sv types.String
	if diags := resp.State.GetAttribute(t.Context(), leaf.AtName("string_value"), &sv); diags.HasError() {
		t.Fatalf("GetAttribute string_value: %v", diags)
	}
	if !sv.IsNull() {
		t.Errorf("string_value: want null, got %v", sv)
	}
	leaf3 := path.Root("membership").AtName("criteria").AtName("children").AtListIndex(2).AtName("children").AtListIndex(1).AtName("values")
	if diags := resp.State.GetAttribute(t.Context(), leaf3, &vals); diags.HasError() {
		t.Fatalf("GetAttribute level-3 values: %v", diags)
	}
	if got := cvListValues(t, vals); !reflect.DeepEqual(got, []string{"ZZ-e", "ZZ-f"}) {
		t.Errorf("level-3 values: got %v", got)
	}
}

func TestCriteriaValues_Schema(t *testing.T) {
	rs := cvResourceSchema(t)
	ds := cvDataSourceSchema(t)
	if diags := rs.ValidateImplementation(t.Context()); diags.HasError() {
		t.Errorf("resource schema: %v", diags)
	}
	if diags := ds.ValidateImplementation(t.Context()); diags.HasError() {
		t.Errorf("data source schema: %v", diags)
	}
	crit := path.Root("membership").AtName("criteria")
	paths := []path.Path{
		crit.AtName("values"),
		crit.AtName("children").AtListIndex(0).AtName("values"),
		crit.AtName("children").AtListIndex(0).AtName("children").AtListIndex(0).AtName("values"),
	}
	want := types.ListType{ElemType: types.StringType}
	for _, p := range paths {
		rt, diags := rs.TypeAtPath(t.Context(), p)
		if diags.HasError() {
			t.Errorf("resource %s: %v", p, diags)
		} else if !rt.Equal(want) {
			t.Errorf("resource %s: want %v, got %v", p, want, rt)
		}
		dt, diags := ds.TypeAtPath(t.Context(), p)
		if diags.HasError() {
			t.Errorf("data source %s: %v", p, diags)
		} else if !dt.Equal(want) {
			t.Errorf("data source %s: want %v, got %v", p, want, dt)
		}
	}
}
