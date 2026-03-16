package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	account := &model.Account{}
	result := r.db.WithContext(ctx).First(account, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return account, result.Error
}

func (r *accountRepository) GetByProvider(ctx context.Context, provider, providerUserID string) (*model.Account, error) {
	account := &model.Account{}
	result := r.db.WithContext(ctx).
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(account)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return account, result.Error
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	account := &model.Account{}
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(account)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return account, result.Error
}

func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *accountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Account{}, "id = ?", id).Error
}
