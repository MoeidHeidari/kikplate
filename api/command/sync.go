package command

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
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
					syncOnePlate(ctx, logger, plateRepo, plateTagRepo, p, defaultSyncInterval)
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
		logger.Warnf("sync: skipped plate %s (%s)", plate.ID, errMsg)
		return
	}

	manifest, _, err := fetchKickplateYAML(*plate.RepoURL, *plate.Branch)
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
		logger.Warnf("sync: plate %s failed: %v", plate.ID, err)
		return
	}

	isVerified := true
	verifiedAt := plate.VerifiedAt
	syncStatus := model.SyncStatusSynced
	var syncError *string

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
			desiredVisibility = model.PlateVisibilityPublic
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
		logger.Warnf("sync: plate %s failed: %s", plate.ID, errText)
		return
	}

	if isVerified {
		if err := syncPlateManifest(ctx, plateRepo, plateTagRepo, plate, manifest, metadataJSON); err != nil {
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

	logger.Infof("sync: plate %s status=%s next=%s", plate.ID, syncStatus, nextSync.Format(time.RFC3339))
}

type syncKickplateYAML struct {
	Owner             string           `yaml:"owner"`
	Name              string           `yaml:"name"`
	Description       string           `yaml:"description"`
	Category          string           `yaml:"category"`
	VerificationToken string           `yaml:"verification_token"`
	Tags              []string         `yaml:"tags"`
	Variables         []map[string]any `yaml:"variables"`
	Dependencies      []map[string]any `yaml:"dependencies"`
}

func syncPlateManifest(
	ctx context.Context,
	plateRepo repository.PlateRepository,
	plateTagRepo repository.PlateTagRepository,
	plate *model.Plate,
	manifest *syncKickplateYAML,
	metadataJSON []byte,
) error {
	plate.Name = manifest.Name
	plate.Category = manifest.Category
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

func fetchKickplateYAML(repoURL, branch string) (*syncKickplateYAML, []byte, error) {
	apiURL := repoURLToContentsURL(repoURL, branch)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, fmt.Errorf("kikplate.yaml not found")
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

	var kp syncKickplateYAML
	if err := yaml.Unmarshal(raw, &kp); err != nil {
		return nil, nil, fmt.Errorf("invalid kickplate.yaml")
	}

	return &kp, raw, nil
}

func repoURLToContentsURL(repoURL, branch string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/contents/kikplate.yaml?ref=%s", extractRepoPath(repoURL), branch)
}

func extractRepoPath(repoURL string) string {
	for _, prefix := range []string{"https://github.com/", "http://github.com/", "github.com/"} {
		if len(repoURL) > len(prefix) && repoURL[:len(prefix)] == prefix {
			return repoURL[len(prefix):]
		}
	}
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
