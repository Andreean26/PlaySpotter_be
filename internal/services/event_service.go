package services

import (
	"errors"
	"playspotter/internal/models"
	"playspotter/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type EventService struct {
	eventRepo       *repositories.EventRepository
	participantRepo *repositories.ParticipantRepository
}

func NewEventService(eventRepo *repositories.EventRepository, participantRepo *repositories.ParticipantRepository) *EventService {
	return &EventService{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
	}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	// Validate event time is in the future
	if event.EventTime.Before(time.Now().UTC()) {
		return errors.New("event time must be in the future")
	}

	// Validate coordinates
	if event.Latitude < -90 || event.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if event.Longitude < -180 || event.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	// Validate capacity
	if event.Capacity < 1 {
		return errors.New("capacity must be at least 1")
	}

	event.Status = "open"
	return s.eventRepo.Create(event)
}

func (s *EventService) GetEvent(id uuid.UUID) (*models.Event, error) {
	return s.eventRepo.FindByID(id)
}

func (s *EventService) UpdateEvent(id uuid.UUID, updates *models.Event, userID uuid.UUID, isAdmin bool) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check permissions
	if !isAdmin && event.CreatorID != userID {
		return errors.New("only event creator or admin can update this event")
	}

	// Validate event time if being updated
	if !updates.EventTime.IsZero() && updates.EventTime.Before(time.Now().UTC()) {
		return errors.New("event time must be in the future")
	}

	// Update fields
	if updates.Title != "" {
		event.Title = updates.Title
	}
	if updates.SportType != "" {
		event.SportType = updates.SportType
	}
	if !updates.EventTime.IsZero() {
		event.EventTime = updates.EventTime
	}
	if updates.LocationName != nil {
		event.LocationName = updates.LocationName
	}
	if updates.Address != nil {
		event.Address = updates.Address
	}
	if updates.Latitude != 0 {
		if updates.Latitude < -90 || updates.Latitude > 90 {
			return errors.New("latitude must be between -90 and 90")
		}
		event.Latitude = updates.Latitude
	}
	if updates.Longitude != 0 {
		if updates.Longitude < -180 || updates.Longitude > 180 {
			return errors.New("longitude must be between -180 and 180")
		}
		event.Longitude = updates.Longitude
	}
	if updates.Capacity > 0 {
		event.Capacity = updates.Capacity
	}
	if updates.Description != nil {
		event.Description = updates.Description
	}

	return s.eventRepo.Update(event)
}

func (s *EventService) DeleteEvent(id uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check permissions
	if !isAdmin && event.CreatorID != userID {
		return errors.New("only event creator or admin can delete this event")
	}

	// Set status to cancelled instead of deleting
	event.Status = "cancelled"
	return s.eventRepo.Update(event)
}

func (s *EventService) JoinEvent(eventID, userID uuid.UUID) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return err
	}

	// Check if event is cancelled
	if event.Status == "cancelled" {
		return errors.New("cannot join cancelled event")
	}

	// Check if event is full
	if event.Status == "full" {
		return errors.New("event is full")
	}

	// Check if event time has passed
	if event.EventTime.Before(time.Now().UTC()) {
		return errors.New("cannot join past event")
	}

	// Check if already joined
	exists, err := s.participantRepo.Exists(eventID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already joined this event")
	}

	// Add participant
	participant := &models.EventParticipant{
		EventID: eventID,
		UserID:  userID,
	}
	if err := s.participantRepo.Create(participant); err != nil {
		return err
	}

	// Check if event is now full
	count, err := s.participantRepo.CountByEvent(eventID)
	if err != nil {
		return err
	}

	if int(count) >= event.Capacity {
		event.Status = "full"
		return s.eventRepo.Update(event)
	}

	return nil
}

func (s *EventService) LeaveEvent(eventID, userID uuid.UUID) error {
	// Check if user is participant
	exists, err := s.participantRepo.Exists(eventID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("you are not a participant of this event")
	}

	// Remove participant
	if err := s.participantRepo.Delete(eventID, userID); err != nil {
		return err
	}

	// Update event status if it was full
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return err
	}

	if event.Status == "full" {
		event.Status = "open"
		return s.eventRepo.Update(event)
	}

	return nil
}

func (s *EventService) ListEvents(filter repositories.EventFilter) ([]map[string]interface{}, int64, error) {
	return s.eventRepo.List(filter)
}

func (s *EventService) ListAllEvents(offset, limit int) ([]models.Event, int64, error) {
	return s.eventRepo.ListAll(offset, limit)
}

func (s *EventService) UpdateStatus(id uuid.UUID, status string) error {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return err
	}

	if status != "open" && status != "full" && status != "cancelled" {
		return errors.New("invalid status")
	}

	event.Status = status
	return s.eventRepo.Update(event)
}
