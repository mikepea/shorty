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

// GroupHandler handles SCIM Group operations
type GroupHandler struct {
	db      *gorm.DB
	baseURL string
}

// NewGroupHandler creates a new SCIM Group handler
func NewGroupHandler(db *gorm.DB, baseURL string) *GroupHandler {
	return &GroupHandler{db: db, baseURL: baseURL}
}

// groupToSCIM converts a models.Group to a SCIM Group
func (h *GroupHandler) groupToSCIM(group *models.Group, includeMembers bool) Group {
	created := group.CreatedAt
	updated := group.UpdatedAt

	scimGroup := Group{
		Schemas:     []string{SchemaGroup},
		ID:          strconv.FormatUint(uint64(group.ID), 10),
		ExternalID:  group.ExternalID,
		DisplayName: group.Name,
		Meta: Meta{
			ResourceType: "Group",
			Created:      &created,
			LastModified: &updated,
			Location:     fmt.Sprintf("%s/scim/v2/Groups/%d", h.baseURL, group.ID),
		},
	}

	if includeMembers {
		var memberships []models.GroupMembership
		h.db.Where("group_id = ?", group.ID).Preload("User").Find(&memberships)

		members := make([]GroupMember, len(memberships))
		for i, m := range memberships {
			members[i] = GroupMember{
				Value:   strconv.FormatUint(uint64(m.UserID), 10),
				Ref:     fmt.Sprintf("%s/scim/v2/Users/%d", h.baseURL, m.UserID),
				Display: m.User.Name,
			}
		}
		scimGroup.Members = members
	}

	return scimGroup
}

// ListGroups returns all groups (GET /scim/v2/Groups)
func (h *GroupHandler) ListGroups(c *gin.Context) {
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

	filter := c.Query("filter")

	var groups []models.Group
	query := h.db.Model(&models.Group{})

	if filter != "" {
		if strings.Contains(filter, "displayName eq") {
			parts := strings.Split(filter, "\"")
			if len(parts) >= 2 {
				name := parts[1]
				query = query.Where("name = ?", name)
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

	query.Offset(startIndex - 1).Limit(count).Find(&groups)

	resources := make([]Group, len(groups))
	for i, group := range groups {
		resources[i] = h.groupToSCIM(&group, false)
	}

	c.JSON(http.StatusOK, ListResponse{
		Schemas:      []string{SchemaListResponse},
		TotalResults: int(totalCount),
		StartIndex:   startIndex,
		ItemsPerPage: len(resources),
		Resources:    resources,
	})
}

// GetGroup returns a single group (GET /scim/v2/Groups/:id)
func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid group ID",
			Status:  "400",
		})
		return
	}

	var group models.Group
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Group not found",
			Status:  "404",
		})
		return
	}

	c.JSON(http.StatusOK, h.groupToSCIM(&group, true))
}

// CreateGroupRequest represents a SCIM group creation request
type CreateGroupRequest struct {
	Schemas     []string      `json:"schemas"`
	ExternalID  string        `json:"externalId"`
	DisplayName string        `json:"displayName"`
	Members     []GroupMember `json:"members"`
}

// CreateGroup creates a new group (POST /scim/v2/Groups)
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  err.Error(),
			Status:  "400",
		})
		return
	}

	if req.DisplayName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas:  []string{SchemaError},
			Detail:   "displayName is required",
			Status:   "400",
			ScimType: "invalidValue",
		})
		return
	}

	group := models.Group{
		ExternalID:  req.ExternalID,
		Name:        req.DisplayName,
		Description: "SCIM-provisioned group",
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		// Add members
		for _, member := range req.Members {
			userID, err := strconv.ParseUint(member.Value, 10, 32)
			if err != nil {
				continue
			}

			membership := models.GroupMembership{
				UserID:  uint(userID),
				GroupID: group.ID,
				Role:    models.GroupRoleMember,
			}
			if err := tx.Create(&membership).Error; err != nil {
				// Ignore duplicate membership errors
				continue
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to create group",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusCreated, h.groupToSCIM(&group, true))
}

// UpdateGroup replaces a group (PUT /scim/v2/Groups/:id)
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid group ID",
			Status:  "400",
		})
		return
	}

	var group models.Group
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Group not found",
			Status:  "404",
		})
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  err.Error(),
			Status:  "400",
		})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Update group
		group.Name = req.DisplayName
		if req.ExternalID != "" {
			group.ExternalID = req.ExternalID
		}
		group.UpdatedAt = time.Now()

		if err := tx.Save(&group).Error; err != nil {
			return err
		}

		// Replace members - remove all and add new
		tx.Where("group_id = ?", group.ID).Delete(&models.GroupMembership{})

		for _, member := range req.Members {
			userID, err := strconv.ParseUint(member.Value, 10, 32)
			if err != nil {
				continue
			}

			membership := models.GroupMembership{
				UserID:  uint(userID),
				GroupID: group.ID,
				Role:    models.GroupRoleMember,
			}
			tx.Create(&membership)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to update group",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusOK, h.groupToSCIM(&group, true))
}

