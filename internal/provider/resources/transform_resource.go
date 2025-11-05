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
	_ resource.Resource                = &transformResource{}
	_ resource.ResourceWithConfigure   = &transformResource{}
	_ resource.ResourceWithImportState = &transformResource{}
)

type transformResource struct {
	client *client.Client
}

func NewTransformResource() resource.Resource {
	return &transformResource{}
}

func (r *transformResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *transformResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

func (r *transformResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.TransformSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Transform.",
		MarkdownDescription: "Manages a SailPoint Transform. Transforms are configurable objects that manipulate attribute data during aggregation or provisioning. See [Transform Documentation](https://developer.sailpoint.com/docs/extensibility/transforms/) for more information.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *transformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Transform resource")

	var plan models.Transform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiTransform, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform",
			fmt.Sprintf("Could not convert transform attributes: %s", err.Error()),
		)
		return
	}

	// Create the transform via API
	createdTransform, err := r.client.CreateTransform(ctx, apiTransform)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Transform",
			fmt.Sprintf("An error occurred while creating the transform: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, createdTransform); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform Response",
			fmt.Sprintf("Could not convert transform response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Transform resource created successfully", map[string]interface{}{
		"transform_id": plan.ID.ValueString(),
	})
}

func (r *transformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Transform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the transform via API
	fetchedTransform, err := r.client.GetTransform(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Transform",
			fmt.Sprintf("Could not read transform ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to Terraform model
	if err := state.ConvertFromSailPointForResource(ctx, fetchedTransform); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform Response",
			fmt.Sprintf("Could not convert transform response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *transformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Transform resource")

	var plan models.Transform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	apiTransform, err := plan.ConvertToSailPoint(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform",
			fmt.Sprintf("Could not convert transform attributes: %s", err.Error()),
		)
		return
	}

	// Update the transform via API
	// Note: Only 'attributes' can be updated; 'name' and 'type' are immutable
	updatedTransform, err := r.client.UpdateTransform(ctx, plan.ID.ValueString(), apiTransform)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Transform",
			fmt.Sprintf("An error occurred while updating the transform: %s", err.Error()),
		)
		return
	}

	// Convert API response back to Terraform model
	if err := plan.ConvertFromSailPointForResource(ctx, updatedTransform); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Transform Response",
			fmt.Sprintf("Could not convert transform response: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Transform resource updated successfully", map[string]interface{}{
		"transform_id": plan.ID.ValueString(),
	})
}

func (r *transformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Transform resource")

	var state models.Transform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the transform via API
	err := r.client.DeleteTransform(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Transform",
			fmt.Sprintf("Could not delete transform ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Transform resource deleted successfully", map[string]interface{}{
		"transform_id": state.ID.ValueString(),
	})
}

func (r *transformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
