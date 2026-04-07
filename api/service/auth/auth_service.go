package auth

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"gorm.io/gorm"
)

type authService struct {
	userRepo     repository.UserRepository
	accountRepo  repository.AccountRepository
	emailVerRepo repository.EmailVerificationRepository
	logger       lib.Logger
	env          lib.Env
	db           *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	emailVerRepo repository.EmailVerificationRepository,
	db lib.Database,
	logger lib.Logger,
	env lib.Env,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		accountRepo:  accountRepo,
		emailVerRepo: emailVerRepo,
		logger:       logger,
		env:          env,
		db:           db.DB,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) error {
	existing, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrEmailTaken
	}

	existing, err = s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         model.UserRoleMember,
		IsActive:     !s.env.EmailVerification.Enabled,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	if !s.env.EmailVerification.Enabled {
		return nil
	}

	if !s.env.SMTP.IsConfigured() {
		return ErrSMTPNotConfigured
	}

	ttl := 24 * time.Hour
	if s.env.EmailVerification.TokenTTL != "" {
		parsed, err := time.ParseDuration(s.env.EmailVerification.TokenTTL)
		if err == nil && parsed > 0 {
			ttl = parsed
		}
	}

	rawToken := uuid.New().String()
	hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))

	ev := &model.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     hashed,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := s.emailVerRepo.Create(ctx, ev); err != nil {
		return err
	}

	if err := s.sendVerificationEmail(input.Email, rawToken); err != nil {
		s.logger.Errorf("failed sending verification email to %s: %v", input.Email, err)
		return ErrVerificationEmailFailed
	}

	return nil
}

func (s *authService) sendVerificationEmail(recipientEmail, rawToken string) error {
	verifyURL, err := s.buildVerificationURL(rawToken)
	if err != nil {
		return err
	}

	subject := "Verify your Kikplate account"
	body := fmt.Sprintf("Welcome to Kikplate!\n\nPlease verify your email by opening this link:\n%s\n\nIf you did not request this account, you can ignore this email.", verifyURL)

	from := s.env.SMTP.FromEmail
	if s.env.SMTP.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.env.SMTP.FromName, s.env.SMTP.FromEmail)
	}

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", recipientEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	return sendSMTPMail(s.env.SMTP, s.env.SMTP.FromEmail, []string{recipientEmail}, []byte(message))
}

func (s *authService) buildVerificationURL(rawToken string) (string, error) {
	base := strings.TrimSpace(s.env.EmailVerification.VerifyURLBase)
	if base == "" {
		base = strings.TrimRight(s.env.FrontendURL, "/") + "/verify-email"
	}

	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("token", rawToken)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func sendSMTPMail(cfg lib.SMTPConfig, from string, to []string, message []byte) error {
	if !cfg.IsConfigured() {
		return ErrSMTPNotConfigured
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	tlsConfig := &tls.Config{ServerName: cfg.Host}

	var (
		client *smtp.Client
		err    error
	)

	if cfg.UseStartTL {
		conn, dialErr := net.Dial("tcp", addr)
		if dialErr != nil {
			return dialErr
		}

		client, err = smtp.NewClient(conn, cfg.Host)
		if err != nil {
			_ = conn.Close()
			return err
		}

		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(tlsConfig); err != nil {
				_ = client.Close()
				return err
			}
		}
	} else {
		tlsConn, dialErr := tls.Dial("tcp", addr, tlsConfig)
		if dialErr != nil {
			return dialErr
		}

		client, err = smtp.NewClient(tlsConn, cfg.Host)
		if err != nil {
			_ = tlsConn.Close()
			return err
		}
	}

	defer client.Close()

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	if err := client.Quit(); err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, rawToken string) (*AuthResult, error) {
	hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))

	ev, err := s.emailVerRepo.GetByToken(ctx, hashed)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, ErrTokenInvalid
	}

	user, err := s.userRepo.GetByID(ctx, ev.UserID)
	if err != nil {
		return nil, err
	}
	user.IsActive = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	if err := s.emailVerRepo.MarkUsed(ctx, ev.ID); err != nil {
		return nil, err
	}

	account, err := s.findOrCreateLocalAccount(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.buildResult(ctx, account)
}

func (s *authService) LoginLocal(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidPassword
	}
	if !user.IsActive {
		return nil, ErrAccountInactive
	}
	if user.PasswordHash == "" {
		return nil, ErrInvalidPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	account, err := s.findOrCreateLocalAccount(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.buildResult(ctx, account)
}

func (s *authService) OAuthRedirect(_ context.Context, input OAuthRedirectInput) (*OAuthRedirectResult, error) {
	provider, ok := s.env.GetOAuthProvider(input.Provider)
	if !ok {
		return nil, ErrProviderNotFound
	}

	cfg := s.oauthConfig(provider)
	state := uuid.New().String()
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)

	return &OAuthRedirectResult{URL: url, State: state}, nil
}

