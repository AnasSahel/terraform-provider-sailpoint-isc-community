// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"
	"testing"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func opPaths(ops []client.JSONPatchOperation) map[string]any {
	paths := make(map[string]any, len(ops))
	for _, op := range ops {
		paths[op.Path] = op.Value
	}
	return paths
}

// After create, attributes the user did not configure are unknown in the plan.
// They must not be patched (an unknown bool would otherwise be sent as false).
func TestToPatchOperationsSkipsUnknownAfterCreate(t *testing.T) {
	ctx := context.Background()
	enabled := false
	matchAll := false
	desc := "d"
	created := &client.ApplicationAPI{
		ID:               "app1",
		Name:             "App",
		Description:      &desc,
		Enabled:          &enabled,
		MatchAllAccounts: &matchAll,
	}
	var state applicationModel
	if diags := state.FromAPI(ctx, created); diags.HasError() {
		t.Fatalf("FromAPI: %v", diags)
	}

	plan := applicationModel{
		Name:                    types.StringValue("App"),
		Description:             types.StringValue("d"),
		Enabled:                 types.BoolValue(true),
		ProvisionRequestEnabled: types.BoolUnknown(),
		MatchAllAccounts:        types.BoolValue(false),
		AppCenterEnabled:        types.BoolUnknown(),
		Owner: &common.ObjectRefModel{
			Type: types.StringValue("IDENTITY"),
			ID:   types.StringValue("owner1"),
			Name: types.StringUnknown(),
		},
	}

	ops, diags := plan.ToPatchOperations(ctx, &state)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	got := opPaths(ops)
	if len(got) != 2 {
		t.Fatalf("ops = %+v, want only /enabled and /owner", ops)
	}
	if v, ok := got["/enabled"].(bool); !ok || !v {
		t.Errorf("/enabled = %v, want true", got["/enabled"])
	}
	owner, ok := got["/owner"].(*client.ObjectRefAPI)
	if !ok || owner.ID != "owner1" || owner.Name != "" {
		t.Errorf("/owner = %#v, want id owner1 without name", got["/owner"])
	}
}

func TestToPatchOperationsIgnoresOwnerName(t *testing.T) {
	ctx := context.Background()
	state := applicationModel{Owner: &common.ObjectRefModel{
		Type: types.StringValue("IDENTITY"), ID: types.StringValue("o1"), Name: types.StringValue("Jane"),
	}}
	plan := applicationModel{Owner: &common.ObjectRefModel{
		Type: types.StringValue("IDENTITY"), ID: types.StringValue("o1"), Name: types.StringUnknown(),
	}}

	ops, diags := plan.ToPatchOperations(ctx, &state)
	if diags.HasError() {
		t.Fatalf("ToPatchOperations: %v", diags)
	}
	if len(ops) != 0 {
		t.Fatalf("ops = %+v, want none", ops)
	}

	plan.Owner = nil
	ops, _ = plan.ToPatchOperations(ctx, &state)
	if len(ops) != 1 || ops[0].Op != "remove" || ops[0].Path != "/owner" {
		t.Fatalf("ops = %+v, want remove /owner", ops)
	}
}

func TestApplicationNameFilterEscapesQuotes(t *testing.T) {
	if got := applicationNameFilter(`My "App"`); got != `name eq "My \"App\""` {
		t.Errorf("filter = %q", got)
	}
}

func TestApplicationAccountSourceFromAPINullPasswordPolicies(t *testing.T) {
	t.Parallel()

	useForPasswordManagement := true
	model, diags := accountSourceFromAPI(context.Background(), &client.ApplicationAccountSourceAPI{
		ID:                       "source-id",
		Type:                     "SOURCE",
		Name:                     "Active Directory",
		UseForPasswordManagement: &useForPasswordManagement,
	})

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if model == nil {
		t.Fatal("expected account source model")
	}
	if !model.PasswordPolicies.IsNull() {
		t.Fatalf("expected null password_policies, got %#v", model.PasswordPolicies)
	}
}

func TestApplicationToPatchOperationsAfterCreate(t *testing.T) {
	t.Parallel()

	created := applicationModel{
		ID:                      types.StringValue("app-id"),
		Name:                    types.StringValue("Test"),
		Description:             types.StringValue("Test application"),
		Enabled:                 types.BoolValue(false),
		ProvisionRequestEnabled: types.BoolValue(false),
		MatchAllAccounts:        types.BoolValue(true),
		AppCenterEnabled:        types.BoolValue(false),
		Owner: &common.ObjectRefModel{
			Type: types.StringValue("IDENTITY"),
			ID:   types.StringValue("ac26a05177f84b249fc8e46ee4148f3e"),
			Name: types.StringValue("andre.faria"),
		},
	}

	plan := applicationModel{
		Name:                    types.StringValue("Test"),
		Description:             types.StringValue("Test application"),
		Enabled:                 types.BoolValue(false),
		ProvisionRequestEnabled: types.BoolValue(true),
		MatchAllAccounts:        types.BoolValue(true),
		AppCenterEnabled:        types.BoolValue(true),
		Owner: &common.ObjectRefModel{
			Type: types.StringValue("IDENTITY"),
			ID:   types.StringValue("2448605646d0486fb6eb3b4c8bada9a8"),
		},
	}

	ops, diags := plan.ToPatchOperations(context.Background(), &created)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 patch operations, got %d: %#v", len(ops), ops)
	}
}

func TestApplicationModelFromAPIWithPasswordPolicies(t *testing.T) {
	t.Parallel()

	var model applicationModel
	diags := model.FromAPI(context.Background(), &client.ApplicationAPI{
		ID:   "app-id",
		Name: "Accounting Functions",
		AccountSource: &client.ApplicationAccountSourceAPI{
			ID:   "source-id",
			Type: "SOURCE",
			Name: "Active Directory",
			PasswordPolicies: []client.ObjectRefAPI{
				{Type: "PASSWORD_POLICY", ID: "policy-id", Name: "Default Policy"},
			},
		},
	})

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if model.AccountSource == nil {
		t.Fatal("expected account source")
	}
	if model.AccountSource.PasswordPolicies.IsNull() || model.AccountSource.PasswordPolicies.IsUnknown() {
		t.Fatal("expected known password_policies list")
	}

	var refs []common.ObjectRefModel
	if err := model.AccountSource.PasswordPolicies.ElementsAs(context.Background(), &refs, false); err != nil {
		t.Fatalf("elements as: %v", err)
	}
	if len(refs) != 1 || refs[0].ID.ValueString() != "policy-id" {
		t.Fatalf("unexpected password policies: %#v", refs)
	}
}
