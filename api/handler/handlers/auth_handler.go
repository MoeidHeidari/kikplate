package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/service/auth"
	"github.com/kickplate/api/service/githubapp"
)

type AuthHandler struct {
	authService auth.AuthService
	githubApp   githubapp.Service
	logger      lib.Logger
	env         lib.Env
}

const githubInstallStateCookieName = "kp_github_install_state"

func NewAuthHandler(
	authService auth.AuthService,
	githubApp githubapp.Service,
	logger lib.Logger,
	env lib.Env,
) AuthHandler {
	return AuthHandler{
		authService: authService,
		githubApp:   githubApp,
		logger:      logger,
		env:         env,
	}
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input auth.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.Register(r.Context(), input); err != nil {
		h.logger.Errorf("register failed: %v", err)
		respondServiceError(w, err)
		return
	}

	message := "registration successful, you can now login"
	if h.env.EmailVerification.Enabled {
		message = "registration successful, check your email to verify your account"
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message": message,
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
	var input auth.LoginInput
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
	h.logger.Infof("=== OAuthRedirect called for provider: %s ===", provider)
	h.logger.Infof("Raw query: %s", r.URL.RawQuery)

	cliRedirectURI := normalizeCLIRedirectURI(r.URL.Query().Get("redirect_uri"))
	if cliRedirectURI == "" {
		rawCliCallback := r.URL.Query().Get("cli_callback")
		if rawCliCallback != "" {
			preview := rawCliCallback
			if len(preview) > 50 {
				preview = preview[:50]
			}
			h.logger.Infof("received cli_callback: %s", preview)
		}
		cliRedirectURI = normalizeCLIRedirectURI(decodeCLICallback(rawCliCallback))
		if cliRedirectURI != "" {
			h.logger.Infof("decoded and normalized CLI redirect URI: %s", cliRedirectURI)
		} else if rawCliCallback != "" {
			h.logger.Warnf("failed to decode/normalize cli_callback")
		}
	}

	result, err := h.authService.OAuthRedirect(r.Context(), auth.OAuthRedirectInput{
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
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	if cliRedirectURI != "" {
		h.logger.Infof("setting oauth_redirect_uri cookie to: %s", cliRedirectURI)
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_redirect_uri",
			Value:    url.QueryEscape(cliRedirectURI),
			MaxAge:   300,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		h.logger.Warnf("no CLI redirect URI found, using default frontend callback")
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_redirect_uri",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
	}

	http.Redirect(w, r, result.URL, http.StatusTemporaryRedirect)
}

func (h AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "github" && strings.TrimSpace(r.URL.Query().Get("installation_id")) != "" {
		h.GitHubCallback(w, r)
		return
	}

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

	frontendBase := h.env.FrontendURL
	redirectBase := strings.TrimRight(frontendBase, "/") + "/auth/callback"
	if redirectCookie, err := r.Cookie("oauth_redirect_uri"); err == nil {
		h.logger.Infof("found oauth_redirect_uri cookie: %s", redirectCookie.Value)
		if decoded, decErr := url.QueryUnescape(redirectCookie.Value); decErr == nil {
			if normalized := normalizeCLIRedirectURI(decoded); normalized != "" {
				h.logger.Infof("using CLI redirect: %s", normalized)
				redirectBase = normalized
			} else {
				h.logger.Warnf("failed to normalize CLI redirect from cookie: %s", decoded)
			}
		} else {
			h.logger.Warnf("failed to unescape cookie value: %v", decErr)
		}
	} else {
		h.logger.Warnf("no oauth_redirect_uri cookie found, using default: %s", redirectBase)
	}

	result, err := h.authService.OAuthCallback(r.Context(), auth.OAuthCallbackInput{
		Provider: provider,
		Code:     code,
		State:    stateCookie.Value,
	})
	if err != nil {
		h.logger.Errorf("oauth callback failed: %v", err)
		http.Redirect(w, r, appendQueryParam(redirectBase, "error", "oauth_failed"), http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, appendQueryParam(redirectBase, "token", result.Token), http.StatusTemporaryRedirect)
}

func (h AuthHandler) GitHubConnect(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	installURL, err := h.githubApp.StartAccountConnect(r.Context(), accountID)
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

func (h AuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	defer clearGitHubInstallStateCookie(w)

	installationID, err := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("installation_id")), 10, 64)
	if err != nil || installationID <= 0 {
		http.Redirect(w, r, h.githubApp.BuildFrontendRedirect("", nil, "error"), http.StatusTemporaryRedirect)
		return
	}

	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if state == "" {
		state = readGitHubInstallStateCookie(r)
	}

	redirectURL, completeErr := h.githubApp.CompleteConnection(r.Context(), installationID, state)
	if completeErr != nil {
		h.logger.Errorf("github app callback failed: %v", completeErr)
	}
	if strings.TrimSpace(redirectURL) == "" {
		redirectURL = strings.TrimRight(h.env.FrontendURL, "/") + "/account?github_account_status=error"
	}

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h AuthHandler) GitHubWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.githubApp.HandleWebhook(r.Context(), r.Header.Get("X-GitHub-Event"), payload, r.Header.Get("X-Hub-Signature-256")); err != nil {
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func absoluteGitHubAppCallbackURL(r *http.Request) string {
	scheme := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	return fmt.Sprintf("%s://%s/auth/github/app/callback", scheme, r.Host)
}

func withGitHubInstallCallback(installURL string, callbackURL string) string {
	parsed, err := url.Parse(installURL)
	if err != nil {
		return installURL
	}
	query := parsed.Query()
	query.Set("redirect_url", callbackURL)
	query.Set("redirect_uri", callbackURL)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func extractGitHubInstallState(installURL string) string {
	parsed, err := url.Parse(installURL)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(parsed.Query().Get("state"))
}

func setGitHubInstallStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     githubInstallStateCookieName,
		Value:    url.QueryEscape(state),
		MaxAge:   600,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func readGitHubInstallStateCookie(r *http.Request) string {
	cookie, err := r.Cookie(githubInstallStateCookieName)
	if err != nil {
		return ""
	}
	decoded, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(decoded)
}

func clearGitHubInstallStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     githubInstallStateCookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func normalizeCLIRedirectURI(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ""
	}

	host := strings.ToLower(parsed.Hostname())
	if host != "127.0.0.1" && host != "localhost" {
		return ""
	}

	if parsed.Path == "" {
		parsed.Path = "/"
	}

	return parsed.String()
}

func appendQueryParam(base, key, value string) string {
	u, err := url.Parse(base)
	if err != nil {
		if strings.Contains(base, "?") {
			return fmt.Sprintf("%s&%s=%s", base, key, url.QueryEscape(value))
		}
		return fmt.Sprintf("%s?%s=%s", base, key, url.QueryEscape(value))
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String()
}

func decodeCLICallback(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return ""
	}

	return string(decoded)
}

func (h AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	if err := h.githubApp.EnsureAccountConnection(r.Context(), accountID); err != nil {
		h.logger.Warnf("ensure account github connection failed account_id=%s err=%v", accountID, err)
	}

	result, err := h.authService.GetMe(r.Context(), accountID)
	if err != nil {
		h.logger.Errorf("get me failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h AuthHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := h.authService.DeleteMe(r.Context(), accountID); err != nil {
		h.logger.Errorf("delete me failed: %v", err)
		respondServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h AuthHandler) SetUsername(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var body struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.SetUsername(r.Context(), accountID, body.Username); err != nil {
		h.logger.Errorf("set username failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "username set successfully",
	})
}

func (h AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	accountID, ok := middleware.GetAccountID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var body struct {
		DisplayName *string `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.authService.UpdateProfile(r.Context(), accountID, auth.UpdateProfileInput{
		DisplayName: body.DisplayName,
		AvatarURL:   body.AvatarURL,
	})
	if err != nil {
		h.logger.Errorf("update profile failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h AuthHandler) Providers(w http.ResponseWriter, r *http.Request) {
	providers := h.env.OAuthProviders
	names := make([]string, 0, len(providers))
	for _, p := range providers {
		names = append(names, p.Name)
	}
	respondJSON(w, http.StatusOK, map[string]any{
		"providers": names,
	})
}

func (h AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.RequestPasswordReset(r.Context(), input.Email); err != nil {
		h.logger.Errorf("request password reset failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "if an account exists with that email, you will receive a password reset link",
	})
}

func (h AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.ResetPassword(r.Context(), input.Token, input.NewPassword); err != nil {
		h.logger.Errorf("reset password failed: %v", err)
		respondServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "password reset successfully",
	})
}
