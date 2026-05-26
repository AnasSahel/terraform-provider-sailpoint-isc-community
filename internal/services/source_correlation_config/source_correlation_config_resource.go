// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_correlation_config

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sourceCorrelationConfigResource{}
	_ resource.ResourceWithConfigure   = &sourceCorrelationConfigResource{}
	_ resource.ResourceWithImportState = &sourceCorrelationConfigResource{}
)

type sourceCorrelationConfigResource struct {
	client *client.Client
}

// NewSourceCorrelationConfigResource creates a new resource for SailPoint source correlation config.
func NewSourceCorrelationConfigResource() resource.Resource {
	return &sourceCorrelationConfigResource{}
}

func (r *sourceCorrelationConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_correlation_config"
}

func (r *sourceCorrelationConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source correlation config resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func attributeAssignmentNestedAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Required: true,
		MarkdownDescription: "List of account-to-identity attribute mappings that drive correlation. " +
			"Each entry maps a source account attribute (`property`) to an identity attribute (`value`).",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"property": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The source account attribute name (e.g. `sAMAccountName`).",
				},
				"value": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The identity attribute name to correlate against (e.g. `uid`).",
				},
				"operation": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Comparison operation. Defaults to `EQ`.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"complex": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Whether the assignment uses a complex expression. Server-resolved when not set.",
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"ignore_case": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Whether the comparison is case-insensitive.",
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"match_mode": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Match mode for string comparisons (e.g. `ANYWHERE`, `START`, `END`, `EXACT`).",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"filter_string": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "An additional filter expression applied to candidate accounts before correlation.",
				},
			},
		},
	}
}

func (r *sourceCorrelationConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for SailPoint Source Correlation Config.",
		MarkdownDescription: "Manages the correlation configuration for a SailPoint source. " +
			"Correlation config determines how incoming account data is matched to known identities. " +
			"\n\n" +
			"This is a singleton-per-source resource using an adopt-only lifecycle: " +
			"Create reads the existing config and applies any declared changes, " +
			"Update replaces the config via PUT, and Delete is a no-op that only removes the resource " +
			"from Terraform state (the config always exists on the server).",
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
				MarkdownDescription: "The ID of the source whose correlation config is being managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Server-assigned name of the correlation config (e.g. `AD Correlation Config`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attribute_assignments": attributeAssignmentNestedAttr(),
		},
	}
}

// Create adopts the existing correlation config for a source. If the plan differs
// from the current server state, it issues a PUT to apply the changes.
func (r *sourceCorrelationConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceCorrelationConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := plan.SourceID.ValueString()
	tflog.Debug(ctx, "Adopting source correlation config", map[string]any{"source_id": sourceID})

	existing, err := r.client.GetSourceCorrelationConfig(ctx, sourceID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Source not found",
				fmt.Sprintf("Source %q does not exist in SailPoint ISC.", sourceID),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error adopting Source Correlation Config",
			fmt.Sprintf("Could not read correlation config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if existing == nil {
		resp.Diagnostics.AddError("Error adopting Source Correlation Config", "Received nil response from SailPoint API")
		return
	}

	// Populate a baseline from the API response so we can diff.
	var current sourceCorrelationConfigModel
	current.SourceID = plan.SourceID
	resp.Diagnostics.Append(current.FromAPI(ctx, existing)...)
	if resp.Diagnostics.HasError() {
		return
	}

	final := existing
	if plan.NeedsUpdate(&current) {
		apiReq := plan.ToAPI()
		updated, err := r.client.PutSourceCorrelationConfig(ctx, sourceID, apiReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error applying initial changes to Source Correlation Config",
				fmt.Sprintf("Could not update correlation config for source %q: %s", sourceID, err.Error()),
			)
			return
		}
		if updated != nil {
			final = updated
		}
	}

	var state sourceCorrelationConfigModel
	state.SourceID = plan.SourceID
	resp.Diagnostics.Append(state.FromAPI(ctx, final)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully adopted source correlation config", map[string]any{"source_id": sourceID})
}

func (r *sourceCorrelationConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceCorrelationConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	apiResp, err := r.client.GetSourceCorrelationConfig(ctx, sourceID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Source not found, removing correlation config from state", map[string]any{"source_id": sourceID})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Source Correlation Config",
			fmt.Sprintf("Could not read correlation config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading Source Correlation Config", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sourceCorrelationConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceCorrelationConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state sourceCorrelationConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	plan.SourceID = state.SourceID

	apiReq := plan.ToAPI()
	apiResp, err := r.client.PutSourceCorrelationConfig(ctx, sourceID, apiReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Correlation Config",
			fmt.Sprintf("Could not update correlation config for source %q: %s", sourceID, err.Error()),
		)
		return
	}
	if apiResp == nil {
		resp.Diagnostics.AddError("Error Updating Source Correlation Config", "Received nil response from SailPoint API")
		return
	}

	var newState sourceCorrelationConfigModel
	newState.SourceID = state.SourceID
	resp.Diagnostics.Append(newState.FromAPI(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Info(ctx, "Successfully updated source correlation config", map[string]any{"source_id": sourceID})
}

// Delete is a no-op — correlation configs always exist on the server;
// this only removes the resource from Terraform state.
func (r *sourceCorrelationConfigResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Source Correlation Config not deleted from SailPoint ISC",
		"Correlation configs always exist for a source and cannot be deleted via the API. "+
			"The resource has been removed from Terraform state, but the config still exists in ISC.",
	)
}

// ImportState imports via source_id.
func (r *sourceCorrelationConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by source_id: set both source_id and id to the provided value.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
