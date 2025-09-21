// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetConnectorsDataSourceSchema returns the schema for the connectors data source
func GetConnectorsDataSourceSchema() datasource_schema.Schema {
	return datasource_schema.Schema{
		Description:         "The `sailpoint_connectors` data source allows you to retrieve a list of available connectors in SailPoint ISC with optional filtering.",
		MarkdownDescription: "The `sailpoint_connectors` data source allows you to retrieve a list of available connectors in SailPoint ISC with optional filtering.",
		Attributes: map[string]datasource_schema.Attribute{
			"id": datasource_schema.StringAttribute{
				Description:         "The unique identifier for this data source.",
				MarkdownDescription: "The unique identifier for this data source.",
				Computed:            true,
			},
			"filters": datasource_schema.StringAttribute{
				Description:         "Filter results using the standard syntax. Supported filters: name (sw, co), type (sw, co, eq), directConnect (eq), category (eq), features (ca), labels (ca).",
				MarkdownDescription: "Filter results using the standard syntax. Supported filters: `name` (sw, co), `type` (sw, co, eq), `directConnect` (eq), `category` (eq), `features` (ca), `labels` (ca).",
				Optional:            true,
			},
			"limit": datasource_schema.Int32Attribute{
				Description:         "Max number of results to return. Defaults to 250.",
				MarkdownDescription: "Max number of results to return. Defaults to 250.",
				Optional:            true,
			},
			"offset": datasource_schema.Int32Attribute{
				Description:         "Offset into the full result set for pagination. Defaults to 0.",
				MarkdownDescription: "Offset into the full result set for pagination. Defaults to 0.",
				Optional:            true,
			},
			"include_count": datasource_schema.BoolAttribute{
				Description:         "If true, populate the X-Total-Count response header with the total number of results.",
				MarkdownDescription: "If true, populate the X-Total-Count response header with the total number of results.",
				Optional:            true,
			},
			"locale": datasource_schema.StringAttribute{
				Description:         "The locale to apply to the config. Defaults to 'en' if not specified.",
				MarkdownDescription: "The locale to apply to the config. Defaults to 'en' if not specified.",
				Optional:            true,
			},
			"paginate_all": datasource_schema.BoolAttribute{
				Description:         "If true, fetch all results using pagination (up to 10,000 records). Overrides limit and offset parameters.",
				MarkdownDescription: "If true, fetch all results using pagination (up to 10,000 records). Overrides `limit` and `offset` parameters.",
				Optional:            true,
			},
			"max_results": datasource_schema.Int32Attribute{
				Description:         "Maximum number of results to fetch when paginate_all is true. Defaults to 10,000. Only applies when paginate_all is true.",
				MarkdownDescription: "Maximum number of results to fetch when `paginate_all` is true. Defaults to 10,000. Only applies when `paginate_all` is true.",
				Optional:            true,
			},
			"page_size": datasource_schema.Int32Attribute{
				Description:         "Number of results per page when using pagination. Defaults to 250. Only applies when paginate_all is true.",
				MarkdownDescription: "Number of results per page when using pagination. Defaults to 250. Only applies when `paginate_all` is true.",
				Optional:            true,
			},
			"connectors": datasource_schema.ListNestedAttribute{
				Description:         "List of connectors matching the specified criteria.",
				MarkdownDescription: "List of connectors matching the specified criteria.",
				Computed:            true,
				NestedObject: datasource_schema.NestedAttributeObject{
					Attributes: map[string]datasource_schema.Attribute{
						"id": datasource_schema.StringAttribute{
							Description:         "The unique identifier of the connector (derived from script name).",
							MarkdownDescription: "The unique identifier of the connector (derived from script name).",
							Computed:            true,
						},
						"name": datasource_schema.StringAttribute{
							Description:         "The display name of the connector.",
							MarkdownDescription: "The display name of the connector.",
							Computed:            true,
						},
						"type": datasource_schema.StringAttribute{
							Description:         "The connector type (e.g., 'active-directory', 'workday').",
							MarkdownDescription: "The connector type (e.g., 'active-directory', 'workday').",
							Computed:            true,
						},
						"script_name": datasource_schema.StringAttribute{
							Description:         "The script name (unique identifier) of the connector.",
							MarkdownDescription: "The script name (unique identifier) of the connector.",
							Computed:            true,
						},
						"class_name": datasource_schema.StringAttribute{
							Description:         "The Java class name that implements the connector.",
							MarkdownDescription: "The Java class name that implements the connector.",
							Computed:            true,
						},
						"direct_connect": datasource_schema.BoolAttribute{
							Description:         "Whether the connector supports direct connection without a VA.",
							MarkdownDescription: "Whether the connector supports direct connection without a VA.",
							Computed:            true,
						},
						"status": datasource_schema.StringAttribute{
							Description:         "The status of the connector (e.g., 'RELEASED', 'BETA').",
							MarkdownDescription: "The status of the connector (e.g., 'RELEASED', 'BETA').",
							Computed:            true,
						},
						"category": datasource_schema.StringAttribute{
							Description:         "The category of the connector (not available in list API).",
							MarkdownDescription: "The category of the connector (not available in list API).",
							Computed:            true,
						},
						"features": datasource_schema.ListAttribute{
							Description:         "List of features supported by the connector.",
							MarkdownDescription: "List of features supported by the connector.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"labels": datasource_schema.ListAttribute{
							Description:         "List of labels associated with the connector (not available in list API).",
							MarkdownDescription: "List of labels associated with the connector (not available in list API).",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

// GetConnectorResourceSchema returns the schema for the connector resource
func GetConnectorResourceSchema() resource_schema.Schema {
	return resource_schema.Schema{
		Description:         "The `sailpoint_connector` resource allows you to manage custom connectors in SailPoint ISC.",
		MarkdownDescription: "The `sailpoint_connector` resource allows you to manage custom connectors in SailPoint ISC.",
		Attributes: map[string]resource_schema.Attribute{
			"id": resource_schema.StringAttribute{
				Description:         "The unique identifier of the connector (script name).",
				MarkdownDescription: "The unique identifier of the connector (script name).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": resource_schema.StringAttribute{
				Description:         "The display name of the connector.",
				MarkdownDescription: "The display name of the connector.",
				Required:            true,
			},
			"type": resource_schema.StringAttribute{
				Description:         "The connector type (e.g., 'active-directory', 'workday').",
				MarkdownDescription: "The connector type (e.g., 'active-directory', 'workday').",
				Required:            true,
			},
			"class_name": resource_schema.StringAttribute{
				Description:         "The Java class name that implements the connector.",
				MarkdownDescription: "The Java class name that implements the connector.",
				Optional:            true,
			},
			"script_name": resource_schema.StringAttribute{
				Description:         "The script name (unique identifier) of the connector. If not provided, will be generated from the name.",
				MarkdownDescription: "The script name (unique identifier) of the connector. If not provided, will be generated from the name.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_xml": resource_schema.StringAttribute{
				Description:         "The connector application XML configuration.",
				MarkdownDescription: "The connector application XML configuration.",
				Optional:            true,
			},
			"correlation_config_xml": resource_schema.StringAttribute{
				Description:         "The connector correlation config XML.",
				MarkdownDescription: "The connector correlation config XML.",
				Optional:            true,
			},
			"source_config_xml": resource_schema.StringAttribute{
				Description:         "The connector source config XML.",
				MarkdownDescription: "The connector source config XML.",
				Optional:            true,
			},
			"source_config": resource_schema.StringAttribute{
				Description:         "The connector source config (JSON format).",
				MarkdownDescription: "The connector source config (JSON format).",
				Optional:            true,
			},
			"s3_location": resource_schema.StringAttribute{
				Description:         "The storage path key for this connector.",
				MarkdownDescription: "The storage path key for this connector.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uploaded_files": resource_schema.ListAttribute{
				Description:         "List of uploaded files supported by the connector.",
				MarkdownDescription: "List of uploaded files supported by the connector.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"file_upload": resource_schema.BoolAttribute{
				Description:         "Whether the connector supports file upload.",
				MarkdownDescription: "Whether the connector supports file upload.",
				Optional:            true,
			},
			"direct_connect": resource_schema.BoolAttribute{
				Description:         "Whether the connector supports direct connection without a VA.",
				MarkdownDescription: "Whether the connector supports direct connection without a VA.",
				Optional:            true,
			},
			"translation_properties": resource_schema.StringAttribute{
				Description:         "Translation attributes by locale key (JSON format).",
				MarkdownDescription: "Translation attributes by locale key (JSON format).",
				Optional:            true,
			},
			"connector_metadata": resource_schema.StringAttribute{
				Description:         "Metadata pertinent to the UI to be used (JSON format).",
				MarkdownDescription: "Metadata pertinent to the UI to be used (JSON format).",
				Optional:            true,
			},
			"status": resource_schema.StringAttribute{
				Description:         "The status of the connector (e.g., 'RELEASED', 'BETA').",
				MarkdownDescription: "The status of the connector (e.g., 'RELEASED', 'BETA').",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
