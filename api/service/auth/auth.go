package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/model"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) error
	VerifyEmail(ctx context.Context, token string) (*AuthResult, error)
	LoginLocal(ctx context.Context, input LoginInput) (*AuthResult, error)
	OAuthRedirect(ctx context.Context, input OAuthRedirectInput) (*OAuthRedirectResult, error)
	OAuthCallback(ctx context.Context, input OAuthCallbackInput) (*AuthResult, error)
	LoginHeader(ctx context.Context, providerUserID string) (*AuthResult, error)
	GetMe(ctx context.Context, accountID uuid.UUID) (*MeResult, error)
	DeleteMe(ctx context.Context, accountID uuid.UUID) error
	SetUsername(ctx context.Context, accountID uuid.UUID, username string) error
	UpdateProfile(ctx context.Context, accountID uuid.UUID, input UpdateProfileInput) (*MeResult, error)
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type UpdateProfileInput struct {
	DisplayName *string
	AvatarURL   *string
}

type MeResult struct {
	AccountID                 string  `json:"account_id"`
	Provider                  string  `json:"provider"`
	DisplayName               *string `json:"display_name,omitempty"`
	AvatarURL                 *string `json:"avatar_url,omitempty"`
	Username                  *string `json:"username,omitempty"`
	Email                     *string `json:"email,omitempty"`
	Role                      *string `json:"role,omitempty"`
	IsActive                  *bool   `json:"is_active,omitempty"`
	GitHubConnected           bool    `json:"github_connected"`
	GitHubInstallationAccount *string `json:"github_installation_account,omitempty"`
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type OAuthRedirectInput struct {
	Provider string
}

type OAuthRedirectResult struct {
	URL   string
	State string
}

type OAuthCallbackInput struct {
	Provider string
	Code     string
	State    string
}

type AuthResult struct {
	Token   string
	Account model.Account
}
