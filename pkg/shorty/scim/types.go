package scim

import "time"

// SCIM 2.0 Schema URIs
const (
	SchemaUser              = "urn:ietf:params:scim:schemas:core:2.0:User"
	SchemaGroup             = "urn:ietf:params:scim:schemas:core:2.0:Group"
	SchemaListResponse      = "urn:ietf:params:scim:api:messages:2.0:ListResponse"
	SchemaError             = "urn:ietf:params:scim:api:messages:2.0:Error"
	SchemaPatchOp           = "urn:ietf:params:scim:api:messages:2.0:PatchOp"
	SchemaServiceProvider   = "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"
	SchemaResourceType      = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
	SchemaSchema            = "urn:ietf:params:scim:schemas:core:2.0:Schema"
)

// Meta contains resource metadata
type Meta struct {
	ResourceType string     `json:"resourceType"`
	Created      *time.Time `json:"created,omitempty"`
	LastModified *time.Time `json:"lastModified,omitempty"`
	Location     string     `json:"location,omitempty"`
	Version      string     `json:"version,omitempty"`
}

// Name represents a user's name
type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

// Email represents a user's email
type Email struct {
	Value   string `json:"value"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

// User represents a SCIM User resource
type User struct {
	Schemas     []string `json:"schemas"`
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId,omitempty"`
	Meta        Meta     `json:"meta"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Emails      []Email  `json:"emails,omitempty"`
	Active      bool     `json:"active"`
}

// GroupMember represents a member in a SCIM Group
type GroupMember struct {
	Value   string `json:"value"`
	Ref     string `json:"$ref,omitempty"`
	Display string `json:"display,omitempty"`
}

// Group represents a SCIM Group resource
type Group struct {
	Schemas    []string      `json:"schemas"`
	ID         string        `json:"id"`
	ExternalID string        `json:"externalId,omitempty"`
	Meta       Meta          `json:"meta"`
	DisplayName string       `json:"displayName"`
	Members    []GroupMember `json:"members,omitempty"`
}

// ListResponse represents a SCIM list response
type ListResponse struct {
	Schemas      []string    `json:"schemas"`
	TotalResults int         `json:"totalResults"`
	StartIndex   int         `json:"startIndex"`
	ItemsPerPage int         `json:"itemsPerPage"`
	Resources    interface{} `json:"Resources"`
}

// ErrorResponse represents a SCIM error response
type ErrorResponse struct {
	Schemas  []string `json:"schemas"`
	Detail   string   `json:"detail"`
	Status   string   `json:"status"`
	ScimType string   `json:"scimType,omitempty"`
}

// PatchOp represents a SCIM PATCH operation
type PatchOp struct {
	Schemas    []string        `json:"schemas"`
	Operations []PatchOperation `json:"Operations"`
}

// PatchOperation represents a single operation in a PATCH request
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// ServiceProviderConfig represents SCIM service provider configuration
type ServiceProviderConfig struct {
	Schemas               []string                  `json:"schemas"`
	DocumentationURI      string                    `json:"documentationUri,omitempty"`
	Patch                 SupportedConfig           `json:"patch"`
	Bulk                  BulkConfig                `json:"bulk"`
	Filter                FilterConfig              `json:"filter"`
	ChangePassword        SupportedConfig           `json:"changePassword"`
	Sort                  SupportedConfig           `json:"sort"`
	Etag                  SupportedConfig           `json:"etag"`
	AuthenticationSchemes []AuthenticationScheme    `json:"authenticationSchemes"`
	Meta                  Meta                      `json:"meta"`
}

// SupportedConfig indicates if a feature is supported
type SupportedConfig struct {
	Supported bool `json:"supported"`
}

// BulkConfig represents bulk operation configuration
type BulkConfig struct {
	Supported      bool `json:"supported"`
	MaxOperations  int  `json:"maxOperations"`
	MaxPayloadSize int  `json:"maxPayloadSize"`
}

// FilterConfig represents filter configuration
type FilterConfig struct {
	Supported  bool `json:"supported"`
	MaxResults int  `json:"maxResults"`
}

// AuthenticationScheme represents an authentication scheme
type AuthenticationScheme struct {
	Type             string `json:"type"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	SpecURI          string `json:"specUri,omitempty"`
	DocumentationURI string `json:"documentationUri,omitempty"`
	Primary          bool   `json:"primary,omitempty"`
}

// ResourceType represents a SCIM resource type
type ResourceType struct {
	Schemas         []string          `json:"schemas"`
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Endpoint        string            `json:"endpoint"`
	Description     string            `json:"description"`
	Schema          string            `json:"schema"`
	SchemaExtensions []SchemaExtension `json:"schemaExtensions,omitempty"`
	Meta            Meta              `json:"meta"`
}

// SchemaExtension represents a schema extension
type SchemaExtension struct {
	Schema   string `json:"schema"`
	Required bool   `json:"required"`
}
