package client

import (
	"context"
	"fmt"
	"net/http"
)

type Source struct {
	ID                        string                           `json:"id,omitempty"`
	Name                      string                           `json:"name"`
	Description               string                           `json:"description,omitempty"`
	Owner                     *ObjectRef                       `json:"owner"`
	Cluster                   *ObjectRef                       `json:"cluster,omitempty"`
	AccountCorrelationConfig  *ObjectRef                       `json:"accountCorrelationConfig,omitempty"`
	AccountCorrelationRule    *ObjectRef                       `json:"accountCorrelationRule,omitempty"`
	ManagerCorrelationMapping *SourceManagerCorrelationMapping `json:"managerCorrelationMapping,omitempty"`
	ManagerCorrelationRule    *ObjectRef                       `json:"managerCorrelationRule,omitempty"`
	BeforeProvisioningRule    *ObjectRef                       `json:"beforeProvisioningRule,omitempty"`
	Schemas                   []ObjectRef                      `json:"schemas,omitempty"`
	PasswordPolicies          []ObjectRef                      `json:"passwordPolicies,omitempty"`
	Features                  []string                         `json:"features,omitempty"`
	Type                      string                           `json:"type,omitempty"`
	Connector                 string                           `json:"connector"`
	ConnectorClass            string                           `json:"connectorClass,omitempty"`
	ConnectorAttributes       map[string]interface{}           `json:"connectorAttributes,omitempty"`
	DeleteThreshold           int32                            `json:"deleteThreshold,omitempty"`
	Authoritative             bool                             `json:"authoritative,omitempty"`
	ManagementWorkgroup       *ObjectRef                       `json:"managementWorkgroup,omitempty"`
	Healthy                   bool                             `json:"healthy,omitempty"`
	Status                    string                           `json:"status,omitempty"`
	Since                     string                           `json:"since,omitempty"`
	ConnectorID               string                           `json:"connectorId,omitempty"`
	ConnectorName             string                           `json:"connectorName,omitempty"`
	ConnectorType             string                           `json:"connectorType,omitempty"`
	ConnectorImplementationID string                           `json:"connectorImplementationId,omitempty"`
	Created                   string                           `json:"created,omitempty"`
	Modified                  string                           `json:"modified,omitempty"`
	CredentialProviderEnabled bool                             `json:"credentialProviderEnabled,omitempty"`
	Category                  string                           `json:"category,omitempty"`
}

type SourceManagerCorrelationMapping struct {
	AccountAttributeName  string `json:"accountAttributeName"`
	IdentityAttributeName string `json:"identityAttributeName"`
}

func (c *Client) GetSource(ctx context.Context, id string) (*Source, error) {
	var result Source

	resp, err := c.HTTPClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(fmt.Sprintf("/v2025/sources/%s", id))

	if err != nil {
		return nil, fmt.Errorf("getting source with ID %q: %v", id, err)
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("source with ID %q not found", id)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	return &result, nil
}

func (c *Client) CreateSource(ctx context.Context, source *Source) (*Source, error) {
	var result Source

	resp, err := c.HTTPClient.R().
		SetContext(ctx).
		SetBody(source).
		SetResult(&result).
		Post("/v2025/sources")

	if err != nil {
		return nil, fmt.Errorf("creating source: %w", err)
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return nil, fmt.Errorf("bad request when creating source")
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized when creating source")
	}

	if resp.StatusCode() == http.StatusForbidden {
		return nil, fmt.Errorf("forbidden when creating source")
	}

	if resp.StatusCode() == http.StatusInternalServerError {
		return nil, fmt.Errorf("internal server error when creating source")
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode())
	}

	return &result, nil
}

// type Owner struct {
// 	Type string `json:"type"`
// 	ID   string `json:"id"`
// 	Name string `json:"name,omitempty"`
// }

// type ErrorResponse struct {
// 	DetailCode string         `json:"detailCode"`
// 	Messages   []ErrorMessage `json:"messages"`
// 	TrackingID string         `json:"trackingId"`
// }

// type ErrorMessage struct {
// 	Locale string `json:"locale"`
// 	Text   string `json:"text"`
// 	Key    string `json:"key"`
// }

// func (c *Client) UpdateSource(ctx context.Context, id string, source *Source) (*Source, error) {
// 	var result Source
// 	var errResp ErrorResponse

// 	resp, err := c.HTTPClient.R().
// 		SetContext(ctx).
// 		SetBody(source).
// 		SetResult(&result).
// 		SetError(&errResp).
// 		Patch(fmt.Sprintf("/v3/sources/%s", id))

// 	if err != nil {
// 		return nil, fmt.Errorf("updating source: %w", err)
// 	}

// 	if resp.StatusCode() == http.StatusNotFound {
// 		return nil, ErrResourceNotFound
// 	}

// 	if resp.IsError() {
// 		return nil, c.handleErrorResponse(resp.StatusCode(), &errResp)
// 	}

// 	return &result, nil
// }

// func (c *Client) DeleteSource(ctx context.Context, id string) error {
// 	var errResp ErrorResponse

// 	resp, err := c.HTTPClient.R().
// 		SetContext(ctx).
// 		SetError(&errResp).
// 		Delete(fmt.Sprintf("/v3/sources/%s", id))

// 	if err != nil {
// 		return fmt.Errorf("deleting source: %w", err)
// 	}

// 	if resp.StatusCode() == http.StatusNotFound {
// 		// Already deleted, not an error
// 		return nil
// 	}

// 	if resp.IsError() {
// 		return c.handleErrorResponse(resp.StatusCode(), &errResp)
// 	}

// 	return nil
// }

// func (c *Client) handleErrorResponse(statusCode int, errResp *ErrorResponse) error {
// 	if errResp.DetailCode != "" {
// 		msg := errResp.DetailCode
// 		if len(errResp.Messages) > 0 {
// 			msg += ": " + errResp.Messages[0].Text
// 		}
// 		return fmt.Errorf("API error (status %d, tracking %s): %s",
// 			statusCode, errResp.TrackingID, msg)
// 	}
// 	return fmt.Errorf("API request failed with status %d", statusCode)
// }
