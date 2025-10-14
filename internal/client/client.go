package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DEFAULT_TIMEOUT = 30 * time.Second
	USER_AGENT      = "terraform-provider-sailpoint"
)

type Client struct {
	baseUrl      string
	clientId     string
	clientSecret string
	tokenUrl     string
	httpClient   *http.Client
	token        *TokenResponse

	// Specialized API clients
	FormDefinitions *APIFormDefinition
}

func NewClient(baseUrl string, clientId string, clientSecret string) *Client {
	client := &Client{
		baseUrl:      baseUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		tokenUrl:     fmt.Sprintf("%s/oauth/token", baseUrl),
		httpClient: &http.Client{
			Timeout: DEFAULT_TIMEOUT,
		},
	}

	client.token, _ = client.getAccessToken()

	// Initialize specialized API clients
	client.FormDefinitions = NewAPIFormDefinition(client)

	return client
}

func (c *Client) getAccessToken() (*TokenResponse, error) {
	// Prepare the request
	req, err := http.PostForm(c.tokenUrl, map[string][]string{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
	})

	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	// Check for HTTP errors
	if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get access token: %s", req.Status)
	}

	// Decode the response
	var tokenResp TokenResponse
	if err := json.NewDecoder(req.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// RequestBuilder provides a fluent interface for building and executing HTTP requests
type RequestBuilder struct {
	client      *Client
	method      string
	path        string
	queryParams url.Values
	headers     map[string]string
	body        interface{}
	err         error
}

// newRequest creates a new RequestBuilder
func (c *Client) newRequest(method, path string) *RequestBuilder {
	return &RequestBuilder{
		client:      c,
		method:      method,
		path:        path,
		queryParams: url.Values{},
		headers:     make(map[string]string),
	}
}

// Get creates a GET request builder
func (c *Client) Get(path string) *RequestBuilder {
	return c.newRequest(http.MethodGet, path)
}

// Post creates a POST request builder
func (c *Client) Post(path string) *RequestBuilder {
	return c.newRequest(http.MethodPost, path)
}

// Put creates a PUT request builder
func (c *Client) Put(path string) *RequestBuilder {
	return c.newRequest(http.MethodPut, path)
}

// Patch creates a PATCH request builder
func (c *Client) Patch(path string) *RequestBuilder {
	return c.newRequest(http.MethodPatch, path)
}

// Delete creates a DELETE request builder
func (c *Client) Delete(path string) *RequestBuilder {
	return c.newRequest(http.MethodDelete, path)
}

// QueryParam adds a query parameter to the request
func (rb *RequestBuilder) QueryParam(key, value string) *RequestBuilder {
	if rb.err != nil {
		return rb
	}
	rb.queryParams.Add(key, value)
	return rb
}

// QueryParams adds multiple query parameters to the request
func (rb *RequestBuilder) QueryParams(params map[string]string) *RequestBuilder {
	if rb.err != nil {
		return rb
	}
	for key, value := range params {
		rb.queryParams.Add(key, value)
	}
	return rb
}

// Header adds a header to the request
func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	if rb.err != nil {
		return rb
	}
	rb.headers[key] = value
	return rb
}

// Headers adds multiple headers to the request
func (rb *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	if rb.err != nil {
		return rb
	}
	for key, value := range headers {
		rb.headers[key] = value
	}
	return rb
}

// Body sets the request body (will be JSON encoded)
func (rb *RequestBuilder) Body(body interface{}) *RequestBuilder {
	if rb.err != nil {
		return rb
	}
	rb.body = body
	return rb
}

// Execute executes the request and returns the response
func (rb *RequestBuilder) Execute() (*http.Response, error) {
	if rb.err != nil {
		return nil, rb.err
	}

	// Ensure we have a valid token
	if rb.client.token == nil {
		token, err := rb.client.getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get access token: %w", err)
		}
		rb.client.token = token
	}

	// Build URL
	requestURL := fmt.Sprintf("%s%s", rb.client.baseUrl, rb.path)
	if len(rb.queryParams) > 0 {
		requestURL = fmt.Sprintf("%s?%s", requestURL, rb.queryParams.Encode())
	}

	// Prepare request body
	var bodyReader io.Reader
	if rb.body != nil {
		jsonData, err := json.Marshal(rb.body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	// Create HTTP request
	req, err := http.NewRequest(rb.method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rb.client.token.AccessToken))
	req.Header.Set("User-Agent", USER_AGENT)
	if rb.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range rb.headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := rb.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// ExecuteJSON executes the request and decodes the JSON response into the provided interface
func (rb *RequestBuilder) ExecuteJSON(result interface{}) error {
	resp, err := rb.Execute()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// ExecuteBytes executes the request and returns the response body as bytes
func (rb *RequestBuilder) ExecuteBytes() ([]byte, error) {
	resp, err := rb.Execute()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// ExecuteNoContent executes the request and expects no content response (204)
func (rb *RequestBuilder) ExecuteNoContent() error {
	resp, err := rb.Execute()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
