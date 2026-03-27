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
		Preload("Organization").
		Preload("Organization.Owner").
		Preload("Owner").
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
		Preload("Organization").
		Preload("Organization.Owner").
		Preload("Owner").
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

	if len(filter.Types) > 0 {
		q = q.Where("type IN ?", filter.Types)
	}
	if len(filter.Categories) > 0 {
		q = q.Where("category IN ?", filter.Categories)
	}
	if filter.OwnerID != nil {
		q = q.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.Search != "" {
		q = q.Where(
			`search_vector @@ websearch_to_tsquery('english', ?)
			 OR name ILIKE ?
			 OR description ILIKE ?
			 OR plate.id IN (SELECT plate_id FROM plate_tag WHERE tag ILIKE ?)`,
			filter.Search,
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
			"%"+filter.Search+"%",
		)
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

	var result *gorm.DB
	if filter.Search != "" {
		result = q.
			Preload("Tags").
			Preload("Organization").
			Preload("Organization.Owner").
			Offset((page - 1) * limit).
			Limit(limit).
			Order(gorm.Expr("ts_rank_cd(search_vector, websearch_to_tsquery('english', ?)) DESC, bookmark_count DESC, created_at DESC", filter.Search)).
			Find(&plates)
	} else {
		result = q.
			Preload("Tags").
			Preload("Organization").
			Preload("Organization.Owner").
			Offset((page - 1) * limit).
			Limit(limit).
			Order("bookmark_count DESC, created_at DESC").
			Find(&plates)
	}

	return plates, int(total), result.Error
}

func (r *plateRepository) Update(ctx context.Context, plate *model.Plate) error {
	return r.db.WithContext(ctx).Save(plate).Error
}

func (r *plateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Plate{}, "id = ?", id).Error
}

func (r *plateRepository) IncrementBookmarkCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Plate{}).
		Where("id = ?", id).
		UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + 1")).Error
}

func (r *plateRepository) DecrementBookmarkCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.Plate{}).
		Where("id = ?", id).
		UpdateColumn("bookmark_count", gorm.Expr("bookmark_count - 1")).Error
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
func (r *plateRepository) GetStats(ctx context.Context) (*repository.PlateStats, error) {
	var stats repository.PlateStats

	r.db.WithContext(ctx).Model(&model.Plate{}).
		Where("status = ?", model.PlateStatusApproved).
		Count(&stats.TotalPlates)

	r.db.WithContext(ctx).Model(&model.Plate{}).
		Where("status = ?", model.PlateStatusApproved).
		Distinct("owner_id").
		Count(&stats.TotalContributors)

	r.db.WithContext(ctx).Model(&model.Plate{}).
		Where("status = ?", model.PlateStatusApproved).
		Distinct("category").
		Count(&stats.TotalCategories)

	r.db.WithContext(ctx).Model(&model.Plate{}).
		Where("status = ?", model.PlateStatusApproved).
		Select("coalesce(sum(bookmark_count), 0)").
		Scan(&stats.TotalBookmarks)

	return &stats, nil
}

func (r *plateRepository) GetFilterOptions(ctx context.Context) (*repository.PlateFilterOptions, error) {
	var categories []string
	if err := r.db.WithContext(ctx).
		Model(&model.Plate{}).
		Distinct("category").
		Order("category ASC").
		Where("category <> ''").
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}

	var tags []string
	if err := r.db.WithContext(ctx).
		Model(&model.PlateTag{}).
		Distinct("tag").
		Order("tag ASC").
		Where("tag <> ''").
		Pluck("tag", &tags).Error; err != nil {
		return nil, err
	}

	return &repository.PlateFilterOptions{
		Categories: categories,
		Tags:       tags,
	}, nil
}

func (r *plateRepository) ListDueForSync(ctx context.Context, limit int) ([]*model.Plate, error) {
	var plates []*model.Plate
	result := r.db.WithContext(ctx).
		Where("type = ? AND (next_sync_at <= NOW() OR (sync_status = ? AND updated_at <= NOW() - INTERVAL '2 minutes'))",
			model.PlateTypeRepository, model.SyncStatusSyncing).
		Order("next_sync_at ASC").
		Limit(limit).
		Find(&plates)
	return plates, result.Error
}
