// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &identityDataSource{}
	_ datasource.DataSourceWithConfigure = &identityDataSource{}
)

type identityDataSource struct {
	client *client.Client
}

// identityDSModel is the datasource model for a SailPoint Identity.
type identityDSModel struct {
	// Inputs — exactly one must be set.
	ID      types.String `tfsdk:"id"`
	Filters types.String `tfsdk:"filters"`

	// Outputs.
	Name            types.String             `tfsdk:"name"`
	Alias           types.String             `tfsdk:"alias"`
	EmailAddress    types.String             `tfsdk:"email_address"`
	Status          types.String             `tfsdk:"status"`
	EmployeeNumber  types.String             `tfsdk:"employee_number"`
	IsManager       types.Bool               `tfsdk:"is_manager"`
	Manager         *common.ObjectRefModel   `tfsdk:"manager"`
	IdentityProfile *identityProfileRefModel `tfsdk:"identity_profile"`
	Source          *common.ObjectRefModel   `tfsdk:"source"`
	LifecycleState  *identityLCSModel        `tfsdk:"lifecycle_state"`
	Attributes      types.Map                `tfsdk:"attributes"`
	Created         types.String             `tfsdk:"created"`
	Modified        types.String             `tfsdk:"modified"`
}

type identityProfileRefModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type identityLCSModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	StateName types.String `tfsdk:"state_name"`
}

// NewIdentityDataSource creates a new data source for SailPoint Identity.
func NewIdentityDataSource() datasource.DataSource {
	return &identityDataSource{}
}

func (d *identityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity"
}

