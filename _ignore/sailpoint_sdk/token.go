// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sailpoint_sdk

import (
	"time"

	"resty.dev/v3"
)

type OAuth2Token struct {
	AccessToken         string `json:"access_token"`
	TokenType           string `json:"token_type"`
	RefreshToken        string `json:"refresh_token"`
	ExpiresIn           int    `json:"expires_in"`
	Scope               string `json:"scope"`
	AccessType          string `json:"access_type"`
	TenantId            string `json:"tenant_id"`
	Internal            bool   `json:"internal"`
	Pod                 string `json:"pod"`
	StrongAuthSupported bool   `json:"strong_auth_supported"`
	Org                 string `json:"org"`
	UserId              string `json:"user_id"`
	IdentityId          string `json:"identity_id"`
	StrongAuth          bool   `json:"strong_auth"`
	Enabled             bool   `json:"enabled"`
	JTI                 string `json:"jti"`

	generatedAt time.Time `json:"-"`
}

type TokenManager interface {
	GetToken() error
	IsTokenExpired() bool
}

var (
	_ TokenManager = &tokenManager{}
)

type tokenManager struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenUrl     string `json:"token_url"`

	Token *OAuth2Token `json:"token"`
}

func NewTokenManager(clientId, clientSecret, tokenUrl string) *tokenManager {
	tm := &tokenManager{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		TokenUrl:     tokenUrl,
	}

	return tm
}

// GetToken implements TokenManager.
func (t *tokenManager) GetToken() error {
	r := resty.New()
	token := &OAuth2Token{}

	defer r.Close()

	_, err := r.R().
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     t.ClientId,
			"client_secret": t.ClientSecret,
		}).
		SetResult(token).
		Post(t.TokenUrl)
	if err != nil {
		return err
	}

	t.Token = token
	return nil
}

// IsTokenExpired implements TokenManager.
func (t *tokenManager) IsTokenExpired() bool {
	if t.Token == nil {
		return true
	}

	// Consider token expired if it will expire in the next 60 seconds
	expiryTime := t.Token.generatedAt.Add(time.Duration(t.Token.ExpiresIn) * time.Second)
	return time.Now().After(expiryTime.Add(-60 * time.Second))
}
