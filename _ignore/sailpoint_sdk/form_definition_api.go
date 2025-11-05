// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sailpoint_sdk

import (
	"context"
	"fmt"
)

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

func (fdapi *FormDefinitionApi) GetFormDefinitionById(ctx context.Context, id string) (map[string]interface{}, error) {
	fd := map[string]interface{}{}

	_, err := fdapi.api.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetResult(&fd).
		Get(fmt.Sprintf("%s/{id}", FORM_DEFINITIONS_ENDPOINT))
	return fd, err

}

func (fdapi *FormDefinitionApi) CreateFormDefinition(ctx context.Context, formDef map[string]interface{}) (map[string]interface{}, error) {
	fd := map[string]interface{}{}

	_, err := fdapi.api.client.R().
		SetContext(ctx).
		SetBody(formDef).
		SetResult(&fd).
		Post(FORM_DEFINITIONS_ENDPOINT)
	return fd, err
}

func (fdapi *FormDefinitionApi) PatchFormDefinition(ctx context.Context, id string, patches []map[string]interface{}) (map[string]interface{}, error) {
	fd := map[string]interface{}{}

	_, err := fdapi.api.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetBody(patches).
		SetResult(&fd).
		Patch(fmt.Sprintf("%s/{id}", FORM_DEFINITIONS_ENDPOINT))
	return fd, err
}

func (fdapi *FormDefinitionApi) DeleteFormDefinition(ctx context.Context, id string) error {
	_, err := fdapi.api.client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		Delete(fmt.Sprintf("%s/{id}", FORM_DEFINITIONS_ENDPOINT))
	return err
}
