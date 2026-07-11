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

	isPrivateRepoPlate, privateRepoErr := s.isPrivateRepositoryPlate(ctx, plate)
	if privateRepoErr != nil {
		s.logger.Warnf("private repository visibility check failed plate_id=%s err=%v", plate.ID, privateRepoErr)
	}
	if isPrivateRepoPlate {
		if requesterID == uuid.Nil {
			return nil, ErrNotFound
		}
		if plate.OrganizationID != nil {
			hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
			if err != nil {
				return nil, err
			}
			if !hasOrgAccess {
				return nil, ErrNotFound
			}
		} else {
			member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
			if err != nil {
				return nil, err
			}
			if plate.OwnerID != requesterID && member == nil {
				return nil, ErrNotFound
			}
		}
	}

	isPrivateOrgPlate, err := s.isPrivateOrganizationPlate(ctx, plate)
	if err != nil {
		return nil, err
	}
	if isPrivateOrgPlate {
		if requesterID == uuid.Nil {
			return nil, ErrNotFound
		}
		hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
		if err != nil {
			return nil, err
		}
		if !hasOrgAccess {
			return nil, ErrNotFound
		}
	}

	if plate.Visibility == model.PlateVisibilityPrivate {
		if plate.OrganizationID != nil && requesterID != uuid.Nil {
			hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
			if err != nil {
				return nil, err
			}
			if !hasOrgAccess {
				return nil, ErrNotFound
			}
		} else {
			member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
			if err != nil {
				return nil, err
			}
			if member == nil {
				return nil, ErrNotFound
			}
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
	if filter.OrganizationID != nil && requesterID != uuid.Nil {
		hasOrgAccess, err := s.hasOrganizationAccess(ctx, *filter.OrganizationID, requesterID)
		if err == nil && hasOrgAccess {
			filter.AccessibleOrganizationIDs = []uuid.UUID{*filter.OrganizationID}
		}
	} else if filter.OrganizationID == nil && requesterID != uuid.Nil && s.orgMembers != nil {
		members, err := s.orgMembers.ListByAccount(ctx, requesterID)
		if err == nil {
			for _, member := range members {
				if member != nil && member.Status == model.OrganizationMemberStatusAccepted {
					filter.AccessibleOrganizationIDs = append(filter.AccessibleOrganizationIDs, member.OrganizationID)
				}
			}
		}
	}
	plates, total, err := s.plates.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var visible []*model.Plate
	for _, plate := range plates {
		isPrivateRepoPlate, privateRepoErr := s.isPrivateRepositoryPlate(ctx, plate)
		if privateRepoErr != nil {
			s.logger.Warnf("private repository visibility check failed plate_id=%s err=%v", plate.ID, privateRepoErr)
		}
		if isPrivateRepoPlate {
			if requesterID == uuid.Nil {
				continue
			}
			if plate.OrganizationID != nil {
				hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
				if err != nil || !hasOrgAccess {
					continue
				}
				visible = append(visible, plate)
				continue
			}
			member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
			if err != nil {
				continue
			}
			if plate.OwnerID == requesterID || member != nil {
				visible = append(visible, plate)
			}
			continue
		}

		isPrivateOrgPlate, err := s.isPrivateOrganizationPlate(ctx, plate)
		if err != nil {
			continue
		}
		if isPrivateOrgPlate {
			if requesterID == uuid.Nil {
				continue
			}
			hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
			if err != nil || !hasOrgAccess {
				continue
			}
			visible = append(visible, plate)
			continue
		}

		if filter.OwnerID == nil {
			if plate.Visibility == model.PlateVisibilityPublic {
				visible = append(visible, plate)
			}
			continue
		}

		if plate.Visibility == model.PlateVisibilityPublic {
			visible = append(visible, plate)
			continue
		}

		if plate.OrganizationID != nil {
			hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, requesterID)
			if err != nil {
				continue
			}
			if hasOrgAccess {
				visible = append(visible, plate)
			}
			continue
		}

		member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
		if err != nil {
			continue
		}
		if plate.OwnerID == requesterID || member != nil {
			visible = append(visible, plate)
		}
	}

	return visible, total, nil
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
	agg, err := s.plates.GetExplorerFilterAggregates(ctx)
	if err != nil {
		return nil, err
	}
	countsByCanonicalSlug := accumulatePlateCategoryTotalsWithConfigCanonicalization(s.env, agg.CategoryCounts)
	categories := explorerCategoryRowsAlignedWithConfigSlugs(s.env, countsByCanonicalSlug)

	return &repository.PlateFilterOptions{
		Categories: categories,
		Tags:       agg.TagOptions,
		Badges:     agg.BadgeOptions,
	}, nil
}

