package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

type transformModel struct {
	Id       types.String `tfsdk:"id"`
	Internal types.Bool   `tfsdk:"internal"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	// Attributes transformAttributesModel `tfsdk:"attributes"`
}

// type transformAttributesModel struct {
// 	RequiresPeriodicRefresh types.Bool `tfsdk:"requires_periodic_refresh"`
// }

type transformsDataSourceModel struct {
	Transforms []transformModel `tfsdk:"transforms"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &transformsDataSource{}
	_ datasource.DataSourceWithConfigure = &transformsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *transformsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// NewTransformsDataSource is a helper function to simplify the provider implementation.
func NewTransformsDataSource() datasource.DataSource {
	return &transformsDataSource{}
}

// transformsDataSource is the data source implementation.
type transformsDataSource struct {
	client *sailpoint.APIClient
}

// Metadata returns the data source type name.
func (d *transformsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transforms"
}

// // Schema defines the schema for the data source.
// func (d *transformsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				Required: true,
// 				// Computed: true,
// 			},
// 		},
// 	}
// 	// resp.Schema = schema.Schema{
// 	// 	Attributes: map[string]schema.Attribute{
// 	// 		"id": schema.StringAttribute{
// 	// 			Computed: true,
// 	// 		},
// 	// 		// "internal": schema.BoolAttribute{
// 	// 		// 	Computed: true,
// 	// 		// },
// 	// 		// "name": schema.StringAttribute{
// 	// 		// 	Computed: true,
// 	// 		// },
// 	// 		// "type": schema.StringAttribute{
// 	// 		// 	Computed: true,
// 	// 		// },
// 	// 		// "attributes": schema.SingleNestedAttribute{
// 	// 		// 	Computed: true,
// 	// 		// 	Attributes: map[string]schema.Attribute{
// 	// 		// 		"requires_periodic_refresh": schema.BoolAttribute{
// 	// 		// 			Required: true,
// 	// 		// 		},
// 	// 		// 	},
// 	// 		// },
// 	// 	},
// 	// }
// }

func (d *transformsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"transforms": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true},
						"internal": schema.BoolAttribute{Computed: true},
						"name":     schema.StringAttribute{Computed: true},
						"type":     schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transformsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state transformsDataSourceModel

	transforms, _, err := d.client.V2025.TransformsAPI.ListTransforms(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read SailPoint Transforms",
			err.Error(),
		)
		return
	}

	for _, transform := range transforms {
		transformState := transformModel{
			Id:       types.StringValue(transform.GetId()),
			Internal: types.BoolValue(transform.GetInternal()),
			Name:     types.StringValue(transform.GetName()),
			Type:     types.StringValue(transform.GetType()),

			// Internal: types.BoolValue(transform.Internal),
			// Name:     types.StringValue(transform.Name),
			// Type:     types.StringValue(transform.Type),
			// Attributes: transformAttributesModel{
			// 	RequiresPeriodicRefresh: types.BoolValue(transform.Attributes["requiresPeriodicRefresh"].(bool)),
			// },
		}

		state.Transforms = append(state.Transforms, transformState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
