// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package source_aggregation_schedule

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sourceAggregationScheduleResource{}
	_ resource.ResourceWithConfigure   = &sourceAggregationScheduleResource{}
	_ resource.ResourceWithImportState = &sourceAggregationScheduleResource{}
)

type sourceAggregationScheduleResource struct {
	client *client.Client
}

// NewSourceAggregationScheduleResource creates a new resource for Source Aggregation Schedule.
func NewSourceAggregationScheduleResource() resource.Resource {
	return &sourceAggregationScheduleResource{}
}

// Metadata implements resource.Resource.
func (r *sourceAggregationScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_aggregation_schedule"
}

// Configure implements resource.ResourceWithConfigure.
func (r *sourceAggregationScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "source aggregation schedule resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

// timeUnitSchema returns the nested attribute schema for hours, days, and months.
func timeUnitSchema(description string, required bool) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: description,
		Required:            required,
		Optional:            !required,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of schedule unit. Valid values: `DAILY`, `WEEKLY`, `MONTHLY`, `CALENDAR`, `LIST`, `RANGE`.",
				Required:            true,
			},
			"values": schema.ListAttribute{
				MarkdownDescription: "The list of values for the schedule unit (e.g. `[\"MON\", \"WED\", \"FRI\"]` for days, `[\"3\"]` for hour 3 AM).",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "The interval for RANGE-type schedules (e.g. every N days).",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Schema implements resource.Resource.
func (r *sourceAggregationScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a SailPoint source aggregation schedule.",
		MarkdownDescription: "Manages a SailPoint source aggregation schedule. " +
			"Aggregation schedules determine when a source performs account and entitlement aggregations. " +
			"Each source supports up to two schedule types: `ACCOUNTS` and `ENTITLEMENTS`.\n\n" +
			"The resource ID is a composite of `<source_id>/<schedule_type>`. " +
			"To import: `terraform import sailpoint_source_aggregation_schedule.example <source_id>/<schedule_type>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The composite resource ID (`<source_id>/<schedule_type>`). Computed by the provider.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the source this schedule belongs to. Changing this forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule_type": schema.StringAttribute{
				MarkdownDescription: "The type of aggregation schedule. Valid values: `ACCOUNTS`, `ENTITLEMENTS`. Changing this forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule": schema.SingleNestedAttribute{
				MarkdownDescription: "The timing configuration for this aggregation schedule.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The schedule recurrence type. Valid values: `DAILY`, `WEEKLY`, `MONTHLY`, `CALENDAR`.",
						Required:            true,
					},
					"hours": timeUnitSchema(
						"The hours component of the schedule — specifies the hour(s) of day when the aggregation runs.",
						true,
					),
					"days": timeUnitSchema(
						"The days component of the schedule — specifies the day(s) of week or month when the aggregation runs. Required for `WEEKLY` and `MONTHLY` schedules.",
						false,
					),
					"months": timeUnitSchema(
						"The months component of the schedule — specifies the month(s) when the aggregation runs. Required for `CALENDAR` schedules.",
						false,
					),
					"time_zone_id": schema.StringAttribute{
						MarkdownDescription: "The IANA time zone identifier for the schedule (e.g. `UTC`, `America/Chicago`). Defaults to `UTC` when not specified.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_threshold": schema.BoolAttribute{
				MarkdownDescription: "When `true`, enables delta aggregation thresholding. If the number of accounts or entitlements changed exceeds the `threshold` percentage, the aggregation is paused for review.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"threshold": schema.Int64Attribute{
				MarkdownDescription: "The percentage threshold for delta aggregation. Only meaningful when `enable_threshold` is `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"credentials_in_body": schema.BoolAttribute{
				MarkdownDescription: "Whether to include credentials in the aggregation request body. Computed from the source configuration.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *sourceAggregationScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := plan.SourceID.ValueString()
	scheduleType := plan.ScheduleType.ValueString()

	tflog.Debug(ctx, "Creating source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	apiBody, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateSourceAggregationSchedule(ctx, sourceID, apiBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SailPoint Source Aggregation Schedule",
			fmt.Sprintf("Could not create aggregation schedule %q for source %q: %s",
				scheduleType, sourceID, err.Error()),
		)
		return
	}

	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(state.FromAPI(ctx, sourceID, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully created SailPoint Source Aggregation Schedule resource", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})
}

// Read implements resource.Resource.
func (r *sourceAggregationScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	scheduleType := state.ScheduleType.ValueString()

	tflog.Debug(ctx, "Reading source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	result, err := r.client.GetSourceAggregationSchedule(ctx, sourceID, scheduleType)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Source aggregation schedule not found, removing from state", map[string]any{
				"source_id":     sourceID,
				"schedule_type": scheduleType,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Source Aggregation Schedule",
			fmt.Sprintf("Could not read aggregation schedule %q for source %q: %s",
				scheduleType, sourceID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, sourceID, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully read SailPoint Source Aggregation Schedule resource", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})
}

// Update implements resource.Resource.
func (r *sourceAggregationScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	scheduleType := state.ScheduleType.ValueString()

	tflog.Debug(ctx, "Updating source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	apiBody, diags := plan.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateSourceAggregationSchedule(ctx, sourceID, scheduleType, apiBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SailPoint Source Aggregation Schedule",
			fmt.Sprintf("Could not update aggregation schedule %q for source %q: %s",
				scheduleType, sourceID, err.Error()),
		)
		return
	}

	var newState sourceAggregationScheduleModel
	resp.Diagnostics.Append(newState.FromAPI(ctx, sourceID, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully updated SailPoint Source Aggregation Schedule resource", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})
}

// Delete implements resource.Resource.
func (r *sourceAggregationScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceAggregationScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceID := state.SourceID.ValueString()
	scheduleType := state.ScheduleType.ValueString()

	tflog.Debug(ctx, "Deleting source aggregation schedule", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})

	err := r.client.DeleteSourceAggregationSchedule(ctx, sourceID, scheduleType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SailPoint Source Aggregation Schedule",
			fmt.Sprintf("Could not delete aggregation schedule %q for source %q: %s",
				scheduleType, sourceID, err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted SailPoint Source Aggregation Schedule resource", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})
}

// ImportState implements resource.ResourceWithImportState.
// Import format: source_id/schedule_type (e.g. abc123/ACCOUNTS).
func (r *sourceAggregationScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing source aggregation schedule resource", map[string]any{
		"import_id": req.ID,
	})

	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: source_id/schedule_type (e.g. abc123/ACCOUNTS), got: %s", req.ID),
		)
		return
	}

	sourceID := parts[0]
	scheduleType := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), sourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("schedule_type"), scheduleType)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)

	tflog.Info(ctx, "Successfully imported SailPoint Source Aggregation Schedule resource", map[string]any{
		"source_id":     sourceID,
		"schedule_type": scheduleType,
	})
}
