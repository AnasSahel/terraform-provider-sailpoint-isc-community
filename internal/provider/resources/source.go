package resources

import (
	"context"
	"fmt"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/client"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/models"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/schemas"
	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sourceResource{}
	_ resource.ResourceWithConfigure   = &sourceResource{}
	_ resource.ResourceWithImportState = &sourceResource{}
)

type sourceResource struct {
	client *client.Client
}

func NewSourceResource() resource.Resource {
	return new(sourceResource)
}

// Then in both source_resource.go and source_data_source.go:
func (r *sourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := utils.ConfigureClient(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.client = client
}

func (r *sourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *sourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaBuilder := schemas.SourceSchemaBuilder{}
	resp.Schema = schema.Schema{
		Description:         "Resource representing a SailPoint Source.",
		MarkdownDescription: "Resource representing a SailPoint Source.",
		Attributes:          schemaBuilder.GetResourceSchema(),
	}
}

func (r *sourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Source resource")
	var plan models.Source
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Source create plan:", map[string]interface{}{
		"plan": plan,
	})

	source := plan.ConvertToCreateRequestPtr(ctx)
	tflog.Debug(ctx, "Converted Source create request:", map[string]interface{}{
		"source": source,
	})

	createdSource, err := r.client.CreateSource(ctx, source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Source",
			fmt.Sprintf("An error occurred while creating the source: %s", err.Error()),
		)
		return
	}

	plan.ConvertFromSailPointForResource(ctx, createdSource)
	tflog.Debug(ctx, fmt.Sprintf("Converted Source state: %+v", plan))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Source resource created successfully", map[string]interface{}{
		"source_id": plan.ID.ValueString(),
	})
}

func (r *sourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fetchedSource, err := r.client.GetSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			fmt.Sprintf("An error occurred while reading the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.ConvertFromSailPointForResource(ctx, fetchedSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// jsonPatches := []client.JSONPatchOperation{}
	var plan, state models.Source
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	patches := state.BuildPatchOptions(ctx, &plan)

	if len(patches) == 0 {
		return
	}

	updatedSource, err := r.client.PatchSource(ctx, state.ID.ValueString(), patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			fmt.Sprintf("An error occurred while updating the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	plan.ConvertFromSailPointForResource(ctx, updatedSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.Source
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSource(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			fmt.Sprintf("An error occurred while deleting the source with ID %q: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func (r *sourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
