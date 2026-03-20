package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/service"
)

type AuthHandler struct {
	authService service.AuthService
	logger      lib.Logger
}

func NewAuthHandler(
	authService service.AuthService,
	logger lib.Logger,
) AuthHandler {
	return AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.Register(r.Context(), input); err != nil {
		h.logger.Errorf("register failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message": "registration successful, check your email to verify your account",
	})
}

func (h AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}

	result, err := h.authService.VerifyEmail(r.Context(), token)
	if err != nil {
		h.logger.Errorf("verify email failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token":   result.Token,
		"account": result.Account,
	})
}

func (h AuthHandler) LoginLocal(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authService.LoginLocal(r.Context(), input)
	if err != nil {
		h.logger.Errorf("login failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token":   result.Token,
		"account": result.Account,
	})
}

func (h AuthHandler) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	result, err := h.authService.OAuthRedirect(r.Context(), service.OAuthRedirectInput{
		Provider: provider,
	})
	if err != nil {
		h.logger.Errorf("oauth redirect failed: %v", err)
		respondServiceError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    result.State,
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, result.URL, http.StatusTemporaryRedirect)
}

func (h AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		respondError(w, http.StatusBadRequest, "invalid oauth state")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "missing oauth code")
		return
	}

	result, err := h.authService.OAuthCallback(r.Context(), service.OAuthCallbackInput{
		Provider: provider,
		Code:     code,
		State:    stateCookie.Value,
	})
	if err != nil {
		h.logger.Errorf("oauth callback failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"token":   result.Token,
		"account": result.Account,
	})
}

func (h AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	result, err := h.authService.GetMe(r.Context(), accountID)
	if err != nil {
		h.logger.Errorf("get me failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}
