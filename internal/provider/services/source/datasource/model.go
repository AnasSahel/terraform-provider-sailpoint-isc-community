// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_datasource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// OwnerReference represents the owner of a source
type OwnerReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ClusterReference represents the cluster associated with a source
type ClusterReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ManagementWorkgroupReference represents the management workgroup associated with a source
type ManagementWorkgroupReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// SourceDataSourceModel represents the data source model for sources
type SourceDataSourceModel struct {
	// Core identifiers
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Required attributes (computed in data source)
	Description types.String    `tfsdk:"description"`
	Owner       *OwnerReference `tfsdk:"owner"`
	Connector   types.String    `tfsdk:"connector"`

	// Core attributes
	Type           types.String      `tfsdk:"type"`
	ConnectorClass types.String      `tfsdk:"connector_class"`
	ConnectionType types.String      `tfsdk:"connection_type"`
	Authoritative  types.Bool        `tfsdk:"authoritative"`
	Cluster        *ClusterReference `tfsdk:"cluster"`

	// Configuration attributes
	ConnectorAttributes types.String `tfsdk:"connector_attributes"` // JSON-encoded map
	DeleteThreshold     types.Int64  `tfsdk:"delete_threshold"`
	Features            types.List   `tfsdk:"features"` // List of strings

	// Management attributes
	ManagementWorkgroup *ManagementWorkgroupReference `tfsdk:"management_workgroup"`

	// Correlation & Rules
	AccountCorrelationConfig  types.String `tfsdk:"account_correlation_config"`  // JSON-encoded CorrelationConfigReference
	AccountCorrelationRule    types.String `tfsdk:"account_correlation_rule"`    // JSON-encoded RuleReference
	ManagerCorrelationRule    types.String `tfsdk:"manager_correlation_rule"`    // JSON-encoded RuleReference
	ManagerCorrelationMapping types.String `tfsdk:"manager_correlation_mapping"` // JSON-encoded CorrelationMapping

	// Provisioning
	BeforeProvisioningRule types.String `tfsdk:"before_provisioning_rule"` // JSON-encoded RuleReference
	PasswordPolicies       types.String `tfsdk:"password_policies"`        // JSON-encoded List of PolicyReference

	// Status & Metadata (Computed)
	Healthy       types.Bool   `tfsdk:"healthy"`
	Status        types.String `tfsdk:"status"`
	Since         types.String `tfsdk:"since"`
	Created       types.String `tfsdk:"created"`
	Modified      types.String `tfsdk:"modified"`
	ConnectorId   types.String `tfsdk:"connector_id"`
	ConnectorName types.String `tfsdk:"connector_name"`
	Schemas       types.String `tfsdk:"schemas"` // JSON-encoded List of SchemaReference

	// Special Parameters
	CredentialProviderEnabled types.Bool   `tfsdk:"credential_provider_enabled"`
	Category                  types.String `tfsdk:"category"`
}

