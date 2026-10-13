// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package application

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationAccountSourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	Name                     types.String `tfsdk:"name"`
	UseForPasswordManagement types.Bool   `tfsdk:"use_for_password_management"`
	PasswordPolicies         types.List   `tfsdk:"password_policies"`
}

type applicationModel struct {
	ID                      types.String                   `tfsdk:"id"`
	CloudAppID              types.String                   `tfsdk:"cloud_app_id"`
	Name                    types.String                   `tfsdk:"name"`
	Description             types.String                   `tfsdk:"description"`
	Enabled                 types.Bool                     `tfsdk:"enabled"`
	ProvisionRequestEnabled types.Bool                     `tfsdk:"provision_request_enabled"`
	MatchAllAccounts        types.Bool                     `tfsdk:"match_all_accounts"`
	AppCenterEnabled        types.Bool                     `tfsdk:"app_center_enabled"`
	AccountSource           *applicationAccountSourceModel `tfsdk:"account_source"`
	Owner                   *common.ObjectRefModel         `tfsdk:"owner"`
	Created                 types.String                   `tfsdk:"created"`
	Modified                types.String                   `tfsdk:"modified"`
}

func (m *applicationModel) FromAPI(ctx context.Context, api *client.ApplicationAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	if api == nil {
		return diagnostics
	}

	m.ID = types.StringValue(api.ID)
	m.CloudAppID = common.StringOrNullIfEmpty(api.CloudAppID)
	m.Name = common.StringOrNullIfEmpty(api.Name)
	m.Description = common.StringOrNull(api.Description)

	if api.Enabled != nil {
		m.Enabled = types.BoolValue(*api.Enabled)
	} else {
		m.Enabled = types.BoolNull()
	}
	if api.ProvisionRequestEnabled != nil {
		m.ProvisionRequestEnabled = types.BoolValue(*api.ProvisionRequestEnabled)
	} else {
		m.ProvisionRequestEnabled = types.BoolNull()
	}
	if api.MatchAllAccounts != nil {
		m.MatchAllAccounts = types.BoolValue(*api.MatchAllAccounts)
	} else {
		m.MatchAllAccounts = types.BoolNull()
	}
	if api.AppCenterEnabled != nil {
		m.AppCenterEnabled = types.BoolValue(*api.AppCenterEnabled)
	} else {
		m.AppCenterEnabled = types.BoolNull()
	}

	m.Created = common.StringOrNull(api.Created)
	m.Modified = common.StringOrNull(api.Modified)

	if api.AccountSource != nil {
		accountSource, diags := accountSourceFromAPI(ctx, api.AccountSource)
		diagnostics.Append(diags...)
		m.AccountSource = accountSource
	} else {
		m.AccountSource = nil
	}

	if api.Owner != nil {
		owner, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
		m.Owner = owner
	} else {
		m.Owner = nil
	}

	return diagnostics
}

func (m *applicationModel) ToCreateAPI(ctx context.Context) (*client.ApplicationCreateAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.ApplicationCreateAPI{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		api.Description = m.Description.ValueString()
	}

	if !m.MatchAllAccounts.IsNull() && !m.MatchAllAccounts.IsUnknown() {
		matchAll := m.MatchAllAccounts.ValueBool()
		api.MatchAllAccounts = &matchAll
	}

	if m.AccountSource == nil || m.AccountSource.ID.IsNull() || m.AccountSource.ID.IsUnknown() {
		diagnostics.AddError("Missing account_source", "account_source.id is required when creating an application.")
		return nil, diagnostics
	}

	accountSource := client.ApplicationCreateAccountSourceAPI{
		ID: m.AccountSource.ID.ValueString(),
	}
	if !m.AccountSource.Type.IsNull() && !m.AccountSource.Type.IsUnknown() {
		sourceType := m.AccountSource.Type.ValueString()
		accountSource.Type = &sourceType
	}
	api.AccountSource = accountSource

	return api, diagnostics
}

func (m *applicationModel) ToPatchOperations(ctx context.Context, state *applicationModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var ops []client.JSONPatchOperation

	if !m.Name.Equal(state.Name) {
		ops = append(ops, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	// Unknown means "not configured, keep whatever the server has": never patch it.
	if !m.Description.IsUnknown() && !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
			ops = append(ops, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			ops = append(ops, client.NewRemovePatch("/description"))
		}
	}

	if isKnown(m.Enabled) && !m.Enabled.Equal(state.Enabled) {
		ops = append(ops, client.NewReplacePatch("/enabled", m.Enabled.ValueBool()))
	}

	if isKnown(m.ProvisionRequestEnabled) && !m.ProvisionRequestEnabled.Equal(state.ProvisionRequestEnabled) {
		ops = append(ops, client.NewReplacePatch("/provisionRequestEnabled", m.ProvisionRequestEnabled.ValueBool()))
	}

	if isKnown(m.MatchAllAccounts) && !m.MatchAllAccounts.Equal(state.MatchAllAccounts) {
		ops = append(ops, client.NewReplacePatch("/matchAllAccounts", m.MatchAllAccounts.ValueBool()))
	}

	if isKnown(m.AppCenterEnabled) && !m.AppCenterEnabled.Equal(state.AppCenterEnabled) {
		ops = append(ops, client.NewReplacePatch("/appCenterEnabled", m.AppCenterEnabled.ValueBool()))
	}

	if ownerChanged(m.Owner, state.Owner) {
		if m.Owner != nil {
			ownerAPI, diags := common.NewObjectRefToAPIPtr(ctx, *m.Owner)
			diagnostics.Append(diags...)
			ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
		} else {
			ops = append(ops, client.NewRemovePatch("/owner"))
		}
	}

	return ops, diagnostics
}

// ownerChanged compares owners by type and id only: name is resolved by the server.
func ownerChanged(plan, state *common.ObjectRefModel) bool {
	if plan == nil || state == nil {
		return (plan == nil) != (state == nil)
	}
	return !plan.Type.Equal(state.Type) || !plan.ID.Equal(state.ID)
}

// isKnown reports whether an Optional+Computed value was set in the configuration.
func isKnown(v types.Bool) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func accountSourceFromAPI(ctx context.Context, api *client.ApplicationAccountSourceAPI) (*applicationAccountSourceModel, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if api == nil {
		return nil, diagnostics
	}

	model := &applicationAccountSourceModel{
		ID:   common.StringOrNullIfEmpty(api.ID),
		Type: common.StringOrNullIfEmpty(api.Type),
		Name: common.StringOrNullIfEmpty(api.Name),
	}

	if api.UseForPasswordManagement != nil {
		model.UseForPasswordManagement = types.BoolValue(*api.UseForPasswordManagement)
	} else {
		model.UseForPasswordManagement = types.BoolNull()
	}

	if len(api.PasswordPolicies) > 0 {
		policies, diags := common.MapListFromAPI(ctx, api.PasswordPolicies, common.ObjectRefObjectType, common.NewObjectRefFromAPI)
		diagnostics.Append(diags...)
		model.PasswordPolicies = policies
	} else {
		model.PasswordPolicies = types.ListNull(common.ObjectRefObjectType)
	}

	return model, diagnostics
}
