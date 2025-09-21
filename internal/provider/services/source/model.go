// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// SourceModel represents the core data structure for a SailPoint source
// This model is shared between resource and data source implementations.
type SourceModel struct {
	// Core identifiers
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Required attributes
	Description types.String `tfsdk:"description"`
	Owner       types.String `tfsdk:"owner"` // JSON-encoded OwnerReference
	Connector   types.String `tfsdk:"connector"`

	// Core attributes
	Type           types.String `tfsdk:"type"`
	ConnectorClass types.String `tfsdk:"connector_class"`
	ConnectionType types.String `tfsdk:"connection_type"`
	Authoritative  types.Bool   `tfsdk:"authoritative"`
	Cluster        types.String `tfsdk:"cluster"` // JSON-encoded ClusterReference

	// Configuration attributes
	ConnectorAttributes types.String `tfsdk:"connector_attributes"` // JSON-encoded map
	DeleteThreshold     types.Int64  `tfsdk:"delete_threshold"`
	Features            types.List   `tfsdk:"features"` // List of strings

	// Management attributes
	ManagementWorkgroup types.String `tfsdk:"management_workgroup"` // JSON-encoded WorkgroupReference

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

// SourceResourceModel extends the base model for resource-specific operations.
type SourceResourceModel struct {
	SourceModel
}

// ToSailPointCreateSourceRequest converts the Terraform model to a SailPoint API create request.
func (r SourceResourceModel) ToSailPointCreateSourceRequest(ctx context.Context) (*api_v2025.Source, diag.Diagnostics) {
	var diags diag.Diagnostics
	source := api_v2025.NewSourceWithDefaults()

	// Validate and set required fields
	if r.Name.IsNull() || r.Name.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The name field is required for creating a source.",
		)
		return nil, diags
	}
	source.SetName(r.Name.ValueString())

	if r.Description.IsNull() || r.Description.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The description field is required for creating a source.",
		)
		return nil, diags
	}
	source.SetDescription(r.Description.ValueString())

	// Handle owner (required) - parse JSON
	if r.Owner.IsNull() || r.Owner.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The owner field is required for creating a source.",
		)
		return nil, diags
	}

	var ownerData map[string]interface{}
	err := json.Unmarshal([]byte(r.Owner.ValueString()), &ownerData)
	if err != nil {
		diags.AddError(
			"Invalid Owner JSON",
			"Failed to parse owner JSON: "+err.Error(),
		)
		return nil, diags
	}

	ownerRef := api_v2025.NewSourceOwnerWithDefaults()
	if id, ok := ownerData["id"].(string); ok {
		ownerRef.SetId(id)
	}
	if ownerType, ok := ownerData["type"].(string); ok {
		ownerRef.SetType(ownerType)
	}
	if name, ok := ownerData["name"].(string); ok {
		ownerRef.SetName(name)
	}
	source.SetOwner(*ownerRef)

	// Handle connector (required)
	if r.Connector.IsNull() || r.Connector.ValueString() == "" {
		diags.AddError(
			"Missing Required Field",
			"The connector field is required for creating a source.",
		)
		return nil, diags
	}
	source.SetConnector(r.Connector.ValueString())

	// Handle optional fields
	if !r.Type.IsNull() && r.Type.ValueString() != "" {
		source.SetType(r.Type.ValueString())
	}

	if !r.ConnectorClass.IsNull() && r.ConnectorClass.ValueString() != "" {
		source.SetConnectorClass(r.ConnectorClass.ValueString())
	}

	if !r.ConnectionType.IsNull() && r.ConnectionType.ValueString() != "" {
		source.SetConnectionType(r.ConnectionType.ValueString())
	}

	if !r.Authoritative.IsNull() {
		source.SetAuthoritative(r.Authoritative.ValueBool())
	}

	// Handle cluster reference (JSON)
	if !r.Cluster.IsNull() && r.Cluster.ValueString() != "" {
		var clusterData map[string]interface{}
		err := json.Unmarshal([]byte(r.Cluster.ValueString()), &clusterData)
		if err != nil {
			diags.AddError(
				"Invalid Cluster JSON",
				"Failed to parse cluster JSON: "+err.Error(),
			)
			return nil, diags
		}

		clusterRef := api_v2025.NewSourceClusterWithDefaults()
		if id, ok := clusterData["id"].(string); ok {
			clusterRef.SetId(id)
		}
		if name, ok := clusterData["name"].(string); ok {
			clusterRef.SetName(name)
		}
		if clusterType, ok := clusterData["type"].(string); ok {
			clusterRef.SetType(clusterType)
		}
		source.SetCluster(*clusterRef)
	}

	// Handle connector attributes (JSON)
	if !r.ConnectorAttributes.IsNull() && r.ConnectorAttributes.ValueString() != "" {
		var connectorAttrs map[string]interface{}
		err := json.Unmarshal([]byte(r.ConnectorAttributes.ValueString()), &connectorAttrs)
		if err != nil {
			diags.AddError(
				"Invalid Connector Attributes",
				"Failed to parse connector_attributes JSON: "+err.Error(),
			)
			return nil, diags
		}
		source.SetConnectorAttributes(connectorAttrs)
	}

	// Handle delete threshold
	if !r.DeleteThreshold.IsNull() {
		source.SetDeleteThreshold(int32(r.DeleteThreshold.ValueInt64()))
	}

	// Handle features list
	if !r.Features.IsNull() {
		var features []string
		featuresDiags := r.Features.ElementsAs(ctx, &features, false)
		if featuresDiags.HasError() {
			diags.Append(featuresDiags...)
			return nil, diags
		}
		source.SetFeatures(features)
	}

	// Handle credential provider enabled
	if !r.CredentialProviderEnabled.IsNull() {
		source.SetCredentialProviderEnabled(r.CredentialProviderEnabled.ValueBool())
	}

	return source, diags
}

