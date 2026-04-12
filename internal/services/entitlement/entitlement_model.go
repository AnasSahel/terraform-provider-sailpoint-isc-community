// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// entitlementModel represents the Terraform state for an Entitlement resource.
type entitlementModel struct {
	ID                     types.String           `tfsdk:"id"`
	Name                   types.String           `tfsdk:"name"`
	Description            types.String           `tfsdk:"description"`
	Attribute              types.String           `tfsdk:"attribute"`
	Value                  types.String           `tfsdk:"value"`
	SourceSchemaObjectType types.String           `tfsdk:"source_schema_object_type"`
	Privileged             types.Bool             `tfsdk:"privileged"`
	CloudGoverned          types.Bool             `tfsdk:"cloud_governed"`
	Requestable            types.Bool             `tfsdk:"requestable"`
	Owner                  *common.ObjectRefModel `tfsdk:"owner"`
	Source                 *common.ObjectRefModel `tfsdk:"source"`
	Segments               types.Set              `tfsdk:"segments"`
	ManuallyUpdatedFields  types.Map              `tfsdk:"manually_updated_fields"`
	Created                types.String           `tfsdk:"created"`
	Modified               types.String           `tfsdk:"modified"`
}

// FromAPI maps the API response into the Terraform state.
func (m *entitlementModel) FromAPI(ctx context.Context, api *client.EntitlementAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNull(api.Description)
	m.Attribute = types.StringValue(api.Attribute)
	m.Value = types.StringValue(api.Value)
	m.SourceSchemaObjectType = types.StringValue(api.SourceSchemaObjectType)

	m.Privileged = boolPtrToTF(api.Privileged)
	m.CloudGoverned = boolPtrToTF(api.CloudGoverned)
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

	if api.Owner != nil {
		owner, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
		m.Owner = owner
	} else {
		m.Owner = nil
	}

	if api.Source != nil {
		source, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Source)
		diagnostics.Append(diags...)
		m.Source = source
	} else {
		m.Source = nil
	}

	if api.Segments != nil {
		segs, diags := types.SetValueFrom(ctx, types.StringType, api.Segments)
		diagnostics.Append(diags...)
		m.Segments = segs
	} else {
		m.Segments = types.SetNull(types.StringType)
	}

	if api.ManuallyUpdatedFields != nil {
		muf, diags := types.MapValueFrom(ctx, types.BoolType, api.ManuallyUpdatedFields)
		diagnostics.Append(diags...)
		m.ManuallyUpdatedFields = muf
	} else {
		m.ManuallyUpdatedFields = types.MapNull(types.BoolType)
	}

	return diagnostics
}

// ToPatchOperations compares the plan (m) against state and returns JSON Patch ops for changed fields.
// Only patchable fields are considered: name, description, requestable, privileged, owner, segments.
func (m *entitlementModel) ToPatchOperations(ctx context.Context, state *entitlementModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var ops []client.JSONPatchOperation

	if !m.Name.Equal(state.Name) && !m.Name.IsNull() && !m.Name.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	if !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() && !m.Description.IsUnknown() {
			ops = append(ops, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			ops = append(ops, client.NewRemovePatch("/description"))
		}
	}

	if !m.Requestable.Equal(state.Requestable) && !m.Requestable.IsNull() && !m.Requestable.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/requestable", m.Requestable.ValueBool()))
	}

	if !m.Privileged.Equal(state.Privileged) && !m.Privileged.IsNull() && !m.Privileged.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/privileged", m.Privileged.ValueBool()))
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

	if !m.Segments.Equal(state.Segments) {
		if !m.Segments.IsNull() && !m.Segments.IsUnknown() {
			var segs []string
			diagnostics.Append(m.Segments.ElementsAs(ctx, &segs, false)...)
			ops = append(ops, client.NewReplacePatch("/segments", segs))
		} else {
			ops = append(ops, client.NewRemovePatch("/segments"))
		}
	}

	return ops, diagnostics
}

// boolPtrToTF converts *bool to types.Bool (nil → null).
func boolPtrToTF(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}
