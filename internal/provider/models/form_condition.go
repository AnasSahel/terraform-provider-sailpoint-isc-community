// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FormCondition represents a conditional logic rule for a form.
// Conditions can show/hide/enable/disable form elements based on other element values.
type FormCondition struct {
	RuleOperator types.String          `tfsdk:"rule_operator"`
	Rules        []FormConditionRule   `tfsdk:"rules"`
	Effects      []FormConditionEffect `tfsdk:"effects"`
}

// FormConditionRule represents a single rule within a form condition.
type FormConditionRule struct {
	SourceType types.String `tfsdk:"source_type"`
	Source     types.String `tfsdk:"source"`
	Operator   types.String `tfsdk:"operator"`
	ValueType  types.String `tfsdk:"value_type"`
	Value      types.String `tfsdk:"value"`
}

// FormConditionEffect represents an effect applied when condition rules are met.
type FormConditionEffect struct {
	EffectType types.String               `tfsdk:"effect_type"`
	Config     *FormConditionEffectConfig `tfsdk:"config"`
}

// FormConditionEffectConfig represents the configuration for a form condition effect.
type FormConditionEffectConfig struct {
	DefaultValueLabel types.String `tfsdk:"default_value_label"`
	Element           types.String `tfsdk:"element"`
}

// ConvertToSailPoint converts the Terraform model to a SailPoint API map.
func (fc *FormCondition) ConvertToSailPoint(ctx context.Context) map[string]interface{} {
	if fc == nil {
		return nil
	}

	result := map[string]interface{}{}

	if !fc.RuleOperator.IsNull() && !fc.RuleOperator.IsUnknown() {
		result["ruleOperator"] = fc.RuleOperator.ValueString()
	}

	if len(fc.Rules) > 0 {
		rules := make([]map[string]interface{}, len(fc.Rules))
		for i, rule := range fc.Rules {
			rules[i] = rule.ConvertToSailPoint(ctx)
		}
		result["rules"] = rules
	}

	if len(fc.Effects) > 0 {
		effects := make([]map[string]interface{}, len(fc.Effects))
		for i, effect := range fc.Effects {
			effects[i] = effect.ConvertToSailPoint(ctx)
		}
		result["effects"] = effects
	}

	return result
}

// ConvertToSailPoint converts a FormConditionRule to a SailPoint API map.
func (fcr *FormConditionRule) ConvertToSailPoint(ctx context.Context) map[string]interface{} {
	if fcr == nil {
		return nil
	}

	result := map[string]interface{}{}

	if !fcr.SourceType.IsNull() && !fcr.SourceType.IsUnknown() {
		result["sourceType"] = fcr.SourceType.ValueString()
	}

	if !fcr.Source.IsNull() && !fcr.Source.IsUnknown() {
		result["source"] = fcr.Source.ValueString()
	}

	if !fcr.Operator.IsNull() && !fcr.Operator.IsUnknown() {
		result["operator"] = fcr.Operator.ValueString()
	}

	if !fcr.ValueType.IsNull() && !fcr.ValueType.IsUnknown() {
		result["valueType"] = fcr.ValueType.ValueString()
	}

	if !fcr.Value.IsNull() && !fcr.Value.IsUnknown() {
		result["value"] = fcr.Value.ValueString()
	}

	return result
}

// ConvertToSailPoint converts a FormConditionEffect to a SailPoint API map.
func (fce *FormConditionEffect) ConvertToSailPoint(ctx context.Context) map[string]interface{} {
	if fce == nil {
		return nil
	}

	result := map[string]interface{}{}

	if !fce.EffectType.IsNull() && !fce.EffectType.IsUnknown() {
		result["effectType"] = fce.EffectType.ValueString()
	}

	if fce.Config != nil {
		configMap := make(map[string]interface{})
		if !fce.Config.DefaultValueLabel.IsNull() && !fce.Config.DefaultValueLabel.IsUnknown() {
			configMap["defaultValueLabel"] = fce.Config.DefaultValueLabel.ValueString()
		}
		if !fce.Config.Element.IsNull() && !fce.Config.Element.IsUnknown() {
			configMap["element"] = fce.Config.Element.ValueString()
		}
		if len(configMap) > 0 {
			result["config"] = configMap
		}
	}

	return result
}

