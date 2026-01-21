// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

// ObjectRef represents a reference to a SailPoint object.
// This is a common structure used across many SailPoint API resources
// to reference related objects like owners, clusters, rules, etc.
type ObjectRef struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// WorkflowTrigger represents a trigger configuration for a workflow.
type WorkflowTrigger struct {
	Type        string                 `json:"type"`
	DisplayName string                 `json:"displayName,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// WorkflowDefinition represents the definition of a workflow's logic.
type WorkflowDefinition struct {
	Start string                 `json:"start"`
	Steps map[string]interface{} `json:"steps"`
}

// ApprovalScheme represents an approval step configuration.
type ApprovalScheme struct {
	ApproverType string  `json:"approverType"` // APP_OWNER, OWNER, SOURCE_OWNER, MANAGER, GOVERNANCE_GROUP, WORKFLOW
	ApproverID   *string `json:"approverId,omitempty"`
}

// AccessRequestConfig (Requestability) represents access request configuration.
type AccessRequestConfig struct {
	CommentsRequired        *bool            `json:"commentsRequired,omitempty"`
	DenialCommentsRequired  *bool            `json:"denialCommentsRequired,omitempty"`
	ReauthorizationRequired *bool            `json:"reauthorizationRequired,omitempty"`
	ApprovalSchemes         []ApprovalScheme `json:"approvalSchemes,omitempty"`
}

// RevocationRequestConfig (Revocability) represents revocation request configuration.
type RevocationRequestConfig struct {
	ApprovalSchemes []ApprovalScheme `json:"approvalSchemes,omitempty"`
}

// ProvisioningCriteria represents provisioning criteria for multi-account selection.
// Supports up to 3 levels of nesting.
type ProvisioningCriteria struct {
	Operation string                  `json:"operation"` // EQUALS, NOT_EQUALS, CONTAINS, HAS, AND, OR
	Attribute *string                 `json:"attribute,omitempty"`
	Value     *string                 `json:"value,omitempty"`
	Children  *[]ProvisioningCriteria `json:"children,omitempty"`
}

// EmailNotificationOption represents email notification settings for a lifecycle state.
type EmailNotificationOption struct {
	NotifyManagers      *bool    `json:"notifyManagers,omitempty"`
	NotifyAllAdmins     *bool    `json:"notifyAllAdmins,omitempty"`
	NotifySpecificUsers *bool    `json:"notifySpecificUsers,omitempty"`
	EmailAddressList    []string `json:"emailAddressList,omitempty"`
}

// AccountAction represents an account action configuration for a lifecycle state.
type AccountAction struct {
	Action           string   `json:"action"` // ENABLE, DISABLE, DELETE
	SourceIds        []string `json:"sourceIds,omitempty"`
	ExcludeSourceIds []string `json:"excludeSourceIds,omitempty"`
	AllSources       *bool    `json:"allSources,omitempty"`
}

// AccessActionConfiguration represents access action settings for a lifecycle state.
type AccessActionConfiguration struct {
	RemoveAllAccessEnabled *bool `json:"removeAllAccessEnabled,omitempty"`
}
