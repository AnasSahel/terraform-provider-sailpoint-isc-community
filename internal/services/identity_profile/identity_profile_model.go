// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity_profile

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// identityExceptionReportRefModel represents the identity exception report reference in Terraform state.
type identityExceptionReportRefModel struct {
	TaskResultID types.String `tfsdk:"task_result_id"`
	ReportName   types.String `tfsdk:"report_name"`
}

func (m *identityExceptionReportRefModel) FromAPI(_ context.Context, api client.IdentityExceptionReportRefAPI) diag.Diagnostics {
	m.TaskResultID = types.StringValue(api.TaskResultID)
	m.ReportName = types.StringValue(api.ReportName)
	return nil
}

// transformDefinitionModel represents a transform definition in Terraform state.
type transformDefinitionModel struct {
	Type       types.String         `tfsdk:"type"`
	Attributes jsontypes.Normalized `tfsdk:"attributes"`
}

func (m *transformDefinitionModel) FromAPI(_ context.Context, api client.TransformDefinitionAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.Type = types.StringValue(api.Type)

	if api.Attributes != nil {
		var diags diag.Diagnostics
		m.Attributes, diags = common.MarshalJSONOrDefault(api.Attributes, "{}")
		diagnostics.Append(diags...)
	} else {
		m.Attributes = jsontypes.NewNormalizedNull()
	}

	return diagnostics
}

func (m *transformDefinitionModel) ToAPI(_ context.Context) (client.TransformDefinitionAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.TransformDefinitionAPI{
		Type: m.Type.ValueString(),
	}

	if attrs, diags := common.UnmarshalJSONField[map[string]interface{}](m.Attributes); attrs != nil {
		api.Attributes = *attrs
		diagnostics.Append(diags...)
	}

	return api, diagnostics
}

// identityAttributeTransformModel represents a transform for an identity attribute.
type identityAttributeTransformModel struct {
	IdentityAttributeName types.String             `tfsdk:"identity_attribute_name"`
	TransformDefinition   transformDefinitionModel `tfsdk:"transform_definition"`
}

func NewIdentityAttributeTransformFromAPI(ctx context.Context, api client.IdentityAttributeTransformAPI) (identityAttributeTransformModel, diag.Diagnostics) {
	var m identityAttributeTransformModel

	diags := m.FromAPI(ctx, api)

	return m, diags
}

func (m *identityAttributeTransformModel) FromAPI(ctx context.Context, api client.IdentityAttributeTransformAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.IdentityAttributeName = types.StringValue(api.IdentityAttributeName)
	diagnostics.Append(m.TransformDefinition.FromAPI(ctx, api.TransformDefinition)...)

	return diagnostics
}

func (m *identityAttributeTransformModel) ToAPI(ctx context.Context) (client.IdentityAttributeTransformAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	transformDef, diags := m.TransformDefinition.ToAPI(ctx)
	diagnostics.Append(diags...)

	return client.IdentityAttributeTransformAPI{
		IdentityAttributeName: m.IdentityAttributeName.ValueString(),
		TransformDefinition:   transformDef,
	}, diagnostics
}

// identityAttributeConfigModel represents the identity attribute configuration in Terraform state.
type identityAttributeConfigModel struct {
	Enabled             types.Bool                        `tfsdk:"enabled"`
	AttributeTransforms []identityAttributeTransformModel `tfsdk:"attribute_transforms"`
}

func (m *identityAttributeConfigModel) FromAPI(ctx context.Context, api client.IdentityAttributeConfigAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.Enabled = types.BoolValue(api.Enabled)

	if len(api.AttributeTransforms) > 0 {
		var diags diag.Diagnostics
		m.AttributeTransforms, diags = common.MapSliceFromAPI(ctx, api.AttributeTransforms, NewIdentityAttributeTransformFromAPI)
		diagnostics.Append(diags...)
	} else {
		m.AttributeTransforms = []identityAttributeTransformModel{}
	}

	return diagnostics
}

func (m *identityAttributeConfigModel) ToAPI(ctx context.Context) (client.IdentityAttributeConfigAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := client.IdentityAttributeConfigAPI{
		Enabled: m.Enabled.ValueBool(),
	}

	transforms := make([]client.IdentityAttributeTransformAPI, len(m.AttributeTransforms))
	for i := range m.AttributeTransforms {
		var diags diag.Diagnostics
		transforms[i], diags = m.AttributeTransforms[i].ToAPI(ctx)
		diagnostics.Append(diags...)
	}
	api.AttributeTransforms = transforms

	return api, diagnostics
}

