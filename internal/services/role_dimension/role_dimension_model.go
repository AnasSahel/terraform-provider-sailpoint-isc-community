// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role_dimension

import (
	"context"
	"reflect"
	"sort"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type dimensionCriteriaKeyModel struct {
	Type     types.String `tfsdk:"type"`
	Property types.String `tfsdk:"property"`
}

// Three-level flattened criteria tree, as on the role.
type dimensionCriteriaLevel3Model struct {
	Operation   types.String               `tfsdk:"operation"`
	Key         *dimensionCriteriaKeyModel `tfsdk:"key"`
	StringValue types.String               `tfsdk:"string_value"`
}

type dimensionCriteriaLevel2Model struct {
	Operation   types.String                   `tfsdk:"operation"`
	Key         *dimensionCriteriaKeyModel     `tfsdk:"key"`
	StringValue types.String                   `tfsdk:"string_value"`
	Children    []dimensionCriteriaLevel3Model `tfsdk:"children"`
}

type dimensionCriteriaModel struct {
	Operation   types.String                   `tfsdk:"operation"`
	Key         *dimensionCriteriaKeyModel     `tfsdk:"key"`
	StringValue types.String                   `tfsdk:"string_value"`
	Children    []dimensionCriteriaLevel2Model `tfsdk:"children"`
}

type dimensionMembershipModel struct {
	Type     types.String            `tfsdk:"type"`
	Criteria *dimensionCriteriaModel `tfsdk:"criteria"`
}

type roleDimensionModel struct {
	ID             types.String              `tfsdk:"id"`
	RoleID         types.String              `tfsdk:"role_id"`
	Name           types.String              `tfsdk:"name"`
	Description    types.String              `tfsdk:"description"`
	Owner          *common.ObjectRefModel    `tfsdk:"owner"`
	AccessProfiles types.Set                 `tfsdk:"access_profiles"`
	Entitlements   types.Set                 `tfsdk:"entitlements"`
	Membership     *dimensionMembershipModel `tfsdk:"membership"`
	Created        types.String              `tfsdk:"created"`
	Modified       types.String              `tfsdk:"modified"`
}

// ---------------------------------------------------------------------------
// FromAPI
// ---------------------------------------------------------------------------

// FromAPI fills the model from an API response. prior is the plan (Create,
// Update) or the state (Read) and decides whether an empty list from the API
// reads as an empty set or as null; it may be nil.
func (m *roleDimensionModel) FromAPI(ctx context.Context, roleID string, api *client.RoleDimensionAPI, prior *roleDimensionModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	priorAccessProfiles := types.SetNull(common.ObjectRefObjectType)
	priorEntitlements := types.SetNull(common.ObjectRefObjectType)
	if prior != nil {
		priorAccessProfiles = prior.AccessProfiles
		priorEntitlements = prior.Entitlements
	}

	m.ID = types.StringValue(api.ID)
	m.RoleID = types.StringValue(roleID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNull(api.Description)
	m.Created = common.StringOrNull(api.Created)
	m.Modified = common.StringOrNull(api.Modified)

	if api.Owner != nil {
		owner, d := common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(d...)
		m.Owner = owner
	} else {
		m.Owner = nil
	}

	var d diag.Diagnostics
	m.AccessProfiles, d = refSetFromAPI(ctx, api.AccessProfiles, priorAccessProfiles)
	diagnostics.Append(d...)
	m.Entitlements, d = refSetFromAPI(ctx, api.Entitlements, priorEntitlements)
	diagnostics.Append(d...)

	m.Membership = membershipFromAPI(api.Membership)

	return diagnostics
}

// ---------------------------------------------------------------------------
// ToAPI
// ---------------------------------------------------------------------------

// ToAPI builds the create request body. Only type and id of refs are sent.
func (m *roleDimensionModel) ToAPI(ctx context.Context) (*client.RoleDimensionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.RoleDimensionAPI{
		Name:        m.Name.ValueString(),
		Description: tfToStringPtr(m.Description),
	}

	if m.Owner != nil {
		owner, d := common.NewObjectRefToAPIPtr(ctx, *m.Owner)
		diagnostics.Append(d...)
		api.Owner = owner
	}

	var d diag.Diagnostics
	api.AccessProfiles, d = refSetToAPI(ctx, m.AccessProfiles)
	diagnostics.Append(d...)
	api.Entitlements, d = refSetToAPI(ctx, m.Entitlements)
	diagnostics.Append(d...)

	api.Membership = membershipToAPI(m.Membership)

	return api, diagnostics
}

// ---------------------------------------------------------------------------
// ToPatchOperations
// ---------------------------------------------------------------------------

// ToPatchOperations returns a JSON Patch with only the fields that differ
// between the plan (m) and the state. Refs are compared by type and id, so a
// different order or an unknown planned name adds no operation.
func (m *roleDimensionModel) ToPatchOperations(ctx context.Context, state *roleDimensionModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var ops []client.JSONPatchOperation

	if !m.Name.Equal(state.Name) {
		ops = append(ops, client.NewReplacePatch("/name", m.Name.ValueString()))
	}
	if !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
			ops = append(ops, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			ops = append(ops, client.NewRemovePatch("/description"))
		}
	}

	switch {
	case m.Owner != nil && (state.Owner == nil ||
		!m.Owner.Type.Equal(state.Owner.Type) || !m.Owner.ID.Equal(state.Owner.ID)):
		owner, d := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(d...)
		ops = append(ops, client.NewReplacePatch("/owner", owner))
	case m.Owner == nil && state.Owner != nil:
		ops = append(ops, client.NewRemovePatch("/owner"))
	}

	for _, f := range []struct {
		path        string
		plan, state types.Set
	}{
		{"/accessProfiles", m.AccessProfiles, state.AccessProfiles},
		{"/entitlements", m.Entitlements, state.Entitlements},
	} {
		planKeys, d := refKeys(ctx, f.plan)
		diagnostics.Append(d...)
		stateKeys, d := refKeys(ctx, f.state)
		diagnostics.Append(d...)
		if reflect.DeepEqual(planKeys, stateKeys) {
			continue
		}
		refs, d := refSetToAPI(ctx, f.plan)
		diagnostics.Append(d...)
		if refs == nil {
			refs = []client.ObjectRefAPI{}
		}
		ops = append(ops, client.NewReplacePatch(f.path, refs))
	}

	if !reflect.DeepEqual(m.Membership, state.Membership) {
		if m.Membership != nil {
			ops = append(ops, client.NewReplacePatch("/membership", membershipToAPI(m.Membership)))
		} else {
			ops = append(ops, client.NewRemovePatch("/membership"))
		}
	}

	return ops, diagnostics
}

// ---------------------------------------------------------------------------
// Conversion helpers: ref sets
// ---------------------------------------------------------------------------

// refSetFromAPI maps API refs to a set. An empty list reads as null, unless
// prior is a known empty set, so `access_profiles = []` stays consistent.
func refSetFromAPI(ctx context.Context, refs []client.ObjectRefAPI, prior types.Set) (types.Set, diag.Diagnostics) {
	if len(refs) == 0 {
		if !prior.IsNull() && !prior.IsUnknown() && len(prior.Elements()) == 0 {
			return types.SetValueMust(common.ObjectRefObjectType, nil), nil
		}
		return types.SetNull(common.ObjectRefObjectType), nil
	}
	var diagnostics diag.Diagnostics
	models := make([]common.ObjectRefModel, 0, len(refs))
	for i := range refs {
		ref, d := common.NewObjectRefFromAPI(ctx, refs[i])
		diagnostics.Append(d...)
		models = append(models, ref)
	}
	set, d := types.SetValueFrom(ctx, common.ObjectRefObjectType, models)
	diagnostics.Append(d...)
	return set, diagnostics
}

// refSetToAPI returns the refs in a set as {type, id}, sorted by id. A null,
// unknown or empty set gives nil.
func refSetToAPI(ctx context.Context, set types.Set) ([]client.ObjectRefAPI, diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() || len(set.Elements()) == 0 {
		return nil, nil
	}
	var diagnostics diag.Diagnostics
	var models []common.ObjectRefModel
	diagnostics.Append(set.ElementsAs(ctx, &models, false)...)
	out := make([]client.ObjectRefAPI, 0, len(models))
	for i := range models {
		ref, d := common.NewObjectRefToAPI(ctx, models[i])
		diagnostics.Append(d...)
		out = append(out, ref)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ID != out[j].ID {
			return out[i].ID < out[j].ID
		}
		return out[i].Type < out[j].Type
	})
	return out, diagnostics
}

