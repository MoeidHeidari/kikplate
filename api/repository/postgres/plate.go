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
	if filter.OrganizationID != nil {
		q = q.Where("organization_id = ?", *filter.OrganizationID)
	}
	if filter.OwnerID == nil {
		q = q.Where("visibility = ?", model.PlateVisibilityPublic)
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
	if len(filter.Badges) > 0 {
		q = q.Where(
			`plate.id IN (SELECT pb.plate_id FROM plate_badge pb INNER JOIN badge b ON b.id = pb.badge_id WHERE b.slug IN ?)`,
			filter.Badges,
		)
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
			Preload("Badges.Badge").
			Preload("Organization").
			Preload("Organization.Owner").
			Offset((page - 1) * limit).
			Limit(limit).
			Order(gorm.Expr("ts_rank_cd(search_vector, websearch_to_tsquery('english', ?)) DESC, bookmark_count DESC, created_at DESC", filter.Search)).
			Find(&plates)
	} else {
		result = q.
			Preload("Tags").
			Preload("Badges.Badge").
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

func (r *plateRepository) GetMonthlyGrowth(ctx context.Context, months int) ([]repository.MonthlyCount, error) {
	var rows []repository.MonthlyCount
	err := r.db.WithContext(ctx).Raw(`
		SELECT to_char(date_trunc('month', created_at), 'YYYY-MM') AS month,
		       COUNT(*) AS count
		FROM plate
		WHERE status = ? AND created_at >= date_trunc('month', NOW()) - INTERVAL '1 month' * ?
		GROUP BY date_trunc('month', created_at)
		ORDER BY month
	`, model.PlateStatusApproved, months).Scan(&rows).Error
	return rows, err
}

func (r *plateRepository) GetCategoryCounts(ctx context.Context) ([]repository.CategoryCount, error) {
	var rows []repository.CategoryCount
	err := r.db.WithContext(ctx).Raw(`
		SELECT category, COUNT(*) AS count
		FROM plate
		WHERE status = ? AND category <> ''
		GROUP BY category
		ORDER BY count DESC
	`, model.PlateStatusApproved).Scan(&rows).Error
	return rows, err
}

func (r *plateRepository) GetTopBookmarked(ctx context.Context, limit int) ([]repository.PlateRanked, error) {
	var rows []repository.PlateRanked
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, slug, name, bookmark_count, avg_rating, category
		FROM plate
		WHERE status = ? AND bookmark_count > 0
		ORDER BY bookmark_count DESC
		LIMIT ?
	`, model.PlateStatusApproved, limit).Scan(&rows).Error
	return rows, err
}

func (r *plateRepository) GetTopRated(ctx context.Context, limit int) ([]repository.PlateRanked, error) {
	var rows []repository.PlateRanked
	err := r.db.WithContext(ctx).Raw(`
		SELECT id, slug, name, bookmark_count, avg_rating, category
		FROM plate
		WHERE status = ? AND avg_rating > 0
		ORDER BY avg_rating DESC, bookmark_count DESC
		LIMIT ?
	`, model.PlateStatusApproved, limit).Scan(&rows).Error
	return rows, err
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

	var badgeRows []struct {
		Slug string `gorm:"column:slug"`
		Name string `gorm:"column:name"`
	}
	if err := r.db.WithContext(ctx).
		Table("badge").
		Select("DISTINCT badge.slug, badge.name").
		Joins("INNER JOIN plate_badge ON plate_badge.badge_id = badge.id").
		Order("badge.name ASC").
		Scan(&badgeRows).Error; err != nil {
		return nil, err
	}
	badges := make([]repository.BadgeOption, 0, len(badgeRows))
	for _, row := range badgeRows {
		if row.Slug == "" {
			continue
		}
		badges = append(badges, repository.BadgeOption{Slug: row.Slug, Name: row.Name})
	}

	return &repository.PlateFilterOptions{
		Categories: categories,
		Tags:       tags,
		Badges:     badges,
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
