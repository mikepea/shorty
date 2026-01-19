package scim

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// UserHandler handles SCIM User operations
type UserHandler struct {
	db      *gorm.DB
	baseURL string
}

// NewUserHandler creates a new SCIM User handler
func NewUserHandler(db *gorm.DB, baseURL string) *UserHandler {
	return &UserHandler{db: db, baseURL: baseURL}
}

// userToSCIM converts a models.User to a SCIM User
func (h *UserHandler) userToSCIM(user *models.User) User {
	emails := []Email{}
	if user.Email != "" {
		emails = append(emails, Email{
			Value:   user.Email,
			Type:    "work",
			Primary: true,
		})
	}

	displayName := user.Name
	if displayName == "" {
		displayName = user.Email
	}

	created := user.CreatedAt
	updated := user.UpdatedAt

	return User{
		Schemas:    []string{SchemaUser},
		ID:         strconv.FormatUint(uint64(user.ID), 10),
		ExternalID: user.ExternalID,
		Meta: Meta{
			ResourceType: "User",
			Created:      &created,
			LastModified: &updated,
			Location:     fmt.Sprintf("%s/scim/v2/Users/%d", h.baseURL, user.ID),
		},
		UserName:    user.Email,
		DisplayName: displayName,
		Name: Name{
			Formatted:  user.Name,
			GivenName:  user.GivenName,
			FamilyName: user.FamilyName,
		},
		Emails: emails,
		Active: user.Active,
	}
}

// ListUsers returns all users (GET /scim/v2/Users)
func (h *UserHandler) ListUsers(c *gin.Context) {
	startIndex, _ := strconv.Atoi(c.DefaultQuery("startIndex", "1"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "100"))

	if startIndex < 1 {
		startIndex = 1
	}
	if count < 1 {
		count = 100
	}
	if count > 1000 {
		count = 1000
	}

	// Parse filter (basic support for userName eq "email")
	filter := c.Query("filter")

	var users []models.User
	query := h.db.Model(&models.User{})

	if filter != "" {
		// Basic filter parsing for userName eq "value"
		if strings.Contains(filter, "userName eq") {
			parts := strings.Split(filter, "\"")
			if len(parts) >= 2 {
				email := parts[1]
				query = query.Where("email = ?", email)
			}
		} else if strings.Contains(filter, "externalId eq") {
			parts := strings.Split(filter, "\"")
			if len(parts) >= 2 {
				externalID := parts[1]
				query = query.Where("external_id = ?", externalID)
			}
		}
	}

	var totalCount int64
	query.Count(&totalCount)

	query.Offset(startIndex - 1).Limit(count).Find(&users)

	resources := make([]User, len(users))
	for i, user := range users {
		resources[i] = h.userToSCIM(&user)
	}

	c.JSON(http.StatusOK, ListResponse{
		Schemas:      []string{SchemaListResponse},
		TotalResults: int(totalCount),
		StartIndex:   startIndex,
		ItemsPerPage: len(resources),
		Resources:    resources,
	})
}

// GetUser returns a single user (GET /scim/v2/Users/:id)
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid user ID",
			Status:  "400",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "User not found",
			Status:  "404",
		})
		return
	}

	c.JSON(http.StatusOK, h.userToSCIM(&user))
}

// CreateUserRequest represents a SCIM user creation request
type CreateUserRequest struct {
	Schemas     []string `json:"schemas"`
	ExternalID  string   `json:"externalId"`
	UserName    string   `json:"userName"`
	Name        Name     `json:"name"`
	DisplayName string   `json:"displayName"`
	Emails      []Email  `json:"emails"`
	Active      *bool    `json:"active"`
}

