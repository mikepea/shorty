package client

import "fmt"

// ListAPIKeys returns all API keys for the authenticated user.
func (c *Client) ListAPIKeys() ([]APIKeyResponse, error) {
	var resp []APIKeyResponse
	if err := c.request("GET", "/api/api-keys", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateAPIKey creates a new API key.
func (c *Client) CreateAPIKey(description string) (*CreateAPIKeyResponse, error) {
	body := map[string]string{"description": description}
	var resp CreateAPIKeyResponse
	if err := c.request("POST", "/api/api-keys", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAPIKey deletes an API key by ID.
func (c *Client) DeleteAPIKey(id string) error {
	return c.request("DELETE", fmt.Sprintf("/api/api-keys/%s", id), nil, nil)
}