func (s *authService) OAuthCallback(ctx context.Context, input OAuthCallbackInput) (*AuthResult, error) {
	provider, ok := s.env.GetOAuthProvider(input.Provider)
	if !ok {
		return nil, ErrProviderNotFound
	}

	cfg := s.oauthConfig(provider)

	token, err := cfg.Exchange(ctx, input.Code)
	if err != nil {
		s.logger.Errorf("oauth token exchange failed for provider %s: %v", input.Provider, err)
		return nil, ErrOAuthFailed
	}

	profile, err := s.fetchProfile(ctx, provider.Name, token)
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.GetByProvider(ctx, provider.Name, profile.ID)
	if err != nil {
		return nil, err
	}

	isNew := account == nil
	if isNew {
		account = &model.Account{
			ID:             uuid.New(),
			Provider:       provider.Name,
			ProviderUserID: profile.ID,
		}
	}

	if profile.Name != "" {
		account.DisplayName = &profile.Name
	}
	if profile.AvatarURL != "" {
		account.AvatarURL = &profile.AvatarURL
	}

	if isNew {
		user, err := s.findOrCreateOAuthUser(ctx, profile)
		if err != nil {
			return nil, err
		}
		account.UserID = &user.ID
		if err := s.accountRepo.Create(ctx, account); err != nil {
			return nil, err
		}
	} else {
		if err := s.accountRepo.Update(ctx, account); err != nil {
			return nil, err
		}
	}

	return s.buildResult(ctx, account)
}

func (s *authService) LoginHeader(ctx context.Context, providerUserID string) (*AuthResult, error) {
	account, err := s.accountRepo.GetByProvider(ctx, "header", providerUserID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		account = &model.Account{
			ID:             uuid.New(),
			Provider:       "header",
			ProviderUserID: providerUserID,
		}
		if err := s.accountRepo.Create(ctx, account); err != nil {
			return nil, err
		}
	}

	return s.buildResult(ctx, account)
}

func (s *authService) SetUsername(ctx context.Context, accountID uuid.UUID, username string) error {
	if username == "" {
		return ErrInvalidUsername
	}

	existing, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUsernameTaken
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrNotFound
	}

	if account.UserID != nil {
		user, err := s.userRepo.GetByID(ctx, *account.UserID)
		if err != nil {
			return err
		}
		if user != nil {
			user.Username = username
			return s.userRepo.Update(ctx, user)
		}
	}

	user := &model.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        "",
		PasswordHash: "",
		Role:         model.UserRoleMember,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	account.UserID = &user.ID
	return s.accountRepo.Update(ctx, account)
}

func (s *authService) UpdateProfile(ctx context.Context, accountID uuid.UUID, input UpdateProfileInput) (*MeResult, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrNotFound
	}

	if input.DisplayName != nil {
		displayName := strings.TrimSpace(*input.DisplayName)
		if displayName == "" {
			account.DisplayName = nil
		} else {
			account.DisplayName = &displayName
		}
	}

	if input.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*input.AvatarURL)
		if avatarURL == "" {
			account.AvatarURL = nil
		} else {
			account.AvatarURL = &avatarURL
		}
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	if account.UserID != nil && input.AvatarURL != nil {
		user, err := s.userRepo.GetByID(ctx, *account.UserID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			user.AvatarURL = account.AvatarURL
			if err := s.userRepo.Update(ctx, user); err != nil {
				return nil, err
			}
		}
	}

	return s.GetMe(ctx, accountID)
}

func (s *authService) findOrCreateOAuthUser(ctx context.Context, profile *oauthProfile) (*model.User, error) {
	if profile.Email != "" {
		existing, err := s.userRepo.GetByEmail(ctx, profile.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	username := profile.Login
	if username == "" {
		username = profile.Email
	}

	taken, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if taken != nil {
		username = fmt.Sprintf("%s_%s", username, uuid.New().String()[:4])
	}

	user := &model.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        profile.Email,
		PasswordHash: "",
		Role:         model.UserRoleMember,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) findOrCreateLocalAccount(ctx context.Context, userID uuid.UUID) (*model.Account, error) {
	account, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if account != nil {
		return account, nil
	}

	account = &model.Account{
		ID:             uuid.New(),
		UserID:         &userID,
		Provider:       "local",
		ProviderUserID: userID.String(),
	}
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *authService) oauthConfig(provider lib.OAuthProvider) *oauth2.Config {
	var endpoint oauth2.Endpoint

	switch provider.Name {
	case "github":
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		}
	case "google":
		endpoint = endpoints.Google
	case "gitlab":
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://gitlab.com/oauth/authorize",
			TokenURL: "https://gitlab.com/oauth/token",
		}
	default:
		endpoint = oauth2.Endpoint{}
	}

	return &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
		Scopes:       provider.Scopes,
		Endpoint:     endpoint,
	}
}

type oauthProfile struct {
	ID        string
	Login     string
	Name      string
	Email     string
	AvatarURL string
}

