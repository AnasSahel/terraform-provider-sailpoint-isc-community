// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"

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
		m.ConnectorAttributes, diags = common.MarshalJSONOrDefault(api.ConnectorAttributes, "{}")
		diagnostics.Append(diags...)
	} else {
		m.ConnectorAttributes = jsontypes.NewNormalizedNull()
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