// FromSailPointSource populates the data source model from a SailPoint API Source object.
func (r *SourceDataSourceModel) FromSailPointSource(ctx context.Context, apiModel *api_v2025.Source) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiModel == nil {
		diags.AddError(
			"Invalid API Response",
			"Received nil source from SailPoint API",
		)
		return diags
	}

	// Set core identifiers
	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Description = types.StringValue(apiModel.GetDescription())
	r.Connector = types.StringValue(apiModel.GetConnector())

	// Handle owner
	owner := apiModel.GetOwner()
	r.Owner = &OwnerReference{
		Type: types.StringValue(owner.GetType()),
		Id:   types.StringValue(owner.GetId()),
		Name: types.StringValue(owner.GetName()),
	}

	// Map all other fields similar to resource model
	if apiModel.GetType() != "" {
		r.Type = types.StringValue(apiModel.GetType())
	}

	if apiModel.GetConnectorClass() != "" {
		r.ConnectorClass = types.StringValue(apiModel.GetConnectorClass())
	}

	if apiModel.GetConnectionType() != "" {
		r.ConnectionType = types.StringValue(apiModel.GetConnectionType())
	}

	if apiModel.Authoritative != nil {
		r.Authoritative = types.BoolValue(*apiModel.Authoritative)
	}

	// Handle cluster - convert to JSON
	if apiModel.Cluster.IsSet() {
		cluster := apiModel.Cluster.Get()
		if cluster != nil {
			r.Cluster = &ClusterReference{
				Type: types.StringValue(cluster.GetType()),
				Id:   types.StringValue(cluster.GetId()),
				Name: types.StringValue(cluster.GetName()),
			}
		}
	}

	// Handle management workgroup
	if apiModel.ManagementWorkgroup.IsSet() {
		workgroup := apiModel.ManagementWorkgroup.Get()
		if workgroup != nil {
			r.ManagementWorkgroup = &ManagementWorkgroupReference{
				Type: types.StringValue(workgroup.GetType()),
				Id:   types.StringValue(workgroup.GetId()),
				Name: types.StringValue(workgroup.GetName()),
			}
		}
	}

	// Handle connector attributes
	if len(apiModel.GetConnectorAttributes()) > 0 {
		connectorAttrs := apiModel.GetConnectorAttributes()
		attrsJson, err := json.Marshal(connectorAttrs)
		if err != nil {
			diags.AddError(
				"JSON Encoding Error",
				"Failed to encode connector attributes: "+err.Error(),
			)
		} else {
			r.ConnectorAttributes = types.StringValue(string(attrsJson))
		}
	}

	// Handle all computed fields
	if apiModel.Healthy != nil {
		r.Healthy = types.BoolValue(*apiModel.Healthy)
	}

	if apiModel.GetStatus() != "" {
		r.Status = types.StringValue(apiModel.GetStatus())
	}

	if apiModel.Created != nil {
		r.Created = types.StringValue(apiModel.Created.String())
	}

	if apiModel.Modified != nil {
		r.Modified = types.StringValue(apiModel.Modified.String())
	}

	if apiModel.GetConnectorName() != "" {
		r.ConnectorName = types.StringValue(apiModel.GetConnectorName())
	}

	return diags
}

// SourcesDataSourceModel represents the data source model for listing sources
type SourcesDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	// Filter parameters
	Filters      types.String `tfsdk:"filters"`
	Sorters      types.String `tfsdk:"sorters"`
	Limit        types.Int32  `tfsdk:"limit"`
	Offset       types.Int32  `tfsdk:"offset"`
	IncludeCount types.Bool   `tfsdk:"include_count"`

	// Pagination parameters
	PaginateAll types.Bool  `tfsdk:"paginate_all"`
	MaxResults  types.Int32 `tfsdk:"max_results"`
	PageSize    types.Int32 `tfsdk:"page_size"`

	// Results
	Sources []SourceSummaryModel `tfsdk:"sources"`
}

// SourceSummaryModel represents a summary of a source for list operations
type SourceSummaryModel struct {
	Id                        types.String      `tfsdk:"id"`
	Name                      types.String      `tfsdk:"name"`
	Description               types.String      `tfsdk:"description"`
	Owner                     *OwnerReference   `tfsdk:"owner"`
	Connector                 types.String      `tfsdk:"connector"`
	Type                      types.String      `tfsdk:"type"`
	ConnectorClass            types.String      `tfsdk:"connector_class"`
	ConnectionType            types.String      `tfsdk:"connection_type"`
	Authoritative             types.Bool        `tfsdk:"authoritative"`
	Cluster                   *ClusterReference `tfsdk:"cluster"`
	ConnectorAttributes       types.String      `tfsdk:"connector_attributes"`
	DeleteThreshold           types.Int64       `tfsdk:"delete_threshold"`
	Features                  types.List        `tfsdk:"features"`
	Healthy                   types.Bool        `tfsdk:"healthy"`
	Status                    types.String      `tfsdk:"status"`
	Since                     types.String      `tfsdk:"since"`
	Created                   types.String      `tfsdk:"created"`
	Modified                  types.String      `tfsdk:"modified"`
	ConnectorId               types.String      `tfsdk:"connector_id"`
	ConnectorName             types.String      `tfsdk:"connector_name"`
	CredentialProviderEnabled types.Bool        `tfsdk:"credential_provider_enabled"`
	Category                  types.String      `tfsdk:"category"`
}

