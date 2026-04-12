// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package access_profile

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------------------------------------------------------------------------
// Terraform models
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

type revokeRequestConfigModel struct {
	ApprovalSchemes []approvalSchemeModel `tfsdk:"approval_schemes"`
}

type additionalOwnerModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Three flattened levels for the provisioning criteria tree (max depth per SailPoint constraint).
type provisioningCriteriaLevel3Model struct {
	Operation types.String `tfsdk:"operation"`
	Attribute types.String `tfsdk:"attribute"`
	Value     types.String `tfsdk:"value"`
}

type provisioningCriteriaLevel2Model struct {
	Operation types.String                      `tfsdk:"operation"`
	Attribute types.String                      `tfsdk:"attribute"`
	Value     types.String                      `tfsdk:"value"`
	Children  []provisioningCriteriaLevel3Model `tfsdk:"children"`
}

type provisioningCriteriaModel struct {
	Operation types.String                      `tfsdk:"operation"`
	Attribute types.String                      `tfsdk:"attribute"`
	Value     types.String                      `tfsdk:"value"`
	Children  []provisioningCriteriaLevel2Model `tfsdk:"children"`
}

type accessProfileModel struct {
	ID                   types.String               `tfsdk:"id"`
	Name                 types.String               `tfsdk:"name"`
	Description          types.String               `tfsdk:"description"`
	Enabled              types.Bool                 `tfsdk:"enabled"`
	Requestable          types.Bool                 `tfsdk:"requestable"`
	Owner                *common.ObjectRefModel     `tfsdk:"owner"`
	Source               *common.ObjectRefModel     `tfsdk:"source"`
	Entitlements         []common.ObjectRefModel    `tfsdk:"entitlements"`
	Segments             types.Set                  `tfsdk:"segments"`
	AdditionalOwners     []additionalOwnerModel     `tfsdk:"additional_owners"`
	AccessRequestConfig  *accessRequestConfigModel  `tfsdk:"access_request_config"`
	RevokeRequestConfig  *revokeRequestConfigModel  `tfsdk:"revoke_request_config"`
	ProvisioningCriteria *provisioningCriteriaModel `tfsdk:"provisioning_criteria"`
	Created              types.String               `tfsdk:"created"`
	Modified             types.String               `tfsdk:"modified"`
}

// ---------------------------------------------------------------------------
// FromAPI
// ---------------------------------------------------------------------------

func (m *accessProfileModel) FromAPI(ctx context.Context, api *client.AccessProfileAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNull(api.Description)
	m.Enabled = boolPtrToTF(api.Enabled)
	m.Requestable = boolPtrToTF(api.Requestable)

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

	source, diags := common.NewObjectRefFromAPIPtr(ctx, api.Source)
	diagnostics.Append(diags...)
	m.Source = source

	if len(api.Entitlements) > 0 {
		m.Entitlements = make([]common.ObjectRefModel, 0, len(api.Entitlements))
		for i := range api.Entitlements {
			ref, d := common.NewObjectRefFromAPI(ctx, api.Entitlements[i])
			diagnostics.Append(d...)
			m.Entitlements = append(m.Entitlements, ref)
		}
	} else {
		m.Entitlements = nil
	}

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
	} else {
		m.AdditionalOwners = nil
	}

	m.AccessRequestConfig = accessRequestConfigFromAPI(api.AccessRequestConfig)
	m.RevokeRequestConfig = revokeRequestConfigFromAPI(api.RevokeRequestConfig)
	m.ProvisioningCriteria = provisioningCriteriaFromAPI(api.ProvisioningCriteria)

	return diagnostics
}

// ---------------------------------------------------------------------------
// ToAPI
// ---------------------------------------------------------------------------

