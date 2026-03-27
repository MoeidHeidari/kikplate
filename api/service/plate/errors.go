package plate

import "errors"

var (
	ErrNotFound                  = errors.New("not found")
	ErrForbidden                 = errors.New("forbidden")
	ErrConflict                  = errors.New("conflict")
	ErrInvalidInput              = errors.New("invalid input")
	ErrNotPendingVerification    = errors.New("plate is not pending verification")
	ErrVerificationTokenMismatch = errors.New("verification token mismatch")
	ErrOrganizationRequired      = errors.New("organization is required")
	ErrNoUsername                = errors.New("repository plates require a username — your account was created via SSO and has no username set")
	ErrOwnerMismatch             = errors.New("owner field does not match the expected owner (username for personal, organization name for org submissions)")
	ErrMissingYAML               = errors.New("kikplate.yaml not found in repository")
	ErrFetchFailed               = errors.New("failed to fetch repository")
	ErrAlreadyReviewed           = errors.New("you have already reviewed this plate")
	ErrCannotReviewOwn           = errors.New("you cannot review your own plate")
)
