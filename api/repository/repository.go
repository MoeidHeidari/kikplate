package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
)

type PlateFilter struct {
	Types      []model.PlateType
	Categories []string
	Tags       []string
	OwnerID    *uuid.UUID
	Search     string
	Page       int
	Limit      int
}

type PlateSyncState struct {
	SyncStatus          model.SyncStatus
	SyncError           *string
	LastSyncedAt        *time.Time
	NextSyncAt          *time.Time
	ConsecutiveFailures int
	IsVerified          bool
	VerifiedAt          *time.Time
	Metadata            []byte
}
type PlateStats struct {
	TotalPlates       int64 `json:"total_plates"`
	TotalContributors int64 `json:"total_contributors"`
	TotalCategories   int64 `json:"total_categories"`
	TotalBookmarks    int64 `json:"total_bookmarks"`
}

type PlateFilterOptions struct {
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
}
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	GetByProvider(ctx context.Context, provider, providerUserID string) (*model.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Account, error)
	Update(ctx context.Context, account *model.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EmailVerificationRepository interface {
	Create(ctx context.Context, ev *model.EmailVerification) error
	GetByToken(ctx context.Context, token string) (*model.EmailVerification, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type PlateRepository interface {
	Create(ctx context.Context, plate *model.Plate) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Plate, error)
	GetBySlug(ctx context.Context, slug string) (*model.Plate, error)
	List(ctx context.Context, filter PlateFilter) ([]*model.Plate, int, error)
	Update(ctx context.Context, plate *model.Plate) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementBookmarkCount(ctx context.Context, id uuid.UUID) error
	DecrementBookmarkCount(ctx context.Context, id uuid.UUID) error
	UpdateSyncState(ctx context.Context, id uuid.UUID, state PlateSyncState) error
	ListDueForSync(ctx context.Context, limit int) ([]*model.Plate, error)
	GetStats(ctx context.Context) (*PlateStats, error)
	GetFilterOptions(ctx context.Context) (*PlateFilterOptions, error)
}

type PlateMemberRepository interface {
	Create(ctx context.Context, member *model.PlateMember) error
	GetByPlateAndAccount(ctx context.Context, plateID, accountID uuid.UUID) (*model.PlateMember, error)
	ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateMember, error)
	ListByAccount(ctx context.Context, accountID uuid.UUID) ([]*model.PlateMember, error)
	SetBookmarked(ctx context.Context, plateID, accountID uuid.UUID, bookmarked bool) error
	Delete(ctx context.Context, plateID, accountID uuid.UUID) error
}

type PlateTagRepository interface {
	CreateMany(ctx context.Context, plateID uuid.UUID, tags []string) error
	ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateTag, error)
	DeleteByPlate(ctx context.Context, plateID uuid.UUID) error
}

type PlateReviewRepository interface {
	Create(ctx context.Context, review *model.PlateReview) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.PlateReview, error)
	GetByPlateAndAccount(ctx context.Context, plateID, accountID uuid.UUID) (*model.PlateReview, error)
	ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateReview, error)
	Update(ctx context.Context, review *model.PlateReview) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BadgeRepository interface {
	Create(ctx context.Context, badge *model.Badge) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Badge, error)
	GetBySlug(ctx context.Context, slug string) (*model.Badge, error)
	List(ctx context.Context) ([]*model.Badge, error)
}

type PlateBadgeRepository interface {
	Grant(ctx context.Context, pb *model.PlateBadge) error
	ListByPlate(ctx context.Context, plateID uuid.UUID) ([]*model.PlateBadge, error)
	Revoke(ctx context.Context, plateID, badgeID uuid.UUID) error
}

type OrganizationRepository interface {
	Create(ctx context.Context, org *model.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error)
	GetByName(ctx context.Context, name string) (*model.Organization, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]*model.Organization, error)
	ListPublic(ctx context.Context, limit, offset int) ([]*model.Organization, int, error)
	Update(ctx context.Context, org *model.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountPlates(ctx context.Context, orgID uuid.UUID) (int, error)
}
