// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type IdentityProfileSchemaBuilder struct{}

var (
	_ SchemaBuilder = &IdentityProfileSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for IdentityProfile resource.
func (sb *IdentityProfileSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]resource_schema.Attribute{
		"id": resource_schema.StringAttribute{
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": resource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Required:            true,
		},
		"created": resource_schema.StringAttribute{
			Description:         desc["created"].description,
			MarkdownDescription: desc["created"].markdown,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"modified": resource_schema.StringAttribute{
			Description:         desc["modified"].description,
			MarkdownDescription: desc["modified"].markdown,
			Computed:            true,
		},
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
		},
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Optional:            true,
			Computed:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description:         desc["owner.type"].description,
					MarkdownDescription: desc["owner.type"].markdown,
					Required:            true,
				},
				"id": resource_schema.StringAttribute{
					Description:         desc["owner.id"].description,
					MarkdownDescription: desc["owner.id"].markdown,
					Required:            true,
				},
				"name": resource_schema.StringAttribute{
					Description:         desc["owner.name"].description,
					MarkdownDescription: desc["owner.name"].markdown,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"priority": resource_schema.Int64Attribute{
			Description:         desc["priority"].description,
			MarkdownDescription: desc["priority"].markdown,
			Optional:            true,
			Computed:            true,
		},
		"authoritative_source": resource_schema.SingleNestedAttribute{
			Description:         desc["authoritative_source"].description,
			MarkdownDescription: desc["authoritative_source"].markdown,
			Required:            true,
			Attributes: map[string]resource_schema.Attribute{
				"type": resource_schema.StringAttribute{
					Description:         desc["authoritative_source.type"].description,
					MarkdownDescription: desc["authoritative_source.type"].markdown,
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"id": resource_schema.StringAttribute{
					Description:         desc["authoritative_source.id"].description,
					MarkdownDescription: desc["authoritative_source.id"].markdown,
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"name": resource_schema.StringAttribute{
					Description:         desc["authoritative_source.name"].description,
					MarkdownDescription: desc["authoritative_source.name"].markdown,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"identity_refresh_required": resource_schema.BoolAttribute{
			Description:         desc["identity_refresh_required"].description,
			MarkdownDescription: desc["identity_refresh_required"].markdown,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"identity_count": resource_schema.Int64Attribute{
			Description:         desc["identity_count"].description,
			MarkdownDescription: desc["identity_count"].markdown,
			Computed:            true,
		},
		"identity_attribute_config": resource_schema.SingleNestedAttribute{
			Description:         desc["identity_attribute_config"].description,
			MarkdownDescription: desc["identity_attribute_config"].markdown,
			Optional:            true,
			Computed:            true,
			Attributes: map[string]resource_schema.Attribute{
				"enabled": resource_schema.BoolAttribute{
					Description:         desc["identity_attribute_config.enabled"].description,
					MarkdownDescription: desc["identity_attribute_config.enabled"].markdown,
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"attribute_transforms": resource_schema.ListNestedAttribute{
					Description:         desc["identity_attribute_config.attribute_transforms"].description,
					MarkdownDescription: desc["identity_attribute_config.attribute_transforms"].markdown,
					Optional:            true,
					NestedObject: resource_schema.NestedAttributeObject{
						Attributes: map[string]resource_schema.Attribute{
							"identity_attribute_name": resource_schema.StringAttribute{
								Description:         desc["identity_attribute_config.attribute_transforms.identity_attribute_name"].description,
								MarkdownDescription: desc["identity_attribute_config.attribute_transforms.identity_attribute_name"].markdown,
								Required:            true,
							},
							"transform_definition": resource_schema.StringAttribute{
								Description:         desc["identity_attribute_config.attribute_transforms.transform_definition"].description,
								MarkdownDescription: desc["identity_attribute_config.attribute_transforms.transform_definition"].markdown,
								Optional:            true,
								CustomType:          jsontypes.NormalizedType{},
							},
						},
					},
				},
			},
		},
		"identity_exception_report_reference": resource_schema.SingleNestedAttribute{
			Description:         desc["identity_exception_report_reference"].description,
			MarkdownDescription: desc["identity_exception_report_reference"].markdown,
			Computed:            true,
			Attributes: map[string]resource_schema.Attribute{
				"task_result_id": resource_schema.StringAttribute{
					Description:         desc["identity_exception_report_reference.task_result_id"].description,
					MarkdownDescription: desc["identity_exception_report_reference.task_result_id"].markdown,
					Computed:            true,
				},
				"report_name": resource_schema.StringAttribute{
					Description:         desc["identity_exception_report_reference.report_name"].description,
					MarkdownDescription: desc["identity_exception_report_reference.report_name"].markdown,
					Computed:            true,
				},
			},
		},
		"has_time_based_attr": resource_schema.BoolAttribute{
			Description:         desc["has_time_based_attr"].description,
			MarkdownDescription: desc["has_time_based_attr"].markdown,
			Optional:            true,
			Computed:            true,
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for IdentityProfile data source.
func (sb *IdentityProfileSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := sb.fieldDescriptions()

	return map[string]datasource_schema.Attribute{
		"id": datasource_schema.StringAttribute{
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
			Required:            true,
		},
		"name": datasource_schema.StringAttribute{
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
			Computed:            true,
		},
		"created": datasource_schema.StringAttribute{
			Description:         desc["created"].description,
			MarkdownDescription: desc["created"].markdown,
			Computed:            true,
		},
		"modified": datasource_schema.StringAttribute{
			Description:         desc["modified"].description,
			MarkdownDescription: desc["modified"].markdown,
			Computed:            true,
		},
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description:         desc["owner.type"].description,
					MarkdownDescription: desc["owner.type"].markdown,
					Computed:            true,
				},
				"id": datasource_schema.StringAttribute{
					Description:         desc["owner.id"].description,
					MarkdownDescription: desc["owner.id"].markdown,
					Computed:            true,
				},
				"name": datasource_schema.StringAttribute{
					Description:         desc["owner.name"].description,
					MarkdownDescription: desc["owner.name"].markdown,
					Computed:            true,
				},
			},
		},
		"priority": datasource_schema.Int64Attribute{
			Description:         desc["priority"].description,
			MarkdownDescription: desc["priority"].markdown,
			Computed:            true,
		},
		"authoritative_source": datasource_schema.SingleNestedAttribute{
			Description:         desc["authoritative_source"].description,
			MarkdownDescription: desc["authoritative_source"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"type": datasource_schema.StringAttribute{
					Description:         desc["authoritative_source.type"].description,
					MarkdownDescription: desc["authoritative_source.type"].markdown,
					Computed:            true,
				},
				"id": datasource_schema.StringAttribute{
					Description:         desc["authoritative_source.id"].description,
					MarkdownDescription: desc["authoritative_source.id"].markdown,
					Computed:            true,
				},
				"name": datasource_schema.StringAttribute{
					Description:         desc["authoritative_source.name"].description,
					MarkdownDescription: desc["authoritative_source.name"].markdown,
					Computed:            true,
				},
			},
		},
		"identity_refresh_required": datasource_schema.BoolAttribute{
			Description:         desc["identity_refresh_required"].description,
			MarkdownDescription: desc["identity_refresh_required"].markdown,
			Computed:            true,
		},
		"identity_count": datasource_schema.Int64Attribute{
			Description:         desc["identity_count"].description,
			MarkdownDescription: desc["identity_count"].markdown,
			Computed:            true,
		},
		"identity_attribute_config": datasource_schema.SingleNestedAttribute{
			Description:         desc["identity_attribute_config"].description,
			MarkdownDescription: desc["identity_attribute_config"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"enabled": datasource_schema.BoolAttribute{
					Description:         desc["identity_attribute_config.enabled"].description,
					MarkdownDescription: desc["identity_attribute_config.enabled"].markdown,
					Computed:            true,
				},
				"attribute_transforms": datasource_schema.ListNestedAttribute{
					Description:         desc["identity_attribute_config.attribute_transforms"].description,
					MarkdownDescription: desc["identity_attribute_config.attribute_transforms"].markdown,
					Computed:            true,
					NestedObject: datasource_schema.NestedAttributeObject{
						Attributes: map[string]datasource_schema.Attribute{
							"identity_attribute_name": datasource_schema.StringAttribute{
								Description:         desc["identity_attribute_config.attribute_transforms.identity_attribute_name"].description,
								MarkdownDescription: desc["identity_attribute_config.attribute_transforms.identity_attribute_name"].markdown,
								Computed:            true,
							},
							"transform_definition": datasource_schema.StringAttribute{
								Description:         desc["identity_attribute_config.attribute_transforms.transform_definition"].description,
								MarkdownDescription: desc["identity_attribute_config.attribute_transforms.transform_definition"].markdown,
								Computed:            true,
								CustomType:          jsontypes.NormalizedType{},
							},
						},
					},
				},
			},
		},
		"identity_exception_report_reference": datasource_schema.SingleNestedAttribute{
			Description:         desc["identity_exception_report_reference"].description,
			MarkdownDescription: desc["identity_exception_report_reference"].markdown,
			Computed:            true,
			Attributes: map[string]datasource_schema.Attribute{
				"task_result_id": datasource_schema.StringAttribute{
					Description:         desc["identity_exception_report_reference.task_result_id"].description,
					MarkdownDescription: desc["identity_exception_report_reference.task_result_id"].markdown,
					Computed:            true,
				},
				"report_name": datasource_schema.StringAttribute{
					Description:         desc["identity_exception_report_reference.report_name"].description,
					MarkdownDescription: desc["identity_exception_report_reference.report_name"].markdown,
					Computed:            true,
				},
			},
		},
		"has_time_based_attr": datasource_schema.BoolAttribute{
			Description:         desc["has_time_based_attr"].description,
			MarkdownDescription: desc["has_time_based_attr"].markdown,
			Computed:            true,
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *IdentityProfileSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "The unique identifier of the identity profile.",
			markdown:    "The unique identifier of the identity profile. This is a system-generated ID.",
		},
		"name": {
			description: "The name of the identity profile.",
			markdown:    "The name of the identity profile. This is the human-readable name used to identify the profile.",
		},
		"created": {
			description: "The timestamp when the identity profile was created.",
			markdown:    "The timestamp when the identity profile was created (ISO 8601 format).",
		},
		"modified": {
			description: "The timestamp when the identity profile was last modified.",
			markdown:    "The timestamp when the identity profile was last modified (ISO 8601 format).",
		},
		"description": {
			description: "A description of the identity profile.",
			markdown:    "A description of the identity profile providing additional context about its purpose and configuration.",
		},
		"owner": {
			description: "The owner of the identity profile.",
			markdown:    "The owner of the identity profile. This is typically an identity that has administrative control over the profile.",
		},
		"owner.type": {
			description: "The type of the owner object.",
			markdown:    "The type of the owner object. Must be `IDENTITY`.",
		},
		"owner.id": {
			description: "The ID of the owner.",
			markdown:    "The unique identifier of the owner identity.",
		},
		"owner.name": {
			description: "The name of the owner.",
			markdown:    "The display name of the owner identity.",
		},
		"priority": {
			description: "The priority of the identity profile.",
			markdown:    "The priority of the identity profile. Lower numbers indicate higher priority. This affects which profile takes precedence when an identity matches multiple profiles.",
		},
		"authoritative_source": {
			description: "The authoritative source for the identity profile.",
			markdown:    "The authoritative source that provides the primary identity data for this profile. This is a required field.",
		},
		"authoritative_source.type": {
			description: "The type of the authoritative source.",
			markdown:    "The type of the authoritative source. Must be `SOURCE`.",
		},
		"authoritative_source.id": {
			description: "The ID of the authoritative source.",
			markdown:    "The unique identifier of the authoritative source.",
		},
		"authoritative_source.name": {
			description: "The name of the authoritative source.",
			markdown:    "The display name of the authoritative source.",
		},
		"identity_refresh_required": {
			description: "Indicates whether an identity refresh is required.",
			markdown:    "Set to `true` if an identity refresh is necessary. You would typically want to trigger an identity refresh when a change has been made on the source. Defaults to `false`.",
		},
		"identity_count": {
			description: "The number of identities belonging to the profile.",
			markdown:    "The number of identities currently associated with this profile. This is a read-only computed value.",
		},
		"identity_attribute_config": {
			description: "Configuration for identity attribute mappings.",
			markdown:    "Configuration that defines how identity attributes are mapped and transformed. This controls the attribute mapping process during identity refresh.",
		},
		"identity_attribute_config.enabled": {
			description: "Indicates whether attribute mapping is enabled.",
			markdown:    "Backend will only promote values if the profile/mapping is enabled. Defaults to `false`.",
		},
		"identity_attribute_config.attribute_transforms": {
			description: "List of attribute transform configurations.",
			markdown:    "List of transforms that define how to generate or collect data for each identity attribute during the identity refresh process.",
		},
		"identity_attribute_config.attribute_transforms.identity_attribute_name": {
			description: "The name of the identity attribute.",
			markdown:    "The name of the identity attribute to which this transform applies (e.g., `email`, `department`).",
		},
		"identity_attribute_config.attribute_transforms.transform_definition": {
			description: "The transform definition for the attribute.",
			markdown:    "The transform definition that specifies how the attribute value should be calculated or derived.",
		},
		"identity_attribute_config.attribute_transforms.transform_definition.type": {
			description: "The type of transform.",
			markdown:    "The type of transform to apply (e.g., `accountAttribute`, `static`, `reference`).",
		},
		"identity_attribute_config.attribute_transforms.transform_definition.attributes": {
			description: "Transform-specific configuration attributes.",
			markdown:    "Transform-specific configuration attributes as a JSON string. The structure varies by transform type.",
		},
		"identity_exception_report_reference": {
			description: "Reference to an identity exception report.",
			markdown:    "Reference to an identity exception report if exceptions occurred during identity processing.",
		},
		"identity_exception_report_reference.task_result_id": {
			description: "The task result ID.",
			markdown:    "The UUID of the task result that generated the exception report.",
		},
		"identity_exception_report_reference.report_name": {
			description: "The name of the report.",
			markdown:    "The name of the exception report.",
		},
		"has_time_based_attr": {
			description: "Indicates if the profile has time-based attributes.",
			markdown:    "Indicates the value of `requiresPeriodicRefresh` attribute for the identity profile. Set to `true` if the profile includes time-based attributes that require periodic refresh. Defaults to `false`.",
		},
	}
}
