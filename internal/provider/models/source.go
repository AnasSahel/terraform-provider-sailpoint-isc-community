package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Source struct {
	ID                        types.String                     `tfsdk:"id"`
	Name                      types.String                     `tfsdk:"name"`
	Description               types.String                     `tfsdk:"description"`
	Owner                     *ObjectRef                       `tfsdk:"owner"`
	Cluster                   *ObjectRef                       `tfsdk:"cluster"`
	AccountCorrelationConfig  *ObjectRef                       `tfsdk:"account_correlation_config"`
	AccountCorrelationRule    *ObjectRef                       `tfsdk:"account_correlation_rule"`
	ManagerCorrelationMapping *SourceManagerCorrelationMapping `tfsdk:"manager_correlation_mapping"`
	ManagerCorrelationRule    *ObjectRef                       `tfsdk:"manager_correlation_rule"`
	BeforeProvisioningRule    *ObjectRef                       `tfsdk:"before_provisioning_rule"`
	Schemas                   []ObjectRef                      `tfsdk:"schemas"`
	PasswordPolicies          []ObjectRef                      `tfsdk:"password_policies"`
	Features                  types.List                       `tfsdk:"features"`
	Type                      types.String                     `tfsdk:"type"`
	Connector                 types.String                     `tfsdk:"connector"`
	ConnectorClass            types.String                     `tfsdk:"connector_class"`
	ConnectorAttributes       jsontypes.Normalized             `tfsdk:"connector_attributes"`
	DeleteThreshold           types.Int32                      `tfsdk:"delete_threshold"`
	Authoritative             types.Bool                       `tfsdk:"authoritative"`
	ManagementWorkgroup       *ObjectRef                       `tfsdk:"management_workgroup"`
	Healthy                   types.Bool                       `tfsdk:"healthy"`
	Status                    types.String                     `tfsdk:"status"`
	Since                     types.String                     `tfsdk:"since"`
	ConnectorID               types.String                     `tfsdk:"connector_id"`
	ConnectorName             types.String                     `tfsdk:"connector_name"`
	ConnectorType             types.String                     `tfsdk:"connector_type"`
	ConnectorImplementationID types.String                     `tfsdk:"connector_implementation_id"`
	Created                   types.String                     `tfsdk:"created"`
	Modified                  types.String                     `tfsdk:"modified"`
	CredentialProviderEnabled types.Bool                       `tfsdk:"credential_provider_enabled"`
	Category                  types.String                     `tfsdk:"category"`
}

type SourceManagerCorrelationMapping struct {
	AccountAttributeName  types.String `tfsdk:"account_attribute_name"`
	IdentityAttributeName types.String `tfsdk:"identity_attribute_name"`
}
