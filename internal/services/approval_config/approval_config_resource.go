// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package approval_config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &approvalConfigResource{}
	_ resource.ResourceWithConfigure      = &approvalConfigResource{}
	_ resource.ResourceWithImportState    = &approvalConfigResource{}
	_ resource.ResourceWithValidateConfig = &approvalConfigResource{}
)

// supportedScopes are the scopes the resource manages (design D2).
var supportedScopes = map[string]struct{}{
	"ROLE":           {},
	"ACCESS_PROFILE": {},
	"ENTITLEMENT":    {},
	"TENANT":         {},
}

const resourceDescription = "Manages the approval reminders, escalation, timeout and fallback approver of a role, access profile, entitlement, or the whole tenant " +
	"(`/generic-approvals/v1/config`, experimental). " +
	"Create is an upsert: if the object already has its own approval config, Terraform takes it over. " +
	"Destroy deletes the object's own config so it inherits the tenant config again; for `scope = \"TENANT\"` destroy only removes the resource from state. " +
	"Every setting is optional: what you leave out keeps the value SailPoint has (or its default). " +
	"Removing a block from the configuration keeps its last value; set `enabled = false` to turn a feature off. " +
	"Requires `ORG_ADMIN`."

type approvalConfigResource struct {
	client *client.Client
}

// NewApprovalConfigResource creates a new resource for generic approval configuration.
func NewApprovalConfigResource() resource.Resource {
	return &approvalConfigResource{}
}

func (r *approvalConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_approval_config"
}

func (r *approvalConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "approval config resource")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = c
}

func chainAttribute(markdown string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: markdown,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"identity_type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Approver type, e.g. `MANAGER_OF`, `IDENTITY`, `GOVERNANCE_GROUP`, `ROLE_OWNER`, `ACCESS_PROFILE_OWNER`, `ENTITLEMENT_OWNER`, `SOURCE_OWNER`.",
				},
				"identity_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "ID of the identity or governance group, when `identity_type` is `IDENTITY` or `GOVERNANCE_GROUP`.",
				},
			},
		},
	}
}

func optStringAttr(markdown string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: markdown,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	}
}

func optBoolAttr(markdown string) schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: markdown,
		PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
	}
}

func optInt64Attr(markdown string) schema.Int64Attribute {
	return schema.Int64Attribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: markdown,
		PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
}

func optObjectPlanModifiers() []planmodifier.Object {
	return []planmodifier.Object{objectplanmodifier.UseStateForUnknown()}
}

func (r *approvalConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         resourceDescription,
		MarkdownDescription: resourceDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "`<scope>/<object_id>`, e.g. `ROLE/2c9180...`. Also the import ID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"scope": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "What `object_id` points to: `ROLE`, `ACCESS_PROFILE`, `ENTITLEMENT` or `TENANT`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"object_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the role, access profile or entitlement. For `TENANT`, the tenant ID.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"reminder_config": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reminders sent to approvers while a request is pending.",
				PlanModifiers:       optObjectPlanModifiers(),
				Attributes: map[string]schema.Attribute{
					"enabled":                   optBoolAttr("Whether reminders are sent."),
					"days_until_first_reminder": optInt64Attr("Days before the first reminder."),
					"reminder_cron_schedule":    optStringAttr("Schedule of the following reminders, in robfig/cron syntax (e.g. `@every 48h`)."),
					"max_reminders":             optInt64Attr("Maximum number of reminders (at most 20)."),
				},
			},
			"escalation_config": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Escalation of pending requests to other approvers.",
				PlanModifiers:       optObjectPlanModifiers(),
				Attributes: map[string]schema.Attribute{
					"enabled":                     optBoolAttr("Whether escalation is enabled."),
					"days_until_first_escalation": optInt64Attr("Days before the first escalation."),
					"escalation_cron_schedule":    optStringAttr("Schedule of the following escalations, in robfig/cron syntax (e.g. `@every 72h`)."),
					"escalation_chain":            chainAttribute("Who the request escalates to, in order (the first item is tier 1)."),
				},
			},
			"timeout_config": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When a pending request expires.",
				PlanModifiers:       optObjectPlanModifiers(),
				Attributes: map[string]schema.Attribute{
					"enabled":            optBoolAttr("Whether requests time out."),
					"days_until_timeout": optInt64Attr("Days before the request times out (at most 90)."),
					"timeout_result":     optStringAttr("Outcome on timeout: `EXPIRED` or `APPROVED`."),
				},
			},
			"cron_timezone": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Timezone for the reminder and escalation schedules.",
				PlanModifiers:       optObjectPlanModifiers(),
				Attributes: map[string]schema.Attribute{
					"location": optStringAttr("IANA timezone, e.g. `Europe/Paris`."),
					"offset":   optStringAttr("Timezone offset, if no location is set."),
				},
			},
			"serial_chain": chainAttribute("Approvers in order, used when the request's approval criteria is `SERIAL`."),
			"requires_comment": optStringAttr(
				"When approvers must comment: `APPROVAL`, `REJECTION`, `ALL` or `OFF`."),
			"fallback_approver": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Approver used when the assigned approver cannot be found and no escalation applies.",
				PlanModifiers:       optObjectPlanModifiers(),
				Attributes: map[string]schema.Attribute{
					"type":        optStringAttr("Approver type, same values as `identity_type`."),
					"identity_id": optStringAttr("ID of the identity or governance group, when `type` is `IDENTITY` or `GOVERNANCE_GROUP`."),
				},
			},
			"machine_identity_manager_assignment": optStringAttr(
				"How `MANAGER_OF` resolves when the requestee is a machine identity (default `MANAGER_OF_REQUESTER`)."),
			"circumvent_approval_process": optBoolAttr(
				"When `true`, approvals are created as passed, skipping approval. Intended for machine account approvals only."),
			"auto_approve": optStringAttr(
				"`OFF`, `DIRECT` or `INDIRECT`: whether steps assigned to the requester are auto-approved."),
		},
	}
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *approvalConfigResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var scope types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("scope"), &scope)...)
	if resp.Diagnostics.HasError() || scope.IsNull() || scope.IsUnknown() {
		return
	}
	if _, ok := supportedScopes[scope.ValueString()]; !ok {
		resp.Diagnostics.AddAttributeError(path.Root("scope"), "Invalid scope",
			fmt.Sprintf("scope must be one of ROLE, ACCESS_PROFILE, ENTITLEMENT, TENANT; got %q.", scope.ValueString()))
	}
}

