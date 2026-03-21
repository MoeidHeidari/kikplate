package handlers

import (
	"net/http"

	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository"
)

type BadgeHandler struct {
	badges repository.BadgeRepository
	logger lib.Logger
}

func NewBadgeHandler(badges repository.BadgeRepository, logger lib.Logger) BadgeHandler {
	return BadgeHandler{badges: badges, logger: logger}
}

func (h BadgeHandler) List(w http.ResponseWriter, r *http.Request) {
	badges, err := h.badges.List(r.Context())
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, badges)
}
