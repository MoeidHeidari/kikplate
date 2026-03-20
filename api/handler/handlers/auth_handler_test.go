package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── mock auth service ────────────────────────────────────────────────────────

type mockAuthService struct {
	registerFn      func(ctx context.Context, input service.RegisterInput) error
	verifyEmailFn   func(ctx context.Context, token string) (*service.AuthResult, error)
	loginLocalFn    func(ctx context.Context, input service.LoginInput) (*service.AuthResult, error)
	oauthRedirectFn func(ctx context.Context, input service.OAuthRedirectInput) (*service.OAuthRedirectResult, error)
	oauthCallbackFn func(ctx context.Context, input service.OAuthCallbackInput) (*service.AuthResult, error)
	loginHeaderFn   func(ctx context.Context, providerUserID string) (*service.AuthResult, error)
	getMeFn         func(ctx context.Context, accountID uuid.UUID) (*service.MeResult, error)
}

func (m *mockAuthService) Register(ctx context.Context, input service.RegisterInput) error {
	return m.registerFn(ctx, input)
}
func (m *mockAuthService) VerifyEmail(ctx context.Context, token string) (*service.AuthResult, error) {
	return m.verifyEmailFn(ctx, token)
}
func (m *mockAuthService) LoginLocal(ctx context.Context, input service.LoginInput) (*service.AuthResult, error) {
	return m.loginLocalFn(ctx, input)
}
func (m *mockAuthService) OAuthRedirect(ctx context.Context, input service.OAuthRedirectInput) (*service.OAuthRedirectResult, error) {
	return m.oauthRedirectFn(ctx, input)
}
func (m *mockAuthService) OAuthCallback(ctx context.Context, input service.OAuthCallbackInput) (*service.AuthResult, error) {
	return m.oauthCallbackFn(ctx, input)
}
func (m *mockAuthService) LoginHeader(ctx context.Context, providerUserID string) (*service.AuthResult, error) {
	return m.loginHeaderFn(ctx, providerUserID)
}

func (m *mockAuthService) GetMe(ctx context.Context, accountID uuid.UUID) (*service.MeResult, error) {
	return m.getMeFn(ctx, accountID)
}

var _ service.AuthService = (*mockAuthService)(nil)

// ─── helpers ──────────────────────────────────────────────────────────────────

func newHandler(svc service.AuthService) handlers.AuthHandler {
	return handlers.NewAuthHandler(svc, lib.GetLogger())
}

func fakeResult() *service.AuthResult {
	return &service.AuthResult{
		Token: "signed.jwt.token",
		Account: model.Account{
			ID:             uuid.New(),
			Provider:       "local",
			ProviderUserID: uuid.New().String(),
		},
	}
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	return body
}

