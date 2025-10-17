package sailpoint_sdk

import (
	"fmt"
	"time"

	"resty.dev/v3"
)

const (
	DEFAULT_TIMEOUT     = 30 * time.Second
	USER_AGENT          = "terraform-provider-sailpoint"
	RETRY_COUNT         = 3
	RETRY_WAIT_TIME     = 5 * time.Second
	RETRY_MAX_WAIT_TIME = 20 * time.Second
	CONTENT_TYPE        = "application/json"
)

type Client struct {
	baseUrl      string
	clientId     string
	clientSecret string

	tokenUrl string

	client       *resty.Client
	tokenManager *tokenManager

	FormDefinitionApi *FormDefinitionApi
}

func NewClient(baseUrl string, clientId string, clientSecret string) *Client {
	localClient := &Client{
		baseUrl:      baseUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		tokenUrl:     fmt.Sprintf("%s/oauth/token", baseUrl),
		client:       resty.New(),
	}

	localClient.tokenManager = NewTokenManager(clientId, clientSecret, localClient.tokenUrl)

	localClient.client.
		SetBaseURL(baseUrl).
		SetTimeout(DEFAULT_TIMEOUT).
		SetHeader("Content-Type", CONTENT_TYPE).
		SetHeader("User-Agent", USER_AGENT).
		SetRetryCount(RETRY_COUNT).
		SetRetryWaitTime(RETRY_WAIT_TIME).
		SetRetryMaxWaitTime(RETRY_MAX_WAIT_TIME)

	localClient.client.AddRequestMiddleware(func(c *resty.Client, r *resty.Request) error {
		if localClient.tokenManager.IsTokenExpired() {
			err := localClient.tokenManager.GetToken()
			if err != nil {
				return err
			}
		}

		r.SetAuthToken(localClient.tokenManager.Token.AccessToken)
		return nil
	})

	// Add APIs initialization here
	localClient.FormDefinitionApi = NewFormDefinitionApi(localClient)

	return localClient
}
