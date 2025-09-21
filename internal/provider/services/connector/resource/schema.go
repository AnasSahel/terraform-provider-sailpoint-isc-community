// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_resource

import (
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetConnectorResourceSchema returns the schema for the connector resource.
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
				Description:         "The display name of the connector. Changes require resource replacement.",
				MarkdownDescription: "The display name of the connector. Changes require resource replacement.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": resource_schema.StringAttribute{
				Description:         "The connector type. If not specified will be defaulted to 'custom ' + name. Changes require resource replacement.",
				MarkdownDescription: "The connector type. If not specified will be defaulted to 'custom ' + name. Changes require resource replacement.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"class_name": resource_schema.StringAttribute{
				Description:         "The Java class name that implements the connector. Changes require resource replacement.",
				MarkdownDescription: "The Java class name that implements the connector. Changes require resource replacement.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"correlation_config_xml": resource_schema.StringAttribute{
				Description:         "The connector correlation config XML.",
				MarkdownDescription: "The connector correlation config XML.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_config_xml": resource_schema.StringAttribute{
				Description:         "The connector source config XML.",
				MarkdownDescription: "The connector source config XML.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_config": resource_schema.StringAttribute{
				Description:         "The connector source config (JSON format).",
				MarkdownDescription: "The connector source config (JSON format).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"file_upload": resource_schema.BoolAttribute{
				Description:         "Whether the connector supports file upload.",
				MarkdownDescription: "Whether the connector supports file upload.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"direct_connect": resource_schema.BoolAttribute{
				Description:         "Whether the connector supports direct connection without a VA. Changes require resource replacement.",
				MarkdownDescription: "Whether the connector supports direct connection without a VA. Changes require resource replacement.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"translation_properties": resource_schema.StringAttribute{
				Description:         "Translation attributes by locale key (JSON format).",
				MarkdownDescription: "Translation attributes by locale key (JSON format).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connector_metadata": resource_schema.StringAttribute{
				Description:         "Metadata pertinent to the UI to be used (JSON format).",
				MarkdownDescription: "Metadata pertinent to the UI to be used (JSON format).",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": resource_schema.StringAttribute{
				Description:         "The status of the connector. Allowed values: 'RELEASED', 'DEVELOPMENT', 'DEMO', 'DEPRECATED'. Can only be set during creation - changes require resource replacement.",
				MarkdownDescription: "The status of the connector. Allowed values: `RELEASED`, `DEVELOPMENT`, `DEMO`, `DEPRECATED`. Can only be set during creation - changes require resource replacement.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
