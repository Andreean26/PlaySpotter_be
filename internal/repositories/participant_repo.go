package repositories

import (
	"playspotter/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) Create(participant *models.EventParticipant) error {
	return r.db.Create(participant).Error
}

func (r *ParticipantRepository) Delete(eventID, userID uuid.UUID) error {
	return r.db.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&models.EventParticipant{}).Error
}

func (r *ParticipantRepository) Exists(eventID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.EventParticipant{}).Where("event_id = ? AND user_id = ?", eventID, userID).Count(&count).Error
	return count > 0, err
}

func (r *ParticipantRepository) CountByEvent(eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.EventParticipant{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
