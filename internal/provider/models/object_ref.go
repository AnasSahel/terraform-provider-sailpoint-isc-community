// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ModelWithSourceTerraformConversionMethods[client.ObjectRef] = &ObjectRef{}
	_ ModelWithSailPointConversionMethods[client.ObjectRef]       = &ObjectRef{}
)

type ObjectRef struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func NewObjectRefFromSailPoint(apiRef *client.ObjectRef) *ObjectRef {
	if apiRef == nil {
		return nil
	}

	return &ObjectRef{
		Type: types.StringValue(apiRef.Type),
		ID:   types.StringValue(apiRef.ID),
		Name: types.StringValue(apiRef.Name),
	}
}

func NewObjectRefFromTerraform(ref *ObjectRef) *client.ObjectRef {
	if ref == nil {
		return nil
	}

	return &client.ObjectRef{
		Type: ref.Type.ValueString(),
		ID:   ref.ID.ValueString(),
		Name: ref.Name.ValueString(),
	}
}

// ConvertToSailPoint implements ModelWithSailPointConversionMethods.
func (o *ObjectRef) ConvertToSailPoint(ctx context.Context) client.ObjectRef {
	if o == nil {
		return client.ObjectRef{}
	}

	return client.ObjectRef{
		Type: o.Type.ValueString(),
		ID:   o.ID.ValueString(),
		Name: o.Name.ValueString(),
	}
}

// ConvertFromSailPointForDataSource implements ModelWithSourceConversionMethods.
func (o *ObjectRef) ConvertFromSailPointForDataSource(ctx context.Context, source *client.ObjectRef) {
	if o == nil || source == nil {
		return
	}

	o.Type = types.StringValue(source.Type)
	o.ID = types.StringValue(source.ID)
	o.Name = types.StringValue(source.Name)
}

// ConvertFromSailPointForResource implements ModelWithSourceConversionMethods.
func (o *ObjectRef) ConvertFromSailPointForResource(ctx context.Context, source *client.ObjectRef) {
	if o == nil || source == nil {
		return
	}

	o.Type = types.StringValue(source.Type)
	o.ID = types.StringValue(source.ID)
	o.Name = types.StringValue(source.Name)
}

func (o *ObjectRef) Equals(other *ObjectRef) bool {
	if o == nil && other == nil {
		return true
	}
	if o == nil || other == nil {
		return false
	}
	return o.Type.ValueString() == other.Type.ValueString() &&
		o.ID.ValueString() == other.ID.ValueString() &&
		o.Name.ValueString() == other.Name.ValueString()
}
