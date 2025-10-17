package datasources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/sailpoint_sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &formDefinitionDataSource{}
	_ datasource.DataSourceWithConfigure = &formDefinitionDataSource{}
)

type formDefinitionDataSource struct {
	client *sailpoint_sdk.Client
}

func NewFormDefinitionDataSource() datasource.DataSource {
	return &formDefinitionDataSource{}
}

func (d *formDefinitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_form_definition"
}

func (d *formDefinitionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring form definition datasource")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint_sdk.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
	tflog.Info(ctx, "Configured form definition datasource", map[string]any{"success": true})
}

func (d *formDefinitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = formDefinitionSchema()
}

func (d *formDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading form definition datasource")
	var data models.FormDefinitionModel

	// Read the configuration data into the model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Form definition datasource config data", map[string]any{"id": data.Id.ValueString()})

	tflog.Debug(ctx, "Fetching form definition by ID", map[string]any{"id": data.Id.ValueString()})
	fd, err := d.client.FormDefinitionApi.GetFormDefinitionById(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve form definition",
			fmt.Sprintf("Failed to retrieve form definition with ID %q: %v", data.Id.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Successfully fetched form definition", map[string]any{
		"id":              data.Id.ValueString(),
		"form_definition": fd,
	})

	diags = data.FromSailPointModel(ctx, fd, models.ConversionOptions[models.FormDefinitionModel]{Plan: nil})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func formDefinitionSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the form definition to retrieve.",
				MarkdownDescription: "The ID of the form definition to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the form definition.",
				MarkdownDescription: "The name of the form definition.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "The description of the form definition.",
				MarkdownDescription: "The description of the form definition.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				Description:         "The creation timestamp of the form definition.",
				MarkdownDescription: "The creation timestamp of the form definition.",
			},
			"modified": schema.StringAttribute{
				Computed:            true,
				Description:         "The last modified timestamp of the form definition.",
				MarkdownDescription: "The last modified timestamp of the form definition.",
			},
			"owner": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "The owner of the form definition.",
				MarkdownDescription: "The owner of the form definition.",
				Attributes:          resourceRefSchema(),
			},
			"used_by": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of entities using this form definition.",
				MarkdownDescription: "List of entities using this form definition.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: resourceRefSchema(),
				},
			},
			// "form_input": schema.ListNestedAttribute{
			// 	Computed:            true,
			// 	Description:         "List of form inputs in this form definition.",
			// 	MarkdownDescription: "List of form inputs in this form definition.",
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: formInputSchema(),
			// 	},
			// },
		},
	}
}

func resourceRefSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Computed:            true,
			Description:         "The type of the owner.",
			MarkdownDescription: "The type of the owner.",
		},
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "The ID of the owner.",
			MarkdownDescription: "The ID of the owner.",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			Description:         "The name of the owner.",
			MarkdownDescription: "The name of the owner.",
		},
	}
}

func formInputSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "The ID of the form input.",
			MarkdownDescription: "The ID of the form input.",
		},
		"type": schema.StringAttribute{
			Computed:            true,
			Description:         "The type of the form input.",
			MarkdownDescription: "The type of the form input.",
		},
		"label": schema.StringAttribute{
			Computed:            true,
			Description:         "The label of the form input.",
			MarkdownDescription: "The label of the form input.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			Description:         "The description of the form input.",
			MarkdownDescription: "The description of the form input.",
		},
	}
}
