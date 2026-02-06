package client

import "fmt"

// ListOrganizations returns all organizations the authenticated user belongs to.
func (c *Client) ListOrganizations() ([]OrgResponse, error) {
	var resp []OrgResponse
	if err := c.request("GET", "/api/organizations", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateOrganization creates a new organization.
func (c *Client) CreateOrganization(name, slug string) (*OrgResponse, error) {
	body := map[string]string{
		"name": name,
		"slug": slug,
	}
	var resp OrgResponse
	if err := c.request("POST", "/api/organizations", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrganization returns a specific organization by ID.
func (c *Client) GetOrganization(id string) (*OrgResponse, error) {
	var resp OrgResponse
	if err := c.request("GET", fmt.Sprintf("/api/organizations/%s", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateOrganization updates an organization's name.
func (c *Client) UpdateOrganization(id, name string) (*OrgResponse, error) {
	body := map[string]string{"name": name}
	var resp OrgResponse
	if err := c.request("PUT", fmt.Sprintf("/api/organizations/%s", id), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteOrganization deletes an organization.
func (c *Client) DeleteOrganization(id string) error {
	return c.request("DELETE", fmt.Sprintf("/api/organizations/%s", id), nil, nil)
}

// ListOrgMembers returns all members of an organization.
func (c *Client) ListOrgMembers(orgID string) ([]OrgMemberResponse, error) {
	var resp []OrgMemberResponse
	if err := c.request("GET", fmt.Sprintf("/api/organizations/%s/members", orgID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddOrgMember adds a user to an organization by email.
func (c *Client) AddOrgMember(orgID, email, role string) (*OrgMemberResponse, error) {
	body := map[string]string{
		"email": email,
		"role":  role,
	}
	var resp OrgMemberResponse
	if err := c.request("POST", fmt.Sprintf("/api/organizations/%s/members", orgID), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateOrgMember updates a member's role in an organization.
func (c *Client) UpdateOrgMember(orgID, userID, role string) (*OrgMemberResponse, error) {
	body := map[string]string{"role": role}
	var resp OrgMemberResponse
	if err := c.request("PUT", fmt.Sprintf("/api/organizations/%s/members/%s", orgID, userID), body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RemoveOrgMember removes a member from an organization.
func (c *Client) RemoveOrgMember(orgID, userID string) error {
	return c.request("DELETE", fmt.Sprintf("/api/organizations/%s/members/%s", orgID, userID), nil, nil)
}
