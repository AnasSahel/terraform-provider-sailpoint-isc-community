// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"fmt"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common/ignorejson"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common/jsonpath"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sourceModel represents the Terraform state for a SailPoint source.
type sourceModel struct {
	ID                        types.String           `tfsdk:"id"`
	Name                      types.String           `tfsdk:"name"`
	Description               types.String           `tfsdk:"description"`
	Owner                     *common.ObjectRefModel `tfsdk:"owner"`
	Cluster                   *common.ObjectRefModel `tfsdk:"cluster"`
	Connector                 types.String           `tfsdk:"connector"`
	ConnectorClass            types.String           `tfsdk:"connector_class"`
	ConnectorAttributes       jsontypes.Normalized   `tfsdk:"connector_attributes"`
	ConnectorAttributesAll    jsontypes.Normalized   `tfsdk:"connector_attributes_all"`
	IgnoreAttributesPaths     types.List             `tfsdk:"ignore_attributes_paths"`
	IgnoreJSONChanges         types.List             `tfsdk:"ignore_json_changes"`
	ConnectionType            types.String           `tfsdk:"connection_type"`
	Type                      types.String           `tfsdk:"type"`
	DeleteThreshold           types.Int64            `tfsdk:"delete_threshold"`
	Authoritative             types.Bool             `tfsdk:"authoritative"`
	Healthy                   types.Bool             `tfsdk:"healthy"`
	Status                    types.String           `tfsdk:"status"`
	Features                  types.Set              `tfsdk:"features"`
	CredentialProviderEnabled types.Bool             `tfsdk:"credential_provider_enabled"`
	Category                  types.String           `tfsdk:"category"`
	ProvisionAsCsv            types.Bool             `tfsdk:"provision_as_csv"`
	Created                   types.String           `tfsdk:"created"`
	Modified                  types.String           `tfsdk:"modified"`
}