func (m *accessProfileModel) ToAPI(ctx context.Context) (*client.AccessProfileAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.AccessProfileAPI{
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

	if m.Owner != nil {
		ownerAPI, diags := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(diags...)
		api.Owner = ownerAPI
	}
	if m.Source != nil {
		sourceAPI, diags := common.NewObjectRefToAPI(ctx, *m.Source)
		diagnostics.Append(diags...)
		api.Source = sourceAPI
	}

	if len(m.Entitlements) > 0 {
		api.Entitlements = make([]client.ObjectRefAPI, 0, len(m.Entitlements))
		for i := range m.Entitlements {
			ref, d := common.NewObjectRefToAPI(ctx, m.Entitlements[i])
			diagnostics.Append(d...)
			api.Entitlements = append(api.Entitlements, ref)
		}
	}

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

	api.AccessRequestConfig = accessRequestConfigToAPI(m.AccessRequestConfig)
	api.RevokeRequestConfig = revokeRequestConfigToAPI(m.RevokeRequestConfig)
	api.ProvisioningCriteria = provisioningCriteriaToAPI(m.ProvisioningCriteria)

	return api, diagnostics
}

// ---------------------------------------------------------------------------
// ToPatchOperations
// ---------------------------------------------------------------------------

func (m *accessProfileModel) ToPatchOperations(ctx context.Context, state *accessProfileModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
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

	if !reflect.DeepEqual(m.Owner, state.Owner) && m.Owner != nil {
		ownerAPI, d := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(d...)
		ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
	}

	// Source + entitlements can be patched together if source changed (API constraint);
	// otherwise each can be patched independently.
	sourceChanged := !reflect.DeepEqual(m.Source, state.Source)
	entsChanged := !reflect.DeepEqual(m.Entitlements, state.Entitlements)

	if sourceChanged && m.Source != nil {
		srcAPI, d := common.NewObjectRefToAPI(ctx, *m.Source)
		diagnostics.Append(d...)
		ops = append(ops, client.NewReplacePatch("/source", srcAPI))
	}
	if entsChanged || sourceChanged {
		ents := make([]client.ObjectRefAPI, 0, len(m.Entitlements))
		for i := range m.Entitlements {
			ref, d := common.NewObjectRefToAPI(ctx, m.Entitlements[i])
			diagnostics.Append(d...)
			ents = append(ents, ref)
		}
		ops = append(ops, client.NewReplacePatch("/entitlements", ents))
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

	if !reflect.DeepEqual(m.ProvisioningCriteria, state.ProvisioningCriteria) {
		if m.ProvisioningCriteria != nil {
			ops = append(ops, client.NewReplacePatch("/provisioningCriteria", provisioningCriteriaToAPI(m.ProvisioningCriteria)))
		} else {
			ops = append(ops, client.NewRemovePatch("/provisioningCriteria"))
		}
	}

	return ops, diagnostics
}

// ---------------------------------------------------------------------------
// Conversion helpers: access_request_config
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
		m.ApprovalSchemes = make([]approvalSchemeModel, 0, len(api.ApprovalSchemes))
		for _, s := range api.ApprovalSchemes {
			m.ApprovalSchemes = append(m.ApprovalSchemes, approvalSchemeModel{
				ApproverType: types.StringValue(s.ApproverType),
				ApproverID:   stringPtrToTF(s.ApproverID),
			})
		}
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
		api.ApprovalSchemes = make([]client.ApprovalSchemeAPI, 0, len(m.ApprovalSchemes))
		for _, s := range m.ApprovalSchemes {
			api.ApprovalSchemes = append(api.ApprovalSchemes, client.ApprovalSchemeAPI{
				ApproverType: s.ApproverType.ValueString(),
				ApproverID:   tfToStringPtr(s.ApproverID),
			})
		}
	}
	return api
}

func revokeRequestConfigFromAPI(api *client.RevocabilityAPI) *revokeRequestConfigModel {
	if api == nil {
		return nil
	}
	m := &revokeRequestConfigModel{}
	if len(api.ApprovalSchemes) > 0 {
		m.ApprovalSchemes = make([]approvalSchemeModel, 0, len(api.ApprovalSchemes))
		for _, s := range api.ApprovalSchemes {
			m.ApprovalSchemes = append(m.ApprovalSchemes, approvalSchemeModel{
				ApproverType: types.StringValue(s.ApproverType),
				ApproverID:   stringPtrToTF(s.ApproverID),
			})
		}
	}
	return m
}

