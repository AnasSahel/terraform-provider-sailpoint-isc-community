package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// String
func NewGoTypeValueIf[TTerraform types.Set | types.String | types.Int32 | types.Bool | types.Int64 | types.Float64, TGo []string | string | int32 | bool | int64 | float64](ctx context.Context, val TTerraform, shouldSet bool) TGo {
	if shouldSet {
		switch v := any(val).(type) {
		case types.String:
			return any(v.ValueString()).(TGo)
		case types.Int32:
			return any(v.ValueInt32()).(TGo)
		case types.Int64:
			return any(v.ValueInt64()).(TGo)
		case types.Bool:
			return any(v.ValueBool()).(TGo)
		case types.Float64:
			return any(v.ValueFloat64()).(TGo)
		case types.Set:
			var result []string
			v.ElementsAs(ctx, &result, false)
			return any(result).(TGo)
		}
	}

	var zero TGo
	return zero
}

func NewTerraformTypeValueIf[TTerraform types.Set | types.String | types.Int32 | types.Bool | types.Int64 | types.Float64, TGo []string | string | int32 | bool | int64 | float64](ctx context.Context, val TGo, shouldSet bool) TTerraform {
	if shouldSet {
		switch any(val).(type) {
		case string:
			return any(types.StringValue(any(val).(string))).(TTerraform)
		case int32:
			return any(types.Int32Value(any(val).(int32))).(TTerraform)
		case int64:
			return any(types.Int64Value(any(val).(int64))).(TTerraform)
		case bool:
			return any(types.BoolValue(any(val).(bool))).(TTerraform)
		case float64:
			return any(types.Float64Value(any(val).(float64))).(TTerraform)
		case []string:
			arr := any(val).([]string)
			setValue, _ := types.SetValueFrom(ctx, types.StringType, arr)
			return any(setValue).(TTerraform)
		}
	}

	switch any(*new(TTerraform)).(type) {
	case types.String:
		return any(types.StringNull()).(TTerraform)
	case types.Int32:
		return any(types.Int32Null()).(TTerraform)
	case types.Int64:
		return any(types.Int64Null()).(TTerraform)
	case types.Bool:
		return any(types.BoolNull()).(TTerraform)
	case types.Float64:
		return any(types.Float64Null()).(TTerraform)
	case types.Set:
		return any(types.SetNull(types.StringType)).(TTerraform)
	}

	var zero TTerraform
	return zero
}

func IsTerraformValueNullOrUnknown[T types.Set | types.String | types.Int32 | types.Bool | types.Int64 | types.Float64](val T) bool {
	switch v := any(val).(type) {
	case types.String:
		return v.IsNull() || v.IsUnknown()
	case types.Int32:
		return v.IsNull() || v.IsUnknown()
	case types.Bool:
		return v.IsNull() || v.IsUnknown()
	case types.Int64:
		return v.IsNull() || v.IsUnknown()
	case types.Float64:
		return v.IsNull() || v.IsUnknown()
	case types.Set:
		return v.IsNull() || v.IsUnknown()
	default:
		return false
	}
}
