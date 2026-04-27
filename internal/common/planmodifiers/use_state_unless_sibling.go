// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseStateForUnknownUnlessSiblingChanges returns a plan modifier that copies
// the prior state value into the plan when the planned value is unknown
// (typical for `Computed` attributes), UNLESS a sibling attribute has changed
// — in which case the planned value is left as unknown so the API can resolve
// it on apply.
//
// Use this on a server-resolved attribute (e.g. `owner.name`) that is derived
// from another, user-controlled attribute (e.g. `owner.id`). Plain
// `UseStateForUnknown` would freeze the stale prior name into the plan even
// when the id changes, and the framework would then reject the apply with
// `inconsistent result after apply` once the API echoes back the new name.
//
// `siblingName` is interpreted relative to the *parent* of the attribute the
// modifier is attached to. For example, attaching this to `owner.name` with
// `siblingName = "id"` looks up `owner.id`.
func UseStateForUnknownUnlessSiblingChanges(siblingName string) planmodifier.String {
	return useStateForUnknownUnlessSiblingChangesModifier{siblingName: siblingName}
}

type useStateForUnknownUnlessSiblingChangesModifier struct {
	siblingName string
}

func (m useStateForUnknownUnlessSiblingChangesModifier) Description(_ context.Context) string {
	return fmt.Sprintf("Use prior state value when planned value is unknown, unless sibling attribute %q changes.", m.siblingName)
}

func (m useStateForUnknownUnlessSiblingChangesModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useStateForUnknownUnlessSiblingChangesModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// On Create the prior state is null — leave the planned value unknown so
	// the API can resolve it on apply.
	if req.StateValue.IsNull() {
		return
	}

	// Only act when the planned value is unknown (the typical Computed case).
	// If the user has somehow set a concrete planned value, respect it.
	if !req.PlanValue.IsUnknown() {
		return
	}

	siblingPath := req.Path.ParentPath().AtName(m.siblingName)

	var planSibling, stateSibling types.String
	if d := req.Plan.GetAttribute(ctx, siblingPath, &planSibling); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	if d := req.State.GetAttribute(ctx, siblingPath, &stateSibling); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// If the sibling has changed, leave the planned value unknown so the API
	// resolves it on apply. Otherwise reuse the prior state to keep no-op
	// plans clean.
	if !planSibling.Equal(stateSibling) {
		return
	}

	resp.PlanValue = req.StateValue
}
