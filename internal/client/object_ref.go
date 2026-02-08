// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

// ObjectRefAPI represents a reference to a SailPoint object (used for Owner and UsedBy).
type ObjectRefAPI struct {
	Type string `json:"type"` // e.g., "IDENTITY", "WORKFLOW"
	ID   string `json:"id"`
	Name string `json:"name,omitempty"` // Optional
}
