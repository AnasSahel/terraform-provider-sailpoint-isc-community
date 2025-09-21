// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetConnectorsDataSourceSchema returns the schema for the connectors data source
func GetConnectorsDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "The `sailpoint_connectors` data source allows you to retrieve a list of available connectors in SailPoint ISC with optional filtering.",
		MarkdownDescription: "The `sailpoint_connectors` data source allows you to retrieve a list of available connectors in SailPoint ISC with optional filtering.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier for this data source.",
				MarkdownDescription: "The unique identifier for this data source.",
				Computed:            true,
			},
			"filters": schema.StringAttribute{
				Description:         "Filter results using the standard syntax. Supported filters: name (sw, co), type (sw, co, eq), directConnect (eq), category (eq), features (ca), labels (ca).",
				MarkdownDescription: "Filter results using the standard syntax. Supported filters: `name` (sw, co), `type` (sw, co, eq), `directConnect` (eq), `category` (eq), `features` (ca), `labels` (ca).",
				Optional:            true,
			},
			"limit": schema.Int32Attribute{
				Description:         "Max number of results to return. Defaults to 250.",
				MarkdownDescription: "Max number of results to return. Defaults to 250.",
				Optional:            true,
			},
			"offset": schema.Int32Attribute{
				Description:         "Offset into the full result set for pagination. Defaults to 0.",
				MarkdownDescription: "Offset into the full result set for pagination. Defaults to 0.",
				Optional:            true,
			},
			"include_count": schema.BoolAttribute{
				Description:         "If true, populate the X-Total-Count response header with the total number of results.",
				MarkdownDescription: "If true, populate the X-Total-Count response header with the total number of results.",
				Optional:            true,
			},
			"locale": schema.StringAttribute{
				Description:         "The locale to apply to the config. Defaults to 'en' if not specified.",
				MarkdownDescription: "The locale to apply to the config. Defaults to 'en' if not specified.",
				Optional:            true,
			},
			"paginate_all": schema.BoolAttribute{
				Description:         "If true, fetch all results using pagination (up to 10,000 records). Overrides limit and offset parameters.",
				MarkdownDescription: "If true, fetch all results using pagination (up to 10,000 records). Overrides `limit` and `offset` parameters.",
				Optional:            true,
			},
			"max_results": schema.Int32Attribute{
				Description:         "Maximum number of results to fetch when paginate_all is true. Defaults to 10,000. Only applies when paginate_all is true.",
				MarkdownDescription: "Maximum number of results to fetch when `paginate_all` is true. Defaults to 10,000. Only applies when `paginate_all` is true.",
				Optional:            true,
			},
			"page_size": schema.Int32Attribute{
				Description:         "Number of results per page when using pagination. Defaults to 250. Only applies when paginate_all is true.",
				MarkdownDescription: "Number of results per page when using pagination. Defaults to 250. Only applies when `paginate_all` is true.",
				Optional:            true,
			},
			"connectors": schema.ListNestedAttribute{
				Description:         "List of connectors matching the specified criteria.",
				MarkdownDescription: "List of connectors matching the specified criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "The unique identifier of the connector (derived from script name).",
							MarkdownDescription: "The unique identifier of the connector (derived from script name).",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "The display name of the connector.",
							MarkdownDescription: "The display name of the connector.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "The connector type (e.g., 'active-directory', 'workday').",
							MarkdownDescription: "The connector type (e.g., 'active-directory', 'workday').",
							Computed:            true,
						},
						"script_name": schema.StringAttribute{
							Description:         "The script name (unique identifier) of the connector.",
							MarkdownDescription: "The script name (unique identifier) of the connector.",
							Computed:            true,
						},
						"class_name": schema.StringAttribute{
							Description:         "The Java class name that implements the connector.",
							MarkdownDescription: "The Java class name that implements the connector.",
							Computed:            true,
						},
						"direct_connect": schema.BoolAttribute{
							Description:         "Whether the connector supports direct connection without a VA.",
							MarkdownDescription: "Whether the connector supports direct connection without a VA.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							Description:         "The status of the connector (e.g., 'RELEASED', 'BETA').",
							MarkdownDescription: "The status of the connector (e.g., 'RELEASED', 'BETA').",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							Description:         "The category of the connector (not available in list API).",
							MarkdownDescription: "The category of the connector (not available in list API).",
							Computed:            true,
						},
						"features": schema.ListAttribute{
							Description:         "List of features supported by the connector.",
							MarkdownDescription: "List of features supported by the connector.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"labels": schema.ListAttribute{
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
