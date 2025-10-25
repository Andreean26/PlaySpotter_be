package repositories

import (
	"playspotter/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwipeRepository struct {
	db *gorm.DB
}

func NewSwipeRepository(db *gorm.DB) *SwipeRepository {
	return &SwipeRepository{db: db}
}

func (r *SwipeRepository) Upsert(swipe *models.EventSwipe) error {
	// Try to find existing swipe
	var existing models.EventSwipe
	err := r.db.Where("event_id = ? AND user_id = ?", swipe.EventID, swipe.UserID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		return r.db.Create(swipe).Error
	} else if err != nil {
		return err
	}

	// Update existing
	existing.Action = swipe.Action
	return r.db.Save(&existing).Error
}

func (r *SwipeRepository) Find(eventID, userID uuid.UUID) (*models.EventSwipe, error) {
	var swipe models.EventSwipe
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&swipe).Error
	if err != nil {
		return nil, err
	}
	return &swipe, nil
}
