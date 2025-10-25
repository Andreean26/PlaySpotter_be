package services

import (
	"playspotter/internal/models"
	"playspotter/internal/repositories"

	"github.com/google/uuid"
)

type SwipeService struct {
	swipeRepo *repositories.SwipeRepository
}

func NewSwipeService(swipeRepo *repositories.SwipeRepository) *SwipeService {
	return &SwipeService{
		swipeRepo: swipeRepo,
	}
}

func (s *SwipeService) RecordSwipe(eventID, userID uuid.UUID, action string) error {
	swipe := &models.EventSwipe{
		EventID: eventID,
		UserID:  userID,
		Action:  action,
	}
	return s.swipeRepo.Upsert(swipe)
}
