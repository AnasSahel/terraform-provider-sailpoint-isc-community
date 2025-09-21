// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api_v2025 "github.com/sailpoint-oss/golang-sdk/v2/api_v2025"
)

// ConnectorsDataSourceModel represents the Terraform data source model for listing connectors
type ConnectorsDataSourceModel struct {
	ID           types.String            `tfsdk:"id"`
	Filters      types.String            `tfsdk:"filters"`
	Limit        types.Int32             `tfsdk:"limit"`
	Offset       types.Int32             `tfsdk:"offset"`
	IncludeCount types.Bool              `tfsdk:"include_count"`
	Locale       types.String            `tfsdk:"locale"`
	PaginateAll  types.Bool              `tfsdk:"paginate_all"`
	MaxResults   types.Int32             `tfsdk:"max_results"`
	PageSize     types.Int32             `tfsdk:"page_size"`
	Connectors   []ConnectorSummaryModel `tfsdk:"connectors"`
}

// ConnectorSummaryModel represents a connector in the list response
type ConnectorSummaryModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	ScriptName    types.String `tfsdk:"script_name"`
	ClassName     types.String `tfsdk:"class_name"`
	DirectConnect types.Bool   `tfsdk:"direct_connect"`
	Status        types.String `tfsdk:"status"`
	Category      types.String `tfsdk:"category"`
	Features      types.List   `tfsdk:"features"`
	Labels        types.List   `tfsdk:"labels"`
}

// FromSailPointV3ConnectorDto updates the model from V3ConnectorDto
func (m *ConnectorSummaryModel) FromSailPointV3ConnectorDto(ctx context.Context, connector *api_v2025.V3ConnectorDto) error {
	// Generate ID from script name since V3ConnectorDto doesn't have ID field
	if connector.HasScriptName() {
		m.ID = types.StringValue(connector.GetScriptName())
		m.ScriptName = types.StringValue(connector.GetScriptName())
	} else {
		// Fallback ID generation from name if no script name
		if connector.HasName() {
			hash := fmt.Sprintf("%x", md5.Sum([]byte(connector.GetName())))
			m.ID = types.StringValue(hash[:8])
		}
	}

	if connector.HasName() {
		m.Name = types.StringValue(connector.GetName())
	} else {
		m.Name = types.StringNull()
	}

	if connector.HasType() {
		m.Type = types.StringValue(connector.GetType())
	} else {
		m.Type = types.StringNull()
	}

	if connector.HasClassName() {
		m.ClassName = types.StringValue(connector.GetClassName())
	} else {
		m.ClassName = types.StringNull()
	}

	if connector.HasDirectConnect() {
		m.DirectConnect = types.BoolValue(connector.GetDirectConnect())
	} else {
		m.DirectConnect = types.BoolNull()
	}

	if connector.HasStatus() {
		m.Status = types.StringValue(connector.GetStatus())
	} else {
		m.Status = types.StringNull()
	}

	// Category and Labels are not available in V3ConnectorDto, set as null
	m.Category = types.StringNull()
	m.Labels = types.ListNull(types.StringType)

	// Handle features list if available
	if connector.HasFeatures() {
		features := connector.GetFeatures()
		featureValues := make([]attr.Value, len(features))
		for i, feature := range features {
			featureValues[i] = types.StringValue(feature)
		}
		featuresListValue, err := types.ListValue(types.StringType, featureValues)
		if err != nil {
			return fmt.Errorf("error creating features list: %v", err)
		}
		m.Features = featuresListValue
	} else {
		m.Features = types.ListNull(types.StringType)
	}

	return nil
}

// ConnectorResourceModel represents the Terraform resource model for managing custom connectors
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

