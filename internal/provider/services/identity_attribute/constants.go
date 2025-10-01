package identity_attribute

// Resource and data source type names.
const (
	ResourceTypeName       = "_identity_attribute"
	DataSourceTypeName     = "_identity_attribute"
	DataSourceListTypeName = "_identity_attribute_list"
)

// Datasource Errors
const (
	ErrProviderDataTitle = "Unexpected Data Source Configure Type"
	ErrProviderDataMsg   = "Expected *sailpoint.APIClient, got: %T. Please report this to the provider developers."

	ErrDataSourceReadTitle = "Error Reading Identity Attribute"
	ErrDataSourceReadMsg   = "Could not read identity attribute %q: %v\n\nFull HTTP response: %v"

	ErrResourceCreateTitle = "Error Creating Identity Attribute"
	ErrResourceCreateMsg   = "Could not create identity attribute %q: %v\n\nFull HTTP response: %v"
	ErrResourceReadTitle   = "Error Reading Identity Attribute"
	ErrResourceReadMsg     = "Could not read identity attribute %q: %v\n\nFull HTTP response: %v"
	ErrResourceUpdateTitle = "Error Updating Identity Attribute"
	ErrResourceUpdateMsg   = "Could not update identity attribute %q: %v\n\nFull HTTP response: %v"
)
