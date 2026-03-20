package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── mock account repo ────────────────────────────────────────────────────────

type mockAccountRepo struct {
	accounts map[string]*model.Account // keyed by "provider:providerUserID"
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: make(map[string]*model.Account)}
}

func (m *mockAccountRepo) key(provider, providerUserID string) string {
	return provider + ":" + providerUserID
}

func (m *mockAccountRepo) seed(a *model.Account) {
	a.CreatedAt = time.Now()
	m.accounts[m.key(a.Provider, a.ProviderUserID)] = a
}

func (m *mockAccountRepo) Create(_ context.Context, a *model.Account) error {
	a.CreatedAt = time.Now()
	m.accounts[m.key(a.Provider, a.ProviderUserID)] = a
	return nil
}

func (m *mockAccountRepo) GetByID(_ context.Context, id uuid.UUID) (*model.Account, error) {
	for _, a := range m.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, nil
}

func (m *mockAccountRepo) GetByProvider(_ context.Context, provider, providerUserID string) (*model.Account, error) {
	return m.accounts[m.key(provider, providerUserID)], nil
}

func (m *mockAccountRepo) GetByUserID(_ context.Context, userID uuid.UUID) (*model.Account, error) {
	for _, a := range m.accounts {
		if a.UserID != nil && *a.UserID == userID {
			return a, nil
		}
	}
	return nil, nil
}

func (m *mockAccountRepo) Update(_ context.Context, a *model.Account) error {
	m.accounts[m.key(a.Provider, a.ProviderUserID)] = a
	return nil
}

func (m *mockAccountRepo) Delete(_ context.Context, id uuid.UUID) error {
	for k, a := range m.accounts {
		if a.ID == id {
			delete(m.accounts, k)
		}
	}
	return nil
}

var _ repository.AccountRepository = (*mockAccountRepo)(nil)

// ─── helpers ──────────────────────────────────────────────────────────────────

const testSecret = "test-jwt-secret"

func testEnv() lib.Env {
	return lib.Env{JWTSecret: testSecret}
}

func testEnvWithHeader(headerName string) lib.Env {
	return lib.Env{JWTSecret: testSecret, AuthHeader: headerName}
}

func makeToken(accountID uuid.UUID, secret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"account_id": accountID.String(),
		"exp":        time.Now().Add(expiry).Unix(),
		"iat":        time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := middleware.GetAccountID(r.Context())
		if ok {
			w.Header().Set("X-Account-ID", id.String())
		}
		w.WriteHeader(http.StatusOK)
	})
}

// ─── context.go ───────────────────────────────────────────────────────────────

func TestSetAndGetAccountID(t *testing.T) {
	id := uuid.New()
	ctx := middleware.SetAccountID(context.Background(), id)

	got, ok := middleware.GetAccountID(ctx)

	require.True(t, ok)
	assert.Equal(t, id, got)
}

func TestGetAccountID_EmptyContext(t *testing.T) {
	id, ok := middleware.GetAccountID(context.Background())

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetAccountID_WrongType(t *testing.T) {
	type otherKey string
	ctx := context.WithValue(context.Background(), otherKey("account_id"), "not-a-uuid")

	id, ok := middleware.GetAccountID(ctx)

	assert.False(t, ok)
	assert.Equal(t, uuid.Nil, id)
}

// ─── auth.go ──────────────────────────────────────────────────────────────────

func TestAuthenticate_ValidToken_SetsAccountID(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, testSecret, time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, accountID.String(), rec.Header().Get("X-Account-ID"))
}

func TestAuthenticate_NoHeader_PassesThrough(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestAuthenticate_MalformedHeader_PassesThrough(t *testing.T) {
	cases := []string{
		"notbearer token",
		"Bearer",
		"Bearer ",
		"token-without-scheme",
	}

	for _, authHeader := range cases {
		t.Run(authHeader, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", authHeader)
			rec := httptest.NewRecorder()

			middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Empty(t, rec.Header().Get("X-Account-ID"))
		})
	}
}

