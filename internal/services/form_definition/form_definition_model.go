// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package form_definition

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Element type definitions for types.List conversions.
var (
	formInputObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"type":        types.StringType,
		"label":       types.StringType,
		"description": types.StringType,
	}}

	formConditionRuleObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"source_type": types.StringType,
		"source":      types.StringType,
		"operator":    types.StringType,
		"value_type":  types.StringType,
		"value":       types.StringType,
	}}

	formConditionEffectConfigObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"default_value_label": types.StringType,
		"element":             types.StringType,
	}}

	formConditionEffectObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"effect_type": types.StringType,
		"config":      formConditionEffectConfigObjectType,
	}}

	formConditionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"rule_operator": types.StringType,
		"rules":         types.ListType{ElemType: formConditionRuleObjectType},
		"effects":       types.ListType{ElemType: formConditionEffectObjectType},
	}}
)

// formInputModel represents a form input field in Terraform state.
type formInputModel struct {
	ID          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
}

func NewFormInputFromAPI(ctx context.Context, api client.FormInputAPI) (formInputModel, diag.Diagnostics) {
	var m formInputModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func NewFormInputToAPI(ctx context.Context, m formInputModel) (client.FormInputAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}

func (m *formInputModel) FromAPI(ctx context.Context, api client.FormInputAPI) diag.Diagnostics {
	m.ID = types.StringValue(api.ID)
	m.Type = types.StringValue(api.Type)
	m.Label = types.StringValue(api.Label)
	m.Description = types.StringValue(api.Description)

	return nil
}

func (m *formInputModel) ToAPI(ctx context.Context) (client.FormInputAPI, diag.Diagnostics) {
	api := client.FormInputAPI{
		ID:          m.ID.ValueString(),
		Type:        m.Type.ValueString(),
		Label:       m.Label.ValueString(),
		Description: m.Description.ValueString(),
	}

	return api, nil
}

// formConditionModel represents a form condition in Terraform state.
type formConditionModel struct {
	RuleOperator types.String               `tfsdk:"rule_operator"`
	Rules        []formConditionRuleModel   `tfsdk:"rules"`
	Effects      []formConditionEffectModel `tfsdk:"effects"`
}

func NewFormConditionFromAPI(ctx context.Context, api client.FormConditionAPI) (formConditionModel, diag.Diagnostics) {
	var m formConditionModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func (m *formConditionModel) FromAPI(ctx context.Context, api client.FormConditionAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.RuleOperator = types.StringValue(api.RuleOperator)

	m.Rules, diags = common.MapSliceFromAPI(ctx, api.Rules, NewFormConditionRuleFromAPI)
	diagnostics.Append(diags...)

	m.Effects, diags = common.MapSliceFromAPI(ctx, api.Effects, NewFormConditionEffectFromAPI)
	diagnostics.Append(diags...)

	return diagnostics
}

func (m *formConditionModel) ToAPI(ctx context.Context) (client.FormConditionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	rules := make([]client.FormConditionRuleAPI, len(m.Rules))
	for i := range m.Rules {
		rules[i], diags = m.Rules[i].ToAPI(ctx)
		diagnostics.Append(diags...)
	}

	effects := make([]client.FormConditionEffectAPI, len(m.Effects))
	for i := range m.Effects {
		effects[i], diags = m.Effects[i].ToAPI(ctx)
		diagnostics.Append(diags...)
	}

	api := client.FormConditionAPI{
		RuleOperator: m.RuleOperator.ValueString(),
		Rules:        rules,
		Effects:      effects,
	}

	return api, diagnostics
}

func FormConditionToAPI(ctx context.Context, m formConditionModel) (client.FormConditionAPI, diag.Diagnostics) {
	return m.ToAPI(ctx)
}

// formConditionRuleModel represents a rule within a form condition.
type formConditionRuleModel struct {
	SourceType types.String `tfsdk:"source_type"`
	Source     types.String `tfsdk:"source"`
	Operator   types.String `tfsdk:"operator"`
	ValueType  types.String `tfsdk:"value_type"`
	Value      types.String `tfsdk:"value"`
}

func NewFormConditionRuleFromAPI(ctx context.Context, api client.FormConditionRuleAPI) (formConditionRuleModel, diag.Diagnostics) {
	var m formConditionRuleModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func (m *formConditionRuleModel) FromAPI(ctx context.Context, api client.FormConditionRuleAPI) diag.Diagnostics {
	m.SourceType = types.StringValue(api.SourceType)
	m.Source = types.StringValue(api.Source)
	m.Operator = types.StringValue(api.Operator)
	m.ValueType = types.StringValue(api.ValueType)
	m.Value = types.StringValue(api.Value)

	return nil
}

func (m *formConditionRuleModel) ToAPI(ctx context.Context) (client.FormConditionRuleAPI, diag.Diagnostics) {
	api := client.FormConditionRuleAPI{
		SourceType: m.SourceType.ValueString(),
		Source:     m.Source.ValueString(),
		Operator:   m.Operator.ValueString(),
		ValueType:  m.ValueType.ValueString(),
		Value:      m.Value.ValueString(),
	}

	return api, nil
}

// formConditionEffectModel represents the effect of a form condition in Terraform state.
type formConditionEffectModel struct {
	EffectType types.String                   `tfsdk:"effect_type"`
	Config     formConditionEffectConfigModel `tfsdk:"config"`
}

func NewFormConditionEffectFromAPI(ctx context.Context, api client.FormConditionEffectAPI) (formConditionEffectModel, diag.Diagnostics) {
	var m formConditionEffectModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func (m *formConditionEffectModel) FromAPI(ctx context.Context, api client.FormConditionEffectAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	m.EffectType = types.StringValue(api.EffectType)
	m.Config, diags = NewFormConditionEffectConfigFromAPI(ctx, api.Config)

	return diags
}

func (m *formConditionEffectModel) ToAPI(ctx context.Context) (client.FormConditionEffectAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	configAPI, diags := m.Config.ToAPI(ctx)
	diagnostics.Append(diags...)

	api := client.FormConditionEffectAPI{
		EffectType: m.EffectType.ValueString(),
		Config:     configAPI,
	}

	return api, diagnostics
}

// formConditionEffectConfigModel represents the configuration for a form condition effect in Terraform state.
type formConditionEffectConfigModel struct {
	DefaultValueLabel types.String `tfsdk:"default_value_label"`
	Element           types.String `tfsdk:"element"`
}

func NewFormConditionEffectConfigFromAPI(ctx context.Context, api client.FormConditionEffectConfigAPI) (formConditionEffectConfigModel, diag.Diagnostics) {
	var m formConditionEffectConfigModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func (m *formConditionEffectConfigModel) FromAPI(ctx context.Context, api client.FormConditionEffectConfigAPI) diag.Diagnostics {
	m.DefaultValueLabel = types.StringValue(api.DefaultValueLabel)
	m.Element = types.StringValue(api.Element)

	return nil
}

func (m *formConditionEffectConfigModel) ToAPI(ctx context.Context) (client.FormConditionEffectConfigAPI, diag.Diagnostics) {
	api := client.FormConditionEffectConfigAPI{
		DefaultValueLabel: m.DefaultValueLabel.ValueString(),
		Element:           m.Element.ValueString(),
	}

	return api, nil
}

// formDefinitionModel represents the Terraform state for a SailPoint form definition.
type formDefinitionModel struct {
	ID             types.String           `tfsdk:"id"`
	Name           types.String           `tfsdk:"name"`
	Description    types.String           `tfsdk:"description"`
	Owner          *common.ObjectRefModel `tfsdk:"owner"`
	UsedBy         types.List             `tfsdk:"used_by"`
	FormInput      types.List             `tfsdk:"form_input"`
	FormElements   jsontypes.Normalized   `tfsdk:"form_elements"`
	FormConditions types.List             `tfsdk:"form_conditions"`
	Created        types.String           `tfsdk:"created"`
	Modified       types.String           `tfsdk:"modified"`
}

// FromAPI maps fields from the API response to the Terraform model.
func (m *formDefinitionModel) FromAPI(ctx context.Context, api client.FormDefinitionAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNullIfEmpty(api.Description)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map owner
	m.Owner, diags = common.NewObjectRefFromAPIPtr(ctx, api.Owner)
	diagnostics.Append(diags...)

	// Map usedBy (Optional only — normalize empty to null)
	if len(api.UsedBy) > 0 {
		m.UsedBy, diags = common.MapListFromAPI(ctx, api.UsedBy, common.ObjectRefObjectType, common.NewObjectRefFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.UsedBy = types.ListNull(common.ObjectRefObjectType)
	}

	// Map formInput (Optional only — normalize empty to null)
	if len(api.FormInput) > 0 {
		m.FormInput, diags = common.MapListFromAPI(ctx, api.FormInput, formInputObjectType, NewFormInputFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.FormInput = types.ListNull(formInputObjectType)
	}

	// Map formConditions (Optional only — normalize empty to null)
	if len(api.FormConditions) > 0 {
		m.FormConditions, diags = common.MapListFromAPI(ctx, api.FormConditions, formConditionObjectType, NewFormConditionFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.FormConditions = types.ListNull(formConditionObjectType)
	}

	// Map formElements (Optional only — normalize empty to null)
	if len(api.FormElements) > 0 {
		m.FormElements, diags = common.MarshalJSONOrDefault(api.FormElements, "[]")
		diagnostics.Append(diags...)
	} else {
		m.FormElements = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API create request.
func (m *formDefinitionModel) ToAPI(ctx context.Context) (client.FormDefinitionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics
	var apiRequest client.FormDefinitionAPI

	apiRequest.Name = m.Name.ValueString()
	apiRequest.Description = m.Description.ValueString()

	// Parse owner
	apiRequest.Owner, diags = m.Owner.ToAPI(ctx)
	diagnostics.Append(diags...)

	// Parse formInput from types.List
	apiRequest.FormInput, diags = common.MapListToAPI(ctx, m.FormInput, NewFormInputToAPI)
	diagnostics.Append(diags...)

	// Parse formConditions from types.List
	apiRequest.FormConditions, diags = common.MapListToAPI(ctx, m.FormConditions, FormConditionToAPI)
	diagnostics.Append(diags...)

	// Parse formElements from JSON
	if elements, diags := common.UnmarshalJSONField[[]client.FormElementAPI](m.FormElements); elements != nil {
		apiRequest.FormElements = *elements
		diagnostics.Append(diags...)
	}

	return apiRequest, diagnostics
}

// ToPatchOperations compares the plan (m) with the current state and generates JSON Patch operations
// for fields that have changed.
func (m *formDefinitionModel) ToPatchOperations(ctx context.Context, state *formDefinitionModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var patchOps []client.JSONPatchOperation

	// Compare name
	if !m.Name.Equal(state.Name) {
		patchOps = append(patchOps, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	// Compare description
	if !m.Description.Equal(state.Description) {
		patchOps = append(patchOps, client.NewReplacePatch("/description", m.Description.ValueString()))
	}

	// Compare owner
	if !reflect.DeepEqual(*m.Owner, *state.Owner) {
		ownerAPI, diags := m.Owner.ToAPI(ctx)
		diagnostics.Append(diags...)
		patchOps = append(patchOps, client.NewReplacePatch("/owner", ownerAPI))
	}

	// Compare usedBy
	if !m.UsedBy.Equal(state.UsedBy) {
		usedBy, diags := common.MapListToAPI(ctx, m.UsedBy, common.NewObjectRefToAPI)
		diagnostics.Append(diags...)
		patchOps = append(patchOps, client.NewReplacePatch("/usedBy", usedBy))
	}

	// Compare formInput
	if !m.FormInput.Equal(state.FormInput) {
		formInput, diags := common.MapListToAPI(ctx, m.FormInput, NewFormInputToAPI)
		diagnostics.Append(diags...)
		patchOps = append(patchOps, client.NewReplacePatch("/formInput", formInput))
	}

	// Compare formConditions
	if !m.FormConditions.Equal(state.FormConditions) {
		formConditions, diags := common.MapListToAPI(ctx, m.FormConditions, FormConditionToAPI)
		diagnostics.Append(diags...)
		patchOps = append(patchOps, client.NewReplacePatch("/formConditions", formConditions))
	}

	// Compare formElements - parse JSON and send as structured data
	if !m.FormElements.Equal(state.FormElements) {
		formElementsAPI, diags := common.UnmarshalJSONField[[]client.FormElementAPI](m.FormElements)
		if formElementsAPI != nil {
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/formElements", *formElementsAPI))
		}
	}

	return patchOps, diagnostics
}
