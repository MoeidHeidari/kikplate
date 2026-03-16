package model

import (
	"time"

	"github.com/google/uuid"
)

type PlateReview struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"          json:"id"`
	PlateID       uuid.UUID `gorm:"type:uuid;not null;index"      json:"plate_id"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null;index"      json:"account_id"`
	Rating        int16     `gorm:"not null"                      json:"rating"`
	Title         *string   `gorm:"size:255"                      json:"title,omitempty"`
	Body          *string   `gorm:"type:text"                     json:"body,omitempty"`
	IsVerifiedUse bool      `gorm:"default:false"                 json:"is_verified_use"`
	CreatedAt     time.Time `gorm:"autoCreateTime"                json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"                json:"updated_at"`

	Plate   *Plate   `gorm:"foreignKey:PlateID"   json:"-"`
	Account *Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (PlateReview) TableName() string {
	return "plate_review"
}
