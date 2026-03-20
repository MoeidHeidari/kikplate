package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kickplate/api/service"
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
	case errors.Is(err, service.ErrEmailTaken):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrUsernameTaken):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidPassword):
		respondError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrAccountInactive):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, service.ErrTokenInvalid):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, service.ErrNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrUnauthorized):
		respondError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrProviderNotFound):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrOAuthFailed):
		respondError(w, http.StatusBadGateway, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
