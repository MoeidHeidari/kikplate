package command

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kickplate/api/events"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"github.com/kickplate/api/service/githubapp"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type SyncCommand struct{}

func (c *SyncCommand) Short() string {
	return "Start the background plate synchronizer"
}

func (c *SyncCommand) Setup(_ *cobra.Command) {}

func (c *SyncCommand) Run() lib.CommandRunner {
	return func(
		env lib.Env,
		logger lib.Logger,
		plateRepo repository.PlateRepository,
		plateTagRepo repository.PlateTagRepository,
		accountRepo repository.AccountRepository,
		userRepo repository.UserRepository,
		githubApps githubapp.Service,
		emitter *events.EventEmitter,
	) {
		pollEvery := parseDurationWithFallback(env.SyncPollInterval, 30*time.Second)
		defaultSyncInterval := parseDurationWithFallback(env.SyncInterval, 6*time.Hour)
		batchSize := env.SyncBatchSize
		if batchSize <= 0 {
			batchSize = 25
		}

		logger.Infof("sync worker started (poll=%s, default_interval=%s, batch_size=%d)", pollEvery, defaultSyncInterval, batchSize)

		ticker := time.NewTicker(pollEvery)
		defer ticker.Stop()
		for {
			ctx := context.Background()
			plates, err := plateRepo.ListDueForSync(ctx, batchSize)
			if err != nil {
				logger.Errorf("sync: failed to list due plates: %v", err)
			} else {
				for _, p := range plates {
					if p == nil {
						continue
					}
					syncOnePlate(ctx, logger, plateRepo, plateTagRepo, accountRepo, userRepo, githubApps, emitter, env, p, defaultSyncInterval)
				}
			}

			<-ticker.C
		}
	}
}

