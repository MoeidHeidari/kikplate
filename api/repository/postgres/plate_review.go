package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type plateReviewRepository struct {
	db *gorm.DB
}

func NewPlateReviewRepository(db *gorm.DB) repository.PlateReviewRepository {
	return &plateReviewRepository{db: db}
}

func (r *plateReviewRepository) Create(ctx context.Context, review *model.PlateReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *plateReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PlateReview, error) {
	review := &model.PlateReview{}
	result := r.db.WithContext(ctx).First(review, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return review, result.Error
}

func (r *plateReviewRepository) GetByPlateAndAccount(ctx context.Context, plateID, accountID uuid.UUID) (*model.PlateReview, error) {
	review := &model.PlateReview{}
	result := r.db.WithContext(ctx).
		Where("plate_id = ? AND account_id = ?", plateID, accountID).
		First(review)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return review, result.Error
}

func (r *plateReviewRepository) ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateReview, error) {
	var reviews []*model.PlateReview
	result := r.db.WithContext(ctx).
		Where("plate_id = ?", plateID).
		Order("created_at DESC").
		Find(&reviews)
	return reviews, result.Error
}

func (r *plateReviewRepository) Update(ctx context.Context, review *model.PlateReview) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *plateReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.PlateReview{}, "id = ?", id).Error
}
