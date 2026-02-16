// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package form_definition

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// formDefinitionOwnerModel represents the owner of a form definition in Terraform state.
type formDefinitionOwnerModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// formDefinitionUsedByModel represents a resource that references the form definition in Terraform state.
type formDefinitionUsedByModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// formDefinitionFormInputModel represents a form input field in Terraform state.
type formDefinitionFormInputModel struct {
	ID          types.String `tfsdk:"id"`
	Type        types.String `tfsdk:"type"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
}

// formDefinitionFormConditionsModel represents a form condition in Terraform state.
type formDefinitionFormConditionModel struct {
	RuleOperator types.String                             `tfsdk:"rule_operator"` // Has either "AND" or "OR"
	Rules        []formDefinitionFormConditionRuleModel   `tfsdk:"rules"`
	Effects      []formDefinitionFormConditionEffectModel `tfsdk:"effects"`
}

// formDefinitionFormConditionRuleModel represents a rule within a form condition.
type formDefinitionFormConditionRuleModel struct {
	SourceType types.String `tfsdk:"source_type"`
	Source     types.String `tfsdk:"source"`
	Operator   types.String `tfsdk:"operator"`
	ValueType  types.String `tfsdk:"value_type"`
	Value      types.String `tfsdk:"value"`
}

// formDefinitionFormConditionEffectModel represents the effect of a form condition in Terraform state.
type formDefinitionFormConditionEffectModel struct {
	EffectType types.String                                 `tfsdk:"effect_type"`
	Config     formDefinitionFormConditionEffectConfigModel `tfsdk:"config"`
}

// formDefinitionFormConditionEffectConfigModel represents the configuration for a form condition effect in Terraform state.
type formDefinitionFormConditionEffectConfigModel struct {
	DefaultValueLabel types.String `tfsdk:"default_value_label"`
	Element           types.String `tfsdk:"element"`
}

// formDefinitionModel represents the Terraform state for a SailPoint form definition.
type formDefinitionModel struct {
	ID             types.String                       `tfsdk:"id"`
	Name           types.String                       `tfsdk:"name"`
	Description    types.String                       `tfsdk:"description"`
	Owner          *formDefinitionOwnerModel          `tfsdk:"owner"`
	UsedBy         []formDefinitionUsedByModel        `tfsdk:"used_by"`
	FormInput      []formDefinitionFormInputModel     `tfsdk:"form_input"`
	FormElements   jsontypes.Normalized               `tfsdk:"form_elements"`
	FormConditions []formDefinitionFormConditionModel `tfsdk:"form_conditions"`
	Created        types.String                       `tfsdk:"created"`
	Modified       types.String                       `tfsdk:"modified"`
}

// FromSailPointAPI maps fields from the API response to the Terraform model.
func (m *formDefinitionModel) FromSailPointAPI(ctx context.Context, api client.FormDefinitionAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = types.StringValue(api.Description)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)

	// Map owner
	m.Owner = &formDefinitionOwnerModel{
		Type: types.StringValue(api.Owner.Type),
		ID:   types.StringValue(api.Owner.ID),
		Name: types.StringValue(api.Owner.Name),
	}

	// Map usedBy
	if len(api.UsedBy) > 0 {
		m.UsedBy = make([]formDefinitionUsedByModel, len(api.UsedBy))
		for i := range api.UsedBy {
			m.UsedBy[i] = formDefinitionUsedByModel{
				Type: types.StringValue(api.UsedBy[i].Type),
				ID:   types.StringValue(api.UsedBy[i].ID),
				Name: types.StringValue(api.UsedBy[i].Name),
			}
		}
	}

	// Map formInput
	if len(api.FormInput) > 0 {
		m.FormInput = make([]formDefinitionFormInputModel, len(api.FormInput))
		for i := range api.FormInput {
			m.FormInput[i] = formDefinitionFormInputModel{
				ID:          types.StringValue(api.FormInput[i].ID),
				Type:        types.StringValue(api.FormInput[i].Type),
				Label:       types.StringValue(api.FormInput[i].Label),
				Description: types.StringValue(api.FormInput[i].Description),
			}
		}
	}

	// Map formConditions
	if len(api.FormConditions) > 0 {
		m.FormConditions = make([]formDefinitionFormConditionModel, len(api.FormConditions))
		for i := range api.FormConditions {
			// Map rules
			rules := make([]formDefinitionFormConditionRuleModel, len(api.FormConditions[i].Rules))
			for j := range api.FormConditions[i].Rules {
				rules[j] = formDefinitionFormConditionRuleModel{
					SourceType: types.StringValue(api.FormConditions[i].Rules[j].SourceType),
					Source:     types.StringValue(api.FormConditions[i].Rules[j].Source),
					Operator:   types.StringValue(api.FormConditions[i].Rules[j].Operator),
					ValueType:  types.StringValue(api.FormConditions[i].Rules[j].ValueType),
					Value:      types.StringValue(api.FormConditions[i].Rules[j].Value),
				}
			}

			// Map effects
			effects := make([]formDefinitionFormConditionEffectModel, len(api.FormConditions[i].Effects))
			for k := range api.FormConditions[i].Effects {
				effects[k] = formDefinitionFormConditionEffectModel{
					EffectType: types.StringValue(api.FormConditions[i].Effects[k].EffectType),
					Config: formDefinitionFormConditionEffectConfigModel{
						DefaultValueLabel: types.StringValue(api.FormConditions[i].Effects[k].Config.DefaultValueLabel),
						Element:           types.StringValue(api.FormConditions[i].Effects[k].Config.Element),
					},
				}
			}

			m.FormConditions[i] = formDefinitionFormConditionModel{
				RuleOperator: types.StringValue(api.FormConditions[i].RuleOperator),
				Rules:        rules,
				Effects:      effects,
			}
		}
	}

	// Map formElements as JSON (use empty array "[]" instead of null for consistency)
	if api.FormElements != nil {
		formElementsBytes, err := json.Marshal(api.FormElements)
		if err != nil {
			diagnostics.AddError(
				"Error Mapping Form Elements",
				"Could not marshal form elements to JSON: "+err.Error(),
			)
			return diagnostics
		}
		m.FormElements = jsontypes.NewNormalizedValue(string(formElementsBytes))
	} else {
		m.FormElements = jsontypes.NewNormalizedValue("[]")
	}

	return diagnostics
}

// ToAPICreateRequest maps fields from the Terraform model to the API create request.
func (m *formDefinitionModel) ToAPICreateRequest(ctx context.Context) (client.FormDefinitionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiRequest := client.FormDefinitionAPI{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
		Owner: client.ObjectRefAPI{
			Type: m.Owner.Type.ValueString(),
			ID:   m.Owner.ID.ValueString(),
		},
	}

	// Parse formInput
	if len(m.FormInput) > 0 {
		var formInput []client.FormInputAPI
		for i := range m.FormInput {
			formInput = append(formInput, client.FormInputAPI{
				ID:          m.FormInput[i].ID.ValueString(),
				Type:        m.FormInput[i].Type.ValueString(),
				Label:       m.FormInput[i].Label.ValueString(),
				Description: m.FormInput[i].Description.ValueString(),
			})
		}
		apiRequest.FormInput = formInput
	}

	// Parse formConditions
	if len(m.FormConditions) > 0 {
		var formConditions []client.FormConditionAPI
		for i := range m.FormConditions {
			// Parse rules
			var rules []client.FormConditionRuleAPI
			for j := range m.FormConditions[i].Rules {
				rules = append(rules, client.FormConditionRuleAPI{
					SourceType: m.FormConditions[i].Rules[j].SourceType.ValueString(),
					Source:     m.FormConditions[i].Rules[j].Source.ValueString(),
					Operator:   m.FormConditions[i].Rules[j].Operator.ValueString(),
					ValueType:  m.FormConditions[i].Rules[j].ValueType.ValueString(),
					Value:      m.FormConditions[i].Rules[j].Value.ValueString(),
				})
			}

			// Parse effects
			var effects []client.FormConditionEffectAPI
			for k := range m.FormConditions[i].Effects {
				effects = append(effects, client.FormConditionEffectAPI{
					EffectType: m.FormConditions[i].Effects[k].EffectType.ValueString(),
					Config: client.FormConditionEffectConfigAPI{
						DefaultValueLabel: m.FormConditions[i].Effects[k].Config.DefaultValueLabel.ValueString(),
						Element:           m.FormConditions[i].Effects[k].Config.Element.ValueString(),
					},
				})
			}

			formConditions = append(formConditions, client.FormConditionAPI{
				RuleOperator: m.FormConditions[i].RuleOperator.ValueString(),
				Rules:        rules,
				Effects:      effects,
			})
		}
		apiRequest.FormConditions = formConditions
	}

	// Parse formElements from JSON
	if !m.FormElements.IsNull() && !m.FormElements.IsUnknown() {
		var formElements []client.FormElementAPI
		if err := json.Unmarshal([]byte(m.FormElements.ValueString()), &formElements); err != nil {
			diagnostics.AddError(
				"Error Parsing Form Elements",
				"Could not parse form elements JSON: "+err.Error(),
			)
			return apiRequest, diagnostics
		}
		apiRequest.FormElements = formElements
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
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/name",
			Value: m.Name.ValueString(),
		})
	}

	// Compare description
	if !m.Description.Equal(state.Description) {
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/description",
			Value: m.Description.ValueString(),
		})
	}

	// Compare owner
	if !reflect.DeepEqual(m.Owner, state.Owner) {
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:   "replace",
			Path: "/owner",
			Value: client.ObjectRefAPI{
				Type: m.Owner.Type.ValueString(),
				ID:   m.Owner.ID.ValueString(),
			},
		})
	}

	// Compare usedBy - this is read-only and should not be updated, so we skip it
	if !reflect.DeepEqual(m.UsedBy, state.UsedBy) {
		usedBy := []client.ObjectRefAPI{}
		for i := range m.UsedBy {
			usedBy = append(usedBy, client.ObjectRefAPI{
				Type: m.UsedBy[i].Type.ValueString(),
				ID:   m.UsedBy[i].ID.ValueString(),
			})
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/usedBy",
			Value: usedBy,
		})
	}

	// Compare formInput
	if !reflect.DeepEqual(m.FormInput, state.FormInput) {
		formInput := []client.FormInputAPI{}
		for i := range m.FormInput {
			formInput = append(formInput, client.FormInputAPI{
				ID:          m.FormInput[i].ID.ValueString(),
				Type:        m.FormInput[i].Type.ValueString(),
				Label:       m.FormInput[i].Label.ValueString(),
				Description: m.FormInput[i].Description.ValueString(),
			})
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formInput",
			Value: formInput,
		})
	}

	// Compare formConditions
	if !reflect.DeepEqual(m.FormConditions, state.FormConditions) {
		formConditions := []client.FormConditionAPI{}
		for i := range m.FormConditions {
			// Parse rules
			rules := []client.FormConditionRuleAPI{}
			for j := range m.FormConditions[i].Rules {
				rules = append(rules, client.FormConditionRuleAPI{
					SourceType: m.FormConditions[i].Rules[j].SourceType.ValueString(),
					Source:     m.FormConditions[i].Rules[j].Source.ValueString(),
					Operator:   m.FormConditions[i].Rules[j].Operator.ValueString(),
					ValueType:  m.FormConditions[i].Rules[j].ValueType.ValueString(),
					Value:      m.FormConditions[i].Rules[j].Value.ValueString(),
				})
			}

			// Parse effects
			effects := []client.FormConditionEffectAPI{}
			for k := range m.FormConditions[i].Effects {
				effects = append(effects, client.FormConditionEffectAPI{
					EffectType: m.FormConditions[i].Effects[k].EffectType.ValueString(),
					Config: client.FormConditionEffectConfigAPI{
						DefaultValueLabel: m.FormConditions[i].Effects[k].Config.DefaultValueLabel.ValueString(),
						Element:           m.FormConditions[i].Effects[k].Config.Element.ValueString(),
					},
				})
			}

			formConditions = append(formConditions, client.FormConditionAPI{
				RuleOperator: m.FormConditions[i].RuleOperator.ValueString(),
				Rules:        rules,
				Effects:      effects,
			})
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formConditions",
			Value: formConditions,
		})
	}

	// Compare formElements - parse JSON and send as structured data
	if !m.FormElements.Equal(state.FormElements) {
		var formElementsAPI []client.FormElementAPI
		if !m.FormElements.IsNull() && !m.FormElements.IsUnknown() {
			if err := json.Unmarshal([]byte(m.FormElements.ValueString()), &formElementsAPI); err != nil {
				diagnostics.AddError(
					"Error Parsing Form Elements",
					"Could not parse form elements JSON: "+err.Error(),
				)
				return patchOps, diagnostics
			}
		}
		// Use empty array if null to avoid API issues
		if formElementsAPI == nil {
			formElementsAPI = []client.FormElementAPI{}
		}
		patchOps = append(patchOps, client.JSONPatchOperation{
			Op:    "replace",
			Path:  "/formElements",
			Value: formElementsAPI,
		})
	}

	return patchOps, diagnostics
}
