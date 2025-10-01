package identity_attribute

import "github.com/hashicorp/terraform-plugin-framework/types"

type IdentityAttributeModel struct {
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Standard    types.Bool   `tfsdk:"standard"`
	Type        types.String `tfsdk:"type"`
	Multi       types.Bool   `tfsdk:"multi"`
	Searchable  types.Bool   `tfsdk:"searchable"`
	System      types.Bool   `tfsdk:"system"`
	Sources     []Source1    `tfsdk:"sources"`
}

type Source1 struct {
	Type       types.String `tfsdk:"type"`
	Properties types.String `tfsdk:"properties"`
}

type IdentityAttributeResourceModel struct {
	IdentityAttributeModel
}

type IdentityAttributeDataSourceModel struct {
	IdentityAttributeModel
}

type IdentityAttributeDataSourceListModel struct {
	IncludeSystem  types.Bool               `tfsdk:"include_system"`
	IncludeSilent  types.Bool               `tfsdk:"include_silent"`
	SearchableOnly types.Bool               `tfsdk:"searchable_only"`
	Items          []IdentityAttributeModel `tfsdk:"items"`
}
