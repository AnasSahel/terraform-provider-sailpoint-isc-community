// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// GetTransformsDataSourceSchema returns the schema definition for the transforms data source.
func GetTransformsDataSourceSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SailPoint ISC Transforms data source with optional filtering",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this data source.",
				Computed:            true,
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "Filter results using the standard syntax. Supported filters: name (sw, co), type (eq), internal (eq).",
				Optional:            true,
			},
			"transforms": schema.ListNestedAttribute{
				MarkdownDescription: "List of transforms",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Transform identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Transform name",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Transform type",
							Computed:            true,
						},
						"internal": schema.BoolAttribute{
							MarkdownDescription: "Indicates if the transform is internal",
							Computed:            true,
						},
						"attributes": schema.StringAttribute{
							MarkdownDescription: "Transform attributes as JSON",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// GetTransformDataSourceSchema returns the schema definition for the single transform data source.
func GetTransformDataSourceSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SailPoint ISC Transform data source for retrieving a single transform by ID or name",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Transform identifier. Either 'id' or 'name' must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Transform name. Either 'id' or 'name' must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Transform type",
				Computed:            true,
			},
			"internal": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the transform is internal",
				Computed:            true,
			},
			"attributes": schema.StringAttribute{
				MarkdownDescription: "Transform attributes as JSON",
				Computed:            true,
			},
		},
	}
}
