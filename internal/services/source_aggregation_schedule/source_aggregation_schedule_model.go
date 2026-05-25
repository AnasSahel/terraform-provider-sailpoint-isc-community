// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_aggregation_schedule

import (
	"context"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// scheduleTimeUnitModel maps the hours/days/months sub-objects.
type scheduleTimeUnitModel struct {
	Type     types.String `tfsdk:"type"`
	Values   types.List   `tfsdk:"values"`
	Interval types.Int64  `tfsdk:"interval"`
}

// scheduleTimeUnitAttrTypes returns the attribute type map for scheduleTimeUnitModel.
func scheduleTimeUnitAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.StringType,
		"values":   types.ListType{ElemType: types.StringType},
		"interval": types.Int64Type,
	}
}

// scheduleModel maps the schedule nested attribute.
type scheduleModel struct {
	Type       types.String `tfsdk:"type"`
	Hours      types.Object `tfsdk:"hours"`
	Days       types.Object `tfsdk:"days"`
	Months     types.Object `tfsdk:"months"`
	TimeZoneID types.String `tfsdk:"time_zone_id"`
}

// scheduleAttrTypes returns the attribute type map for scheduleModel.
func scheduleAttrTypes() map[string]attr.Type {
	timeUnitType := types.ObjectType{AttrTypes: scheduleTimeUnitAttrTypes()}
	return map[string]attr.Type{
		"type":         types.StringType,
		"hours":        timeUnitType,
		"days":         timeUnitType,
		"months":       timeUnitType,
		"time_zone_id": types.StringType,
	}
}

// sourceAggregationScheduleModel is the Terraform state model for sailpoint_source_aggregation_schedule.
type sourceAggregationScheduleModel struct {
	ID                types.String `tfsdk:"id"`
	SourceID          types.String `tfsdk:"source_id"`
	ScheduleType      types.String `tfsdk:"schedule_type"`
	Schedule          types.Object `tfsdk:"schedule"`
	EnableThreshold   types.Bool   `tfsdk:"enable_threshold"`
	Threshold         types.Int64  `tfsdk:"threshold"`
	CredentialsInBody types.Bool   `tfsdk:"credentials_in_body"`
}

// timeUnitFromAPI converts a ScheduleHoursAPI-like struct to a types.Object.
func timeUnitFromAPI(ctx context.Context, apiType string, apiValues []string, apiInterval *int64) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	valList, d := types.ListValueFrom(ctx, types.StringType, apiValues)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(scheduleTimeUnitAttrTypes()), diags
	}

	interval := types.Int64Null()
	if apiInterval != nil {
		interval = types.Int64Value(*apiInterval)
	}

	obj, d := types.ObjectValue(scheduleTimeUnitAttrTypes(), map[string]attr.Value{
		"type":     types.StringValue(apiType),
		"values":   valList,
		"interval": interval,
	})
	diags.Append(d...)
	return obj, diags
}

