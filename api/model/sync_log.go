package model

import (
	"time"

	"github.com/google/uuid"
)

type SyncLogStatus string

const (
	SyncLogStatusSuccess SyncLogStatus = "success"
	SyncLogStatusFailed  SyncLogStatus = "failed"
	SyncLogStatusSkipped SyncLogStatus = "skipped"
)

type SyncLog struct {
	ID              uuid.UUID     `gorm:"type:uuid;primaryKey"      json:"id"`
	PlateID         uuid.UUID     `gorm:"type:uuid;not null;index"  json:"plate_id"`
	Status          SyncLogStatus `gorm:"type:varchar(20);not null" json:"status"`
	ChangesDetected bool          `gorm:"default:false"             json:"changes_detected"`
	Error           *string       `gorm:"type:text"                 json:"error,omitempty"`
	DurationMs      int           `gorm:"default:0"                 json:"duration_ms"`
	SyncedAt        time.Time     `gorm:"autoCreateTime"            json:"synced_at"`

	Plate *Plate `gorm:"foreignKey:PlateID" json:"-"`
}

func (SyncLog) TableName() string {
	return "sync_log"
}