// ConvertFromSailPoint converts a SailPoint API map to the Terraform model.
func (fc *FormCondition) ConvertFromSailPoint(ctx context.Context, condition map[string]interface{}) {
	if fc == nil || condition == nil {
		return
	}

	if ruleOperator, ok := condition["ruleOperator"].(string); ok {
		fc.RuleOperator = types.StringValue(ruleOperator)
	} else {
		fc.RuleOperator = types.StringNull()
	}

	if rules, ok := condition["rules"].([]interface{}); ok {
		fcRules := make([]FormConditionRule, len(rules))
		for i, ruleInterface := range rules {
			if ruleMap, ok := ruleInterface.(map[string]interface{}); ok {
				fcRules[i].ConvertFromSailPoint(ctx, ruleMap)
			}
		}
		fc.Rules = fcRules
	} else {
		fc.Rules = []FormConditionRule{}
	}

	if effects, ok := condition["effects"].([]interface{}); ok {
		fcEffects := make([]FormConditionEffect, len(effects))
		for i, effectInterface := range effects {
			if effectMap, ok := effectInterface.(map[string]interface{}); ok {
				fcEffects[i].ConvertFromSailPoint(ctx, effectMap)
			}
		}
		fc.Effects = fcEffects
	} else {
		fc.Effects = []FormConditionEffect{}
	}
}

// ConvertFromSailPoint converts a SailPoint API map to FormConditionRule.
func (fcr *FormConditionRule) ConvertFromSailPoint(ctx context.Context, rule map[string]interface{}) {
	if fcr == nil || rule == nil {
		return
	}

	if sourceType, ok := rule["sourceType"].(string); ok {
		fcr.SourceType = types.StringValue(sourceType)
	} else {
		fcr.SourceType = types.StringNull()
	}

	if source, ok := rule["source"].(string); ok {
		fcr.Source = types.StringValue(source)
	} else {
		fcr.Source = types.StringNull()
	}

	if operator, ok := rule["operator"].(string); ok {
		fcr.Operator = types.StringValue(operator)
	} else {
		fcr.Operator = types.StringNull()
	}

	if valueType, ok := rule["valueType"].(string); ok {
		fcr.ValueType = types.StringValue(valueType)
	} else {
		fcr.ValueType = types.StringNull()
	}

	if value, ok := rule["value"].(string); ok {
		fcr.Value = types.StringValue(value)
	} else {
		fcr.Value = types.StringNull()
	}
}

// ConvertFromSailPoint converts a SailPoint API map to FormConditionEffect.
func (fce *FormConditionEffect) ConvertFromSailPoint(ctx context.Context, effect map[string]interface{}) {
	if fce == nil || effect == nil {
		return
	}

	if effectType, ok := effect["effectType"].(string); ok {
		fce.EffectType = types.StringValue(effectType)
	} else {
		fce.EffectType = types.StringNull()
	}

	if config, ok := effect["config"].(map[string]interface{}); ok {
		fce.Config = &FormConditionEffectConfig{}
		if defaultValueLabel, ok := config["defaultValueLabel"].(string); ok {
			fce.Config.DefaultValueLabel = types.StringValue(defaultValueLabel)
		} else {
			fce.Config.DefaultValueLabel = types.StringNull()
		}
		// Element can be either a string or a number
		if element, ok := config["element"].(string); ok {
			fce.Config.Element = types.StringValue(element)
		} else if elementNum, ok := config["element"].(float64); ok {
			// Convert number to string
			fce.Config.Element = types.StringValue(fmt.Sprintf("%.0f", elementNum))
		} else {
			fce.Config.Element = types.StringNull()
		}
	}
}
