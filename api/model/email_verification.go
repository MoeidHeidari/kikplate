package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"          json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"      json:"user_id"`
	Token     string    `gorm:"size:255;not null;uniqueIndex" json:"-"`
	IsUsed    bool      `gorm:"default:false"                 json:"is_used"`
	ExpiresAt time.Time `gorm:"not null"                      json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime"                json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (EmailVerification) TableName() string {
	return "email_verification"
}
