package client

import "time"

// AuthResponse represents the response from register/login endpoints.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents user data in API responses.
type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	SystemRole  string `json:"system_role"`
	HasPassword bool   `json:"has_password"`
}

// OrgResponse represents an organization in API responses.
type OrgResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	IsGlobal    bool   `json:"is_global"`
	Role        string `json:"role,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// OrgMemberResponse represents an organization member in API responses.
type OrgMemberResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// GroupResponse represents a group in API responses.
type GroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
}

// GroupMemberResponse represents a group member in API responses.
type GroupMemberResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// CreateGroupRequest is the request body for creating a group.
type CreateGroupRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// LinkResponse represents a link in API responses.
type LinkResponse struct {
	ID          string `json:"id"`
	GroupID     string `json:"group_id"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsUnread    bool   `json:"is_unread"`
	ClickCount  uint   `json:"click_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateLinkRequest is the request body for creating a link.
type CreateLinkRequest struct {
	URL         string `json:"url"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public,omitempty"`
	IsUnread    bool   `json:"is_unread,omitempty"`
}

// UpdateLinkRequest is the request body for updating a link.
type UpdateLinkRequest struct {
	URL         string `json:"url,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublic    *bool  `json:"is_public,omitempty"`
	IsUnread    *bool  `json:"is_unread,omitempty"`
}

// TagResponse represents a tag in API responses.
type TagResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LinkCount int    `json:"link_count,omitempty"`
}

// CreateAPIKeyResponse is returned when creating an API key (includes full key).
type CreateAPIKeyResponse struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	KeyPrefix   string    `json:"key_prefix"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// APIKeyResponse represents an API key in list responses.
type APIKeyResponse struct {
	ID          string     `json:"id"`
	KeyPrefix   string     `json:"key_prefix"`
	Description string     `json:"description"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// SearchParams holds query parameters for the link search endpoint.
type SearchParams struct {
	Query    string
	IsUnread *bool
	IsPublic *bool
	GroupID  string
	Tag      string
	Limit    int
	Offset   int
}

// MessageResponse is a generic response containing just a message field.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorBody is the JSON body returned by the API on errors.
type ErrorBody struct {
	Error string `json:"error"`
}
