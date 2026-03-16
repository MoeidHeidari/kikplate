package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleMember UserRole = "member"
	UserRoleAdmin  UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"                 json:"id"`
	Username     string    `gorm:"uniqueIndex;size:32;not null"         json:"username"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"        json:"email"`
	PasswordHash string    `gorm:"size:255;not null"                    json:"-"`
	AvatarURL    *string   `gorm:"size:2048"                            json:"avatar_url,omitempty"`
	Role         UserRole  `gorm:"type:varchar(20);default:member"      json:"role"`
	IsActive     bool      `gorm:"default:false"                        json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime"                       json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"                       json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