func (d *identityDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "identity data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *identityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Identity.",
		MarkdownDescription: "Data source for SailPoint Identity. " +
			"Look up a full identity by `id` or by `filters` (ISC filter expression). " +
			"Exactly one of `id` or `filters` must be provided. " +
			"Identities are not created by Terraform — they are derived from source aggregation. " +
			"Use this data source to resolve an identity by alias or name and reference its ID in " +
			"`owner` fields on other resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the identity. Mutually exclusive with `filters`.",
			},
			"filters": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ISC filter expression (e.g. `alias eq \"jsmith\"` or `name eq \"Jane Smith\"`). Mutually exclusive with `id`. Returns an error if zero or more than one identity matches.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the identity.",
			},
			"alias": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The login alias (username) of the identity.",
			},
			"email_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The email address of the identity.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the identity (e.g. `ACTIVE`, `INACTIVE`).",
			},
			"employee_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The employee number of the identity.",
			},
			"is_manager": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the identity is a manager.",
			},
			"manager": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The manager of the identity.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"identity_profile": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The identity profile this identity belongs to.",
				Attributes: map[string]schema.Attribute{
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"source": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The authoritative source for this identity.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"lifecycle_state": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The current lifecycle state of the identity.",
				Attributes: map[string]schema.Attribute{
					"id":         schema.StringAttribute{Computed: true},
					"name":       schema.StringAttribute{Computed: true},
					"state_name": schema.StringAttribute{Computed: true},
				},
			},
			"attributes": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				MarkdownDescription: "Additional identity attributes as key-value pairs. Keys depend on the org's identity schema. " +
					"Scalar values are returned as plain strings; multi-valued or complex values are JSON-encoded strings.",
			},
			"created":  schema.StringAttribute{Computed: true},
			"modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *identityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config identityDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueString() != ""
	hasFilters := !config.Filters.IsNull() && !config.Filters.IsUnknown() && config.Filters.ValueString() != ""

	if !hasID && !hasFilters {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"One of `id` or `filters` must be provided.",
		)
		return
	}
	if hasID && hasFilters {
		resp.Diagnostics.AddError(
			"Conflicting arguments",
			"`id` and `filters` are mutually exclusive. Provide only one.",
		)
		return
	}

	var apiResp *client.IdentityAPI

	if hasID {
		tflog.Debug(ctx, "Reading identity data source by ID", map[string]any{"id": config.ID.ValueString()})
		identity, err := d.client.GetIdentity(ctx, config.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Identity",
				fmt.Sprintf("Could not read identity %q: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}
		apiResp = identity
	} else {
		tflog.Debug(ctx, "Reading identity data source by filters", map[string]any{"filters": config.Filters.ValueString()})
		// Fetch up to 2 results — we only accept exactly 1.
		identities, err := d.client.ListIdentities(ctx, config.Filters.ValueString(), 2)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Identities",
				fmt.Sprintf("Could not list identities with filters %q: %s", config.Filters.ValueString(), err.Error()),
			)
			return
		}
		switch len(identities) {
		case 0:
			resp.Diagnostics.AddError(
				"No identity found",
				fmt.Sprintf("No identity matched filters %q. Verify the filter expression.", config.Filters.ValueString()),
			)
			return
		case 1:
			apiResp = &identities[0]
		default:
			resp.Diagnostics.AddError(
				"Multiple identities found",
				fmt.Sprintf("filters %q matched more than one identity. Use a more specific expression or supply `id` directly.", config.Filters.ValueString()),
			)
			return
		}
	}

	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Identity", "Received nil response from SailPoint API")
		return
	}

	resp.Diagnostics.Append(populateIdentityDSModel(ctx, &config, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func populateIdentityDSModel(ctx context.Context, m *identityDSModel, api *client.IdentityAPI) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	m.ID = types.StringValue(api.ID)
	m.Name = types.StringValue(api.Name)
	m.Alias = types.StringValue(api.Alias)
	m.EmailAddress = stringPtrToTF(api.EmailAddress)
	m.Status = stringPtrToTF(api.Status)
	m.EmployeeNumber = stringPtrToTF(api.EmployeeNumber)
	m.Created = stringPtrToTF(api.Created)
	m.Modified = stringPtrToTF(api.Modified)

	if api.IsManager != nil {
		m.IsManager = types.BoolValue(*api.IsManager)
	} else {
		m.IsManager = types.BoolNull()
	}

	if api.Manager != nil {
		ref, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Manager)
		diagnostics.Append(diags...)
		m.Manager = ref
	} else {
		m.Manager = nil
	}

	if api.IdentityProfile != nil {
		m.IdentityProfile = &identityProfileRefModel{
			ID:   types.StringValue(api.IdentityProfile.ID),
			Name: types.StringValue(api.IdentityProfile.Name),
		}
	} else {
		m.IdentityProfile = nil
	}

	if api.Source != nil {
		ref, diags := common.NewObjectRefFromAPIPtr(ctx, *api.Source)
		diagnostics.Append(diags...)
		m.Source = ref
	} else {
		m.Source = nil
	}

	if api.LifecycleState != nil {
		lcs := &identityLCSModel{}
		lcs.ID = stringPtrToTF(api.LifecycleState.ID)
		lcs.Name = stringPtrToTF(api.LifecycleState.Name)
		lcs.StateName = stringPtrToTF(api.LifecycleState.StateName)
		m.LifecycleState = lcs
	} else {
		m.LifecycleState = nil
	}

	if len(api.Attributes) > 0 {
		attrStrings, diags := identityAttributesToStringMap(api.Attributes)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return diagnostics
		}
		attrMap, diags := types.MapValueFrom(ctx, types.StringType, attrStrings)
		diagnostics.Append(diags...)
		m.Attributes = attrMap
	} else {
		m.Attributes = types.MapNull(types.StringType)
	}

	return diagnostics
}

func identityAttributesToStringMap(attrs map[string]interface{}) (map[string]string, diag.Diagnostics) {
	if len(attrs) == 0 {
		return nil, diag.Diagnostics{}
	}

	result := make(map[string]string, len(attrs))
	var diagnostics diag.Diagnostics

	for key, value := range attrs {
		str, diags := identityAttributeValueToString(value)
		diagnostics.Append(diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}
		result[key] = str
	}

	return result, diagnostics
}

func identityAttributeValueToString(value interface{}) (string, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value == nil {
		return "", diagnostics
	}

	switch val := value.(type) {
	case string:
		return val, diagnostics
	case float64:
		if val == math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val) {
			return strconv.FormatInt(int64(val), 10), diagnostics
		}
		return strconv.FormatFloat(val, 'f', -1, 64), diagnostics
	case bool:
		return strconv.FormatBool(val), diagnostics
	default:
		bytes, err := json.Marshal(val)
		if err != nil {
			diagnostics.AddError(
				"Error Converting Identity Attribute",
				fmt.Sprintf("Could not serialize attribute value: %s", err.Error()),
			)
			return "", diagnostics
		}
		return string(bytes), diagnostics
	}
}

func stringPtrToTF(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
