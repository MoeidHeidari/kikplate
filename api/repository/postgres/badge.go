package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type badgeRepository struct {
	db *gorm.DB
}

func NewBadgeRepository(db lib.Database) repository.BadgeRepository {
	return &badgeRepository{db: db.DB}
}

func (r *badgeRepository) Create(ctx context.Context, badge *model.Badge) error {
	return r.db.WithContext(ctx).Create(badge).Error
}

func (r *badgeRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Badge, error) {
	badge := &model.Badge{}
	result := r.db.WithContext(ctx).First(badge, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return badge, result.Error
}

func (r *badgeRepository) GetBySlug(ctx context.Context, slug string) (*model.Badge, error) {
	badge := &model.Badge{}
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(badge)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return badge, result.Error
}

func (r *badgeRepository) List(ctx context.Context) ([]*model.Badge, error) {
	var badges []*model.Badge
	result := r.db.WithContext(ctx).Order("tier, name").Find(&badges)
	return badges, result.Error
}
