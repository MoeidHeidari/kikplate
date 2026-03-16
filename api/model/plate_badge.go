package model

import (
	"time"

	"github.com/google/uuid"
)

type PlateBadge struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"      json:"id"`
	PlateID   uuid.UUID `gorm:"type:uuid;not null;index"  json:"plate_id"`
	BadgeID   uuid.UUID `gorm:"type:uuid;not null;index"  json:"badge_id"`
	GrantedBy string    `gorm:"size:255;not null"         json:"granted_by"`
	Reason    *string   `gorm:"type:text"                 json:"reason,omitempty"`
	GrantedAt time.Time `gorm:"autoCreateTime"            json:"granted_at"`

	Plate *Plate `gorm:"foreignKey:PlateID" json:"-"`
	Badge *Badge `gorm:"foreignKey:BadgeID" json:"-"`
}

func (PlateBadge) TableName() string {
	return "plate_badge"
}
