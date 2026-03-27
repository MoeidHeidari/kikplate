package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
	users  repository.UserRepository
	logger lib.Logger
}

func NewPlateHandler(plates plateservice.PlateService, users repository.UserRepository, logger lib.Logger) PlateHandler {
	return PlateHandler{plates: plates, users: users, logger: logger}
}

func (h PlateHandler) SubmitRepository(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		RepoURL        string `json:"repo_url"`
		Branch         string `json:"branch"`
		OrganizationID string `json:"organization_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	var orgID *uuid.UUID
	if body.OrganizationID != "" {
		parsedID, err := uuid.Parse(body.OrganizationID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid organization id")
			return
		}
		orgID = &parsedID
	}

	p, err := h.plates.SubmitRepository(r.Context(), accountID, plateservice.SubmitRepositoryInput{
		RepoURL:        body.RepoURL,
		Branch:         body.Branch,
		OrganizationID: orgID,
	})
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, p)
}

func (h PlateHandler) VerifyRepository(w http.ResponseWriter, r *http.Request) {
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

	p, err := h.plates.VerifyRepository(r.Context(), plateID, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h PlateHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	accountID, _ := middleware.GetAccountID(r.Context())
	slug := chi.URLParam(r, "slug")

	p, err := h.plates.GetBySlug(r.Context(), slug, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	type OwnerInfo struct {
		ID          string  `json:"id"`
		Username    *string `json:"username,omitempty"`
		DisplayName *string `json:"display_name,omitempty"`
		AvatarURL   *string `json:"avatar_url,omitempty"`
	}

	type OrganizationInfo struct {
		ID          string     `json:"id"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		LogoURL     *string    `json:"logo_url,omitempty"`
		OwnerID     string     `json:"owner_id"`
		Owner       *OwnerInfo `json:"owner,omitempty"`
	}

	type PlateResponse struct {
		*model.Plate
		Owner        *OwnerInfo        `json:"owner,omitempty"`
		Organization *OrganizationInfo `json:"organization,omitempty"`
		IsBookmarked bool              `json:"is_bookmarked"`
	}

	resp := PlateResponse{Plate: p, IsBookmarked: false}

	// Check if current user has bookmarked this plate
	if accountID != uuid.Nil {
		member, err := h.plates.GetMember(r.Context(), p.ID, accountID)
		if err == nil && member != nil && member.IsBookmarked {
			resp.IsBookmarked = true
		}
	}

	if p.Owner != nil {
		info := &OwnerInfo{
			ID:          p.Owner.ID.String(),
			DisplayName: p.Owner.DisplayName,
			AvatarURL:   p.Owner.AvatarURL,
		}
		if h.users != nil && p.Owner.UserID != nil {
			user, _ := h.users.GetByID(r.Context(), *p.Owner.UserID)
			if user != nil {
				info.Username = &user.Username
				if user.AvatarURL != nil {
					info.AvatarURL = user.AvatarURL
				}
			}
		}
		resp.Owner = info
	}

	if p.Organization != nil {
		org := &OrganizationInfo{
			ID:          p.Organization.ID.String(),
			Name:        p.Organization.Name,
			Description: p.Organization.Description,
			LogoURL:     p.Organization.LogoURL,
			OwnerID:     p.Organization.OwnerID.String(),
		}

		if p.Organization.Owner != nil {
			ownerInfo := &OwnerInfo{
				ID:          p.Organization.Owner.ID.String(),
				DisplayName: p.Organization.Owner.DisplayName,
				AvatarURL:   p.Organization.Owner.AvatarURL,
			}

			if h.users != nil && p.Organization.Owner.UserID != nil {
				user, _ := h.users.GetByID(r.Context(), *p.Organization.Owner.UserID)
				if user != nil {
					ownerInfo.Username = &user.Username
					if user.AvatarURL != nil {
						ownerInfo.AvatarURL = user.AvatarURL
					}
				}
			}

			org.Owner = ownerInfo
		}

		resp.Organization = org
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h PlateHandler) List(w http.ResponseWriter, r *http.Request) {
	accountID, _ := middleware.GetAccountID(r.Context())
	q := r.URL.Query()

	filter := repository.PlateFilter{
		Search: q.Get("search"),
	}

	for _, key := range []string{"categories", "category"} {
		for _, c := range q[key] {
			for _, part := range splitComma(c) {
				filter.Categories = append(filter.Categories, part)
			}
		}
	}

	for _, key := range []string{"types", "type"} {
		for _, t := range q[key] {
			for _, part := range splitComma(t) {
				filter.Types = append(filter.Types, model.PlateType(part))
			}
		}
	}

	for _, key := range []string{"tags", "tag"} {
		for _, t := range q[key] {
			for _, part := range splitComma(t) {
				filter.Tags = append(filter.Tags, part)
			}
		}
	}

	if ownerID := q.Get("owner_id"); ownerID != "" {
		if id, err := uuid.Parse(ownerID); err == nil {
			filter.OwnerID = &id
		}
	}

	if p := q.Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			filter.Page = n
		}
	}
	if l := q.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			filter.Limit = n
		}
	}

	plates, total, err := h.plates.List(r.Context(), filter, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"data": plates, "total": total})
}

func splitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func (h PlateHandler) ListBookmarked(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := 48
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	plates, err := h.plates.ListBookmarked(r.Context(), accountID, limit)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"data": plates, "total": len(plates)})
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

func (h PlateHandler) MoveToOrganization(w http.ResponseWriter, r *http.Request) {
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
		OrganizationID *string `json:"organization_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	var orgID *uuid.UUID
	if body.OrganizationID != nil {
		trimmed := strings.TrimSpace(*body.OrganizationID)
		if trimmed != "" {
			parsed, err := uuid.Parse(trimmed)
			if err != nil {
				respondError(w, http.StatusBadRequest, "invalid organization id")
				return
			}
			orgID = &parsed
		}
	}

	p, err := h.plates.MoveToOrganization(r.Context(), plateID, accountID, orgID)
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

func (h PlateHandler) Remove(w http.ResponseWriter, r *http.Request) {
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

	if err := h.plates.Remove(r.Context(), plateID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) SetBookmark(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Bookmarked bool `json:"bookmarked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.plates.SetBookmark(r.Context(), plateID, accountID, input.Bookmarked); err != nil {
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PlateHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
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
		Rating int16   `json:"rating"`
		Title  *string `json:"title,omitempty"`
		Body   *string `json:"body,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if body.Rating < 1 || body.Rating > 5 {
		respondError(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}

	review, err := h.plates.SubmitReview(r.Context(), plateID, accountID, plateservice.SubmitReviewInput{
		Rating: body.Rating,
		Title:  body.Title,
		Body:   body.Body,
	})
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, review)
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

func (h PlateHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.plates.GetStats(r.Context())
	if err != nil {
		respondServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

func (h PlateHandler) FilterOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.plates.GetFilterOptions(r.Context())
	if err != nil {
		respondServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, options)
}