// refKeys returns the sorted "type|id" keys of the refs in a set. A null,
// unknown or empty set gives nil.
func refKeys(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	refs, diagnostics := refSetToAPI(ctx, set)
	if len(refs) == 0 {
		return nil, diagnostics
	}
	keys := make([]string, 0, len(refs))
	for _, r := range refs {
		keys = append(keys, r.Type+"|"+r.ID)
	}
	sort.Strings(keys)
	return keys, diagnostics
}

// ---------------------------------------------------------------------------
// Conversion helpers: membership + criteria tree
// ---------------------------------------------------------------------------

func membershipFromAPI(api *client.RoleMembershipAPI) *dimensionMembershipModel {
	if api == nil {
		return nil
	}
	return &dimensionMembershipModel{
		Type:     types.StringValue(api.Type),
		Criteria: criteriaFromAPI(api.Criteria),
	}
}

func membershipToAPI(m *dimensionMembershipModel) *client.RoleMembershipAPI {
	if m == nil {
		return nil
	}
	return &client.RoleMembershipAPI{
		Type:     m.Type.ValueString(),
		Criteria: criteriaToAPI(m.Criteria),
	}
}

func criteriaFromAPI(api *client.RoleCriteriaAPI) *dimensionCriteriaModel {
	if api == nil {
		return nil
	}
	m := &dimensionCriteriaModel{
		Operation:   types.StringValue(api.Operation),
		Key:         criteriaKeyFromAPI(api.Key),
		StringValue: criteriaStringValue(api.Operation, api.StringValue),
	}
	if len(api.Children) > 0 {
		m.Children = make([]dimensionCriteriaLevel2Model, 0, len(api.Children))
		for i := range api.Children {
			child := &api.Children[i]
			lvl2 := dimensionCriteriaLevel2Model{
				Operation:   types.StringValue(child.Operation),
				Key:         criteriaKeyFromAPI(child.Key),
				StringValue: criteriaStringValue(child.Operation, child.StringValue),
			}
			if len(child.Children) > 0 {
				lvl2.Children = make([]dimensionCriteriaLevel3Model, 0, len(child.Children))
				for j := range child.Children {
					leaf := &child.Children[j]
					lvl2.Children = append(lvl2.Children, dimensionCriteriaLevel3Model{
						Operation:   types.StringValue(leaf.Operation),
						Key:         criteriaKeyFromAPI(leaf.Key),
						StringValue: criteriaStringValue(leaf.Operation, leaf.StringValue),
					})
				}
			}
			m.Children = append(m.Children, lvl2)
		}
	}
	return m
}

