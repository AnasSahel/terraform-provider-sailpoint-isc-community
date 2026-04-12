// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// segmentValueModel represents a typed value in an EQUALS expression.
type segmentValueModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// segmentExpressionLeafModel represents a leaf expression (max 1 additional nesting level).
// Leaf nodes cannot have further children.
type segmentExpressionLeafModel struct {
	Operator  types.String       `tfsdk:"operator"`
	Attribute types.String       `tfsdk:"attribute"`
	Value     *segmentValueModel `tfsdk:"value"`
}

// segmentExpressionModel represents the root expression node. It may be:
//   - a leaf (operator=EQUALS with attribute+value, no children), or
//   - a branch (operator=AND with children, no attribute/value).
type segmentExpressionModel struct {
	Operator  types.String                 `tfsdk:"operator"`
	Attribute types.String                 `tfsdk:"attribute"`
	Value     *segmentValueModel           `tfsdk:"value"`
	Children  []segmentExpressionLeafModel `tfsdk:"children"`
}

// visibilityCriteriaModel wraps the segment's root expression.
type visibilityCriteriaModel struct {
	Expression *segmentExpressionModel `tfsdk:"expression"`
}

// segmentModel represents the Terraform state for a Segment resource.
type segmentModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	Description        types.String             `tfsdk:"description"`
	Active             types.Bool               `tfsdk:"active"`
	Owner              *common.ObjectRefModel   `tfsdk:"owner"`
	VisibilityCriteria *visibilityCriteriaModel `tfsdk:"visibility_criteria"`
	Created            types.String             `tfsdk:"created"`
	Modified           types.String             `tfsdk:"modified"`
}

// FromAPI maps the API response into the Terraform state.
func (m *segmentModel) FromAPI(ctx context.Context, api *client.SegmentAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNull(api.Description)

	if api.Active != nil {
		m.Active = types.BoolValue(*api.Active)
	} else {
		m.Active = types.BoolNull()
	}

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

	if api.Owner != nil {
		owner, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
		m.Owner = owner
	} else {
		m.Owner = nil
	}

	if api.VisibilityCriteria != nil {
		m.VisibilityCriteria = visibilityCriteriaFromAPI(api.VisibilityCriteria)
	} else {
		m.VisibilityCriteria = nil
	}

	return diagnostics
}

// ToAPI maps the Terraform state into an API create/replace payload.
func (m *segmentModel) ToAPI(ctx context.Context) (*client.SegmentAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.SegmentAPI{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		desc := m.Description.ValueString()
		api.Description = &desc
	}

	if !m.Active.IsNull() && !m.Active.IsUnknown() {
		active := m.Active.ValueBool()
		api.Active = &active
	}

	if m.Owner != nil {
		owner, diags := common.NewObjectRefToAPIPtr(ctx, *m.Owner)
		diagnostics.Append(diags...)
		api.Owner = owner
	}

	if m.VisibilityCriteria != nil {
		api.VisibilityCriteria = visibilityCriteriaToAPI(m.VisibilityCriteria)
	}

	return api, diagnostics
}

// ToPatchOperations compares the plan (m) against state and returns JSON Patch ops for changed fields.
func (m *segmentModel) ToPatchOperations(ctx context.Context, state *segmentModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
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

	if !m.Active.Equal(state.Active) {
		if !m.Active.IsNull() {
			ops = append(ops, client.NewReplacePatch("/active", m.Active.ValueBool()))
		} else {
			ops = append(ops, client.NewRemovePatch("/active"))
		}
	}

	if !reflect.DeepEqual(m.Owner, state.Owner) {
		if m.Owner != nil {
			ownerAPI, diags := common.NewObjectRefToAPIPtr(ctx, *m.Owner)
			diagnostics.Append(diags...)
			ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
		} else {
			ops = append(ops, client.NewRemovePatch("/owner"))
		}
	}

	if !reflect.DeepEqual(m.VisibilityCriteria, state.VisibilityCriteria) {
		if m.VisibilityCriteria != nil {
			ops = append(ops, client.NewReplacePatch("/visibilityCriteria", visibilityCriteriaToAPI(m.VisibilityCriteria)))
		} else {
			ops = append(ops, client.NewRemovePatch("/visibilityCriteria"))
		}
	}

	return ops, diagnostics
}

// visibilityCriteriaFromAPI converts the API visibility criteria into the Terraform model.
func visibilityCriteriaFromAPI(api *client.VisibilityCriteriaAPI) *visibilityCriteriaModel {
	if api == nil {
		return nil
	}
	vc := &visibilityCriteriaModel{}
	if api.Expression != nil {
		vc.Expression = expressionFromAPI(api.Expression)
	}
	return vc
}

// expressionFromAPI converts a root expression node into the Terraform model.
// Children (if any) are flattened into the leaf-only model type.
func expressionFromAPI(api *client.SegmentExpressionAPI) *segmentExpressionModel {
	if api == nil {
		return nil
	}
	expr := &segmentExpressionModel{
		Operator:  types.StringValue(api.Operator),
		Attribute: stringPtrToTF(api.Attribute),
		Value:     valueFromAPI(api.Value),
	}
	if len(api.Children) > 0 {
		expr.Children = make([]segmentExpressionLeafModel, 0, len(api.Children))
		for i := range api.Children {
			child := &api.Children[i]
			expr.Children = append(expr.Children, segmentExpressionLeafModel{
				Operator:  types.StringValue(child.Operator),
				Attribute: stringPtrToTF(child.Attribute),
				Value:     valueFromAPI(child.Value),
			})
		}
	}
	return expr
}

// valueFromAPI converts a typed value into the Terraform model.
func valueFromAPI(api *client.SegmentValueAPI) *segmentValueModel {
	if api == nil {
		return nil
	}
	return &segmentValueModel{
		Type:  types.StringValue(api.Type),
		Value: types.StringValue(api.Value),
	}
}

// visibilityCriteriaToAPI converts the Terraform visibility criteria into an API payload.
func visibilityCriteriaToAPI(m *visibilityCriteriaModel) *client.VisibilityCriteriaAPI {
	if m == nil {
		return nil
	}
	return &client.VisibilityCriteriaAPI{
		Expression: expressionToAPI(m.Expression),
	}
}

// expressionToAPI converts the root expression into an API payload.
func expressionToAPI(m *segmentExpressionModel) *client.SegmentExpressionAPI {
	if m == nil {
		return nil
	}
	api := &client.SegmentExpressionAPI{
		Operator:  m.Operator.ValueString(),
		Attribute: tfToStringPtr(m.Attribute),
		Value:     valueToAPI(m.Value),
	}
	if len(m.Children) > 0 {
		api.Children = make([]client.SegmentExpressionAPI, 0, len(m.Children))
		for i := range m.Children {
			child := &m.Children[i]
			api.Children = append(api.Children, client.SegmentExpressionAPI{
				Operator:  child.Operator.ValueString(),
				Attribute: tfToStringPtr(child.Attribute),
				Value:     valueToAPI(child.Value),
			})
		}
	}
	return api
}

// valueToAPI converts a typed value into the API payload.
func valueToAPI(m *segmentValueModel) *client.SegmentValueAPI {
	if m == nil {
		return nil
	}
	return &client.SegmentValueAPI{
		Type:  m.Type.ValueString(),
		Value: m.Value.ValueString(),
	}
}

// stringPtrToTF converts *string to types.String (nil → null).
func stringPtrToTF(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

// tfToStringPtr converts types.String to *string (null/unknown → nil).
func tfToStringPtr(s types.String) *string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	v := s.ValueString()
	return &v
}
