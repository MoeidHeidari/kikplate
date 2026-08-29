package githubapp

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

var (
	ErrNotConfigured            = errors.New("github app is not configured")
	ErrUnauthorized             = errors.New("not authorized to manage github access for this target")
	ErrInvalidState             = errors.New("invalid github installation state")
	ErrInstallationNotFound     = errors.New("github installation not found")
	ErrInstallationTargetType   = errors.New("github installation target type is invalid")
	ErrInstallationLookupFailed = errors.New("failed to read github installation")
	ErrAccessTokenFailed        = errors.New("failed to create github installation token")
	ErrWebhookSignatureInvalid  = errors.New("github webhook signature is invalid")
	ErrWebhookPayloadInvalid    = errors.New("github webhook payload is invalid")
)

const (
	stateTargetAccount      = "account"
	stateTargetOrganization = "organization"
)

type Service interface {
	IsConfigured() bool
	StartAccountConnect(ctx context.Context, accountID uuid.UUID) (string, error)
	StartOrganizationConnect(ctx context.Context, requesterID, organizationID uuid.UUID) (string, error)
	CompleteConnection(ctx context.Context, installationID int64, state string) (string, error)
	GetTokenForOwner(ctx context.Context, ownerAccountID uuid.UUID, organizationID *uuid.UUID) (string, error)
	BuildFrontendRedirect(targetType string, organizationID *uuid.UUID, status string) string
	HandleWebhook(ctx context.Context, event string, payload []byte, signature string) error
	EnsureAccountConnection(ctx context.Context, accountID uuid.UUID) error
	EnsureOrganizationConnection(ctx context.Context, organizationID uuid.UUID) error
}

type service struct {
	env        lib.Env
	logger     lib.Logger
	db         *gorm.DB
	accounts   repository.AccountRepository
	orgs       repository.OrganizationRepository
	orgMembers repository.OrganizationMemberRepository
	httpClient *http.Client
	privateKey []byte
}

type stateClaims struct {
	TargetType     string  `json:"target_type"`
	AccountID      string  `json:"account_id"`
	OrganizationID *string `json:"organization_id,omitempty"`
	jwt.RegisteredClaims
}