// FromSailPointConnectorDetail updates the model from ConnectorDetail
func (m *ConnectorResourceModel) FromSailPointConnectorDetail(ctx context.Context, connector *api_v2025.ConnectorDetail) error {
	// Use script name as ID since it's the unique identifier for custom connectors
	if connector.HasScriptName() {
		m.ID = types.StringValue(connector.GetScriptName())
		m.ScriptName = types.StringValue(connector.GetScriptName())
	}

	if connector.HasName() {
		m.Name = types.StringValue(connector.GetName())
	} else {
		m.Name = types.StringNull()
	}

	if connector.HasType() {
		m.Type = types.StringValue(connector.GetType())
	} else {
		m.Type = types.StringNull()
	}

	if connector.HasClassName() {
		m.ClassName = types.StringValue(connector.GetClassName())
	} else {
		m.ClassName = types.StringNull()
	}

	if connector.HasApplicationXml() {
		m.ApplicationXml = types.StringValue(connector.GetApplicationXml())
	} else {
		m.ApplicationXml = types.StringNull()
	}

	if connector.HasCorrelationConfigXml() {
		m.CorrelationConfigXml = types.StringValue(connector.GetCorrelationConfigXml())
	} else {
		m.CorrelationConfigXml = types.StringNull()
	}

	if connector.HasSourceConfigXml() {
		m.SourceConfigXml = types.StringValue(connector.GetSourceConfigXml())
	} else {
		m.SourceConfigXml = types.StringNull()
	}

	if connector.HasSourceConfig() {
		m.SourceConfig = types.StringValue(connector.GetSourceConfig())
	} else {
		m.SourceConfig = types.StringNull()
	}

	if connector.HasS3Location() {
		m.S3Location = types.StringValue(connector.GetS3Location())
	} else {
		m.S3Location = types.StringNull()
	}

	if connector.HasFileUpload() {
		m.FileUpload = types.BoolValue(connector.GetFileUpload())
	} else {
		m.FileUpload = types.BoolNull()
	}

	if connector.HasDirectConnect() {
		m.DirectConnect = types.BoolValue(connector.GetDirectConnect())
	} else {
		m.DirectConnect = types.BoolNull()
	}

	if connector.HasStatus() {
		m.Status = types.StringValue(connector.GetStatus())
	} else {
		m.Status = types.StringNull()
	}

	// Handle uploaded files list
	if connector.HasUploadedFiles() {
		files := connector.GetUploadedFiles()
		fileValues := make([]attr.Value, len(files))
		for i, file := range files {
			fileValues[i] = types.StringValue(file)
		}
		filesListValue, err := types.ListValue(types.StringType, fileValues)
		if err != nil {
			return fmt.Errorf("error creating uploaded_files list: %v", err)
		}
		m.UploadedFiles = filesListValue
	} else {
		m.UploadedFiles = types.ListNull(types.StringType)
	}

	// Handle translation properties JSON
	if connector.HasTranslationProperties() {
		translationPropsBytes, err := json.Marshal(connector.GetTranslationProperties())
		if err != nil {
			return fmt.Errorf("error marshaling translation properties: %v", err)
		}
		m.TranslationProperties = types.StringValue(string(translationPropsBytes))
	} else {
		m.TranslationProperties = types.StringNull()
	}

	// Handle connector metadata JSON
	if connector.HasConnectorMetadata() {
		metadataBytes, err := json.Marshal(connector.GetConnectorMetadata())
		if err != nil {
			return fmt.Errorf("error marshaling connector metadata: %v", err)
		}
		m.ConnectorMetadata = types.StringValue(string(metadataBytes))
	} else {
		m.ConnectorMetadata = types.StringNull()
	}

	return nil
}

// ToSailPointV3CreateConnectorDto converts the model to V3CreateConnectorDto for create operations
func (m *ConnectorResourceModel) ToSailPointV3CreateConnectorDto() (*api_v2025.V3CreateConnectorDto, error) {
	// Name and ClassName are required fields for V3CreateConnectorDto
	name := m.Name.ValueString()
	className := "sailpoint.connector.OpenConnectorAdapter" // Default OpenConnector class
	if !m.ClassName.IsNull() && !m.ClassName.IsUnknown() {
		className = m.ClassName.ValueString()
	}

	connector := api_v2025.NewV3CreateConnectorDto(name, className)

	if !m.Type.IsNull() && !m.Type.IsUnknown() {
		connector.SetType(m.Type.ValueString())
	}

	if !m.DirectConnect.IsNull() && !m.DirectConnect.IsUnknown() {
		connector.SetDirectConnect(m.DirectConnect.ValueBool())
	}

	if !m.Status.IsNull() && !m.Status.IsUnknown() {
		connector.SetStatus(m.Status.ValueString())
	}

	return connector, nil
}

// ToSailPointV3ConnectorDto converts the model to V3ConnectorDto for create operations
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
			return nil, fmt.Errorf("error unmarshaling connector metadata: %v", err)
		}
		connector.SetConnectorMetadata(metadata)
	}

	return connector, nil
}
