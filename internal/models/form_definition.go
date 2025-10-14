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

type FormOwner struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