func syncOnePlate(
	ctx context.Context,
	logger lib.Logger,
	plateRepo repository.PlateRepository,
	plateTagRepo repository.PlateTagRepository,
	accountRepo repository.AccountRepository,
	userRepo repository.UserRepository,
	githubApps githubapp.Service,
	emitter *events.EventEmitter,
	env lib.Env,
	plate *model.Plate,
	defaultSyncInterval time.Duration,
) {
	now := time.Now()
	interval := defaultSyncInterval
	if plate.SyncInterval != nil && strings.TrimSpace(*plate.SyncInterval) != "" {
		interval = parseDurationWithFallback(*plate.SyncInterval, defaultSyncInterval)
	}
	nextSync := now.Add(interval)
	desiredVisibility := plate.Visibility

	setVisibility := func(v model.PlateVisibility) {
		if plate.Visibility == v {
			return
		}
		plate.Visibility = v
		if err := plateRepo.Update(ctx, plate); err != nil {
			logger.Warnf("sync: failed to update visibility for plate %s: %v", plate.ID, err)
		}
	}

	emitSyncIssue := func(issue string) {
		if emitter == nil {
			return
		}
		account, err := accountRepo.GetByID(ctx, plate.OwnerID)
		if err != nil || account == nil || account.UserID == nil {
			return
		}
		user, err := userRepo.GetByID(ctx, *account.UserID)
		if err != nil || user == nil || strings.TrimSpace(user.Email) == "" {
			return
		}
		emitter.Emit(events.PlateSyncIssue, events.PlateSyncIssuePayload{
			Email:     user.Email,
			PlateName: plate.Name,
			Issue:     issue,
		})
	}

	syncing := model.SyncStatusSyncing
	if err := plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
		SyncStatus:          syncing,
		SyncError:           plate.SyncError,
		LastSyncedAt:        plate.LastSyncedAt,
		NextSyncAt:          plate.NextSyncAt,
		ConsecutiveFailures: plate.ConsecutiveFailures,
		IsVerified:          plate.IsVerified,
		VerifiedAt:          plate.VerifiedAt,
		Metadata:            plate.Metadata,
	}); err != nil {
		logger.Errorf("sync: failed to mark syncing for plate %s: %v", plate.ID, err)
		return
	}

	if plate.RepoURL == nil || plate.Branch == nil {
		errMsg := "plate missing repository information"
		failed := model.SyncStatusFailed
		desiredVisibility = model.PlateVisibilityPrivate
		isVerified := false
		_ = plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
			SyncStatus:          failed,
			SyncError:           &errMsg,
			LastSyncedAt:        &now,
			NextSyncAt:          &nextSync,
			ConsecutiveFailures: plate.ConsecutiveFailures + 1,
			IsVerified:          isVerified,
			VerifiedAt:          nil,
		})
		emitSyncIssue(errMsg)
		logger.Warnf("sync: skipped plate %s (%s)", plate.ID, errMsg)
		return
	}

	githubToken := strings.TrimSpace(env.GitHubToken)
	if githubApps != nil {
		resolvedToken, tokenErr := githubApps.GetTokenForOwner(ctx, plate.OwnerID, plate.OrganizationID)
		if tokenErr != nil {
			logger.Warnf("sync: failed to resolve github installation token for plate %s: %v", plate.ID, tokenErr)
		} else if strings.TrimSpace(resolvedToken) != "" {
			githubToken = resolvedToken
		}
	}

	manifest, _, err := fetchPlateManifestYAML(*plate.RepoURL, *plate.Branch, githubToken)
	if err != nil {
		failed := model.SyncStatusFailed
		errText := err.Error()
		desiredVisibility = model.PlateVisibilityPrivate
		isVerified := false
		_ = plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
			SyncStatus:          failed,
			SyncError:           &errText,
			LastSyncedAt:        &now,
			NextSyncAt:          &nextSync,
			ConsecutiveFailures: plate.ConsecutiveFailures + 1,
			IsVerified:          isVerified,
			VerifiedAt:          nil,
		})
		emitSyncIssue(errText)
		logger.Warnf("sync: plate %s failed: %v", plate.ID, err)
		return
	}
	repoPrivate, repoVisibilityErr := fetchRepositoryPrivate(*plate.RepoURL, githubToken)
	if repoVisibilityErr != nil {
		logger.Warnf("sync: failed to fetch repository visibility for plate %s: %v", plate.ID, repoVisibilityErr)
	}

	isVerified := true
	verifiedAt := plate.VerifiedAt
	syncStatus := model.SyncStatusSynced
	var syncError *string
	privateOrganizationPlate := env.Features.PrivateOrganizationsEnabled &&
		plate.Organization != nil &&
		plate.Organization.Visibility == model.OrganizationVisibilityPrivate

	if plate.VerificationToken == nil || strings.TrimSpace(*plate.VerificationToken) == "" {
		isVerified = false
		verifiedAt = nil
		syncStatus = model.SyncStatusUnverified
		errText := "missing verification token on plate"
		syncError = &errText
		desiredVisibility = model.PlateVisibilityPrivate
	} else {
		expected := strings.ToLower(strings.TrimSpace(*plate.VerificationToken))
		actual := strings.ToLower(strings.TrimSpace(manifest.VerificationToken))
		if expected == "" || actual != expected {
			isVerified = false
			verifiedAt = nil
			syncStatus = model.SyncStatusUnverified
			errText := fmt.Sprintf("verification token mismatch during sync (expected=%q, found=%q)", expected, actual)
			syncError = &errText
			desiredVisibility = model.PlateVisibilityPrivate
		} else {
			isVerified = true
			if verifiedAt == nil {
				verifiedAt = &now
			}
			if privateOrganizationPlate || repoPrivate {
				desiredVisibility = model.PlateVisibilityPrivate
			} else {
				desiredVisibility = model.PlateVisibilityPublic
			}
		}
	}

	metadataJSON, err := json.Marshal(manifest)
	if err != nil {
		failed := model.SyncStatusFailed
		errText := "failed to encode manifest metadata"
		desiredVisibility = model.PlateVisibilityPrivate
		isVerified = false
		verifiedAt = nil
		_ = plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
			SyncStatus:          failed,
			SyncError:           &errText,
			LastSyncedAt:        &now,
			NextSyncAt:          &nextSync,
			ConsecutiveFailures: plate.ConsecutiveFailures + 1,
			IsVerified:          isVerified,
			VerifiedAt:          verifiedAt,
		})
		setVisibility(desiredVisibility)
		emitSyncIssue(errText)
		logger.Warnf("sync: plate %s failed: %s", plate.ID, errText)
		return
	}

	if isVerified {
		if err := syncPlateManifest(ctx, plateRepo, plateTagRepo, env, plate, manifest, metadataJSON); err != nil {
			failed := model.SyncStatusFailed
			errText := fmt.Sprintf("failed to persist manifest changes: %v", err)
			desiredVisibility = model.PlateVisibilityPrivate
			isVerified = false
			verifiedAt = nil
			_ = plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
				SyncStatus:          failed,
				SyncError:           &errText,
				LastSyncedAt:        &now,
				NextSyncAt:          &nextSync,
				ConsecutiveFailures: plate.ConsecutiveFailures + 1,
				IsVerified:          isVerified,
				VerifiedAt:          verifiedAt,
				Metadata:            metadataJSON,
			})
			setVisibility(desiredVisibility)
			emitSyncIssue(errText)
			logger.Warnf("sync: plate %s failed: %s", plate.ID, errText)
			return
		}
	}

	if err := plateRepo.UpdateSyncState(ctx, plate.ID, repository.PlateSyncState{
		SyncStatus:          syncStatus,
		SyncError:           syncError,
		LastSyncedAt:        &now,
		NextSyncAt:          &nextSync,
		ConsecutiveFailures: 0,
		IsVerified:          isVerified,
		VerifiedAt:          verifiedAt,
		Metadata:            metadataJSON,
	}); err != nil {
		logger.Errorf("sync: failed to persist state for plate %s: %v", plate.ID, err)
		return
	}

	setVisibility(desiredVisibility)
	if syncError != nil {
		emitSyncIssue(*syncError)
	}

	logger.Infof("sync: plate %s status=%s next=%s", plate.ID, syncStatus, nextSync.Format(time.RFC3339))
}

