package plate

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("conflict")
	ErrInvalidInput    = errors.New("invalid input")
	ErrOwnerMismatch   = errors.New("owner field does not match your username")
	ErrMissingYAML     = errors.New("kikplate.yaml not found in repository")
	ErrFetchFailed     = errors.New("failed to fetch repository")
	ErrAlreadyReviewed = errors.New("you have already reviewed this plate")
	ErrCannotReviewOwn = errors.New("you cannot review your own plate")
)
