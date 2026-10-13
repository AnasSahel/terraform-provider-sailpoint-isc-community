// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package approval_config

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const fullAPIJSON = `{
  "id": "r1",
  "scope": "DOMAIN_OBJECT",
  "reminderConfig": {"enabled": true, "daysUntilFirstReminder": 2, "reminderCronSchedule": "@every 48h", "maxReminders": 3},
  "escalationConfig": {"enabled": true, "daysUntilFirstEscalation": 7, "escalationCronSchedule": "@every 72h",
    "escalationChain": [
      {"tier": 1, "identityId": "", "identityType": "MANAGER_OF"},
      {"tier": 2, "identityId": "g1", "identityType": "GOVERNANCE_GROUP"}
    ]},
  "timeoutConfig": {"enabled": true, "daysUntilTimeout": 90, "timeoutResult": "EXPIRED"},
  "cronTimezone": {"location": "Europe/Paris", "offset": ""},
  "serialChain": [{"identityId": "", "identityType": "MANAGER_OF"}],
  "requiresComment": "OFF",
  "fallbackApprover": {"identityID": "abc", "type": "GOVERNANCE_GROUP"},
  "machineIdentityManagerAssignment": "MANAGER_OF_REQUESTER",
  "circumventApprovalProcess": false,
  "autoApprove": "OFF"
}`

func decodeAPI(t *testing.T, raw string) *client.ApprovalConfigAPI {
	t.Helper()
	var api client.ApprovalConfigAPI
	if err := json.Unmarshal([]byte(raw), &api); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return &api
}

func newModel() approvalConfigModel {
	return approvalConfigModel{
		Scope:    types.StringValue("ROLE"),
		ObjectID: types.StringValue("r1"),
	}
}

func TestFromAPIRoundTrip(t *testing.T) {
	ctx := context.Background()
	m := newModel()
	if diags := m.FromAPI(ctx, decodeAPI(t, fullAPIJSON)); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}
	body, diags := m.ToAPI(ctx)
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags)
	}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	var got, want map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(fullAPIJSON), &want); err != nil {
		t.Fatal(err)
	}
	delete(want, "id")
	delete(want, "scope")
	esc, _ := want["escalationConfig"].(map[string]any)
	chain, _ := esc["escalationChain"].([]any)
	for _, it := range chain {
		item, _ := it.(map[string]any)
		delete(item, "tier")
		if item["identityId"] == "" {
			delete(item, "identityId")
		}
	}
	serial, _ := want["serialChain"].([]any)
	for _, it := range serial {
		item, _ := it.(map[string]any)
		if item["identityId"] == "" {
			delete(item, "identityId")
		}
	}
	tz, _ := want["cronTimezone"].(map[string]any)
	if tz["offset"] == "" {
		delete(tz, "offset")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip mismatch\n got: %v\nwant: %v", got, want)
	}
}

func TestFromAPIEmptyStringsBecomeNull(t *testing.T) {
	ctx := context.Background()
	m := newModel()
	if diags := m.FromAPI(ctx, decodeAPI(t, fullAPIJSON)); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}

	var ec escalationConfigModel
	if diags := m.EscalationConfig.As(ctx, &ec, objectAsOptions); diags.HasError() {
		t.Fatal(diags)
	}
	var chain []chainItemModel
	if diags := ec.EscalationChain.ElementsAs(ctx, &chain, false); diags.HasError() {
		t.Fatal(diags)
	}
	if !chain[0].IdentityID.IsNull() || chain[1].IdentityID.ValueString() != "g1" {
		t.Errorf("chain identity ids = %v, %v", chain[0].IdentityID, chain[1].IdentityID)
	}

	api := decodeAPI(t, `{"fallbackApprover":{"identityID":"","type":"MANAGER_OF"},"cronTimezone":{"location":"UTC","offset":""}}`)
	m2 := newModel()
	if diags := m2.FromAPI(ctx, api); diags.HasError() {
		t.Fatal(diags)
	}
	var fb fallbackApproverModel
	if diags := m2.FallbackApprover.As(ctx, &fb, objectAsOptions); diags.HasError() {
		t.Fatal(diags)
	}
	if !fb.IdentityID.IsNull() {
		t.Errorf("fallback identity_id = %v, want null", fb.IdentityID)
	}
	var tz cronTimezoneModel
	if diags := m2.CronTimezone.As(ctx, &tz, objectAsOptions); diags.HasError() {
		t.Fatal(diags)
	}
	if !tz.Offset.IsNull() {
		t.Errorf("offset = %v, want null", tz.Offset)
	}
}

