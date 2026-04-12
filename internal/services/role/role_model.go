// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package role

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type approvalSchemeModel struct {
	ApproverType types.String `tfsdk:"approver_type"`
	ApproverID   types.String `tfsdk:"approver_id"`
}

type accessDurationModel struct {
	Value    types.Int64  `tfsdk:"value"`
	TimeUnit types.String `tfsdk:"time_unit"`
}

type accessRequestConfigModel struct {
	CommentsRequired           types.Bool            `tfsdk:"comments_required"`
	DenialCommentsRequired     types.Bool            `tfsdk:"denial_comments_required"`
	ReauthorizationRequired    types.Bool            `tfsdk:"reauthorization_required"`
	RequireEndDate             types.Bool            `tfsdk:"require_end_date"`
	MaxPermittedAccessDuration *accessDurationModel  `tfsdk:"max_permitted_access_duration"`
	ApprovalSchemes            []approvalSchemeModel `tfsdk:"approval_schemes"`
}

// Role revoke config has extra comment fields vs access profile revoke config.
type revokeRequestConfigModel struct {
	CommentsRequired       types.Bool            `tfsdk:"comments_required"`
	DenialCommentsRequired types.Bool            `tfsdk:"denial_comments_required"`
	ApprovalSchemes        []approvalSchemeModel `tfsdk:"approval_schemes"`
}

type additionalOwnerModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type membershipIdentityModel struct {
	ID        types.String `tfsdk:"id"`
	Type      types.String `tfsdk:"type"`
	Name      types.String `tfsdk:"name"`
	AliasName types.String `tfsdk:"alias_name"`
}

type criteriaKeyModel struct {
	Type     types.String `tfsdk:"type"`
	Property types.String `tfsdk:"property"`
	SourceID types.String `tfsdk:"source_id"`
}

// Three-level flattened criteria tree.
type criteriaLevel3Model struct {
	Operation   types.String      `tfsdk:"operation"`
	Key         *criteriaKeyModel `tfsdk:"key"`
	StringValue types.String      `tfsdk:"string_value"`
}

type criteriaLevel2Model struct {
	Operation   types.String          `tfsdk:"operation"`
	Key         *criteriaKeyModel     `tfsdk:"key"`
	StringValue types.String          `tfsdk:"string_value"`
	Children    []criteriaLevel3Model `tfsdk:"children"`
}

type roleCriteriaModel struct {
	Operation   types.String          `tfsdk:"operation"`
	Key         *criteriaKeyModel     `tfsdk:"key"`
	StringValue types.String          `tfsdk:"string_value"`
	Children    []criteriaLevel2Model `tfsdk:"children"`
}

type membershipModel struct {
	Type       types.String              `tfsdk:"type"`
	Criteria   *roleCriteriaModel        `tfsdk:"criteria"`
	Identities []membershipIdentityModel `tfsdk:"identities"`
}

type roleModel struct {
	ID                  types.String              `tfsdk:"id"`
	Name                types.String              `tfsdk:"name"`
	Description         types.String              `tfsdk:"description"`
	Enabled             types.Bool                `tfsdk:"enabled"`
	Requestable         types.Bool                `tfsdk:"requestable"`
	Dimensional         types.Bool                `tfsdk:"dimensional"`
	Owner               *common.ObjectRefModel    `tfsdk:"owner"`
	AccessProfiles      []common.ObjectRefModel   `tfsdk:"access_profiles"`
	Entitlements        []common.ObjectRefModel   `tfsdk:"entitlements"`
	Segments            types.Set                 `tfsdk:"segments"`
	AdditionalOwners    []additionalOwnerModel    `tfsdk:"additional_owners"`
	Membership          *membershipModel          `tfsdk:"membership"`
	AccessRequestConfig *accessRequestConfigModel `tfsdk:"access_request_config"`
	RevokeRequestConfig *revokeRequestConfigModel `tfsdk:"revoke_request_config"`
	DimensionRefs       []common.ObjectRefModel   `tfsdk:"dimension_refs"`
	Created             types.String              `tfsdk:"created"`
	Modified            types.String              `tfsdk:"modified"`
}

// ---------------------------------------------------------------------------
// FromAPI
// ---------------------------------------------------------------------------

