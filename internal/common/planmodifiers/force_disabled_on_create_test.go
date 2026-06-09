// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestForceDisabledOnCreate_ForcesToFalseOnCreate(t *testing.T) {
	t.Parallel()

	modifier := ForceDisabledOnCreate()
	req := planmodifier.BoolRequest{
		StateValue: types.BoolNull(), // null state value = create (no prior state)
		PlanValue:  types.BoolValue(true),
	}
	resp := &planmodifier.BoolResponse{PlanValue: req.PlanValue}
	modifier.PlanModifyBool(context.Background(), req, resp)

	if resp.PlanValue.ValueBool() {
		t.Error("expected PlanValue to be forced to false on create, got true")
	}
}

func TestForceDisabledOnCreate_PassesThroughOnUpdate(t *testing.T) {
	t.Parallel()

	modifier := ForceDisabledOnCreate()
	req := planmodifier.BoolRequest{
		StateValue: types.BoolValue(false), // non-null state value = update
		PlanValue:  types.BoolValue(true),
	}
	resp := &planmodifier.BoolResponse{PlanValue: req.PlanValue}
	modifier.PlanModifyBool(context.Background(), req, resp)

	if !resp.PlanValue.ValueBool() {
		t.Error("expected PlanValue to remain true on update, got false")
	}
}

func TestForceDisabledOnCreate_AlreadyFalseOnCreate(t *testing.T) {
	t.Parallel()

	modifier := ForceDisabledOnCreate()
	req := planmodifier.BoolRequest{
		StateValue: types.BoolNull(),
		PlanValue:  types.BoolValue(false),
	}
	resp := &planmodifier.BoolResponse{PlanValue: req.PlanValue}
	modifier.PlanModifyBool(context.Background(), req, resp)

	if resp.PlanValue.ValueBool() {
		t.Error("expected PlanValue to remain false on create")
	}
}
