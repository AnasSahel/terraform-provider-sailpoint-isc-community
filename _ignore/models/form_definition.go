// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ BaseModel[FormDefinitionModel] = &FormDefinitionModel{}
)

type FormDefinitionModel struct {
	Id          types.String       `tfsdk:"id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	Created     types.String       `tfsdk:"created"`
	Modified    types.String       `tfsdk:"modified"`
	Owner       *ResourceRefModel  `tfsdk:"owner"`
	UsedBy      []ResourceRefModel `tfsdk:"used_by"`
}

// func (f *FormDefinitionModel) SetIdFromAnyIf(ctx context.Context, id any, shouldSet bool) *FormDefinitionModel {
// 	SetTerraformStringFromAnyIf(ctx, id, &f.Id, shouldSet)
// 	return f
// }

// func (f *FormDefinitionModel) SetIdFromAny(ctx context.Context, id any) *FormDefinitionModel {
// 	return f.SetIdFromAnyIf(ctx, id, true)
// }

// func (f *FormDefinitionModel) SetNameFromAnyIf(ctx context.Context, name any, shouldSet bool) *FormDefinitionModel {
// 	SetTerraformStringFromAnyIf(ctx, name, &f.Name, shouldSet)
// 	return f
// }

// func (f *FormDefinitionModel) SetNameFromAny(ctx context.Context, name any) *FormDefinitionModel {
// 	return f.SetNameFromAnyIf(ctx, name, true)
// }

// func (f *FormDefinitionModel) SetDescriptionFromAnyIf(ctx context.Context, description any, shouldSet bool) *FormDefinitionModel {
// 	SetTerraformStringFromAnyIf(ctx, description, &f.Description, shouldSet)
// 	return f
// }

// func (f *FormDefinitionModel) SetDescriptionFromAny(ctx context.Context, description any) *FormDefinitionModel {
// 	return f.SetDescriptionFromAnyIf(ctx, description, true)
// }

// func (f *FormDefinitionModel) SetCreatedFromAnyIf(ctx context.Context, created any, shouldSet bool) *FormDefinitionModel {
// 	SetTerraformStringFromAnyIf(ctx, created, &f.Created, shouldSet)
// 	return f
// }

// func (f *FormDefinitionModel) SetCreatedFromAny(ctx context.Context, created any) *FormDefinitionModel {
// 	return f.SetCreatedFromAnyIf(ctx, created, true)
// }

// func (f *FormDefinitionModel) SetModifiedFromAnyIf(ctx context.Context, modified any, shouldSet bool) *FormDefinitionModel {
// 	SetTerraformStringFromAnyIf(ctx, modified, &f.Modified, shouldSet)
// 	return f
// }

// func (f *FormDefinitionModel) SetModifiedFromAny(ctx context.Context, modified any) *FormDefinitionModel {
// 	return f.SetModifiedFromAnyIf(ctx, modified, true)
// }

// func (f *FormDefinitionModel) SetOwnerFromAnyIf(ctx context.Context, owner any, shouldSet bool) *FormDefinitionModel {
// 	if shouldSet {
// 		if ownerMap, ok := owner.(map[string]interface{}); ok {
// 			f.Owner = &ResourceRefModel{}
// 			f.Owner.FromSailPointModel(ctx, ownerMap, ConversionOptions[ResourceRefModel]{})
// 		} else {
// 			f.Owner = nil
// 		}
// 	}
// 	return f
// }

// func (f *FormDefinitionModel) SetOwnerFromAny(ctx context.Context, owner any) *FormDefinitionModel {
// 	return f.SetOwnerFromAnyIf(ctx, owner, true)
// }

// func (f *FormDefinitionModel) SetUsedByFromAnyIf(ctx context.Context, usedBy any, shouldSet bool) *FormDefinitionModel {
// 	if shouldSet {
// 		// Check if usedBy is a slice of interfaces
// 		if usedByList, ok := usedBy.([]interface{}); ok {
// 			f.UsedBy = make([]ResourceRefModel, len(usedByList)) // Initialize the slice with the correct length
// 			for i, v := range usedByList {
// 				f.UsedBy[i].FromSailPointModel(ctx, v.(map[string]interface{}), ConversionOptions[ResourceRefModel]{})
// 			}
// 		} else {
// 			f.UsedBy = nil
// 		}
// 	}
// 	return f
// }

// func (f *FormDefinitionModel) SetUsedByFromAny(ctx context.Context, usedBy any) *FormDefinitionModel {
// 	return f.SetUsedByFromAnyIf(ctx, usedBy, true)
// }

// FromSailPointGetResponse implements BaseModel.
func (f *FormDefinitionModel) FromSailPointModel(ctx context.Context, input any, opts ConversionOptions[FormDefinitionModel]) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Trace(ctx, "Converting from SailPoint model to FormDefinitionModel", map[string]any{"input": input})

	f.SetIdFromAny(ctx, input.(map[string]interface{})["id"])
	f.SetNameFromAny(ctx, input.(map[string]interface{})["name"])

	f.SetCreatedFromAny(ctx, input.(map[string]interface{})["created"])
	f.SetModifiedFromAny(ctx, input.(map[string]interface{})["modified"])

	f.SetOwnerFromAny(ctx, input.(map[string]interface{})["owner"])

	f.SetDescriptionFromAnyIf(ctx, input.(map[string]interface{})["description"], opts.Plan == nil || !opts.Plan.Description.IsNull())
	f.SetUsedByFromAnyIf(ctx, input.(map[string]interface{})["usedBy"], opts.Plan == nil || opts.Plan.UsedBy != nil)

	return diags
}

// ToSailPointCreateRequest implements BaseModel.
func (f *FormDefinitionModel) ToSailPointCreateRequest(ctx context.Context) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	owner, d := f.Owner.ToSailPointCreateRequest(ctx)
	diags.Append(d...)

	req := map[string]interface{}{
		"id":    f.Id.ValueString(),
		"name":  f.Name.ValueString(),
		"owner": owner,
	}

	if !f.Description.IsNull() {
		req["description"] = f.Description.ValueString()
	}

	if len(f.UsedBy) > 0 {
		req["usedBy"] = make([]map[string]interface{}, 0, len(f.UsedBy))
		for _, v := range f.UsedBy {
			usedByReq, d := v.ToSailPointCreateRequest(ctx)
			diags.Append(d...)
			req["usedBy"] = append(req["usedBy"].([]map[string]interface{}), usedByReq.(map[string]interface{}))
		}
	}

	return req, diags
}
