package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kickplate/api/service/auth"
	"github.com/kickplate/api/service/githubapp"
	organizationservice "github.com/kickplate/api/service/organization"
	plateservice "github.com/kickplate/api/service/plate"
)

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func respondServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrEmailTaken):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, auth.ErrUsernameTaken):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, auth.ErrInvalidUsername):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, auth.ErrWeakPassword):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, auth.ErrInvalidPassword):
		respondError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrAccountInactive):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, auth.ErrTokenInvalid):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, auth.ErrSMTPNotConfigured):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, auth.ErrVerificationEmailFailed):
		respondError(w, http.StatusBadGateway, err.Error())
	case errors.Is(err, auth.ErrNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, auth.ErrUnauthorized):
		respondError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrProviderNotFound):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, auth.ErrOAuthFailed):
		respondError(w, http.StatusBadGateway, err.Error())
	case errors.Is(err, plateservice.ErrNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, plateservice.ErrForbidden):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, plateservice.ErrConflict):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, plateservice.ErrAlreadyReviewed):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, plateservice.ErrCannotReviewOwn):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, plateservice.ErrNoUsername):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrOwnerMismatch):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrNotPendingVerification):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, plateservice.ErrVerificationTokenMismatch):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrMissingYAML):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrFetchFailed):
		respondError(w, http.StatusBadGateway, err.Error())
	case errors.Is(err, plateservice.ErrPrivateRepositoryScope):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrAccountGitHubAccess):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrOrganizationGitHubAccess):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, plateservice.ErrInvalidInput):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, plateservice.ErrOrganizationRequired):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, organizationservice.ErrNameTaken):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, organizationservice.ErrNotOwner):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, organizationservice.ErrHasPlates):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, organizationservice.ErrNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, organizationservice.ErrNameRequired):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, organizationservice.ErrInvalidVisibility):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, organizationservice.ErrNotMember):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, organizationservice.ErrOwnerCannotLeave):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, organizationservice.ErrOwnerCannotBeRemoved):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, organizationservice.ErrInviteEmailInvalid):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, organizationservice.ErrInviteNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, organizationservice.ErrInviteExpired):
		respondError(w, http.StatusGone, err.Error())
	case errors.Is(err, organizationservice.ErrInviteForbidden):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, organizationservice.ErrInvalidInputRole):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, githubapp.ErrNotConfigured):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, githubapp.ErrUnauthorized):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, githubapp.ErrInvalidState):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, githubapp.ErrInstallationNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, githubapp.ErrInstallationTargetType):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, githubapp.ErrInstallationLookupFailed), errors.Is(err, githubapp.ErrAccessTokenFailed):
		respondError(w, http.StatusBadGateway, err.Error())
	case errors.Is(err, githubapp.ErrWebhookSignatureInvalid):
		respondError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, githubapp.ErrWebhookPayloadInvalid):
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
