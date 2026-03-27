package plate

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
)

type SubmitRepositoryInput struct {
	RepoURL        string
	Branch         string
	OrganizationID *uuid.UUID
}

type UpdatePlateInput struct {
	Name        *string
	Description *string
	Category    *string
	Visibility  *model.PlateVisibility
}

type SubmitReviewInput struct {
	Rating int16
	Title  *string
	Body   *string
}

type PlateService interface {
	SubmitRepository(ctx context.Context, accountID uuid.UUID, input SubmitRepositoryInput) (*model.Plate, error)
	VerifyRepository(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) (*model.Plate, error)

	GetBySlug(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error)
	List(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error)

	GetMember(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) (*model.PlateMember, error)
	ListBookmarked(ctx context.Context, accountID uuid.UUID, limit int) ([]*model.Plate, error)
	Update(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input UpdatePlateInput) (*model.Plate, error)
	MoveToOrganization(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, organizationID *uuid.UUID) (*model.Plate, error)
	Archive(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	Remove(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	SetBookmark(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, bookmarked bool) error
	ReplaceTags(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, tags []string) error
	SubmitReview(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input SubmitReviewInput) (*model.PlateReview, error)

	Approve(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	Reject(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	GrantBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string, reason *string) error
	RevokeBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string) error
	GetStats(ctx context.Context) (*repository.PlateStats, error)
	GetFilterOptions(ctx context.Context) (*repository.PlateFilterOptions, error)
}
