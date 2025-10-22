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
	Type            types.String `tfsdk:"type"`
	Connector       types.String `tfsdk:"connector"`
	ConnectorClass  types.String `tfsdk:"connector_class"`
	DeleteThreshold types.Int32  `tfsdk:"delete_threshold"`
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
	// Features                  types.List                       `tfsdk:"features"`
	// ConnectorAttributes       jsontypes.Normalized             `tfsdk:"connector_attributes"`
	// Authoritative             types.Bool                       `tfsdk:"authoritative"`
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

func (s *Source) ConvertToSailPoint(_ context.Context) client.Source {
	if s == nil {
		return client.Source{}
	}

	source := client.Source{
		Name:            s.Name.ValueString(),
		Description:     NewGoTypeValueIf[types.String, string](s.Description, !IsTerraformValueNullOrUnknown(s.Description)),
		Owner:           NewObjectRefFromTerraform(s.Owner),
		Cluster:         NewObjectRefFromTerraform(s.Cluster),
		Type:            s.Type.ValueString(),
		Connector:       s.Connector.ValueString(),
		ConnectorClass:  NewGoTypeValueIf[types.String, string](s.ConnectorClass, !IsTerraformValueNullOrUnknown(s.ConnectorClass)),
		DeleteThreshold: NewGoTypeValueIf[types.Int32, int32](s.DeleteThreshold, !IsTerraformValueNullOrUnknown(s.DeleteThreshold)),
	}

	return source
}

func (s *Source) ConvertToCreateRequestPtr(_ context.Context) *client.Source {
	source := s.ConvertToSailPoint(context.Background())
	return &source
}

func (s *Source) ConvertFromSailPoint(_ context.Context, source *client.Source, includeNull bool) {
	if s == nil || source == nil {
		return
	}

	s.ID = types.StringValue(source.ID)
	s.Name = types.StringValue(source.Name)
	s.Description = types.StringValue(source.Description)
	s.Owner = NewObjectRefFromSailPoint(source.Owner)
	s.Cluster = NewObjectRefFromSailPoint(source.Cluster)
	s.Connector = types.StringValue(source.Connector)
	s.Created = types.StringValue(source.Created)
	s.Modified = types.StringValue(source.Modified)

	s.Type = NewTerraformTypeValueIf[types.String](source.Type, includeNull || !IsTerraformValueNullOrUnknown(s.Type))
	s.ConnectorClass = NewTerraformTypeValueIf[types.String](source.ConnectorClass, includeNull || !IsTerraformValueNullOrUnknown(s.ConnectorClass))
	s.DeleteThreshold = NewTerraformTypeValueIf[types.Int32](source.DeleteThreshold, includeNull || !IsTerraformValueNullOrUnknown(s.DeleteThreshold))
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

	return ops
}