func (m *roleModel) FromAPI(ctx context.Context, api *client.RoleAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNull(api.Description)
	m.Enabled = boolPtrToTF(api.Enabled)
	m.Requestable = boolPtrToTF(api.Requestable)
	m.Dimensional = boolPtrToTF(api.Dimensional)

	if api.Created != nil {
		m.Created = types.StringValue(*api.Created)
	} else {
		m.Created = types.StringNull()
	}
	if api.Modified != nil {
		m.Modified = types.StringValue(*api.Modified)
	} else {
		m.Modified = types.StringNull()
	}

	owner, diags := common.NewObjectRefFromAPIPtr(ctx, api.Owner)
	diagnostics.Append(diags...)
	m.Owner = owner

	m.AccessProfiles = objectRefSliceFromAPI(ctx, api.AccessProfiles, &diagnostics)
	m.Entitlements = objectRefSliceFromAPI(ctx, api.Entitlements, &diagnostics)
	m.DimensionRefs = objectRefSliceFromAPI(ctx, api.DimensionRefs, &diagnostics)

	if api.Segments != nil {
		segs, d := types.SetValueFrom(ctx, types.StringType, api.Segments)
		diagnostics.Append(d...)
		m.Segments = segs
	} else {
		m.Segments = types.SetNull(types.StringType)
	}

	if len(api.AdditionalOwners) > 0 {
		m.AdditionalOwners = make([]additionalOwnerModel, 0, len(api.AdditionalOwners))
		for i := range api.AdditionalOwners {
			m.AdditionalOwners = append(m.AdditionalOwners, additionalOwnerModel{
				Type: types.StringValue(api.AdditionalOwners[i].Type),
				ID:   types.StringValue(api.AdditionalOwners[i].ID),
				Name: types.StringValue(api.AdditionalOwners[i].Name),
			})
		}
	}

	m.Membership = membershipFromAPI(api.Membership)
	m.AccessRequestConfig = accessRequestConfigFromAPI(api.AccessRequestConfig)
	m.RevokeRequestConfig = revokeRequestConfigFromAPI(api.RevokeRequestConfig)

	return diagnostics
}

// ---------------------------------------------------------------------------
// ToAPI
// ---------------------------------------------------------------------------

func (m *roleModel) ToAPI(ctx context.Context) (*client.RoleAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.RoleAPI{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		v := m.Description.ValueString()
		api.Description = &v
	}
	if !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		v := m.Enabled.ValueBool()
		api.Enabled = &v
	}
	if !m.Requestable.IsNull() && !m.Requestable.IsUnknown() {
		v := m.Requestable.ValueBool()
		api.Requestable = &v
	}
	if !m.Dimensional.IsNull() && !m.Dimensional.IsUnknown() {
		v := m.Dimensional.ValueBool()
		api.Dimensional = &v
	}

	if m.Owner != nil {
		ownerAPI, diags := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(diags...)
		api.Owner = ownerAPI
	}

	api.AccessProfiles = objectRefSliceToAPI(ctx, m.AccessProfiles, &diagnostics)
	api.Entitlements = objectRefSliceToAPI(ctx, m.Entitlements, &diagnostics)
	api.DimensionRefs = objectRefSliceToAPI(ctx, m.DimensionRefs, &diagnostics)

	if !m.Segments.IsNull() && !m.Segments.IsUnknown() {
		var segs []string
		diagnostics.Append(m.Segments.ElementsAs(ctx, &segs, false)...)
		api.Segments = segs
	}

	if len(m.AdditionalOwners) > 0 {
		api.AdditionalOwners = make([]client.ObjectRefAPI, 0, len(m.AdditionalOwners))
		for i := range m.AdditionalOwners {
			api.AdditionalOwners = append(api.AdditionalOwners, client.ObjectRefAPI{
				Type: m.AdditionalOwners[i].Type.ValueString(),
				ID:   m.AdditionalOwners[i].ID.ValueString(),
			})
		}
	}

	api.Membership = membershipToAPI(m.Membership)
	api.AccessRequestConfig = accessRequestConfigToAPI(m.AccessRequestConfig)
	api.RevokeRequestConfig = revokeRequestConfigToAPI(m.RevokeRequestConfig)

	return api, diagnostics
}

// ---------------------------------------------------------------------------
// ToPatchOperations
// ---------------------------------------------------------------------------

