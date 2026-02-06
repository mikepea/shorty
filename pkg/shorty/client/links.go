package client

import (
	"fmt"
	"net/url"
	"strconv"
)

// SearchLinks searches for links across the user's groups.
func (c *Client) SearchLinks(params SearchParams) ([]LinkResponse, error) {
	q := url.Values{}
	if params.Query != "" {
		q.Set("q", params.Query)
	}
	if params.IsUnread != nil {
		q.Set("is_unread", strconv.FormatBool(*params.IsUnread))
	}
	if params.IsPublic != nil {
		q.Set("is_public", strconv.FormatBool(*params.IsPublic))
	}
	if params.GroupID != "" {
		q.Set("group_id", params.GroupID)
	}
	if params.Tag != "" {
		q.Set("tag", params.Tag)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}

	path := "/api/links"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var resp []LinkResponse
	if err := c.request("GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListGroupLinks returns all links in a group.
func (c *Client) ListGroupLinks(groupID string) ([]LinkResponse, error) {
	var resp []LinkResponse
	if err := c.request("GET", fmt.Sprintf("/api/groups/%s/links", groupID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateLink creates a new link in a group.
func (c *Client) CreateLink(groupID string, req CreateLinkRequest) (*LinkResponse, error) {
	var resp LinkResponse
	if err := c.request("POST", fmt.Sprintf("/api/groups/%s/links", groupID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLink returns a link by its slug.
func (c *Client) GetLink(slug string) (*LinkResponse, error) {
	var resp LinkResponse
	if err := c.request("GET", fmt.Sprintf("/api/links/%s", slug), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateLink updates a link by its slug.
func (c *Client) UpdateLink(slug string, req UpdateLinkRequest) (*LinkResponse, error) {
	var resp LinkResponse
	if err := c.request("PUT", fmt.Sprintf("/api/links/%s", slug), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteLink deletes a link by its slug.
func (c *Client) DeleteLink(slug string) error {
	return c.request("DELETE", fmt.Sprintf("/api/links/%s", slug), nil, nil)
}
