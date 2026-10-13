// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package approval_config

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type approvalConfigModel struct {
	ID                               types.String `tfsdk:"id"`
	Scope                            types.String `tfsdk:"scope"`
	ObjectID                         types.String `tfsdk:"object_id"`
	ReminderConfig                   types.Object `tfsdk:"reminder_config"`
	EscalationConfig                 types.Object `tfsdk:"escalation_config"`
	TimeoutConfig                    types.Object `tfsdk:"timeout_config"`
	CronTimezone                     types.Object `tfsdk:"cron_timezone"`
	SerialChain                      types.List   `tfsdk:"serial_chain"`
	RequiresComment                  types.String `tfsdk:"requires_comment"`
	FallbackApprover                 types.Object `tfsdk:"fallback_approver"`
	MachineIdentityManagerAssignment types.String `tfsdk:"machine_identity_manager_assignment"`
	CircumventApprovalProcess        types.Bool   `tfsdk:"circumvent_approval_process"`
	AutoApprove                      types.String `tfsdk:"auto_approve"`
}

type reminderConfigModel struct {
	Enabled                types.Bool   `tfsdk:"enabled"`
	DaysUntilFirstReminder types.Int64  `tfsdk:"days_until_first_reminder"`
	ReminderCronSchedule   types.String `tfsdk:"reminder_cron_schedule"`
	MaxReminders           types.Int64  `tfsdk:"max_reminders"`
}

type escalationConfigModel struct {
	Enabled                  types.Bool   `tfsdk:"enabled"`
	DaysUntilFirstEscalation types.Int64  `tfsdk:"days_until_first_escalation"`
	EscalationCronSchedule   types.String `tfsdk:"escalation_cron_schedule"`
	EscalationChain          types.List   `tfsdk:"escalation_chain"`
}

type chainItemModel struct {
	IdentityType types.String `tfsdk:"identity_type"`
	IdentityID   types.String `tfsdk:"identity_id"`
}

type timeoutConfigModel struct {
	Enabled          types.Bool   `tfsdk:"enabled"`
	DaysUntilTimeout types.Int64  `tfsdk:"days_until_timeout"`
	TimeoutResult    types.String `tfsdk:"timeout_result"`
}

type cronTimezoneModel struct {
	Location types.String `tfsdk:"location"`
	Offset   types.String `tfsdk:"offset"`
}

type fallbackApproverModel struct {
	Type       types.String `tfsdk:"type"`
	IdentityID types.String `tfsdk:"identity_id"`
}

var chainItemAttrTypes = map[string]attr.Type{
	"identity_type": types.StringType,
	"identity_id":   types.StringType,
}

var chainItemObjectType = types.ObjectType{AttrTypes: chainItemAttrTypes}

var reminderConfigAttrTypes = map[string]attr.Type{
	"enabled":                   types.BoolType,
	"days_until_first_reminder": types.Int64Type,
	"reminder_cron_schedule":    types.StringType,
	"max_reminders":             types.Int64Type,
}

var escalationConfigAttrTypes = map[string]attr.Type{
	"enabled":                     types.BoolType,
	"days_until_first_escalation": types.Int64Type,
	"escalation_cron_schedule":    types.StringType,
	"escalation_chain":            types.ListType{ElemType: chainItemObjectType},
}

var timeoutConfigAttrTypes = map[string]attr.Type{
	"enabled":            types.BoolType,
	"days_until_timeout": types.Int64Type,
	"timeout_result":     types.StringType,
}

var cronTimezoneAttrTypes = map[string]attr.Type{
	"location": types.StringType,
	"offset":   types.StringType,
}

var fallbackApproverAttrTypes = map[string]attr.Type{
	"type":        types.StringType,
	"identity_id": types.StringType,
}

