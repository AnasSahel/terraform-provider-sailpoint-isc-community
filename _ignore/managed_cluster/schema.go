package managedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The ID of the managed cluster.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the managed cluster.",
			},
			"pod": schema.StringAttribute{
				Computed:    true,
				Description: "The pod of the managed cluster.",
			},
			"org": schema.StringAttribute{
				Computed:    true,
				Description: "The organization of the managed cluster.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the managed cluster.",
			},
			"configuration": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The configuration of the managed cluster.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the managed cluster.",
			},
			"client_type": schema.StringAttribute{
				Computed:    true,
				Description: "The client type of the managed cluster.",
			},
			"ccg_version": schema.StringAttribute{
				Computed:    true,
				Description: "The CCG version of the managed cluster.",
			},
			"pinned_config": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the managed cluster is using a pinned configuration.",
			},
			"operational": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the managed cluster is operational.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the managed cluster.",
			},
			"public_key_certificate": schema.StringAttribute{
				Computed:    true,
				Description: "The public key certificate of the managed cluster.",
			},
			"public_key_thumbprint": schema.StringAttribute{
				Computed:    true,
				Description: "The public key thumbprint of the managed cluster.",
			},
			"public_key_type": schema.StringAttribute{
				Computed:    true,
				Description: "The public key type of the managed cluster.",
			},
			"alert_key": schema.StringAttribute{
				Computed:    true,
				Description: "The alert key of the managed cluster.",
			},
			"client_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The client IDs associated with the managed cluster.",
			},
			"service_count": schema.Int32Attribute{
				Computed:    true,
				Description: "The number of services running in the managed cluster.",
			},
			"cc_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Cloud Control ID of the managed cluster.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The creation timestamp of the managed cluster.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The last updated timestamp of the managed cluster.",
			},
		},
	}
}

func (r *ManagedClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
