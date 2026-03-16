package model

import (
	"time"

	"github.com/google/uuid"
)

type BadgeTier string

const (
	BadgeTierCommunity BadgeTier = "community"
	BadgeTierVerified  BadgeTier = "verified"
	BadgeTierOfficial  BadgeTier = "official"
	BadgeTierSponsored BadgeTier = "sponsored"
)

type Badge struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"          json:"id"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Name        string    `gorm:"size:255;not null"             json:"name"`
	Description string    `gorm:"type:text"                     json:"description"`
	Icon        string    `gorm:"size:255"                      json:"icon"`
	Tier        BadgeTier `gorm:"type:varchar(20);not null"     json:"tier"`
	CreatedAt   time.Time `gorm:"autoCreateTime"                json:"created_at"`
}

func (Badge) TableName() string {
	return "badge"
}
