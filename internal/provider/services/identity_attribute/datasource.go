package identity_attribute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

var (
	_ datasource.DataSource              = &identityAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &identityAttributeDataSource{}
)

type identityAttributeDataSource struct {
	client *sailpoint.APIClient
}

func NewIdentityAttributeDataSource() datasource.DataSource {
	return &identityAttributeDataSource{}
}

func (d *identityAttributeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)
	if !ok {
		resp.Diagnostics.AddError(ErrProviderDataTitle, fmt.Sprintf(ErrProviderDataMsg, req.ProviderData))
		return
	}

	d.client = client
}

func (d *identityAttributeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + DataSourceTypeName
}

func (d *identityAttributeDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema()
}

func (d *identityAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IdentityAttributeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	identityAttribute, httpResponse, err := d.client.V2025.IdentityAttributesAPI.GetIdentityAttribute(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(ErrDataSourceReadTitle, fmt.Sprintf(ErrDataSourceReadMsg, data.Name.ValueString(), err, httpResponse))
		return
	}

	state := MapIdentityAttributeToDataSourceModel(*identityAttribute)

	// Map sources
	data.Sources = []Source1{}
	for _, source := range identityAttribute.GetSources() {
		sourceJson, _ := json.Marshal(source.GetProperties())
		data.Sources = append(data.Sources, Source1{
			Type:       types.StringValue(source.GetType()),
			Properties: types.StringValue(string(sourceJson)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
