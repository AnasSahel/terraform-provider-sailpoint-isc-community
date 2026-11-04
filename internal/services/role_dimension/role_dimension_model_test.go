// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// refSet builds a set of refs; an empty name argument means an unknown name, as in a plan.
func refSet(t *testing.T, refs ...[3]string) types.Set {
	t.Helper()
	models := make([]common.ObjectRefModel, 0, len(refs))
	for _, r := range refs {
		name := types.StringUnknown()
		if r[2] != "" {
			name = types.StringValue(r[2])
		}
		models = append(models, common.ObjectRefModel{Type: types.StringValue(r[0]), ID: types.StringValue(r[1]), Name: name})
	}
	set, diags := types.SetValueFrom(context.Background(), common.ObjectRefObjectType, models)
	if diags.HasError() {
		t.Fatalf("set: %v", diags)
	}
	return set
}

func decodeDimension(t *testing.T, body string) *client.RoleDimensionAPI {
	t.Helper()
	var api client.RoleDimensionAPI
	if err := json.Unmarshal([]byte(body), &api); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return &api
}

func fromAPI(t *testing.T, body string, prior *roleDimensionModel) roleDimensionModel {
	t.Helper()
	var m roleDimensionModel
	if diags := m.FromAPI(context.Background(), "r1", decodeDimension(t, body), prior); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	return m
}

func TestFromAPIAccessProfileOrderIgnored(t *testing.T) {
	a := fromAPI(t, `{"id":"d1","name":"n","accessProfiles":[
		{"type":"ACCESS_PROFILE","id":"ap1","name":"AP 1"},{"type":"ACCESS_PROFILE","id":"ap2","name":"AP 2"}]}`, nil)
	b := fromAPI(t, `{"id":"d1","name":"n","accessProfiles":[
		{"type":"ACCESS_PROFILE","id":"ap2","name":"AP 2"},{"type":"ACCESS_PROFILE","id":"ap1","name":"AP 1"}]}`, nil)
	if !a.AccessProfiles.Equal(b.AccessProfiles) {
		t.Errorf("sets differ: %v vs %v", a.AccessProfiles, b.AccessProfiles)
	}
	if a.RoleID.ValueString() != "r1" || a.ID.ValueString() != "d1" {
		t.Errorf("role_id/id = %s/%s", a.RoleID, a.ID)
	}
}

func TestFromAPINormalisesEmptyConnectorValues(t *testing.T) {
	m := fromAPI(t, `{"id":"d1","name":"n","membership":{"type":"STANDARD","criteria":{
		"operation":"AND","stringValue":"","children":[
			{"operation":"OR","stringValue":"","children":[
				{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"stringValue":""}]},
			{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"stringValue":"Engineering"}]}}}`, nil)
	c := m.Membership.Criteria
	if !c.StringValue.IsNull() {
		t.Errorf("AND string_value = %v, want null", c.StringValue)
	}
	if !c.Children[0].StringValue.IsNull() {
		t.Errorf("OR string_value = %v, want null", c.Children[0].StringValue)
	}
	if got := c.Children[0].Children[0].StringValue; got.IsNull() || got.ValueString() != "" {
		t.Errorf("EQUALS leaf string_value = %v, want \"\"", got)
	}
	if got := c.Children[1].StringValue.ValueString(); got != "Engineering" {
		t.Errorf("EQUALS string_value = %q", got)
	}
}

func TestFromAPIOwner(t *testing.T) {
	if m := fromAPI(t, `{"id":"d1","name":"n","owner":null}`, nil); m.Owner != nil {
		t.Errorf("owner = %v, want nil", m.Owner)
	}
	m := fromAPI(t, `{"id":"d1","name":"n","owner":{"type":"IDENTITY","id":"i1","name":null}}`, nil)
	if m.Owner == nil || m.Owner.ID.ValueString() != "i1" || m.Owner.Name.IsNull() || m.Owner.Name.ValueString() != "" {
		t.Errorf("owner = %+v, want id i1 and name \"\"", m.Owner)
	}
}

