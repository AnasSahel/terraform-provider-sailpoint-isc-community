package models

import (
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ObjectRef struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func ObjectRefDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Computed:            true,
			Description:         "The type of the referenced object.",
			MarkdownDescription: "The type of the referenced object.",
		},
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "The unique identifier of the referenced object.",
			MarkdownDescription: "The unique identifier (UUID) of the referenced object.",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			Description:         "The name of the referenced object.",
			MarkdownDescription: "The human-readable name of the referenced object.",
		},
	}
}

func NewObjectRefFromAPI(apiRef *client.ObjectRef) *ObjectRef {
	if apiRef == nil {
		return nil
	}
	return &ObjectRef{
		Type: types.StringValue(apiRef.Type),
		ID:   types.StringValue(apiRef.ID),
		Name: types.StringValue(apiRef.Name),
	}
}
