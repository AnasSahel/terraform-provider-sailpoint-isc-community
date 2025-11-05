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
	_ datasource.DataSource              = &transformDataSource{}
	_ datasource.DataSourceWithConfigure = &transformDataSource{}
)

type transformDataSource struct {
	client *client.Client
}

func NewTransformDataSource() datasource.DataSource {
	return &transformDataSource{}
}

func (d *transformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *transformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

func (d *transformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.TransformSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Transform.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific SailPoint Transform. Transforms are configurable objects that manipulate attribute data. See [Transform Documentation](https://developer.sailpoint.com/docs/extensibility/transforms/) for more information.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *transformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Transform data source")

	var config models.Transform
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the transform via API
	fetchedTransform, err := d.client.GetTransform(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Transform",
			fmt.Sprintf("Could not read transform ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	var state models.Transform
	if err := state.ConvertFromSailPointForDataSource(ctx, fetchedTransform); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform Response",
			fmt.Sprintf("Could not convert transform response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Transform data source read successfully", map[string]interface{}{
		"transform_id": state.ID.ValueString(),
	})
}
