package provider

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			Required:            true,
			Description:         "The type of the referenced object.",
			MarkdownDescription: "The type of the referenced object.",
		},
		"id": resource_schema.StringAttribute{
			Required:            true,
			Description:         "The unique identifier of the referenced object.",
			MarkdownDescription: "The unique identifier (UUID) of the referenced object.",
		},
		"name": resource_schema.StringAttribute{
			Optional:            true,
			Description:         "The name of the referenced object.",
			MarkdownDescription: "The human-readable name of the referenced object.",
		},
	}
}