// PatchGroup patches a group (PATCH /scim/v2/Groups/:id)
func (h *GroupHandler) PatchGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid group ID",
			Status:  "400",
		})
		return
	}

	var group models.Group
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Group not found",
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

	err = h.db.Transaction(func(tx *gorm.DB) error {
		for _, op := range patch.Operations {
			switch strings.ToLower(op.Op) {
			case "replace":
				h.applyReplaceOp(tx, &group, op)
			case "add":
				h.applyAddOp(tx, &group, op)
			case "remove":
				h.applyRemoveOp(tx, &group, op)
			}
		}

		group.UpdatedAt = time.Now()
		return tx.Save(&group).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to update group",
			Status:  "500",
		})
		return
	}

	c.JSON(http.StatusOK, h.groupToSCIM(&group, true))
}

func (h *GroupHandler) applyReplaceOp(tx *gorm.DB, group *models.Group, op PatchOperation) {
	path := strings.ToLower(op.Path)

	switch {
	case path == "displayname":
		if v, ok := op.Value.(string); ok {
			group.Name = v
		}
	case path == "externalid":
		if v, ok := op.Value.(string); ok {
			group.ExternalID = v
		}
	case path == "members":
		// Replace all members
		tx.Where("group_id = ?", group.ID).Delete(&models.GroupMembership{})
		if members, ok := op.Value.([]interface{}); ok {
			for _, m := range members {
				if memberMap, ok := m.(map[string]interface{}); ok {
					if value, ok := memberMap["value"].(string); ok {
						userID, err := strconv.ParseUint(value, 10, 32)
						if err == nil {
							tx.Create(&models.GroupMembership{
								UserID:  uint(userID),
								GroupID: group.ID,
								Role:    models.GroupRoleMember,
							})
						}
					}
				}
			}
		}
	case path == "":
		// No path means attributes in value
		if attrs, ok := op.Value.(map[string]interface{}); ok {
			if displayName, ok := attrs["displayName"].(string); ok {
				group.Name = displayName
			}
		}
	}
}

func (h *GroupHandler) applyAddOp(tx *gorm.DB, group *models.Group, op PatchOperation) {
	path := strings.ToLower(op.Path)

	if path == "members" || strings.HasPrefix(path, "members") {
		// Add members
		switch v := op.Value.(type) {
		case []interface{}:
			for _, m := range v {
				if memberMap, ok := m.(map[string]interface{}); ok {
					if value, ok := memberMap["value"].(string); ok {
						userID, err := strconv.ParseUint(value, 10, 32)
						if err == nil {
							tx.FirstOrCreate(&models.GroupMembership{
								UserID:  uint(userID),
								GroupID: group.ID,
							}, models.GroupMembership{
								UserID:  uint(userID),
								GroupID: group.ID,
								Role:    models.GroupRoleMember,
							})
						}
					}
				}
			}
		case map[string]interface{}:
			if value, ok := v["value"].(string); ok {
				userID, err := strconv.ParseUint(value, 10, 32)
				if err == nil {
					tx.FirstOrCreate(&models.GroupMembership{
						UserID:  uint(userID),
						GroupID: group.ID,
					}, models.GroupMembership{
						UserID:  uint(userID),
						GroupID: group.ID,
						Role:    models.GroupRoleMember,
					})
				}
			}
		}
	}
}

func (h *GroupHandler) applyRemoveOp(tx *gorm.DB, group *models.Group, op PatchOperation) {
	path := strings.ToLower(op.Path)

	// Handle paths like members[value eq "123"]
	if strings.HasPrefix(path, "members[") {
		// Extract user ID from filter
		// Format: members[value eq "123"]
		start := strings.Index(path, "\"")
		end := strings.LastIndex(path, "\"")
		if start != -1 && end != -1 && start < end {
			userIDStr := path[start+1 : end]
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err == nil {
				tx.Where("group_id = ? AND user_id = ?", group.ID, userID).Delete(&models.GroupMembership{})
			}
		}
	} else if path == "members" {
		// Remove all members
		tx.Where("group_id = ?", group.ID).Delete(&models.GroupMembership{})
	} else if path == "externalid" {
		group.ExternalID = ""
	}
}

// DeleteGroup deletes a group (DELETE /scim/v2/Groups/:id)
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Invalid group ID",
			Status:  "400",
		})
		return
	}

	var group models.Group
	if err := h.db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Group not found",
			Status:  "404",
		})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("group_id = ?", group.ID).Delete(&models.GroupMembership{})
		tx.Where("group_id = ?", group.ID).Delete(&models.Link{})
		return tx.Delete(&group).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Schemas: []string{SchemaError},
			Detail:  "Failed to delete group",
			Status:  "500",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers SCIM Group routes
func (h *GroupHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/Groups", h.ListGroups)
	rg.GET("/Groups/:id", h.GetGroup)
	rg.POST("/Groups", h.CreateGroup)
	rg.PUT("/Groups/:id", h.UpdateGroup)
	rg.PATCH("/Groups/:id", h.PatchGroup)
	rg.DELETE("/Groups/:id", h.DeleteGroup)
}
