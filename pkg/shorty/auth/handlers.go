package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikepea/shorty/pkg/shorty/models"
	"gorm.io/gorm"
)

// Handler handles authentication requests
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a new auth handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	SystemRole  string `json:"system_role"`
	HasPassword bool   `json:"has_password"`
}

// ChangePasswordRequest represents the password change request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account and receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 409 {object} map[string]string "Email already registered"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user and personal group in a transaction
	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		SystemRole:   models.SystemRoleUser,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create personal group
		personalGroup := models.Group{
			Name:        req.Name + "'s Links",
			Description: "Personal links for " + req.Name,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Email, string(user.SystemRole))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			SystemRole:  string(user.SystemRole),
			HasPassword: user.PasswordHash != "",
		},
	})
}

// Login handles user login
// @Summary Login
// @Description Authenticate with email and password to receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if !CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate token
	token, err := GenerateToken(user.ID, user.Email, string(user.SystemRole))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			SystemRole:  string(user.SystemRole),
			HasPassword: user.PasswordHash != "",
		},
	})
}

// Me returns the current authenticated user
// @Summary Get current user
// @Description Get the authenticated user's profile
// @Tags auth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string "Authentication required"
// @Security BearerAuth
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		SystemRole:  string(user.SystemRole),
		HasPassword: user.PasswordHash != "",
	})
}

// Logout handles user logout (client-side token invalidation)
// @Summary Logout
// @Description Logout the current user (client-side token invalidation)
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ChangePassword handles password change for authenticated users
// @Summary Change password
// @Description Change the password for the authenticated user (requires existing password)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]string "Validation error or SSO-only user"
// @Failure 401 {object} map[string]string "Authentication required or wrong current password"
// @Security BearerAuth
// @Router /auth/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user has a password (not SSO-only)
	if user.PasswordHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change password for SSO-only accounts"})
		return
	}

	// Verify current password
	if !CheckPassword(req.CurrentPassword, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Update password
	if err := h.db.Model(&user).Update("password_hash", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// RegisterRoutes registers auth routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
	rg.POST("/logout", h.Logout)
	rg.GET("/me", AuthMiddleware(), h.Me)
	rg.PUT("/password", AuthMiddleware(), h.ChangePassword)
}