type syncPlateManifestYAML struct {
	Owner             string           `yaml:"owner" json:"owner"`
	Name              string           `yaml:"name" json:"name"`
	Description       string           `yaml:"description" json:"description"`
	Category          string           `yaml:"category" json:"category"`
	VerificationToken string           `yaml:"verification_token" json:"verification_token"`
	Tags              []string         `yaml:"tags" json:"tags"`
	Schema            map[string]any   `yaml:"schema" json:"schema"`
	Variables         []map[string]any `yaml:"variables" json:"variables"`
	Dependencies      []map[string]any `yaml:"dependencies" json:"dependencies"`
}

func syncPlateManifest(
	ctx context.Context,
	plateRepo repository.PlateRepository,
	plateTagRepo repository.PlateTagRepository,
	env lib.Env,
	plate *model.Plate,
	manifest *syncPlateManifestYAML,
	metadataJSON []byte,
) error {
	plate.Name = manifest.Name
	plate.Category = lib.NormalizePlateCategory(env, manifest.Category)
	plate.Metadata = metadataJSON

	if strings.TrimSpace(manifest.Description) == "" {
		plate.Description = nil
	} else {
		description := manifest.Description
		plate.Description = &description
	}

	if err := plateRepo.Update(ctx, plate); err != nil {
		return err
	}

	tags := normalizedManifestTags(manifest.Tags)
	if err := plateTagRepo.DeleteByPlate(ctx, plate.ID); err != nil {
		return err
	}
	if len(tags) == 0 {
		return nil
	}

	return plateTagRepo.CreateMany(ctx, plate.ID, tags)
}

func normalizedManifestTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}

	return result
}

