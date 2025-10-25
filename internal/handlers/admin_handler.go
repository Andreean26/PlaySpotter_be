package handlers

import (
	"net/http"
	"playspotter/internal/services"
	"playspotter/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userService  *services.UserService
	eventService *services.EventService
}

func NewAdminHandler(userService *services.UserService, eventService *services.EventService) *AdminHandler {
	return &AdminHandler{
		userService:  userService,
		eventService: eventService,
	}
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=open full cancelled"`
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Description Get paginated list of all users
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	pagination := utils.NewPaginationParams(
		c.GetInt("page"),
		c.GetInt("limit"),
	)

	// Bind query params
	var query struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}
	if err := c.ShouldBindQuery(&query); err == nil {
		pagination = utils.NewPaginationParams(query.Page, query.Limit)
	}

	users, total, err := h.userService.ListUsers(pagination.GetOffset(), pagination.Limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to fetch users")
		return
	}

	meta := pagination.GetMeta(total)
	utils.RespondSuccessWithMeta(c, users, &meta)
}

// UpdateUserRole godoc
// @Summary Update user role (admin only)
// @Description Change a user's role between user and admin
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateRoleRequest true "New role"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/users/{id}/role [put]
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid user ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.userService.UpdateUserRole(id, req.Role); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	user, _ := h.userService.GetUser(id)
	utils.RespondSuccess(c, user)
}

// ListAllEvents godoc
// @Summary List all events (admin only)
// @Description Get paginated list of all events regardless of status
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /admin/events [get]
func (h *AdminHandler) ListAllEvents(c *gin.Context) {
	pagination := utils.NewPaginationParams(
		c.GetInt("page"),
		c.GetInt("limit"),
	)

	// Bind query params
	var query struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}
	if err := c.ShouldBindQuery(&query); err == nil {
		pagination = utils.NewPaginationParams(query.Page, query.Limit)
	}

	events, total, err := h.eventService.ListAllEvents(pagination.GetOffset(), pagination.Limit)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to fetch events")
		return
	}

	meta := pagination.GetMeta(total)
	utils.RespondSuccessWithMeta(c, events, &meta)
}

// UpdateEventStatus godoc
// @Summary Update event status (admin only)
// @Description Change an event's status
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param request body UpdateStatusRequest true "New status"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/events/{id}/status [put]
func (h *AdminHandler) UpdateEventStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.eventService.UpdateStatus(id, req.Status); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "update_failed", err.Error())
		return
	}

	event, _ := h.eventService.GetEvent(id)
	utils.RespondSuccess(c, event)
}