// save PUTs the planned config, reads it back and stores the result in the model.
// summary is the diagnostic title ("Error Creating Approval Config" / "Error Updating Approval Config").
func (r *approvalConfigResource) save(ctx context.Context, plan *approvalConfigModel, summary string) diag.Diagnostics {
	var diags diag.Diagnostics
	scope, objectID := plan.Scope.ValueString(), plan.ObjectID.ValueString()

	body, d := plan.ToAPI(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	if _, err := r.client.PutApprovalConfig(ctx, objectID, scope, body); err != nil {
		diags.AddError(summary, fmt.Sprintf("Could not set approval config for %s %q: %s", scope, objectID, err))
		return diags
	}

	got, err := r.client.GetApprovalConfig(ctx, objectID, scope)
	if err != nil {
		diags.AddError(summary, fmt.Sprintf("Could not read approval config back after saving it for %s %q: %s", scope, objectID, err))
		return diags
	}
	if !isOwnConfig(got, objectID) {
		diags.AddError(summary, fmt.Sprintf(
			"SailPoint accepted the approval config for %s %q but returns an inherited config (scope %s, id %s) when reading it back. Check that object_id is the ID of a %s.",
			scope, objectID, deref(got.Scope), deref(got.ID), scope))
		return diags
	}

	diags.Append(plan.FromAPI(ctx, got)...)
	return diags
}

func (r *approvalConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan approvalConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.save(ctx, &plan, "Error Creating Approval Config")...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Created approval config", map[string]any{"scope": plan.Scope.ValueString(), "object_id": plan.ObjectID.ValueString()})
}

func (r *approvalConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state approvalConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	scope, objectID := state.Scope.ValueString(), state.ObjectID.ValueString()

	got, err := r.client.GetApprovalConfig(ctx, objectID, scope)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			tflog.Info(ctx, "Approval config not found, removing from state", map[string]any{"scope": scope, "object_id": objectID})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Approval Config",
			fmt.Sprintf("Could not read approval config for %s %q: %s", scope, objectID, err))
		return
	}

	if !isOwnConfig(got, objectID) {
		tflog.Info(ctx, "Object has no approval config of its own anymore, removing from state", map[string]any{
			"scope":          scope,
			"object_id":      objectID,
			"returned_scope": deref(got.Scope),
			"returned_id":    deref(got.ID),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(state.FromAPI(ctx, got)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *approvalConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state approvalConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Scope = state.Scope
	plan.ObjectID = state.ObjectID

	resp.Diagnostics.Append(r.save(ctx, &plan, "Error Updating Approval Config")...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Updated approval config", map[string]any{"scope": plan.Scope.ValueString(), "object_id": plan.ObjectID.ValueString()})
}

func (r *approvalConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state approvalConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	scope, objectID := state.Scope.ValueString(), state.ObjectID.ValueString()

	if scope == "TENANT" {
		resp.Diagnostics.AddWarning(
			"Tenant approval config not deleted from SailPoint ISC",
			"The tenant approval config is the default every role, access profile and entitlement inherits, so Terraform does not delete it. "+
				"The resource has been removed from Terraform state; the config is unchanged in ISC.")
		return
	}

	for _, s := range []string{scope, "DOMAIN_OBJECT"} {
		if err := r.client.DeleteApprovalConfig(ctx, objectID, s); err != nil { // 404 already returns nil
			resp.Diagnostics.AddError("Error Deleting Approval Config",
				fmt.Sprintf("Could not delete approval config for %s %q: %s", scope, objectID, err))
			return
		}
		got, err := r.client.GetApprovalConfig(ctx, objectID, scope)
		if errors.Is(err, client.ErrNotFound) {
			return // the object itself is gone
		}
		if err != nil {
			resp.Diagnostics.AddError("Error Deleting Approval Config",
				fmt.Sprintf("Could not check approval config for %s %q after deleting it: %s", scope, objectID, err))
			return
		}
		if !isOwnConfig(got, objectID) {
			tflog.Info(ctx, "Deleted approval config", map[string]any{"scope": scope, "object_id": objectID})
			return
		}
		tflog.Debug(ctx, "Approval config still present after delete, retrying with DOMAIN_OBJECT scope", map[string]any{"scope": scope, "object_id": objectID})
	}
	resp.Diagnostics.AddError("Error Deleting Approval Config",
		fmt.Sprintf("SailPoint still returns an approval config of its own for %s %q after deleting it with scopes %s and DOMAIN_OBJECT.", scope, objectID, scope))
}

// ImportState implements resource.ResourceWithImportState.
// Import format: <scope>/<object_id> (e.g. ROLE/2c9180...).
func (r *approvalConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("Expected import ID format: <scope>/<object_id> (e.g. ROLE/2c9180835d2e5168015d32f890ca1581), got: %s", req.ID))
		return
	}
	if _, ok := supportedScopes[parts[0]]; !ok {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("Unsupported scope %q in import ID; expected ROLE, ACCESS_PROFILE, ENTITLEMENT or TENANT.", parts[0]))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
