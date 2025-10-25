package handlers

import (
	"net/http"
	"playspotter/internal/services"
	"playspotter/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user,omitempty"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.RespondError(c, http.StatusConflict, "email_exists", err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to register user")
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Data: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utils.RespondError(c, http.StatusUnauthorized, "invalid_credentials", err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to login")
		return
	}

	utils.RespondSuccess(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "invalid_refresh_token", err.Error())
		return
	}

	utils.RespondSuccess(c, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Revoke refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest false "Refresh token to revoke"
// @Success 200 {object} utils.SuccessResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		_ = h.authService.Logout(req.RefreshToken)
	}

	utils.RespondSuccess(c, gin.H{"message": "Logged out successfully"})
}

// BootstrapAdmin godoc
// @Summary Bootstrap admin user
// @Description Create initial admin user (only if no admin exists)
// @Tags internal
// @Accept json
// @Produce json
// @Security SetupToken
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /internal/bootstrap-admin [post]
func (h *AuthHandler) BootstrapAdmin(c *gin.Context, adminEmail, adminPassword, bootstrapToken string) {
	// Verify setup token
	token := c.GetHeader("X-Setup-Token")
	if token != bootstrapToken {
		utils.RespondError(c, http.StatusForbidden, "invalid_token", "Invalid setup token")
		return
	}

	err := h.authService.BootstrapAdmin(adminEmail, adminPassword)
	if err != nil {
		if err.Error() == "admin already exists" {
			utils.RespondError(c, http.StatusConflict, "admin_exists", err.Error())
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to create admin")
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Data: gin.H{"message": "Admin created successfully"},
	})
}