func criteriaToAPI(m *dimensionCriteriaModel) *client.RoleCriteriaAPI {
	if m == nil {
		return nil
	}
	api := &client.RoleCriteriaAPI{
		Operation:   m.Operation.ValueString(),
		Key:         criteriaKeyToAPI(m.Key),
		StringValue: tfToStringPtr(m.StringValue),
	}
	if len(m.Children) > 0 {
		api.Children = make([]client.RoleCriteriaAPI, 0, len(m.Children))
		for i := range m.Children {
			child := &m.Children[i]
			lvl2 := client.RoleCriteriaAPI{
				Operation:   child.Operation.ValueString(),
				Key:         criteriaKeyToAPI(child.Key),
				StringValue: tfToStringPtr(child.StringValue),
			}
			if len(child.Children) > 0 {
				lvl2.Children = make([]client.RoleCriteriaAPI, 0, len(child.Children))
				for j := range child.Children {
					leaf := &child.Children[j]
					lvl2.Children = append(lvl2.Children, client.RoleCriteriaAPI{
						Operation:   leaf.Operation.ValueString(),
						Key:         criteriaKeyToAPI(leaf.Key),
						StringValue: tfToStringPtr(leaf.StringValue),
					})
				}
			}
			api.Children = append(api.Children, lvl2)
		}
	}
	return api
}

// criteriaStringValue maps a node's stringValue. The UI writes "" on AND and
// OR nodes, which carry no value, so that reads as null.
func criteriaStringValue(operation string, value *string) types.String {
	if (operation == "AND" || operation == "OR") && value != nil && *value == "" {
		return types.StringNull()
	}
	return common.StringOrNull(value)
}

func criteriaKeyFromAPI(api *client.RoleCriteriaKeyAPI) *dimensionCriteriaKeyModel {
	if api == nil {
		return nil
	}
	return &dimensionCriteriaKeyModel{
		Type:     types.StringValue(api.Type),
		Property: types.StringValue(api.Property),
	}
}

func criteriaKeyToAPI(m *dimensionCriteriaKeyModel) *client.RoleCriteriaKeyAPI {
	if m == nil {
		return nil
	}
	return &client.RoleCriteriaKeyAPI{
		Type:     m.Type.ValueString(),
		Property: m.Property.ValueString(),
	}
}

// ---------------------------------------------------------------------------
// Primitive helpers
// ---------------------------------------------------------------------------

func tfToStringPtr(s types.String) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}
