package client

import "fmt"

// ListTags returns all tags across the user's groups.
func (c *Client) ListTags() ([]TagResponse, error) {
	var resp []TagResponse
	if err := c.request("GET", "/api/tags", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListGroupTags returns all tags used in a specific group.
func (c *Client) ListGroupTags(groupID string) ([]TagResponse, error) {
	var resp []TagResponse
	if err := c.request("GET", fmt.Sprintf("/api/groups/%s/tags", groupID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLinkTags returns the tags for a specific link.
func (c *Client) GetLinkTags(slug string) ([]TagResponse, error) {
	var resp []TagResponse
	if err := c.request("GET", fmt.Sprintf("/api/links/%s/tags", slug), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetLinkTags replaces all tags on a link.
func (c *Client) SetLinkTags(slug string, tags []string) ([]TagResponse, error) {
	body := map[string][]string{"tags": tags}
	var resp []TagResponse
	if err := c.request("PUT", fmt.Sprintf("/api/links/%s/tags", slug), body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddLinkTag adds a single tag to a link.
func (c *Client) AddLinkTag(slug, tagName string) (*TagResponse, error) {
	var resp TagResponse
	if err := c.request("POST", fmt.Sprintf("/api/links/%s/tags/%s", slug, tagName), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveLinkTag removes a tag from a link.
func (c *Client) RemoveLinkTag(slug, tagName string) error {
	return c.request("DELETE", fmt.Sprintf("/api/links/%s/tags/%s", slug, tagName), nil, nil)
}