type installationResponse struct {
	ID      int64 `json:"id"`
	Account struct {
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"account"`
}

type accessTokenResponse struct {
	Token string `json:"token"`
}

func NewService(
	db lib.Database,
	env lib.Env,
	logger lib.Logger,
	accounts repository.AccountRepository,
	orgs repository.OrganizationRepository,
	orgMembers repository.OrganizationMemberRepository,
) Service {
	svc := &service{
		env:        env,
		logger:     logger,
		db:         db.DB,
		accounts:   accounts,
		orgs:       orgs,
		orgMembers: orgMembers,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}

	if key := loadPrivateKey(env.GitHubApp); len(key) > 0 {
		svc.privateKey = key
	}

	return svc
}

func (s *service) IsConfigured() bool {
	return s.env.GitHubApp.IsConfigured() && len(s.privateKey) > 0 && strings.TrimSpace(s.env.JWTSecret) != ""
}

func (s *service) notConfiguredError() error {
	missing := append([]string{}, s.env.GitHubApp.MissingFields()...)
	if strings.TrimSpace(s.env.JWTSecret) == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) == 0 && len(s.privateKey) == 0 {
		missing = append(missing, "github.app.private_key or github.app.private_key_path")
	}
	if len(missing) == 0 {
		return ErrNotConfigured
	}
	return fmt.Errorf("%w: missing %s", ErrNotConfigured, strings.Join(missing, ", "))
}

func (s *service) StartAccountConnect(ctx context.Context, accountID uuid.UUID) (string, error) {
	if !s.IsConfigured() {
		return "", s.notConfiguredError()
	}

	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil {
		return "", err
	}
	if account == nil {
		return "", ErrUnauthorized
	}

	state, err := s.signStateToken(stateClaims{
		TargetType: stateTargetAccount,
		AccountID:  accountID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})
	if err != nil {
		return "", err
	}

	return appendQueryValue(s.env.GitHubApp.InstallURLValue(), "state", state), nil
}

func (s *service) StartOrganizationConnect(ctx context.Context, requesterID, organizationID uuid.UUID) (string, error) {
	if !s.IsConfigured() {
		return "", s.notConfiguredError()
	}

	org, err := s.orgs.GetByID(ctx, organizationID)
	if err != nil || org == nil {
		return "", ErrUnauthorized
	}
	if !s.canManageOrganization(ctx, org, requesterID) {
		return "", ErrUnauthorized
	}

	orgID := organizationID.String()
	state, err := s.signStateToken(stateClaims{
		TargetType:     stateTargetOrganization,
		AccountID:      requesterID.String(),
		OrganizationID: &orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})
	if err != nil {
		return "", err
	}

	return appendQueryValue(s.env.GitHubApp.InstallURLValue(), "state", state), nil
}

func (s *service) CompleteConnection(ctx context.Context, installationID int64, state string) (string, error) {
	claims, err := s.parseStateToken(state)
	if err != nil {
		return s.BuildFrontendRedirect("", nil, "error"), err
	}

	installation, err := s.getInstallation(ctx, installationID)
	if err != nil {
		return s.BuildFrontendRedirect(claims.TargetType, parseOptionalUUID(claims.OrganizationID), "error"), err
	}

	now := time.Now()
	switch claims.TargetType {
	case stateTargetAccount:
		if !strings.EqualFold(installation.Account.Type, "User") {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "type-mismatch"), ErrInstallationTargetType
		}
		accountID, parseErr := uuid.Parse(claims.AccountID)
		if parseErr != nil {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "error"), ErrInvalidState
		}
		account, getErr := s.accounts.GetByID(ctx, accountID)
		if getErr != nil || account == nil {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "error"), ErrUnauthorized
		}
		account.GitHubInstallationID = &installation.ID
		account.GitHubInstallationAccount = normalizeOptionalString(installation.Account.Login)
		account.GitHubConnectedAt = &now
		if updateErr := s.accounts.Update(ctx, account); updateErr != nil {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "error"), updateErr
		}
		return s.BuildFrontendRedirect(claims.TargetType, nil, "connected"), nil
	case stateTargetOrganization:
		if !strings.EqualFold(installation.Account.Type, "Organization") {
			return s.BuildFrontendRedirect(claims.TargetType, parseOptionalUUID(claims.OrganizationID), "type-mismatch"), ErrInstallationTargetType
		}
		if claims.OrganizationID == nil {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "error"), ErrInvalidState
		}
		orgID, parseErr := uuid.Parse(*claims.OrganizationID)
		if parseErr != nil {
			return s.BuildFrontendRedirect(claims.TargetType, nil, "error"), ErrInvalidState
		}
		org, getErr := s.orgs.GetByID(ctx, orgID)
		if getErr != nil || org == nil {
			return s.BuildFrontendRedirect(claims.TargetType, &orgID, "error"), ErrUnauthorized
		}
		org.GitHubInstallationID = &installation.ID
		org.GitHubInstallationAccount = normalizeOptionalString(installation.Account.Login)
		org.GitHubConnectedAt = &now
		if updateErr := s.orgs.Update(ctx, org); updateErr != nil {
			return s.BuildFrontendRedirect(claims.TargetType, &orgID, "error"), updateErr
		}
		return s.BuildFrontendRedirect(claims.TargetType, &orgID, "connected"), nil
	default:
		return s.BuildFrontendRedirect("", nil, "error"), ErrInvalidState
	}
}

func (s *service) GetTokenForOwner(ctx context.Context, ownerAccountID uuid.UUID, organizationID *uuid.UUID) (string, error) {
	if !s.IsConfigured() {
		return "", nil
	}

	var installationID *int64
	var clearStale func() error
	if organizationID != nil {
		org, err := s.orgs.GetByID(ctx, *organizationID)
		if err != nil || org == nil {
			return "", err
		}
		installationID = org.GitHubInstallationID
		clearStale = func() error {
			org.GitHubInstallationID = nil
			org.GitHubInstallationAccount = nil
			org.GitHubConnectedAt = nil
			return s.orgs.Update(ctx, org)
		}
	} else {
		account, err := s.accounts.GetByID(ctx, ownerAccountID)
		if err != nil || account == nil {
			return "", err
		}
		installationID = account.GitHubInstallationID
		clearStale = func() error {
			account.GitHubInstallationID = nil
			account.GitHubInstallationAccount = nil
			account.GitHubConnectedAt = nil
			return s.accounts.Update(ctx, account)
		}
	}

	if installationID == nil || *installationID <= 0 {
		return "", nil
	}

	token, err := s.createInstallationToken(ctx, *installationID)
	if errors.Is(err, ErrInstallationNotFound) {
		if clearStale != nil {
			if clearErr := clearStale(); clearErr != nil {
				s.logger.Warnf("failed to clear stale github installation_id=%d err=%v", *installationID, clearErr)
			}
		}
		return "", nil
	}
	return token, err
}

func (s *service) BuildFrontendRedirect(targetType string, organizationID *uuid.UUID, status string) string {
	base := strings.TrimRight(s.env.FrontendURL, "/") + "/account"
	query := url.Values{}

	switch targetType {
	case stateTargetOrganization:
		query.Set("tab", "organizations")
		query.Set("github_org_status", status)
		if organizationID != nil {
			query.Set("github_org_id", organizationID.String())
		}
	case stateTargetAccount:
		query.Set("tab", "profile")
		query.Set("github_account_status", status)
	default:
		query.Set("github_account_status", status)
	}

	return base + "?" + query.Encode()
}

func (s *service) EnsureAccountConnection(ctx context.Context, accountID uuid.UUID) error {
	if !s.IsConfigured() {
		return nil
	}
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil || account == nil {
		return err
	}
	if account.GitHubInstallationID == nil || *account.GitHubInstallationID <= 0 {
		return nil
	}

	_, err = s.getInstallation(ctx, *account.GitHubInstallationID)
	if errors.Is(err, ErrInstallationNotFound) {
		account.GitHubInstallationID = nil
		account.GitHubInstallationAccount = nil
		account.GitHubConnectedAt = nil
		return s.accounts.Update(ctx, account)
	}
	if err != nil {
		s.logger.Warnf("ensure account github connection failed account_id=%s installation_id=%d err=%v", accountID, *account.GitHubInstallationID, err)
	}
	return nil
}

func (s *service) EnsureOrganizationConnection(ctx context.Context, organizationID uuid.UUID) error {
	if !s.IsConfigured() {
		return nil
	}
	org, err := s.orgs.GetByID(ctx, organizationID)
	if err != nil || org == nil {
		return err
	}
	if org.GitHubInstallationID == nil || *org.GitHubInstallationID <= 0 {
		return nil
	}

	_, err = s.getInstallation(ctx, *org.GitHubInstallationID)
	if errors.Is(err, ErrInstallationNotFound) {
		org.GitHubInstallationID = nil
		org.GitHubInstallationAccount = nil
		org.GitHubConnectedAt = nil
		return s.orgs.Update(ctx, org)
	}
	if err != nil {
		s.logger.Warnf("ensure org github connection failed org_id=%s installation_id=%d err=%v", organizationID, *org.GitHubInstallationID, err)
	}
	return nil
}

func (s *service) HandleWebhook(ctx context.Context, event string, payload []byte, signature string) error {
	if strings.TrimSpace(s.env.GitHubApp.WebhookSecret) == "" {
		return ErrNotConfigured
	}
	if !verifyWebhookSignature(payload, signature, s.env.GitHubApp.WebhookSecret) {
		return ErrWebhookSignatureInvalid
	}

	s.logger.Infof("github webhook received event=%s", strings.TrimSpace(event))

	if strings.TrimSpace(event) != "installation" {
		return nil
	}

	var body struct {
		Action       string `json:"action"`
		Installation struct {
			ID int64 `json:"id"`
		} `json:"installation"`
	}
	if err := json.Unmarshal(payload, &body); err != nil {
		return ErrWebhookPayloadInvalid
	}
	if body.Installation.ID <= 0 {
		return ErrWebhookPayloadInvalid
	}

	action := strings.ToLower(strings.TrimSpace(body.Action))
	s.logger.Infof("github installation webhook action=%s installation_id=%d", action, body.Installation.ID)
	if action != "deleted" && action != "suspend" {
		return nil
	}

	if s.db == nil {
		return nil
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		accountResult := tx.Model(&model.Account{}).
			Where("git_hub_installation_id = ?", body.Installation.ID).
			Updates(map[string]any{
				"git_hub_installation_id":      gorm.Expr("NULL"),
				"git_hub_installation_account": gorm.Expr("NULL"),
				"git_hub_connected_at":         gorm.Expr("NULL"),
			})
		if accountResult.Error != nil {
			return accountResult.Error
		}

		organizationResult := tx.Model(&model.Organization{}).
			Where("git_hub_installation_id = ?", body.Installation.ID).
			Updates(map[string]any{
				"git_hub_installation_id":      gorm.Expr("NULL"),
				"git_hub_installation_account": gorm.Expr("NULL"),
				"git_hub_connected_at":         gorm.Expr("NULL"),
			})
		if organizationResult.Error != nil {
			return organizationResult.Error
		}

		s.logger.Infof(
			"github uninstall disconnect applied installation_id=%d cleared_accounts=%d cleared_organizations=%d",
			body.Installation.ID,
			accountResult.RowsAffected,
			organizationResult.RowsAffected,
		)

		return nil
	})
}

func (s *service) canManageOrganization(ctx context.Context, org *model.Organization, requesterID uuid.UUID) bool {
	if org.OwnerID == requesterID {
		return true
	}
	member, err := s.orgMembers.GetByOrganizationAndAccount(ctx, org.ID, requesterID)
	if err != nil || member == nil {
		return false
	}
	if member.Status != model.OrganizationMemberStatusAccepted {
		return false
	}
	return member.Role == model.OrganizationMemberRoleOwner || member.Role == model.OrganizationMemberRoleAdmin
}

func (s *service) signStateToken(claims stateClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.env.JWTSecret))
}

func (s *service) parseStateToken(raw string) (*stateClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &stateClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.env.JWTSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidState
	}
	claims, ok := token.Claims.(*stateClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidState
	}
	if claims.TargetType == "" || claims.AccountID == "" {
		return nil, ErrInvalidState
	}
	return claims, nil
}

func (s *service) getInstallation(ctx context.Context, installationID int64) (*installationResponse, error) {
	appJWT, err := s.createAppJWT()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.github.com/app/installations/%d", installationID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("User-Agent", "kickplate-api")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, ErrInstallationLookupFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrInstallationNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		s.logger.Warnf("github installation lookup failed: %d %s", resp.StatusCode, strings.TrimSpace(string(body)))
		return nil, ErrInstallationLookupFailed
	}

	var installation installationResponse
	if err := json.NewDecoder(resp.Body).Decode(&installation); err != nil {
		return nil, ErrInstallationLookupFailed
	}

	return &installation, nil
}

func (s *service) createInstallationToken(ctx context.Context, installationID int64) (string, error) {
	appJWT, err := s.createAppJWT()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID), bytes.NewReader([]byte(`{}`)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+appJWT)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "kickplate-api")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", ErrAccessTokenFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusNotFound {
			return "", ErrInstallationNotFound
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		s.logger.Warnf("github installation token failed: %d %s", resp.StatusCode, strings.TrimSpace(string(body)))
		return "", ErrAccessTokenFailed
	}

	var payload accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", ErrAccessTokenFailed
	}

	if strings.TrimSpace(payload.Token) == "" {
		return "", ErrAccessTokenFailed
	}

	return payload.Token, nil
}

func (s *service) createAppJWT() (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(s.privateKey)
	if err != nil {
		return "", ErrNotConfigured
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now.Add(-1 * time.Minute)),
		ExpiresAt: jwt.NewNumericDate(now.Add(9 * time.Minute)),
		Issuer:    strconv.FormatInt(s.env.GitHubApp.AppID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}

func loadPrivateKey(cfg lib.GitHubAppConfig) []byte {
	if raw := strings.TrimSpace(cfg.PrivateKey); raw != "" {
		return []byte(strings.ReplaceAll(raw, `\n`, "\n"))
	}
	if path := strings.TrimSpace(cfg.PrivateKeyPath); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return data
		}
	}
	return nil
}

func appendQueryValue(rawURL, key, value string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := parsed.Query()
	query.Set(key, value)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func normalizeOptionalString(raw string) *string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func parseOptionalUUID(raw *string) *uuid.UUID {
	if raw == nil {
		return nil
	}
	parsed, err := uuid.Parse(*raw)
	if err != nil {
		return nil
	}
	return &parsed
}

func verifyWebhookSignature(payload []byte, signatureHeader string, secret string) bool {
	signatureHeader = strings.TrimSpace(signatureHeader)
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return false
	}
	receivedHex := strings.TrimPrefix(signatureHeader, "sha256=")
	received, err := hex.DecodeString(receivedHex)
	if err != nil {
		return false
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	expected := h.Sum(nil)

	return hmac.Equal(received, expected)
}
