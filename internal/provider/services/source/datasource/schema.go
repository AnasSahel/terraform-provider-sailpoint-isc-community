// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package source_datasource

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSourceDataSourceSchema returns the schema for the single source data source.
func GetSourceDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Use this data source to get information about a specific SailPoint Identity Security Cloud (ISC) Source.",
		MarkdownDescription: "Use this data source to get information about a specific SailPoint Identity Security Cloud (ISC) Source.",

		Attributes: map[string]schema.Attribute{
			// Core identifiers - required for lookup
			"id": schema.StringAttribute{
				Description:         "The unique ID of the source.",
				MarkdownDescription: "The unique ID of the source.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Human-readable display name of the source.",
				MarkdownDescription: "Human-readable display name of the source.",
				Optional:            true,
				Computed:            true,
			},

			// Required attributes
			"description": schema.StringAttribute{
				Description:         "Human-readable description of the source.",
				MarkdownDescription: "Human-readable description of the source.",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				Description:         "Reference to identity object who owns the source.",
				MarkdownDescription: "Reference to identity object who owns the source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         "Type of object being referenced.",
						MarkdownDescription: "Type of object being referenced.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Owner identity's ID.",
						MarkdownDescription: "Owner identity's ID.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Owner identity's human-readable display name.",
						MarkdownDescription: "Owner identity's human-readable display name.",
						Computed:            true,
					},
				},
			},
			"connector": schema.StringAttribute{
				Description:         "Connector type used to connect to the source.",
				MarkdownDescription: "Connector type used to connect to the source.",
				Computed:            true,
			},

			// Core attributes
			"type": schema.StringAttribute{
				Description:         "Type of the source (e.g., 'SOURCE', 'TARGET').",
				MarkdownDescription: "Type of the source (e.g., 'SOURCE', 'TARGET').",
				Computed:            true,
			},
			"connector_class": schema.StringAttribute{
				Description:         "The connector class name.",
				MarkdownDescription: "The connector class name.",
				Computed:            true,
			},
			"connection_type": schema.StringAttribute{
				Description:         "Type of connection (e.g., 'file', 'account').",
				MarkdownDescription: "Type of connection (e.g., 'file', 'account').",
				Computed:            true,
			},
			"authoritative": schema.BoolAttribute{
				Description:         "Whether this source is authoritative.",
				MarkdownDescription: "Whether this source is authoritative.",
				Computed:            true,
			},
			"cluster": schema.SingleNestedAttribute{
				Description:         "Reference to the cluster where the source resides.",
				MarkdownDescription: "Reference to the cluster where the source resides.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         "Type of object being referenced (always 'CLUSTER').",
						MarkdownDescription: "Type of object being referenced (always 'CLUSTER').",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Cluster's unique identifier.",
						MarkdownDescription: "Cluster's unique identifier.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Cluster's human-readable display name.",
						MarkdownDescription: "Cluster's human-readable display name.",
						Computed:            true,
					},
				},
			},

			// Configuration attributes
			"connector_attributes": schema.StringAttribute{
				Description:         "JSON representation of connector-specific attributes.",
				MarkdownDescription: "JSON representation of connector-specific attributes.",
				Computed:            true,
			},
			"delete_threshold": schema.Int64Attribute{
				Description:         "Number of accounts to delete at once (caps mass deletion).",
				MarkdownDescription: "Number of accounts to delete at once (caps mass deletion).",
				Computed:            true,
			},
			"features": schema.ListAttribute{
				Description:         "List of features supported by this source.",
				MarkdownDescription: "List of features supported by this source.",
				ElementType:         types.StringType,
				Computed:            true,
			},

			// Management attributes
			"management_workgroup": schema.SingleNestedAttribute{
				Description:         "Reference to the management workgroup for this source.",
				MarkdownDescription: "Management workgroup reference that controls who can manage this source.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         "Type of object being referenced (typically 'GOVERNANCE_GROUP').",
						MarkdownDescription: "Type of object being referenced (typically 'GOVERNANCE_GROUP').",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Management workgroup's unique identifier.",
						MarkdownDescription: "Management workgroup's unique identifier.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Management workgroup's human-readable display name.",
						MarkdownDescription: "Management workgroup's human-readable display name.",
						Computed:            true,
					},
				},
			},

			// Correlation & Rules
			"account_correlation_config": schema.StringAttribute{
				Description:         "JSON representation of account correlation configuration reference.",
				MarkdownDescription: "JSON representation of account correlation configuration reference.",
				Computed:            true,
			},
			"account_correlation_rule": schema.StringAttribute{
				Description:         "JSON representation of account correlation rule reference.",
				MarkdownDescription: "JSON representation of account correlation rule reference.",
				Computed:            true,
			},
			"manager_correlation_rule": schema.StringAttribute{
				Description:         "JSON representation of manager correlation rule reference.",
				MarkdownDescription: "JSON representation of manager correlation rule reference.",
				Computed:            true,
			},
			"manager_correlation_mapping": schema.StringAttribute{
				Description:         "JSON representation of manager correlation mapping configuration.",
				MarkdownDescription: "JSON representation of manager correlation mapping configuration.",
				Computed:            true,
			},

			// Provisioning
			"before_provisioning_rule": schema.StringAttribute{
				Description:         "JSON representation of before provisioning rule reference.",
				MarkdownDescription: "JSON representation of before provisioning rule reference.",
				Computed:            true,
			},
			"password_policies": schema.StringAttribute{
				Description:         "JSON representation of password policies list.",
				MarkdownDescription: "JSON representation of password policies list.",
				Computed:            true,
			},

			// Status & Metadata (Computed)
			"healthy": schema.BoolAttribute{
				Description:         "Whether the source is healthy.",
				MarkdownDescription: "Whether the source is healthy.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				Description:         "Status of the source (e.g., 'SOURCE_STATE_HEALTHY').",
				MarkdownDescription: "Status of the source (e.g., 'SOURCE_STATE_HEALTHY').",
				Computed:            true,
			},
			"since": schema.StringAttribute{
				Description:         "Timestamp when the source entered its current status.",
				MarkdownDescription: "Timestamp when the source entered its current status.",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				Description:         "Timestamp when the source was created.",
				MarkdownDescription: "Timestamp when the source was created.",
				Computed:            true,
			},
			"modified": schema.StringAttribute{
				Description:         "Timestamp when the source was last modified.",
				MarkdownDescription: "Timestamp when the source was last modified.",
				Computed:            true,
			},
			"connector_id": schema.StringAttribute{
				Description:         "The unique ID of the connector.",
				MarkdownDescription: "The unique ID of the connector.",
				Computed:            true,
			},
			"connector_name": schema.StringAttribute{
				Description:         "The display name of the connector.",
				MarkdownDescription: "The display name of the connector.",
				Computed:            true,
			},
			"schemas": schema.StringAttribute{
				Description:         "JSON representation of schemas associated with this source.",
				MarkdownDescription: "JSON representation of schemas associated with this source.",
				Computed:            true,
			},

			// Special Parameters
			"credential_provider_enabled": schema.BoolAttribute{
				Description:         "Whether credential provider is enabled for this source.",
				MarkdownDescription: "Whether credential provider is enabled for this source.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				Description:         "The category of the source.",
				MarkdownDescription: "The category of the source.",
				Computed:            true,
			},
		},
	}
}

