package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type plateBadgeRepository struct {
	db *gorm.DB
}

func NewPlateBadgeRepository(db *gorm.DB) repository.PlateBadgeRepository {
	return &plateBadgeRepository{db: db}
}

func (r *plateBadgeRepository) Grant(ctx context.Context, pb *model.PlateBadge) error {
	return r.db.WithContext(ctx).Create(pb).Error
}

func (r *plateBadgeRepository) ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateBadge, error) {
	var badges []*model.PlateBadge
	result := r.db.WithContext(ctx).
		Preload("Badge").
		Where("plate_id = ?", plateID).
		Find(&badges)
	return badges, result.Error
}

func (r *plateBadgeRepository) Revoke(ctx context.Context, plateID, badgeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("plate_id = ? AND badge_id = ?", plateID, badgeID).
		Delete(&model.PlateBadge{}).Error
}
