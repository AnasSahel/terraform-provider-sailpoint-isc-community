package client

type TokenResponse struct {
	AccessToken         string `json:"access_token"`
	TokenType           string `json:"token_type"`
	RefreshToken        string `json:"refresh_token"`
	ExpiresIn           int    `json:"expires_in"`
	Scope               string `json:"scope"`
	AccessType          string `json:"accessType"`
	TenantID            string `json:"tenant_id"`
	Internal            bool   `json:"internal"`
	Pod                 string `json:"pod"`
	StrongAuthSupported bool   `json:"strong_auth_supported"`
	Org                 string `json:"org"`
	UserID              string `json:"user_id"`
	IdentityID          string `json:"identity_id"`
	StrongAuth          bool   `json:"strong_auth"`
	Enabled             bool   `json:"enabled"`
	JTI                 string `json:"jti"`
}
