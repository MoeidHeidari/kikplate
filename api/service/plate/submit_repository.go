package plate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
	"gorm.io/gorm"
)

func (s *plateService) SubmitRepository(ctx context.Context, accountID uuid.UUID, input SubmitRepositoryInput) (*model.Plate, error) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
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

	var ownerName string
	if input.OrganizationID != nil {
		if s.orgs == nil {
			return nil, ErrInvalidInput
		}

		org, orgErr := s.orgs.GetByID(ctx, *input.OrganizationID)
		if orgErr != nil || org == nil {
			return nil, ErrInvalidInput
		}
		if org.OwnerID != accountID {
			return nil, ErrForbidden
		}
		ownerName = org.Name
	} else {
		if account.UserID == nil || s.users == nil {
			return nil, ErrNoUsername
		}

		user, err := s.users.GetByID(ctx, *account.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil || strings.TrimSpace(user.Username) == "" {
			return nil, ErrNoUsername
		}
		ownerName = user.Username
	}

	if kp.Owner != ownerName {
		return nil, ErrOwnerMismatch
	}

	metadata, err := json.Marshal(kp)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	now := time.Now()
	verificationToken := uuid.New().String()

	plate := &model.Plate{
		ID:                     uuid.New(),
		OwnerID:                accountID,
		OrganizationID:         input.OrganizationID,
		Type:                   model.PlateTypeRepository,
		Slug:                   slugify(kp.Name),
		Name:                   kp.Name,
		Category:               kp.Category,
		Status:                 model.PlateStatusPending,
		Visibility:             model.PlateVisibilityPrivate,
		Metadata:               metadata,
		RepoURL:                &input.RepoURL,
		Branch:                 &branch,
		IsVerified:             false,
		VerificationToken:      &verificationToken,
		VerificationTokenSetAt: &now,
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