// GetSourcesDataSourceSchema returns the schema for the sources list data source.
func GetSourcesDataSourceSchema() schema.Schema {
	return schema.Schema{
		Description:         "Use this data source to get a list of SailPoint Identity Security Cloud (ISC) Sources with optional filtering.",
		MarkdownDescription: "Use this data source to get a list of SailPoint Identity Security Cloud (ISC) Sources with optional filtering.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for this data source.",
				MarkdownDescription: "Unique identifier for this data source.",
				Computed:            true,
			},

			// Filter parameters
			"filters": schema.StringAttribute{
				Description:         "JSON string representing search filters to apply.",
				MarkdownDescription: "JSON string representing search filters to apply.",
				Optional:            true,
			},
			"sorters": schema.StringAttribute{
				Description:         "JSON string representing sort criteria.",
				MarkdownDescription: "JSON string representing sort criteria.",
				Optional:            true,
			},
			"limit": schema.Int32Attribute{
				Description:         "Maximum number of results to return.",
				MarkdownDescription: "Maximum number of results to return.",
				Optional:            true,
			},
			"offset": schema.Int32Attribute{
				Description:         "Offset for pagination.",
				MarkdownDescription: "Offset for pagination.",
				Optional:            true,
			},
			"include_count": schema.BoolAttribute{
				Description:         "Whether to include count in the response.",
				MarkdownDescription: "Whether to include count in the response.",
				Optional:            true,
			},

			// Pagination parameters
			"paginate_all": schema.BoolAttribute{
				Description:         "Whether to paginate through all results automatically.",
				MarkdownDescription: "Whether to paginate through all results automatically.",
				Optional:            true,
			},
			"max_results": schema.Int32Attribute{
				Description:         "Maximum number of results to return when paginating all.",
				MarkdownDescription: "Maximum number of results to return when paginating all.",
				Optional:            true,
			},
			"page_size": schema.Int32Attribute{
				Description:         "Size of each page when paginating.",
				MarkdownDescription: "Size of each page when paginating.",
				Optional:            true,
			},

			// Results
			"sources": schema.ListNestedAttribute{
				Description:         "List of sources matching the criteria.",
				MarkdownDescription: "List of sources matching the criteria.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "The unique ID of the source.",
							MarkdownDescription: "The unique ID of the source.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Human-readable display name of the source.",
							MarkdownDescription: "Human-readable display name of the source.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Human-readable description of the source.",
							MarkdownDescription: "Human-readable description of the source.",
							Computed:            true,
						},
						"owner": schema.SingleNestedAttribute{
							Description:         "Reference to identity object who owns the source.",
							MarkdownDescription: "Reference to identity object who owns the source.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description:         "Type of object being referenced.",
									MarkdownDescription: "Type of object being referenced.",
									Computed:            true,
								},
								"id": schema.StringAttribute{
									Description:         "Owner identity's ID.",
									MarkdownDescription: "Owner identity's ID.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Owner identity's human-readable display name.",
									MarkdownDescription: "Owner identity's human-readable display name.",
									Computed:            true,
								},
							},
						},
						"connector": schema.StringAttribute{
							Description:         "Connector type used to connect to the source.",
							MarkdownDescription: "Connector type used to connect to the source.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Type of the source.",
							MarkdownDescription: "Type of the source.",
							Computed:            true,
						},
						"connector_class": schema.StringAttribute{
							Description:         "The connector class name.",
							MarkdownDescription: "The connector class name.",
							Computed:            true,
						},
						"connection_type": schema.StringAttribute{
							Description:         "Type of connection.",
							MarkdownDescription: "Type of connection.",
							Computed:            true,
						},
						"authoritative": schema.BoolAttribute{
							Description:         "Whether this source is authoritative.",
							MarkdownDescription: "Whether this source is authoritative.",
							Computed:            true,
						},
						"cluster": schema.SingleNestedAttribute{
							Description:         "Reference to the cluster where the source resides.",
							MarkdownDescription: "Reference to the cluster where the source resides.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description:         "Type of object being referenced (always 'CLUSTER').",
									MarkdownDescription: "Type of object being referenced (always 'CLUSTER').",
									Computed:            true,
								},
								"id": schema.StringAttribute{
									Description:         "Cluster's unique identifier.",
									MarkdownDescription: "Cluster's unique identifier.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Cluster's human-readable display name.",
									MarkdownDescription: "Cluster's human-readable display name.",
									Computed:            true,
								},
							},
						},
						"connector_attributes": schema.StringAttribute{
							Description:         "JSON representation of connector-specific attributes.",
							MarkdownDescription: "JSON representation of connector-specific attributes.",
							Computed:            true,
						},
						"delete_threshold": schema.Int64Attribute{
							Description:         "Number of accounts to delete at once.",
							MarkdownDescription: "Number of accounts to delete at once.",
							Computed:            true,
						},
						"features": schema.ListAttribute{
							Description:         "List of features supported by this source.",
							MarkdownDescription: "List of features supported by this source.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"healthy": schema.BoolAttribute{
							Description:         "Whether the source is healthy.",
							MarkdownDescription: "Whether the source is healthy.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							Description:         "Status of the source.",
							MarkdownDescription: "Status of the source.",
							Computed:            true,
						},
						"since": schema.StringAttribute{
							Description:         "Timestamp when the source entered its current status.",
							MarkdownDescription: "Timestamp when the source entered its current status.",
							Computed:            true,
						},
						"created": schema.StringAttribute{
							Description:         "Timestamp when the source was created.",
							MarkdownDescription: "Timestamp when the source was created.",
							Computed:            true,
						},
						"modified": schema.StringAttribute{
							Description:         "Timestamp when the source was last modified.",
							MarkdownDescription: "Timestamp when the source was last modified.",
							Computed:            true,
						},
						"connector_id": schema.StringAttribute{
							Description:         "The unique ID of the connector.",
							MarkdownDescription: "The unique ID of the connector.",
							Computed:            true,
						},
						"connector_name": schema.StringAttribute{
							Description:         "The display name of the connector.",
							MarkdownDescription: "The display name of the connector.",
							Computed:            true,
						},
						"credential_provider_enabled": schema.BoolAttribute{
							Description:         "Whether credential provider is enabled.",
							MarkdownDescription: "Whether credential provider is enabled.",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							Description:         "The category of the source.",
							MarkdownDescription: "The category of the source.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}
