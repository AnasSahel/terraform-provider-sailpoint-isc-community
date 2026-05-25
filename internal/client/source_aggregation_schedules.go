// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	sourceScheduleEndpointList   = "/v2025/sources/{sourceId}/schedules"
	sourceScheduleEndpointGet    = "/v2025/sources/{sourceId}/schedules/{scheduleType}"
	sourceScheduleEndpointCreate = "/v2025/sources/{sourceId}/schedules"
	sourceScheduleEndpointUpdate = "/v2025/sources/{sourceId}/schedules/{scheduleType}"
	sourceScheduleEndpointDelete = "/v2025/sources/{sourceId}/schedules/{scheduleType}"
)

// ScheduleHoursAPI represents the hours component of an aggregation schedule.
type ScheduleHoursAPI struct {
	Type     string   `json:"type"`
	Values   []string `json:"values,omitempty"`
	Interval *int64   `json:"interval,omitempty"`
}

// ScheduleDaysAPI represents the days component of an aggregation schedule.
type ScheduleDaysAPI struct {
	Type     string   `json:"type"`
	Values   []string `json:"values,omitempty"`
	Interval *int64   `json:"interval,omitempty"`
}

// ScheduleMonthsAPI represents the months component of an aggregation schedule.
type ScheduleMonthsAPI struct {
	Type     string   `json:"type"`
	Values   []string `json:"values,omitempty"`
	Interval *int64   `json:"interval,omitempty"`
}

// ScheduleAPI represents the schedule timing configuration.
type ScheduleAPI struct {
	Type       string             `json:"type"`
	Hours      ScheduleHoursAPI   `json:"hours"`
	Days       *ScheduleDaysAPI   `json:"days,omitempty"`
	Months     *ScheduleMonthsAPI `json:"months,omitempty"`
	TimeZoneID *string            `json:"timeZoneId,omitempty"`
}

// SourceAggregationScheduleAPI represents a SailPoint source aggregation schedule from the API.
type SourceAggregationScheduleAPI struct {
	Type              string      `json:"type"`
	Schedule          ScheduleAPI `json:"schedule"`
	EnableThreshold   bool        `json:"enableThreshold"`
	Threshold         *int64      `json:"threshold,omitempty"`
	CredentialsInBody *bool       `json:"credentialsInBody,omitempty"`
}

// sourceScheduleErrorContext provides context for error messages.
type sourceScheduleErrorContext struct {
	Operation    string
	SourceID     string
	ScheduleType string
	ResponseBody string
}

// ListSourceAggregationSchedules retrieves all aggregation schedules for a source.
func (c *Client) ListSourceAggregationSchedules(ctx context.Context, sourceID string) ([]SourceAggregationScheduleAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	tflog.Debug(ctx, "Listing source aggregation schedules", map[string]any{
		"source_id": sourceID,
	})

	var schedules []SourceAggregationScheduleAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&schedules).
		SetPathParam("sourceId", sourceID).
		Get(sourceScheduleEndpointList)

	if resp != nil && resp.IsError() {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "list", SourceID: sourceID, ResponseBody: string(resp.Bytes())},
			nil,
			resp.StatusCode(),
		)
	}

	if err != nil {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "list", SourceID: sourceID},
			err,
			0,
		)
	}

	tflog.Debug(ctx, "Successfully listed source aggregation schedules", map[string]any{
		"source_id": sourceID,
		"count":     len(schedules),
	})

	return schedules, nil
}

// GetSourceAggregationSchedule retrieves a specific aggregation schedule by type.
func (c *Client) GetSourceAggregationSchedule(ctx context.Context, sourceID, scheduleType string) (*SourceAggregationScheduleAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if scheduleType == "" {
		return nil, fmt.Errorf("schedule type cannot be empty")
	}

	tflog.Debug(ctx, "Getting source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	var schedule SourceAggregationScheduleAPI

	resp, err := c.prepareRequest(ctx).
		SetResult(&schedule).
		SetPathParam("sourceId", sourceID).
		SetPathParam("scheduleType", scheduleType).
		Get(sourceScheduleEndpointGet)

	if resp != nil && resp.IsError() {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{
				Operation: "get", SourceID: sourceID, ScheduleType: scheduleType,
				ResponseBody: string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	if err != nil {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "get", SourceID: sourceID, ScheduleType: scheduleType},
			err,
			0,
		)
	}

	tflog.Debug(ctx, "Successfully retrieved source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	return &schedule, nil
}

// CreateSourceAggregationSchedule creates a new aggregation schedule for a source.
func (c *Client) CreateSourceAggregationSchedule(ctx context.Context, sourceID string, schedule *SourceAggregationScheduleAPI) (*SourceAggregationScheduleAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if schedule == nil {
		return nil, fmt.Errorf("schedule cannot be nil")
	}

	tflog.Debug(ctx, "Creating source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": schedule.Type,
	})

	var result SourceAggregationScheduleAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(schedule).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		Post(sourceScheduleEndpointCreate)

	if resp != nil && resp.IsError() {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{
				Operation: "create", SourceID: sourceID, ScheduleType: schedule.Type,
				ResponseBody: string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	if err != nil {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "create", SourceID: sourceID, ScheduleType: schedule.Type},
			err,
			0,
		)
	}

	tflog.Info(ctx, "Successfully created source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": result.Type,
	})

	return &result, nil
}

