package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SourceSchemaBuilder struct{}

var (
	_ SchemaBuilder = &SourceSchemaBuilder{}
)

func (sb *SourceSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	desc := sb.fieldDescriptions()
	objectRefSchemaBuilder := ObjectRefSchemaBuilder{}

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
		"description": resource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"owner": resource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Required:            true,
			Attributes:          objectRefSchemaBuilder.GetResourceSchema(),
		},
		"cluster": resource_schema.SingleNestedAttribute{
			Description:         desc["cluster"].description,
			MarkdownDescription: desc["cluster"].markdown,
			Optional:            true,
			Attributes:          objectRefSchemaBuilder.GetResourceSchema(),
		},
		"features": resource_schema.SetAttribute{
			Description:         desc["features"].description,
			MarkdownDescription: desc["features"].markdown,
			Optional:            true,
			ElementType:         types.StringType,
		},
		"type": resource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"connector": resource_schema.StringAttribute{
			Description:         desc["connector"].description,
			MarkdownDescription: desc["connector"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"connector_class": resource_schema.StringAttribute{
			Description:         desc["connector_class"].description,
			MarkdownDescription: desc["connector_class"].markdown,
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"delete_threshold": resource_schema.Int32Attribute{
			Description:         desc["delete_threshold"].description,
			MarkdownDescription: desc["delete_threshold"].markdown,
			Optional:            true,
			PlanModifiers: []planmodifier.Int32{
				int32planmodifier.UseStateForUnknown(),
			},
		},
		"authoritative": resource_schema.BoolAttribute{
			Description:         desc["authoritative"].description,
			MarkdownDescription: desc["authoritative"].markdown,
			Computed:            true,
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
	}
}

func (sb *SourceSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := sb.fieldDescriptions()
	objectRefSchemaBuilder := ObjectRefSchemaBuilder{}

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
		"description": datasource_schema.StringAttribute{
			Description:         desc["description"].description,
			MarkdownDescription: desc["description"].markdown,
			Computed:            true,
		},
		"owner": datasource_schema.SingleNestedAttribute{
			Description:         desc["owner"].description,
			MarkdownDescription: desc["owner"].markdown,
			Computed:            true,
			Attributes:          objectRefSchemaBuilder.GetDataSourceSchema(),
		},
		"cluster": datasource_schema.SingleNestedAttribute{
			Description:         desc["cluster"].description,
			MarkdownDescription: desc["cluster"].markdown,
			Computed:            true,
			Attributes:          objectRefSchemaBuilder.GetDataSourceSchema(),
		},
		"features": datasource_schema.SetAttribute{
			Description:         desc["features"].description,
			MarkdownDescription: desc["features"].markdown,
			Computed:            true,
			ElementType:         types.StringType,
		},
		"type": datasource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Computed:            true,
		},
		"connector": datasource_schema.StringAttribute{
			Description:         desc["connector"].description,
			MarkdownDescription: desc["connector"].markdown,
			Computed:            true,
		},
		"connector_class": datasource_schema.StringAttribute{
			Description:         desc["connector_class"].description,
			MarkdownDescription: desc["connector_class"].markdown,
			Computed:            true,
		},
		"delete_threshold": datasource_schema.Int32Attribute{
			Description:         desc["delete_threshold"].description,
			MarkdownDescription: desc["delete_threshold"].markdown,
			Computed:            true,
		},
		"authoritative": datasource_schema.BoolAttribute{
			Description:         desc["authoritative"].description,
			MarkdownDescription: desc["authoritative"].markdown,
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
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *SourceSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct{ description, markdown string }{
		"id":               {description: "The ID of the Source.", markdown: "The unique identifier for the Source."},
		"name":             {description: "The name of the Source.", markdown: "The name of the Source resource."},
		"description":      {description: "The description of the Source.", markdown: "A brief description of the Source."},
		"owner":            {description: "The owner of the Source.", markdown: "Reference to the owner of the Source."},
		"cluster":          {description: "The cluster of the Source.", markdown: "Reference to the cluster associated with the Source."},
		"features":         {description: "A list of features enabled for the source.", markdown: "An array of features that are enabled or supported by this source."},
		"type":             {description: "The type of the source (e.g., 'Application', 'Database').", markdown: "The category or type of the source within SailPoint ISC."},
		"connector":        {description: "The connector of the Source.", markdown: "The connector of the Source."},
		"connector_class":  {description: "The class of the connector used by the source.", markdown: "The specific class name of the connector implementation for this source."},
		"delete_threshold": {description: "The delete threshold for the source.", markdown: "The threshold value that determines when accounts are deleted from the source."},
		"authoritative":    {description: "Indicates if the source is authoritative.", markdown: "A boolean flag indicating whether this source is considered authoritative for its data."},
		"created":          {description: "The creation timestamp of the Source.", markdown: "The creation timestamp of the Source."},
		"modified":         {description: "The last modified timestamp of the Source.", markdown: "The last modified timestamp of the Source."},
	}
}
