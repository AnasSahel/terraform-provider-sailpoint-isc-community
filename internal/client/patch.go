// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

// JSONPatchOperation represents a JSON Patch operation (RFC 6902).
type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// NewReplacePatch creates a JSON Patch "replace" operation for the given path and value.
func NewReplacePatch(path string, value any) JSONPatchOperation {
	return JSONPatchOperation{
		Op:    "replace",
		Path:  path,
		Value: value,
	}
}

// NewRemovePatch creates a JSON Patch "remove" operation for the given path.
func NewRemovePatch(path string) JSONPatchOperation {
	return JSONPatchOperation{
		Op:   "remove",
		Path: path,
	}
}
