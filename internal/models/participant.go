package models

import (
	"time"

	"github.com/google/uuid"
)

type EventParticipant struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EventID  uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	JoinedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"joined_at"`

	// Relations
	Event *Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (EventParticipant) TableName() string {
	return "event_participants"
}
