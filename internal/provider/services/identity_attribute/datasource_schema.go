package identity_attribute

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func datasourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Identity Attribute data source schema",
		MarkdownDescription: "Identity Attribute data source schema",
		Attributes:          datasourceModelSchema(),
	}
}

func datasourceListSchema() schema.Schema {
	return schema.Schema{
		Description:         "Identity Attribute List data source schema",
		MarkdownDescription: "Identity Attribute List data source schema",
		Attributes: map[string]schema.Attribute{
			"include_system": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether to include system attributes in the list.",
				MarkdownDescription: "Whether to include system attributes in the list.",
			},
			"include_silent": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether to include silent attributes in the list.",
				MarkdownDescription: "Whether to include silent attributes in the list.",
			},
			"searchable_only": schema.BoolAttribute{
				Optional:            true,
				Description:         "Whether to include only searchable attributes in the list.",
				MarkdownDescription: "Whether to include only searchable attributes in the list.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of identity attributes.",
				MarkdownDescription: "List of identity attributes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: datasourceModelSchema(),
				},
			},
		},
	}
}

func datasourceModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required:            true,
			Description:         "The name of the identity attribute.",
			MarkdownDescription: "The name of the identity attribute.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			Description:         "The display name of the identity attribute.",
			MarkdownDescription: "The display name of the identity attribute.",
		},
		"standard": schema.BoolAttribute{
			Computed:            true,
			Description:         "Indicates if the identity attribute is a standard attribute.",
			MarkdownDescription: "Indicates if the identity attribute is a standard attribute.",
		},
		"type": schema.StringAttribute{
			Computed:            true,
			Description:         "The data type of the identity attribute.",
			MarkdownDescription: "The data type of the identity attribute.",
		},
		"multi": schema.BoolAttribute{
			Computed:            true,
			Description:         "Indicates if the identity attribute can have multiple values.",
			MarkdownDescription: "Indicates if the identity attribute can have multiple values.",
		},
		"searchable": schema.BoolAttribute{
			Computed:            true,
			Description:         "Indicates if the identity attribute is searchable.",
			MarkdownDescription: "Indicates if the identity attribute is searchable.",
		},
		"system": schema.BoolAttribute{
			Computed:            true,
			Description:         "Indicates if the identity attribute is a system attribute.",
			MarkdownDescription: "Indicates if the identity attribute is a system attribute.",
		},
		"sources": schema.ListNestedAttribute{
			Computed:            true,
			Description:         "The sources of the identity attribute.",
			MarkdownDescription: "The sources of the identity attribute.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "The type of the source.",
						MarkdownDescription: "The type of the source.",
					},
					"properties": schema.StringAttribute{
						Computed:            true,
						Description:         "The properties of the source in JSON format.",
						MarkdownDescription: "The properties of the source in JSON format.",
					},
				},
			},
		},
	}
}
