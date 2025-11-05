package models

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ ModelWithSourceTerraformConversionMethods[client.Source] = &Source{}
	_ ModelWithSailPointConversionMethods[client.Source]       = &Source{}
	_ ModelWithSailPointPatchMethods[Source]                   = &Source{}
)

type SourceManagerCorrelationMapping struct {
	AccountAttributeName  types.String `tfsdk:"account_attribute_name"`
	IdentityAttributeName types.String `tfsdk:"identity_attribute_name"`
}

type Source struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Owner           *ObjectRef   `tfsdk:"owner"`
	Cluster         *ObjectRef   `tfsdk:"cluster"`
	Features        types.Set    `tfsdk:"features"`
	Type            types.String `tfsdk:"type"`
	Connector       types.String `tfsdk:"connector"`
	ConnectorClass  types.String `tfsdk:"connector_class"`
	DeleteThreshold types.Int32  `tfsdk:"delete_threshold"`
	Authoritative   types.Bool   `tfsdk:"authoritative"`
	Created         types.String `tfsdk:"created"`
	Modified        types.String `tfsdk:"modified"`

	// ConnectorName types.String `tfsdk:"connector_name"`

	// AccountCorrelationConfig  *ObjectRef                       `tfsdk:"account_correlation_config"`
	// AccountCorrelationRule    *ObjectRef                       `tfsdk:"account_correlation_rule"`
	// ManagerCorrelationMapping *SourceManagerCorrelationMapping `tfsdk:"manager_correlation_mapping"`
	// ManagerCorrelationRule    *ObjectRef                       `tfsdk:"manager_correlation_rule"`
	// BeforeProvisioningRule    *ObjectRef                       `tfsdk:"before_provisioning_rule"`
	// Schemas                   []ObjectRef                      `tfsdk:"schemas"`
	// PasswordPolicies          []ObjectRef                      `tfsdk:"password_policies"`
	// ConnectorAttributes       jsontypes.Normalized             `tfsdk:"connector_attributes"`
	// ManagementWorkgroup       *ObjectRef                       `tfsdk:"management_workgroup"`
	// Healthy                   types.Bool                       `tfsdk:"healthy"`
	// Status                    types.String                     `tfsdk:"status"`
	// Since                     types.String                     `tfsdk:"since"`
	// ConnectorID               types.String                     `tfsdk:"connector_id"`
	// ConnectorType             types.String                     `tfsdk:"connector_type"`
	// ConnectorImplementationID types.String                     `tfsdk:"connector_implementation_id"`
	// CredentialProviderEnabled types.Bool                       `tfsdk:"credential_provider_enabled"`
	// Category                  types.String                     `tfsdk:"category"`
}

func (s *Source) ConvertToSailPoint(ctx context.Context) client.Source {
	if s == nil {
		return client.Source{}
	}

	source := client.Source{
		Name:      s.Name.ValueString(),
		Type:      s.Type.ValueString(),
		Connector: s.Connector.ValueString(),
		// Authoritative: s.Authoritative.ValueBool(),

		Owner:   NewObjectRefFromTerraform(s.Owner),
		Cluster: NewObjectRefFromTerraform(s.Cluster),

		Description:     NewGoTypeValueIf[types.String, string](ctx, s.Description, !IsTerraformValueNullOrUnknown(s.Description)),
		Features:        NewGoTypeValueIf[types.Set, []string](ctx, s.Features, !IsTerraformValueNullOrUnknown(s.Features)),
		ConnectorClass:  NewGoTypeValueIf[types.String, string](ctx, s.ConnectorClass, !IsTerraformValueNullOrUnknown(s.ConnectorClass)),
		DeleteThreshold: NewGoTypeValueIf[types.Int32, int32](ctx, s.DeleteThreshold, !IsTerraformValueNullOrUnknown(s.DeleteThreshold)),
	}

	return source
}

func (s *Source) ConvertToCreateRequestPtr(_ context.Context) *client.Source {
	source := s.ConvertToSailPoint(context.Background())
	return &source
}

func (s *Source) ConvertFromSailPoint(ctx context.Context, source *client.Source, includeNull bool) {
	if s == nil || source == nil {
		return
	}

	s.ID = types.StringValue(source.ID)
	s.Name = types.StringValue(source.Name)
	s.Description = types.StringValue(source.Description)
	s.Connector = types.StringValue(source.Connector)
	s.Created = types.StringValue(source.Created)
	s.Modified = types.StringValue(source.Modified)

	s.Owner = NewObjectRefFromSailPoint(source.Owner)
	s.Cluster = NewObjectRefFromSailPoint(source.Cluster)

	s.Features = NewTerraformTypeValueIf[types.Set](ctx, source.Features, includeNull || !IsTerraformValueNullOrUnknown(s.Features))
	s.Type = NewTerraformTypeValueIf[types.String](ctx, source.Type, includeNull || !IsTerraformValueNullOrUnknown(s.Type))
	s.ConnectorClass = NewTerraformTypeValueIf[types.String](ctx, source.ConnectorClass, includeNull || !IsTerraformValueNullOrUnknown(s.ConnectorClass))
	s.DeleteThreshold = NewTerraformTypeValueIf[types.Int32](ctx, source.DeleteThreshold, includeNull || !IsTerraformValueNullOrUnknown(s.DeleteThreshold))
	s.Authoritative = NewTerraformTypeValueIf[types.Bool](ctx, source.Authoritative, includeNull || !IsTerraformValueNullOrUnknown(s.Authoritative))
}

func (s *Source) ConvertFromSailPointForResource(ctx context.Context, source *client.Source) {
	s.ConvertFromSailPoint(ctx, source, false)
}

func (s *Source) ConvertFromSailPointForDataSource(ctx context.Context, source *client.Source) {
	s.ConvertFromSailPoint(ctx, source, true)
}

func (s *Source) BuildPatchOptions(ctx context.Context, desired *Source) []client.JSONPatchOperation {
	var ops []client.JSONPatchOperation

	if s == nil || desired == nil {
		return ops
	}

	if s.Name.ValueString() != desired.Name.ValueString() {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/name", desired.Name.ValueString()))
	}

	if s.Description.ValueString() != desired.Description.ValueString() {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/description", desired.Description.ValueString()))
	}

	if !s.Owner.Equals(desired.Owner) {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/owner", desired.Owner.ConvertToSailPoint(ctx)))
	}

	if !s.Cluster.Equals(desired.Cluster) {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/cluster", desired.Cluster.ConvertToSailPoint(ctx)))
	}

	if s.ConnectorClass.ValueString() != desired.ConnectorClass.ValueString() {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/connectorClass", desired.ConnectorClass.ValueString()))
	}

	if s.DeleteThreshold.ValueInt32() != desired.DeleteThreshold.ValueInt32() {
		ops = append(ops, client.NewReplaceJSONPatchOperation("/deleteThreshold", desired.DeleteThreshold.ValueInt32()))
	}

	return ops
}
