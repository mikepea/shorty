package client

import "fmt"

// ListGroups returns all groups the authenticated user belongs to.
func (c *Client) ListGroups() ([]GroupResponse, error) {
	var resp []GroupResponse
	if err := c.request("GET", "/api/groups", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateGroup creates a new group.
func (c *Client) CreateGroup(req CreateGroupRequest) (*GroupResponse, error) {
	var resp GroupResponse
	if err := c.request("POST", "/api/groups", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGroup returns a specific group by ID.
func (c *Client) GetGroup(id string) (*GroupResponse, error) {
	var resp GroupResponse
	if err := c.request("GET", fmt.Sprintf("/api/groups/%s", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateGroup updates a group's name and/or description.
func (c *Client) UpdateGroup(id, name, description string) (*GroupResponse, error) {
	body := map[string]string{}
	if name != "" {
		body["name"] = name
	}
	if description != "" {
		body["description"] = description
	}
	var resp GroupResponse
	if err := c.request("PUT", fmt.Sprintf("/api/groups/%s", id), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteGroup deletes a group.
func (c *Client) DeleteGroup(id string) error {
	return c.request("DELETE", fmt.Sprintf("/api/groups/%s", id), nil, nil)
}

// ListGroupMembers returns all members of a group.
func (c *Client) ListGroupMembers(groupID string) ([]GroupMemberResponse, error) {
	var resp []GroupMemberResponse
	if err := c.request("GET", fmt.Sprintf("/api/groups/%s/members", groupID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddGroupMember adds a user to a group by email.
func (c *Client) AddGroupMember(groupID, email, role string) (*GroupMemberResponse, error) {
	body := map[string]string{
		"email": email,
		"role":  role,
	}
	var resp GroupMemberResponse
	if err := c.request("POST", fmt.Sprintf("/api/groups/%s/members", groupID), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateGroupMember updates a member's role in a group.
func (c *Client) UpdateGroupMember(groupID, userID, role string) (*GroupMemberResponse, error) {
	body := map[string]string{"role": role}
	var resp GroupMemberResponse
	if err := c.request("PUT", fmt.Sprintf("/api/groups/%s/members/%s", groupID, userID), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveGroupMember removes a member from a group.
func (c *Client) RemoveGroupMember(groupID, userID string) error {
	return c.request("DELETE", fmt.Sprintf("/api/groups/%s/members/%s", groupID, userID), nil, nil)
}
