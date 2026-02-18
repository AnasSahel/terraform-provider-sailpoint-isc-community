// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package form_definition

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &formDefinitionDataSource{}
	_ datasource.DataSourceWithConfigure = &formDefinitionDataSource{}
)

type formDefinitionDataSource struct {
	client *client.Client
}

func NewFormDefinitionDataSource() datasource.DataSource {
	return &formDefinitionDataSource{}
}

func (d *formDefinitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

func (d *formDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "form definition data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *formDefinitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Data source for SailPoint Form Definition.",
		MarkdownDescription: "Data source for SailPoint Form Definition. Forms are used to collect data in access requests and workflows.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the form definition.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the form definition.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the form definition.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the form definition.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner (e.g., IDENTITY).",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier of the owner.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
					},
				},
			},
			"used_by": schema.ListNestedAttribute{
				MarkdownDescription: "List of objects that use this form definition.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the referencing object (WORKFLOW, SOURCE, MySailPoint).",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the referencing object.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the referencing object.",
							Computed:            true,
						},
					},
				},
			},
			"form_input": schema.ListNestedAttribute{
				MarkdownDescription: "List of form inputs that can be passed into the form for use in conditional logic.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the form input.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the form input (STRING, ARRAY).",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The label of the form input.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the form input.",
							Computed:            true,
						},
					},
				},
			},
			"form_elements": schema.StringAttribute{
				MarkdownDescription: "JSON array of form elements (fields, sections, etc.). Elements must be wrapped in SECTION elements. Each element object has: id, elementType (TEXT, TOGGLE, TEXTAREA, HIDDEN, PHONE, EMAIL, SELECT, DATE, SECTION, COLUMN_SET, IMAGE, DESCRIPTION), config, key, validations.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"form_conditions": schema.ListNestedAttribute{
				MarkdownDescription: "List of conditions for the form definition. Conditions control the visibility and behavior of form elements based on form inputs and other conditions.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_operator": schema.StringAttribute{
							MarkdownDescription: "The operator for the condition (AND, OR).",
							Computed:            true,
						},
						"rules": schema.ListNestedAttribute{
							MarkdownDescription: "List of rules for the condition.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source_type": schema.StringAttribute{
										MarkdownDescription: "The type of the source for the rule.",
										Computed:            true,
									},
									"source": schema.StringAttribute{
										MarkdownDescription: "The source for the rule.",
										Computed:            true,
									},
									"operator": schema.StringAttribute{
										MarkdownDescription: "The operator for the rule.",
										Computed:            true,
									},
									"value_type": schema.StringAttribute{
										MarkdownDescription: "The type of the value for the rule.",
										Computed:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "The value for the rule.",
										Computed:            true,
									},
								},
							},
						},
						"effects": schema.ListNestedAttribute{
							MarkdownDescription: "List of effects for the condition.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect_type": schema.StringAttribute{
										MarkdownDescription: "The type of the effect.",
										Computed:            true,
									},
									"config": schema.SingleNestedAttribute{
										MarkdownDescription: "The configuration for the effect.",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"default_value_label": schema.StringAttribute{
												MarkdownDescription: "The default value label for the effect.",
												Computed:            true,
											},
											"element": schema.StringAttribute{
												MarkdownDescription: "The element targeted by the effect.",
												Computed:            true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time when the form definition was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the form definition was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *formDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SailPoint Form Definition data source")

	var config formDefinitionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the form definition from SailPoint
	tflog.Debug(ctx, "Fetching form definition from SailPoint", map[string]any{
		"id": config.ID.ValueString(),
	})
	formDefinitionResponse, err := d.client.GetFormDefinition(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Form Definition",
			fmt.Sprintf("Could not read SailPoint Form Definition %q: %s", config.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Form Definition", map[string]any{
			"id":    config.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if formDefinitionResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Form Definition",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the data source model
	var state formDefinitionModel
	resp.Diagnostics.Append(state.FromAPI(ctx, *formDefinitionResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Form Definition data source", map[string]any{
		"id":   config.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}
