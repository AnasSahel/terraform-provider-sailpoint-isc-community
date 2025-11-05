package datasources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = &sourceDataSource{}
)

type sourceDataSource struct {
	client *client.Client
}

func NewSourceDataSource() datasource.DataSource {
	return new(sourceDataSource)
}

func (d *sourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := utils.ConfigureClient(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = client
}

func (d *sourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (d *sourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaBuilder := schemas.SourceSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Data source for retrieving information about a SailPoint Identity Security Cloud (ISC) source.",
		MarkdownDescription: "Use this data source to retrieve detailed information about a specific source in SailPoint ISC. Sources represent systems or applications from which identity data is aggregated.",
		Attributes:          schemaBuilder.GetDataSourceSchema(),
	}
}

func (d *sourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching source by ID", map[string]any{"id": state.ID.ValueString()})

	source, err := d.client.GetSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Fetch Source",
			fmt.Sprintf("Failed to fetch source with ID %q: %v", state.ID.ValueString(), err),
		)
		return
	}

	state.ConvertFromSailPointForDataSource(ctx, source)
	tflog.Debug(ctx, "Successfully fetched source", map[string]any{"source_id": source.ID})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Helper functions
