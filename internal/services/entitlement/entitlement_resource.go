// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package entitlement

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &entitlementResource{}
	_ resource.ResourceWithConfigure   = &entitlementResource{}
	_ resource.ResourceWithImportState = &entitlementResource{}
)

type entitlementResource struct {
	client *client.Client
}

// NewEntitlementResource creates a new Entitlement resource with adopt-only lifecycle.
func NewEntitlementResource() resource.Resource {
	return &entitlementResource{}
}

func (r *entitlementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entitlement"
}

func (r *entitlementResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "entitlement resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func (r *entitlementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Adopts an existing SailPoint Entitlement and manages its patchable metadata.",
		MarkdownDescription: "Adopts an existing SailPoint Entitlement and manages its patchable metadata." +
			" Entitlements cannot be created or deleted via the API — they are managed by source aggregation." +
			" This resource uses an adopt-only lifecycle: Create reads the existing entitlement, Update patches metadata," +
			" and Delete is a no-op that only removes the resource from Terraform state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of an existing entitlement.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitlement. Patchable; overrides the source-aggregated name when set.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the entitlement. Patchable.",
				Optional:            true,
				Computed:            true,
			},
			"attribute": schema.StringAttribute{
				MarkdownDescription: "Source attribute name (e.g., `memberOf`). Read-only from aggregation.",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Source attribute value (e.g., a group DN). Read-only from aggregation.",
				Computed:            true,
			},
			"source_schema_object_type": schema.StringAttribute{
				MarkdownDescription: "Type of the entitlement in the source schema (e.g., `group`). Read-only.",
				Computed:            true,
			},
			"privileged": schema.BoolAttribute{
				MarkdownDescription: "Whether the entitlement grants elevated access. Patchable.",
				Optional:            true,
				Computed:            true,
			},
			"cloud_governed": schema.BoolAttribute{
				MarkdownDescription: "Whether the entitlement is cloud-governed. Read-only.",
				Computed:            true,
			},
			"requestable": schema.BoolAttribute{
				MarkdownDescription: "Whether users can request this entitlement directly. Patchable.",
				Optional:            true,
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the entitlement. Patchable.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Owner type. Must be `IDENTITY`.",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "Identity ID of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name of the owner identity. Server-resolved.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"source": schema.SingleNestedAttribute{
				MarkdownDescription: "Source the entitlement was aggregated from. Read-only.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"segments": schema.SetAttribute{
				MarkdownDescription: "Segment UUIDs the entitlement is assigned to. Patchable.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"manually_updated_fields": schema.MapAttribute{
				MarkdownDescription: "Tracks which fields were manually overridden (protected from aggregation overwrites). Read-only.",
				Computed:            true,
				ElementType:         types.BoolType,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the entitlement was first aggregated.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "When the entitlement was last modified.",
				Computed:            true,
			},
		},
	}
}

// Create adopts an existing entitlement by ID. The entitlement must already exist in ISC —
// entitlements are managed via source aggregation, not via Terraform.
func (r *entitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan entitlementModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	tflog.Debug(ctx, "Adopting entitlement", map[string]any{"id": id})

	existing, err := r.client.GetEntitlement(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Entitlement not found",
				fmt.Sprintf("Entitlement %q does not exist in SailPoint ISC. Entitlements are created through source aggregation, not via Terraform.", id),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error adopting SailPoint Entitlement",
			fmt.Sprintf("Could not read entitlement %q: %s", id, err.Error()),
		)
		return
	}
	if existing == nil {
		resp.Diagnostics.AddError("Error adopting SailPoint Entitlement", "Received nil response from SailPoint API")
		return
	}

	// Populate a baseline state from the existing entitlement so we can diff the plan against it.
	var current entitlementModel
	resp.Diagnostics.Append(current.FromAPI(ctx, existing)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Diff plan vs current: any patchable fields the user set that differ from the current API
	// state are applied in a single PATCH.
	ops, diags := plan.ToPatchOperations(ctx, &current)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	final := existing
	if len(ops) > 0 {
		updated, err := r.client.PatchEntitlement(ctx, id, ops)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error applying initial patches to adopted Entitlement",
				fmt.Sprintf("Could not patch entitlement %q: %s", id, err.Error()),
			)
			return
		}
		if updated != nil {
			final = updated
		}
	}

	var state entitlementModel
	resp.Diagnostics.Append(state.FromAPI(ctx, final)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully adopted entitlement", map[string]any{
		"id":      state.ID.ValueString(),
		"name":    state.Name.ValueString(),
		"patches": len(ops),
	})
}

func (r *entitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state entitlementModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	apiResp, err := r.client.GetEntitlement(ctx, id)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Entitlement not found, removing from state", map[string]any{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Entitlement",
			fmt.Sprintf("Could not read entitlement %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Entitlement", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *entitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan entitlementModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state entitlementModel
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

	apiResp, err := r.client.PatchEntitlement(ctx, id, ops)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Entitlement",
			fmt.Sprintf("Could not update entitlement %q: %s", id, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating SailPoint Entitlement", "Received nil response from SailPoint API")
		return
	}

	var newState entitlementModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated entitlement", map[string]any{
		"id":      newState.ID.ValueString(),
		"patches": len(ops),
	})
}

// Delete is a no-op for entitlements — they are managed by source aggregation
// and cannot be removed via the API. Terraform state tracking is dropped, but
// the entitlement persists in ISC.
func (r *entitlementResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Intentional no-op.
	resp.Diagnostics.AddWarning(
		"Entitlement not deleted from SailPoint ISC",
		"Entitlements are managed by source aggregation and cannot be deleted via the API. "+
			"The resource has been removed from Terraform state, but the entitlement still exists in ISC.",
	)
}

func (r *entitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