func TestFromAPIEmptyLists(t *testing.T) {
	body := `{"id":"d1","name":"n","accessProfiles":[],"entitlements":[],"membership":null}`
	for name, prior := range map[string]*roleDimensionModel{
		"nil prior":  nil,
		"null prior": {AccessProfiles: types.SetNull(common.ObjectRefObjectType), Entitlements: types.SetNull(common.ObjectRefObjectType)},
	} {
		m := fromAPI(t, body, prior)
		if !m.AccessProfiles.IsNull() || !m.Entitlements.IsNull() {
			t.Errorf("%s: sets = %v / %v, want null", name, m.AccessProfiles, m.Entitlements)
		}
		if m.Membership != nil {
			t.Errorf("%s: membership = %v, want nil", name, m.Membership)
		}
	}
	empty := types.SetValueMust(common.ObjectRefObjectType, nil)
	m := fromAPI(t, body, &roleDimensionModel{AccessProfiles: empty, Entitlements: empty})
	if m.AccessProfiles.IsNull() || len(m.AccessProfiles.Elements()) != 0 {
		t.Errorf("access_profiles = %v, want empty set", m.AccessProfiles)
	}
	if m.Entitlements.IsNull() || len(m.Entitlements.Elements()) != 0 {
		t.Errorf("entitlements = %v, want empty set", m.Entitlements)
	}
}

func TestToAPI(t *testing.T) {
	m := roleDimensionModel{
		ID:          types.StringValue("d1"),
		RoleID:      types.StringValue("r1"),
		Name:        types.StringValue("n"),
		Description: types.StringNull(),
		Owner:       &common.ObjectRefModel{Type: types.StringValue("IDENTITY"), ID: types.StringValue("i1"), Name: types.StringValue("Owner")},
		AccessProfiles: refSet(t,
			[3]string{"ACCESS_PROFILE", "ap2", "AP 2"},
			[3]string{"ACCESS_PROFILE", "ap1", ""},
		),
		Entitlements: types.SetNull(common.ObjectRefObjectType),
	}
	api, diags := m.ToAPI(context.Background())
	if diags.HasError() {
		t.Fatal(diags)
	}
	body, _ := json.Marshal(api)
	want := `{"name":"n","owner":{"type":"IDENTITY","id":"i1"},` +
		`"accessProfiles":[{"type":"ACCESS_PROFILE","id":"ap1"},{"type":"ACCESS_PROFILE","id":"ap2"}]}`
	if string(body) != want {
		t.Errorf("body = %s\nwant   %s", body, want)
	}
}

func baseModel(t *testing.T) roleDimensionModel {
	t.Helper()
	return roleDimensionModel{
		ID:          types.StringValue("d1"),
		RoleID:      types.StringValue("r1"),
		Name:        types.StringValue("n"),
		Description: types.StringValue("desc"),
		Owner:       &common.ObjectRefModel{Type: types.StringValue("IDENTITY"), ID: types.StringValue("i1"), Name: types.StringValue("")},
		AccessProfiles: refSet(t,
			[3]string{"ACCESS_PROFILE", "ap1", "AP 1"},
			[3]string{"ACCESS_PROFILE", "ap2", "AP 2"},
		),
		Entitlements: refSet(t, [3]string{"ENTITLEMENT", "e1", "E 1"}),
		Membership: &dimensionMembershipModel{
			Type: types.StringValue("STANDARD"),
			Criteria: &dimensionCriteriaModel{
				Operation:   types.StringValue("EQUALS"),
				Key:         &dimensionCriteriaKeyModel{Type: types.StringValue("IDENTITY"), Property: types.StringValue("attribute.department")},
				StringValue: types.StringValue("Engineering"),
			},
		},
		Created:  types.StringValue("2026-01-01T00:00:00Z"),
		Modified: types.StringValue("2026-01-01T00:00:00Z"),
	}
}

