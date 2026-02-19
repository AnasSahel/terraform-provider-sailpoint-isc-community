// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"errors"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sourceResource{}
	_ resource.ResourceWithConfigure   = &sourceResource{}
	_ resource.ResourceWithImportState = &sourceResource{}
)

type sourceResource struct {
	client *client.Client
}

func NewSourceResource() resource.Resource {
	return &sourceResource{}
}

// Metadata implements resource.Resource.
func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

// Configure implements resource.ResourceWithConfigure.
func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// Schema implements resource.Resource.
func (r *sourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource for SailPoint Source.",
		MarkdownDescription: "Resource for SailPoint Source. Sources represent managed systems (e.g., Active Directory, Workday) in Identity Security Cloud.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the source.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the source.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the source.",
				Optional:            true,
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the source.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner (e.g., `IDENTITY`).",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the owner.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the owner. Resolved by the server from the owner ID.",
						Computed:            true,
					},
				},
			},
			"cluster": schema.SingleNestedAttribute{
				MarkdownDescription: "The cluster associated with this source. Required for on-premise sources.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the cluster (e.g., `CLUSTER`).",
						Required:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "The ID of the cluster.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the cluster. Resolved by the server from the cluster ID.",
						Computed:            true,
					},
				},
			},
			"connector": schema.StringAttribute{
				MarkdownDescription: "The connector script name. Cannot be changed after creation.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_class": schema.StringAttribute{
				MarkdownDescription: "The fully qualified name of the Java class that implements the connector interface. Cannot be changed after creation.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_attributes": schema.StringAttribute{
				MarkdownDescription: "A JSON object containing connector-specific configuration. " +
					"The server may add extra fields (e.g., `beforeProvisioningRule`, `since`) and modify values (e.g., `cloudDisplayName`) on creation and updates. " +
					"Only the keys you specify in your configuration are managed by Terraform; server-added keys will appear in state after the first refresh.",
				Optional:   true,
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"connection_type": schema.StringAttribute{
				MarkdownDescription: "The connection type (e.g., `direct`, `file`).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of system being managed. Cannot be changed after creation.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_threshold": schema.Int64Attribute{
				MarkdownDescription: "The percentage threshold for skipping the delete phase (0-100).",
				Optional:            true,
				Computed:            true,
			},
			"authoritative": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is referenced by an identity profile. Cannot be changed after creation.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"healthy": schema.BoolAttribute{
				MarkdownDescription: "Whether the source is healthy.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the source (e.g., `SOURCE_STATE_HEALTHY`, `SOURCE_STATE_ERROR_ACCOUNT_FILE_IMPORT`).",
				Computed:            true,
			},
			"features": schema.ListAttribute{
				MarkdownDescription: "The list of features enabled for the source (e.g., `PROVISIONING`, `SYNC_PROVISIONING`, `AUTHENTICATE`).",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"credential_provider_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether credential provider is enabled for the source.",
				Optional:            true,
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The source category (e.g., `CredentialProvider`).",
				Optional:            true,
			},
			"provision_as_csv": schema.BoolAttribute{
				MarkdownDescription: "If `true`, configures the source as a Delimited File (CSV) source during creation. This is a create-only parameter and cannot be changed after creation.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the source was last modified.",
				Computed:            true,
			},
		},
	}
}

// Create implements resource.Resource.
func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceModel
	tflog.Debug(ctx, "Getting plan for source resource")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCreateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	provisionAsCsv := !plan.ProvisionAsCsv.IsNull() && plan.ProvisionAsCsv.ValueBool()

	tflog.Debug(ctx, "Creating source via SailPoint API", map[string]any{
		"name":             plan.Name.ValueString(),
		"provision_as_csv": provisionAsCsv,
	})
	sourceAPIResponse, err := r.client.CreateSource(ctx, &apiCreateRequest, provisionAsCsv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Source",
			fmt.Sprintf("Could not create SailPoint Source %q: %s", plan.Name.ValueString(), err.Error()),
		)
		return
	}

	if sourceAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Source",
			"Received nil response from SailPoint API",
		)
		return
	}

	var state sourceModel
	resp.Diagnostics.Append(state.FromAPI(ctx, *sourceAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the planned connector_attributes to avoid "inconsistent result after apply".
	// The API enriches this field with server-populated defaults (e.g., beforeProvisioningRule,
	// since, status) and may modify user-provided values (e.g., cloudDisplayName).
	if !plan.ConnectorAttributes.IsNull() && !plan.ConnectorAttributes.IsUnknown() {
		state.ConnectorAttributes = plan.ConnectorAttributes
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully created SailPoint Source resource", map[string]any{
		"name": plan.Name.ValueString(),
		"id":   state.ID.ValueString(),
	})
}

// Read implements resource.Resource.
func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceModel
	tflog.Debug(ctx, "Getting state for source resource read")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching source from SailPoint", map[string]any{
		"id": state.ID.ValueString(),
	})
	sourceResponse, err := r.client.GetSource(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "SailPoint Source not found, removing from state", map[string]any{
				"id": state.ID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source",
			fmt.Sprintf("Could not read SailPoint Source %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	if sourceResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source",
			"Received nil response from SailPoint API",
		)
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, *sourceResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Source resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

// Update implements resource.Resource.
func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceModel
	tflog.Debug(ctx, "Getting plan for source resource update")
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state sourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiUpdateRequest, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating source via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	sourceAPIResponse, err := r.client.UpdateSource(ctx, state.ID.ValueString(), &apiUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Source",
			fmt.Sprintf("Could not update SailPoint Source %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	if sourceAPIResponse == nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Source",
			"Received nil response from SailPoint API",
		)
		return
	}

	var newState sourceModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, *sourceAPIResponse)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully updated SailPoint Source resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": newState.Name.ValueString(),
	})
}

// Delete implements resource.Resource.
func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceModel
	tflog.Debug(ctx, "Getting state for source resource deletion")
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting source via SailPoint API", map[string]any{
		"id": state.ID.ValueString(),
	})
	err := r.client.DeleteSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Source",
			fmt.Sprintf("Could not delete SailPoint Source %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Info(ctx, "Successfully deleted SailPoint Source resource", map[string]any{
		"id":   state.ID.ValueString(),
		"name": state.Name.ValueString(),
	})
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing source resource", map[string]any{
		"id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Successfully imported SailPoint Source resource", map[string]any{
		"id": req.ID,
	})
}