// FromSailPointSource populates the summary model from a SailPoint API Source object.
func (r *SourceSummaryModel) FromSailPointSource(ctx context.Context, apiModel *api_v2025.Source) error {
	if apiModel == nil {
		return fmt.Errorf("received nil source from SailPoint API")
	}

	// Map required fields
	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Connector = types.StringValue(apiModel.GetConnector())

	if apiModel.GetDescription() != "" {
		r.Description = types.StringValue(apiModel.GetDescription())
	}

	// Handle owner
	owner := apiModel.GetOwner()
	r.Owner = &OwnerReference{
		Type: types.StringValue(owner.GetType()),
		Id:   types.StringValue(owner.GetId()),
		Name: types.StringValue(owner.GetName()),
	}

	// Map optional fields
	if apiModel.GetType() != "" {
		r.Type = types.StringValue(apiModel.GetType())
	}

	if apiModel.GetConnectorClass() != "" {
		r.ConnectorClass = types.StringValue(apiModel.GetConnectorClass())
	}

	if apiModel.GetConnectionType() != "" {
		r.ConnectionType = types.StringValue(apiModel.GetConnectionType())
	}

	if apiModel.Authoritative != nil {
		r.Authoritative = types.BoolValue(*apiModel.Authoritative)
	}

	// Handle cluster
	if apiModel.Cluster.IsSet() {
		cluster := apiModel.Cluster.Get()
		if cluster != nil {
			r.Cluster = &ClusterReference{
				Type: types.StringValue(cluster.GetType()),
				Id:   types.StringValue(cluster.GetId()),
				Name: types.StringValue(cluster.GetName()),
			}
		}
	}

	// Handle connector attributes
	if len(apiModel.GetConnectorAttributes()) > 0 {
		connectorAttrs := apiModel.GetConnectorAttributes()
		attrsJson, err := json.Marshal(connectorAttrs)
		if err != nil {
			return fmt.Errorf("failed to encode connector attributes: %v", err)
		}
		r.ConnectorAttributes = types.StringValue(string(attrsJson))
	}

	// Handle delete threshold
	if apiModel.DeleteThreshold != nil {
		r.DeleteThreshold = types.Int64Value(int64(*apiModel.DeleteThreshold))
	}

	// Handle features
	if len(apiModel.GetFeatures()) > 0 {
		featuresValue, featuresDiags := types.ListValueFrom(ctx, types.StringType, apiModel.GetFeatures())
		if featuresDiags.HasError() {
			return fmt.Errorf("failed to convert features")
		}
		r.Features = featuresValue
	}

	// Handle computed fields
	if apiModel.Healthy != nil {
		r.Healthy = types.BoolValue(*apiModel.Healthy)
	}

	if apiModel.GetStatus() != "" {
		r.Status = types.StringValue(apiModel.GetStatus())
	}

	if apiModel.GetSince() != "" {
		r.Since = types.StringValue(apiModel.GetSince())
	}

	if apiModel.Created != nil {
		r.Created = types.StringValue(apiModel.Created.String())
	}

	if apiModel.Modified != nil {
		r.Modified = types.StringValue(apiModel.Modified.String())
	}

	if apiModel.GetConnectorId() != "" {
		r.ConnectorId = types.StringValue(apiModel.GetConnectorId())
	}

	if apiModel.GetConnectorName() != "" {
		r.ConnectorName = types.StringValue(apiModel.GetConnectorName())
	}

	// Handle credential provider enabled
	if apiModel.CredentialProviderEnabled != nil {
		r.CredentialProviderEnabled = types.BoolValue(*apiModel.CredentialProviderEnabled)
	}

	// Handle category
	if apiModel.Category.IsSet() {
		categoryPtr := apiModel.Category.Get()
		if categoryPtr != nil {
			r.Category = types.StringValue(*categoryPtr)
		}
	}

	return nil
}
