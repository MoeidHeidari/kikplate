package plate

import (
	"context"
	"net/url"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
)

func (s *plateService) GetBySlug(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error) {
	seen := map[string]struct{}{}
	baseCandidates := []string{slug}

	if unescaped, err := url.PathUnescape(slug); err == nil {
		baseCandidates = append(baseCandidates, unescaped)
	}
	if unescaped, err := url.QueryUnescape(slug); err == nil {
		baseCandidates = append(baseCandidates, unescaped)
	}

	candidates := append([]string{}, baseCandidates...)
	for _, v := range baseCandidates {
		candidates = append(candidates,
			strings.ReplaceAll(v, " ", "+"),
			strings.ReplaceAll(v, "%2B", "+"),
			strings.ReplaceAll(v, "%2b", "+"),
		)
	}

	var (
		plate *model.Plate
		err   error
	)

	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}

		plate, err = s.plates.GetBySlug(ctx, candidate)
		if err != nil {
			return nil, err
		}
		if plate != nil {
			break
		}
	}

	if plate == nil {
		return nil, ErrNotFound
	}

	if plate.Visibility == model.PlateVisibilityPrivate {
		member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
		if err != nil {
			return nil, err
		}
		if member == nil {
			return nil, ErrNotFound
		}
	}

	if requesterID != uuid.Nil {
		review, err := s.reviews.GetByPlateAndAccount(ctx, plate.ID, requesterID)
		if err == nil && review != nil {
			plate.UserRating = &review.Rating
		}
	}

	return plate, nil
}

func (s *plateService) List(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error) {
	filter.Categories = lib.NormalizePlateCategoryFilter(s.env, filter.Categories)
	plates, total, err := s.plates.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if filter.OwnerID == nil {
		// Repository List already restricts to public plates when OwnerID is nil.
		return plates, total, nil
	}

	var visible []*model.Plate
	for _, plate := range plates {
		if plate.Visibility == model.PlateVisibilityPublic {
			visible = append(visible, plate)
			continue
		}

		if plate.OwnerID == requesterID {
			visible = append(visible, plate)
			continue
		}

		member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
		if err != nil {
			continue
		}
		if member != nil {
			visible = append(visible, plate)
		}
	}

	adjustedTotal := total
	if total > len(visible) {
		adjustedTotal = len(visible)
	}

	return visible, adjustedTotal, nil
}

func (s *plateService) GetStats(ctx context.Context) (*repository.PlateStats, error) {
	st, err := s.plates.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	counts, err := s.GetCategoryCounts(ctx)
	if err != nil {
		return nil, err
	}
	st.TotalCategories = int64(len(counts))
	return st, nil
}

func (s *plateService) GetFilterOptions(ctx context.Context) (*repository.PlateFilterOptions, error) {
	opts, err := s.plates.GetFilterOptions(ctx)
	if err != nil {
		return nil, err
	}
	opts.Categories = lib.PlateCategorySlugs(s.env)
	return opts, nil
}

func (s *plateService) GetMonthlyGrowth(ctx context.Context, months int) ([]repository.MonthlyCount, error) {
	return s.plates.GetMonthlyGrowth(ctx, months)
}

func (s *plateService) GetCategoryCounts(ctx context.Context) ([]repository.CategoryCount, error) {
	rows, err := s.plates.GetCategoryCounts(ctx)
	if err != nil {
		return nil, err
	}
	merged := make(map[string]int64)
	for _, row := range rows {
		k := lib.NormalizePlateCategory(s.env, row.Category)
		merged[k] += row.Count
	}
	out := make([]repository.CategoryCount, 0, len(merged))
	for k, v := range merged {
		if v > 0 {
			out = append(out, repository.CategoryCount{Category: k, Count: v})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Category < out[j].Category
	})
	return out, nil
}

func (s *plateService) GetTopBookmarked(ctx context.Context, limit int) ([]repository.PlateRanked, error) {
	return s.plates.GetTopBookmarked(ctx, limit)
}

func (s *plateService) GetTopRated(ctx context.Context, limit int) ([]repository.PlateRanked, error) {
	return s.plates.GetTopRated(ctx, limit)
}

func (s *plateService) ListBookmarked(ctx context.Context, accountID uuid.UUID, limit int) ([]*model.Plate, error) {
	members, err := s.members.ListByAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	var plates []*model.Plate
	for _, member := range members {
		if !member.IsBookmarked {
			continue
		}
		plate, err := s.plates.GetByID(ctx, member.PlateID)
		if err != nil {
			continue
		}
		if plate != nil {
			plates = append(plates, plate)
		}
	}

	if limit > 0 && len(plates) > limit {
		plates = plates[:limit]
	}

	return plates, nil
}
