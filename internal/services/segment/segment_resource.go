// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package segment

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &segmentResource{}
	_ resource.ResourceWithConfigure   = &segmentResource{}
	_ resource.ResourceWithImportState = &segmentResource{}
)

type segmentResource struct {
	client *client.Client
}

// NewSegmentResource creates a new Segment resource.
func NewSegmentResource() resource.Resource {
	return &segmentResource{}
}

func (r *segmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (r *segmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "segment resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func (r *segmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Segment.",
		MarkdownDescription: "Resource for SailPoint Segment. Segments control visibility of access items (access profiles, roles, entitlements) by restricting which identities can see and request them based on expression criteria.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the segment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the segment.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the segment.",
				Optional:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the segment is operational. Inactive segments do not restrict visibility.",
				Optional:            true,
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the segment.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner object. Must be `IDENTITY`.",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner. Resolved by the server from the owner ID.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"visibility_criteria": schema.SingleNestedAttribute{
				MarkdownDescription: "Visibility rules that determine which identities the segment applies to.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"expression": schema.SingleNestedAttribute{
						MarkdownDescription: "Root expression node. Either an `EQUALS` leaf, or an `AND` branch with `EQUALS` children.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"operator": schema.StringAttribute{
								MarkdownDescription: "Operator for this node. One of `AND`, `EQUALS`.",
								Required:            true,
							},
							"attribute": schema.StringAttribute{
								MarkdownDescription: "Identity attribute to compare. Required when `operator` is `EQUALS`.",
								Optional:            true,
							},
							"value": schema.SingleNestedAttribute{
								MarkdownDescription: "Typed value. Required when `operator` is `EQUALS`.",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "Value type (e.g., `STRING`).",
										Required:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "Value to compare against.",
										Required:            true,
									},
								},
							},
							"children": schema.ListNestedAttribute{
								MarkdownDescription: "Child leaf expressions. Required when `operator` is `AND`. Children cannot have further children.",
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"operator": schema.StringAttribute{
											MarkdownDescription: "Operator for this leaf. Typically `EQUALS`.",
											Required:            true,
										},
										"attribute": schema.StringAttribute{
											MarkdownDescription: "Identity attribute to compare.",
											Optional:            true,
										},
										"value": schema.SingleNestedAttribute{
											MarkdownDescription: "Typed value.",
											Optional:            true,
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													MarkdownDescription: "Value type (e.g., `STRING`).",
													Required:            true,
												},
												"value": schema.StringAttribute{
													MarkdownDescription: "Value to compare against.",
													Required:            true,
												},
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
				MarkdownDescription: "The date and time the segment was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the segment was last modified.",
				Computed:            true,
			},
		},
	}
}

func (r *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan segmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating segment", map[string]any{"name": plan.Name.ValueString()})
	apiResp, err := r.client.CreateSegment(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Segment",
			fmt.Sprintf("Could not create SailPoint Segment %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Creating SailPoint Segment", "Received nil response from SailPoint API")
		return
	}

	var state segmentModel
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully created segment", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	apiResp, err := r.client.GetSegment(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Segment not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Segment",
			fmt.Sprintf("Could not read SailPoint Segment %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Segment", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan segmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	ops, diags := plan.ToPatchOperations(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.PatchSegment(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Segment",
			fmt.Sprintf("Could not update SailPoint Segment %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Segment", "Received nil response from SailPoint API")
		return
	}

	var newState segmentModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated segment", map[string]any{
		"id":   newState.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

func (r *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if err := r.client.DeleteSegment(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Segment",
			fmt.Sprintf("Could not delete SailPoint Segment %q: %s", id, err.Error()),
		)
		return
	}
	tflog.Info(ctx, "Successfully deleted segment", map[string]any{"id": id})
}

func (r *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