// routeRequest wires a chi URL param into the request context
// so handlers that call chi.URLParam() work correctly in tests.
func routeRequest(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// ─── Register ─────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(_ context.Context, _ service.RegisterInput) error { return nil },
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, service.RegisterInput{
		Username: "moeidheidari",
		Email:    "moe@example.com",
		Password: "strongpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).Register(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	body := decodeBody(t, rec)
	assert.Contains(t, body["message"], "verify your account")
}

func TestRegister_InvalidBody(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("not-json"))
	rec := httptest.NewRecorder()

	newHandler(svc).Register(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	body := decodeBody(t, rec)
	assert.Equal(t, "invalid request body", body["error"])
}

func TestRegister_EmailTaken(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(_ context.Context, _ service.RegisterInput) error { return service.ErrEmailTaken },
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, service.RegisterInput{
		Username: "moeidheidari", Email: "moe@example.com", Password: "strongpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).Register(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Equal(t, service.ErrEmailTaken.Error(), decodeBody(t, rec)["error"])
}

func TestRegister_UsernameTaken(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(_ context.Context, _ service.RegisterInput) error { return service.ErrUsernameTaken },
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, service.RegisterInput{
		Username: "moeidheidari", Email: "moe@example.com", Password: "strongpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).Register(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Equal(t, service.ErrUsernameTaken.Error(), decodeBody(t, rec)["error"])
}

func TestRegister_UnexpectedError_Returns500(t *testing.T) {
	svc := &mockAuthService{
		registerFn: func(_ context.Context, _ service.RegisterInput) error {
			return errors.New("database exploded")
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, service.RegisterInput{
		Username: "moeidheidari", Email: "moe@example.com", Password: "strongpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).Register(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error", decodeBody(t, rec)["error"])
}

// ─── VerifyEmail ──────────────────────────────────────────────────────────────

func TestVerifyEmail_Success(t *testing.T) {
	result := fakeResult()
	svc := &mockAuthService{
		verifyEmailFn: func(_ context.Context, token string) (*service.AuthResult, error) {
			assert.Equal(t, "valid-raw-token", token)
			return result, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/verify-email?token=valid-raw-token", nil)
	rec := httptest.NewRecorder()

	newHandler(svc).VerifyEmail(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := decodeBody(t, rec)
	assert.Equal(t, result.Token, body["token"])
	assert.NotNil(t, body["account"])
}

func TestVerifyEmail_MissingToken(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodGet, "/auth/verify-email", nil)
	rec := httptest.NewRecorder()

	newHandler(svc).VerifyEmail(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "token is required", decodeBody(t, rec)["error"])
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	svc := &mockAuthService{
		verifyEmailFn: func(_ context.Context, _ string) (*service.AuthResult, error) {
			return nil, service.ErrTokenInvalid
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/verify-email?token=bad-token", nil)
	rec := httptest.NewRecorder()

	newHandler(svc).VerifyEmail(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, service.ErrTokenInvalid.Error(), decodeBody(t, rec)["error"])
}

// ─── LoginLocal ───────────────────────────────────────────────────────────────

func TestLoginLocal_Success(t *testing.T) {
	result := fakeResult()
	svc := &mockAuthService{
		loginLocalFn: func(_ context.Context, input service.LoginInput) (*service.AuthResult, error) {
			assert.Equal(t, "moe@example.com", input.Email)
			return result, nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, service.LoginInput{
		Email: "moe@example.com", Password: "correctpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).LoginLocal(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := decodeBody(t, rec)
	assert.Equal(t, result.Token, body["token"])
	assert.NotNil(t, body["account"])
}

func TestLoginLocal_InvalidBody(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("not-json"))
	rec := httptest.NewRecorder()

	newHandler(svc).LoginLocal(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginLocal_WrongPassword(t *testing.T) {
	svc := &mockAuthService{
		loginLocalFn: func(_ context.Context, _ service.LoginInput) (*service.AuthResult, error) {
			return nil, service.ErrInvalidPassword
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, service.LoginInput{
		Email: "moe@example.com", Password: "wrongpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).LoginLocal(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, service.ErrInvalidPassword.Error(), decodeBody(t, rec)["error"])
}

func TestLoginLocal_InactiveAccount(t *testing.T) {
	svc := &mockAuthService{
		loginLocalFn: func(_ context.Context, _ service.LoginInput) (*service.AuthResult, error) {
			return nil, service.ErrAccountInactive
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, service.LoginInput{
		Email: "moe@example.com", Password: "correctpassword",
	}))
	rec := httptest.NewRecorder()

	newHandler(svc).LoginLocal(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, service.ErrAccountInactive.Error(), decodeBody(t, rec)["error"])
}

// ─── OAuthRedirect ────────────────────────────────────────────────────────────

func TestOAuthRedirect_Success(t *testing.T) {
	svc := &mockAuthService{
		oauthRedirectFn: func(_ context.Context, input service.OAuthRedirectInput) (*service.OAuthRedirectResult, error) {
			assert.Equal(t, "github", input.Provider)
			return &service.OAuthRedirectResult{
				URL:   "https://github.com/login/oauth/authorize?client_id=abc&state=xyz",
				State: "xyz",
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/redirect", nil)
	req = routeRequest(req, "provider", "github")
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthRedirect(rec, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Contains(t, rec.Header().Get("Location"), "github.com")

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "oauth_state", cookies[0].Name)
	assert.Equal(t, "xyz", cookies[0].Value)
}

func TestOAuthRedirect_UnknownProvider(t *testing.T) {
	svc := &mockAuthService{
		oauthRedirectFn: func(_ context.Context, _ service.OAuthRedirectInput) (*service.OAuthRedirectResult, error) {
			return nil, service.ErrProviderNotFound
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/unknown/redirect", nil)
	req = routeRequest(req, "provider", "unknown")
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthRedirect(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, service.ErrProviderNotFound.Error(), decodeBody(t, rec)["error"])
}

// ─── OAuthCallback ────────────────────────────────────────────────────────────

func TestOAuthCallback_Success(t *testing.T) {
	result := fakeResult()
	result.Account.Provider = "github"

	svc := &mockAuthService{
		oauthCallbackFn: func(_ context.Context, input service.OAuthCallbackInput) (*service.AuthResult, error) {
			assert.Equal(t, "github", input.Provider)
			assert.Equal(t, "valid-code", input.Code)
			return result, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?code=valid-code&state=test-state", nil)
	req = routeRequest(req, "provider", "github")
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "test-state"})
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthCallback(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := decodeBody(t, rec)
	assert.Equal(t, result.Token, body["token"])
}

func TestOAuthCallback_MissingCode(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?state=test-state", nil)
	req = routeRequest(req, "provider", "github")
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "test-state"})
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthCallback(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "missing oauth code", decodeBody(t, rec)["error"])
}

func TestOAuthCallback_InvalidState(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?code=valid-code&state=wrong-state", nil)
	req = routeRequest(req, "provider", "github")
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "correct-state"})
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthCallback(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "invalid oauth state", decodeBody(t, rec)["error"])
}

func TestOAuthCallback_MissingStateCookie(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?code=valid-code&state=some-state", nil)
	req = routeRequest(req, "provider", "github")
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthCallback(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "invalid oauth state", decodeBody(t, rec)["error"])
}

func TestOAuthCallback_OAuthFailed(t *testing.T) {
	svc := &mockAuthService{
		oauthCallbackFn: func(_ context.Context, _ service.OAuthCallbackInput) (*service.AuthResult, error) {
			return nil, service.ErrOAuthFailed
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/github/callback?code=bad-code&state=test-state", nil)
	req = routeRequest(req, "provider", "github")
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "test-state"})
	rec := httptest.NewRecorder()

	newHandler(svc).OAuthCallback(rec, req)

	assert.Equal(t, http.StatusBadGateway, rec.Code)
	assert.Equal(t, service.ErrOAuthFailed.Error(), decodeBody(t, rec)["error"])
}

// ─── content-type and error leaking ──────────────────────────────────────────

func TestAllHandlers_AlwaysReturnJSON(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		body    *bytes.Buffer
		handler func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:    "register bad body",
			method:  http.MethodPost,
			body:    bytes.NewBufferString("bad"),
			handler: newHandler(&mockAuthService{}).Register,
		},
		{
			name:    "verify email missing token",
			method:  http.MethodGet,
			body:    nil,
			handler: newHandler(&mockAuthService{}).VerifyEmail,
		},
		{
			name:    "login bad body",
			method:  http.MethodPost,
			body:    bytes.NewBufferString("bad"),
			handler: newHandler(&mockAuthService{}).LoginLocal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := tc.body
			if body == nil {
				body = &bytes.Buffer{}
			}
			req := httptest.NewRequest(tc.method, "/", body)
			rec := httptest.NewRecorder()
			tc.handler(rec, req)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}

func TestInternalErrors_NeverLeakDetails(t *testing.T) {
	internalErr := errors.New("pq: connection refused at 10.0.0.5:5432")

	cases := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		req     *http.Request
	}{
		{
			name: "register",
			handler: newHandler(&mockAuthService{
				registerFn: func(_ context.Context, _ service.RegisterInput) error { return internalErr },
			}).Register,
			req: httptest.NewRequest(http.MethodPost, "/", jsonBody(t, service.RegisterInput{
				Username: "u", Email: "e@e.com", Password: "p",
			})),
		},
		{
			name: "login",
			handler: newHandler(&mockAuthService{
				loginLocalFn: func(_ context.Context, _ service.LoginInput) (*service.AuthResult, error) {
					return nil, internalErr
				},
			}).LoginLocal,
			req: httptest.NewRequest(http.MethodPost, "/", jsonBody(t, service.LoginInput{
				Email: "e@e.com", Password: "p",
			})),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.handler(rec, tc.req)
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, "internal server error", decodeBody(t, rec)["error"])
			assert.NotContains(t, rec.Body.String(), "pq:")
			assert.NotContains(t, rec.Body.String(), "10.0.0.5")
		})
	}
}

func TestMe_Success(t *testing.T) {
	accountID := uuid.New()
	username := "moeidheidari"
	email := "moe@example.com"
	role := "member"
	isActive := true

	svc := &mockAuthService{
		getMeFn: func(_ context.Context, id uuid.UUID) (*service.MeResult, error) {
			assert.Equal(t, accountID, id)
			return &service.MeResult{
				AccountID: accountID.String(),
				Provider:  "local",
				Username:  &username,
				Email:     &email,
				Role:      &role,
				IsActive:  &isActive,
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req = req.WithContext(middleware.SetAccountID(req.Context(), accountID))
	rec := httptest.NewRecorder()

	newHandler(svc).Me(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := decodeBody(t, rec)
	assert.Equal(t, accountID.String(), body["account_id"])
	assert.Equal(t, "moeidheidari", body["username"])
}

func TestMe_NoAccountInContext_Returns401(t *testing.T) {
	svc := &mockAuthService{}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	newHandler(svc).Me(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