func TestAuthenticate_ExpiredToken_PassesThrough(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, testSecret, -time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestAuthenticate_WrongSecret_PassesThrough(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, "wrong-secret", time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestAuthenticate_TokenMissingAccountIDClaim_PassesThrough(t *testing.T) {
	claims := jwt.MapClaims{
		"sub": "someone",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	rec := httptest.NewRecorder()

	middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestAuthenticate_BearerCaseInsensitive(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, testSecret, time.Hour)

	for _, scheme := range []string{"Bearer", "bearer", "BEARER"} {
		t.Run(scheme, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", scheme+" "+token)
			rec := httptest.NewRecorder()

			middleware.Authenticate(testEnv(), lib.GetLogger())(okHandler()).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, accountID.String(), rec.Header().Get("X-Account-ID"))
		})
	}
}

// ─── header_auth.go ───────────────────────────────────────────────────────────

func TestHeaderAuth_KnownHeader_SetsAccountID(t *testing.T) {
	accountRepo := newMockAccountRepo()
	accountID := uuid.New()
	accountRepo.seed(&model.Account{
		ID:             accountID,
		Provider:       "header",
		ProviderUserID: "internal-user-abc",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Auth-User", "internal-user-abc")
	rec := httptest.NewRecorder()

	mw := middleware.HeaderAuth(testEnvWithHeader("X-Auth-User"), accountRepo, lib.GetLogger())
	mw(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, accountID.String(), rec.Header().Get("X-Account-ID"))
}

func TestHeaderAuth_HeaderAbsent_PassesThrough(t *testing.T) {
	accountRepo := newMockAccountRepo()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mw := middleware.HeaderAuth(testEnvWithHeader("X-Auth-User"), accountRepo, lib.GetLogger())
	mw(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestHeaderAuth_UnknownHeaderValue_PassesThrough(t *testing.T) {
	accountRepo := newMockAccountRepo()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Auth-User", "unknown-user")
	rec := httptest.NewRecorder()

	mw := middleware.HeaderAuth(testEnvWithHeader("X-Auth-User"), accountRepo, lib.GetLogger())
	mw(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

func TestHeaderAuth_AuthHeaderNotConfigured_PassesThrough(t *testing.T) {
	accountRepo := newMockAccountRepo()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Auth-User", "internal-user-abc")
	rec := httptest.NewRecorder()

	mw := middleware.HeaderAuth(testEnv(), accountRepo, lib.GetLogger())
	mw(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("X-Account-ID"))
}

// ─── require_auth.go ──────────────────────────────────────────────────────────

func TestRequireAuth_WithAccountID_Passes(t *testing.T) {
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetAccountID(req.Context(), accountID))
	rec := httptest.NewRecorder()

	middleware.RequireAuth(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireAuth_NoAccountID_Returns401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.RequireAuth(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, "authentication required", body["error"])
}

func TestRequireAuth_ContentTypeJSON_On401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	middleware.RequireAuth(okHandler()).ServeHTTP(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// ─── integration: Authenticate → RequireAuth chain ───────────────────────────

func TestChain_ValidToken_Passes(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, testSecret, time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler := middleware.Authenticate(testEnv(), lib.GetLogger())(
		middleware.RequireAuth(okHandler()),
	)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, accountID.String(), rec.Header().Get("X-Account-ID"))
}

func TestChain_NoToken_Returns401(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := middleware.Authenticate(testEnv(), lib.GetLogger())(
		middleware.RequireAuth(okHandler()),
	)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestChain_ExpiredToken_Returns401(t *testing.T) {
	accountID := uuid.New()
	token := makeToken(accountID, testSecret, -time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler := middleware.Authenticate(testEnv(), lib.GetLogger())(
		middleware.RequireAuth(okHandler()),
	)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestChain_HeaderAuth_ValidHeader_Passes(t *testing.T) {
	accountRepo := newMockAccountRepo()
	accountID := uuid.New()
	accountRepo.seed(&model.Account{
		ID:             accountID,
		Provider:       "header",
		ProviderUserID: "proxy-user-xyz",
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Auth-User", "proxy-user-xyz")
	rec := httptest.NewRecorder()

	env := testEnvWithHeader("X-Auth-User")
	handler := middleware.HeaderAuth(env, accountRepo, lib.GetLogger())(
		middleware.RequireAuth(okHandler()),
	)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, accountID.String(), rec.Header().Get("X-Account-ID"))
}

func TestChain_HeaderAuth_MissingHeader_Returns401(t *testing.T) {
	accountRepo := newMockAccountRepo()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	env := testEnvWithHeader("X-Auth-User")
	handler := middleware.HeaderAuth(env, accountRepo, lib.GetLogger())(
		middleware.RequireAuth(okHandler()),
	)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