// FromSailPointSource populates the Terraform model from a SailPoint API response.
func (r *SourceResourceModel) FromSailPointSource(ctx context.Context, apiModel *api_v2025.Source) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiModel == nil {
		diags.AddError(
			"Invalid API Response",
			"Received nil source from SailPoint API",
		)
		return diags
	}

	// Map required fields
	r.Id = types.StringValue(apiModel.GetId())
	r.Name = types.StringValue(apiModel.GetName())
	r.Connector = types.StringValue(apiModel.GetConnector())

	if apiModel.GetDescription() != "" {
		r.Description = types.StringValue(apiModel.GetDescription())
	}

	// Handle owner - convert to JSON
	owner := apiModel.GetOwner()
	ownerData := map[string]interface{}{
		"id":   owner.GetId(),
		"type": owner.GetType(),
		"name": owner.GetName(),
	}
	ownerJson, err := json.Marshal(ownerData)
	if err != nil {
		diags.AddError(
			"JSON Encoding Error",
			"Failed to encode owner: "+err.Error(),
		)
	} else {
		r.Owner = types.StringValue(string(ownerJson))
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

	// Handle cluster - convert to JSON
	if apiModel.Cluster.IsSet() {
		cluster := apiModel.Cluster.Get()
		if cluster != nil {
			clusterData := map[string]interface{}{
				"id":   cluster.GetId(),
				"name": cluster.GetName(),
				"type": cluster.GetType(),
			}
			clusterJson, err := json.Marshal(clusterData)
			if err != nil {
				diags.AddError(
					"JSON Encoding Error",
					"Failed to encode cluster: "+err.Error(),
				)
			} else {
				r.Cluster = types.StringValue(string(clusterJson))
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

	// Handle delete threshold
	if apiModel.DeleteThreshold != nil {
		r.DeleteThreshold = types.Int64Value(int64(*apiModel.DeleteThreshold))
	}

	// Handle features
	if len(apiModel.GetFeatures()) > 0 {
		featuresValue, featuresDiags := types.ListValueFrom(ctx, types.StringType, apiModel.GetFeatures())
		if featuresDiags.HasError() {
			diags.Append(featuresDiags...)
		} else {
			r.Features = featuresValue
		}
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

	// Handle schemas - convert to JSON
	if len(apiModel.GetSchemas()) > 0 {
		schemas := apiModel.GetSchemas()
		var schemasList []map[string]interface{}
		for _, schema := range schemas {
			schemaData := map[string]interface{}{
				"id":   schema.GetId(),
				"name": schema.GetName(),
				"type": schema.GetType(),
			}
			schemasList = append(schemasList, schemaData)
		}
		schemasJson, err := json.Marshal(schemasList)
		if err != nil {
			diags.AddError(
				"JSON Encoding Error",
				"Failed to encode schemas: "+err.Error(),
			)
		} else {
			r.Schemas = types.StringValue(string(schemasJson))
		}
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

	return diags
}

// SourceDataSourceModel represents the data source model for sources
type SourceDataSourceModel struct {
	SourceModel
}

// FromSailPointSource populates the data source model from a SailPoint API Source object.
func (r *SourceDataSourceModel) FromSailPointSource(ctx context.Context, apiModel *api_v2025.Source) diag.Diagnostics {
	return r.FromSailPointSourceDataSource(ctx, apiModel)
}

// FromSailPointSourceDataSource populates the data source model from a SailPoint API Source object.
func (r *SourceDataSourceModel) FromSailPointSourceDataSource(ctx context.Context, apiModel *api_v2025.Source) diag.Diagnostics {
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

	// Handle owner - convert to JSON
	owner := apiModel.GetOwner()
	ownerData := map[string]interface{}{
		"id":   owner.GetId(),
		"type": owner.GetType(),
		"name": owner.GetName(),
	}
	ownerJson, err := json.Marshal(ownerData)
	if err != nil {
		diags.AddError(
			"JSON Encoding Error",
			"Failed to encode owner: "+err.Error(),
		)
	} else {
		r.Owner = types.StringValue(string(ownerJson))
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
			clusterData := map[string]interface{}{
				"id":   cluster.GetId(),
				"name": cluster.GetName(),
				"type": cluster.GetType(),
			}
			clusterJson, err := json.Marshal(clusterData)
			if err != nil {
				diags.AddError(
					"JSON Encoding Error",
					"Failed to encode cluster: "+err.Error(),
				)
			} else {
				r.Cluster = types.StringValue(string(clusterJson))
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