func (m *roleModel) ToPatchOperations(ctx context.Context, state *roleModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
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
	if !m.Enabled.Equal(state.Enabled) && !m.Enabled.IsNull() && !m.Enabled.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/enabled", m.Enabled.ValueBool()))
	}
	if !m.Requestable.Equal(state.Requestable) && !m.Requestable.IsNull() && !m.Requestable.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/requestable", m.Requestable.ValueBool()))
	}
	if !m.Dimensional.Equal(state.Dimensional) && !m.Dimensional.IsNull() && !m.Dimensional.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/dimensional", m.Dimensional.ValueBool()))
	}

	if !reflect.DeepEqual(m.Owner, state.Owner) && m.Owner != nil {
		ownerAPI, d := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(d...)
		ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
	}

	if !reflect.DeepEqual(m.AccessProfiles, state.AccessProfiles) {
		refs := objectRefSliceToAPI(ctx, m.AccessProfiles, &diagnostics)
		if refs == nil {
			refs = []client.ObjectRefAPI{}
		}
		ops = append(ops, client.NewReplacePatch("/accessProfiles", refs))
	}
	if !reflect.DeepEqual(m.Entitlements, state.Entitlements) {
		refs := objectRefSliceToAPI(ctx, m.Entitlements, &diagnostics)
		if refs == nil {
			refs = []client.ObjectRefAPI{}
		}
		ops = append(ops, client.NewReplacePatch("/entitlements", refs))
	}
	if !reflect.DeepEqual(m.DimensionRefs, state.DimensionRefs) {
		refs := objectRefSliceToAPI(ctx, m.DimensionRefs, &diagnostics)
		if refs == nil {
			refs = []client.ObjectRefAPI{}
		}
		ops = append(ops, client.NewReplacePatch("/dimensionRefs", refs))
	}

	if !m.Segments.Equal(state.Segments) {
		if !m.Segments.IsNull() && !m.Segments.IsUnknown() {
			var segs []string
			diagnostics.Append(m.Segments.ElementsAs(ctx, &segs, false)...)
			ops = append(ops, client.NewReplacePatch("/segments", segs))
		} else {
			ops = append(ops, client.NewRemovePatch("/segments"))
		}
	}

	if !reflect.DeepEqual(m.AdditionalOwners, state.AdditionalOwners) {
		if m.AdditionalOwners != nil {
			refs := make([]client.ObjectRefAPI, 0, len(m.AdditionalOwners))
			for i := range m.AdditionalOwners {
				refs = append(refs, client.ObjectRefAPI{
					Type: m.AdditionalOwners[i].Type.ValueString(),
					ID:   m.AdditionalOwners[i].ID.ValueString(),
				})
			}
			ops = append(ops, client.NewReplacePatch("/additionalOwners", refs))
		} else {
			ops = append(ops, client.NewRemovePatch("/additionalOwners"))
		}
	}

	if !reflect.DeepEqual(m.Membership, state.Membership) {
		if m.Membership != nil {
			ops = append(ops, client.NewReplacePatch("/membership", membershipToAPI(m.Membership)))
		} else {
			ops = append(ops, client.NewRemovePatch("/membership"))
		}
	}

	if !reflect.DeepEqual(m.AccessRequestConfig, state.AccessRequestConfig) {
		if m.AccessRequestConfig != nil {
			ops = append(ops, client.NewReplacePatch("/accessRequestConfig", accessRequestConfigToAPI(m.AccessRequestConfig)))
		} else {
			ops = append(ops, client.NewRemovePatch("/accessRequestConfig"))
		}
	}

	if !reflect.DeepEqual(m.RevokeRequestConfig, state.RevokeRequestConfig) {
		if m.RevokeRequestConfig != nil {
			ops = append(ops, client.NewReplacePatch("/revokeRequestConfig", revokeRequestConfigToAPI(m.RevokeRequestConfig)))
		} else {
			ops = append(ops, client.NewRemovePatch("/revokeRequestConfig"))
		}
	}

	return ops, diagnostics
}

// ---------------------------------------------------------------------------
// Conversion helpers: ObjectRef slices
// ---------------------------------------------------------------------------

func objectRefSliceFromAPI(ctx context.Context, in []client.ObjectRefAPI, diags *diag.Diagnostics) []common.ObjectRefModel {
	if len(in) == 0 {
		return nil
	}
	out := make([]common.ObjectRefModel, 0, len(in))
	for i := range in {
		ref, d := common.NewObjectRefFromAPI(ctx, in[i])
		diags.Append(d...)
		out = append(out, ref)
	}
	return out
}

func objectRefSliceToAPI(ctx context.Context, in []common.ObjectRefModel, diags *diag.Diagnostics) []client.ObjectRefAPI {
	if len(in) == 0 {
		return nil
	}
	out := make([]client.ObjectRefAPI, 0, len(in))
	for i := range in {
		ref, d := common.NewObjectRefToAPI(ctx, in[i])
		diags.Append(d...)
		out = append(out, ref)
	}
	return out
}

