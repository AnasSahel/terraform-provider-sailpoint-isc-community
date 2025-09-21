// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_datasource

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
)

// ConnectorsDataSourceModel represents the Terraform data source model for listing connectors.
type ConnectorsDataSourceModel struct {
	ID           types.String            `tfsdk:"id"`
	Filters      types.String            `tfsdk:"filters"`
	Limit        types.Int32             `tfsdk:"limit"`
	Offset       types.Int32             `tfsdk:"offset"`
	IncludeCount types.Bool              `tfsdk:"include_count"`
	Locale       types.String            `tfsdk:"locale"`
	PaginateAll  types.Bool              `tfsdk:"paginate_all"`
	MaxResults   types.Int32             `tfsdk:"max_results"`
	PageSize     types.Int32             `tfsdk:"page_size"`
	Connectors   []ConnectorSummaryModel `tfsdk:"connectors"`
}

// ConnectorSummaryModel represents a connector in the list response.
type ConnectorSummaryModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	ScriptName    types.String `tfsdk:"script_name"`
	ClassName     types.String `tfsdk:"class_name"`
	DirectConnect types.Bool   `tfsdk:"direct_connect"`
	Status        types.String `tfsdk:"status"`
	Features      types.List   `tfsdk:"features"`
}

// FromSailPointV3ConnectorDto updates the model from V3ConnectorDto.
func (m *ConnectorSummaryModel) FromSailPointV3ConnectorDto(ctx context.Context, connector *api_v2025.V3ConnectorDto) error {
	// Generate ID from script name since V3ConnectorDto doesn't have ID field
	if connector.HasScriptName() {
		m.ID = types.StringValue(connector.GetScriptName())
		m.ScriptName = types.StringValue(connector.GetScriptName())
	} else {
		// Fallback ID generation from name if no script name
		if connector.HasName() {
			hash := fmt.Sprintf("%x", md5.Sum([]byte(connector.GetName())))
			m.ID = types.StringValue(hash[:8])
		}
	}

	m.Name = utils.StringOrNull(connector.HasName(), connector.GetName())
	m.Type = utils.StringOrNull(connector.HasType(), connector.GetType())
	m.ClassName = utils.StringOrNull(connector.HasClassName(), connector.GetClassName())
	m.DirectConnect = utils.BoolOrNull(connector.HasDirectConnect(), connector.GetDirectConnect())
	m.Status = utils.StringOrNull(connector.HasStatus(), connector.GetStatus())

	// Handle features list if available
	if connector.HasFeatures() {
		features := connector.GetFeatures()
		featureValues := make([]attr.Value, len(features))
		for i, feature := range features {
			featureValues[i] = types.StringValue(feature)
		}
		featuresListValue, diags := types.ListValue(types.StringType, featureValues)
		if diags.HasError() {
			return fmt.Errorf("error creating features list: %s", diags.Errors())
		}
		m.Features = featuresListValue
	} else {
		m.Features = types.ListNull(types.StringType)
	}

	return nil
}
