package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash string    `gorm:"type:text;not null" json:"-"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	Revoked   bool      `gorm:"type:boolean;not null;default:false" json:"revoked"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
