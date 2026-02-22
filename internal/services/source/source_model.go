// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
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
	ConnectionType            types.String           `tfsdk:"connection_type"`
	Type                      types.String           `tfsdk:"type"`
	DeleteThreshold           types.Int64            `tfsdk:"delete_threshold"`
	Authoritative             types.Bool             `tfsdk:"authoritative"`
	Healthy                   types.Bool             `tfsdk:"healthy"`
	Status                    types.String           `tfsdk:"status"`
	Features                  types.List             `tfsdk:"features"`
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

	// Map connector attributes
	if api.ConnectorAttributes != nil {
		normalized, d := common.MarshalJSONOrDefault(api.ConnectorAttributes, "{}")
		diagnostics.Append(d...)
		m.ConnectorAttributes = normalized
		m.ConnectorAttributesAll = normalized
	} else {
		m.ConnectorAttributes = jsontypes.NewNormalizedNull()
		m.ConnectorAttributesAll = jsontypes.NewNormalizedNull()
	}

	// Map features
	if len(api.Features) > 0 {
		m.Features, diags = types.ListValueFrom(ctx, types.StringType, api.Features)
		diagnostics.Append(diags...)
	} else {
		m.Features = types.ListNull(types.StringType)
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

	// Connector Attributes
	if !m.ConnectorAttributes.Equal(state.ConnectorAttributes) {
		if !m.ConnectorAttributes.IsNull() && !m.ConnectorAttributes.IsUnknown() {
			if connAttrs, diags := common.UnmarshalJSONField[map[string]interface{}](m.ConnectorAttributes); connAttrs != nil {
				diagnostics.Append(diags...)
				patchOps = append(patchOps, client.NewReplacePatch("/connectorAttributes", *connAttrs))
			} else {
				diagnostics.Append(diags...)
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
