// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

// EntitlementDirectPermission represents a permission with rights and target.
type EntitlementDirectPermission struct {
	Rights []string `json:"rights,omitempty"`
	Target *string  `json:"target,omitempty"`
}

// EntitlementManuallyUpdatedFields tracks which fields have been manually updated.
type EntitlementManuallyUpdatedFields struct {
	DisplayName *bool `json:"DISPLAY_NAME,omitempty"`
	Description *bool `json:"DESCRIPTION,omitempty"`
}

// EntitlementAccessModelMetadataValue represents a value in access model metadata.
type EntitlementAccessModelMetadataValue struct {
	Value  *string `json:"value,omitempty"`
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

// EntitlementAccessModelMetadataAttribute represents an attribute in access model metadata.
type EntitlementAccessModelMetadataAttribute struct {
	Key         *string                               `json:"key,omitempty"`
	Name        *string                               `json:"name,omitempty"`
	Multiselect *bool                                 `json:"multiselect,omitempty"`
	Status      *string                               `json:"status,omitempty"`
	Type        *string                               `json:"type,omitempty"`
	ObjectTypes []string                              `json:"objectTypes,omitempty"`
	Description *string                               `json:"description,omitempty"`
	Values      []EntitlementAccessModelMetadataValue `json:"values,omitempty"`
}

// EntitlementAccessModelMetadata represents access model metadata for an entitlement.
type EntitlementAccessModelMetadata struct {
	Attributes []EntitlementAccessModelMetadataAttribute `json:"attributes,omitempty"`
}

// Entitlement represents a SailPoint Entitlement.
type Entitlement struct {
	ID                     string                            `json:"id,omitempty"`
	Name                   string                            `json:"name"`
	Created                *string                           `json:"created,omitempty"`
	Modified               *string                           `json:"modified,omitempty"`
	Description            *string                           `json:"description,omitempty"`
	Attribute              *string                           `json:"attribute,omitempty"`
	Value                  *string                           `json:"value,omitempty"`
	SourceSchemaObjectType *string                           `json:"sourceSchemaObjectType,omitempty"`
	Privileged             *bool                             `json:"privileged,omitempty"`
	Requestable            *bool                             `json:"requestable,omitempty"`
	CloudGoverned          *bool                             `json:"cloudGoverned,omitempty"`
	Source                 *ObjectRef                        `json:"source,omitempty"`
	Owner                  *ObjectRef                        `json:"owner,omitempty"`
	Attributes             map[string]interface{}            `json:"attributes,omitempty"`
	DirectPermissions      []EntitlementDirectPermission     `json:"directPermissions,omitempty"`
	Segments               []string                          `json:"segments,omitempty"`
	ManuallyUpdatedFields  *EntitlementManuallyUpdatedFields `json:"manuallyUpdatedFields,omitempty"`
	AccessModelMetadata    *EntitlementAccessModelMetadata   `json:"accessModelMetadata,omitempty"`
}

// GetEntitlement retrieves an entitlement by ID from SailPoint.
func (c *Client) GetEntitlement(ctx context.Context, id string) (*Entitlement, error) {
	var result Entitlement
	path := fmt.Sprintf("/v2025/entitlements/%s", id)

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "entitlement",
			ResourceID: id,
		}, err)
	}

	if resp.IsError() {
		return nil, c.formatErrorWithBody(ErrorContext{
			Operation:  "read",
			Resource:   "entitlement",
			ResourceID: id,
		}, resp.StatusCode(), resp.String())
	}

	return &result, nil
}
