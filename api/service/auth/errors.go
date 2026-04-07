package auth

import "errors"

var (
	ErrEmailTaken              = errors.New("email already registered")
	ErrUsernameTaken           = errors.New("username already taken")
	ErrInvalidPassword         = errors.New("invalid credentials")
	ErrAccountInactive         = errors.New("email not verified")
	ErrTokenInvalid            = errors.New("token invalid or expired")
	ErrSMTPNotConfigured       = errors.New("email verification enabled but smtp is not configured")
	ErrVerificationEmailFailed = errors.New("failed to send verification email")
	ErrNotFound                = errors.New("not found")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrProviderNotFound        = errors.New("oauth provider not configured")
	ErrOAuthFailed             = errors.New("oauth authentication failed")
	ErrInvalidUsername         = errors.New("username cannot be empty")
)
