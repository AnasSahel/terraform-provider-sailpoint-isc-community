// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_resource

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// AccountCorrelationConfigReference represents the account correlation configuration associated with a source
type AccountCorrelationConfigReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// AccountCorrelationRuleReference represents the account correlation rule associated with a source
type AccountCorrelationRuleReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ManagerCorrelationRuleReference represents the manager correlation rule associated with a source
type ManagerCorrelationRuleReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// BeforeProvisioningRuleReference represents the before provisioning rule associated with a source
type BeforeProvisioningRuleReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ManagerCorrelationMappingReference represents the manager correlation mapping associated with a source
type ManagerCorrelationMappingReference struct {
	AccountAttributeName  types.String `tfsdk:"account_attribute_name"`
	IdentityAttributeName types.String `tfsdk:"identity_attribute_name"`
}

// PasswordPolicyReference represents a password policy associated with a source
type PasswordPolicyReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// SchemaReference represents a schema associated with a source
type SchemaReference struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// SourceResourceModel extends the base model for resource-specific operations.
type SourceResourceModel struct {
	// Core identifiers
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Required attributes
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
	AccountCorrelationConfig  *AccountCorrelationConfigReference  `tfsdk:"account_correlation_config"`
	AccountCorrelationRule    *AccountCorrelationRuleReference    `tfsdk:"account_correlation_rule"`
	ManagerCorrelationRule    *ManagerCorrelationRuleReference    `tfsdk:"manager_correlation_rule"`
	ManagerCorrelationMapping *ManagerCorrelationMappingReference `tfsdk:"manager_correlation_mapping"`

	// Provisioning
	BeforeProvisioningRule *BeforeProvisioningRuleReference `tfsdk:"before_provisioning_rule"`
	PasswordPolicies       types.List                       `tfsdk:"password_policies"`

	// Status & Metadata (Computed)
	Healthy       types.Bool   `tfsdk:"healthy"`
	Status        types.String `tfsdk:"status"`
	Since         types.String `tfsdk:"since"`
	Created       types.String `tfsdk:"created"`
	Modified      types.String `tfsdk:"modified"`
	ConnectorId   types.String `tfsdk:"connector_id"`
	ConnectorName types.String `tfsdk:"connector_name"`
	Schemas       types.List   `tfsdk:"schemas"`

	// Special Parameters
	CredentialProviderEnabled types.Bool   `tfsdk:"credential_provider_enabled"`
	Category                  types.String `tfsdk:"category"`
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

	// Handle owner (required)
	if r.Owner == nil {
		diags.AddError(
			"Missing Required Field",
			"The owner field is required for creating a source.",
		)
		return nil, diags
	}

	ownerRef := api_v2025.NewSourceOwnerWithDefaults()
	if !r.Owner.Id.IsNull() && r.Owner.Id.ValueString() != "" {
		ownerRef.SetId(r.Owner.Id.ValueString())
	} else {
		diags.AddError(
			"Missing Required Field",
			"Owner ID is required.",
		)
		return nil, diags
	}

	if !r.Owner.Type.IsNull() && r.Owner.Type.ValueString() != "" {
		ownerRef.SetType(r.Owner.Type.ValueString())
	} else {
		diags.AddError(
			"Missing Required Field",
			"Owner type is required.",
		)
		return nil, diags
	}

	if !r.Owner.Name.IsNull() && r.Owner.Name.ValueString() != "" {
		ownerRef.SetName(r.Owner.Name.ValueString())
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

	// Handle cluster reference (optional)
	if r.Cluster != nil {
		clusterRef := api_v2025.NewSourceClusterWithDefaults()

		if !r.Cluster.Id.IsNull() && r.Cluster.Id.ValueString() != "" {
			clusterRef.SetId(r.Cluster.Id.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Cluster ID is required when cluster is specified.",
			)
			return nil, diags
		}

		if !r.Cluster.Type.IsNull() && r.Cluster.Type.ValueString() != "" {
			clusterRef.SetType(r.Cluster.Type.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Cluster type is required when cluster is specified.",
			)
			return nil, diags
		}

		if !r.Cluster.Name.IsNull() && r.Cluster.Name.ValueString() != "" {
			clusterRef.SetName(r.Cluster.Name.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Cluster name is required when cluster is specified.",
			)
			return nil, diags
		}

		source.SetCluster(*clusterRef)
	}

	// Handle management workgroup reference (optional)
	if r.ManagementWorkgroup != nil {
		workgroupRef := api_v2025.NewSourceManagementWorkgroupWithDefaults()

		if !r.ManagementWorkgroup.Id.IsNull() && r.ManagementWorkgroup.Id.ValueString() != "" {
			workgroupRef.SetId(r.ManagementWorkgroup.Id.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Management workgroup ID is required when management workgroup is specified.",
			)
			return nil, diags
		}

		if !r.ManagementWorkgroup.Type.IsNull() && r.ManagementWorkgroup.Type.ValueString() != "" {
			workgroupRef.SetType(r.ManagementWorkgroup.Type.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Management workgroup type is required when management workgroup is specified.",
			)
			return nil, diags
		}

		if !r.ManagementWorkgroup.Name.IsNull() && r.ManagementWorkgroup.Name.ValueString() != "" {
			workgroupRef.SetName(r.ManagementWorkgroup.Name.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Management workgroup name is required when management workgroup is specified.",
			)
			return nil, diags
		}

		source.SetManagementWorkgroup(*workgroupRef)
	}

	// Handle account correlation config reference (optional)
	if r.AccountCorrelationConfig != nil {
		correlationConfigRef := api_v2025.NewSourceAccountCorrelationConfigWithDefaults()

		if !r.AccountCorrelationConfig.Id.IsNull() && r.AccountCorrelationConfig.Id.ValueString() != "" {
			correlationConfigRef.SetId(r.AccountCorrelationConfig.Id.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Account correlation config ID is required when account correlation config is specified.",
			)
			return nil, diags
		}

		if !r.AccountCorrelationConfig.Type.IsNull() && r.AccountCorrelationConfig.Type.ValueString() != "" {
			correlationConfigRef.SetType(r.AccountCorrelationConfig.Type.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Account correlation config type is required when account correlation config is specified.",
			)
			return nil, diags
		}

		if !r.AccountCorrelationConfig.Name.IsNull() && r.AccountCorrelationConfig.Name.ValueString() != "" {
			correlationConfigRef.SetName(r.AccountCorrelationConfig.Name.ValueString())
		} else {
			diags.AddError(
				"Missing Required Field",
				"Account correlation config name is required when account correlation config is specified.",
			)
			return nil, diags
		}

		source.SetAccountCorrelationConfig(*correlationConfigRef)
	}

	// Handle account correlation rule reference (optional)
	if r.AccountCorrelationRule != nil {
		correlationRuleRef := api_v2025.NewSourceAccountCorrelationRuleWithDefaults()

		if !r.AccountCorrelationRule.Id.IsNull() && r.AccountCorrelationRule.Id.ValueString() != "" {
			correlationRuleRef.SetId(r.AccountCorrelationRule.Id.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Account correlation rule ID is required when account correlation rule is specified.")
			return nil, diags
		}

		if !r.AccountCorrelationRule.Type.IsNull() && r.AccountCorrelationRule.Type.ValueString() != "" {
			correlationRuleRef.SetType(r.AccountCorrelationRule.Type.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Account correlation rule type is required when account correlation rule is specified.")
			return nil, diags
		}

		if !r.AccountCorrelationRule.Name.IsNull() && r.AccountCorrelationRule.Name.ValueString() != "" {
			correlationRuleRef.SetName(r.AccountCorrelationRule.Name.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Account correlation rule name is required when account correlation rule is specified.")
			return nil, diags
		}

		source.SetAccountCorrelationRule(*correlationRuleRef)
	}

	// Handle manager correlation rule reference (optional)
	if r.ManagerCorrelationRule != nil {
		managerRuleRef := api_v2025.NewSourceManagerCorrelationRuleWithDefaults()

		if !r.ManagerCorrelationRule.Id.IsNull() && r.ManagerCorrelationRule.Id.ValueString() != "" {
			managerRuleRef.SetId(r.ManagerCorrelationRule.Id.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Manager correlation rule ID is required when manager correlation rule is specified.")
			return nil, diags
		}

		if !r.ManagerCorrelationRule.Type.IsNull() && r.ManagerCorrelationRule.Type.ValueString() != "" {
			managerRuleRef.SetType(r.ManagerCorrelationRule.Type.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Manager correlation rule type is required when manager correlation rule is specified.")
			return nil, diags
		}

		if !r.ManagerCorrelationRule.Name.IsNull() && r.ManagerCorrelationRule.Name.ValueString() != "" {
			managerRuleRef.SetName(r.ManagerCorrelationRule.Name.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Manager correlation rule name is required when manager correlation rule is specified.")
			return nil, diags
		}

		source.SetManagerCorrelationRule(*managerRuleRef)
	}

	// Handle before provisioning rule reference (optional)
	if r.BeforeProvisioningRule != nil {
		beforeRuleRef := api_v2025.NewSourceBeforeProvisioningRuleWithDefaults()

		if !r.BeforeProvisioningRule.Id.IsNull() && r.BeforeProvisioningRule.Id.ValueString() != "" {
			beforeRuleRef.SetId(r.BeforeProvisioningRule.Id.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Before provisioning rule ID is required when before provisioning rule is specified.")
			return nil, diags
		}

		if !r.BeforeProvisioningRule.Type.IsNull() && r.BeforeProvisioningRule.Type.ValueString() != "" {
			beforeRuleRef.SetType(r.BeforeProvisioningRule.Type.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Before provisioning rule type is required when before provisioning rule is specified.")
			return nil, diags
		}

		if !r.BeforeProvisioningRule.Name.IsNull() && r.BeforeProvisioningRule.Name.ValueString() != "" {
			beforeRuleRef.SetName(r.BeforeProvisioningRule.Name.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Before provisioning rule name is required when before provisioning rule is specified.")
			return nil, diags
		}

		source.SetBeforeProvisioningRule(*beforeRuleRef)
	}

	// Handle manager correlation mapping (optional)
	if r.ManagerCorrelationMapping != nil {
		mappingRef := api_v2025.NewSourceManagerCorrelationMappingWithDefaults()

		if !r.ManagerCorrelationMapping.AccountAttributeName.IsNull() && r.ManagerCorrelationMapping.AccountAttributeName.ValueString() != "" {
			mappingRef.SetAccountAttributeName(r.ManagerCorrelationMapping.AccountAttributeName.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Account attribute name is required when manager correlation mapping is specified.")
			return nil, diags
		}

		if !r.ManagerCorrelationMapping.IdentityAttributeName.IsNull() && r.ManagerCorrelationMapping.IdentityAttributeName.ValueString() != "" {
			mappingRef.SetIdentityAttributeName(r.ManagerCorrelationMapping.IdentityAttributeName.ValueString())
		} else {
			diags.AddError("Missing Required Field", "Identity attribute name is required when manager correlation mapping is specified.")
			return nil, diags
		}

		source.SetManagerCorrelationMapping(*mappingRef)
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
	if !r.Features.IsNull() && !r.Features.IsUnknown() {
		var features []string
		featuresDiags := r.Features.ElementsAs(ctx, &features, false)
		if featuresDiags.HasError() {
			diags.Append(featuresDiags...)
			return nil, diags
		}
		source.SetFeatures(features)
	}

	// Handle password policies list
	if !r.PasswordPolicies.IsNull() && !r.PasswordPolicies.IsUnknown() {
		var passwordPolicyRefs []PasswordPolicyReference
		passwordPoliciesDiags := r.PasswordPolicies.ElementsAs(ctx, &passwordPolicyRefs, false)
		if passwordPoliciesDiags.HasError() {
			diags.Append(passwordPoliciesDiags...)
			return nil, diags
		}

		if len(passwordPolicyRefs) > 0 {
			var passwordPolicies []api_v2025.SourcePasswordPoliciesInner
			for _, policy := range passwordPolicyRefs {
				policyRef := api_v2025.NewSourcePasswordPoliciesInnerWithDefaults()

				if !policy.Id.IsNull() && policy.Id.ValueString() != "" {
					policyRef.SetId(policy.Id.ValueString())
				} else {
					diags.AddError("Missing Required Field", "Password policy ID is required when password policies are specified.")
					return nil, diags
				}

				if !policy.Type.IsNull() && policy.Type.ValueString() != "" {
					policyRef.SetType(policy.Type.ValueString())
				} else {
					diags.AddError("Missing Required Field", "Password policy type is required when password policies are specified.")
					return nil, diags
				}

				if !policy.Name.IsNull() && policy.Name.ValueString() != "" {
					policyRef.SetName(policy.Name.ValueString())
				} else {
					diags.AddError("Missing Required Field", "Password policy name is required when password policies are specified.")
					return nil, diags
				}

				passwordPolicies = append(passwordPolicies, *policyRef)
			}
			source.SetPasswordPolicies(passwordPolicies)
		}
	}

	// Handle credential provider enabled
	if !r.CredentialProviderEnabled.IsNull() {
		source.SetCredentialProviderEnabled(r.CredentialProviderEnabled.ValueBool())
	}

	return source, diags
}

// Define object types for list elements
var (
	passwordPolicyObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
			"name": types.StringType,
		},
	}

	schemaObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type": types.StringType,
			"id":   types.StringType,
			"name": types.StringType,
		},
	}
) // FromSailPointSource populates the Terraform model from a SailPoint API response.
func (r *SourceResourceModel) FromSailPointSource(ctx context.Context, apiModel *api_v2025.Source) diag.Diagnostics {
	var diags diag.Diagnostics

	if apiModel == nil {
		diags.AddError(
			"Invalid API Response",
			"Received nil source from SailPoint API",
		)
		return diags
	}

	// Initialize empty lists with proper element types
	r.PasswordPolicies = types.ListNull(passwordPolicyObjectType)
	r.Schemas = types.ListNull(schemaObjectType)

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

	// Handle account correlation config - skip auto-generated ones to prevent inconsistent results
	// Only populate if it was explicitly configured in Terraform
	// Note: This will be nil by default, which is what Terraform expects for optional nested attributes

	// Handle account correlation rule
	if apiModel.AccountCorrelationRule.IsSet() {
		correlationRule := apiModel.AccountCorrelationRule.Get()
		if correlationRule != nil {
			r.AccountCorrelationRule = &AccountCorrelationRuleReference{
				Type: types.StringValue(correlationRule.GetType()),
				Id:   types.StringValue(correlationRule.GetId()),
				Name: types.StringValue(correlationRule.GetName()),
			}
		}
	}

	// Handle manager correlation rule
	if apiModel.ManagerCorrelationRule.IsSet() {
		managerRule := apiModel.ManagerCorrelationRule.Get()
		if managerRule != nil {
			r.ManagerCorrelationRule = &ManagerCorrelationRuleReference{
				Type: types.StringValue(managerRule.GetType()),
				Id:   types.StringValue(managerRule.GetId()),
				Name: types.StringValue(managerRule.GetName()),
			}
		}
	}

	// Handle before provisioning rule
	if apiModel.BeforeProvisioningRule.IsSet() {
		beforeRule := apiModel.BeforeProvisioningRule.Get()
		if beforeRule != nil {
			r.BeforeProvisioningRule = &BeforeProvisioningRuleReference{
				Type: types.StringValue(beforeRule.GetType()),
				Id:   types.StringValue(beforeRule.GetId()),
				Name: types.StringValue(beforeRule.GetName()),
			}
		}
	}

	// Handle manager correlation mapping
	if apiModel.ManagerCorrelationMapping != nil {
		mapping := apiModel.GetManagerCorrelationMapping()
		r.ManagerCorrelationMapping = &ManagerCorrelationMappingReference{
			AccountAttributeName:  types.StringValue(mapping.GetAccountAttributeName()),
			IdentityAttributeName: types.StringValue(mapping.GetIdentityAttributeName()),
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

	// Handle password policies
	if len(apiModel.GetPasswordPolicies()) > 0 {
		passwordPolicies := apiModel.GetPasswordPolicies()
		var passwordPolicyRefs []PasswordPolicyReference
		for _, policy := range passwordPolicies {
			passwordPolicyRefs = append(passwordPolicyRefs, PasswordPolicyReference{
				Type: types.StringValue(policy.GetType()),
				Id:   types.StringValue(policy.GetId()),
				Name: types.StringValue(policy.GetName()),
			})
		}
		passwordPoliciesValue, passwordPoliciesDiags := types.ListValueFrom(ctx, passwordPolicyObjectType, passwordPolicyRefs)
		if passwordPoliciesDiags.HasError() {
			diags.Append(passwordPoliciesDiags...)
		} else {
			r.PasswordPolicies = passwordPoliciesValue
		}
	}

	// Handle schemas
	if len(apiModel.GetSchemas()) > 0 {
		schemas := apiModel.GetSchemas()
		var schemaRefs []SchemaReference
		for _, schema := range schemas {
			schemaRefs = append(schemaRefs, SchemaReference{
				Type: types.StringValue(schema.GetType()),
				Id:   types.StringValue(schema.GetId()),
				Name: types.StringValue(schema.GetName()),
			})
		}
		schemasValue, schemasDiags := types.ListValueFrom(ctx, schemaObjectType, schemaRefs)
		if schemasDiags.HasError() {
			diags.Append(schemasDiags...)
		} else {
			r.Schemas = schemasValue
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
