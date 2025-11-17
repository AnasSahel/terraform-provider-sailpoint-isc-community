// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &identityAttributeResource{}
	_ resource.ResourceWithConfigure   = &identityAttributeResource{}
	_ resource.ResourceWithImportState = &identityAttributeResource{}
)

type identityAttributeResource struct {
	client *client.Client
}

func NewIdentityAttributeResource() resource.Resource {
	return &identityAttributeResource{}
}

func (r *identityAttributeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *identityAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_attribute"
}

func (r *identityAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.IdentityAttributeSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Identity Attribute.",
		MarkdownDescription: "Manages a SailPoint Identity Attribute. Identity attributes are configurable fields on identity objects. See [Identity Attributes API](https://developer.sailpoint.com/docs/api/v2025/list-identity-attributes/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *identityAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating IdentityAttribute resource")

	var plan models.IdentityAttribute
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiAttribute, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute",
			fmt.Sprintf("Could not convert identity attribute: %s", err.Error()),
		)
		return
	}

	// Create the identity attribute via API
	createdAttribute, err := r.client.CreateIdentityAttribute(ctx, apiAttribute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Identity Attribute",
			fmt.Sprintf("An error occurred while creating the identity attribute: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdAttribute); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute Response",
			fmt.Sprintf("Could not convert identity attribute response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Attribute resource created successfully", map[string]interface{}{
		"attribute_name": plan.Name.ValueString(),
	})
}

func (r *identityAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.IdentityAttribute
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the identity attribute via API
	fetchedAttribute, err := r.client.GetIdentityAttribute(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Identity Attribute",
			fmt.Sprintf("Could not read identity attribute %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedAttribute); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute Response",
			fmt.Sprintf("Could not convert identity attribute response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identityAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Identity Attribute resource")

	var plan models.IdentityAttribute
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiAttribute, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute",
			fmt.Sprintf("Could not convert identity attribute: %s", err.Error()),
		)
		return
	}

	// Update the identity attribute via API (full update with PUT)
	updatedAttribute, err := r.client.UpdateIdentityAttribute(ctx, plan.Name.ValueString(), apiAttribute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Identity Attribute",
			fmt.Sprintf("An error occurred while updating the identity attribute: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedAttribute); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Identity Attribute Response",
			fmt.Sprintf("Could not convert identity attribute response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Identity Attribute resource updated successfully", map[string]interface{}{
		"attribute_name": plan.Name.ValueString(),
	})
}

func (r *identityAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Identity Attribute resource")

	var state models.IdentityAttribute
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the identity attribute via API
	err := r.client.DeleteIdentityAttribute(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Attribute",
			fmt.Sprintf("Could not delete identity attribute %s: %s", state.Name.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Identity Attribute resource deleted successfully", map[string]interface{}{
		"attribute_name": state.Name.ValueString(),
	})
}

func (r *identityAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Identity attributes use "name" as the identifier
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
