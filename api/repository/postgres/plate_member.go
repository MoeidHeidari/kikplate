package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type plateMemberRepository struct {
	db *gorm.DB
}

func NewPlateMemberRepository(db *gorm.DB) repository.PlateMemberRepository {
	return &plateMemberRepository{db: db}
}

func (r *plateMemberRepository) Create(ctx context.Context, member *model.PlateMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *plateMemberRepository) GetByPlateAndAccount(ctx context.Context, plateID, accountID uuid.UUID) (*model.PlateMember, error) {
	member := &model.PlateMember{}
	result := r.db.WithContext(ctx).
		Where("plate_id = ? AND account_id = ?", plateID, accountID).
		First(member)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return member, result.Error
}

func (r *plateMemberRepository) ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateMember, error) {
	var members []*model.PlateMember
	result := r.db.WithContext(ctx).Where("plate_id = ?", plateID).Find(&members)
	return members, result.Error
}

func (r *plateMemberRepository) ListByAccount(ctx context.Context, accountID uuid.UUID) ([]*model.PlateMember, error) {
	var members []*model.PlateMember
	result := r.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&members)
	return members, result.Error
}

func (r *plateMemberRepository) UpdateLastUsedAt(ctx context.Context, plateID, accountID uuid.UUID, t time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.PlateMember{}).
		Where("plate_id = ? AND account_id = ?", plateID, accountID).
		Update("last_used_at", t).Error
}

func (r *plateMemberRepository) Delete(ctx context.Context, plateID, accountID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("plate_id = ? AND account_id = ?", plateID, accountID).
		Delete(&model.PlateMember{}).Error
}
