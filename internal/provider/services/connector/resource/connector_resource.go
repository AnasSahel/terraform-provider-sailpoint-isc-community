// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/iancoleman/strcase"
	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ConnectorResource{}
	_ resource.ResourceWithConfigure   = &ConnectorResource{}
	_ resource.ResourceWithImportState = &ConnectorResource{}
)

func NewConnectorResource() resource.Resource {
	return &ConnectorResource{}
}

// ConnectorResource defines the resource implementation.
type ConnectorResource struct {
	client *api_v2025.APIClient
}

func (r *ConnectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector"
}

func (r *ConnectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetConnectorResourceSchema()
}

func (r *ConnectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api_v2025.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api_v2025.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ConnectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConnectorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate script name if not provided
	if data.ScriptName.IsNull() || data.ScriptName.IsUnknown() {
		// Convert name to snake_case for script name
		scriptName := strcase.ToSnake(data.Name.ValueString())
		scriptName = strings.ReplaceAll(scriptName, " ", "_")
		scriptName = strings.ToLower(scriptName)
		data.ScriptName = types.StringValue(scriptName)
	}

	// Convert to SailPoint API model for create
	createDto, err := data.ToSailPointV3CreateConnectorDto()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Connector Data",
			fmt.Sprintf("Could not convert connector data: %v", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating custom connector with script name: %s", data.ScriptName.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("CreateDTO: Name=%s, ClassName=%s", createDto.GetName(), createDto.GetClassName()))
	if createDto.HasType() {
		tflog.Debug(ctx, fmt.Sprintf("CreateDTO Type: %s", createDto.GetType()))
	}
	if createDto.HasDirectConnect() {
		tflog.Debug(ctx, fmt.Sprintf("CreateDTO DirectConnect: %v", createDto.GetDirectConnect()))
	}

	// Log the full DTO for debugging
	tflog.Debug(ctx, fmt.Sprintf("Full CreateDTO: %+v", createDto))

	// Create custom connector
	apiReq := r.client.ConnectorsAPI.CreateCustomConnector(ctx).V3CreateConnectorDto(*createDto)
	createdConnector, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Connector",
			fmt.Sprintf("Could not create connector %s: %v\nHTTP Response: %+v", data.Name.ValueString(), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully created custom connector: %s", createdConnector.GetScriptName()))

	// Update model with created connector data - create a minimal ConnectorDetail from V3ConnectorDto
	connectorDetail := &api_v2025.ConnectorDetail{
		Name:       createdConnector.Name,
		Type:       createdConnector.Type,
		ScriptName: createdConnector.ScriptName,
		Status:     createdConnector.Status,
	}

	if createdConnector.HasClassName() {
		className := createdConnector.GetClassName()
		connectorDetail.ClassName = &className
	}

	if createdConnector.DirectConnect != nil {
		connectorDetail.DirectConnect = createdConnector.DirectConnect
	}

	if createdConnector.ConnectorMetadata != nil {
		connectorDetail.ConnectorMetadata = createdConnector.ConnectorMetadata
	}

	if err := data.FromSailPointConnectorDetail(ctx, connectorDetail); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Created Connector Data",
			fmt.Sprintf("Could not convert created connector data: %v", err),
		)
		return
	}

	// Set ID to script name
	data.ID = data.ScriptName

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConnectorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptName := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading custom connector: %s", scriptName))

	// Get connector details
	apiReq := r.client.ConnectorsAPI.GetConnector(ctx, scriptName)
	connector, httpResp, err := apiReq.Execute()
	if err != nil {
		// Handle 404 case - connector was deleted outside of Terraform
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Warn(ctx, fmt.Sprintf("Connector %s not found, removing from state", scriptName))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Connector",
			fmt.Sprintf("Could not read connector %s: %v\nHTTP Response: %+v", scriptName, err, httpResp),
		)
		return
	}

	// Update model with current connector data
	if err := data.FromSailPointConnectorDetail(ctx, connector); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Connector Data",
			fmt.Sprintf("Could not convert connector data: %v", err),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConnectorResourceModel
	var currentState ConnectorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptName := currentState.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating custom connector: %s", scriptName))

	// Create JSON patch operations for the changes
	var patches []api_v2025.JsonPatchOperation

	// Only patch fields that are supported by the SailPoint API according to documentation:
	// https://developer.sailpoint.com/docs/api/v2025/update-connector
	// Patchable fields: connectorMetadata, applicationXml, correlationConfigXml, sourceConfigXml

	if !data.ConnectorMetadata.Equal(currentState.ConnectorMetadata) && !data.ConnectorMetadata.IsNull() {
		metadataValue := data.ConnectorMetadata.ValueString()
		patchValue := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&metadataValue)
		patches = append(patches, api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/connectorMetadata",
			Value: &patchValue,
		})
	}

	if !data.ApplicationXml.Equal(currentState.ApplicationXml) && !data.ApplicationXml.IsNull() {
		xmlValue := data.ApplicationXml.ValueString()
		patchValue := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&xmlValue)
		patches = append(patches, api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/applicationXml",
			Value: &patchValue,
		})
	}

	if !data.CorrelationConfigXml.Equal(currentState.CorrelationConfigXml) && !data.CorrelationConfigXml.IsNull() {
		xmlValue := data.CorrelationConfigXml.ValueString()
		patchValue := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&xmlValue)
		patches = append(patches, api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/correlationConfigXml",
			Value: &patchValue,
		})
	}

	if !data.SourceConfigXml.Equal(currentState.SourceConfigXml) && !data.SourceConfigXml.IsNull() {
		xmlValue := data.SourceConfigXml.ValueString()
		patchValue := api_v2025.StringAsUpdateMultiHostSourcesRequestInnerValue(&xmlValue)
		patches = append(patches, api_v2025.JsonPatchOperation{
			Op:    "replace",
			Path:  "/sourceConfigXml",
			Value: &patchValue,
		})
	}

	// Note: name, type, className, directConnect, and status are create-only fields and require resource replacement

	// Only proceed with update if there are actual changes
	if len(patches) == 0 {
		tflog.Debug(ctx, "No changes detected for connector update")
		// Just refresh the state by reading current data
		// Get current connector details
		getReq := r.client.ConnectorsAPI.GetConnector(ctx, scriptName)
		connector, httpResp, err := getReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Connector During Update",
				fmt.Sprintf("Could not read connector %s: %v\nHTTP Response: %+v", scriptName, err, httpResp),
			)
			return
		}

		// Update model with current connector data
		if err := data.FromSailPointConnectorDetail(ctx, connector); err != nil {
			resp.Diagnostics.AddError(
				"Error Converting Connector Data During Update",
				fmt.Sprintf("Could not convert connector data: %v", err),
			)
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	// Apply patches to update connector
	apiReq := r.client.ConnectorsAPI.UpdateConnector(ctx, scriptName).JsonPatchOperation(patches)
	updatedConnector, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Connector",
			fmt.Sprintf("Could not update connector %s: %v\nHTTP Response: %+v", scriptName, err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated custom connector: %s", scriptName))

	// Update model with updated connector data
	if err := data.FromSailPointConnectorDetail(ctx, updatedConnector); err != nil {
		resp.Diagnostics.AddError(
			"Error Converting Updated Connector Data",
			fmt.Sprintf("Could not convert updated connector data: %v", err),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConnectorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptName := data.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting custom connector: %s", scriptName))

	// Delete connector
	apiReq := r.client.ConnectorsAPI.DeleteCustomConnector(ctx, scriptName)
	httpResp, err := apiReq.Execute()
	if err != nil {
		// Handle 404 case - connector was already deleted
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Warn(ctx, fmt.Sprintf("Connector %s already deleted", scriptName))
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting Connector",
			fmt.Sprintf("Could not delete connector %s: %v\nHTTP Response: %+v", scriptName, err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted custom connector: %s", scriptName))
}

func (r *ConnectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by script name
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
