package plate

import "errors"

var (
	ErrNotFound                  = errors.New("not found")
	ErrForbidden                 = errors.New("forbidden")
	ErrConflict                  = errors.New("conflict")
	ErrInvalidInput              = errors.New("invalid input")
	ErrPrivateRepositoryScope    = errors.New("private repositories are only supported for private organizations")
	ErrAccountGitHubAccess       = errors.New("private personal repositories require account GitHub access")
	ErrOrganizationGitHubAccess  = errors.New("private organization repositories require organization GitHub access")
	ErrNotPendingVerification    = errors.New("plate is not pending verification")
	ErrVerificationTokenMismatch = errors.New("verification token mismatch")
	ErrOrganizationRequired      = errors.New("organization is required")
	ErrNoUsername                = errors.New("repository plates require a username — go to Account settings to set one before submitting")
	ErrOwnerMismatch             = errors.New("owner field does not match the expected owner (username for personal, organization name for org submissions)")
	ErrMissingYAML               = errors.New("manifest not found in repository (expected plate.yaml or kikplate.yaml)")
	ErrFetchFailed               = errors.New("failed to fetch repository")
	ErrAlreadyReviewed           = errors.New("you have already reviewed this plate")
	ErrCannotReviewOwn           = errors.New("you cannot review your own plate")
)