// CreateUser creates a new user (POST /scim/v2/Users)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  err.Error(),
			Status:  "400",
		})
		return
	}

	// Determine email
	email := req.UserName
	if email == "" && len(req.Emails) > 0 {
		for _, e := range req.Emails {
			if e.Primary || email == "" {
				email = e.Value
			}
		}
	}

	if email == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas:  []string{SchemaError},
			Detail:   "userName or email is required",
			Status:   "400",
			ScimType: "invalidValue",
		})
		return
	}

	// Check for existing user
	var existingUser models.User
	if err := h.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, ErrorResponse{
			Schemas:  []string{SchemaError},
			Detail:   "User with this email already exists",
			Status:   "409",
			ScimType: "uniqueness",
		})
		return
	}

	// Determine name
	name := req.DisplayName
	if name == "" {
		name = req.Name.Formatted
	}
	if name == "" && (req.Name.GivenName != "" || req.Name.FamilyName != "") {
		name = strings.TrimSpace(req.Name.GivenName + " " + req.Name.FamilyName)
	}
	if name == "" {
		name = strings.Split(email, "@")[0]
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	user := models.User{
		ExternalID: req.ExternalID,
		Email:      email,
		Name:       name,
		GivenName:  req.Name.GivenName,
		FamilyName: req.Name.FamilyName,
		Active:     active,
		SystemRole: models.SystemRoleUser,
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create personal group
		personalGroup := models.Group{
			Name:        user.Name + "'s Links",
			Description: "Personal links for " + user.Name,
		}
		if err := tx.Create(&personalGroup).Error; err != nil {
			return err
		}

		// Add user as admin of personal group
		membership := models.GroupMembership{
			UserID:  user.ID,
			GroupID: personalGroup.ID,
			Role:    models.GroupRoleAdmin,
		}
		return tx.Create(&membership).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to create user",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusCreated, h.userToSCIM(&user))
}

// UpdateUser replaces a user (PUT /scim/v2/Users/:id)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid user ID",
			Status:  "400",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "User not found",
			Status:  "404",
		})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  err.Error(),
			Status:  "400",
		})
		return
	}

	// Update fields
	if req.ExternalID != "" {
		user.ExternalID = req.ExternalID
	}

	email := req.UserName
	if email == "" && len(req.Emails) > 0 {
		for _, e := range req.Emails {
			if e.Primary || email == "" {
				email = e.Value
			}
		}
	}
	if email != "" {
		user.Email = email
	}

	if req.DisplayName != "" {
		user.Name = req.DisplayName
	} else if req.Name.Formatted != "" {
		user.Name = req.Name.Formatted
	} else if req.Name.GivenName != "" || req.Name.FamilyName != "" {
		user.Name = strings.TrimSpace(req.Name.GivenName + " " + req.Name.FamilyName)
	}

	user.GivenName = req.Name.GivenName
	user.FamilyName = req.Name.FamilyName

	if req.Active != nil {
		user.Active = *req.Active
	}

	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to update user",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusOK, h.userToSCIM(&user))
}

// PatchUser patches a user (PATCH /scim/v2/Users/:id)
func (h *UserHandler) PatchUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid user ID",
			Status:  "400",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "User not found",
			Status:  "404",
		})
		return
	}

	var patch PatchOp
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  err.Error(),
			Status:  "400",
		})
		return
	}

	for _, op := range patch.Operations {
		switch strings.ToLower(op.Op) {
		case "replace":
			h.applyReplaceOp(&user, op)
		case "add":
			h.applyReplaceOp(&user, op) // For users, add is similar to replace
		case "remove":
			h.applyRemoveOp(&user, op)
		}
	}

	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to update user",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusOK, h.userToSCIM(&user))
}

func (h *UserHandler) applyReplaceOp(user *models.User, op PatchOperation) {
	path := strings.ToLower(op.Path)

	switch path {
	case "active":
		if v, ok := op.Value.(bool); ok {
			user.Active = v
		}
	case "username":
		if v, ok := op.Value.(string); ok {
			user.Email = v
		}
	case "displayname":
		if v, ok := op.Value.(string); ok {
			user.Name = v
		}
	case "externalid":
		if v, ok := op.Value.(string); ok {
			user.ExternalID = v
		}
	case "name.givenname":
		if v, ok := op.Value.(string); ok {
			user.GivenName = v
		}
	case "name.familyname":
		if v, ok := op.Value.(string); ok {
			user.FamilyName = v
		}
	case "name":
		if nameMap, ok := op.Value.(map[string]interface{}); ok {
			if gn, ok := nameMap["givenName"].(string); ok {
				user.GivenName = gn
			}
			if fn, ok := nameMap["familyName"].(string); ok {
				user.FamilyName = fn
			}
			if f, ok := nameMap["formatted"].(string); ok {
				user.Name = f
			}
		}
	case "":
		// No path means the value contains the attributes to update
		if attrs, ok := op.Value.(map[string]interface{}); ok {
			if active, ok := attrs["active"].(bool); ok {
				user.Active = active
			}
			if userName, ok := attrs["userName"].(string); ok {
				user.Email = userName
			}
			if displayName, ok := attrs["displayName"].(string); ok {
				user.Name = displayName
			}
		}
	}
}

func (h *UserHandler) applyRemoveOp(user *models.User, op PatchOperation) {
	path := strings.ToLower(op.Path)

	switch path {
	case "externalid":
		user.ExternalID = ""
	case "name.givenname":
		user.GivenName = ""
	case "name.familyname":
		user.FamilyName = ""
	}
}

// DeleteUser deletes a user (DELETE /scim/v2/Users/:id)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid user ID",
			Status:  "400",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "User not found",
			Status:  "404",
		})
		return
	}

	// Delete user and related data
	err = h.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", user.ID).Delete(&models.APIKey{})
		tx.Where("user_id = ?", user.ID).Delete(&models.GroupMembership{})
		tx.Where("user_id = ?", user.ID).Delete(&models.OIDCIdentity{})
		tx.Where("created_by_id = ?", user.ID).Delete(&models.Link{})
		return tx.Delete(&user).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to delete user",
			Status:  "500",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers SCIM User routes
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/Users", h.ListUsers)
	rg.GET("/Users/:id", h.GetUser)
	rg.POST("/Users", h.CreateUser)
	rg.PUT("/Users/:id", h.UpdateUser)
	rg.PATCH("/Users/:id", h.PatchUser)
	rg.DELETE("/Users/:id", h.DeleteUser)
}
