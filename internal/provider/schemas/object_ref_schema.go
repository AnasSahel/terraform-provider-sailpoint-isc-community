package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type ObjectRefSchemaBuilder struct{}

var (
	_ SchemaBuilder = &ObjectRefSchemaBuilder{}
)

// GetDataSourceSchema implements SchemaBuilder.
func (o *ObjectRefSchemaBuilder) GetDataSourceSchema() map[string]datasource_schema.Attribute {
	desc := o.fieldDescriptions()

	return map[string]datasource_schema.Attribute{
		"type": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
		},
		"id": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
		},
		"name": datasource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
		},
	}
}

// GetResourceSchema implements SchemaBuilder.
func (o *ObjectRefSchemaBuilder) GetResourceSchema() map[string]resource_schema.Attribute {
	desc := o.fieldDescriptions()

	return map[string]resource_schema.Attribute{
		"type": resource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["type"].description,
			MarkdownDescription: desc["type"].markdown,
		},
		"id": resource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["id"].description,
			MarkdownDescription: desc["id"].markdown,
		},
		"name": resource_schema.StringAttribute{
			Computed:            true,
			Description:         desc["name"].description,
			MarkdownDescription: desc["name"].markdown,
		},
	}
}

// fieldDescriptions implements SchemaBuilder.
func (o *ObjectRefSchemaBuilder) fieldDescriptions() map[string]struct {
	description string
	markdown    string
} {
	return map[string]struct {
		description string
		markdown    string
	}{
		"type": {description: "The type of the referenced object.", markdown: "The type of the referenced object."},
		"id":   {description: "The unique identifier of the referenced object.", markdown: "The unique identifier (UUID) of the referenced object."},
		"name": {description: "The name of the referenced object.", markdown: "The human-readable name of the referenced object."},
	}
}
