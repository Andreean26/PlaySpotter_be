package models

import (
	"time"

	"github.com/google/uuid"
)

type EventSwipe struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Action    string    `gorm:"type:text;not null;check:action IN ('like','skip')" json:"action"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	// Relations
	Event *Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (EventSwipe) TableName() string {
	return "event_swipes"
}
