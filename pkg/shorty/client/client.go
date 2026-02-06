package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents a structured error returned by the Shorty API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Client is an HTTP client for the Shorty API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	orgID      string
}

// New creates a new Shorty API client.
// If httpClient is nil, http.DefaultClient is used.
func New(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// SetToken sets the authentication token (JWT or API key).
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetOrgID sets the X-Organization-ID header for subsequent requests.
func (c *Client) SetOrgID(orgID string) {
	c.orgID = orgID
}

// ClearOrgID removes the X-Organization-ID header.
func (c *Client) ClearOrgID() {
	c.orgID = ""
}

// request performs an HTTP request and decodes the JSON response.
// If result is nil, the response body is discarded.
func (c *Client) request(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.orgID != "" {
		req.Header.Set("X-Organization-ID", c.orgID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errBody ErrorBody
		if json.Unmarshal(respBody, &errBody) == nil && errBody.Error != "" {
			return &APIError{StatusCode: resp.StatusCode, Message: errBody.Error}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: string(respBody)}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}
