package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	plateservice "github.com/kickplate/api/service/plate"
)

type PlateHandler struct {
	plates plateservice.PlateService
	logger lib.Logger
}

func NewPlateHandler(plates plateservice.PlateService, logger lib.Logger) PlateHandler {
	return PlateHandler{plates: plates, logger: logger}
}

func (h PlateHandler) SubmitRepository(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		RepoURL string `json:"repo_url"`
		Branch  string `json:"branch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	p, err := h.plates.SubmitRepository(r.Context(), accountID, plateservice.SubmitRepositoryInput{
		RepoURL: body.RepoURL,
		Branch:  body.Branch,
	})
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, p)
}

func (h PlateHandler) SubmitFile(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Visibility  string   `json:"visibility"`
		Filename    string   `json:"filename"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	p, err := h.plates.SubmitFile(r.Context(), accountID, plateservice.SubmitFileInput{
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		Visibility:  model.PlateVisibility(body.Visibility),
		Filename:    body.Filename,
		Content:     body.Content,
		Tags:        body.Tags,
	})
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, p)
}

func (h PlateHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	accountID, _ := middleware.GetAccountID(r.Context())
	slug := chi.URLParam(r, "slug")

	p, err := h.plates.GetBySlug(r.Context(), slug, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h PlateHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, _ := middleware.GetAccountID(r.Context())
	q := r.URL.Query()

	filter := repository.PlateFilter{
		Category: q.Get("category"),
		Search:   q.Get("search"),
	}
	if tag := q.Get("tag"); tag != "" {
		filter.Tags = []string{tag}
	}

	plates, total, err := h.plates.List(r.Context(), filter, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"data": plates, "total": total})
}

func (h PlateHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		Visibility  *string `json:"visibility"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	input := plateservice.UpdatePlateInput{
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
	}
	if body.Visibility != nil {
		v := model.PlateVisibility(*body.Visibility)
		input.Visibility = &v
	}

	p, err := h.plates.Update(r.Context(), plateID, accountID, input)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h PlateHandler) Archive(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	if err := h.plates.Archive(r.Context(), plateID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) RecordUse(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	if err := h.plates.RecordUse(r.Context(), plateID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) ReplaceTags(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	var body struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.plates.ReplaceTags(r.Context(), plateID, accountID, body.Tags); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) Approve(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	if err := h.plates.Approve(r.Context(), plateID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) Reject(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	if err := h.plates.Reject(r.Context(), plateID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) GrantBadge(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	var body struct {
		BadgeSlug string  `json:"badge_slug"`
		Reason    *string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.plates.GrantBadge(r.Context(), plateID, accountID, body.BadgeSlug, body.Reason); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) RevokeBadge(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	plateID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid plate id")
		return
	}

	badgeSlug := chi.URLParam(r, "slug")

	if err := h.plates.RevokeBadge(r.Context(), plateID, accountID, badgeSlug); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
