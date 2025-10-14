package client

import (
	"fmt"
)

const (
	formDefinitionsBasePath = "/v3/form-definitions"
)

// APIFormDefinition handles all form definition operations
type APIFormDefinition struct {
	client *Client
}

// NewAPIFormDefinition creates a new form definition API client
func NewAPIFormDefinition(client *Client) *APIFormDefinition {
	return &APIFormDefinition{
		client: client,
	}
}

// FormDefinition represents a SailPoint form definition
type FormDefinition struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	Owner          *FormOwner      `json:"owner,omitempty"`
	UsedBy         []FormUsedBy    `json:"usedBy,omitempty"`
	FormInput      []FormElement   `json:"formInput,omitempty"`
	FormElements   []FormElement   `json:"formElements,omitempty"`
	FormConditions []FormCondition `json:"formConditions,omitempty"`
	Created        string          `json:"created,omitempty"`
	Modified       string          `json:"modified,omitempty"`
}

// FormOwner represents the owner of a form definition
type FormOwner struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FormUsedBy represents where a form is used
type FormUsedBy struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// FormElement represents a form field element
type FormElement struct {
	ID                 string                 `json:"id,omitempty"`
	ElementType        string                 `json:"elementType,omitempty"`
	Config             map[string]interface{} `json:"config,omitempty"`
	Key                string                 `json:"key,omitempty"`
	Label              string                 `json:"label,omitempty"`
	Required           bool                   `json:"required,omitempty"`
	HelpText           string                 `json:"helpText,omitempty"`
	ValidationMessages map[string]string      `json:"validationMessages,omitempty"`
}

// FormCondition represents a conditional logic in a form
type FormCondition struct {
	RuleOperator string                   `json:"ruleOperator,omitempty"`
	Rules        []map[string]interface{} `json:"rules,omitempty"`
	Effects      []FormEffect             `json:"effects,omitempty"`
}

// FormEffect represents the effect of a condition
type FormEffect struct {
	EffectType string                 `json:"effectType,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// ListFormDefinitionsOptions contains options for listing form definitions
type ListFormDefinitionsOptions struct {
	Offset  *int
	Limit   *int
	Filters *string
	Sorters *string
}

// List retrieves all form definitions
func (api *APIFormDefinition) List(options *ListFormDefinitionsOptions) ([]FormDefinition, error) {
	builder := api.client.Get(formDefinitionsBasePath)

	if options != nil {
		if options.Offset != nil {
			builder.QueryParam("offset", fmt.Sprintf("%d", *options.Offset))
		}
		if options.Limit != nil {
			builder.QueryParam("limit", fmt.Sprintf("%d", *options.Limit))
		}
		if options.Filters != nil {
			builder.QueryParam("filters", *options.Filters)
		}
		if options.Sorters != nil {
			builder.QueryParam("sorters", *options.Sorters)
		}
	}

	var formDefinitions []FormDefinition
	if err := builder.ExecuteJSON(&formDefinitions); err != nil {
		return nil, fmt.Errorf("failed to list form definitions: %w", err)
	}

	return formDefinitions, nil
}

// Get retrieves a specific form definition by ID
func (api *APIFormDefinition) Get(id string) (*FormDefinition, error) {
	path := fmt.Sprintf("%s/%s", formDefinitionsBasePath, id)

	var formDefinition FormDefinition
	if err := api.client.Get(path).ExecuteJSON(&formDefinition); err != nil {
		return nil, fmt.Errorf("failed to get form definition %s: %w", id, err)
	}

	return &formDefinition, nil
}

// Create creates a new form definition
func (api *APIFormDefinition) Create(formDefinition *FormDefinition) (*FormDefinition, error) {
	var result FormDefinition
	if err := api.client.Post(formDefinitionsBasePath).
		Body(formDefinition).
		ExecuteJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to create form definition: %w", err)
	}

	return &result, nil
}

// Update updates an existing form definition
func (api *APIFormDefinition) Update(id string, formDefinition *FormDefinition) (*FormDefinition, error) {
	path := fmt.Sprintf("%s/%s", formDefinitionsBasePath, id)

	var result FormDefinition
	if err := api.client.Put(path).
		Body(formDefinition).
		ExecuteJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to update form definition %s: %w", id, err)
	}

	return &result, nil
}

// Patch partially updates a form definition using JSON patch operations
func (api *APIFormDefinition) Patch(id string, operations []map[string]interface{}) (*FormDefinition, error) {
	path := fmt.Sprintf("%s/%s", formDefinitionsBasePath, id)

	var result FormDefinition
	if err := api.client.Patch(path).
		Body(operations).
		ExecuteJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to patch form definition %s: %w", id, err)
	}

	return &result, nil
}

// Delete deletes a form definition
func (api *APIFormDefinition) Delete(id string) error {
	path := fmt.Sprintf("%s/%s", formDefinitionsBasePath, id)

	if err := api.client.Delete(path).ExecuteNoContent(); err != nil {
		return fmt.Errorf("failed to delete form definition %s: %w", id, err)
	}

	return nil
}

// Export exports a form definition (returns the JSON payload)
func (api *APIFormDefinition) Export(id string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s/export", formDefinitionsBasePath, id)

	data, err := api.client.Get(path).ExecuteBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to export form definition %s: %w", id, err)
	}

	return data, nil
}

// Import imports a form definition from JSON payload
func (api *APIFormDefinition) Import(data []byte) (*FormDefinition, error) {
	path := fmt.Sprintf("%s/import", formDefinitionsBasePath)

	var result FormDefinition
	if err := api.client.Post(path).
		Header("Content-Type", "application/json").
		Body(data).
		ExecuteJSON(&result); err != nil {
		return nil, fmt.Errorf("failed to import form definition: %w", err)
	}

	return &result, nil
}
