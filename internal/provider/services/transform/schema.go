// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transform

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// GetTransformResourceSchema returns the schema definition for the transform resource.
func GetTransformResourceSchema() resourceschema.Schema {
	return resourceschema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SailPoint ISC Transform resource",

		Attributes: map[string]resourceschema.Attribute{
			"id": resourceschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Transform identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": resourceschema.StringAttribute{
				MarkdownDescription: "Transform name",
				Required:            true,
			},
			"type": resourceschema.StringAttribute{
				MarkdownDescription: "Transform type",
				Required:            true,
			},
			"internal": resourceschema.BoolAttribute{
				MarkdownDescription: "Indicates if the transform is internal",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"attributes": resourceschema.StringAttribute{
				MarkdownDescription: "Transform attributes as JSON",
				Required:            true,
			},
		},
	}
}

// GetTransformsDataSourceSchema returns the schema definition for the transforms data source.
func GetTransformsDataSourceSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SailPoint ISC Transforms data source",

		Attributes: map[string]schema.Attribute{
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
