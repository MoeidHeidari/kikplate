package plate

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
)

type SubmitRepositoryInput struct {
	RepoURL string
	Branch  string
}

type SubmitFileInput struct {
	Name        string
	Description string
	Category    string
	Visibility  model.PlateVisibility
	Filename    string
	Content     string
	Tags        []string
}

type UpdatePlateInput struct {
	Name        *string
	Description *string
	Category    *string
	Visibility  *model.PlateVisibility
}

type PlateService interface {
	SubmitRepository(ctx context.Context, accountID uuid.UUID, input SubmitRepositoryInput) (*model.Plate, error)
	SubmitFile(ctx context.Context, accountID uuid.UUID, input SubmitFileInput) (*model.Plate, error)

	GetBySlug(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error)
	List(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error)

	Update(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input UpdatePlateInput) (*model.Plate, error)
	Archive(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	RecordUse(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	ReplaceTags(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, tags []string) error

	Approve(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	Reject(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	GrantBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string, reason *string) error
	RevokeBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string) error
}
