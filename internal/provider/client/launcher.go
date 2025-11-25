// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
)

// LauncherReference represents the reference to a workflow or other resource.
type LauncherReference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Launcher represents a SailPoint Launcher.
type Launcher struct {
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	Disabled    bool               `json:"disabled"`
	Reference   *LauncherReference `json:"reference"`
	Config      string             `json:"config"`
	Owner       *ObjectRef         `json:"owner,omitempty"`
	Created     *string            `json:"created,omitempty"`
	Modified    *string            `json:"modified,omitempty"`
}

// CreateLauncher creates a new launcher in SailPoint.
func (c *Client) CreateLauncher(ctx context.Context, launcher *Launcher) (*Launcher, error) {
	var result Launcher
	resp, err := c.doRequest(ctx, http.MethodPost, "/v2025/launchers", launcher, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "launcher",
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "launcher",
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// GetLauncher retrieves a launcher by ID from SailPoint.
func (c *Client) GetLauncher(ctx context.Context, id string) (*Launcher, error) {
	var result Launcher
	path := fmt.Sprintf("/v2025/launchers/%s", id)

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "launcher",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "read",
			Resource:   "launcher",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// UpdateLauncher performs a full update (PUT) of a launcher.
func (c *Client) UpdateLauncher(ctx context.Context, id string, launcher *Launcher) (*Launcher, error) {
	var result Launcher
	path := fmt.Sprintf("/v2025/launchers/%s", id)

	resp, err := c.doRequest(ctx, http.MethodPut, path, launcher, &result)
	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "launcher",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "launcher",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return &result, nil
}

// DeleteLauncher deletes a launcher by ID.
func (c *Client) DeleteLauncher(ctx context.Context, id string) error {
	path := fmt.Sprintf("/v2025/launchers/%s", id)

	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "launcher",
			ResourceID: id,
		}, err, 0)
	}

	if resp.IsError() {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "launcher",
			ResourceID: id,
		}, nil, resp.StatusCode())
	}

	return nil
}
