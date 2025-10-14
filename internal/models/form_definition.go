package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type FormDefinitionModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Created     types.String `tfsdk:"created"`
	Modified    types.String `tfsdk:"modified"`

	Owner *FormOwner `tfsdk:"owner"`
}

func (m *FormDefinitionModel) FromApiModel(apiModel map[string]interface{}) {
	if apiModel == nil {
		return
	}

	SetStringValue(apiModel, "id", &m.Id)
	SetStringValue(apiModel, "name", &m.Name)
	SetStringValue(apiModel, "description", &m.Description)
	SetStringValue(apiModel, "created", &m.Created)
	SetStringValue(apiModel, "modified", &m.Modified)

	if owner, ok := apiModel["owner"].(map[string]interface{}); ok {
		m.Owner = &FormOwner{}
		m.Owner.FromApiModel(owner)
	}
}

func (m *FormDefinitionModel) ToCreateApiModel() map[string]interface{} {
	apiModel := map[string]interface{}{
		"name": m.Name.ValueString(),
	}

	// Only include description if it's not null
	if !m.Description.IsNull() {
		apiModel["description"] = m.Description.ValueString()
	}

	if m.Owner != nil {
		apiModel["owner"] = map[string]interface{}{
			"type": m.Owner.Type.ValueString(),
			"id":   m.Owner.Id.ValueString(),
		}
		// Only include name if it's not null
		if !m.Owner.Name.IsNull() {
			apiModel["owner"].(map[string]interface{})["name"] = m.Owner.Name.ValueString()
		}
	}

	return apiModel
}
