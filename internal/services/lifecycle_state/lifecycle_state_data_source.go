// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &lifecycleStateDataSource{}
	_ datasource.DataSourceWithConfigure = &lifecycleStateDataSource{}
)

type lifecycleStateDataSource struct {
	client *client.Client
}

// NewLifecycleStateDataSource creates a new data source for Lifecycle State.
func NewLifecycleStateDataSource() datasource.DataSource {
	return &lifecycleStateDataSource{}
}

// Metadata implements datasource.DataSource.
func (d *lifecycleStateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lifecycle_state"
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *lifecycleStateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "lifecycle_state data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

// Schema implements datasource.DataSource.
func (d *lifecycleStateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a SailPoint Lifecycle State by ID.",
		MarkdownDescription: "Retrieves a SailPoint Lifecycle State by ID. Lifecycle states define the different stages an identity can be in within an identity profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the lifecycle state.",
				Required:            true,
			},
			"identity_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the identity profile this lifecycle state belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the lifecycle state.",
				Computed:            true,
			},
			"technical_name": schema.StringAttribute{
				MarkdownDescription: "The technical name of the lifecycle state. This is used for internal purposes.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the lifecycle state.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the lifecycle state is enabled.",
				Computed:            true,
			},
			"identity_count": schema.Int32Attribute{
				MarkdownDescription: "The number of identities that have this lifecycle state.",
				Computed:            true,
			},
			"identity_state": schema.StringAttribute{
				MarkdownDescription: "The identity state associated with this lifecycle state. Possible values: `ACTIVE`, `INACTIVE_SHORT_TERM`, `INACTIVE_LONG_TERM`.",
				Computed:            true,
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "The priority of the lifecycle state. Lower numbers appear first when listing with `?sorters=priority`.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date and time the lifecycle state was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the lifecycle state was last modified.",
				Computed:            true,
			},
			"email_notification_option": schema.SingleNestedAttribute{
				MarkdownDescription: "Email notification configuration for the lifecycle state.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"notify_managers": schema.BoolAttribute{
						MarkdownDescription: "If true, managers are notified of lifecycle state changes.",
						Computed:            true,
					},
					"notify_all_admins": schema.BoolAttribute{
						MarkdownDescription: "If true, all admins are notified of lifecycle state changes.",
						Computed:            true,
					},
					"notify_specific_users": schema.BoolAttribute{
						MarkdownDescription: "If true, users specified in `email_address_list` are notified of lifecycle state changes.",
						Computed:            true,
					},
					"email_address_list": schema.ListAttribute{
						MarkdownDescription: "List of email addresses to notify when `notify_specific_users` is true.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"account_actions": schema.ListNestedAttribute{
				MarkdownDescription: "List of account actions to perform when an identity enters this lifecycle state.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							MarkdownDescription: "The action to perform. Possible values: `ENABLE`, `DISABLE`, `DELETE`.",
							Computed:            true,
						},
						"source_ids": schema.ListAttribute{
							MarkdownDescription: "List of source IDs to apply the action to. Required if `all_sources` is false.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"exclude_source_ids": schema.ListAttribute{
							MarkdownDescription: "List of source IDs to exclude from the action. Cannot be used with `source_ids`.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"all_sources": schema.BoolAttribute{
							MarkdownDescription: "If true, the action applies to all sources. If true, `source_ids` must not be provided.",
							Computed:            true,
						},
					},
				},
			},
			"access_profile_ids": schema.ListAttribute{
				MarkdownDescription: "List of access profile IDs associated with this lifecycle state.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"access_action_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "Access action configuration for the lifecycle state.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"remove_all_access_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, all access is removed when an identity enters this lifecycle state.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Read implements datasource.DataSource.
func (d *lifecycleStateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config lifecycleStateModel
	tflog.Debug(ctx, "Getting config for lifecycle state data source")
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identityProfileID := config.IdentityProfileID.ValueString()
	lifecycleStateID := config.ID.ValueString()

	// Read the lifecycle state from SailPoint
	tflog.Debug(ctx, "Fetching lifecycle state from SailPoint", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	lifecycleStateResponse, err := d.client.GetLifecycleState(ctx, identityProfileID, lifecycleStateID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Lifecycle State",
			fmt.Sprintf("Could not read SailPoint Lifecycle State %q in identity profile %q: %s",
				lifecycleStateID, identityProfileID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Lifecycle State", map[string]any{
			"identity_profile_id": identityProfileID,
			"lifecycle_state_id":  lifecycleStateID,
			"error":               err.Error(),
		})
		return
	}

	if lifecycleStateResponse == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Lifecycle State",
			"Received nil response from SailPoint API",
		)
		return
	}

	// Map the response to the data source model
	var state lifecycleStateModel
	tflog.Debug(ctx, "Mapping SailPoint Lifecycle State API response to data source model", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  lifecycleStateID,
	})
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, *lifecycleStateResponse, identityProfileID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	tflog.Debug(ctx, "Setting state for lifecycle state data source", map[string]any{
		"id": state.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Successfully read SailPoint Lifecycle State data source", map[string]any{
		"identity_profile_id": identityProfileID,
		"lifecycle_state_id":  state.ID.ValueString(),
		"name":                state.Name.ValueString(),
	})
}
