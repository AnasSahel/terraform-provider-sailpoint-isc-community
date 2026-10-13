// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	managedClientEndpointList = "/managed-clients/v1"
	managedClientEndpointGet  = "/managed-clients/v1/{id}"

	managedClientsPageSize = 250
)

// ManagedClientAPI represents a SailPoint managed client (virtual appliance or CCG).
// The secret, alertKey, vaDownloadUrl, clientId, cookbook and ccId fields returned by the API are
// deliberately not mapped, so they never reach Terraform state.
type ManagedClientAPI struct {
	ID                string          `json:"id,omitempty"`
	ClusterID         string          `json:"clusterId,omitempty"`
	Name              *string         `json:"name,omitempty"`
	Description       *string         `json:"description,omitempty"`
	Type              *string         `json:"type,omitempty"`
	ClusterType       *string         `json:"clusterType,omitempty"`
	Status            *string         `json:"status,omitempty"`
	IPAddress         *string         `json:"ipAddress,omitempty"`
	VaVersion         *string         `json:"vaVersion,omitempty"`
	LastSeen          *string         `json:"lastSeen,omitempty"`
	SinceLastSeen     json.RawMessage `json:"sinceLastSeen,omitempty"`
	ProvisionStatus   *string         `json:"provisionStatus,omitempty"`
	APIGatewayBaseURL *string         `json:"apiGatewayBaseUrl,omitempty"`
	CreatedAt         *string         `json:"createdAt,omitempty"`
	UpdatedAt         *string         `json:"updatedAt,omitempty"`
}

type managedClientErrorContext struct {
	Operation    string
	ID           string
	ResponseBody string
}

// GetManagedClient retrieves a managed client by ID.
func (c *Client) GetManagedClient(ctx context.Context, id string) (*ManagedClientAPI, error) {
	if id == "" {
		return nil, fmt.Errorf("managed client ID cannot be empty")
	}

	tflog.Debug(ctx, "Getting managed client", map[string]any{"id": id})

	var mc ManagedClientAPI
	resp, err := c.prepareRequest(ctx).
		SetResult(&mc).
		SetPathParam("id", id).
		Get(managedClientEndpointGet)

	if err != nil {
		return nil, c.formatManagedClientError(managedClientErrorContext{Operation: "get", ID: id}, err, 0)
	}
	if resp.IsStatusFailure() {
		return nil, c.formatManagedClientError(
			managedClientErrorContext{Operation: "get", ID: id, ResponseBody: string(resp.Bytes())},
			nil, resp.StatusCode(),
		)
	}

	tflog.Debug(ctx, "Successfully retrieved managed client", map[string]any{"id": id})
	return &mc, nil
}

// ListManagedClients retrieves all managed clients, following pagination.
// The endpoint has no reliable server-side filter on name or cluster, so callers
// match on the returned list.
func (c *Client) ListManagedClients(ctx context.Context) ([]ManagedClientAPI, error) {
	tflog.Debug(ctx, "Listing managed clients")

	var all []ManagedClientAPI
	for offset := 0; ; offset += managedClientsPageSize {
		var page []ManagedClientAPI
		resp, err := c.prepareRequest(ctx).
			SetResult(&page).
			SetQueryParam("offset", strconv.Itoa(offset)).
			SetQueryParam("limit", strconv.Itoa(managedClientsPageSize)).
			Get(managedClientEndpointList)

		if err != nil {
			return nil, c.formatManagedClientError(managedClientErrorContext{Operation: "list"}, err, 0)
		}
		if resp.IsStatusFailure() {
			return nil, c.formatManagedClientError(
				managedClientErrorContext{Operation: "list", ResponseBody: string(resp.Bytes())},
				nil, resp.StatusCode(),
			)
		}

		all = append(all, page...)
		if len(page) < managedClientsPageSize {
			break
		}
	}

	tflog.Debug(ctx, "Successfully listed managed clients", map[string]any{"count": len(all)})
	return all, nil
}

func (c *Client) formatManagedClientError(ctx managedClientErrorContext, err error, statusCode int) error {
	var baseMsg string
	switch ctx.Operation {
	case "get":
		baseMsg = fmt.Sprintf("failed to get managed client '%s'", ctx.ID)
	case "list":
		baseMsg = "failed to list managed clients"
	default:
		baseMsg = "managed client API error"
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
