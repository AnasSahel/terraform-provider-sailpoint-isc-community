// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	ObjectRefObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"id":   types.StringType,
		"name": types.StringType,
	}}
)

type ObjectRefModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewObjectRefFromAPI(ctx context.Context, api client.ObjectRefAPI) (ObjectRefModel, diag.Diagnostics) {
	var m ObjectRefModel
	diags := m.FromAPI(ctx, api)
	return m, diags
}

func NewObjectRefFromAPIPtr(ctx context.Context, api client.ObjectRefAPI) (*ObjectRefModel, diag.Diagnostics) {
	var m ObjectRefModel
	diags := m.FromAPI(ctx, api)
	return &m, diags
}

func (m *ObjectRefModel) FromAPI(ctx context.Context, api client.ObjectRefAPI) diag.Diagnostics {
	m.Type = types.StringValue(api.Type)
	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	return nil
}

func (m *ObjectRefModel) ToAPI(ctx context.Context) (client.ObjectRefAPI, diag.Diagnostics) {
	return client.ObjectRefAPI{
		Type: m.Type.ValueString(),
		ID:   m.ID.ValueString(),
		Name: m.Name.ValueString(),
	}, nil
}
