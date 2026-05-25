// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package governance_group

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &governanceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &governanceGroupDataSource{}
)

type governanceGroupDataSource struct {
	client *client.Client
}

// governanceGroupDSModel is the datasource-specific model. It extends the
// resource model with an optional filters input field.
type governanceGroupDSModel struct {
	ID          types.String                 `tfsdk:"id"`
	Filters     types.String                 `tfsdk:"filters"`
	Name        types.String                 `tfsdk:"name"`
	Description types.String                 `tfsdk:"description"`
	Owner       *common.ObjectRefModel       `tfsdk:"owner"`
	Members     []governanceGroupMemberModel `tfsdk:"members"`
	MemberCount types.Int64                  `tfsdk:"member_count"`
	Created     types.String                 `tfsdk:"created"`
	Modified    types.String                 `tfsdk:"modified"`
}

// NewGovernanceGroupDataSource creates a new data source for SailPoint Governance Group.
func NewGovernanceGroupDataSource() datasource.DataSource {
	return &governanceGroupDataSource{}
}

func (d *governanceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_governance_group"
}

func (d *governanceGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	c, diags := common.ConfigureClient(ctx, req.ProviderData, "governance group data source")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.client = c
}

func (d *governanceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for SailPoint Governance Group.",
		MarkdownDescription: "Data source for SailPoint Governance Group. " +
			"Look up a governance group by `id` or by `filters` (ISC filter expression). " +
			"Exactly one of `id` or `filters` must be provided.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the governance group. Mutually exclusive with `filters`.",
			},
			"filters": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ISC filter expression to look up the group (e.g. `name eq \"example-group\"`). Mutually exclusive with `id`.",
			},
			"name":        schema.StringAttribute{Computed: true},
			"description": schema.StringAttribute{Computed: true},
			"owner": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"id":   schema.StringAttribute{Computed: true},
					"name": schema.StringAttribute{Computed: true},
				},
			},
			"members": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{Computed: true},
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
			"member_count": schema.Int64Attribute{Computed: true},
			"created":      schema.StringAttribute{Computed: true},
			"modified":     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *governanceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config governanceGroupDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate: exactly one of id or filters must be set.
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

	var apiResp *client.GovernanceGroupAPI

	if hasID {
		tflog.Debug(ctx, "Reading governance group data source by ID", map[string]any{"id": config.ID.ValueString()})
		gg, err := d.client.GetGovernanceGroup(ctx, config.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading SailPoint Governance Group",
				fmt.Sprintf("Could not read governance group %q: %s", config.ID.ValueString(), err.Error()),
			)
			return
		}
		apiResp = gg
	} else {
		tflog.Debug(ctx, "Reading governance group data source by filters", map[string]any{"filters": config.Filters.ValueString()})
		groups, err := d.client.ListGovernanceGroups(ctx, config.Filters.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Listing SailPoint Governance Groups",
				fmt.Sprintf("Could not list governance groups with filters %q: %s", config.Filters.ValueString(), err.Error()),
			)
			return
		}
		switch len(groups) {
		case 0:
			resp.Diagnostics.AddError(
				"No governance group found",
				fmt.Sprintf("No governance group matched filters %q. Verify the filter expression.", config.Filters.ValueString()),
			)
			return
		case 1:
			apiResp = &groups[0]
		default:
			resp.Diagnostics.AddError(
				"Multiple governance groups found",
				fmt.Sprintf("filters %q matched %d governance groups. Use a more specific expression or supply `id` directly.", config.Filters.ValueString(), len(groups)),
			)
			return
		}
	}

	if apiResp == nil {
		resp.Diagnostics.AddError("Error Reading SailPoint Governance Group", "Received nil response from SailPoint API")
		return
	}

	members, err := d.client.ListGovernanceGroupMembers(ctx, apiResp.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SailPoint Governance Group Members",
			fmt.Sprintf("Could not read members for governance group %q: %s", apiResp.ID, err.Error()),
		)
		return
	}

	// Populate the datasource model.
	config.ID = types.StringValue(apiResp.ID)
	config.Name = types.StringValue(apiResp.Name)
	config.Description = stringPtrToTF(apiResp.Description)
	config.Created = stringPtrToTF(apiResp.Created)
	config.Modified = stringPtrToTF(apiResp.Modified)

	if apiResp.MemberCount != nil {
		config.MemberCount = types.Int64Value(*apiResp.MemberCount)
	} else {
		config.MemberCount = types.Int64Value(int64(len(members)))
	}

	owner, diags := common.NewObjectRefFromAPIPtr(ctx, apiResp.Owner)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Owner = owner

	if len(members) > 0 {
		config.Members = make([]governanceGroupMemberModel, 0, len(members))
		for _, m := range members {
			config.Members = append(config.Members, governanceGroupMemberModel{
				Type: types.StringValue(m.Type),
				ID:   types.StringValue(m.ID),
				Name: types.StringValue(m.Name),
			})
		}
	} else {
		config.Members = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
