package identity_attribute

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func resourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Identity Attribute resource schema",
		MarkdownDescription: "Identity Attribute resource schema",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "The name of the identity attribute.",
				MarkdownDescription: "The name of the identity attribute.",
			},
			"display_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "The display name of the identity attribute.",
				MarkdownDescription: "The display name of the identity attribute.",
			},
			"standard": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description:         "Indicates if the identity attribute is a standard attribute.",
				MarkdownDescription: "Indicates if the identity attribute is a standard attribute.",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "The data type of the identity attribute.",
				MarkdownDescription: "The data type of the identity attribute.",
			},
			"multi": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description:         "Indicates if the identity attribute can have multiple values.",
				MarkdownDescription: "Indicates if the identity attribute can have multiple values.",
			},
			"searchable": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description:         "Indicates if the identity attribute is searchable.",
				MarkdownDescription: "Indicates if the identity attribute is searchable.",
			},
			"system": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description:         "Indicates if the identity attribute is a system attribute.",
				MarkdownDescription: "Indicates if the identity attribute is a system attribute.",
			},
			"sources": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "The sources of the identity attribute.",
				MarkdownDescription: "The sources of the identity attribute.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							// Computed:            true,
							Optional:            true,
							Description:         "The type of the source.",
							MarkdownDescription: "The type of the source.",
						},
						"properties": schema.StringAttribute{
							// Computed:            true,
							Optional:            true,
							Description:         "The properties of the source in JSON format.",
							MarkdownDescription: "The properties of the source in JSON format.",
						},
					},
				},
			},
		},
	}
}