func revokeRequestConfigToAPI(m *revokeRequestConfigModel) *client.RevocabilityAPI {
	if m == nil {
		return nil
	}
	api := &client.RevocabilityAPI{}
	if len(m.ApprovalSchemes) > 0 {
		api.ApprovalSchemes = make([]client.ApprovalSchemeAPI, 0, len(m.ApprovalSchemes))
		for _, s := range m.ApprovalSchemes {
			api.ApprovalSchemes = append(api.ApprovalSchemes, client.ApprovalSchemeAPI{
				ApproverType: s.ApproverType.ValueString(),
				ApproverID:   tfToStringPtr(s.ApproverID),
			})
		}
	}
	return api
}

// ---------------------------------------------------------------------------
// Conversion helpers: provisioning_criteria (3-level tree)
// ---------------------------------------------------------------------------

func provisioningCriteriaFromAPI(api *client.ProvisioningCriteriaAPI) *provisioningCriteriaModel {
	if api == nil {
		return nil
	}
	m := &provisioningCriteriaModel{
		Operation: types.StringValue(api.Operation),
		Attribute: stringPtrToTF(api.Attribute),
		Value:     stringPtrToTF(api.Value),
	}
	if len(api.Children) > 0 {
		m.Children = make([]provisioningCriteriaLevel2Model, 0, len(api.Children))
		for i := range api.Children {
			child := &api.Children[i]
			lvl2 := provisioningCriteriaLevel2Model{
				Operation: types.StringValue(child.Operation),
				Attribute: stringPtrToTF(child.Attribute),
				Value:     stringPtrToTF(child.Value),
			}
			if len(child.Children) > 0 {
				lvl2.Children = make([]provisioningCriteriaLevel3Model, 0, len(child.Children))
				for j := range child.Children {
					leaf := &child.Children[j]
					lvl2.Children = append(lvl2.Children, provisioningCriteriaLevel3Model{
						Operation: types.StringValue(leaf.Operation),
						Attribute: stringPtrToTF(leaf.Attribute),
						Value:     stringPtrToTF(leaf.Value),
					})
				}
			}
			m.Children = append(m.Children, lvl2)
		}
	}
	return m
}

func provisioningCriteriaToAPI(m *provisioningCriteriaModel) *client.ProvisioningCriteriaAPI {
	if m == nil {
		return nil
	}
	api := &client.ProvisioningCriteriaAPI{
		Operation: m.Operation.ValueString(),
		Attribute: tfToStringPtr(m.Attribute),
		Value:     tfToStringPtr(m.Value),
	}
	if len(m.Children) > 0 {
		api.Children = make([]client.ProvisioningCriteriaAPI, 0, len(m.Children))
		for i := range m.Children {
			child := &m.Children[i]
			lvl2 := client.ProvisioningCriteriaAPI{
				Operation: child.Operation.ValueString(),
				Attribute: tfToStringPtr(child.Attribute),
				Value:     tfToStringPtr(child.Value),
			}
			if len(child.Children) > 0 {
				lvl2.Children = make([]client.ProvisioningCriteriaAPI, 0, len(child.Children))
				for j := range child.Children {
					leaf := &child.Children[j]
					lvl2.Children = append(lvl2.Children, client.ProvisioningCriteriaAPI{
						Operation: leaf.Operation.ValueString(),
						Attribute: tfToStringPtr(leaf.Attribute),
						Value:     tfToStringPtr(leaf.Value),
					})
				}
			}
			api.Children = append(api.Children, lvl2)
		}
	}
	return api
}

// ---------------------------------------------------------------------------
// Primitive conversion helpers
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