// FromAPI maps fields from the API response to the Terraform model.
func (m *sourceModel) FromAPI(ctx context.Context, api client.SourceAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = common.StringOrNullIfEmpty(api.Description)
	m.Connector = types.StringValue(api.Connector)
	m.ConnectorClass = common.StringOrNullIfEmpty(api.ConnectorClass)
	m.ConnectionType = common.StringOrNullIfEmpty(api.ConnectionType)
	m.Type = common.StringOrNullIfEmpty(api.Type)
	m.Status = common.StringOrNullIfEmpty(api.Status)
	m.Healthy = types.BoolValue(api.Healthy)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)
	m.Category = common.StringOrNull(api.Category)

	// Map delete threshold
	if api.DeleteThreshold != nil {
		m.DeleteThreshold = types.Int64Value(*api.DeleteThreshold)
	} else {
		m.DeleteThreshold = types.Int64Null()
	}

	// Map authoritative
	if api.Authoritative != nil {
		m.Authoritative = types.BoolValue(*api.Authoritative)
	} else {
		m.Authoritative = types.BoolNull()
	}

	// Map credential provider enabled
	if api.CredentialProviderEnabled != nil {
		m.CredentialProviderEnabled = types.BoolValue(*api.CredentialProviderEnabled)
	} else {
		m.CredentialProviderEnabled = types.BoolNull()
	}

	// Map owner
	if api.Owner != nil {
		m.Owner, diags = common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
	} else {
		m.Owner = nil
	}

	// Map cluster
	if api.Cluster != nil {
		m.Cluster, diags = common.NewObjectRefFromAPIPtr(ctx, *api.Cluster)
		diagnostics.Append(diags...)
	} else {
		m.Cluster = nil
	}

	// Map connector attributes.
	// ConnectorAttributesAll always receives the full API response.
	// ConnectorAttributes is left null here; the Read/Create/Update callers
	// are responsible for projecting it down to the user-declared key set,
	// because FromAPI has no access to prior state or config.
	if api.ConnectorAttributes != nil {
		normalized, d := common.MarshalJSONOrDefault(api.ConnectorAttributes, "{}")
		diagnostics.Append(d...)
		m.ConnectorAttributesAll = normalized
		m.ConnectorAttributes = jsontypes.NewNormalizedNull()
	} else {
		m.ConnectorAttributes = jsontypes.NewNormalizedNull()
		m.ConnectorAttributesAll = jsontypes.NewNormalizedNull()
	}

	// Map features
	if len(api.Features) > 0 {
		m.Features, diags = types.SetValueFrom(ctx, types.StringType, api.Features)
		diagnostics.Append(diags...)
	} else {
		m.Features = types.SetNull(types.StringType)
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API create request.
func (m *sourceModel) ToAPI(ctx context.Context) (client.SourceAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics
	var apiRequest client.SourceAPI

	apiRequest.Name = m.Name.ValueString()
	apiRequest.Description = m.Description.ValueString()
	apiRequest.Connector = m.Connector.ValueString()
	apiRequest.ConnectorClass = m.ConnectorClass.ValueString()
	apiRequest.ConnectionType = m.ConnectionType.ValueString()
	apiRequest.Type = m.Type.ValueString()

	// Map delete threshold
	if !m.DeleteThreshold.IsNull() && !m.DeleteThreshold.IsUnknown() {
		v := m.DeleteThreshold.ValueInt64()
		apiRequest.DeleteThreshold = &v
	}

	// Map authoritative
	if !m.Authoritative.IsNull() && !m.Authoritative.IsUnknown() {
		v := m.Authoritative.ValueBool()
		apiRequest.Authoritative = &v
	}

	// Map credential provider enabled
	if !m.CredentialProviderEnabled.IsNull() && !m.CredentialProviderEnabled.IsUnknown() {
		v := m.CredentialProviderEnabled.ValueBool()
		apiRequest.CredentialProviderEnabled = &v
	}

	// Map category
	if !m.Category.IsNull() && !m.Category.IsUnknown() {
		v := m.Category.ValueString()
		apiRequest.Category = &v
	}

	// Map owner
	if m.Owner != nil {
		ownerAPI, d := m.Owner.ToAPI(ctx)
		diagnostics.Append(d...)
		apiRequest.Owner = &ownerAPI
	}

	// Map cluster
	if m.Cluster != nil {
		clusterAPI, d := m.Cluster.ToAPI(ctx)
		diagnostics.Append(d...)
		apiRequest.Cluster = &clusterAPI
	}

	// Map connector attributes
	if connAttrs, d := common.UnmarshalJSONField[map[string]interface{}](m.ConnectorAttributes); connAttrs != nil {
		apiRequest.ConnectorAttributes = *connAttrs
		diagnostics.Append(d...)
	} else {
		diagnostics.Append(d...)
	}

	// Map features
	if !m.Features.IsNull() && !m.Features.IsUnknown() {
		var features []string
		diags = m.Features.ElementsAs(ctx, &features, false)
		diagnostics.Append(diags...)
		apiRequest.Features = features
	}

	return apiRequest, diagnostics
}

// ToPatchOperations compares the plan (m) with the current state and generates JSON Patch operations
// for mutable fields that have changed. Immutable fields (connector, connector_class, type, authoritative)
// use RequiresReplace and are not included here.
func (m *sourceModel) ToPatchOperations(ctx context.Context, state *sourceModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var patchOps []client.JSONPatchOperation

	// Name
	if !m.Name.Equal(state.Name) {
		patchOps = append(patchOps, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	// Description
	if !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
			patchOps = append(patchOps, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			patchOps = append(patchOps, client.NewReplacePatch("/description", ""))
		}
	}

	// Owner
	if !reflect.DeepEqual(m.Owner, state.Owner) {
		if m.Owner != nil {
			ownerAPI, diags := common.NewObjectRefToAPIPtr(ctx, *m.Owner)
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/owner", ownerAPI))
		}
	}

	// Cluster
	if !reflect.DeepEqual(m.Cluster, state.Cluster) {
		if m.Cluster != nil {
			clusterAPI, diags := common.NewObjectRefToAPIPtr(ctx, *m.Cluster)
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/cluster", clusterAPI))
		} else {
			patchOps = append(patchOps, client.NewRemovePatch("/cluster"))
		}
	}

	// Connector Attributes — merge user-managed keys onto the full current
	// attributes so server-managed keys (healthy, status, etc.) are preserved.
	if !m.ConnectorAttributes.Equal(state.ConnectorAttributes) {
		if !m.ConnectorAttributes.IsNull() && !m.ConnectorAttributes.IsUnknown() {
			userAttrs, diags := common.UnmarshalJSONField[map[string]interface{}](m.ConnectorAttributes)
			diagnostics.Append(diags...)
			if userAttrs != nil {
				// Start from the full current attributes, then overlay user keys
				merged := make(map[string]interface{})
				var serverAttrs map[string]interface{}
				if fullAttrs, d := common.UnmarshalJSONField[map[string]interface{}](state.ConnectorAttributesAll); fullAttrs != nil {
					diagnostics.Append(d...)
					serverAttrs = *fullAttrs
					for k, v := range serverAttrs {
						merged[k] = v
					}
				} else {
					diagnostics.Append(d...)
				}
				for k, v := range *userAttrs {
					merged[k] = v
				}

				// Re-inject server values at paths the user wants to ignore, so that
				// nested server-managed fields (e.g. passwords inside domainSettings)
				// are not wiped by the top-level key replacement. Covers both the
				// deprecated ignore_attributes_paths and ignore_json_changes.
				if serverAttrs != nil {
					ignorePaths, d := m.resolvedIgnorePaths(ctx)
					diagnostics.Append(d...)
					if !diagnostics.HasError() && len(ignorePaths) > 0 {
						if err := jsonpath.PreservePaths(merged, serverAttrs, ignorePaths); err != nil {
							diagnostics.AddError(
								"Invalid ignore paths",
								fmt.Sprintf("Failed to apply ignore paths: %s", err),
							)
						}
					}
				}

				patchOps = append(patchOps, client.NewReplacePatch("/connectorAttributes", merged))
			}
		}
	}

	// Delete Threshold
	if !m.DeleteThreshold.Equal(state.DeleteThreshold) {
		if !m.DeleteThreshold.IsNull() && !m.DeleteThreshold.IsUnknown() {
			patchOps = append(patchOps, client.NewReplacePatch("/deleteThreshold", m.DeleteThreshold.ValueInt64()))
		}
	}

	// Features
	if !m.Features.Equal(state.Features) {
		if !m.Features.IsNull() && !m.Features.IsUnknown() {
			var features []string
			diags := m.Features.ElementsAs(ctx, &features, false)
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/features", features))
		}
	}

	// Credential Provider Enabled
	if !m.CredentialProviderEnabled.Equal(state.CredentialProviderEnabled) {
		if !m.CredentialProviderEnabled.IsNull() && !m.CredentialProviderEnabled.IsUnknown() {
			patchOps = append(patchOps, client.NewReplacePatch("/credentialProviderEnabled", m.CredentialProviderEnabled.ValueBool()))
		}
	}

	// Category
	if !m.Category.Equal(state.Category) {
		if !m.Category.IsNull() {
			patchOps = append(patchOps, client.NewReplacePatch("/category", m.Category.ValueString()))
		} else {
			patchOps = append(patchOps, client.NewRemovePatch("/category"))
		}
	}

	return patchOps, diagnostics
}

// resolvedIgnorePaths returns the union of ignore paths as jsonpath remainders
// rooted at "$", drawn from both the deprecated ignore_attributes_paths (already
// "$"-rooted, e.g. "$.domainSettings[*].password") and ignore_json_changes
// (rooted at the resource field, e.g. "connector_attributes.domainSettings[*].password",
// whose remainder is "$.domainSettings[*].password").
//
// These feed jsonpath.RemovePaths (Read / ModifyPlan, to prune the ignored leaf
// from the projected connector_attributes so it never produces a diff) and
// jsonpath.PreservePaths (apply, to keep the server value out of the PATCH wipe).
func (m *sourceModel) resolvedIgnorePaths(ctx context.Context) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var paths []string

	if !m.IgnoreAttributesPaths.IsNull() && !m.IgnoreAttributesPaths.IsUnknown() {
		var legacy []string
		diags.Append(m.IgnoreAttributesPaths.ElementsAs(ctx, &legacy, false)...)
		paths = append(paths, legacy...)
	}

	remainders, d := ignorejson.Remainders(ctx, m.IgnoreJSONChanges, "connector_attributes")
	diags.Append(d...)
	paths = append(paths, remainders...)

	return paths, diags
}

// pruneIgnoredPaths removes every ignored path from attrs, returning the pruned
// value. It is a no-op when attrs is null/unknown or there are no paths. Used by
// Read and ModifyPlan so a masked/server-managed nested field the practitioner
// chose to ignore is absent from the managed connector_attributes — matching the
// config (which omits it) and reaching No changes.
func pruneIgnoredPaths(attrs jsontypes.Normalized, paths []string) (jsontypes.Normalized, diag.Diagnostics) {
	var diags diag.Diagnostics
	if attrs.IsNull() || attrs.IsUnknown() || len(paths) == 0 {
		return attrs, diags
	}
	pruned, err := jsonpath.RemovePathsInJSON(attrs.ValueString(), paths)
	if err != nil {
		diags.AddError(
			"Invalid ignore paths",
			fmt.Sprintf("Failed to prune ignored paths from connector_attributes: %s", err),
		)
		return attrs, diags
	}
	return jsontypes.NewNormalizedValue(pruned), diags
}

// projectConnectorAttributes returns a jsontypes.Normalized containing only
// the top-level keys present in managedKeys, with values taken from fullAttrs.
//
// Purpose: on every Read, the API returns the full connector-attribute set
// (including server-managed keys such as "status", "healthy", "since",
// credential blobs, etc.). To honour the documented partial-management
// contract — "only the keys you specify are managed by Terraform" — we must
// project the API response down to the keys the user declared before writing
// to state. This suppresses phantom drift from server-added keys while still
// surfacing real drift on declared keys (e.g., a value the server normalised).
//
// If a key that is present in managedKeys is absent from fullAttrs (the server
// deleted it), the managed value is kept as-is so that the plan shows drift.
func projectConnectorAttributes(managedKeys, fullAttrs jsontypes.Normalized) (jsontypes.Normalized, diag.Diagnostics) {
	managed, diags := common.UnmarshalJSONField[map[string]interface{}](managedKeys)
	if diags.HasError() || managed == nil {
		return managedKeys, diags
	}
	full, d := common.UnmarshalJSONField[map[string]interface{}](fullAttrs)
	diags.Append(d...)
	if diags.HasError() || full == nil {
		// If fullAttrs is unparseable, fall back to managed keys unchanged
		// so we do not silently swallow the user's config.
		return managedKeys, diags
	}
	projected := make(map[string]interface{}, len(*managed))
	for k := range *managed {
		if v, ok := (*full)[k]; ok {
			// Server has this key: use the current server value so real
			// drift (server-side normalisation) is visible in the plan.
			projected[k] = v
		} else {
			// Server no longer has this key: keep the managed value so
			// the plan surfaces the divergence.
			projected[k] = (*managed)[k]
		}
	}
	return common.MarshalJSONOrDefault(projected, "{}")
}
