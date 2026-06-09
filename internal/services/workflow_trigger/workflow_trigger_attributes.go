// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package workflow_trigger

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

// stripTopLevelNullKeys removes top-level keys whose value is JSON null from a
// trigger `attributes` object. SailPoint injects null-valued keys (e.g.
// `integrationId`) into the returned attributes that the practitioner's config
// cannot express; left in state they cause spurious drift on refresh. Returns
// the input unchanged when it is null/unknown, not a JSON object, or has no
// null keys (#136).
func stripTopLevelNullKeys(attrs jsontypes.Normalized) jsontypes.Normalized {
	if attrs.IsNull() || attrs.IsUnknown() {
		return attrs
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(attrs.ValueString()), &obj); err != nil {
		// Not a JSON object (or unparseable) — leave it untouched.
		return attrs
	}

	changed := false
	for k, v := range obj {
		if string(v) == "null" {
			delete(obj, k)
			changed = true
		}
	}
	if !changed {
		return attrs
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return attrs
	}
	return jsontypes.NewNormalizedValue(string(b))
}
