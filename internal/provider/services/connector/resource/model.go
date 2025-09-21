// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector_resource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"

	"github.com/AnasSahel/terraform-provider-sailpoint-isc-community/internal/provider/utils"
)

// ConnectorResourceModel represents the Terraform resource model for managing custom connectors.
type ConnectorResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	ClassName             types.String `tfsdk:"class_name"`
	ScriptName            types.String `tfsdk:"script_name"`
	ApplicationXml        types.String `tfsdk:"application_xml"`
	CorrelationConfigXml  types.String `tfsdk:"correlation_config_xml"`
	SourceConfigXml       types.String `tfsdk:"source_config_xml"`
	SourceConfig          types.String `tfsdk:"source_config"`
	S3Location            types.String `tfsdk:"s3_location"`
	UploadedFiles         types.List   `tfsdk:"uploaded_files"`
	FileUpload            types.Bool   `tfsdk:"file_upload"`
	DirectConnect         types.Bool   `tfsdk:"direct_connect"`
	TranslationProperties types.String `tfsdk:"translation_properties"` // JSON string
	ConnectorMetadata     types.String `tfsdk:"connector_metadata"`     // JSON string
	Status                types.String `tfsdk:"status"`
}

// FromSailPointConnectorDetail updates the model from ConnectorDetail.
func (m *ConnectorResourceModel) FromSailPointConnectorDetail(ctx context.Context, connector *api_v2025.ConnectorDetail) error {
	// Use script name as ID since it's the unique identifier for custom connectors
	if connector.HasScriptName() {
		m.ID = types.StringValue(connector.GetScriptName())
		m.ScriptName = types.StringValue(connector.GetScriptName())
	}

	m.Name = utils.StringOrNull(connector.HasName(), connector.GetName())
	m.Type = utils.StringOrNull(connector.HasType(), connector.GetType())
	m.ClassName = utils.StringOrNull(connector.HasClassName(), connector.GetClassName())

	m.ApplicationXml = utils.StringOrNull(connector.HasApplicationXml(), connector.GetApplicationXml())
	m.CorrelationConfigXml = utils.StringOrNull(connector.HasCorrelationConfigXml(), connector.GetCorrelationConfigXml())
	m.SourceConfigXml = utils.StringOrNull(connector.HasSourceConfigXml(), connector.GetSourceConfigXml())
	m.SourceConfig = utils.StringOrNull(connector.HasSourceConfig(), connector.GetSourceConfig())
	m.S3Location = utils.StringOrNull(connector.HasS3Location(), connector.GetS3Location())

	m.FileUpload = utils.BoolOrNull(connector.HasFileUpload(), connector.GetFileUpload())

	m.DirectConnect = utils.BoolOrNull(connector.HasDirectConnect(), connector.GetDirectConnect())
	m.Status = utils.StringOrNull(connector.HasStatus(), connector.GetStatus())

	// Handle uploaded files list
	if connector.HasUploadedFiles() {
		files := connector.GetUploadedFiles()
		fileValues := make([]attr.Value, len(files))
		for i, file := range files {
			fileValues[i] = types.StringValue(file)
		}
		filesListValue, diags := types.ListValue(types.StringType, fileValues)
		if diags.HasError() {
			return fmt.Errorf("error creating uploaded_files list: %s", diags.Errors())
		}
		m.UploadedFiles = filesListValue
	} else {
		m.UploadedFiles = types.ListNull(types.StringType)
	}

	// Handle translation properties JSON
	if connector.HasTranslationProperties() {
		translationPropsBytes, err := json.Marshal(connector.GetTranslationProperties())
		if err != nil {
			return fmt.Errorf("error marshaling translation properties: %w", err)
		}
		m.TranslationProperties = types.StringValue(string(translationPropsBytes))
	} else {
		m.TranslationProperties = types.StringNull()
	}

	// Handle connector metadata JSON
	if connector.HasConnectorMetadata() {
		metadataBytes, err := json.Marshal(connector.GetConnectorMetadata())
		if err != nil {
			return fmt.Errorf("error marshaling connector metadata: %w", err)
		}
		m.ConnectorMetadata = types.StringValue(string(metadataBytes))
	} else {
		m.ConnectorMetadata = types.StringNull()
	}

	return nil
}

// ToSailPointV3CreateConnectorDto converts the model to V3CreateConnectorDto for create operations.
func (m *ConnectorResourceModel) ToSailPointV3CreateConnectorDto() (*api_v2025.V3CreateConnectorDto, error) {
	// Name and ClassName are required fields for V3CreateConnectorDto
	name := m.Name.ValueString()
	className := "sailpoint.connector.OpenConnectorAdapter" // Default OpenConnector class
	if !m.ClassName.IsNull() && !m.ClassName.IsUnknown() {
		className = m.ClassName.ValueString()
	}

	connector := api_v2025.NewV3CreateConnectorDto(name, className)

	// Note: Script name is auto-generated from the name by SailPoint API
	// DirectConnect defaults to true in the DTO constructor

	// Only set type if explicitly provided, otherwise let API default it to 'custom ' + name
	if !m.Type.IsNull() && !m.Type.IsUnknown() && m.Type.ValueString() != "" {
		connector.SetType(m.Type.ValueString())
	}

	// Only override DirectConnect if explicitly set to false
	if !m.DirectConnect.IsNull() && !m.DirectConnect.IsUnknown() && !m.DirectConnect.ValueBool() {
		connector.SetDirectConnect(false)
	}

	// Set status if provided (defaults to DEVELOPMENT for custom connectors)
	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		connector.SetStatus(m.Status.ValueString())
	}

	return connector, nil
}

// ToSailPointV3ConnectorDto converts the model to V3ConnectorDto for create operations.
func (m *ConnectorResourceModel) ToSailPointV3ConnectorDto() (*api_v2025.V3ConnectorDto, error) {
	connector := api_v2025.NewV3ConnectorDto()

	if !m.Name.IsNull() && !m.Name.IsUnknown() {
		connector.SetName(m.Name.ValueString())
	}

	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		connector.SetType(m.Type.ValueString())
	}

	if !m.ScriptName.IsNull() && !m.ScriptName.IsUnknown() {
		connector.SetScriptName(m.ScriptName.ValueString())
	}

	if !m.ClassName.IsNull() && !m.ClassName.IsUnknown() {
		connector.SetClassName(m.ClassName.ValueString())
	}

	if !m.DirectConnect.IsNull() && !m.DirectConnect.IsUnknown() {
		connector.SetDirectConnect(m.DirectConnect.ValueBool())
	}

	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		connector.SetStatus(m.Status.ValueString())
	}

	// Handle connector metadata
	if !m.ConnectorMetadata.IsNull() && !m.ConnectorMetadata.IsUnknown() {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(m.ConnectorMetadata.ValueString()), &metadata); err != nil {
			return nil, fmt.Errorf("error unmarshaling connector metadata: %w", err)
		}
		connector.SetConnectorMetadata(metadata)
	}

	return connector, nil
}
