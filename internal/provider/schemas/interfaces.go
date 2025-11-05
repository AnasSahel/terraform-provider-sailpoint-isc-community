// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemas

import (
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SchemaBuilder interface {
	GetResourceSchema() map[string]resource_schema.Attribute
	GetDataSourceSchema() map[string]datasource_schema.Attribute

	fieldDescriptions() map[string]struct{ description, markdown string }
}
