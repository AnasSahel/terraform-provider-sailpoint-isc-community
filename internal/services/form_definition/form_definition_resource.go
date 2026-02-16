// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package form_definition

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &formDefinitionResource{}
	_ resource.ResourceWithConfigure   = &formDefinitionResource{}
	_ resource.ResourceWithImportState = &formDefinitionResource{}
)

type formDefinitionResource struct {
	client *client.Client
}

func NewFormDefinitionResource() resource.Resource {
	return &formDefinitionResource{}
}

// Metadata implements resource.Resource.
func (r *formDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

// Configure implements resource.ResourceWithConfigure.
func (r *formDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "form definition resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *formDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Form Definition.",
		MarkdownDescription: "Resource for SailPoint Form Definition. Forms are used to collect data in access requests and workflows.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the form definition.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the form definition.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the form definition.",
				Optional:            true,
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the form definition.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner (e.g., IDENTITY).",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"used_by": schema.ListNestedAttribute{
				MarkdownDescription: "List of objects that use this form definition.",
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
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
				Optional:            true,
				Computed:            true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"id":          types.StringType,
						"type":        types.StringType,
						"label":       types.StringType,
						"description": types.StringType,
					}},
					[]attr.Value{},
				)),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the form input.",
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the form input (STRING, ARRAY).",
							Required:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The label of the form input.",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the form input.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"form_elements": schema.StringAttribute{
				MarkdownDescription: "JSON array of form elements (fields, sections, etc.). Elements must be wrapped in SECTION elements. Each element object has: id, elementType (TEXT, TOGGLE, TEXTAREA, HIDDEN, PHONE, EMAIL, SELECT, DATE, SECTION, COLUMN_SET, IMAGE, DESCRIPTION), config, key, validations.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"form_conditions": schema.ListNestedAttribute{
				MarkdownDescription: "List of conditions for the form definition. Conditions control the visibility and behavior of form elements based on form inputs and other conditions.",
				Optional:            true,
				Computed:            true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.ObjectType{AttrTypes: map[string]attr.Type{
						"rule_operator": types.StringType,
						"rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
							"source_type": types.StringType,
							"source":      types.StringType,
							"operator":    types.StringType,
							"value_type":  types.StringType,
							"value":       types.StringType,
						}}},
						"effects": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
							"effect_type": types.StringType,
							"config": types.ObjectType{AttrTypes: map[string]attr.Type{
								"default_value_label": types.StringType,
								"element":             types.StringType,
							}},
						}}},
					}},
					[]attr.Value{},
				)),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_operator": schema.StringAttribute{
							MarkdownDescription: "The operator for the condition (AND, OR).",
							Required:            true,
						},
						"rules": schema.ListNestedAttribute{
							MarkdownDescription: "List of rules for the condition.",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source_type": schema.StringAttribute{
										MarkdownDescription: "The type of the source for the rule (INPUT, ELEMENT).",
										Required:            true,
									},
									"source": schema.StringAttribute{
										MarkdownDescription: "The source for the rule.",
										Required:            true,
									},
									"operator": schema.StringAttribute{
										MarkdownDescription: "The operator for the rule (EQ, NE, CO, NOT_CO, IN, NOT_IN, EM, NOT_EM, SW, NOT_SW, EW, NOT_EW).",
										Required:            true,
									},
									"value_type": schema.StringAttribute{
										MarkdownDescription: "The type of the value for the rule (STRING, STRING_LIST, INPUT, ELEMENT, LIST, BOOLEAN).",
										Required:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "The value for the rule.",
										Optional:            true,
										Computed:            true,
									},
								},
							},
						},
						"effects": schema.ListNestedAttribute{
							MarkdownDescription: "List of effects for the condition.",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect_type": schema.StringAttribute{
										MarkdownDescription: "The type of the effect (HIDE, SHOW, DISABLE, ENABLE, REQUIRE, OPTIONAL, SUBMIT_MESSAGE, SUBMIT_NOTIFICATION, SET_DEFAULT_VALUE).",
										Required:            true,
									},
									"config": schema.SingleNestedAttribute{
										MarkdownDescription: "The configuration for the effect.",
										Required:            true,
										Attributes: map[string]schema.Attribute{
											"default_value_label": schema.StringAttribute{
												MarkdownDescription: "The default value label for the effect.",
												Optional:            true,
												Computed:            true,
											},
											"element": schema.StringAttribute{
												MarkdownDescription: "The element targeted by the effect.",
												Optional:            true,
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the form definition was last modified.",
				Computed:            true,
			},
		},
	}
}

