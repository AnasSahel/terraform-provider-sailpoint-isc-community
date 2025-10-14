package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type FormOwner struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (o *FormOwner) ToApiModel() map[string]interface{} {
	if o == nil {
		return nil
	}
	return map[string]interface{}{
		"type": o.Type.ValueString(),
		"id":   o.Id.ValueString(),
		"name": o.Name.ValueString(),
	}
}

func (o *FormOwner) FromApiModel(apiModel map[string]interface{}) {
	if apiModel == nil {
		return
	}
	SetStringValue(apiModel, "type", &o.Type)
	SetStringValue(apiModel, "id", &o.Id)
	SetStringValue(apiModel, "name", &o.Name)
}
