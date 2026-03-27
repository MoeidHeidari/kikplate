package plate

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
)

func (s *plateService) SubmitReview(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input SubmitReviewInput) (*model.PlateReview, error) {
	p, err := s.plates.GetByID(ctx, plateID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrNotFound
	}

	if p.OwnerID == accountID {
		return nil, ErrCannotReviewOwn
	}

	existing, err := s.reviews.GetByPlateAndAccount(ctx, plateID, accountID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyReviewed
	}

	review := &model.PlateReview{
		ID:        uuid.New(),
		PlateID:   plateID,
		AccountID: accountID,
		Rating:    input.Rating,
		Title:     input.Title,
		Body:      input.Body,
	}

	if err := s.reviews.Create(ctx, review); err != nil {
		return nil, err
	}

	reviews, err := s.reviews.ListByPlate(ctx, plateID)
	if err == nil && len(reviews) > 0 {
		sum := 0
		for _, r := range reviews {
			sum += int(r.Rating)
		}
		p.AvgRating = float64(sum) / float64(len(reviews))
		_ = s.plates.Update(ctx, p)
	}

	return review, nil
}

func (s *plateService) GrantBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string, reason *string) error {
	badge, err := s.badgeCatalog.GetBySlug(ctx, badgeSlug)
	if err != nil || badge == nil {
		return ErrNotFound
	}

	return s.badges.Grant(ctx, &model.PlateBadge{
		ID:        uuid.New(),
		PlateID:   plateID,
		BadgeID:   badge.ID,
		GrantedBy: adminAccountID.String(),
		Reason:    reason,
	})
}

func (s *plateService) RevokeBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string) error {
	badge, err := s.badgeCatalog.GetBySlug(ctx, badgeSlug)
	if err != nil || badge == nil {
		return ErrNotFound
	}

	return s.badges.Revoke(ctx, plateID, badge.ID)
}

func (s *plateService) requireOwnerOrMember(ctx context.Context, plateID, accountID uuid.UUID, requiredRole model.PlateMemberRole) error {
	member, err := s.members.GetByPlateAndAccount(ctx, plateID, accountID)
	if err != nil {
		return err
	}
	if member == nil || member.Role != requiredRole {
		return ErrForbidden
	}
	return nil
}
