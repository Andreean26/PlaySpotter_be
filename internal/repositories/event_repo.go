package repositories

import (
	"fmt"
	"playspotter/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) FindByID(id uuid.UUID) (*models.Event, error) {
	var event models.Event
	err := r.db.Preload("Creator").Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *EventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Event{}, id).Error
}

type EventFilter struct {
	Lat         *float64
	Lng         *float64
	MaxDistance *float64
	SportType   string
	DateFrom    *time.Time
	DateTo      *time.Time
	Status      string
	Offset      int
	Limit       int
}

func (r *EventRepository) List(filter EventFilter) ([]map[string]interface{}, int64, error) {
	var total int64
	var results []map[string]interface{}

	// Base query for counting
	countQuery := r.db.Model(&models.Event{})

	// Base query for selecting
	query := r.db.Table("events e").
		Select("e.*, u.name as creator_name, u.email as creator_email")

	// Add join for creator
	query = query.Joins("LEFT JOIN users u ON e.creator_id = u.id")
	countQuery = countQuery.Joins("LEFT JOIN users u ON e.creator_id = u.id")

	// Add distance calculation if lat/lng provided
	if filter.Lat != nil && filter.Lng != nil {
		distanceFormula := fmt.Sprintf(`
			2 * 6371 * ASIN(
				SQRT(
					POWER(SIN(RADIANS(%f - e.latitude)/2), 2) +
					COS(RADIANS(e.latitude)) * COS(RADIANS(%f)) *
					POWER(SIN(RADIANS(%f - e.longitude)/2), 2)
				)
			)
		`, *filter.Lat, *filter.Lat, *filter.Lng)

		query = query.Select("e.*, u.name as creator_name, u.email as creator_email, " + distanceFormula + " as distance_km")

		// Add max distance filter if provided
		if filter.MaxDistance != nil {
			query = query.Where(distanceFormula+" <= ?", *filter.MaxDistance)
			countQuery = countQuery.Where(distanceFormula+" <= ?", *filter.MaxDistance)
		}
	}

	// Apply filters
	if filter.Status != "" {
		query = query.Where("e.status = ?", filter.Status)
		countQuery = countQuery.Where("status = ?", filter.Status)
	} else {
		// Default: only open events
		query = query.Where("e.status = ?", "open")
		countQuery = countQuery.Where("status = ?", "open")
	}

	// Filter future events
	query = query.Where("e.event_time > ?", time.Now().UTC())
	countQuery = countQuery.Where("event_time > ?", time.Now().UTC())

	if filter.SportType != "" {
		query = query.Where("e.sport_type = ?", filter.SportType)
		countQuery = countQuery.Where("sport_type = ?", filter.SportType)
	}

	if filter.DateFrom != nil {
		query = query.Where("e.event_time >= ?", filter.DateFrom)
		countQuery = countQuery.Where("event_time >= ?", filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("e.event_time <= ?", filter.DateTo)
		countQuery = countQuery.Where("event_time <= ?", filter.DateTo)
	}

	// Count total
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Order by
	if filter.Lat != nil && filter.Lng != nil {
		query = query.Order("distance_km ASC, e.event_time ASC")
	} else {
		query = query.Order("e.event_time ASC")
	}

	// Apply pagination
	query = query.Offset(filter.Offset).Limit(filter.Limit)

	// Execute query
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *EventRepository) ListAll(offset, limit int) ([]models.Event, int64, error) {
	var events []models.Event
	var total int64

	if err := r.db.Model(&models.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Creator").Offset(offset).Limit(limit).Order("event_time DESC").Find(&events).Error
	return events, total, err
}

func (r *EventRepository) GetParticipantCount(eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.EventParticipant{}).Where("event_id = ?", eventID).Count(&count).Error
	return count, err
}