// Create implements resource.Resource.
func (r *formDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the plan
	var plan formDefinitionModel
	tflog.Debug(ctx, "Getting plan for form definition resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map resource model to API model
	tflog.Debug(ctx, "Mapping form definition resource model to API create request", map[string]any{
		"name": plan.Name.ValueString(),
	})
	apiCreateRequest, diags := plan.ToAPICreateRequest(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the form definition via the API client
	tflog.Debug(ctx, "Creating form definition via SailPoint API", map[string]any{
		"name": plan.Name.ValueString(),
	})
	formDefinitionAPIResponse, err := r.client.CreateFormDefinition(ctx, &apiCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Form Definition",
			fmt.Sprintf("Could not create SailPoint Form Definition %q: %s", plan.Name.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to create SailPoint Form Definition", map[string]any{
			"name":  plan.Name.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if formDefinitionAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Form Definition",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var state formDefinitionModel
	tflog.Debug(ctx, "Mapping SailPoint Form Definition API response to resource model", map[string]any{
		"name": plan.Name.ValueString(),
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *formDefinitionAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for form definition resource", map[string]any{
		"name": plan.Name.ValueString(),
		"id":   state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Form Definition resource", map[string]any{
		"name": plan.Name.ValueString(),
		"id":   state.ID.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *formDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state formDefinitionModel
	tflog.Debug(ctx, "Getting state for form definition resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the form definition from SailPoint
	tflog.Debug(ctx, "Fetching form definition from SailPoint", map[string]any{
		"id": state.ID.ValueString(),
	})
	formDefinitionResponse, err := r.client.GetFormDefinition(ctx, state.ID.ValueString())
	if err != nil {
		// If resource was deleted outside of Terraform, remove it from state
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Form Definition not found, removing from state", map[string]any{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Form Definition",
			fmt.Sprintf("Could not read SailPoint Form Definition %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Form Definition", map[string]any{
			"id":    state.ID.ValueString(),
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

	// Map the response to the resource model
	tflog.Debug(ctx, "Mapping SailPoint Form Definition API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *formDefinitionResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for form definition resource", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Form Definition resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *formDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state formDefinitionModel
	tflog.Debug(ctx, "Getting state for form definition resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting form definition via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	err := r.client.DeleteFormDefinition(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Form Definition",
			fmt.Sprintf("Could not delete SailPoint Form Definition %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to delete SailPoint Form Definition", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Form Definition resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *formDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan formDefinitionModel
	tflog.Debug(ctx, "Getting plan for form definition resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to retrieve the ID and compare for changes
	var state formDefinitionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate patch operations for changed fields only
	tflog.Debug(ctx, "Generating patch operations for form definition update", map[string]any{
		"id": state.ID.ValueString(),
	})
	patchOps, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log each patch operation for debugging
	for i, op := range patchOps {
		tflog.Debug(ctx, "Patch operation", map[string]any{
			"index": i,
			"op":    op.Op,
			"path":  op.Path,
		})
	}

	// Update the form definition via the API client
	tflog.Debug(ctx, "Updating form definition via SailPoint API", map[string]any{
		"id":          state.ID.ValueString(),
		"patch_count": len(patchOps),
	})
	formDefinitionAPIResponse, err := r.client.PatchFormDefinition(ctx, state.ID.ValueString(), patchOps)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Form Definition",
			fmt.Sprintf("Could not update SailPoint Form Definition %q: %s", state.ID.ValueString(), err.Error()),
		)
		tflog.Error(ctx, "Failed to update SailPoint Form Definition", map[string]any{
			"id":    state.ID.ValueString(),
			"error": err.Error(),
		})
		return
	}

	if formDefinitionAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Form Definition",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the API response back to the resource model
	var newState formDefinitionModel
	tflog.Debug(ctx, "Mapping SailPoint Form Definition API response to resource model", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(newState.FromSailPointAPI(ctx, *formDefinitionAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for form definition resource", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Form Definition resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

func (r *formDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by setting the "id" attribute
	tflog.Debug(ctx, "Importing form definition resource", map[string]any{
		"id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Form Definition resource", map[string]any{
		"id": req.ID,
	})
}
