package handlers

import (
	"net/http"
	"playspotter/internal/middlewares"
	"playspotter/internal/services"
	"playspotter/internal/utils"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	userService *services.UserService
}

func NewMeHandler(userService *services.UserService) *MeHandler {
	return &MeHandler{
		userService: userService,
	}
}

type UpdateMeRequest struct {
	Name     string `json:"name"`
	Password string `json:"password" binding:"omitempty,min=8"`
}

// GetMe godoc
// @Summary Get current user info
// @Description Get authenticated user information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /me [get]
func (h *MeHandler) GetMe(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "user_not_found", "User not found")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// UpdateMe godoc
// @Summary Update current user
// @Description Update authenticated user information
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateMeRequest true "Update details"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /me [put]
func (h *MeHandler) UpdateMe(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.userService.UpdateUser(userID, req.Name, req.Password); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to update user")
		return
	}

	user, _ := h.userService.GetUser(userID)
	utils.RespondSuccess(c, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"updated_at": user.UpdatedAt,
	})
}
