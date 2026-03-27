package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db lib.Database) repository.OrganizationRepository {
	return &organizationRepository{db: db.DB}
}

func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	err := r.db.WithContext(ctx).Preload("Owner").First(&org, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetByName(ctx context.Context, name string) (*model.Organization, error) {
	var org model.Organization
	err := r.db.WithContext(ctx).Preload("Owner").First(&org, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*model.Organization, error) {
	var orgs []*model.Organization
	err := r.db.WithContext(ctx).Preload("Owner").Where("owner_id = ?", ownerID).Order("created_at DESC").Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepository) ListPublic(ctx context.Context, limit, offset int) ([]*model.Organization, int, error) {
	var orgs []*model.Organization
	var total int64

	q := r.db.WithContext(ctx).Preload("Owner")
	q.Model(&model.Organization{}).Count(&total)

	err := q.Limit(limit).Offset(offset).Order("created_at DESC").Find(&orgs).Error
	return orgs, int(total), err
}

func (r *organizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, "id = ?", id).Error
}

func (r *organizationRepository) CountPlates(ctx context.Context, orgID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Plate{}).Where("organization_id = ?", orgID).Count(&count).Error
	return int(count), err
}
