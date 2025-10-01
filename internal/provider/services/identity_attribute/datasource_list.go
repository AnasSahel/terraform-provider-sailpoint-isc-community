package identity_attribute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

var (
	_ datasource.DataSource              = &identityAttributeDataSourceList{}
	_ datasource.DataSourceWithConfigure = &identityAttributeDataSourceList{}
)

type identityAttributeDataSourceList struct {
	client *sailpoint.APIClient
}

func NewIdentityAttributeDataSourceList() datasource.DataSource {
	return &identityAttributeDataSourceList{}
}

func (d *identityAttributeDataSourceList) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *identityAttributeDataSourceList) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + DataSourceListTypeName
}

func (d *identityAttributeDataSourceList) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceListSchema()
}

func (d *identityAttributeDataSourceList) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IdentityAttributeDataSourceListModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	request := d.client.V2025.IdentityAttributesAPI.ListIdentityAttributes(ctx)
	if !data.IncludeSilent.IsNull() {
		request = request.IncludeSilent(data.IncludeSilent.ValueBool())
	}
	if !data.IncludeSystem.IsNull() {
		request = request.IncludeSystem(data.IncludeSystem.ValueBool())
	}
	if !data.SearchableOnly.IsNull() {
		request = request.SearchableOnly(data.SearchableOnly.ValueBool())
	}

	identityAttributes, httpResponse, err := request.Execute()
	if err != nil {
		resp.Diagnostics.AddError(ErrDataSourceReadTitle, fmt.Sprintf(ErrDataSourceReadMsg, "list", err, httpResponse))
		return
	}

	data.Items = []IdentityAttributeModel{}
	for _, identityAttribute := range identityAttributes {
		data.Items = append(data.Items, MapIdentityAttributeToDataSourceModel(identityAttribute))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
