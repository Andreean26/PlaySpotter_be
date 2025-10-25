package handlers

import (
	"net/http"
	"playspotter/internal/middlewares"
	"playspotter/internal/models"
	"playspotter/internal/repositories"
	"playspotter/internal/services"
	"playspotter/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *services.EventService
	swipeService *services.SwipeService
}

func NewEventHandler(eventService *services.EventService, swipeService *services.SwipeService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		swipeService: swipeService,
	}
}

type CreateEventRequest struct {
	Title        string  `json:"title" binding:"required,max=120"`
	SportType    string  `json:"sport_type" binding:"required,max=50"`
	EventTime    string  `json:"event_time" binding:"required"`
	LocationName *string `json:"location_name" binding:"omitempty,max=160"`
	Address      *string `json:"address"`
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Capacity     int     `json:"capacity" binding:"required,min=1"`
	Description  *string `json:"description"`
}

type UpdateEventRequest struct {
	Title        string   `json:"title" binding:"omitempty,max=120"`
	SportType    string   `json:"sport_type" binding:"omitempty,max=50"`
	EventTime    string   `json:"event_time"`
	LocationName *string  `json:"location_name" binding:"omitempty,max=160"`
	Address      *string  `json:"address"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	Capacity     *int     `json:"capacity" binding:"omitempty,min=1"`
	Description  *string  `json:"description"`
}

type SwipeRequest struct {
	Action string `json:"action" binding:"required,oneof=like skip"`
}

type EventFeedQuery struct {
	Lat         *float64 `form:"lat" binding:"omitempty,min=-90,max=90"`
	Lng         *float64 `form:"lng" binding:"omitempty,min=-180,max=180"`
	MaxDistance *float64 `form:"max_distance_km" binding:"omitempty,min=0"`
	SportType   string   `form:"sport_type"`
	DateFrom    string   `form:"date_from"`
	DateTo      string   `form:"date_to"`
	Page        int      `form:"page" binding:"omitempty,min=1"`
	Limit       int      `form:"limit" binding:"omitempty,min=1,max=100"`
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new sports event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEventRequest true "Event details"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Parse event time
	eventTime, err := time.Parse(time.RFC3339, req.EventTime)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_time_format", "Event time must be in RFC3339 format")
		return
	}

	event := &models.Event{
		CreatorID:    userID,
		Title:        req.Title,
		SportType:    req.SportType,
		EventTime:    eventTime.UTC(),
		LocationName: req.LocationName,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Capacity:     req.Capacity,
		Description:  req.Description,
	}

	if err := h.eventService.CreateEvent(event); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse{
		Data: event,
	})
}

// GetEvent godoc
// @Summary Get event by ID
// @Description Get detailed information about a specific event
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	event, err := h.eventService.GetEvent(id)
	if err != nil {
		utils.RespondError(c, http.StatusNotFound, "event_not_found", "Event not found")
		return
	}

	utils.RespondSuccess(c, event)
}

// ListEvents godoc
// @Summary List events
// @Description Get a list of events with optional filters
// @Tags events
// @Accept json
// @Produce json
// @Param lat query number false "Latitude"
// @Param lng query number false "Longitude"
// @Param max_distance_km query number false "Maximum distance in km"
// @Param sport_type query string false "Sport type"
// @Param date_from query string false "Date from (RFC3339)"
// @Param date_to query string false "Date to (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	var query EventFeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Setup pagination
	pagination := utils.NewPaginationParams(query.Page, query.Limit)

	// Parse dates if provided
	var dateFrom, dateTo *time.Time
	if query.DateFrom != "" {
		t, err := time.Parse(time.RFC3339, query.DateFrom)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "invalid_date_format", "date_from must be in RFC3339 format")
			return
		}
		dateFrom = &t
	}
	if query.DateTo != "" {
		t, err := time.Parse(time.RFC3339, query.DateTo)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "invalid_date_format", "date_to must be in RFC3339 format")
			return
		}
		dateTo = &t
	}

	filter := repositories.EventFilter{
		Lat:         query.Lat,
		Lng:         query.Lng,
		MaxDistance: query.MaxDistance,
		SportType:   query.SportType,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Status:      "open",
		Offset:      pagination.GetOffset(),
		Limit:       pagination.Limit,
	}

	events, total, err := h.eventService.ListEvents(filter)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "internal_error", "Failed to fetch events")
		return
	}

	meta := pagination.GetMeta(total)
	utils.RespondSuccessWithMeta(c, events, &meta)
}

// UpdateEvent godoc
// @Summary Update event
// @Description Update event details (creator or admin only)
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param request body UpdateEventRequest true "Update details"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	role, _ := middlewares.GetUserRole(c)
	isAdmin := role == "admin"

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	updates := &models.Event{
		Title:        req.Title,
		SportType:    req.SportType,
		LocationName: req.LocationName,
		Address:      req.Address,
		Description:  req.Description,
	}

	if req.EventTime != "" {
		eventTime, err := time.Parse(time.RFC3339, req.EventTime)
		if err != nil {
			utils.RespondError(c, http.StatusBadRequest, "invalid_time_format", "Event time must be in RFC3339 format")
			return
		}
		updates.EventTime = eventTime.UTC()
	}

	if req.Latitude != nil {
		updates.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		updates.Longitude = *req.Longitude
	}
	if req.Capacity != nil {
		updates.Capacity = *req.Capacity
	}

	if err := h.eventService.UpdateEvent(id, updates, userID, isAdmin); err != nil {
		if err.Error() == "only event creator or admin can update this event" {
			utils.RespondError(c, http.StatusForbidden, "forbidden", err.Error())
			return
		}
		utils.RespondError(c, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	event, _ := h.eventService.GetEvent(id)
	utils.RespondSuccess(c, event)
}

// DeleteEvent godoc
// @Summary Delete/Cancel event
// @Description Cancel an event (creator or admin only)
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	role, _ := middlewares.GetUserRole(c)
	isAdmin := role == "admin"

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	if err := h.eventService.DeleteEvent(id, userID, isAdmin); err != nil {
		if err.Error() == "only event creator or admin can delete this event" {
			utils.RespondError(c, http.StatusForbidden, "forbidden", err.Error())
			return
		}
		utils.RespondError(c, http.StatusBadRequest, "delete_failed", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "Event cancelled successfully"})
}

// JoinEvent godoc
// @Summary Join an event
// @Description Join as a participant in an event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /events/{id}/join [post]
func (h *EventHandler) JoinEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	if err := h.eventService.JoinEvent(id, userID); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "join_failed", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "Joined event successfully"})
}

// LeaveEvent godoc
// @Summary Leave an event
// @Description Remove yourself as a participant from an event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /events/{id}/leave [post]
func (h *EventHandler) LeaveEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	if err := h.eventService.LeaveEvent(id, userID); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "leave_failed", err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "Left event successfully"})
}

// SwipeEvent godoc
// @Summary Swipe on an event
// @Description Record a like or skip action for an event
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Event ID"
// @Param request body SwipeRequest true "Swipe action"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /events/{id}/swipe [post]
func (h *EventHandler) SwipeEvent(c *gin.Context) {
	userID, ok := middlewares.GetUserID(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized", "User ID not found")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid_id", "Invalid event ID")
		return
	}

	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.swipeService.RecordSwipe(id, userID, req.Action); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "swipe_failed", "Failed to record swipe")
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "Swipe recorded successfully"})
}