func accumulatePlateCategoryTotalsWithConfigCanonicalization(env lib.Env, rows []repository.CategoryCount) map[string]int64 {
	out := make(map[string]int64)
	for _, row := range rows {
		k := lib.NormalizePlateCategory(env, row.Category)
		out[k] += row.Count
	}
	return out
}

func explorerCategoryRowsAlignedWithConfigSlugs(env lib.Env, countsByCanonicalSlug map[string]int64) []repository.CategoryFilterOption {
	slugs := lib.PlateCategorySlugs(env)
	opts := make([]repository.CategoryFilterOption, 0, len(slugs))
	for _, slug := range slugs {
		opts = append(opts, repository.CategoryFilterOption{
			Slug:  slug,
			Count: countsByCanonicalSlug[slug],
		})
	}
	return opts
}

func (s *plateService) GetMonthlyGrowth(ctx context.Context, months int) ([]repository.MonthlyCount, error) {
	return s.plates.GetMonthlyGrowth(ctx, months)
}

func (s *plateService) GetCategoryCounts(ctx context.Context) ([]repository.CategoryCount, error) {
	rows, err := s.plates.GetCategoryCounts(ctx)
	if err != nil {
		return nil, err
	}
	merged := accumulatePlateCategoryTotalsWithConfigCanonicalization(s.env, rows)
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
		if plate == nil {
			continue
		}

		isPrivateOrgPlate, err := s.isPrivateOrganizationPlate(ctx, plate)
		if err != nil {
			continue
		}
		if isPrivateOrgPlate {
			hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, accountID)
			if err != nil || !hasOrgAccess {
				continue
			}
			plates = append(plates, plate)
			continue
		}

		if plate.Visibility == model.PlateVisibilityPrivate {
			if plate.OrganizationID != nil {
				hasOrgAccess, err := s.hasOrganizationAccess(ctx, *plate.OrganizationID, accountID)
				if err != nil || !hasOrgAccess {
					continue
				}
			} else if plate.OwnerID != accountID {
				pm, err := s.members.GetByPlateAndAccount(ctx, plate.ID, accountID)
				if err != nil || pm == nil {
					continue
				}
			}
		}

		plates = append(plates, plate)
	}

	if limit > 0 && len(plates) > limit {
		plates = plates[:limit]
	}

	return plates, nil
}

func (s *plateService) isPrivateOrganizationPlate(ctx context.Context, plate *model.Plate) (bool, error) {
	if !s.env.Features.PrivateOrganizationsEnabled || plate == nil || plate.OrganizationID == nil {
		return false, nil
	}
	if plate.Organization != nil {
		return plate.Organization.Visibility == model.OrganizationVisibilityPrivate, nil
	}
	if s.orgs == nil {
		return false, nil
	}
	org, err := s.orgs.GetByID(ctx, *plate.OrganizationID)
	if err != nil {
		return false, err
	}
	if org == nil {
		return false, nil
	}
	plate.Organization = org
	return org.Visibility == model.OrganizationVisibilityPrivate, nil
}

func (s *plateService) isPrivateRepositoryPlate(ctx context.Context, plate *model.Plate) (bool, error) {
	if plate == nil || plate.Type != model.PlateTypeRepository || plate.RepoURL == nil {
		return false, nil
	}
	ownerAccountID := plate.OwnerID.String()
	var organizationIDStr *string
	if plate.OrganizationID != nil {
		orgID := plate.OrganizationID.String()
		organizationIDStr = &orgID
	}
	return s.fetchRepositoryVisibility(ctx, *plate.RepoURL, ownerAccountID, organizationIDStr)
}
