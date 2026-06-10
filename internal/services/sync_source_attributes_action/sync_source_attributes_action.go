// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package sync_source_attributes_action

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = &syncSourceAttributesAction{}
	_ action.ActionWithConfigure = &syncSourceAttributesAction{}
)

type syncSourceAttributesAction struct {
	client *client.Client
}

type syncSourceAttributesActionModel struct {
	SourceID types.String `tfsdk:"source_id"`
}

// NewSyncSourceAttributesAction creates a new provider action that triggers a source attribute sync.
func NewSyncSourceAttributesAction() action.Action {
	return &syncSourceAttributesAction{}
}

func (a *syncSourceAttributesAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync_source_attributes"
}

func (a *syncSourceAttributesAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		Description: "Triggers a one-time attribute synchronization for a SailPoint source.",
		MarkdownDescription: "Triggers a one-time attribute synchronization for a SailPoint source. " +
			"Attribute sync pushes identity attribute values from ISC to the target source accounts." +
			"\n\n" +
			"**Experimental**: This action uses the `X-SailPoint-Experimental: true` header on the " +
			"`POST /v2025/sources/{id}/synchronize-attributes` endpoint. Behavior may change without notice." +
			"\n\n" +
			"**Asynchronous**: The sync is triggered immediately but attribute propagation happens " +
			"asynchronously in ISC. The action returns once the request is accepted." +
			"\n\n" +
			"**ORG_ADMIN required**." +
			"\n\n" +
			"**Terraform ≥ 1.14 required** to use actions.",
		Attributes: map[string]actionschema.Attribute{
			"source_id": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the source to synchronize attributes for.",
			},
		},
	}
}

func (a *syncSourceAttributesAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client type for provider data but got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	a.client = c
}

func (a *syncSourceAttributesAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config syncSourceAttributesActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := config.SourceID.ValueString()
	tflog.Debug(ctx, "Invoking sync_source_attributes action", map[string]any{"source_id": sourceID})

	resp.SendProgress(ctx, fmt.Sprintf("Triggering attribute sync for source %q...", sourceID))

	if err := a.client.SyncSourceAttributes(ctx, sourceID); err != nil {
		resp.Diagnostics.AddError(
			"Error Triggering Source Attribute Sync",
			fmt.Sprintf("Could not trigger attribute sync for source %q: %s", sourceID, err.Error()),
		)
		return
	}

	resp.SendProgress(ctx, fmt.Sprintf("Attribute sync triggered for source %q. Propagation is asynchronous.", sourceID))
	tflog.Info(ctx, "Successfully triggered source attribute sync", map[string]any{"source_id": sourceID})
}