// FromAPI populates the Terraform model from the API response.
func (m *sourceAggregationScheduleModel) FromAPI(ctx context.Context, sourceID string, api *client.SourceAggregationScheduleAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(sourceID + "/" + api.Type)
	m.SourceID = types.StringValue(sourceID)
	m.ScheduleType = types.StringValue(api.Type)
	m.EnableThreshold = types.BoolValue(api.EnableThreshold)

	if api.Threshold != nil {
		m.Threshold = types.Int64Value(*api.Threshold)
	} else {
		m.Threshold = types.Int64Null()
	}

	if api.CredentialsInBody != nil {
		m.CredentialsInBody = types.BoolValue(*api.CredentialsInBody)
	} else {
		m.CredentialsInBody = types.BoolNull()
	}

	// Map hours
	hoursObj, d := timeUnitFromAPI(ctx, api.Schedule.Hours.Type, api.Schedule.Hours.Values, api.Schedule.Hours.Interval)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	// Map days (optional)
	daysObj := types.ObjectNull(scheduleTimeUnitAttrTypes())
	if api.Schedule.Days != nil {
		daysObj, d = timeUnitFromAPI(ctx, api.Schedule.Days.Type, api.Schedule.Days.Values, api.Schedule.Days.Interval)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	// Map months (optional)
	monthsObj := types.ObjectNull(scheduleTimeUnitAttrTypes())
	if api.Schedule.Months != nil {
		monthsObj, d = timeUnitFromAPI(ctx, api.Schedule.Months.Type, api.Schedule.Months.Values, api.Schedule.Months.Interval)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	// Map time_zone_id (optional)
	tzID := types.StringNull()
	if api.Schedule.TimeZoneID != nil {
		tzID = types.StringValue(*api.Schedule.TimeZoneID)
	}

	schedObj, d := types.ObjectValue(scheduleAttrTypes(), map[string]attr.Value{
		"type":         types.StringValue(api.Schedule.Type),
		"hours":        hoursObj,
		"days":         daysObj,
		"months":       monthsObj,
		"time_zone_id": tzID,
	})
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.Schedule = schedObj
	return diags
}

// timeUnitToAPI extracts an API hours/days/months struct from a types.Object.
// Returns nil if the object is null or unknown.
func timeUnitToAPI(ctx context.Context, obj types.Object) (*struct {
	Type     string
	Values   []string
	Interval *int64
}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var m scheduleTimeUnitModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	var values []string
	diags.Append(m.Values.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil, diags
	}

	var interval *int64
	if !m.Interval.IsNull() && !m.Interval.IsUnknown() {
		v := m.Interval.ValueInt64()
		interval = &v
	}

	result := &struct {
		Type     string
		Values   []string
		Interval *int64
	}{
		Type:     m.Type.ValueString(),
		Values:   values,
		Interval: interval,
	}
	return result, diags
}

// ToAPI converts the Terraform model to the API request body.
func (m *sourceAggregationScheduleModel) ToAPI(ctx context.Context) (*client.SourceAggregationScheduleAPI, diag.Diagnostics) {
	var diags diag.Diagnostics

	var sched scheduleModel
	diags.Append(m.Schedule.As(ctx, &sched, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Hours (required)
	hoursRaw, d := timeUnitToAPI(ctx, sched.Hours)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	hours := client.ScheduleHoursAPI{}
	if hoursRaw != nil {
		hours.Type = hoursRaw.Type
		hours.Values = hoursRaw.Values
		hours.Interval = hoursRaw.Interval
	}

	apiSchedule := client.ScheduleAPI{
		Type:  sched.Type.ValueString(),
		Hours: hours,
	}

	// Days (optional)
	daysRaw, d := timeUnitToAPI(ctx, sched.Days)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	if daysRaw != nil {
		apiSchedule.Days = &client.ScheduleDaysAPI{
			Type:     daysRaw.Type,
			Values:   daysRaw.Values,
			Interval: daysRaw.Interval,
		}
	}

	// Months (optional)
	monthsRaw, d := timeUnitToAPI(ctx, sched.Months)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	if monthsRaw != nil {
		apiSchedule.Months = &client.ScheduleMonthsAPI{
			Type:     monthsRaw.Type,
			Values:   monthsRaw.Values,
			Interval: monthsRaw.Interval,
		}
	}

	// TimeZoneID (optional)
	if !sched.TimeZoneID.IsNull() && !sched.TimeZoneID.IsUnknown() {
		tz := sched.TimeZoneID.ValueString()
		apiSchedule.TimeZoneID = &tz
	}

	api := &client.SourceAggregationScheduleAPI{
		Type:            m.ScheduleType.ValueString(),
		Schedule:        apiSchedule,
		EnableThreshold: m.EnableThreshold.ValueBool(),
	}

	if !m.Threshold.IsNull() && !m.Threshold.IsUnknown() {
		v := m.Threshold.ValueInt64()
		api.Threshold = &v
	}

	if !m.CredentialsInBody.IsNull() && !m.CredentialsInBody.IsUnknown() {
		v := m.CredentialsInBody.ValueBool()
		api.CredentialsInBody = &v
	}

	return api, diags
}
