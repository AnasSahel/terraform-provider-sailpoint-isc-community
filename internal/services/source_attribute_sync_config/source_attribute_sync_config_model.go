// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sourceAttrSyncAttributeModel represents a single attribute entry in the sync config.
type sourceAttrSyncAttributeModel struct {
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	DisplayName types.String `tfsdk:"display_name"`
	Target      types.String `tfsdk:"target"`
}

// sourceAttrSyncConfigModel is the Terraform state for the sailpoint_source_attribute_sync_config resource.
type sourceAttrSyncConfigModel struct {
	ID         types.String                   `tfsdk:"id"`
	SourceID   types.String                   `tfsdk:"source_id"`
	Source     types.Object                   `tfsdk:"source"`
	Attributes []sourceAttrSyncAttributeModel `tfsdk:"attributes"`
}

// FromAPI populates the model from the API response.
//
// When cfgNames is non-empty (resource Read path), only the named attributes are written to
// state — this prevents a perpetual plan diff when managing a subset of attributes.
// When cfgNames is nil or empty (data source path), all attributes are populated.
func (m *sourceAttrSyncConfigModel) FromAPI(ctx context.Context, api *client.SourceAttributeSyncConfigAPI, cfgNames map[string]struct{}) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	// id mirrors source_id for import convenience.
	m.ID = m.SourceID

	source, diags := common.NewObjectRefFromAPI(ctx, api.Source)
	diagnostics.Append(diags...)
	sourceObj, diags := types.ObjectValueFrom(ctx, common.ObjectRefObjectType.AttrTypes, source)
	diagnostics.Append(diags...)
	m.Source = sourceObj

	var attrs []sourceAttrSyncAttributeModel
	for _, a := range api.Attributes {
		if len(cfgNames) > 0 {
			if _, ok := cfgNames[a.Name]; !ok {
				continue
			}
		}
		attrs = append(attrs, sourceAttrSyncAttributeModel{
			Name:        types.StringValue(a.Name),
			Enabled:     types.BoolValue(a.Enabled),
			DisplayName: types.StringValue(a.DisplayName),
			Target:      types.StringValue(a.Target),
		})
	}
	m.Attributes = attrs

	return diagnostics
}

// ToAPI builds the PUT request body by merging the model's enabled flags into the full remote
// config. Attributes not listed in the model retain their current `enabled` value from the
// remote, preventing accidental disabling of unmanaged attributes.
//
// displayNames maps identity attribute names to their current display name (from
// /identity-attributes/v1). The sync config GET can return a stale display name after an
// identity attribute is renamed, and the PUT rejects it, so each entry's display name is
// replaced when a current one is known.
//
// An attribute declared with enabled = false that is absent on the server is skipped: ISC drops
// entries whose target is not mapped by the CREATE provisioning policy when they are disabled,
// so the desired state (not synced) is already reached. Any other declared name absent on the
// server is an error.
func (m *sourceAttrSyncConfigModel) ToAPI(
	remote *client.SourceAttributeSyncConfigAPI,
	displayNames map[string]string,
) (*client.SourceAttributeSyncConfigAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	// Build a map of user-specified names → desired enabled value.
	userMap := make(map[string]bool, len(m.Attributes))
	for _, a := range m.Attributes {
		userMap[a.Name.ValueString()] = a.Enabled.ValueBool()
	}

	// Merge: start from the full remote list, override only user-specified names.
	found := make(map[string]struct{}, len(m.Attributes))
	result := make([]client.SourceAttributeSyncAttributeAPI, 0, len(remote.Attributes))
	for _, r := range remote.Attributes {
		attr := r
		if enabled, ok := userMap[r.Name]; ok {
			attr.Enabled = enabled
			found[r.Name] = struct{}{}
		}
		if dn, ok := displayNames[r.Name]; ok && dn != "" {
			attr.DisplayName = dn
		}
		result = append(result, attr)
	}

	// Validate: a declared name not found on the source is an error, unless it is declared
	// disabled (absent == not synced, which is the desired state).
	for _, a := range m.Attributes {
		name := a.Name.ValueString()
		if _, ok := found[name]; ok || !a.Enabled.ValueBool() {
			continue
		}
		diagnostics.AddError(
			"Invalid attribute name",
			"Attribute \""+name+"\" does not exist in the source's attribute sync config. "+
				"Only attributes present in the source's Create Account definition can be synced. "+
				"Check the source's Create Account provisioning policy and remove or fix the attribute name.",
		)
	}

	if diagnostics.HasError() {
		return nil, diagnostics
	}

	return &client.SourceAttributeSyncConfigAPI{
		Source:     remote.Source,
		Attributes: result,
	}, diagnostics
}

// keepAbsentDisabled adds back to the state the declared attributes that have enabled = false
// and are absent from the server response. ISC removes disabled entries that the CREATE
// provisioning policy does not map, so "absent" means "not synced": the declared state is
// reached and keeping the entry avoids a perpetual diff. Computed fields keep their prior
// value when known, otherwise they are null.
func (m *sourceAttrSyncConfigModel) keepAbsentDisabled(declared []sourceAttrSyncAttributeModel, api *client.SourceAttributeSyncConfigAPI) {
	onServer := make(map[string]struct{}, len(api.Attributes))
	for _, a := range api.Attributes {
		onServer[a.Name] = struct{}{}
	}
	for _, d := range declared {
		if d.Enabled.IsUnknown() || d.Enabled.ValueBool() {
			continue
		}
		if _, ok := onServer[d.Name.ValueString()]; ok {
			continue
		}
		entry := sourceAttrSyncAttributeModel{
			Name:        d.Name,
			Enabled:     types.BoolValue(false),
			DisplayName: types.StringNull(),
			Target:      types.StringNull(),
		}
		if !d.DisplayName.IsUnknown() {
			entry.DisplayName = d.DisplayName
		}
		if !d.Target.IsUnknown() {
			entry.Target = d.Target
		}
		m.Attributes = append(m.Attributes, entry)
	}
}

// configuredAttributeNames returns the set of attribute names currently tracked in Terraform state.
// Used by Read to filter the API response to only the managed subset.
func configuredAttributeNames(attrs []sourceAttrSyncAttributeModel) map[string]struct{} {
	names := make(map[string]struct{}, len(attrs))
	for _, a := range attrs {
		names[a.Name.ValueString()] = struct{}{}
	}
	return names
}
