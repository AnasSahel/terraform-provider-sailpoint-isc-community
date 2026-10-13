// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	managedClusterEndpointList = "/managed-clusters/v1"
	managedClusterEndpointGet  = "/managed-clusters/v1/{id}"

	managedClustersPageSize = 250
)

// ManagedClusterAPI represents a SailPoint managed cluster (virtual appliance cluster).
// Key material returned by the API (keyPair, attributes.keystore, public keys) is
// intentionally not decoded so it never reaches Terraform state.
type ManagedClusterAPI struct {
	ID                                 string             `json:"id"`
	Name                               string             `json:"name,omitempty"`
	Pod                                string             `json:"pod,omitempty"`
	Org                                string             `json:"org,omitempty"`
	Type                               string             `json:"type,omitempty"`
	Configuration                      map[string]*string `json:"configuration,omitempty"`
	Description                        *string            `json:"description,omitempty"`
	ClientType                         *string            `json:"clientType,omitempty"`
	CcgVersion                         string             `json:"ccgVersion,omitempty"`
	PinnedConfig                       *bool              `json:"pinnedConfig,omitempty"`
	Operational                        *bool              `json:"operational,omitempty"`
	Status                             *string            `json:"status,omitempty"`
	AlertKey                           *string            `json:"alertKey,omitempty"`
	ClientIDs                          []string           `json:"clientIds,omitempty"`
	ServiceCount                       *int32             `json:"serviceCount,omitempty"`
	CurrentInstalledReleaseVersion     *string            `json:"currentInstalledReleaseVersion,omitempty"`
	ConsolidatedHealthIndicatorsStatus *string            `json:"consolidatedHealthIndicatorsStatus,omitempty"`
	CreatedAt                          *string            `json:"createdAt,omitempty"`
	UpdatedAt                          *string            `json:"updatedAt,omitempty"`
}

type managedClusterErrorContext struct {
	Operation    string
	ID           string
	ResponseBody string
}

// GetManagedCluster retrieves a managed cluster by ID.
func (c *Client) GetManagedCluster(ctx context.Context, id string) (*ManagedClusterAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("managed cluster ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting managed cluster", map[string]any{"id": id})

	var cluster ManagedClusterAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&cluster).
		SetPathParam("id", id).
		Get(managedClusterEndpointGet)

	if err != nil {
		return nil, c.formatManagedClusterError(managedClusterErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatManagedClusterError(
			managedClusterErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved managed cluster", map[string]any{
		"id":   id,
		"name": cluster.Name,
	})
	return &cluster, nil
}

// ListManagedClusters lists every managed cluster of the tenant, following pagination.
func (c *Client) ListManagedClusters(ctx context.Context) ([]ManagedClusterAPI, error) {
	tflog.Debug(ctx, "Listing managed clusters")

	var all []ManagedClusterAPI
	for offset := 0; ; offset += managedClustersPageSize {
		var page []ManagedClusterAPI
		resp, err := c.prepareRequest(ctx).
			SetResult(&page).
			SetQueryParam("offset", strconv.Itoa(offset)).
			SetQueryParam("limit", strconv.Itoa(managedClustersPageSize)).
			Get(managedClusterEndpointList)

		if err != nil {
			return nil, c.formatManagedClusterError(managedClusterErrorContext{Operation: "list"}, err, 0)
		}
		if resp.IsStatusFailure() {
			return nil, c.formatManagedClusterError(
				managedClusterErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
				nil, resp.StatusCode(),
			)
		}

		all = append(all, page...)
		if len(page) < managedClustersPageSize {
			break
		}
	}

	tflog.Debug(ctx, "Successfully listed managed clusters", map[string]any{"count": len(all)})
	return all, nil
}

func (c *Client) formatManagedClusterError(ctx managedClusterErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch ctx.Operation {
	case "get":
		baseMsg = fmt.Sprintf("failed to get managed cluster '%s'", ctx.ID)
	case "list":
		baseMsg = "failed to list managed clusters"
	default:
		baseMsg = "managed cluster API error"
	}

	if err != nil {
		return fmt.Errorf("%s: %w", baseMsg, err)
	}

	detail := ""
	if ctx.ResponseBody != "" {
		detail = fmt.Sprintf(" - response: %s", ctx.ResponseBody)
	}

	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%s: invalid request (400)%s", baseMsg, detail)
	case http.StatusUnauthorized:
		return fmt.Errorf("%s: authentication failed (401)%s", baseMsg, detail)
	case http.StatusForbidden:
		return fmt.Errorf("%s: access denied (403)%s", baseMsg, detail)
	case http.StatusNotFound:
		return fmt.Errorf("%s: %w", baseMsg, ErrNotFound)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%s: rate limit exceeded (429)%s", baseMsg, detail)
	case http.StatusInternalServerError:
		return fmt.Errorf("%s: server error (500)%s", baseMsg, detail)
	default:
		if statusCode > 0 {
			return fmt.Errorf("%s: unexpected status code %d%s", baseMsg, statusCode, detail)
		}
	}

	return fmt.Errorf("%s: unknown error", baseMsg)
}
