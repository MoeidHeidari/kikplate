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

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db lib.Database) repository.EmailVerificationRepository {
	return &emailVerificationRepository{db: db.DB}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *model.EmailVerification) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *emailVerificationRepository) GetByToken(ctx context.Context, token string) (*model.EmailVerification, error) {
	ev := &model.EmailVerification{}
	result := r.db.WithContext(ctx).
		Where("token = ? AND is_used = false AND expires_at > ?", token, time.Now()).
		First(ev)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ev, result.Error
}

func (r *emailVerificationRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.EmailVerification{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *emailVerificationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR is_used = true", time.Now()).
		Delete(&model.EmailVerification{}).Error
}
