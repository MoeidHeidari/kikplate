package plate

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/uuid"
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
	plates, total, err := s.plates.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	if filter.OwnerID == nil {
		visible := make([]*model.Plate, 0, len(plates))
		for _, plate := range plates {
			if plate.Visibility == model.PlateVisibilityPublic {
				visible = append(visible, plate)
			}
		}
		return visible, len(visible), nil
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
	return s.plates.GetStats(ctx)
}

func (s *plateService) GetFilterOptions(ctx context.Context) (*repository.PlateFilterOptions, error) {
	return s.plates.GetFilterOptions(ctx)
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
