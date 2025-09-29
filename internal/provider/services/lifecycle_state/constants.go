// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lifecycle_state

// Resource and data source type names.
const (
	ResourceTypeName       = "_lifecycle_state"
	DataSourceTypeName     = "_lifecycle_state"
	ListDataSourceTypeName = "_lifecycle_state_list"
)

// Error message constants for consistent error handling.
const (
	ErrUnexpectedConfigureType  = "Unexpected Configure Type"
	ErrMissingIdentityProfileID = "Missing Identity Profile ID"
	ErrInvalidIdentityProfileID = "Invalid Identity Profile ID"
	ErrInvalidLifecycleStateID  = "Invalid Lifecycle State ID"
	ErrInvalidImportID          = "Invalid Import ID"
	ErrCreatingLifecycleState   = "Error Creating Lifecycle State"
	ErrReadingLifecycleState    = "Error Reading Lifecycle State"
	ErrReadingLifecycleStates   = "Error Reading Lifecycle States"
	ErrUpdatingLifecycleState   = "Error Updating Lifecycle State"
	ErrDeletingLifecycleState   = "Error Deleting Lifecycle State"
)

// Common messages.
const (
	MsgExpectedAPIClient       = "Expected *sailpoint.APIClient. Please report this issue to the provider developers."
	MsgIdentityProfileRequired = "The identity profile ID must be specified."
	MsgLifecycleStateRequired  = "The lifecycle state ID must be specified."
	MsgImportFormatRequired    = "Import ID must be in format 'identity_profile_id:lifecycle_state_id'"
)

// Import format constants.
const (
	ImportIDSeparator = ":"
	ImportIDParts     = 2
)

// JSON Patch paths.
const (
	PatchPathEnabled     = "/enabled"
	PatchPathDescription = "/description"
	PatchPathPriority    = "/priority"
)

// Valid identity states for validation.
var ValidIdentityStates = []string{
	"ACTIVE",
	"INACTIVE_SHORT_TERM",
	"INACTIVE_LONG_TERM",
}

// Valid account actions for validation.
var ValidAccountActions = []string{
	"ENABLE",
	"DISABLE",
	"DELETE",
}
