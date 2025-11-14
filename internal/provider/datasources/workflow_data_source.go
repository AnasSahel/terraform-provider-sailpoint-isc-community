// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &workflowDataSource{}
	_ datasource.DataSourceWithConfigure = &workflowDataSource{}
)

type workflowDataSource struct {
	client *client.Client
}

func NewWorkflowDataSource() datasource.DataSource {
	return &workflowDataSource{}
}

func (d *workflowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *workflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *workflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.WorkflowSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Workflow.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Workflow. Workflows are custom automation scripts that respond to event triggers. See [Workflow Documentation](https://developer.sailpoint.com/docs/extensibility/workflows/) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Workflow data source")

	var config models.Workflow
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the workflow via API
	fetchedWorkflow, err := d.client.GetWorkflow(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workflow",
			fmt.Sprintf("Could not read workflow ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.Workflow
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedWorkflow); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Workflow Response",
			fmt.Sprintf("Could not convert workflow response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Workflow data source read successfully", map[string]interface{}{
		"workflow_id": state.ID.ValueString(),
	})
}
