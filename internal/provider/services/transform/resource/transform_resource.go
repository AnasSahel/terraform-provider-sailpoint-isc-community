// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &TransformResource{}
	_ resource.ResourceWithConfigure   = &TransformResource{}
	_ resource.ResourceWithImportState = &TransformResource{}
)

// TransformResource is the resource implementation.
type TransformResource struct {
	client *sailpoint.APIClient
}

// NewTransformResource is a helper function to simplify the provider implementation.
func NewTransformResource() resource.Resource {
	return &TransformResource{}
}

// Configure adds the provider configured client to the resource.
func (r *TransformResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sailpoint.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sailpoint.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *TransformResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transform"
}

// Schema defines the schema for the resource.
func (r *TransformResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetTransformResourceSchema()
}

// Create creates the resource and sets the initial Terraform state.
func (r *TransformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TransformResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to SailPoint API object
	transform, diags := plan.ToSailPointTransform()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call SailPoint API to create the transform
	transformResponse, response, err := r.client.V2025.TransformsAPI.CreateTransform(context.Background()).Transform(*transform).Execute()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to create transform '%s'", plan.Name.ValueString())
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if response != nil {
			switch response.StatusCode {
			case 400:
				detailMsg = fmt.Sprintf("Bad Request - The transform configuration is invalid. Please check the 'type' and 'attributes' fields. API error: %s", err.Error())
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to create transforms. Please check your user permissions in SailPoint."
			case 409:
				detailMsg = fmt.Sprintf("Conflict - A transform with name '%s' already exists. Choose a different name.", plan.Name.ValueString())
			case 429:
				detailMsg = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
			default:
				detailMsg = fmt.Sprintf("HTTP %d - %s", response.StatusCode, err.Error())
			}
			detailMsg += fmt.Sprintf("\nHTTP Response: %v", response)
		}

		resp.Diagnostics.AddError(errorMsg, detailMsg)
		return
	}

	// Map response to Terraform model
	diags = plan.FromSailPointTransformRead(ctx, *transformResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *TransformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TransformResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	transform, response, err := r.client.V2025.TransformsAPI.GetTransform(context.Background(), state.Id.ValueString()).Execute()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to read transform with ID '%s'", state.Id.ValueString())
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if response != nil {
			switch response.StatusCode {
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to read transforms. Please check your user permissions in SailPoint."
			case 404:
				detailMsg = fmt.Sprintf("Transform with ID '%s' not found. It may have been deleted outside of Terraform.", state.Id.ValueString())
			case 429:
				detailMsg = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
			default:
				detailMsg = fmt.Sprintf("HTTP %d - %s", response.StatusCode, err.Error())
			}
			detailMsg += fmt.Sprintf("\nHTTP Response: %v", response)
		}

		resp.Diagnostics.AddError(errorMsg, detailMsg)
		return
	}

	// Map response to Terraform model
	diags = state.FromSailPointTransformRead(ctx, *transform)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *TransformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TransformResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to SailPoint API object
	transform, diags := plan.ToSailPointTransform()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	transformResponse, httpResponse, err := r.client.V2025.TransformsAPI.UpdateTransform(context.Background(), plan.Id.ValueString()).Transform(*transform).Execute()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to update transform '%s'", plan.Name.ValueString())
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if httpResponse != nil {
			switch httpResponse.StatusCode {
			case 400:
				detailMsg = fmt.Sprintf("Bad Request - The transform configuration is invalid. Please check the 'attributes' field. Note that 'name' and 'type' cannot be changed after creation. API error: %s", err.Error())
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to update transforms. Please check your user permissions in SailPoint."
			case 404:
				detailMsg = fmt.Sprintf("Transform with ID '%s' not found. It may have been deleted outside of Terraform.", plan.Id.ValueString())
			case 429:
				detailMsg = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
			default:
				detailMsg = fmt.Sprintf("HTTP %d - %s", httpResponse.StatusCode, err.Error())
			}
			detailMsg += fmt.Sprintf("\nHTTP Response: %v", httpResponse)
		}

		resp.Diagnostics.AddError(errorMsg, detailMsg)
		return
	}

	// Map response to Terraform model
	diags = plan.FromSailPointTransformRead(ctx, *transformResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *TransformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TransformResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the transform via SailPoint API
	httpResponse, err := r.client.V2025.TransformsAPI.DeleteTransform(context.Background(), state.Id.ValueString()).Execute()
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to delete transform '%s'", state.Name.ValueString())
		detailMsg := fmt.Sprintf("SailPoint API error: %s", err.Error())

		// Add specific handling for common error scenarios
		if httpResponse != nil {
			switch httpResponse.StatusCode {
			case 401:
				detailMsg = "Unauthorized - Please check your SailPoint credentials and API access."
			case 403:
				detailMsg = "Forbidden - Insufficient permissions to delete transforms. Please check your user permissions in SailPoint."
			case 404:
				// 404 on delete is not necessarily an error - the resource may already be deleted
				detailMsg = fmt.Sprintf("Transform with ID '%s' not found. It may have already been deleted.", state.Id.ValueString())
				resp.Diagnostics.AddWarning(errorMsg, detailMsg)
				return // Don't treat 404 as error for delete operations
			case 409:
				detailMsg = fmt.Sprintf("Conflict - Transform '%s' is still in use and cannot be deleted. Remove references before deleting.", state.Name.ValueString())
			case 429:
				detailMsg = "Rate Limit Exceeded - Too many API requests. Please retry after a few moments."
			default:
				detailMsg = fmt.Sprintf("HTTP %d - %s", httpResponse.StatusCode, err.Error())
			}
			detailMsg += fmt.Sprintf("\nHTTP Response: %v", httpResponse)
		}

		resp.Diagnostics.AddError(errorMsg, detailMsg)
		return
	}
}

// ImportState enables importing existing transforms by ID.
func (r *TransformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
