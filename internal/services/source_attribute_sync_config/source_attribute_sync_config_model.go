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
	Source     *common.ObjectRefModel         `tfsdk:"source"`
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
	m.Source = &source

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
// attribute list. Attributes not listed in the model retain their current `enabled` value from
// the remote, preventing accidental disabling of unmanaged attributes.
//
// Returns a diagnostic error if the model lists an attribute name not present on the source.
func (m *sourceAttrSyncConfigModel) ToAPI(
	ctx context.Context,
	remoteAttributes []client.SourceAttributeSyncAttributeAPI,
) (*client.SourceAttributeSyncConfigAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	// Build a map of user-specified names → desired enabled value.
	userMap := make(map[string]bool, len(m.Attributes))
	for _, a := range m.Attributes {
		userMap[a.Name.ValueString()] = a.Enabled.ValueBool()
	}

	// Merge: start from the full remote list, override only user-specified names.
	found := make(map[string]struct{}, len(m.Attributes))
	result := make([]client.SourceAttributeSyncAttributeAPI, 0, len(remoteAttributes))
	for _, remote := range remoteAttributes {
		attr := remote
		if enabled, ok := userMap[remote.Name]; ok {
			attr.Enabled = enabled
			found[remote.Name] = struct{}{}
		}
		result = append(result, attr)
	}

	// Validate: any user-specified name not found on the source is an error.
	for _, a := range m.Attributes {
		name := a.Name.ValueString()
		if _, ok := found[name]; !ok {
			diagnostics.AddError(
				"Invalid attribute name",
				"Attribute \""+name+"\" does not exist in the source's attribute sync config. "+
					"Only attributes present in the source's Create Account definition can be synced. "+
					"Check the source's Create Account provisioning policy and remove or fix the attribute name.",
			)
		}
	}

	if diagnostics.HasError() {
		return nil, diagnostics
	}

	source := client.ObjectRefAPI{}
	if m.Source != nil {
		src, diags := m.Source.ToAPI(ctx)
		diagnostics.Append(diags...)
		source = src
	}

	return &client.SourceAttributeSyncConfigAPI{
		Source:     source,
		Attributes: result,
	}, diagnostics
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
