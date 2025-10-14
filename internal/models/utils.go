package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// SetStringValue is a helper function to safely extract a string value from a map
// and assign it to a types.String field. If the value is an empty string,
// it preserves null to avoid inconsistent state errors.
func SetStringValue(apiModel map[string]interface{}, key string, target *types.String) {
	if v, ok := apiModel[key].(string); ok {
		if v == "" {
			// Preserve null for empty strings if target is currently null
			if target.IsNull() {
				return
			}
		}
		*target = types.StringValue(v)
	}
}
