package plate

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type KickplateYAML struct {
	Owner        string           `yaml:"owner"`
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Category     string           `yaml:"category"`
	Tags         []string         `yaml:"tags"`
	Variables    []map[string]any `yaml:"variables"`
	Dependencies []map[string]any `yaml:"dependencies"`
}

type plateService struct {
	db           *gorm.DB
	plates       repository.PlateRepository
	tags         repository.PlateTagRepository
	members      repository.PlateMemberRepository
	badges       repository.PlateBadgeRepository
	badgeCatalog repository.BadgeRepository
	accounts     repository.AccountRepository
	users        repository.UserRepository
	logger       lib.Logger
}

func NewPlateService(
	db lib.Database,
	plates repository.PlateRepository,
	tags repository.PlateTagRepository,
	members repository.PlateMemberRepository,
	badges repository.PlateBadgeRepository,
	badgeCatalog repository.BadgeRepository,
	accounts repository.AccountRepository,
	users repository.UserRepository,
	logger lib.Logger,
) PlateService {
	return &plateService{
		db:           db.DB,
		plates:       plates,
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		logger:       logger,
	}
}

func (s *plateService) SubmitRepository(ctx context.Context, accountID uuid.UUID, input SubmitRepositoryInput) (*model.Plate, error) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrNotFound
	}

	if account.UserID == nil {
		return nil, ErrForbidden
	}

	user, err := s.users.GetByID(ctx, *account.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	branch := input.Branch
	if branch == "" {
		branch = "main"
	}

	kp, err := s.fetchKickplateYAML(input.RepoURL, branch)
	if err != nil {
		return nil, err
	}

	if kp.Owner != user.Username {
		return nil, ErrOwnerMismatch
	}

	metadata, err := json.Marshal(kp)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	now := time.Now()
	syncStatus := model.SyncStatusSynced
	syncInterval := "6 hours"
	nextSync := now.Add(6 * time.Hour)

	plate := &model.Plate{
		ID:           uuid.New(),
		OwnerID:      accountID,
		Type:         model.PlateTypeRepository,
		Slug:         slugify(kp.Name),
		Name:         kp.Name,
		Category:     kp.Category,
		Status:       model.PlateStatusApproved,
		PublishedAt:  &now,
		Visibility:   model.PlateVisibilityPublic,
		Metadata:     metadata,
		RepoURL:      &input.RepoURL,
		Branch:       &branch,
		SyncStatus:   &syncStatus,
		SyncInterval: &syncInterval,
		NextSyncAt:   &nextSync,
		LastSyncedAt: &now,
		IsVerified:   true,
		VerifiedAt:   &now,
	}
	if kp.Description != "" {
		plate.Description = &kp.Description
	}

	do := func() error {
		if err := s.plates.Create(ctx, plate); err != nil {
			if strings.Contains(err.Error(), "idx_plate_slug") {
				return ErrConflict
			}
			s.logger.Errorf("plates.Create failed: %v", err)
			return err
		}
		s.logger.Infof("plate created: %s", plate.ID)

		if len(kp.Tags) > 0 {
			if err := s.tags.CreateMany(ctx, plate.ID, kp.Tags); err != nil {
				s.logger.Errorf("tags.CreateMany failed: %v", err)
				return err
			}
			s.logger.Infof("tags created for plate: %s", plate.ID)
		}

		if err := s.members.Create(ctx, &model.PlateMember{
			ID:        uuid.New(),
			PlateID:   plate.ID,
			AccountID: accountID,
			Role:      model.PlateMemberRoleOwner,
		}); err != nil {
			s.logger.Errorf("members.Create failed: %v", err)
			return err
		}
		s.logger.Infof("owner member created for plate: %s", plate.ID)

		return nil
	}

	if s.db != nil {
		if err := s.db.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
			return do()
		}); err != nil {
			s.logger.Errorf("transaction failed: %v", err)
			return nil, err
		}
	} else {
		if err := do(); err != nil {
			return nil, err
		}
	}

	return plate, nil
}