// identityProfileModel represents the Terraform state for an Identity Profile.
type identityProfileModel struct {
	ID                               types.String                     `tfsdk:"id"`
	Name                             types.String                     `tfsdk:"name"`
	Description                      types.String                     `tfsdk:"description"`
	Owner                            *common.ObjectRefModel           `tfsdk:"owner"`
	Priority                         types.Int64                      `tfsdk:"priority"`
	AuthoritativeSource              common.ObjectRefModel            `tfsdk:"authoritative_source"`
	IdentityRefreshRequired          types.Bool                       `tfsdk:"identity_refresh_required"`
	IdentityCount                    types.Int32                      `tfsdk:"identity_count"`
	IdentityAttributeConfig          *identityAttributeConfigModel    `tfsdk:"identity_attribute_config"`
	IdentityExceptionReportReference *identityExceptionReportRefModel `tfsdk:"identity_exception_report_reference"`
	HasTimeBasedAttr                 types.Bool                       `tfsdk:"has_time_based_attr"`
	Created                          types.String                     `tfsdk:"created"`
	Modified                         types.String                     `tfsdk:"modified"`
}

// FromAPI maps fields from the API response to the Terraform model.
func (m *identityProfileModel) FromAPI(ctx context.Context, api client.IdentityProfileAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Created = types.StringValue(api.Created)
	m.Modified = types.StringValue(api.Modified)
	m.Priority = types.Int64Value(api.Priority)
	m.IdentityRefreshRequired = types.BoolValue(api.IdentityRefreshRequired)
	m.IdentityCount = types.Int32Value(api.IdentityCount)
	m.HasTimeBasedAttr = types.BoolValue(api.HasTimeBasedAttr)

	// Map Description (nullable)
	m.Description = common.StringOrNull(api.Description)

	// Map Owner (nullable)
	if api.Owner != nil {
		var diags diag.Diagnostics
		m.Owner, diags = common.NewObjectRefFromAPIPtr(ctx, *api.Owner)
		diagnostics.Append(diags...)
	} else {
		m.Owner = nil
	}

	// Map AuthoritativeSource
	diagnostics.Append(m.AuthoritativeSource.FromAPI(ctx, api.AuthoritativeSource)...)

	// Map IdentityAttributeConfig
	m.IdentityAttributeConfig = &identityAttributeConfigModel{}
	diagnostics.Append(m.IdentityAttributeConfig.FromAPI(ctx, api.IdentityAttributeConfig)...)

	// Map IdentityExceptionReportReference (nullable)
	if api.IdentityExceptionReportReference != nil {
		m.IdentityExceptionReportReference = &identityExceptionReportRefModel{}
		diagnostics.Append(m.IdentityExceptionReportReference.FromAPI(ctx, *api.IdentityExceptionReportReference)...)
	} else {
		m.IdentityExceptionReportReference = nil
	}

	return diagnostics
}

// ToAPI maps fields from the Terraform model to the API create request.
func (m *identityProfileModel) ToAPI(ctx context.Context) (client.IdentityProfileCreateAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var diags diag.Diagnostics

	apiRequest := client.IdentityProfileCreateAPI{
		Name: m.Name.ValueString(),
	}

	// Map Description (optional, nullable)
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		desc := m.Description.ValueString()
		apiRequest.Description = &desc
	}

	// Map Owner (optional, nullable)
	if m.Owner != nil {
		ownerAPI, d := m.Owner.ToAPI(ctx)
		diagnostics.Append(d...)
		apiRequest.Owner = &ownerAPI
	}

	// Map Priority (optional)
	if !m.Priority.IsNull() && !m.Priority.IsUnknown() {
		apiRequest.Priority = m.Priority.ValueInt64()
	}

	// Map AuthoritativeSource
	apiRequest.AuthoritativeSource, diags = m.AuthoritativeSource.ToAPI(ctx)
	diagnostics.Append(diags...)

	// Map IdentityAttributeConfig (optional)
	if m.IdentityAttributeConfig != nil {
		apiRequest.IdentityAttributeConfig, diags = m.IdentityAttributeConfig.ToAPI(ctx)
		diagnostics.Append(diags...)
	}

	return apiRequest, diagnostics
}

// ToPatchOperations compares the plan (m) with the current state and generates JSON Patch operations
// for fields that have changed.
func (m *identityProfileModel) ToPatchOperations(ctx context.Context, state *identityProfileModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
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
			patchOps = append(patchOps, client.NewRemovePatch("/description"))
		}
	}

	// Priority
	if !m.Priority.Equal(state.Priority) {
		patchOps = append(patchOps, client.NewReplacePatch("/priority", m.Priority.ValueInt64()))
	}

	// Owner
	if !reflect.DeepEqual(m.Owner, state.Owner) {
		if m.Owner != nil {
			ownerAPI, diags := m.Owner.ToAPI(ctx)
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/owner", ownerAPI))
		} else {
			patchOps = append(patchOps, client.NewRemovePatch("/owner"))
		}
	}

	// IdentityAttributeConfig
	if !reflect.DeepEqual(m.IdentityAttributeConfig, state.IdentityAttributeConfig) {
		if m.IdentityAttributeConfig != nil {
			configAPI, diags := m.IdentityAttributeConfig.ToAPI(ctx)
			diagnostics.Append(diags...)
			patchOps = append(patchOps, client.NewReplacePatch("/identityAttributeConfig", configAPI))
		} else {
			patchOps = append(patchOps, client.NewRemovePatch("/identityAttributeConfig"))
		}
	}

	return patchOps, diagnostics
}
