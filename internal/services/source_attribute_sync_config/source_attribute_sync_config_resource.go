// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_attribute_sync_config

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
	_ resource.Resource                = &sourceAttrSyncConfigResource{}
	_ resource.ResourceWithConfigure   = &sourceAttrSyncConfigResource{}
	_ resource.ResourceWithImportState = &sourceAttrSyncConfigResource{}
)

type sourceAttrSyncConfigResource struct {
	client *client.Client
}

// NewSourceAttributeSyncConfigResource creates a new adopt-only resource for source attribute sync config.
func NewSourceAttributeSyncConfigResource() resource.Resource {
	return &sourceAttrSyncConfigResource{}
}

func (r *sourceAttrSyncConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_attribute_sync_config"
}

func (r *sourceAttrSyncConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source attribute sync config resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func (r *sourceAttrSyncConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the attribute synchronization configuration for a SailPoint source (Beta API). " +
			"Adopt-only lifecycle: Create reads the existing config and applies declared enabled flags via PUT, " +
			"Update sends a full PUT, and Delete is a no-op. " +
			"Partial management is supported: unlisted attributes keep their current server value. " +
			"Requires ORG_ADMIN. Terraform >= 1.14 required for lifecycle-triggered actions.",
		MarkdownDescription: "Manages the attribute synchronization configuration for a SailPoint source (Beta API). " +
			"Adopt-only lifecycle: Create reads the existing config and applies declared enabled flags via PUT, " +
			"Update sends a full PUT, and Delete is a no-op. " +
			"Partial management is supported: unlisted attributes keep their current server value. " +
			"Requires ORG_ADMIN. Terraform >= 1.14 required for lifecycle-triggered actions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier — mirrors `source_id` for import convenience.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the source whose attribute sync config is being managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Reference to the source. Server-populated.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"attributes": schema.SetNestedAttribute{
				Required: true,
				MarkdownDescription: "The set of identity attributes to manage sync enablement for. " +
					"Only specify attributes you want to control; unspecified attributes retain their current server value.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The identity attribute name (e.g. `email`, `department`). Must exist in the source's Create Account definition.",
						},
						"enabled": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether this attribute should be synced to the source.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The human-readable display name of the attribute. Server-provided.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"target": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The account attribute name on the target source (e.g. `mail`). Server-provided.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

// Create adopts the existing attribute sync config for a source, then applies the plan's `enabled` flags.
func (r *sourceAttrSyncConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceAttrSyncConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := plan.SourceID.ValueString()
	tflog.Debug(ctx, "Adopting source attribute sync config", map[string]any{"source_id": sourceID})

	existing, err := r.client.GetSourceAttributeSyncConfig(ctx, sourceID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Source not found",
				fmt.Sprintf("Source %q does not exist in SailPoint ISC.", sourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error adopting Source Attribute Sync Config",
			fmt.Sprintf("Could not read attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if existing == nil {
		resp.Diagnostics.AddError("Error adopting Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	// Build PUT body: merge plan's enabled flags into the full remote list.
	apiReq, diags := plan.ToAPI(ctx, existing.Attributes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.PutSourceAttributeSyncConfig(ctx, sourceID, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying initial changes to Source Attribute Sync Config",
			fmt.Sprintf("Could not update attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if updated == nil {
		resp.Diagnostics.AddError("Error applying initial changes to Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	var state sourceAttrSyncConfigModel
	state.SourceID = plan.SourceID
	cfgNames := configuredAttributeNames(plan.Attributes)
	resp.Diagnostics.Append(state.FromAPI(ctx, updated, cfgNames)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully adopted source attribute sync config", map[string]any{"source_id": sourceID})
}

func (r *sourceAttrSyncConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceAttrSyncConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	apiResp, err := r.client.GetSourceAttributeSyncConfig(ctx, sourceID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Source not found, removing attribute sync config from state", map[string]any{"source_id": sourceID})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Source Attribute Sync Config",
			fmt.Sprintf("Could not read attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	// Filter to the configured subset to avoid a perpetual diff when managing only some attributes.
	cfgNames := configuredAttributeNames(state.Attributes)
	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp, cfgNames)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceAttrSyncConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceAttrSyncConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state sourceAttrSyncConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	plan.SourceID = state.SourceID

	// Fetch the full current remote state to use as the merge base.
	existing, err := r.client.GetSourceAttributeSyncConfig(ctx, sourceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Attribute Sync Config",
			fmt.Sprintf("Could not read current attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if existing == nil {
		resp.Diagnostics.AddError("Error Updating Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	apiReq, diags := plan.ToAPI(ctx, existing.Attributes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.PutSourceAttributeSyncConfig(ctx, sourceID, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Attribute Sync Config",
			fmt.Sprintf("Could not update attribute sync config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating Source Attribute Sync Config", "Received nil response from SailPoint API")
		return
	}

	var newState sourceAttrSyncConfigModel
	newState.SourceID = state.SourceID
	cfgNames := configuredAttributeNames(plan.Attributes)
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp, cfgNames)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated source attribute sync config", map[string]any{"source_id": sourceID})
}

// Delete is a no-op — attribute sync configs always exist on the server for a source;
// this only removes the resource from Terraform state.
func (r *sourceAttrSyncConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Source Attribute Sync Config not deleted from SailPoint ISC",
		"Attribute sync configs always exist for a source and cannot be deleted via the API. "+
			"The resource has been removed from Terraform state, but the config still exists in ISC.",
	)
}

// ImportState imports by source_id.
func (r *sourceAttrSyncConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
