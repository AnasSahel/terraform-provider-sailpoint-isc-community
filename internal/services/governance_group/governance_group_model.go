// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package governance_group

import (
	"context"
	"reflect"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type governanceGroupMemberModel struct {
	Type types.String `tfsdk:"type"`
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type governanceGroupModel struct {
	ID          types.String                 `tfsdk:"id"`
	Name        types.String                 `tfsdk:"name"`
	Description types.String                 `tfsdk:"description"`
	Owner       *common.ObjectRefModel       `tfsdk:"owner"`
	Members     []governanceGroupMemberModel `tfsdk:"members"`
	MemberCount types.Int64                  `tfsdk:"member_count"`
	Created     types.String                 `tfsdk:"created"`
	Modified    types.String                 `tfsdk:"modified"`
}

// ---------------------------------------------------------------------------
// FromAPI
// ---------------------------------------------------------------------------

func (m *governanceGroupModel) FromAPI(ctx context.Context, api *client.GovernanceGroupAPI, members []client.GovernanceGroupMemberAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Description = stringPtrToTF(api.Description)
	m.Created = stringPtrToTF(api.Created)
	m.Modified = stringPtrToTF(api.Modified)

	if api.MemberCount != nil {
		m.MemberCount = types.Int64Value(*api.MemberCount)
	} else {
		m.MemberCount = types.Int64Value(int64(len(members)))
	}

	owner, diags := common.NewObjectRefFromAPIPtr(ctx, api.Owner)
	diagnostics.Append(diags...)
	m.Owner = owner

	// Populate members from the separate members API response.
	if len(members) > 0 {
		m.Members = make([]governanceGroupMemberModel, 0, len(members))
		for _, mem := range members {
			m.Members = append(m.Members, governanceGroupMemberModel{
				Type: types.StringValue(mem.Type),
				ID:   types.StringValue(mem.ID),
				Name: types.StringValue(mem.Name),
			})
		}
	} else {
		m.Members = nil
	}

	return diagnostics
}

// ---------------------------------------------------------------------------
// ToAPI
// ---------------------------------------------------------------------------

func (m *governanceGroupModel) ToAPI(ctx context.Context) (*client.GovernanceGroupAPI, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	api := &client.GovernanceGroupAPI{
		Name: m.Name.ValueString(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		v := m.Description.ValueString()
		api.Description = &v
	}

	if m.Owner != nil {
		ownerAPI, diags := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(diags...)
		api.Owner = ownerAPI
	}

	return api, diagnostics
}

// ---------------------------------------------------------------------------
// ToPatchOperations
// ---------------------------------------------------------------------------

func (m *governanceGroupModel) ToPatchOperations(ctx context.Context, state *governanceGroupModel) ([]client.JSONPatchOperation, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	var ops []client.JSONPatchOperation

	if !m.Name.Equal(state.Name) {
		ops = append(ops, client.NewReplacePatch("/name", m.Name.ValueString()))
	}

	if !m.Description.Equal(state.Description) {
		if !m.Description.IsNull() {
			ops = append(ops, client.NewReplacePatch("/description", m.Description.ValueString()))
		} else {
			ops = append(ops, client.NewRemovePatch("/description"))
		}
	}

	if !reflect.DeepEqual(m.Owner, state.Owner) && m.Owner != nil {
		ownerAPI, diags := common.NewObjectRefToAPI(ctx, *m.Owner)
		diagnostics.Append(diags...)
		ops = append(ops, client.NewReplacePatch("/owner", ownerAPI))
	}

	return ops, diagnostics
}

// ---------------------------------------------------------------------------
// Member diff helpers
// ---------------------------------------------------------------------------

// MemberDiff computes which members to add and which to remove to reconcile
// the plan members against the current state members.
func MemberDiff(plan []governanceGroupMemberModel, state []governanceGroupMemberModel) (toAdd []client.GovernanceGroupMemberAPI, toRemoveIDs []string) {
	stateByID := make(map[string]bool, len(state))
	for _, m := range state {
		stateByID[m.ID.ValueString()] = true
	}

	planByID := make(map[string]bool, len(plan))
	for _, m := range plan {
		planByID[m.ID.ValueString()] = true
	}

	for _, m := range plan {
		id := m.ID.ValueString()
		if !stateByID[id] {
			toAdd = append(toAdd, client.GovernanceGroupMemberAPI{
				Type: m.Type.ValueString(),
				ID:   id,
			})
		}
	}

	for _, m := range state {
		id := m.ID.ValueString()
		if !planByID[id] {
			toRemoveIDs = append(toRemoveIDs, id)
		}
	}

	return toAdd, toRemoveIDs
}

// ---------------------------------------------------------------------------
// Primitive helpers
// ---------------------------------------------------------------------------

func stringPtrToTF(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
