package plate

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
)

func (s *plateService) VerifyRepository(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) (*model.Plate, error) {
	plate := &model.Plate{}
	if err := s.db.WithContext(ctx).First(plate, "id = ?", plateID).Error; err != nil {
		s.logger.Errorf("plate not found: %v", err)
		return nil, ErrNotFound
	}

	if plate.OwnerID != accountID {
		return nil, ErrForbidden
	}

	if plate.Status != model.PlateStatusPending || plate.VerificationToken == nil {
		return nil, ErrNotPendingVerification
	}

	if plate.RepoURL == nil || plate.Branch == nil {
		return nil, fmt.Errorf("%w: plate missing repository information", ErrInvalidInput)
	}

	expectedToken := strings.ToLower(strings.TrimSpace(*plate.VerificationToken))
	ownerAccountID := plate.OwnerID.String()
	var organizationIDStr *string
	if plate.OrganizationID != nil {
		orgID := plate.OrganizationID.String()
		organizationIDStr = &orgID
	}

	var kp *KickplateYAML
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		kp, err = s.fetchKickplateYAMLWithOptions(ctx, *plate.RepoURL, *plate.Branch, ownerAccountID, organizationIDStr, attempt > 0)
		if err != nil {
			if attempt == 0 && (errors.Is(err, ErrFetchFailed) || errors.Is(err, ErrMissingYAML)) {
				time.Sleep(1200 * time.Millisecond)
				continue
			}
			s.logger.Errorf("fetch kickplate.yaml failed: %v", err)
			return nil, fmt.Errorf("failed to fetch repository: %w", err)
		}

		manifestToken := strings.ToLower(strings.TrimSpace(kp.VerificationToken))
		if manifestToken == expectedToken {
			break
		}

		if attempt == 0 {
			time.Sleep(1200 * time.Millisecond)
			continue
		}

		return nil, fmt.Errorf(
			"%w: expected verification_token %q, found %q in branch %q",
			ErrVerificationTokenMismatch,
			expectedToken,
			manifestToken,
			*plate.Branch,
		)
	}

	now := time.Now()
	syncStatus := model.SyncStatusSynced
	syncIntervalDuration := parseSyncDuration(s.env.SyncInterval)
	syncInterval := syncIntervalDuration.String()
	nextSync := now.Add(syncIntervalDuration)
	visibility := model.PlateVisibilityPublic

	repoPrivate, visibilityErr := s.fetchRepositoryVisibility(ctx, *plate.RepoURL, ownerAccountID, organizationIDStr)
	if visibilityErr != nil {
		return nil, visibilityErr
	}
	if repoPrivate {
		visibility = model.PlateVisibilityPrivate
	}

	if plate.OrganizationID != nil && s.orgs != nil {
		org, orgErr := s.orgs.GetByID(ctx, *plate.OrganizationID)
		if orgErr != nil {
			return nil, orgErr
		}
		if s.env.Features.PrivateOrganizationsEnabled && org != nil && org.Visibility == model.OrganizationVisibilityPrivate {
			visibility = model.PlateVisibilityPrivate
		}
	}

	plate.Status = model.PlateStatusApproved
	plate.Visibility = visibility
	plate.IsVerified = true
	plate.VerifiedAt = &now
	plate.PublishedAt = &now
	plate.SyncStatus = &syncStatus
	plate.SyncInterval = &syncInterval
	plate.NextSyncAt = &nextSync
	plate.LastSyncedAt = &now

	if err := s.db.WithContext(ctx).Save(plate).Error; err != nil {
		s.logger.Errorf("failed to update plate status: %v", err)
		return nil, err
	}

	s.emitPlateVerifiedEvent(ctx, plate)

	s.logger.Infof("plate verified and published: %s", plate.ID)
	return plate, nil
}
