package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatorID    uuid.UUID `gorm:"type:uuid;not null" json:"creator_id"`
	Title        string    `gorm:"type:varchar(120);not null" json:"title"`
	SportType    string    `gorm:"type:varchar(50);not null" json:"sport_type"`
	EventTime    time.Time `gorm:"type:timestamptz;not null" json:"event_time"`
	LocationName *string   `gorm:"type:varchar(160)" json:"location_name,omitempty"`
	Address      *string   `gorm:"type:text" json:"address,omitempty"`
	Latitude     float64   `gorm:"type:decimal(9,6);not null" json:"latitude"`
	Longitude    float64   `gorm:"type:decimal(9,6);not null" json:"longitude"`
	Capacity     int       `gorm:"type:int;not null;check:capacity >= 1" json:"capacity"`
	Description  *string   `gorm:"type:text" json:"description,omitempty"`
	Status       string    `gorm:"type:text;not null;default:'open';check:status IN ('open','full','cancelled')" json:"status"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relations (not stored in DB)
	Creator      *User  `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	Participants []User `gorm:"many2many:event_participants;" json:"participants,omitempty"`
}

func (Event) TableName() string {
	return "events"
}

// EventWithDistance extends Event with distance information
type EventWithDistance struct {
	Event
	DistanceKM *float64 `json:"distance_km,omitempty"`
}