// knownBool / knownInt64 / knownString return nil for null or unknown values, else a pointer.
func knownBool(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

func knownInt64(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := v.ValueInt64()
	return &i
}

func knownString(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

// isOwnConfig reports whether api is objectID's own config rather than an inherited one (design D5).
func isOwnConfig(api *client.ApprovalConfigAPI, objectID string) bool {
	return api != nil && api.ID != nil && *api.ID == objectID
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func chainFromAPI(ctx context.Context, items *[]client.ApprovalChainItemAPI) (types.List, diag.Diagnostics) {
	if items == nil || len(*items) == 0 {
		return types.ListValueMust(chainItemObjectType, []attr.Value{}), nil
	}
	models := make([]chainItemModel, 0, len(*items))
	for _, it := range *items {
		models = append(models, chainItemModel{
			IdentityType: types.StringValue(it.IdentityType),
			IdentityID:   common.StringOrNullIfEmpty(it.IdentityID),
		})
	}
	return types.ListValueFrom(ctx, chainItemObjectType, models)
}

func chainToAPI(ctx context.Context, list types.List) (*[]client.ApprovalChainItemAPI, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var models []chainItemModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	items := make([]client.ApprovalChainItemAPI, 0, len(models))
	for _, m := range models {
		items = append(items, client.ApprovalChainItemAPI{
			IdentityID:   m.IdentityID.ValueString(),
			IdentityType: m.IdentityType.ValueString(),
		})
	}
	return &items, nil
}

// FromAPI maps an API config onto the model. It never touches Scope or ObjectID (design D4).
func (m *approvalConfigModel) FromAPI(ctx context.Context, api *client.ApprovalConfigAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(m.Scope.ValueString() + "/" + m.ObjectID.ValueString())

	// reminder_config
	m.ReminderConfig = types.ObjectNull(reminderConfigAttrTypes)
	if api.ReminderConfig != nil {
		v, d := types.ObjectValueFrom(ctx, reminderConfigAttrTypes, reminderConfigModel{
			Enabled:                types.BoolPointerValue(api.ReminderConfig.Enabled),
			DaysUntilFirstReminder: types.Int64PointerValue(api.ReminderConfig.DaysUntilFirstReminder),
			ReminderCronSchedule:   types.StringPointerValue(api.ReminderConfig.ReminderCronSchedule),
			MaxReminders:           types.Int64PointerValue(api.ReminderConfig.MaxReminders),
		})
		diags.Append(d...)
		m.ReminderConfig = v
	}

	// escalation_config
	m.EscalationConfig = types.ObjectNull(escalationConfigAttrTypes)
	if api.EscalationConfig != nil {
		chain, d := chainFromAPI(ctx, api.EscalationConfig.EscalationChain)
		diags.Append(d...)
		v, d := types.ObjectValueFrom(ctx, escalationConfigAttrTypes, escalationConfigModel{
			Enabled:                  types.BoolPointerValue(api.EscalationConfig.Enabled),
			DaysUntilFirstEscalation: types.Int64PointerValue(api.EscalationConfig.DaysUntilFirstEscalation),
			EscalationCronSchedule:   types.StringPointerValue(api.EscalationConfig.EscalationCronSchedule),
			EscalationChain:          chain,
		})
		diags.Append(d...)
		m.EscalationConfig = v
	}

	// timeout_config
	m.TimeoutConfig = types.ObjectNull(timeoutConfigAttrTypes)
	if api.TimeoutConfig != nil {
		v, d := types.ObjectValueFrom(ctx, timeoutConfigAttrTypes, timeoutConfigModel{
			Enabled:          types.BoolPointerValue(api.TimeoutConfig.Enabled),
			DaysUntilTimeout: types.Int64PointerValue(api.TimeoutConfig.DaysUntilTimeout),
			TimeoutResult:    types.StringPointerValue(api.TimeoutConfig.TimeoutResult),
		})
		diags.Append(d...)
		m.TimeoutConfig = v
	}

	// cron_timezone
	m.CronTimezone = types.ObjectNull(cronTimezoneAttrTypes)
	if api.CronTimezone != nil {
		v, d := types.ObjectValueFrom(ctx, cronTimezoneAttrTypes, cronTimezoneModel{
			Location: types.StringPointerValue(api.CronTimezone.Location),
			Offset:   common.StringOrNullIfEmpty(deref(api.CronTimezone.Offset)),
		})
		diags.Append(d...)
		m.CronTimezone = v
	}

	// serial_chain
	chain, d := chainFromAPI(ctx, api.SerialChain)
	diags.Append(d...)
	m.SerialChain = chain

	// fallback_approver
	m.FallbackApprover = types.ObjectNull(fallbackApproverAttrTypes)
	if api.FallbackApprover != nil {
		v, d := types.ObjectValueFrom(ctx, fallbackApproverAttrTypes, fallbackApproverModel{
			Type:       types.StringValue(api.FallbackApprover.Type),
			IdentityID: common.StringOrNullIfEmpty(api.FallbackApprover.IdentityID),
		})
		diags.Append(d...)
		m.FallbackApprover = v
	}

	m.RequiresComment = types.StringPointerValue(api.RequiresComment)
	m.MachineIdentityManagerAssignment = types.StringPointerValue(api.MachineIdentityManagerAssignment)
	m.CircumventApprovalProcess = types.BoolPointerValue(api.CircumventApprovalProcess)
	m.AutoApprove = types.StringPointerValue(api.AutoApprove)

	return diags
}

// ToAPI builds the PUT body. ID and Scope stay nil (design D15).
func (m *approvalConfigModel) ToAPI(ctx context.Context) (*client.ApprovalConfigAPI, diag.Diagnostics) {
	var diags diag.Diagnostics
	api := &client.ApprovalConfigAPI{
		RequiresComment:                  knownString(m.RequiresComment),
		MachineIdentityManagerAssignment: knownString(m.MachineIdentityManagerAssignment),
		CircumventApprovalProcess:        knownBool(m.CircumventApprovalProcess),
		AutoApprove:                      knownString(m.AutoApprove),
	}

	if !m.ReminderConfig.IsNull() && !m.ReminderConfig.IsUnknown() {
		var rc reminderConfigModel
		diags.Append(m.ReminderConfig.As(ctx, &rc, basetypes.ObjectAsOptions{})...)
		api.ReminderConfig = &client.ApprovalReminderConfigAPI{
			Enabled:                knownBool(rc.Enabled),
			DaysUntilFirstReminder: knownInt64(rc.DaysUntilFirstReminder),
			ReminderCronSchedule:   knownString(rc.ReminderCronSchedule),
			MaxReminders:           knownInt64(rc.MaxReminders),
		}
	}

	if !m.EscalationConfig.IsNull() && !m.EscalationConfig.IsUnknown() {
		var ec escalationConfigModel
		diags.Append(m.EscalationConfig.As(ctx, &ec, basetypes.ObjectAsOptions{})...)
		chain, d := chainToAPI(ctx, ec.EscalationChain)
		diags.Append(d...)
		api.EscalationConfig = &client.ApprovalEscalationConfigAPI{
			Enabled:                  knownBool(ec.Enabled),
			DaysUntilFirstEscalation: knownInt64(ec.DaysUntilFirstEscalation),
			EscalationCronSchedule:   knownString(ec.EscalationCronSchedule),
			EscalationChain:          chain,
		}
	}

	if !m.TimeoutConfig.IsNull() && !m.TimeoutConfig.IsUnknown() {
		var tc timeoutConfigModel
		diags.Append(m.TimeoutConfig.As(ctx, &tc, basetypes.ObjectAsOptions{})...)
		api.TimeoutConfig = &client.ApprovalTimeoutConfigAPI{
			Enabled:          knownBool(tc.Enabled),
			DaysUntilTimeout: knownInt64(tc.DaysUntilTimeout),
			TimeoutResult:    knownString(tc.TimeoutResult),
		}
	}

	if !m.CronTimezone.IsNull() && !m.CronTimezone.IsUnknown() {
		var tz cronTimezoneModel
		diags.Append(m.CronTimezone.As(ctx, &tz, basetypes.ObjectAsOptions{})...)
		api.CronTimezone = &client.ApprovalCronTimezoneAPI{
			Location: knownString(tz.Location),
			Offset:   knownString(tz.Offset),
		}
	}

	serial, d := chainToAPI(ctx, m.SerialChain)
	diags.Append(d...)
	api.SerialChain = serial

	if !m.FallbackApprover.IsNull() && !m.FallbackApprover.IsUnknown() {
		var fb fallbackApproverModel
		diags.Append(m.FallbackApprover.As(ctx, &fb, basetypes.ObjectAsOptions{})...)
		api.FallbackApprover = &client.ApprovalFallbackApproverAPI{
			Type:       fb.Type.ValueString(),
			IdentityID: fb.IdentityID.ValueString(),
		}
	}

	return api, diags
}
