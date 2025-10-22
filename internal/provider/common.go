package provider

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ObjectRefDataSourceSchema() map[string]datasource_schema.Attribute {
	return map[string]datasource_schema.Attribute{
		"type": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         "The type of the referenced object.",
			MarkdownDescription: "The type of the referenced object.",
		},
		"id": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         "The unique identifier of the referenced object.",
			MarkdownDescription: "The unique identifier (UUID) of the referenced object.",
		},
		"name": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         "The name of the referenced object.",
			MarkdownDescription: "The human-readable name of the referenced object.",
		},
	}
}

func ObjectRefResourceSchema() map[string]resource_schema.Attribute {
	return map[string]resource_schema.Attribute{
		"type": resource_schema.StringAttribute{
			Description:         "The type of the referenced object.",
			MarkdownDescription: "The type of the referenced object.",
			Required:            true,
		},
		"id": resource_schema.StringAttribute{
			Description:         "The unique identifier of the referenced object.",
			MarkdownDescription: "The unique identifier (UUID) of the referenced object.",
			Required:            true,
		},
		"name": resource_schema.StringAttribute{
			Description:         "The name of the referenced object.",
			MarkdownDescription: "The human-readable name of the referenced object.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}
