// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceSchemaDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceSchemaDataSource{}
)

type sourceSchemaDataSource struct {
	client *client.Client
}

func NewSourceSchemaDataSource() datasource.DataSource {
	return &sourceSchemaDataSource{}
}

func (d *sourceSchemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_schema"
}

func (d *sourceSchemaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source schema data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *sourceSchemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a single SailPoint source schema.",
		MarkdownDescription: "Retrieves a single SailPoint source schema. " +
			"Use `include_types` or `include_names` to filter the schemas returned by the API. " +
			"The data source returns the first schema from the filtered results. " +
			"Schemas are created automatically when a source is created.",
		Attributes: map[string]schema.Attribute{
			// Input parameters
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the source to retrieve the schema from.",
				Required:            true,
			},
			"include_types": schema.StringAttribute{
				MarkdownDescription: "If set to `group`, only group schemas are returned. If set to `user`, only user schemas are returned.",
				Optional:            true,
			},
			"include_names": schema.StringAttribute{
				MarkdownDescription: "A comma-separated list of schema names to filter results (e.g., `account`, `group`).",
				Optional:            true,
			},

			// Output attributes
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the schema.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the schema (e.g., `account`, `group`).",
				Computed:            true,
			},
			"native_object_type": schema.StringAttribute{
				MarkdownDescription: "The name of the object type on the native system that the schema represents (e.g., `User`, `Group`).",
				Computed:            true,
			},
			"identity_attribute": schema.StringAttribute{
				MarkdownDescription: "The name of the attribute used to calculate the unique identifier for an object in the schema.",
				Computed:            true,
			},
			"display_attribute": schema.StringAttribute{
				MarkdownDescription: "The name of the attribute used to calculate the display value for an object in the schema.",
				Computed:            true,
			},
			"hierarchy_attribute": schema.StringAttribute{
				MarkdownDescription: "The name of the attribute whose values represent other objects in a hierarchy. Only relevant to group schemas.",
				Computed:            true,
			},
			"include_permissions": schema.BoolAttribute{
				MarkdownDescription: "Flag indicating whether to include permissions with the object data when aggregating the schema.",
				Computed:            true,
			},
			"features": schema.ListAttribute{
				MarkdownDescription: "Optional features supported by the source.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"configuration": schema.StringAttribute{
				MarkdownDescription: "Extra configuration data for the schema as a JSON object.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "The attribute definitions for the schema.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the attribute.",
							Computed:            true,
						},
						"native_name": schema.StringAttribute{
							MarkdownDescription: "The native name of the attribute on the source system.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the attribute (e.g., `STRING`).",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the attribute.",
							Computed:            true,
						},
						"is_multi": schema.BoolAttribute{
							MarkdownDescription: "Whether the attribute supports multiple values.",
							Computed:            true,
						},
						"is_entitlement": schema.BoolAttribute{
							MarkdownDescription: "Whether the attribute is an entitlement.",
							Computed:            true,
						},
						"is_group": schema.BoolAttribute{
							MarkdownDescription: "Whether the attribute represents a group.",
							Computed:            true,
						},
						"schema": schema.SingleNestedAttribute{
							MarkdownDescription: "A reference to another schema, if applicable (e.g., group membership references the group schema).",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "The type of the schema reference (e.g., `CONNECTOR_SCHEMA`).",
									Computed:            true,
								},
								"id": schema.StringAttribute{
									MarkdownDescription: "The ID of the referenced schema.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the referenced schema.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "The date the schema was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				MarkdownDescription: "The date the schema was last modified.",
				Computed:            true,
			},
		},
	}
}

func (d *sourceSchemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SailPoint Source Schema data source")

	var config sourceSchemaModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := config.SourceID.ValueString()
	includeTypes := ""
	if !config.IncludeTypes.IsNull() && !config.IncludeTypes.IsUnknown() {
		includeTypes = config.IncludeTypes.ValueString()
	}
	includeNames := ""
	if !config.IncludeNames.IsNull() && !config.IncludeNames.IsUnknown() {
		includeNames = config.IncludeNames.ValueString()
	}

	tflog.Debug(ctx, "Fetching source schemas from SailPoint", map[string]any{
		"source_id":     sourceID,
		"include_types": includeTypes,
		"include_names": includeNames,
	})

	schemas, err := d.client.ListSourceSchemas(ctx, sourceID, includeTypes, includeNames)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Schema",
			fmt.Sprintf("Could not read source schemas for source %q: %s", sourceID, err.Error()),
		)
		tflog.Error(ctx, "Failed to read SailPoint Source Schema", map[string]any{
			"source_id": sourceID,
			"error":     err.Error(),
		})
		return
	}

	if schemas == nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Schema",
			"Received nil response from SailPoint API",
		)
		return
	}

	if len(schemas) == 0 {
		resp.Diagnostics.AddError(
			"No Schema Found",
			fmt.Sprintf("No schemas found for source %q with the specified filters (include_types=%q, include_names=%q).",
				sourceID, includeTypes, includeNames),
		)
		return
	}

	// Use the first schema from the results
	var state sourceSchemaModel
	resp.Diagnostics.Append(state.FromSailPointAPI(ctx, schemas[0])...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve input parameters in state
	state.SourceID = config.SourceID
	state.IncludeTypes = config.IncludeTypes
	state.IncludeNames = config.IncludeNames

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully read SailPoint Source Schema data source", map[string]any{
		"source_id":   sourceID,
		"schema_id":   state.ID.ValueString(),
		"schema_name": state.Name.ValueString(),
	})
}
