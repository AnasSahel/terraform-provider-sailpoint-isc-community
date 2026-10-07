// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// entitlementModel represents the Terraform state for an Entitlement resource.
type entitlementModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	Attribute              types.String `tfsdk:"attribute"`
	Value                  types.String `tfsdk:"value"`
	SourceSchemaObjectType types.String `tfsdk:"source_schema_object_type"`
	Privileged             types.Bool   `tfsdk:"privileged"`
	CloudGoverned          types.Bool   `tfsdk:"cloud_governed"`
	Requestable            types.Bool   `tfsdk:"requestable"`
	Owner                  types.Object `tfsdk:"owner"`
	Source                 types.Object `tfsdk:"source"`
	Segments               types.Set    `tfsdk:"segments"`
	ManuallyUpdatedFields  types.Map    `tfsdk:"manually_updated_fields"`
	Created                types.String `tfsdk:"created"`
	Modified               types.String `tfsdk:"modified"`
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

	owner, diags := objectRefFromAPI(ctx, api.Owner)
	diagnostics.Append(diags...)
	m.Owner = owner

	source, diags := objectRefFromAPI(ctx, api.Source)
	diagnostics.Append(diags...)
	m.Source = source

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
//
// Unknown plan values mean "not set in config, let the server decide" and never produce a patch.
// A remove patch is only emitted when the plan is explicitly null and the state still holds a value.
func (m *entitlementModel) ToPatchOperations(ctx context.Context, state *entitlementModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var ops []client.JSONPatchOperation

	if !m.Name.Equal(state.Name) && !m.Name.IsNull() && !m.Name.IsUnknown() {
		ops = append(ops, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	if !m.Description.IsUnknown() && !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
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

	if !m.Owner.IsUnknown() {
		switch {
		case m.Owner.IsNull():
			if !state.Owner.IsNull() && !state.Owner.IsUnknown() {
				ops = append(ops, client.NewRemovePatch("/owner"))
			}
		default:
			var planOwner common.ObjectRefModel
			diagnostics.Append(m.Owner.As(ctx, &planOwner, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
			if ownerChanged(ctx, planOwner, state.Owner, &diagnostics) {
				ownerAPI, diags := common.NewObjectRefToAPIPtr(ctx, planOwner)
				diagnostics.Append(diags...)
				ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
			}
		}
	}

	if !m.Segments.IsUnknown() && !m.Segments.Equal(state.Segments) {
		if !m.Segments.IsNull() {
			var segs []string
			diagnostics.Append(m.Segments.ElementsAs(ctx, &segs, false)...)
			ops = append(ops, client.NewReplacePatch("/segments", segs))
		} else {
			ops = append(ops, client.NewRemovePatch("/segments"))
		}
	}

	return ops, diagnostics
}

// ownerChanged reports whether the planned owner differs from the owner in state.
// Only type and id are compared: name is server-resolved and may be unknown in the plan.
func ownerChanged(ctx context.Context, plan common.ObjectRefModel, state types.Object, diagnostics *diag.Diagnostics) bool {
	if state.IsNull() || state.IsUnknown() {
		return true
	}
	var stateOwner common.ObjectRefModel
	diagnostics.Append(state.As(ctx, &stateOwner, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
	return !plan.Type.Equal(stateOwner.Type) || !plan.ID.Equal(stateOwner.ID)
}

// objectRefFromAPI converts an optional API object reference to a types.Object (nil → null).
func objectRefFromAPI(ctx context.Context, api *client.ObjectRefAPI) (types.Object, diag.Diagnostics) {
	if api == nil {
		return types.ObjectNull(common.ObjectRefObjectType.AttrTypes), nil
	}
	ref, diags := common.NewObjectRefFromAPI(ctx, *api)
	if diags.HasError() {
		return types.ObjectNull(common.ObjectRefObjectType.AttrTypes), diags
	}
	obj, objDiags := types.ObjectValueFrom(ctx, common.ObjectRefObjectType.AttrTypes, ref)
	diags.Append(objDiags...)
	return obj, diags
}

// boolPtrToTF converts *bool to types.Bool (nil → null).
func boolPtrToTF(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}
