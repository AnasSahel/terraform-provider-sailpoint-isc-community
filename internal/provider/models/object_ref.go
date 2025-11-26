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

	objRef := &client.ObjectRef{
		Type: ref.Type.ValueString(),
		ID:   ref.ID.ValueString(),
	}

	// Only include Name if it's not null/unknown
	if !ref.Name.IsNull() && !ref.Name.IsUnknown() {
		objRef.Name = ref.Name.ValueString()
	}

	return objRef
}

// ConvertToSailPoint implements ModelWithSailPointConversionMethods.
func (o *ObjectRef) ConvertToSailPoint(ctx context.Context) client.ObjectRef {
	if o == nil {
		return client.ObjectRef{}
	}

	ref := client.ObjectRef{
		Type: o.Type.ValueString(),
		ID:   o.ID.ValueString(),
	}

	// Only include Name if it's not null/unknown
	if !o.Name.IsNull() && !o.Name.IsUnknown() {
		ref.Name = o.Name.ValueString()
	}

	return ref
}

// ConvertFromSailPointForDataSource implements ModelWithSourceConversionMethods.
func (o *ObjectRef) ConvertFromSailPointForDataSource(ctx context.Context, source *client.ObjectRef) {
	if o == nil || source == nil {
		return
	}

	o.Type = types.StringValue(source.Type)
	o.ID = types.StringValue(source.ID)
	// Only set Name if it's not empty to avoid inconsistent state errors
	if source.Name != "" {
		o.Name = types.StringValue(source.Name)
	} else {
		o.Name = types.StringNull()
	}
}

// ConvertFromSailPointForResource implements ModelWithSourceConversionMethods.
func (o *ObjectRef) ConvertFromSailPointForResource(ctx context.Context, source *client.ObjectRef) {
	if o == nil || source == nil {
		return
	}

	o.Type = types.StringValue(source.Type)
	o.ID = types.StringValue(source.ID)
	// Only set Name if it's not empty to avoid inconsistent state errors
	if source.Name != "" {
		o.Name = types.StringValue(source.Name)
	} else {
		o.Name = types.StringNull()
	}
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
