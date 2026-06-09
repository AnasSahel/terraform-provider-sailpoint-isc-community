// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ForceDisabledOnCreate returns a plan modifier that sets the planned value of a
// bool attribute to false on resource creation (when there is no prior state).
// On updates the planned value is left unchanged.
//
// Use this for the workflow enabled attribute: the SailPoint API cannot create a
// workflow in the enabled state. Declaring enabled = true at create time will
// produce enabled = false in state; the next apply (an update) will enable the
// workflow ("converge over two applies").
func ForceDisabledOnCreate() planmodifier.Bool {
	return forceDisabledOnCreateModifier{}
}

type forceDisabledOnCreateModifier struct{}

func (m forceDisabledOnCreateModifier) Description(_ context.Context) string {
	return "Forces the planned value to false on create; on update the declared value is used as-is."
}

func (m forceDisabledOnCreateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m forceDisabledOnCreateModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// On create there is no prior state, so StateValue is null.
	// Force the planned value to false: the API cannot create an enabled workflow.
	if req.StateValue.IsNull() {
		resp.PlanValue = types.BoolValue(false)
	}
}
