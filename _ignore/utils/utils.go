package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

func SetStringFromAny(v any, target *types.String) {
	if str, ok := v.(string); ok {
		*target = types.StringValue(str)
	} else {
		*target = types.StringNull()
	}
}

func SetStringFromAnyIfNotNull(v any, target *types.String) {
	if !target.IsNull() {
		SetStringFromAny(v, target)
	}
}
