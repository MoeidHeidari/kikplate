package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type plateRepository struct {
	db *gorm.DB
}

func NewPlateRepository(db lib.Database) repository.PlateRepository {
	return &plateRepository{db: db.DB}
}

func (r *plateRepository) Create(ctx context.Context, plate *model.Plate) error {
	return r.db.WithContext(ctx).Create(plate).Error
}

func (r *plateRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Plate, error) {
	plate := &model.Plate{}
	result := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Badges.Badge").
		First(plate, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return plate, result.Error
}

func (r *plateRepository) GetBySlug(ctx context.Context, slug string) (*model.Plate, error) {
	plate := &model.Plate{}
	result := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Badges.Badge").
		Where("slug = ?", slug).
		First(plate)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return plate, result.Error
}

func (r *plateRepository) List(ctx context.Context, filter repository.PlateFilter) ([]*model.Plate, int, error) {
	var plates []*model.Plate
	var total int64

	q := r.db.WithContext(ctx).Model(&model.Plate{})

	if filter.Type != nil {
		q = q.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.Visibility != nil {
		q = q.Where("visibility = ?", *filter.Visibility)
	}
	if filter.Category != "" {
		q = q.Where("category = ?", filter.Category)
	}
	if filter.OwnerID != nil {
		q = q.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.Search != "" {
		q = q.Where("search_vector @@ plainto_tsquery('english', ?)", filter.Search)
	}
	if len(filter.Tags) > 0 {
		q = q.Joins("JOIN plate_tag ON plate_tag.plate_id = plate.id").
			Where("plate_tag.tag IN ?", filter.Tags).
			Group("plate.id")
	}

	q.Count(&total)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}

	result := q.
		Preload("Tags").
		Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Find(&plates)

	return plates, int(total), result.Error
}

func (r *plateRepository) Update(ctx context.Context, plate *model.Plate) error {
	return r.db.WithContext(ctx).Save(plate).Error
}

func (r *plateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Plate{}, "id = ?", id).Error
}

func (r *plateRepository) IncrementUseCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Plate{}).
		Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

func (r *plateRepository) UpdateSyncState(ctx context.Context, id uuid.UUID, state repository.PlateSyncState) error {
	return r.db.WithContext(ctx).
		Model(&model.Plate{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"sync_status":          state.SyncStatus,
			"sync_error":           state.SyncError,
			"last_synced_at":       state.LastSyncedAt,
			"next_sync_at":         state.NextSyncAt,
			"consecutive_failures": state.ConsecutiveFailures,
			"is_verified":          state.IsVerified,
			"verified_at":          state.VerifiedAt,
			"metadata":             state.Metadata,
		}).Error
}

func (r *plateRepository) ListDueForSync(ctx context.Context, limit int) ([]*model.Plate, error) {
	var plates []*model.Plate
	result := r.db.WithContext(ctx).
		Where("type = ? AND next_sync_at <= NOW() AND sync_status != ?",
			model.PlateTypeRepository, model.SyncStatusSyncing).
		Order("next_sync_at ASC").
		Limit(limit).
		Find(&plates)
	return plates, result.Error
}
