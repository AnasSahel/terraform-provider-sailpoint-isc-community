package client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	sourcesEndpointV2025 = "/v2025/sources"
)

type Source struct {
	// Core identification
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Type      string `json:"type,omitempty"`
	Connector string `json:"connector"`

	// Descriptive fields
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`

	// References
	Owner               *ObjectRef `json:"owner"`
	Cluster             *ObjectRef `json:"cluster,omitempty"`
	ManagementWorkgroup *ObjectRef `json:"managementWorkgroup,omitempty"`

	// Correlation configuration
	AccountCorrelationConfig  *ObjectRef                       `json:"accountCorrelationConfig,omitempty"`
	AccountCorrelationRule    *ObjectRef                       `json:"accountCorrelationRule,omitempty"`
	ManagerCorrelationMapping *SourceManagerCorrelationMapping `json:"managerCorrelationMapping,omitempty"`
	ManagerCorrelationRule    *ObjectRef                       `json:"managerCorrelationRule,omitempty"`

	// Provisioning rules
	BeforeProvisioningRule *ObjectRef  `json:"beforeProvisioningRule,omitempty"`
	DeleteThreshold        int32       `json:"deleteThreshold,omitempty"`
	PasswordPolicies       []ObjectRef `json:"passwordPolicies,omitempty"`

	// Connector details
	ConnectorClass            string                 `json:"connectorClass,omitempty"`
	ConnectorID               string                 `json:"connectorId,omitempty"`
	ConnectorName             string                 `json:"connectorName,omitempty"`
	ConnectorType             string                 `json:"connectorType,omitempty"`
	ConnectorImplementationID string                 `json:"connectorImplementationId,omitempty"`
	ConnectorAttributes       map[string]interface{} `json:"connectorAttributes,omitempty"`

	// Schema and policies
	Schemas  []ObjectRef `json:"schemas,omitempty"`
	Features []string    `json:"features,omitempty"`

	// Status and metadata
	Authoritative             bool   `json:"authoritative,omitempty"`
	CredentialProviderEnabled bool   `json:"credentialProviderEnabled,omitempty"`
	Healthy                   bool   `json:"healthy,omitempty"`
	Status                    string `json:"status,omitempty"`
	Since                     string `json:"since,omitempty"`
	Created                   string `json:"created,omitempty"`
	Modified                  string `json:"modified,omitempty"`
}

type SourceManagerCorrelationMapping struct {
	AccountAttributeName  string `json:"accountAttributeName"`
	IdentityAttributeName string `json:"identityAttributeName"`
}

func (c *Client) GetSource(ctx context.Context, id string) (*Source, error) {
	var result Source

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", sourcesEndpointV2025, id), nil, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "get",
			Resource:   "source",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatError(ErrorContext{
		Operation:  "get",
		Resource:   "source",
		ResourceID: id,
	}, nil, resp.StatusCode())
}

func (c *Client) CreateSource(ctx context.Context, source *Source) (*Source, error) {
	var result Source

	resp, err := c.doRequest(ctx, http.MethodPost, sourcesEndpointV2025, source, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation: "create",
			Resource:  "source",
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusCreated {
		return &result, nil
	}

	return nil, c.formatError(ErrorContext{
		Operation: "create",
		Resource:  "source",
	}, nil, resp.StatusCode())
}

func (c *Client) DeleteSource(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", sourcesEndpointV2025, id), nil, nil)

	if err != nil {
		return c.formatError(ErrorContext{
			Operation:  "delete",
			Resource:   "source",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusAccepted {
		return nil
	}

	return c.formatError(ErrorContext{
		Operation:  "delete",
		Resource:   "source",
		ResourceID: id,
	}, nil, resp.StatusCode())
}

func (c *Client) PatchSource(ctx context.Context, id string, patches []JSONPatchOperation) (*Source, error) {
	var result Source

	resp, err := c.doRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/%s", sourcesEndpointV2025, id), patches, &result)

	if err != nil {
		return nil, c.formatError(ErrorContext{
			Operation:  "update",
			Resource:   "source",
			ResourceID: id,
		}, err, 0)
	}

	if resp.StatusCode() == http.StatusOK {
		return &result, nil
	}

	return nil, c.formatError(ErrorContext{
		Operation:  "update",
		Resource:   "source",
		ResourceID: id,
	}, nil, resp.StatusCode())
}
