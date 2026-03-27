package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex" json:"name"`
	Description string    `json:"description"`
	LogoURL     *string   `gorm:"size:2048" json:"logo_url,omitempty"`
	OwnerID     uuid.UUID `gorm:"index" json:"owner_id"`
	Owner       *Account  `json:"owner,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