func TestToPatchOperations(t *testing.T) {
	cases := map[string]struct {
		change func(m *roleDimensionModel)
		want   string
	}{
		"reordered access profiles and unknown names": {
			change: func(m *roleDimensionModel) {
				m.AccessProfiles = refSet(t, [3]string{"ACCESS_PROFILE", "ap2", ""}, [3]string{"ACCESS_PROFILE", "ap1", ""})
				m.Owner.Name = types.StringUnknown()
				m.Modified = types.StringUnknown()
			},
			want: `null`,
		},
		"name only": {
			change: func(m *roleDimensionModel) { m.Name = types.StringValue("m") },
			want:   `[{"op":"replace","path":"/name","value":"m"}]`,
		},
		"description cleared": {
			change: func(m *roleDimensionModel) { m.Description = types.StringNull() },
			want:   `[{"op":"remove","path":"/description"}]`,
		},
		"owner changed": {
			change: func(m *roleDimensionModel) {
				m.Owner = &common.ObjectRefModel{Type: types.StringValue("IDENTITY"), ID: types.StringValue("i2"), Name: types.StringUnknown()}
			},
			want: `[{"op":"replace","path":"/owner","value":{"type":"IDENTITY","id":"i2"}}]`,
		},
		"owner removed": {
			change: func(m *roleDimensionModel) { m.Owner = nil },
			want:   `[{"op":"remove","path":"/owner"}]`,
		},
		"access profile added": {
			change: func(m *roleDimensionModel) {
				m.AccessProfiles = refSet(t,
					[3]string{"ACCESS_PROFILE", "ap3", ""},
					[3]string{"ACCESS_PROFILE", "ap2", "AP 2"},
					[3]string{"ACCESS_PROFILE", "ap1", "AP 1"},
				)
			},
			want: `[{"op":"replace","path":"/accessProfiles","value":[` +
				`{"type":"ACCESS_PROFILE","id":"ap1"},{"type":"ACCESS_PROFILE","id":"ap2"},{"type":"ACCESS_PROFILE","id":"ap3"}]}]`,
		},
		"access profiles removed": {
			change: func(m *roleDimensionModel) { m.AccessProfiles = types.SetNull(common.ObjectRefObjectType) },
			want:   `[{"op":"replace","path":"/accessProfiles","value":[]}]`,
		},
		"entitlements changed": {
			change: func(m *roleDimensionModel) { m.Entitlements = refSet(t, [3]string{"ENTITLEMENT", "e2", ""}) },
			want:   `[{"op":"replace","path":"/entitlements","value":[{"type":"ENTITLEMENT","id":"e2"}]}]`,
		},
		"membership changed": {
			change: func(m *roleDimensionModel) {
				m.Membership = &dimensionMembershipModel{
					Type: types.StringValue("STANDARD"),
					Criteria: &dimensionCriteriaModel{
						Operation:   types.StringValue("EQUALS"),
						Key:         &dimensionCriteriaKeyModel{Type: types.StringValue("IDENTITY"), Property: types.StringValue("attribute.department")},
						StringValue: types.StringValue("Sales"),
					},
				}
			},
			want: `[{"op":"replace","path":"/membership","value":{"type":"STANDARD","criteria":` +
				`{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"stringValue":"Sales"}}}]`,
		},
		"membership removed": {
			change: func(m *roleDimensionModel) { m.Membership = nil },
			want:   `[{"op":"remove","path":"/membership"}]`,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			state := baseModel(t)
			plan := baseModel(t)
			tc.change(&plan)
			ops, diags := plan.ToPatchOperations(context.Background(), &state)
			if diags.HasError() {
				t.Fatal(diags)
			}
			got, _ := json.Marshal(ops)
			if string(got) != tc.want {
				t.Errorf("ops = %s\nwant  %s", got, tc.want)
			}
		})
	}
}

func TestMembershipRoundTrip(t *testing.T) {
	api := decodeDimension(t, `{"id":"d1","name":"n","membership":{"type":"STANDARD","criteria":{
		"operation":"AND","children":[
			{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.department"},"stringValue":"Engineering"},
			{"operation":"OR","children":[
				{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"stringValue":"A"},
				{"operation":"EQUALS","key":{"type":"IDENTITY","property":"attribute.location"},"stringValue":"B"}]}]}}}`)
	got := membershipToAPI(membershipFromAPI(api.Membership))
	if !reflect.DeepEqual(got, api.Membership) {
		a, _ := json.Marshal(got)
		b, _ := json.Marshal(api.Membership)
		t.Errorf("round trip = %s\nwant         %s", a, b)
	}
}
