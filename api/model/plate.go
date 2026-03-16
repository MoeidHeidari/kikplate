package model

import (
	"time"

	"github.com/google/uuid"
)

type PlateType string
type PlateStatus string
type PlateVisibility string
type SyncStatus string

const (
	PlateTypeRepository PlateType = "repository"
	PlateTypeFile       PlateType = "file"
)

const (
	PlateStatusPending  PlateStatus = "pending"
	PlateStatusApproved PlateStatus = "approved"
	PlateStatusRejected PlateStatus = "rejected"
	PlateStatusArchived PlateStatus = "archived"
)

const (
	PlateVisibilityPublic   PlateVisibility = "public"
	PlateVisibilityPrivate  PlateVisibility = "private"
	PlateVisibilityUnlisted PlateVisibility = "unlisted"
)

const (
	SyncStatusPending    SyncStatus = "pending"
	SyncStatusSyncing    SyncStatus = "syncing"
	SyncStatusSynced     SyncStatus = "synced"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusUnverified SyncStatus = "unverified"
)

type Plate struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey"                        json:"id"`
	OwnerID     uuid.UUID       `gorm:"type:uuid;not null;index"                    json:"owner_id"`
	Type        PlateType       `gorm:"type:varchar(20);not null"                   json:"type"`
	Slug        string          `gorm:"size:255;not null;uniqueIndex"               json:"slug"`
	Name        string          `gorm:"size:255;not null"                           json:"name"`
	Description *string         `gorm:"type:text"                                   json:"description,omitempty"`
	Category    string          `gorm:"size:100;not null;index"                     json:"category"`
	Status      PlateStatus     `gorm:"type:varchar(20);default:pending;index"      json:"status"`
	Visibility  PlateVisibility `gorm:"type:varchar(20);default:public"             json:"visibility"`
	Metadata    []byte          `gorm:"type:jsonb"                                  json:"metadata,omitempty"`
	UseCount    int             `gorm:"default:0"                                   json:"use_count"`
	StarCount   int             `gorm:"default:0"                                   json:"star_count"`
	AvgRating   float64         `gorm:"type:numeric(3,2);default:0"                 json:"avg_rating"`
	IsVerified  bool            `gorm:"default:false"                               json:"is_verified"`
	VerifiedAt  *time.Time      `                                                   json:"verified_at,omitempty"`
	PublishedAt *time.Time      `                                                   json:"published_at,omitempty"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"                              json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"                              json:"updated_at"`

	RepoURL             *string     `gorm:"type:text"                               json:"repo_url,omitempty"`
	Branch              *string     `gorm:"size:255"                                json:"branch,omitempty"`
	SyncStatus          *SyncStatus `gorm:"type:varchar(20)"                        json:"sync_status,omitempty"`
	SyncError           *string     `gorm:"type:text"                               json:"sync_error,omitempty"`
	SyncInterval        *string     `gorm:"type:interval"                           json:"sync_interval,omitempty"`
	NextSyncAt          *time.Time  `gorm:"index"                                   json:"next_sync_at,omitempty"`
	LastSyncedAt        *time.Time  `                                               json:"last_synced_at,omitempty"`
	ConsecutiveFailures int         `gorm:"default:0"                               json:"consecutive_failures"`

	Content  *string `gorm:"type:text"   json:"content,omitempty"`
	Filename *string `gorm:"size:255"    json:"filename,omitempty"`

	// associations
	Owner   *Account      `gorm:"foreignKey:OwnerID"  json:"-"`
	Tags    []PlateTag    `gorm:"foreignKey:PlateID"  json:"tags,omitempty"`
	Members []PlateMember `gorm:"foreignKey:PlateID"  json:"-"`
	Reviews []PlateReview `gorm:"foreignKey:PlateID"  json:"-"`
	Badges  []PlateBadge  `gorm:"foreignKey:PlateID"  json:"-"`
}

func (Plate) TableName() string {
	return "plate"
}