// ---------------------------------------------------------------------------
// Conversion helpers: access/revoke request configs
// ---------------------------------------------------------------------------

func accessRequestConfigFromAPI(api *client.RequestabilityAPI) *accessRequestConfigModel {
	if api == nil {
		return nil
	}
	m := &accessRequestConfigModel{
		CommentsRequired:        boolPtrToTF(api.CommentsRequired),
		DenialCommentsRequired:  boolPtrToTF(api.DenialCommentsRequired),
		ReauthorizationRequired: boolPtrToTF(api.ReauthorizationRequired),
		RequireEndDate:          boolPtrToTF(api.RequireEndDate),
	}
	if api.MaxPermittedAccessDuration != nil {
		m.MaxPermittedAccessDuration = &accessDurationModel{
			Value:    int64PtrToTF(api.MaxPermittedAccessDuration.Value),
			TimeUnit: stringPtrToTF(api.MaxPermittedAccessDuration.TimeUnit),
		}
	}
	if len(api.ApprovalSchemes) > 0 {
		m.ApprovalSchemes = approvalSchemesFromAPI(api.ApprovalSchemes)
	}
	return m
}

func accessRequestConfigToAPI(m *accessRequestConfigModel) *client.RequestabilityAPI {
	if m == nil {
		return nil
	}
	api := &client.RequestabilityAPI{}
	if !m.CommentsRequired.IsNull() && !m.CommentsRequired.IsUnknown() {
		v := m.CommentsRequired.ValueBool()
		api.CommentsRequired = &v
	}
	if !m.DenialCommentsRequired.IsNull() && !m.DenialCommentsRequired.IsUnknown() {
		v := m.DenialCommentsRequired.ValueBool()
		api.DenialCommentsRequired = &v
	}
	if !m.ReauthorizationRequired.IsNull() && !m.ReauthorizationRequired.IsUnknown() {
		v := m.ReauthorizationRequired.ValueBool()
		api.ReauthorizationRequired = &v
	}
	if !m.RequireEndDate.IsNull() && !m.RequireEndDate.IsUnknown() {
		v := m.RequireEndDate.ValueBool()
		api.RequireEndDate = &v
	}
	if m.MaxPermittedAccessDuration != nil {
		api.MaxPermittedAccessDuration = &client.AccessDurationAPI{
			Value:    tfToInt64Ptr(m.MaxPermittedAccessDuration.Value),
			TimeUnit: tfToStringPtr(m.MaxPermittedAccessDuration.TimeUnit),
		}
	}
	if len(m.ApprovalSchemes) > 0 {
		api.ApprovalSchemes = approvalSchemesToAPI(m.ApprovalSchemes)
	}
	return api
}

func revokeRequestConfigFromAPI(api *client.RevocabilityForRoleAPI) *revokeRequestConfigModel {
	if api == nil {
		return nil
	}
	m := &revokeRequestConfigModel{
		CommentsRequired:       boolPtrToTF(api.CommentsRequired),
		DenialCommentsRequired: boolPtrToTF(api.DenialCommentsRequired),
	}
	if len(api.ApprovalSchemes) > 0 {
		m.ApprovalSchemes = approvalSchemesFromAPI(api.ApprovalSchemes)
	}
	return m
}

func revokeRequestConfigToAPI(m *revokeRequestConfigModel) *client.RevocabilityForRoleAPI {
	if m == nil {
		return nil
	}
	api := &client.RevocabilityForRoleAPI{}
	if !m.CommentsRequired.IsNull() && !m.CommentsRequired.IsUnknown() {
		v := m.CommentsRequired.ValueBool()
		api.CommentsRequired = &v
	}
	if !m.DenialCommentsRequired.IsNull() && !m.DenialCommentsRequired.IsUnknown() {
		v := m.DenialCommentsRequired.ValueBool()
		api.DenialCommentsRequired = &v
	}
	if len(m.ApprovalSchemes) > 0 {
		api.ApprovalSchemes = approvalSchemesToAPI(m.ApprovalSchemes)
	}
	return api
}

func approvalSchemesFromAPI(in []client.ApprovalSchemeAPI) []approvalSchemeModel {
	out := make([]approvalSchemeModel, 0, len(in))
	for _, s := range in {
		out = append(out, approvalSchemeModel{
			ApproverType: types.StringValue(s.ApproverType),
			ApproverID:   stringPtrToTF(s.ApproverID),
		})
	}
	return out
}