func (s *plateService) SubmitFile(ctx context.Context, accountID uuid.UUID, input SubmitFileInput) (*model.Plate, error) {
	now := time.Now()

	metadata, _ := json.Marshal(map[string]any{
		"file_type": strings.ToLower(input.Filename),
	})

	plate := &model.Plate{
		ID:         uuid.New(),
		OwnerID:    accountID,
		Type:       model.PlateTypeFile,
		Slug:       slugify(input.Name),
		Name:       input.Name,
		Category:   input.Category,
		Status:     model.PlateStatusPending,
		Visibility: input.Visibility,
		Metadata:   metadata,
		Content:    &input.Content,
		Filename:   &input.Filename,
		IsVerified: true,
		VerifiedAt: &now,
	}
	if input.Description != "" {
		plate.Description = &input.Description
	}

	do := func() error {
		if err := s.plates.Create(ctx, plate); err != nil {
			return err
		}
		if len(input.Tags) > 0 {
			if err := s.tags.CreateMany(ctx, plate.ID, input.Tags); err != nil {
				return err
			}
		}
		return s.members.Create(ctx, &model.PlateMember{
			ID:        uuid.New(),
			PlateID:   plate.ID,
			AccountID: accountID,
			Role:      model.PlateMemberRoleOwner,
		})
	}

	if s.db != nil {
		if err := s.db.WithContext(ctx).Transaction(func(_ *gorm.DB) error {
			return do()
		}); err != nil {
			return nil, err
		}
	} else {
		if err := do(); err != nil {
			return nil, err
		}
	}

	return plate, nil
}

func (s *plateService) GetBySlug(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error) {
	plate, err := s.plates.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if plate == nil {
		return nil, ErrNotFound
	}

	if plate.Visibility == model.PlateVisibilityPrivate {
		member, err := s.members.GetByPlateAndAccount(ctx, plate.ID, requesterID)
		if err != nil {
			return nil, err
		}
		if member == nil {
			return nil, ErrNotFound
		}
	}

	return plate, nil
}

func (s *plateService) List(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error) {
	if filter.Visibility == nil {
		pub := model.PlateVisibilityPublic
		filter.Visibility = &pub
	}
	return s.plates.List(ctx, filter)
}

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

func (s *plateService) RecordUse(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error {
	now := time.Now()

	member, err := s.members.GetByPlateAndAccount(ctx, plateID, accountID)
	if err != nil {
		return err
	}

	if member == nil {
		if err := s.members.Create(ctx, &model.PlateMember{
			ID:         uuid.New(),
			PlateID:    plateID,
			AccountID:  accountID,
			Role:       model.PlateMemberRoleMember,
			LastUsedAt: &now,
		}); err != nil {
			return err
		}
	} else {
		if err := s.members.UpdateLastUsedAt(ctx, plateID, accountID, now); err != nil {
			return err
		}
	}

	return s.plates.IncrementUseCount(ctx, plateID)
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

func (s *plateService) fetchKickplateYAML(repoURL, branch string) (*KickplateYAML, error) {
	apiURL := repoURLToContentsURL(repoURL, branch)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, ErrFetchFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrMissingYAML
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrFetchFailed
	}

	var ghResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, ErrFetchFailed
	}

	var raw []byte
	if ghResp.Encoding == "base64" {
		raw, err = base64.StdEncoding.DecodeString(ghResp.Content)
		if err != nil {
			return nil, ErrFetchFailed
		}
	} else {
		raw = []byte(ghResp.Content)
	}

	var kp KickplateYAML
	if err := yaml.Unmarshal(raw, &kp); err != nil {
		return nil, ErrFetchFailed
	}
	if kp.Owner == "" {
		return nil, ErrMissingYAML
	}

	return &kp, nil
}

func repoURLToContentsURL(repoURL, branch string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/kikplate.yaml?ref=%s",
		extractRepoPath(repoURL), branch)
}

func extractRepoPath(repoURL string) string {
	for _, prefix := range []string{
		"https://github.com/",
		"http://github.com/",
		"github.com/",
	} {
		if len(repoURL) > len(prefix) && repoURL[:len(prefix)] == prefix {
			return repoURL[len(prefix):]
		}
	}
	return repoURL
}

func slugify(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
func NewPlateServiceForTest(
	plates repository.PlateRepository,
	tags repository.PlateTagRepository,
	members repository.PlateMemberRepository,
	badges repository.PlateBadgeRepository,
	badgeCatalog repository.BadgeRepository,
	accounts repository.AccountRepository,
	users repository.UserRepository,
	logger lib.Logger,
) PlateService {
	return &plateService{
		db:           nil,
		plates:       plates,
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		logger:       logger,
	}
}
