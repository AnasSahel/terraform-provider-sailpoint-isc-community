// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// GetTransformResourceSchema returns the schema definition for the transform resource.
func GetTransformResourceSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SailPoint ISC Transform resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Transform identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Transform name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Transform type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"accountAttribute",
						"base64Decode",
						"base64Encode",
						"concatenation",
						"conditional",
						"dateCompare",
						"dateFormat",
						"dateMath",
						"decompose",
						"displayName",
						"e164phone",
						"firstValid",
						"getReference",
						"getReferenceIdentityAttribute",
						"identityAttribute",
						"indexOf",
						"iso3166",
						"lastIndexOf",
						"leftPad",
						"lookup",
						"lower",
						"normalizeNames",
						"randomAlphaNumeric",
						"randomNumeric",
						"replace",
						"replaceAll",
						"rightPad",
						"rule",
						"split",
						"static",
						"substring",
						"trim",
						"upper",
						"uuid",
					),
				},
			},
			"internal": schema.BoolAttribute{
				MarkdownDescription: "Indicates if the transform is internal",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"attributes": schema.StringAttribute{
				MarkdownDescription: "Transform attributes as JSON",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\s*\{.*\}\s*$`),
						"must be valid JSON object",
					),
				},
			},
		},
	}
}