func (s *authService) fetchProfile(ctx context.Context, providerName string, token *oauth2.Token) (*oauthProfile, error) {
	switch providerName {
	case "github":
		return s.fetchGitHubProfile(ctx, token)
	case "google":
		return s.fetchGoogleProfile(ctx, token)
	case "gitlab":
		return s.fetchGitLabProfile(ctx, token)
	default:
		return nil, ErrProviderNotFound
	}
}

func (s *authService) fetchGitHubProfile(ctx context.Context, token *oauth2.Token) (*oauthProfile, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, ErrOAuthFailed
	}
	defer resp.Body.Close()

	var gh struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gh); err != nil {
		return nil, ErrOAuthFailed
	}

	return &oauthProfile{
		ID:        fmt.Sprintf("%d", gh.ID),
		Login:     gh.Login,
		Name:      gh.Name,
		Email:     gh.Email,
		AvatarURL: gh.AvatarURL,
	}, nil
}

func (s *authService) fetchGoogleProfile(ctx context.Context, token *oauth2.Token) (*oauthProfile, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrOAuthFailed
	}
	defer resp.Body.Close()

	var g struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		return nil, ErrOAuthFailed
	}

	emailPrefix := g.Email
	if idx := len(g.Email); idx > 0 {
		for i, c := range g.Email {
			if c == '@' {
				emailPrefix = g.Email[:i]
				break
			}
		}
	}

	return &oauthProfile{
		ID:        g.ID,
		Login:     emailPrefix,
		Name:      g.Name,
		Email:     g.Email,
		AvatarURL: g.AvatarURL,
	}, nil
}

func (s *authService) fetchGitLabProfile(ctx context.Context, token *oauth2.Token) (*oauthProfile, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://gitlab.com/api/v4/user")
	if err != nil {
		return nil, ErrOAuthFailed
	}
	defer resp.Body.Close()

	var gl struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&gl); err != nil {
		return nil, ErrOAuthFailed
	}

	return &oauthProfile{
		ID:        fmt.Sprintf("%d", gl.ID),
		Login:     gl.Username,
		Name:      gl.Name,
		Email:     gl.Email,
		AvatarURL: gl.AvatarURL,
	}, nil
}

func (s *authService) buildResult(_ context.Context, account *model.Account) (*AuthResult, error) {
	token, err := s.issueToken(account.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, Account: *account}, nil
}

func (s *authService) issueToken(accountID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"account_id": accountID.String(),
		"exp":        time.Now().Add(72 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.env.JWTSecret))
}

func (s *authService) GetMe(ctx context.Context, accountID uuid.UUID) (*MeResult, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrNotFound
	}

	result := &MeResult{
		AccountID:   account.ID.String(),
		Provider:    account.Provider,
		DisplayName: account.DisplayName,
		AvatarURL:   account.AvatarURL,
	}

	if account.UserID != nil {
		user, err := s.userRepo.GetByID(ctx, *account.UserID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			if user.AvatarURL != nil {
				result.AvatarURL = user.AvatarURL
			}
			result.Username = &user.Username
			result.Email = &user.Email
			role := string(user.Role)
			result.Role = &role
			result.IsActive = &user.IsActive
		}
	}

	return result, nil
}

func (s *authService) DeleteMe(ctx context.Context, accountID uuid.UUID) error {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrNotFound
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ownedPlateIDs []uuid.UUID
		if err := tx.Model(&model.Plate{}).
			Where("owner_id = ?", accountID).
			Pluck("id", &ownedPlateIDs).Error; err != nil {
			return err
		}

		if len(ownedPlateIDs) > 0 {
			if err := tx.Where("plate_id IN ?", ownedPlateIDs).Delete(&model.PlateBadge{}).Error; err != nil {
				return err
			}
			if err := tx.Where("plate_id IN ?", ownedPlateIDs).Delete(&model.PlateReview{}).Error; err != nil {
				return err
			}
			if err := tx.Where("plate_id IN ?", ownedPlateIDs).Delete(&model.PlateMember{}).Error; err != nil {
				return err
			}
			if err := tx.Where("plate_id IN ?", ownedPlateIDs).Delete(&model.PlateTag{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN ?", ownedPlateIDs).Delete(&model.Plate{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("owner_id = ?", accountID).Delete(&model.Organization{}).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id = ?", accountID).Delete(&model.PlateReview{}).Error; err != nil {
			return err
		}
		if err := tx.Where("account_id = ?", accountID).Delete(&model.PlateMember{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", accountID).Delete(&model.Account{}).Error; err != nil {
			return err
		}

		if account.UserID != nil {
			var remaining int64
			if err := tx.Model(&model.Account{}).Where("user_id = ?", *account.UserID).Count(&remaining).Error; err != nil {
				return err
			}
			if remaining == 0 {
				if err := tx.Where("user_id = ?", *account.UserID).Delete(&model.EmailVerification{}).Error; err != nil {
					return err
				}
				if err := tx.Where("id = ?", *account.UserID).Delete(&model.User{}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
