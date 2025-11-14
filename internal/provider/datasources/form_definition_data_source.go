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
	_ datasource.DataSource              = &formDefinitionDataSource{}
	_ datasource.DataSourceWithConfigure = &formDefinitionDataSource{}
)

type formDefinitionDataSource struct {
	client *client.Client
}

func NewFormDefinitionDataSource() datasource.DataSource {
	return &formDefinitionDataSource{}
}

func (d *formDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *formDefinitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

func (d *formDefinitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.FormDefinitionSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Form Definition.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Form Definition. Forms are composed of sections and fields for data collection in workflows. See [Custom Forms Documentation](https://developer.sailpoint.com/docs/api/v2025/custom-forms/) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *formDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Form Definition data source")

	var config models.FormDefinition
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the form definition via API
	fetchedForm, err := d.client.GetFormDefinition(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Form Definition",
			fmt.Sprintf("Could not read form definition ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.FormDefinition
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedForm); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Form Definition Response",
			fmt.Sprintf("Could not convert form definition response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Form Definition data source read successfully", map[string]interface{}{
		"form_id": state.ID.ValueString(),
	})
}
