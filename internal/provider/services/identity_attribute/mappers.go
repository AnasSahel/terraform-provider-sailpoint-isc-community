package identity_attribute

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// DataSource
func MapIdentityAttributeToDataSourceModel(identityAttribute api_v2025.IdentityAttribute) IdentityAttributeModel {
	model := IdentityAttributeModel{
		Name:        types.StringValue(identityAttribute.GetName()),
		DisplayName: types.StringValue(identityAttribute.GetDisplayName()),
		Standard:    types.BoolValue(identityAttribute.GetStandard()),
		Type:        types.StringValue(identityAttribute.GetType()),
		Multi:       types.BoolValue(identityAttribute.GetMulti()),
		Searchable:  types.BoolValue(identityAttribute.GetSearchable()),
		System:      types.BoolValue(identityAttribute.GetSystem()),
	}

	model.Sources = MapSourcesToModel(identityAttribute.GetSources())

	return model
}

func MapSourcesToModel(sources []api_v2025.Source1) []Source1 {
	result := make([]Source1, 0, len(sources))

	for _, source := range sources {
		sourceJson, _ := json.Marshal(source)
		result = append(result, Source1{
			Type:       types.StringValue(source.GetType()),
			Properties: types.StringValue(string(sourceJson)),
		})
	}

	return result
}

// Resource
func MapModelToIdentityAttribute(model IdentityAttributeModel) api_v2025.IdentityAttribute {
	identityAttribute := api_v2025.NewIdentityAttributeWithDefaults()
	identityAttribute.SetName(model.Name.ValueString())
	identityAttribute.SetDisplayName(model.DisplayName.ValueString())
	identityAttribute.SetStandard(model.Standard.ValueBool())
	identityAttribute.SetType(model.Type.ValueString())
	identityAttribute.SetMulti(model.Multi.ValueBool())
	identityAttribute.SetSearchable(model.Searchable.ValueBool())
	identityAttribute.SetSystem(model.System.ValueBool())

	return *identityAttribute
}

func MapIdentityAttributeToResourceModel(identityAttribute api_v2025.IdentityAttribute) IdentityAttributeModel {
	model := IdentityAttributeModel{
		Name:        types.StringValue(identityAttribute.GetName()),
		DisplayName: types.StringValue(identityAttribute.GetDisplayName()),
		Standard:    types.BoolValue(identityAttribute.GetStandard()),
		Type:        types.StringValue(identityAttribute.GetType()),
		Multi:       types.BoolValue(identityAttribute.GetMulti()),
		Searchable:  types.BoolValue(identityAttribute.GetSearchable()),
		System:      types.BoolValue(identityAttribute.GetSystem()),
	}

	return model
}

// Resource mappers
func MapIdentityAttributeToTerraformResource(identityAttribute api_v2025.IdentityAttribute) IdentityAttributeResourceModel {
	model := IdentityAttributeResourceModel{
		IdentityAttributeModel: IdentityAttributeModel{
			Name:        types.StringValue(identityAttribute.GetName()),
			DisplayName: types.StringValue(identityAttribute.GetDisplayName()),
			Standard:    types.BoolValue(identityAttribute.GetStandard()),
			Type:        types.StringValue(identityAttribute.GetType()),
			Multi:       types.BoolValue(identityAttribute.GetMulti()),
			Searchable:  types.BoolValue(identityAttribute.GetSearchable()),
			System:      types.BoolValue(identityAttribute.GetSystem()),
		},
	}

	return model
}

func MapTerraformResourceToIdentityAttribute(model IdentityAttributeResourceModel) api_v2025.IdentityAttribute {
	identityAttribute := api_v2025.NewIdentityAttributeWithDefaults()
	identityAttribute.SetName(model.Name.ValueString())
	identityAttribute.SetDisplayName(model.DisplayName.ValueString())
	identityAttribute.SetStandard(model.Standard.ValueBool())
	identityAttribute.SetType(model.Type.ValueString())
	identityAttribute.SetMulti(model.Multi.ValueBool())
	identityAttribute.SetSearchable(model.Searchable.ValueBool())
	identityAttribute.SetSystem(model.System.ValueBool())

	return *identityAttribute
}
