package organization

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
)

var (
	ErrNameTaken    = errors.New("organization name already taken")
	ErrNotOwner     = errors.New("only organization owner can perform this action")
	ErrHasPlates    = errors.New("cannot delete organization that contains plates")
	ErrNotFound     = errors.New("organization not found")
	ErrNameRequired = errors.New("organization name is required")
)

type OrganizationService interface {
	Create(ctx context.Context, input CreateOrganizationInput, ownerID uuid.UUID) (*model.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error)
	GetByName(ctx context.Context, name string) (*model.Organization, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*model.Organization, error)
	ListPublic(ctx context.Context, limit, offset int) ([]*model.Organization, int, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateOrganizationInput, requesterID uuid.UUID) (*model.Organization, error)
	Delete(ctx context.Context, id uuid.UUID, requesterID uuid.UUID) error
}

type CreateOrganizationInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type UpdateOrganizationInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	LogoURL     *string `json:"logo_url,omitempty"`
}

type organizationService struct {
	orgRepo repository.OrganizationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}

func (s *organizationService) Create(ctx context.Context, input CreateOrganizationInput, ownerID uuid.UUID) (*model.Organization, error) {
	if input.Name == "" {
		return nil, ErrNameRequired
	}

	existing, _ := s.orgRepo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, ErrNameTaken
	}

	org := &model.Organization{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		LogoURL:     normalizeOptionalString(input.LogoURL),
		OwnerID:     ownerID,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *organizationService) GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	return org, nil
}

func (s *organizationService) GetByName(ctx context.Context, name string) (*model.Organization, error) {
	org, err := s.orgRepo.GetByName(ctx, name)
	if err != nil {
		return nil, ErrNotFound
	}

	return org, nil
}

func (s *organizationService) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*model.Organization, error) {
	return s.orgRepo.ListByOwner(ctx, ownerID)
}

func (s *organizationService) ListPublic(ctx context.Context, limit, offset int) ([]*model.Organization, int, error) {
	return s.orgRepo.ListPublic(ctx, limit, offset)
}

func (s *organizationService) Update(ctx context.Context, id uuid.UUID, input UpdateOrganizationInput, requesterID uuid.UUID) (*model.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	if org.OwnerID != requesterID {
		return nil, ErrNotOwner
	}

	if input.Name != nil && *input.Name != "" {
		if *input.Name != org.Name {
			existing, _ := s.orgRepo.GetByName(ctx, *input.Name)
			if existing != nil {
				return nil, ErrNameTaken
			}
		}

		org.Name = *input.Name
	}

	if input.Description != nil {
		org.Description = *input.Description
	}

	if input.LogoURL != nil {
		org.LogoURL = normalizeOptionalString(*input.LogoURL)
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func normalizeOptionalString(v string) *string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (s *organizationService) Delete(ctx context.Context, id uuid.UUID, requesterID uuid.UUID) error {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if org.OwnerID != requesterID {
		return ErrNotOwner
	}

	count, err := s.orgRepo.CountPlates(ctx, id)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrHasPlates
	}

	return s.orgRepo.Delete(ctx, id)
}
