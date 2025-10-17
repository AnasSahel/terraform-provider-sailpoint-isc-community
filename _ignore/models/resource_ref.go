package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ BaseModel[ResourceRefModel] = &ResourceRefModel{}
)

type ResourceRefModel struct {
	Type types.String `tfsdk:"type"`
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// FromSailPointGetResponse implements BaseModel.
func (r *ResourceRefModel) FromSailPointModel(ctx context.Context, input any, opts ConversionOptions[ResourceRefModel]) diag.Diagnostics {
	SetTerraformStringFromAny(ctx, input.(map[string]interface{})["type"], &r.Type)
	SetTerraformStringFromAny(ctx, input.(map[string]interface{})["id"], &r.Id)
	SetTerraformStringFromAny(ctx, input.(map[string]interface{})["name"], &r.Name)

	return nil
}

// ToSailPointCreateRequest implements BaseModel.
func (r *ResourceRefModel) ToSailPointCreateRequest(ctx context.Context) (any, diag.Diagnostics) {
	req := map[string]interface{}{
		"type": r.Type.ValueString(),
		"id":   r.Id.ValueString(),
		"name": r.Name.ValueString(),
	}

	return req, nil
}
