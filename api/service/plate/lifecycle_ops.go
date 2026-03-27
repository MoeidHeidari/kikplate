package plate

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"gorm.io/gorm"
)

func (s *plateService) Update(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input UpdatePlateInput) (*model.Plate, error) {
	plate, err := s.plates.GetByID(ctx, plateID)
	if err != nil || plate == nil {
		return nil, ErrNotFound
	}

	if err := s.requireOwnerOrMember(ctx, plateID, accountID, model.PlateMemberRoleOwner); err != nil {
		return nil, err
	}

	if input.Name != nil {
		plate.Name = *input.Name
	}
	if input.Description != nil {
		plate.Description = input.Description
	}
	if input.Category != nil {
		plate.Category = *input.Category
	}
	if input.Visibility != nil {
		plate.Visibility = *input.Visibility
	}

	if err := s.plates.Update(ctx, plate); err != nil {
		return nil, err
	}

	return plate, nil
}

func (s *plateService) MoveToOrganization(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, organizationID *uuid.UUID) (*model.Plate, error) {
	plate, err := s.plates.GetByID(ctx, plateID)
	if err != nil || plate == nil {
		return nil, ErrNotFound
	}

	if err := s.requireOwnerOrMember(ctx, plateID, accountID, model.PlateMemberRoleOwner); err != nil {
		return nil, err
	}

	if organizationID != nil {
		if s.orgs == nil {
			return nil, ErrInvalidInput
		}

		org, err := s.orgs.GetByID(ctx, *organizationID)
		if err != nil || org == nil {
			return nil, ErrInvalidInput
		}
		if org.OwnerID != accountID {
			return nil, ErrForbidden
		}
	}

	if s.db != nil {
		if err := s.db.WithContext(ctx).
			Model(&model.Plate{}).
			Where("id = ?", plateID).
			Update("organization_id", organizationID).Error; err != nil {
			return nil, err
		}

		updated, err := s.plates.GetByID(ctx, plateID)
		if err != nil || updated == nil {
			return nil, ErrNotFound
		}
		return updated, nil
	}

	plate.OrganizationID = organizationID
	if err := s.plates.Update(ctx, plate); err != nil {
		return nil, err
	}

	return plate, nil
}

func (s *plateService) Archive(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error {
	p, err := s.plates.GetByID(ctx, plateID)
	if err != nil || p == nil {
		return ErrNotFound
	}

	if err := s.requireOwnerOrMember(ctx, plateID, accountID, model.PlateMemberRoleOwner); err != nil {
		return err
	}

	p.Status = model.PlateStatusArchived
	p.Slug = fmt.Sprintf("%s-archived-%s", p.Slug, p.ID.String()[:6])
	return s.plates.Update(ctx, p)
}

func (s *plateService) Remove(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error {
	p, err := s.plates.GetByID(ctx, plateID)
	if err != nil || p == nil {
		return ErrNotFound
	}

	if err := s.requireOwnerOrMember(ctx, plateID, accountID, model.PlateMemberRoleOwner); err != nil {
		return err
	}

	if err := s.tags.DeleteByPlate(ctx, plateID); err != nil {
		return err
	}

	members, err := s.members.ListByPlate(ctx, plateID)
	if err != nil {
		return err
	}
	for _, member := range members {
		if err := s.members.Delete(ctx, plateID, member.AccountID); err != nil {
			return err
		}
	}

	plateBadges, err := s.badges.ListByPlate(ctx, plateID)
	if err != nil {
		return err
	}
	for _, pb := range plateBadges {
		if err := s.badges.Revoke(ctx, plateID, pb.BadgeID); err != nil {
			return err
		}
	}

	reviews, err := s.reviews.ListByPlate(ctx, plateID)
	if err != nil {
		return err
	}
	for _, review := range reviews {
		if err := s.reviews.Delete(ctx, review.ID); err != nil {
			return err
		}
	}

	return s.plates.Delete(ctx, plateID)
}

func (s *plateService) SetBookmark(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, bookmarked bool) error {
	member, err := s.members.GetByPlateAndAccount(ctx, plateID, accountID)
	if err != nil {
		return err
	}

	if member == nil {
		if !bookmarked {
			return nil
		}
		if err := s.members.Create(ctx, &model.PlateMember{
			ID:           uuid.New(),
			PlateID:      plateID,
			AccountID:    accountID,
			Role:         model.PlateMemberRoleMember,
			IsBookmarked: true,
			BookmarkedAt: func(t time.Time) *time.Time { return &t }(time.Now()),
		}); err != nil {
			return err
		}
		return s.plates.IncrementBookmarkCount(ctx, plateID)
	}

	wasBookmarked := member.IsBookmarked
	if wasBookmarked == bookmarked {
		return nil
	}

	if err := s.members.SetBookmarked(ctx, plateID, accountID, bookmarked); err != nil {
		return err
	}

	if bookmarked {
		return s.plates.IncrementBookmarkCount(ctx, plateID)
	} else {
		return s.plates.DecrementBookmarkCount(ctx, plateID)
	}
}

func (s *plateService) GetMember(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) (*model.PlateMember, error) {
	return s.members.GetByPlateAndAccount(ctx, plateID, accountID)
}

func (s *plateService) ReplaceTags(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, tags []string) error {
	if err := s.requireOwnerOrMember(ctx, plateID, accountID, model.PlateMemberRoleOwner); err != nil {
		return err
	}

	do := func() error {
		if err := s.tags.DeleteByPlate(ctx, plateID); err != nil {
			return err
		}
		if len(tags) == 0 {
			return nil
		}
		return s.tags.CreateMany(ctx, plateID, tags)
	}

	if s.db != nil {
		return s.db.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
			return do()
		})
	}
	return do()
}

func (s *plateService) Approve(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error {
	plate, err := s.plates.GetByID(ctx, plateID)
	if err != nil || plate == nil {
		return ErrNotFound
	}

	now := time.Now()
	plate.Status = model.PlateStatusApproved
	plate.PublishedAt = &now
	return s.plates.Update(ctx, plate)
}

func (s *plateService) Reject(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error {
	plate, err := s.plates.GetByID(ctx, plateID)
	if err != nil || plate == nil {
		return ErrNotFound
	}

	plate.Status = model.PlateStatusRejected
	return s.plates.Update(ctx, plate)
}
