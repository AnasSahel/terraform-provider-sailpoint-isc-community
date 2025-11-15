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