func fetchPlateManifestYAML(repoURL, branch string, githubToken string) (*syncPlateManifestYAML, []byte, error) {
	manifest, raw, err := fetchManifestFileYAML(repoURL, branch, "plate.yaml", githubToken)
	if err == nil && strings.TrimSpace(manifest.Owner) != "" {
		return manifest, raw, nil
	}
	if err != nil && !errors.Is(err, ErrSyncMissingManifest) {
		return nil, nil, err
	}

	legacyManifest, legacyRaw, legacyErr := fetchManifestFileYAML(repoURL, branch, "kikplate.yaml", githubToken)
	if legacyErr != nil {
		if err == nil {
			return nil, nil, ErrSyncMissingManifest
		}
		if errors.Is(err, ErrSyncMissingManifest) && errors.Is(legacyErr, ErrSyncMissingManifest) {
			return nil, nil, ErrSyncMissingManifest
		}
		return nil, nil, legacyErr
	}
	if strings.TrimSpace(legacyManifest.Owner) == "" {
		return nil, nil, ErrSyncMissingManifest
	}

	return legacyManifest, legacyRaw, nil
}

var ErrSyncMissingManifest = errors.New("manifest not found")

func fetchManifestFileYAML(repoURL, branch, filename string, githubToken string) (*syncPlateManifestYAML, []byte, error) {
	apiURL := repoURLToContentsURL(repoURL, branch, filename)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch failed")
	}
	if githubToken != "" {
		req.Header.Set("Authorization", "token "+githubToken)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, ErrSyncMissingManifest
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("fetch failed with status %d", resp.StatusCode)
	}

	var ghResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, nil, fmt.Errorf("invalid github response")
	}

	var raw []byte
	if ghResp.Encoding == "base64" {
		raw, err = base64.StdEncoding.DecodeString(ghResp.Content)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid base64 content")
		}
	} else {
		raw = []byte(ghResp.Content)
	}

	var manifest syncPlateManifestYAML
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		return nil, nil, fmt.Errorf("invalid manifest yaml")
	}

	return &manifest, raw, nil
}

func fetchRepositoryPrivate(repoURL string, githubToken string) (bool, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", extractRepoPath(repoURL))
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, err
	}
	if githubToken != "" {
		req.Header.Set("Authorization", "token "+githubToken)
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("visibility fetch failed with status %d", resp.StatusCode)
	}
	var payload struct {
		Private bool `json:"private"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return false, err
	}
	return payload.Private, nil
}

func repoURLToContentsURL(repoURL, branch, filename string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", extractRepoPath(repoURL), filename, branch)
}

func extractRepoPath(repoURL string) string {
	repoURL = strings.TrimSpace(repoURL)
	for _, prefix := range []string{"https://github.com/", "http://github.com/", "github.com/"} {
		if len(repoURL) > len(prefix) && repoURL[:len(prefix)] == prefix {
			repoURL = repoURL[len(prefix):]
			break
		}
	}
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.Trim(repoURL, "/")
	return repoURL
}

func parseDurationWithFallback(value string, fallback time.Duration) time.Duration {
	s := strings.TrimSpace(strings.ToLower(value))
	if s == "" {
		return fallback
	}

	if d, err := time.ParseDuration(s); err == nil && d > 0 {
		return d
	}

	parts := strings.Fields(s)
	if len(parts) == 2 {
		if n, err := strconv.Atoi(parts[0]); err == nil && n > 0 {
			switch parts[1] {
			case "h", "hour", "hours":
				return time.Duration(n) * time.Hour
			case "m", "minute", "minutes":
				return time.Duration(n) * time.Minute
			case "s", "second", "seconds":
				return time.Duration(n) * time.Second
			}
		}
	}

	if hh, mm, ss, ok := parseClockDuration(s); ok {
		return time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute + time.Duration(ss)*time.Second
	}

	return fallback
}

func parseClockDuration(value string) (int, int, int, bool) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	h, errH := strconv.Atoi(parts[0])
	m, errM := strconv.Atoi(parts[1])
	s, errS := strconv.Atoi(parts[2])
	if errH != nil || errM != nil || errS != nil || h < 0 || m < 0 || s < 0 {
		return 0, 0, 0, false
	}
	return h, m, s, true
}

func NewSyncCommand() *SyncCommand {
	return &SyncCommand{}
}
