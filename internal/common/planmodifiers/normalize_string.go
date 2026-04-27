// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

// Package planmodifiers contains reusable plan modifiers for the SailPoint ISC
// Terraform provider. These exist mainly to absorb server-side normalizations
// at plan time so apply does not fail with "inconsistent result after apply".
package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NormalizeString returns a plan modifier that rewrites the planned value of a
// string attribute when it matches one of the keys in `aliases`, replacing it
// with the canonical value. Use this when the SailPoint API silently rewrites
// an input value server-side: aligning the plan with the canonical form
// prevents the framework from rejecting the apply with "inconsistent result
// after apply" once the server echoes the rewritten value back.
func NormalizeString(aliases map[string]string) planmodifier.String {
	return normalizeStringModifier{aliases: aliases}
}

type normalizeStringModifier struct {
	aliases map[string]string
}

func (m normalizeStringModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Normalizes alias values to their canonical form: %v", m.aliases)
}

func (m normalizeStringModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m normalizeStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if canonical, ok := m.aliases[req.PlanValue.ValueString()]; ok {
		resp.PlanValue = types.StringValue(canonical)
	}
}
