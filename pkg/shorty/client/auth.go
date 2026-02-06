package client

// Register creates a new user account and stores the returned JWT token.
func (c *Client) Register(email, password, name string) (*AuthResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
		"name":     name,
	}
	var resp AuthResponse
	if err := c.request("POST", "/api/auth/register", body, &resp); err != nil {
		return nil, err
	}
	c.token = resp.Token
	return &resp, nil
}

// Login authenticates a user and stores the returned JWT token.
func (c *Client) Login(email, password string) (*AuthResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	var resp AuthResponse
	if err := c.request("POST", "/api/auth/login", body, &resp); err != nil {
		return nil, err
	}
	c.token = resp.Token
	return &resp, nil
}

// Logout clears the local token (server-side logout is a no-op).
func (c *Client) Logout() error {
	err := c.request("POST", "/api/auth/logout", nil, nil)
	c.token = ""
	return err
}

// Me returns the authenticated user's profile.
func (c *Client) Me() (*UserResponse, error) {
	var resp UserResponse
	if err := c.request("GET", "/api/auth/me", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangePassword changes the authenticated user's password.
func (c *Client) ChangePassword(currentPassword, newPassword string) error {
	body := map[string]string{
		"current_password": currentPassword,
		"new_password":     newPassword,
	}
	return c.request("PUT", "/api/auth/password", body, nil)
}