func approvalSchemesToAPI(in []approvalSchemeModel) []client.ApprovalSchemeAPI {
	out := make([]client.ApprovalSchemeAPI, 0, len(in))
	for _, s := range in {
		out = append(out, client.ApprovalSchemeAPI{
			ApproverType: s.ApproverType.ValueString(),
			ApproverID:   tfToStringPtr(s.ApproverID),
		})
	}
	return out
}

// ---------------------------------------------------------------------------
// Conversion helpers: membership + criteria tree
// ---------------------------------------------------------------------------

func membershipFromAPI(api *client.RoleMembershipAPI) *membershipModel {
	if api == nil {
		return nil
	}
	m := &membershipModel{
		Type: types.StringValue(api.Type),
	}
	if api.Criteria != nil {
		m.Criteria = criteriaFromAPI(api.Criteria)
	}
	if len(api.Identities) > 0 {
		m.Identities = make([]membershipIdentityModel, 0, len(api.Identities))
		for _, id := range api.Identities {
			m.Identities = append(m.Identities, membershipIdentityModel{
				ID:        types.StringValue(id.ID),
				Type:      types.StringValue(id.Type),
				Name:      types.StringValue(id.Name),
				AliasName: types.StringValue(id.AliasName),
			})
		}
	}
	return m
}

func membershipToAPI(m *membershipModel) *client.RoleMembershipAPI {
	if m == nil {
		return nil
	}
	api := &client.RoleMembershipAPI{
		Type: m.Type.ValueString(),
	}
	if m.Criteria != nil {
		api.Criteria = criteriaToAPI(m.Criteria)
	}
	if len(m.Identities) > 0 {
		api.Identities = make([]client.RoleMembershipIdentityAPI, 0, len(m.Identities))
		for _, id := range m.Identities {
			api.Identities = append(api.Identities, client.RoleMembershipIdentityAPI{
				ID:   id.ID.ValueString(),
				Type: id.Type.ValueString(),
			})
		}
	}
	return api
}

func criteriaFromAPI(api *client.RoleCriteriaAPI) *roleCriteriaModel {
	if api == nil {
		return nil
	}
	m := &roleCriteriaModel{
		Operation:   types.StringValue(api.Operation),
		Key:         criteriaKeyFromAPI(api.Key),
		StringValue: stringPtrToTF(api.StringValue),
	}
	if len(api.Children) > 0 {
		m.Children = make([]criteriaLevel2Model, 0, len(api.Children))
		for i := range api.Children {
			child := &api.Children[i]
			lvl2 := criteriaLevel2Model{
				Operation:   types.StringValue(child.Operation),
				Key:         criteriaKeyFromAPI(child.Key),
				StringValue: stringPtrToTF(child.StringValue),
			}
			if len(child.Children) > 0 {
				lvl2.Children = make([]criteriaLevel3Model, 0, len(child.Children))
				for j := range child.Children {
					leaf := &child.Children[j]
					lvl2.Children = append(lvl2.Children, criteriaLevel3Model{
						Operation:   types.StringValue(leaf.Operation),
						Key:         criteriaKeyFromAPI(leaf.Key),
						StringValue: stringPtrToTF(leaf.StringValue),
					})
				}
			}
			m.Children = append(m.Children, lvl2)
		}
	}
	return m
}

func criteriaToAPI(m *roleCriteriaModel) *client.RoleCriteriaAPI {
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

func criteriaKeyFromAPI(api *client.RoleCriteriaKeyAPI) *criteriaKeyModel {
	if api == nil {
		return nil
	}
	return &criteriaKeyModel{
		Type:     types.StringValue(api.Type),
		Property: types.StringValue(api.Property),
		SourceID: stringPtrToTF(api.SourceID),
	}
}

func criteriaKeyToAPI(m *criteriaKeyModel) *client.RoleCriteriaKeyAPI {
	if m == nil {
		return nil
	}
	return &client.RoleCriteriaKeyAPI{
		Type:     m.Type.ValueString(),
		Property: m.Property.ValueString(),
		SourceID: tfToStringPtr(m.SourceID),
	}
}

// ---------------------------------------------------------------------------
// Primitive helpers
// ---------------------------------------------------------------------------

func boolPtrToTF(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

func int64PtrToTF(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

func stringPtrToTF(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

func tfToStringPtr(s types.String) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}

func tfToInt64Ptr(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	val := v.ValueInt64()
	return &val
}
