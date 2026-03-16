package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"                          json:"id"`
	UserID         *uuid.UUID `gorm:"type:uuid;index"                               json:"user_id,omitempty"`
	Provider       string     `gorm:"size:64;not null"                              json:"provider"`
	ProviderUserID string     `gorm:"size:255;not null"                             json:"provider_user_id"`
	DisplayName    *string    `gorm:"size:255"                                      json:"display_name,omitempty"`
	AvatarURL      *string    `gorm:"size:2048"                                     json:"avatar_url,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"                                json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Account) TableName() string {
	return "account"
}
