package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// String
func NewGoTypeValueIf[TTerraform types.String | types.Int32 | types.Bool | types.Int64 | types.Float64, TGo string | int32 | bool | int64 | float64](val TTerraform, shouldSet bool) TGo {
	if shouldSet {
		switch v := any(val).(type) {
		case types.String:
			return any(v.ValueString()).(TGo)
		case types.Int32:
			return any(v.ValueInt32()).(TGo)
		case types.Bool:
			return any(v.ValueBool()).(TGo)
		case types.Int64:
			return any(v.ValueInt64()).(TGo)
		case types.Float64:
			return any(v.ValueFloat64()).(TGo)
		}
	}

	var zero TGo
	return zero
}

func NewTerraformTypeValueIf[TTerraform types.String | types.Int32 | types.Bool | types.Int64 | types.Float64, TGo string | int32 | bool | int64 | float64](val TGo, shouldSet bool) TTerraform {
	if shouldSet {
		switch any(val).(type) {
		case string:
			return any(types.StringValue(any(val).(string))).(TTerraform)
		case int32:
			return any(types.Int32Value(any(val).(int32))).(TTerraform)
		case bool:
			return any(types.BoolValue(any(val).(bool))).(TTerraform)
		case int64:
			return any(types.Int64Value(any(val).(int64))).(TTerraform)
		case float64:
			return any(types.Float64Value(any(val).(float64))).(TTerraform)
		}
	}

	var zero TTerraform
	return zero
}

func IsTerraformValueNullOrUnknown[T types.String | types.Int32 | types.Bool | types.Int64 | types.Float64](val T) bool {
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
	default:
		return false
	}
}
