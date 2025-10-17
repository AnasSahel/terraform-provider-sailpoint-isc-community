package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConversionOptions[T any] struct {
	Plan *T // For resource Create/Update to compare with state
}

type BaseModel[T any] interface {
	FromSailPointModel(ctx context.Context, input any, opts ConversionOptions[T]) diag.Diagnostics

	ToSailPointCreateRequest(ctx context.Context) (any, diag.Diagnostics)
}

func SetTerraformStringFromAny(ctx context.Context, v any, target *types.String) {
	if str, ok := v.(string); ok {
		*target = types.StringValue(str)
	} else {
		*target = types.StringNull()
	}
}

func SetTerraformStringFromAnyIf(ctx context.Context, v any, target *types.String, shouldSet bool) {
	if shouldSet {
		SetTerraformStringFromAny(ctx, v, target)
	}
}

// SetNestedModelFromAnyIf sets a pointer to a nested model from any type
func SetNestedModelFromAnyIf[T BaseModel[T]](ctx context.Context, value any, target **T, shouldSet bool) {
	if !shouldSet {
		return
	}

	if valueMap, ok := value.(map[string]interface{}); ok {
		var model T
		model.FromSailPointModel(ctx, valueMap, ConversionOptions[T]{})
		*target = &model
	} else {
		*target = nil
	}
}

// SetNestedModelSliceFromAnyIf sets a slice of nested models from any type
func SetNestedModelSliceFromAnyIf[T BaseModel[T]](ctx context.Context, value any, target *[]T, shouldSet bool) {
	if !shouldSet {
		return
	}

	if valueList, ok := value.([]interface{}); ok {
		*target = make([]T, len(valueList))
		for i, v := range valueList {
			if vMap, ok := v.(map[string]interface{}); ok {
				(*target)[i].FromSailPointModel(ctx, vMap, ConversionOptions[T]{})
			}
		}
	} else {
		*target = nil
	}
}
