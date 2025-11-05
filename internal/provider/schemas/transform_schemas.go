package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

type TransformSchemaBuilder struct{}

var (
	_ SchemaBuilder = &TransformSchemaBuilder{}
)

// GetResourceSchema implements SchemaBuilder for Transform resource.
func (sb *TransformSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
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
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"type": resource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"attributes": resource_schema.StringAttribute{
			Description:         desc["attributes"].description,
			MarkdownDescription: desc["attributes"].markdown,
			Required:            true,
		},
		"internal": resource_schema.BoolAttribute{
			Description:         desc["internal"].description,
			MarkdownDescription: desc["internal"].markdown,
			Computed:            true,
		},
	}
}

// GetDataSourceSchema implements SchemaBuilder for Transform data source.
func (sb *TransformSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
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
		"type": datasource_schema.StringAttribute{
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
			Computed:            true,
		},
		"attributes": datasource_schema.StringAttribute{
			Description:         desc["attributes"].description,
			MarkdownDescription: desc["attributes"].markdown,
			Computed:            true,
		},
		"internal": datasource_schema.BoolAttribute{
			Description:         desc["internal"].description,
			MarkdownDescription: desc["internal"].markdown,
			Computed:            true,
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (sb *TransformSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"id": {
			description: "Unique identifier of the transform.",
			markdown:    "Unique identifier (UUID) of the transform.",
		},
		"name": {
			description: "Name of the transform.",
			markdown:    "Name of the transform as it appears in the UI. This field is immutable after creation.",
		},
		"type": {
			description: "Type of the transform operation.",
			markdown:    "Type of the transform operation (e.g., `lower`, `upper`, `lookup`, `static`). This field is immutable after creation. See [Transform Operations](https://developer.sailpoint.com/docs/extensibility/transforms/operations) for available types.",
		},
		"attributes": {
			description: "Configuration attributes for the transform as a JSON string.",
			markdown:    "Configuration attributes for the transform as a JSON string. The structure varies by transform type. This is the only field that can be updated after creation.",
		},
		"internal": {
			description: "Indicates whether this is an internal SailPoint transform.",
			markdown:    "Indicates whether this is an internal SailPoint transform (true) or a custom transform (false). Only SailPoint employees can create internal transforms.",
		},
	}
}