func TestFromAPIChainsNilAndEmptyAreEmptyList(t *testing.T) {
	ctx := context.Background()
	m := newModel()
	api := decodeAPI(t, `{"escalationConfig":{"enabled":false},"serialChain":[]}`)
	if diags := m.FromAPI(ctx, api); diags.HasError() {
		t.Fatal(diags)
	}
	var ec escalationConfigModel
	if diags := m.EscalationConfig.As(ctx, &ec, objectAsOptions); diags.HasError() {
		t.Fatal(diags)
	}
	for name, l := range map[string]types.List{"escalation_chain": ec.EscalationChain, "serial_chain": m.SerialChain} {
		if l.IsNull() || l.IsUnknown() || len(l.Elements()) != 0 {
			t.Errorf("%s = %v, want known empty list", name, l)
		}
	}
}

func TestFromAPIKeepsScopeAndObjectID(t *testing.T) {
	m := newModel()
	if diags := m.FromAPI(context.Background(), decodeAPI(t, `{"id":"r1","scope":"DOMAIN_OBJECT"}`)); diags.HasError() {
		t.Fatal(diags)
	}
	if m.Scope.ValueString() != "ROLE" || m.ObjectID.ValueString() != "r1" || m.ID.ValueString() != "ROLE/r1" {
		t.Errorf("scope/object_id/id = %s/%s/%s", m.Scope, m.ObjectID, m.ID)
	}
}

func unknownModel() approvalConfigModel {
	m := newModel()
	m.ReminderConfig = types.ObjectUnknown(reminderConfigAttrTypes)
	m.EscalationConfig = types.ObjectNull(escalationConfigAttrTypes)
	m.TimeoutConfig = types.ObjectUnknown(timeoutConfigAttrTypes)
	m.CronTimezone = types.ObjectNull(cronTimezoneAttrTypes)
	m.FallbackApprover = types.ObjectUnknown(fallbackApproverAttrTypes)
	m.SerialChain = types.ListUnknown(chainItemObjectType)
	m.RequiresComment = types.StringUnknown()
	m.MachineIdentityManagerAssignment = types.StringNull()
	m.CircumventApprovalProcess = types.BoolUnknown()
	m.AutoApprove = types.StringNull()
	return m
}

func TestToAPIOmitsUnknownAndNull(t *testing.T) {
	m := unknownModel()
	body, diags := m.ToAPI(context.Background())
	if diags.HasError() {
		t.Fatal(diags)
	}
	raw, _ := json.Marshal(body)
	if string(raw) != "{}" {
		t.Errorf("body = %s, want {}", raw)
	}
}

func TestToAPIUnknownInnerValueOmitted(t *testing.T) {
	ctx := context.Background()
	m := unknownModel()
	rc, diags := types.ObjectValue(reminderConfigAttrTypes, map[string]attr.Value{
		"enabled":                   types.BoolValue(true),
		"days_until_first_reminder": types.Int64Unknown(),
		"reminder_cron_schedule":    types.StringUnknown(),
		"max_reminders":             types.Int64Unknown(),
	})
	if diags.HasError() {
		t.Fatal(diags)
	}
	m.ReminderConfig = rc
	body, diags := m.ToAPI(ctx)
	if diags.HasError() {
		t.Fatal(diags)
	}
	raw, _ := json.Marshal(body)
	if string(raw) != `{"reminderConfig":{"enabled":true}}` {
		t.Errorf("body = %s", raw)
	}
}

func TestIsOwnConfig(t *testing.T) {
	own, inherited := "r1", "tenant-1"
	cases := []struct {
		name string
		api  *client.ApprovalConfigAPI
		want bool
	}{
		{"own", &client.ApprovalConfigAPI{ID: &own}, true},
		{"inherited", &client.ApprovalConfigAPI{ID: &inherited}, false},
		{"nil id", &client.ApprovalConfigAPI{}, false},
		{"nil api", nil, false},
	}
	for _, tc := range cases {
		if got := isOwnConfig(tc.api, "r1"); got != tc.want {
			t.Errorf("%s: isOwnConfig = %v, want %v", tc.name, got, tc.want)
		}
	}
}

var objectAsOptions = basetypes.ObjectAsOptions{}
