// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_correlation_config

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type correlationAttributeAssignmentModel struct {
	Property     types.String `tfsdk:"property"`
	Value        types.String `tfsdk:"value"`
	Operation    types.String `tfsdk:"operation"`
	Complex      types.Bool   `tfsdk:"complex"`
	IgnoreCase   types.Bool   `tfsdk:"ignore_case"`
	MatchMode    types.String `tfsdk:"match_mode"`
	FilterString types.String `tfsdk:"filter_string"`
}

type sourceCorrelationConfigModel struct {
	ID                   types.String                          `tfsdk:"id"`
	SourceID             types.String                          `tfsdk:"source_id"`
	Name                 types.String                          `tfsdk:"name"`
	AttributeAssignments []correlationAttributeAssignmentModel `tfsdk:"attribute_assignments"`
}

// ---------------------------------------------------------------------------
// FromAPI
// ---------------------------------------------------------------------------

func (m *sourceCorrelationConfigModel) FromAPI(_ context.Context, api *client.SourceCorrelationConfigAPI) diag.Diagnostics {
	// id mirrors source_id for import convenience.
	m.ID = m.SourceID

	if api.Name != nil {
		m.Name = types.StringValue(*api.Name)
	} else {
		m.Name = types.StringNull()
	}

	if len(api.AttributeAssignments) > 0 {
		m.AttributeAssignments = make([]correlationAttributeAssignmentModel, 0, len(api.AttributeAssignments))
		for _, a := range api.AttributeAssignments {
			m.AttributeAssignments = append(m.AttributeAssignments, assignmentFromAPI(a))
		}
	} else {
		m.AttributeAssignments = nil
	}

	return nil
}

// ---------------------------------------------------------------------------
// ToAPI
// ---------------------------------------------------------------------------

func (m *sourceCorrelationConfigModel) ToAPI() *client.SourceCorrelationConfigAPI {
	cfg := &client.SourceCorrelationConfigAPI{}

	assignments := make([]client.CorrelationAttributeAssignmentAPI, 0, len(m.AttributeAssignments))
	for _, a := range m.AttributeAssignments {
		assignments = append(assignments, assignmentToAPI(a))
	}
	cfg.AttributeAssignments = assignments

	return cfg
}

// ---------------------------------------------------------------------------
// NeedsUpdate compares plan to state — returns true when attribute_assignments differ.
// ---------------------------------------------------------------------------

func (m *sourceCorrelationConfigModel) NeedsUpdate(state *sourceCorrelationConfigModel) bool {
	return !reflect.DeepEqual(m.AttributeAssignments, state.AttributeAssignments)
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func assignmentFromAPI(a client.CorrelationAttributeAssignmentAPI) correlationAttributeAssignmentModel {
	m := correlationAttributeAssignmentModel{
		Property:     types.StringValue(a.Property),
		Value:        types.StringValue(a.Value),
		Operation:    stringPtrToTF(a.Operation),
		Complex:      boolPtrToTF(a.Complex),
		IgnoreCase:   boolPtrToTF(a.IgnoreCase),
		MatchMode:    stringPtrToTF(a.MatchMode),
		FilterString: stringPtrToTF(a.FilterString),
	}
	return m
}

func assignmentToAPI(m correlationAttributeAssignmentModel) client.CorrelationAttributeAssignmentAPI {
	a := client.CorrelationAttributeAssignmentAPI{
		Property: m.Property.ValueString(),
		Value:    m.Value.ValueString(),
	}
	if !m.Operation.IsNull() && !m.Operation.IsUnknown() {
		v := m.Operation.ValueString()
		a.Operation = &v
	}
	if !m.Complex.IsNull() && !m.Complex.IsUnknown() {
		v := m.Complex.ValueBool()
		a.Complex = &v
	}
	if !m.IgnoreCase.IsNull() && !m.IgnoreCase.IsUnknown() {
		v := m.IgnoreCase.ValueBool()
		a.IgnoreCase = &v
	}
	if !m.MatchMode.IsNull() && !m.MatchMode.IsUnknown() {
		v := m.MatchMode.ValueString()
		a.MatchMode = &v
	}
	if !m.FilterString.IsNull() && !m.FilterString.IsUnknown() {
		v := m.FilterString.ValueString()
		a.FilterString = &v
	}
	return a
}

func stringPtrToTF(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

func boolPtrToTF(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}
