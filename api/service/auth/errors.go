package auth

import "errors"

var (
	ErrEmailTaken       = errors.New("email already registered")
	ErrUsernameTaken    = errors.New("username already taken")
	ErrInvalidPassword  = errors.New("invalid credentials")
	ErrAccountInactive  = errors.New("email not verified")
	ErrTokenInvalid     = errors.New("token invalid or expired")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrProviderNotFound = errors.New("oauth provider not configured")
	ErrOAuthFailed      = errors.New("oauth authentication failed")
	ErrInvalidUsername  = errors.New("username cannot be empty")
)
