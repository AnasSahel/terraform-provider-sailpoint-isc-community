// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iancoleman/strcase"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// ManagedClusterModel represents the core data structure for a SailPoint managed cluster
// This model is shared between resource and data source implementations.
type ManagedClusterModel struct {
	// Core identifiers
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Required attributes
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`

	// Optional configuration
	Configuration types.Map `tfsdk:"configuration"`

	// Computed organizational attributes
	Pod types.String `tfsdk:"pod"`
	Org types.String `tfsdk:"org"`

	// Computed cluster information
	ClientType   types.String `tfsdk:"client_type"`
	CcgVersion   types.String `tfsdk:"ccg_version"`
	PinnedConfig types.Bool   `tfsdk:"pinned_config"`
	Operational  types.Bool   `tfsdk:"operational"`
	Status       types.String `tfsdk:"status"`
	AlertKey     types.String `tfsdk:"alert_key"`

	// Computed security attributes
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKey            types.String `tfsdk:"public_key"`

	// Computed metrics and relationships
	ClientIds    types.List   `tfsdk:"client_ids"`
	ServiceCount types.Int32  `tfsdk:"service_count"`
	CcId         types.String `tfsdk:"cc_id"`

	// Computed timestamps
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// ManagedClusterResourceModel extends the base model for resource-specific operations.
type ManagedClusterResourceModel struct {
	ManagedClusterModel
}

// ToSailPointCreateManagedClusterRequest converts the Terraform model to a SailPoint API create request.
func (r ManagedClusterResourceModel) ToSailPointCreateManagedClusterRequest(ctx context.Context) (*api_v2025.ManagedClusterRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	managedClusterRequest := api_v2025.NewManagedClusterRequestWithDefaults()

	// Validate and set required fields
	if r.Name.IsNull() || r.Name.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The name field is required for creating a managed cluster.",
		)
		return nil, diags
	}
	managedClusterRequest.SetName(r.Name.ValueString())

	if r.Type.IsNull() || r.Type.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The type field is required for creating a managed cluster.",
		)
		return nil, diags
	}
	managedClusterRequest.SetType(api_v2025.ManagedClusterTypes(r.Type.ValueString()))

	if r.Description.IsNull() || r.Description.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The description field is required for creating a managed cluster.",
		)
		return nil, diags
	}
	managedClusterRequest.SetDescription(r.Description.ValueString())

	// Handle optional configuration
	if !r.Configuration.IsNull() {
		terraformConfig := make(map[string]string)
		configDiags := r.Configuration.ElementsAs(ctx, &terraformConfig, false)
		if configDiags.HasError() {
			diags.Append(configDiags...)
			return nil, diags
		}

		// Convert Terraform field names (snake_case) to SailPoint API field names (camelCase)
		sailpointConfig := make(map[string]string)
		for k, v := range terraformConfig {
			camelKey := strcase.ToLowerCamel(k)
			sailpointConfig[camelKey] = v
		}

		managedClusterRequest.SetConfiguration(sailpointConfig)
	}

	return managedClusterRequest, diags
}

// FromSailPointManagedCluster populates the Terraform model from a SailPoint API response.
func (r *ManagedClusterResourceModel) FromSailPointManagedCluster(ctx context.Context, apiModel *api_v2025.ManagedCluster) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiModel == nil {
		diags.AddError(
			"Invalid API Response",
			"Received nil managed cluster from SailPoint API",
		)
		return diags
	}

	// Convert configuration from camelCase to snake_case
	sailpointConfig := apiModel.GetConfiguration()
	terraformConfig := make(map[string]string)

	for k, v := range sailpointConfig {
		snakeKey := strcase.ToSnake(k)
		terraformConfig[snakeKey] = v
	}

	// Convert configuration to Terraform Map type
	conf, configDiags := types.MapValueFrom(ctx, types.StringType, terraformConfig)
	if configDiags.HasError() {
		diags.Append(configDiags...)
	}

	// Convert client IDs to Terraform List type
	clientIds, clientIdsDiags := types.ListValueFrom(ctx, types.StringType, apiModel.GetClientIds())
	if clientIdsDiags.HasError() {
		diags.Append(clientIdsDiags...)
	}

	// Map required fields
	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Type = types.StringValue(string(apiModel.GetType()))
	r.Description = types.StringValue(apiModel.GetDescription())
	r.Configuration = conf

	// Map computed organizational fields
	r.Pod = types.StringValue(apiModel.GetPod())
	r.Org = types.StringValue(apiModel.GetOrg())

	// Map computed cluster information fields
	r.ClientType = types.StringValue(string(apiModel.GetClientType()))
	r.CcgVersion = types.StringValue(apiModel.GetCcgVersion())
	r.PinnedConfig = types.BoolValue(apiModel.GetPinnedConfig())
	r.Operational = types.BoolValue(apiModel.GetOperational())
	r.Status = types.StringValue(apiModel.GetStatus())
	r.AlertKey = types.StringValue(apiModel.GetAlertKey())

	// Map computed security fields
	r.PublicKeyCertificate = types.StringValue(apiModel.GetPublicKeyCertificate())
	r.PublicKeyThumbprint = types.StringValue(apiModel.GetPublicKeyThumbprint())
	r.PublicKey = types.StringValue(apiModel.GetPublicKey())

	// Map computed metrics and relationships
	r.ClientIds = clientIds
	r.ServiceCount = types.Int32Value(apiModel.GetServiceCount())
	r.CcId = types.StringValue(apiModel.GetCcId())

	// Map computed timestamps
	r.CreatedAt = types.StringValue(apiModel.GetCreatedAt().String())
	r.UpdatedAt = types.StringValue(apiModel.GetUpdatedAt().String())

	return diags
}

// UpdateSelectiveFields updates only the fields that were changed from an API response
// This method helps prevent inconsistency errors by preserving unchanged computed fields.
func (r *ManagedClusterResourceModel) UpdateSelectiveFields(ctx context.Context, apiModel *api_v2025.ManagedCluster, plan *ManagedClusterResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiModel == nil {
		diags.AddError(
			"Invalid API Response",
			"Received nil managed cluster from SailPoint API during update",
		)
		return diags
	}

	// Always update the ID (required field)
	r.Id = types.StringValue(apiModel.GetId())

	// Update fields that were changed and returned by the API
	if apiModel.HasName() {
		r.Name = types.StringValue(apiModel.GetName())
	}

	if apiModel.HasDescription() {
		r.Description = types.StringValue(apiModel.GetDescription())
	}

	if apiModel.HasType() {
		r.Type = types.StringValue(string(apiModel.GetType()))
	}

	// For configuration, preserve the planned configuration to avoid inconsistency errors
	if !plan.Configuration.IsNull() {
		r.Configuration = plan.Configuration
	}

	// Update computed fields only if they have meaningful values in the response
	if apiModel.HasPod() && apiModel.GetPod() != "" {
		r.Pod = types.StringValue(apiModel.GetPod())
	}

	if apiModel.HasOrg() && apiModel.GetOrg() != "" {
		r.Org = types.StringValue(apiModel.GetOrg())
	}

	// ClientType is always present but may be nullable - check if it's valid
	if clientType, ok := apiModel.GetClientTypeOk(); ok && clientType != nil {
		r.ClientType = types.StringValue(string(*clientType))
	}

	// CcgVersion is required but check if it has meaningful value
	if apiModel.GetCcgVersion() != "" && apiModel.GetCcgVersion() != "Undefined" {
		r.CcgVersion = types.StringValue(apiModel.GetCcgVersion())
	}

	if apiModel.HasUpdatedAt() {
		r.UpdatedAt = types.StringValue(apiModel.GetUpdatedAt().String())
	}

	return diags
}

// ManagedClusterDataSourceModel represents the data source model for managed clusters
// It's identical to ManagedClusterModel but used for data sources specifically.
type ManagedClusterDataSourceModel struct {
	// Core identifiers
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Required attributes
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`

	// Configuration
	Configuration types.Map `tfsdk:"configuration"`

	// Organizational attributes
	Pod types.String `tfsdk:"pod"`
	Org types.String `tfsdk:"org"`

	// Cluster information
	ClientType   types.String `tfsdk:"client_type"`
	CcgVersion   types.String `tfsdk:"ccg_version"`
	PinnedConfig types.Bool   `tfsdk:"pinned_config"`
	Operational  types.Bool   `tfsdk:"operational"`
	Status       types.String `tfsdk:"status"`
	AlertKey     types.String `tfsdk:"alert_key"`

	// Metrics and relationships
	ClientIds    types.List   `tfsdk:"client_ids"`
	ServiceCount types.Int32  `tfsdk:"service_count"`
	CcId         types.String `tfsdk:"cc_id"`

	// Security/Key attributes
	PublicKeyCertificate types.String `tfsdk:"public_key_certificate"`
	PublicKeyThumbprint  types.String `tfsdk:"public_key_thumbprint"`
	PublicKey            types.String `tfsdk:"public_key"`

	// Timestamps
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// FromSailPointManagedClusterDataSource populates the data source model from a SailPoint API ManagedCluster object.
func (r *ManagedClusterDataSourceModel) FromSailPointManagedClusterDataSource(ctx context.Context, apiModel *api_v2025.ManagedCluster) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set core identifiers
	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())

	// Set type and description
	r.Type = types.StringValue(string(apiModel.GetType()))

	if apiModel.HasDescription() {
		r.Description = types.StringValue(apiModel.GetDescription())
	}

	// Convert configuration map from camelCase to snake_case
	if apiModel.HasConfiguration() {
		configMap := make(map[string]string)
		for key, value := range apiModel.GetConfiguration() {
			snakeKey := strcase.ToSnake(key)
			configMap[snakeKey] = value
		}

		configValue, configDiags := types.MapValueFrom(ctx, types.StringType, configMap)
		diags.Append(configDiags...)
		if !configDiags.HasError() {
			r.Configuration = configValue
		}
	} else {
		// Create empty map if no configuration
		configValue, configDiags := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		diags.Append(configDiags...)
		if !configDiags.HasError() {
			r.Configuration = configValue
		}
	}

	// Set organizational attributes
	if apiModel.HasPod() {
		r.Pod = types.StringValue(apiModel.GetPod())
	}

	if apiModel.HasOrg() {
		r.Org = types.StringValue(apiModel.GetOrg())
	}

	// Set cluster information
	r.ClientType = types.StringValue(string(apiModel.GetClientType()))

	// CcgVersion is required but check if it has meaningful value
	if apiModel.GetCcgVersion() != "" && apiModel.GetCcgVersion() != "Undefined" {
		r.CcgVersion = types.StringValue(apiModel.GetCcgVersion())
	}

	if apiModel.HasPinnedConfig() {
		r.PinnedConfig = types.BoolValue(apiModel.GetPinnedConfig())
	}

	if apiModel.HasOperational() {
		r.Operational = types.BoolValue(apiModel.GetOperational())
	}

	if apiModel.HasStatus() {
		r.Status = types.StringValue(apiModel.GetStatus())
	}

	if apiModel.HasAlertKey() {
		r.AlertKey = types.StringValue(apiModel.GetAlertKey())
	}

	// Set metrics and relationships
	if apiModel.HasClientIds() {
		clientIdsValue, clientIdsDiags := types.ListValueFrom(ctx, types.StringType, apiModel.GetClientIds())
		diags.Append(clientIdsDiags...)
		if !clientIdsDiags.HasError() {
			r.ClientIds = clientIdsValue
		}
	} else {
		// Create empty list if no client IDs
		clientIdsValue, clientIdsDiags := types.ListValueFrom(ctx, types.StringType, []string{})
		diags.Append(clientIdsDiags...)
		if !clientIdsDiags.HasError() {
			r.ClientIds = clientIdsValue
		}
	}

	if apiModel.HasServiceCount() {
		r.ServiceCount = types.Int32Value(apiModel.GetServiceCount())
	}

	if apiModel.HasCcId() {
		r.CcId = types.StringValue(apiModel.GetCcId())
	}

	// Set security/key attributes
	if apiModel.HasPublicKeyCertificate() {
		r.PublicKeyCertificate = types.StringValue(apiModel.GetPublicKeyCertificate())
	}

	if apiModel.HasPublicKeyThumbprint() {
		r.PublicKeyThumbprint = types.StringValue(apiModel.GetPublicKeyThumbprint())
	}

	if apiModel.HasPublicKey() {
		r.PublicKey = types.StringValue(apiModel.GetPublicKey())
	}

	// Set timestamps
	if apiModel.HasCreatedAt() {
		r.CreatedAt = types.StringValue(apiModel.GetCreatedAt().String())
	}

	if apiModel.HasUpdatedAt() {
		r.UpdatedAt = types.StringValue(apiModel.GetUpdatedAt().String())
	}

	return diags
}
