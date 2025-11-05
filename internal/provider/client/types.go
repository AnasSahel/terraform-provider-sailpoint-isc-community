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