// UpdateSourceAggregationSchedule performs a full replacement (PUT) of an aggregation schedule.
func (c *Client) UpdateSourceAggregationSchedule(ctx context.Context, sourceID, scheduleType string, schedule *SourceAggregationScheduleAPI) (*SourceAggregationScheduleAPI, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("source ID cannot be empty")
	}

	if scheduleType == "" {
		return nil, fmt.Errorf("schedule type cannot be empty")
	}

	if schedule == nil {
		return nil, fmt.Errorf("schedule cannot be nil")
	}

	tflog.Debug(ctx, "Updating source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	var result SourceAggregationScheduleAPI

	resp, err := c.prepareRequest(ctx).
		SetBody(schedule).
		SetResult(&result).
		SetPathParam("sourceId", sourceID).
		SetPathParam("scheduleType", scheduleType).
		Put(sourceScheduleEndpointUpdate)

	if resp != nil && resp.IsError() {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{
				Operation: "update", SourceID: sourceID, ScheduleType: scheduleType,
				ResponseBody: string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	if err != nil {
		return nil, c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "update", SourceID: sourceID, ScheduleType: scheduleType},
			err,
			0,
		)
	}

	tflog.Info(ctx, "Successfully updated source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	return &result, nil
}

// DeleteSourceAggregationSchedule removes an aggregation schedule from a source.
func (c *Client) DeleteSourceAggregationSchedule(ctx context.Context, sourceID, scheduleType string) error {
	if sourceID == "" {
		return fmt.Errorf("source ID cannot be empty")
	}

	if scheduleType == "" {
		return fmt.Errorf("schedule type cannot be empty")
	}

	tflog.Debug(ctx, "Deleting source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	resp, err := c.prepareRequest(ctx).
		SetPathParam("sourceId", sourceID).
		SetPathParam("scheduleType", scheduleType).
		Delete(sourceScheduleEndpointDelete)

	if resp != nil && resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			tflog.Debug(ctx, "Source aggregation schedule not found, treating as already deleted", map[string]any{
				"source_id":     sourceID,
				"schedule_type": scheduleType,
			})
			return nil
		}

		return c.formatSourceScheduleError(
			sourceScheduleErrorContext{
				Operation: "delete", SourceID: sourceID, ScheduleType: scheduleType,
				ResponseBody: string(resp.Bytes()),
			},
			nil,
			resp.StatusCode(),
		)
	}

	if err != nil {
		return c.formatSourceScheduleError(
			sourceScheduleErrorContext{Operation: "delete", SourceID: sourceID, ScheduleType: scheduleType},
			err,
			0,
		)
	}

	tflog.Info(ctx, "Successfully deleted source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	return nil
}

// formatSourceScheduleError formats errors with appropriate context for aggregation schedule operations.
func (c *Client) formatSourceScheduleError(errCtx sourceScheduleErrorContext, err error, statusCode int) error {
	var baseMsg string

	if errCtx.ScheduleType != "" {
		baseMsg = fmt.Sprintf("failed to %s aggregation schedule %q for source %q",
			errCtx.Operation, errCtx.ScheduleType, errCtx.SourceID)
	} else {
		baseMsg = fmt.Sprintf("failed to %s aggregation schedules for source %q",
			errCtx.Operation, errCtx.SourceID)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	if statusCode != 0 {
		detail := ""
		if errCtx.ResponseBody != "" {
			detail = fmt.Sprintf(" - response: %s", errCtx.ResponseBody)
		}

		switch statusCode {
		case http.StatusBadRequest:
			return fmt.Errorf("%s: invalid request (400)%s", baseMsg, detail)
		case http.StatusUnauthorized:
			return fmt.Errorf("%s: authentication failed (401)%s", baseMsg, detail)
		case http.StatusForbidden:
			return fmt.Errorf("%s: access denied (403)%s", baseMsg, detail)
		case http.StatusNotFound:
			return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
		case http.StatusConflict:
			return fmt.Errorf("%s: conflict (409)%s", baseMsg, detail)
		case http.StatusTooManyRequests:
			return fmt.Errorf("%s: rate limit exceeded (429)%s", baseMsg, detail)
		case http.StatusInternalServerError:
			return fmt.Errorf("%s: server error (500)%s", baseMsg, detail)
		default:
			return fmt.Errorf("%s: unexpected status code %d%s", baseMsg, statusCode, detail)
		}
	}

	return fmt.Errorf("%s: unknown error", baseMsg)
}
