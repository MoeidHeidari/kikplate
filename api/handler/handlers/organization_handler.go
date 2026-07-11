package handlers

import (
	"context"
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
	"github.com/kickplate/api/service/githubapp"
	organizationservice "github.com/kickplate/api/service/organization"
)

type OrganizationHandler struct {
	orgs      organizationservice.OrganizationService
	githubApp githubapp.Service
	accounts  repository.AccountRepository
	users     repository.UserRepository
	logger    lib.Logger
}

func NewOrganizationHandler(orgs organizationservice.OrganizationService, githubApp githubapp.Service, accounts repository.AccountRepository, users repository.UserRepository, logger lib.Logger) OrganizationHandler {
	return OrganizationHandler{orgs: orgs, githubApp: githubApp, accounts: accounts, users: users, logger: logger}
}

type orgOwnerInfo struct {
	ID          string  `json:"id"`
	Username    *string `json:"username,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type orgMemberResponse struct {
	ID        uuid.UUID     `json:"id"`
	AccountID uuid.UUID     `json:"account_id"`
	Role      string        `json:"role"`
	Status    string        `json:"status"`
	JoinedAt  interface{}   `json:"joined_at"`
	Profile   *orgOwnerInfo `json:"profile,omitempty"`
}

type orgResponse struct {
	ID                            uuid.UUID     `json:"id"`
	Name                          string        `json:"name"`
	Visibility                    string        `json:"visibility"`
	Description                   string        `json:"description"`
	LogoURL                       *string       `json:"logo_url,omitempty"`
	OwnerID                       uuid.UUID     `json:"owner_id"`
	GitHubConnected               bool          `json:"github_connected"`
	GitHubInstallationAccountName *string       `json:"github_installation_account,omitempty"`
	MembershipRole                *string       `json:"membership_role,omitempty"`
	Owner                         *orgOwnerInfo `json:"owner,omitempty"`
	CreatedAt                     interface{}   `json:"created_at"`
	UpdatedAt                     interface{}   `json:"updated_at"`
}

func (h OrganizationHandler) enrichOrg(ctx context.Context, org *model.Organization, membershipRole *string) orgResponse {
	resp := orgResponse{
		ID:                            org.ID,
		Name:                          org.Name,
		Visibility:                    org.Visibility,
		Description:                   org.Description,
		LogoURL:                       org.LogoURL,
		OwnerID:                       org.OwnerID,
		GitHubConnected:               org.GitHubInstallationID != nil,
		GitHubInstallationAccountName: org.GitHubInstallationAccount,
		MembershipRole:                membershipRole,
		CreatedAt:                     org.CreatedAt,
		UpdatedAt:                     org.UpdatedAt,
	}
	if org.Owner != nil {
		info := &orgOwnerInfo{
			ID:          org.Owner.ID.String(),
			DisplayName: org.Owner.DisplayName,
			AvatarURL:   org.Owner.AvatarURL,
		}
		if org.Owner.UserID != nil {
			user, _ := h.users.GetByID(ctx, *org.Owner.UserID)
			if user != nil {
				info.Username = &user.Username
				if user.AvatarURL != nil {
					info.AvatarURL = user.AvatarURL
				}
			}
		}
		resp.Owner = info
	}
	return resp
}

func (h OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input organizationservice.CreateOrganizationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	org, err := h.orgs.Create(r.Context(), input, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, org)
}

func (h OrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	accountID, _ := middleware.GetAccountID(r.Context())
	org, err := h.orgs.GetByID(r.Context(), id, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, h.enrichOrg(r.Context(), org, nil))
}

func (h OrganizationHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "invalid organization name")
		return
	}

	accountID, _ := middleware.GetAccountID(r.Context())
	org, err := h.orgs.GetByName(r.Context(), name, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, h.enrichOrg(r.Context(), org, nil))
}

func (h OrganizationHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgs, err := h.orgs.ListByAccount(r.Context(), accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}
	for _, org := range orgs {
		if org == nil {
			continue
		}
		if err := h.githubApp.EnsureOrganizationConnection(r.Context(), org.ID); err != nil {
			h.logger.Warnf("ensure org github connection failed org_id=%s err=%v", org.ID, err)
		}
	}
	orgs, err = h.orgs.ListByAccount(r.Context(), accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	result := make([]orgResponse, len(orgs))
	for i, org := range orgs {
		result[i] = h.enrichOrg(r.Context(), org, org.MembershipRole)
	}

	respondJSON(w, http.StatusOK, result)
}

func (h OrganizationHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	orgs, total, err := h.orgs.ListPublic(r.Context(), limit, offset)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	result := make([]orgResponse, len(orgs))
	for i, org := range orgs {
		result[i] = h.enrichOrg(r.Context(), org, nil)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"organizations": result,
		"total":         total,
		"limit":         limit,
		"offset":        offset,
	})
}

func (h OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	var input organizationservice.UpdateOrganizationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Visibility != nil {
		trim := strings.TrimSpace(*input.Visibility)
		input.Visibility = &trim
	}

	org, err := h.orgs.Update(r.Context(), id, input, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, org)
}

func (h OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	if err := h.orgs.Delete(r.Context(), id, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "organization deleted"})
}

func (h OrganizationHandler) Leave(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	if err := h.orgs.Leave(r.Context(), id, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "left organization"})
}

func (h OrganizationHandler) ConnectGitHub(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	installURL, err := h.githubApp.StartOrganizationConnect(r.Context(), accountID, orgID)
	if err != nil {
		respondServiceError(w, err)
		return
	}
	if state := extractGitHubInstallState(installURL); strings.TrimSpace(state) != "" {
		setGitHubInstallStateCookie(w, state)
	}
	installURL = withGitHubInstallCallback(installURL, absoluteGitHubAppCallbackURL(r))

	http.Redirect(w, r, installURL, http.StatusTemporaryRedirect)
}

func (h OrganizationHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	var input organizationservice.InviteOrganizationMemberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inv, err := h.orgs.InviteMember(r.Context(), orgID, accountID, input)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, inv)
}

func (h OrganizationHandler) ListMyInvitations(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invs, err := h.orgs.ListMyInvitations(r.Context(), accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, invs)
}

func (h OrganizationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invID, err := uuid.Parse(chi.URLParam(r, "invitationId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid invitation id")
		return
	}

	if err := h.orgs.AcceptInvitation(r.Context(), invID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "invitation accepted"})
}

func (h OrganizationHandler) DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	invID, err := uuid.Parse(chi.URLParam(r, "invitationId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid invitation id")
		return
	}

	if err := h.orgs.DeclineInvitation(r.Context(), invID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "invitation declined"})
}

func (h OrganizationHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	members, err := h.orgs.ListMembers(r.Context(), orgID, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	resp := make([]orgMemberResponse, 0, len(members))
	for _, m := range members {
		item := orgMemberResponse{
			ID:        m.ID,
			AccountID: m.AccountID,
			Role:      string(m.Role),
			Status:    string(m.Status),
			JoinedAt:  m.CreatedAt,
		}
		if acct, acctErr := h.orgsAccountByID(r.Context(), m.AccountID); acctErr == nil && acct != nil {
			profile := &orgOwnerInfo{
				ID:          acct.ID.String(),
				DisplayName: acct.DisplayName,
				AvatarURL:   acct.AvatarURL,
			}
			if acct.UserID != nil {
				if user, userErr := h.users.GetByID(r.Context(), *acct.UserID); userErr == nil && user != nil {
					profile.Username = &user.Username
					if user.AvatarURL != nil {
						profile.AvatarURL = user.AvatarURL
					}
				}
			}
			item.Profile = profile
		}
		resp = append(resp, item)
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h OrganizationHandler) ListInvitations(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	invs, err := h.orgs.ListInvitations(r.Context(), orgID, accountID)
	if err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, invs)
}

func (h OrganizationHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	memberAccountID, err := uuid.Parse(chi.URLParam(r, "accountId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	if err := h.orgs.RemoveMember(r.Context(), orgID, memberAccountID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

func (h OrganizationHandler) RevokeInvitation(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orgID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	invitationID, err := uuid.Parse(chi.URLParam(r, "invitationId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid invitation id")
		return
	}

	if err := h.orgs.RevokeInvitation(r.Context(), orgID, invitationID, accountID); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "invitation revoked"})
}

func (h OrganizationHandler) orgsAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	return h.accounts.GetByID(ctx, id)
}
