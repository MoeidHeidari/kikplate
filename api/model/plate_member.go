package model

import (
	"time"

	"github.com/google/uuid"
)

type PlateMemberRole string

const (
	PlateMemberRoleOwner  PlateMemberRole = "owner"
	PlateMemberRoleMember PlateMemberRole = "member"
)

type PlateMember struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey"              json:"id"`
	PlateID      uuid.UUID       `gorm:"type:uuid;not null;index"          json:"plate_id"`
	AccountID    uuid.UUID       `gorm:"type:uuid;not null;index"          json:"account_id"`
	Role         PlateMemberRole `gorm:"type:varchar(20);not null"         json:"role"`
	IsBookmarked bool            `gorm:"default:false"                     json:"is_bookmarked"`
	BookmarkedAt *time.Time      `                                          json:"bookmarked_at,omitempty"`
	JoinedAt     time.Time       `gorm:"autoCreateTime"                    json:"joined_at"`

	Plate   *Plate   `gorm:"foreignKey:PlateID"   json:"-"`
	Account *Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (PlateMember) TableName() string {
	return "plate_member"
}
