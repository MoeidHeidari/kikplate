package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type plateTagRepository struct {
	db *gorm.DB
}

func NewPlateTagRepository(db lib.Database) repository.PlateTagRepository {
	return &plateTagRepository{db: db.DB}
}

func (r *plateTagRepository) CreateMany(ctx context.Context, plateID uuid.UUID, tags []string) error {
	records := make([]*model.PlateTag, len(tags))
	for i, tag := range tags {
		records[i] = &model.PlateTag{
			ID:      uuid.New(),
			PlateID: plateID,
			Tag:     tag,
		}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *plateTagRepository) ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateTag, error) {
	var tags []*model.PlateTag
	result := r.db.WithContext(ctx).Where("plate_id = ?", plateID).Find(&tags)
	return tags, result.Error
}

func (r *plateTagRepository) DeleteByPlate(ctx context.Context, plateID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("plate_id = ?", plateID).
		Delete(&model.PlateTag{}).Error
}
