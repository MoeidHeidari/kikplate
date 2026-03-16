package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	result := r.db.WithContext(ctx).First(user, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, result.Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	result := r.db.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, result.Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return user, result.Error
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}
