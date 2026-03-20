package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type syncLogRepository struct {
	db *gorm.DB
}

func NewSyncLogRepository(db lib.Database) repository.SyncLogRepository {
	return &syncLogRepository{db: db.DB}
}

func (r *syncLogRepository) Create(ctx context.Context, log *model.SyncLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *syncLogRepository) ListByPlate(ctx context.Context, plateID uuid.UUID, limit int) ([]*model.SyncLog, error) {
	var logs []*model.SyncLog
	result := r.db.WithContext(ctx).
		Where("plate_id = ?", plateID).
		Order("synced_at DESC").
		Limit(limit).
		Find(&logs)
	return logs, result.Error
}

func (r *syncLogRepository) DeleteOlderThan(ctx context.Context, t time.Time) error {
	return r.db.WithContext(ctx).
		Where("synced_at < ?", t).
		Delete(&model.SyncLog{}).Error
}
