package sailpoint_sdk

import "fmt"

const (
	FORM_DEFINITIONS_ENDPOINT = "/v2025/form-definitions"
)

type FormDefinitionApi struct {
	api *Client
}

func NewFormDefinitionApi(client *Client) *FormDefinitionApi {
	return &FormDefinitionApi{
		api: client,
	}
}

func (fdapi *FormDefinitionApi) GetFormDefinitionById(id string) (FormDefinition, error) {
	fd := FormDefinition{}

	_, err := fdapi.api.client.R().
		SetPathParam("id", id).
		SetResult(&fd).
		Get(fmt.Sprintf("%s/{id}", FORM_DEFINITIONS_ENDPOINT))

	return fd, err
}
