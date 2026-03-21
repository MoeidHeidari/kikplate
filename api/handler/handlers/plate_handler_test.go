package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	plateservice "github.com/kickplate/api/service/plate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── mock service ─────────────────────────────────────────────────────────────

type mockPlateService struct {
	submitRepositoryFn func(ctx context.Context, accountID uuid.UUID, input plateservice.SubmitRepositoryInput) (*model.Plate, error)
	submitFileFn       func(ctx context.Context, accountID uuid.UUID, input plateservice.SubmitFileInput) (*model.Plate, error)
	getBySlugFn        func(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error)
	listFn             func(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error)
	updateFn           func(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input plateservice.UpdatePlateInput) (*model.Plate, error)
	archiveFn          func(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	recordUseFn        func(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error
	replaceTagsFn      func(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, tags []string) error
	approveFn          func(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	rejectFn           func(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error
	grantBadgeFn       func(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string, reason *string) error
	revokeBadgeFn      func(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string) error
}

func (m *mockPlateService) SubmitFile(ctx context.Context, accountID uuid.UUID, input plateservice.SubmitFileInput) (*model.Plate, error) {
	if m.submitFileFn == nil {
		return nil, nil
	}
	return m.submitFileFn(ctx, accountID, input)
}

func (m *mockPlateService) SubmitRepository(ctx context.Context, accountID uuid.UUID, input plateservice.SubmitRepositoryInput) (*model.Plate, error) {
	if m.submitRepositoryFn == nil {
		return nil, nil
	}
	return m.submitRepositoryFn(ctx, accountID, input)
}

func (m *mockPlateService) GetBySlug(ctx context.Context, slug string, requesterID uuid.UUID) (*model.Plate, error) {
	if m.getBySlugFn == nil {
		return nil, nil
	}
	return m.getBySlugFn(ctx, slug, requesterID)
}

func (m *mockPlateService) List(ctx context.Context, filter repository.PlateFilter, requesterID uuid.UUID) ([]*model.Plate, int, error) {
	if m.listFn == nil {
		return nil, 0, nil
	}
	return m.listFn(ctx, filter, requesterID)
}

func (m *mockPlateService) Update(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, input plateservice.UpdatePlateInput) (*model.Plate, error) {
	if m.updateFn == nil {
		return nil, nil
	}
	return m.updateFn(ctx, plateID, accountID, input)
}

func (m *mockPlateService) Archive(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error {
	if m.archiveFn == nil {
		return nil
	}
	return m.archiveFn(ctx, plateID, accountID)
}

func (m *mockPlateService) RecordUse(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID) error {
	if m.recordUseFn == nil {
		return nil
	}
	return m.recordUseFn(ctx, plateID, accountID)
}

func (m *mockPlateService) ReplaceTags(ctx context.Context, plateID uuid.UUID, accountID uuid.UUID, tags []string) error {
	if m.replaceTagsFn == nil {
		return nil
	}
	return m.replaceTagsFn(ctx, plateID, accountID, tags)
}

func (m *mockPlateService) Approve(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error {
	if m.approveFn == nil {
		return nil
	}
	return m.approveFn(ctx, plateID, adminAccountID)
}

func (m *mockPlateService) Reject(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID) error {
	if m.rejectFn == nil {
		return nil
	}
	return m.rejectFn(ctx, plateID, adminAccountID)
}

func (m *mockPlateService) GrantBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string, reason *string) error {
	if m.grantBadgeFn == nil {
		return nil
	}
	return m.grantBadgeFn(ctx, plateID, adminAccountID, badgeSlug, reason)
}

func (m *mockPlateService) RevokeBadge(ctx context.Context, plateID uuid.UUID, adminAccountID uuid.UUID, badgeSlug string) error {
	if m.revokeBadgeFn == nil {
		return nil
	}
	return m.revokeBadgeFn(ctx, plateID, adminAccountID, badgeSlug)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func newPlateHandler(svc plateservice.PlateService) handlers.PlateHandler {
	return handlers.NewPlateHandler(svc, lib.GetLogger())
}

func withAccountID(r *http.Request, id uuid.UUID) *http.Request {
	ctx := middleware.SetAccountID(r.Context(), id)
	return r.WithContext(ctx)
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx, ok := r.Context().Value(chi.RouteCtxKey).(*chi.Context)
	if !ok || rctx == nil {
		rctx = chi.NewRouteContext()
	}
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func decodeJSON(t *testing.T, body *bytes.Buffer, target any) {
	t.Helper()
	require.NoError(t, json.NewDecoder(body).Decode(target))
}

func stubPlate() *model.Plate {
	slug := "go-dockerfile"
	return &model.Plate{
		ID:         uuid.New(),
		Slug:       slug,
		Name:       "Go Dockerfile",
		Type:       model.PlateTypeFile,
		Status:     model.PlateStatusPending,
		Visibility: model.PlateVisibilityPublic,
		IsVerified: true,
	}
}

// ─── SubmitFile ───────────────────────────────────────────────────────────────

func TestPlateHandler_SubmitFile_Success(t *testing.T) {
	plate := stubPlate()
	svc := &mockPlateService{
		submitFileFn: func(_ context.Context, _ uuid.UUID, _ plateservice.SubmitFileInput) (*model.Plate, error) {
			return plate, nil
		},
	}

	h := newPlateHandler(svc)
	body := jsonBody(t, map[string]any{
		"name": "Go Dockerfile", "category": "devops",
		"visibility": "public", "filename": "Dockerfile",
		"content": "FROM golang:1.22-alpine", "tags": []string{"go"},
	})

	r := httptest.NewRequest(http.MethodPost, "/plates/file", body)
	r = withAccountID(r, uuid.New())
	w := httptest.NewRecorder()

	h.SubmitFile(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp model.Plate
	decodeJSON(t, w.Body, &resp)
	assert.Equal(t, plate.ID, resp.ID)
}

func TestPlateHandler_SubmitFile_InvalidBody(t *testing.T) {
	h := newPlateHandler(&mockPlateService{})

	r := httptest.NewRequest(http.MethodPost, "/plates/file", bytes.NewBufferString("not-json"))
	r = withAccountID(r, uuid.New())
	w := httptest.NewRecorder()

	h.SubmitFile(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlateHandler_SubmitFile_Unauthenticated(t *testing.T) {
	h := newPlateHandler(&mockPlateService{})

	r := httptest.NewRequest(http.MethodPost, "/plates/file", jsonBody(t, map[string]any{}))
	w := httptest.NewRecorder()

	h.SubmitFile(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── SubmitRepository ─────────────────────────────────────────────────────────

func TestPlateHandler_SubmitRepository_OwnerMismatch(t *testing.T) {
	svc := &mockPlateService{
		submitRepositoryFn: func(_ context.Context, _ uuid.UUID, _ plateservice.SubmitRepositoryInput) (*model.Plate, error) {
			return nil, plateservice.ErrOwnerMismatch
		},
	}

	h := newPlateHandler(svc)
	body := jsonBody(t, map[string]any{
		"repo_url": "https://github.com/someone/repo",
		"branch":   "main",
		"username": "wronguser",
	})

	r := httptest.NewRequest(http.MethodPost, "/plates/repository", body)
	r = withAccountID(r, uuid.New())
	w := httptest.NewRecorder()

	h.SubmitRepository(w, r)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	decodeJSON(t, w.Body, &resp)
	assert.Equal(t, plateservice.ErrOwnerMismatch.Error(), resp["error"])
}

func TestPlateHandler_SubmitRepository_Unauthenticated(t *testing.T) {
	h := newPlateHandler(&mockPlateService{})

	r := httptest.NewRequest(http.MethodPost, "/plates/repository", jsonBody(t, map[string]any{}))
	w := httptest.NewRecorder()

	h.SubmitRepository(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── GetBySlug ────────────────────────────────────────────────────────────────

func TestPlateHandler_GetBySlug_Success(t *testing.T) {
	plate := stubPlate()
	svc := &mockPlateService{
		getBySlugFn: func(_ context.Context, slug string, _ uuid.UUID) (*model.Plate, error) {
			return plate, nil
		},
	}

	h := newPlateHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/plates/go-dockerfile", nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "slug", "go-dockerfile")
	w := httptest.NewRecorder()

	h.GetBySlug(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.Plate
	decodeJSON(t, w.Body, &resp)
	assert.Equal(t, plate.ID, resp.ID)
}

func TestPlateHandler_GetBySlug_NotFound(t *testing.T) {
	svc := &mockPlateService{
		getBySlugFn: func(_ context.Context, _ string, _ uuid.UUID) (*model.Plate, error) {
			return nil, plateservice.ErrNotFound
		},
	}

	h := newPlateHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/plates/nope", nil)
	r = withChiParam(r, "slug", "nope")
	w := httptest.NewRecorder()

	h.GetBySlug(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestPlateHandler_List_Success(t *testing.T) {
	plates := []*model.Plate{stubPlate(), stubPlate()}
	svc := &mockPlateService{
		listFn: func(_ context.Context, _ repository.PlateFilter, _ uuid.UUID) ([]*model.Plate, int, error) {
			return plates, 2, nil
		},
	}

	h := newPlateHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/plates", nil)
	r = withAccountID(r, uuid.New())
	w := httptest.NewRecorder()

	h.List(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	decodeJSON(t, w.Body, &resp)
	assert.EqualValues(t, 2, resp["total"])
	assert.Len(t, resp["data"], 2)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestPlateHandler_Update_Success(t *testing.T) {
	plate := stubPlate()
	desc := "Updated"
	plate.Description = &desc

	svc := &mockPlateService{
		updateFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ plateservice.UpdatePlateInput) (*model.Plate, error) {
			return plate, nil
		},
	}

	h := newPlateHandler(svc)
	r := httptest.NewRequest(http.MethodPatch, "/plates/"+plate.ID.String(), jsonBody(t, map[string]any{"description": "Updated"}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plate.ID.String())
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPlateHandler_Update_Forbidden(t *testing.T) {
	svc := &mockPlateService{
		updateFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ plateservice.UpdatePlateInput) (*model.Plate, error) {
			return nil, plateservice.ErrForbidden
		},
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPatch, "/plates/"+plateID.String(), jsonBody(t, map[string]any{}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPlateHandler_Update_InvalidUUID(t *testing.T) {
	h := newPlateHandler(&mockPlateService{})

	r := httptest.NewRequest(http.MethodPatch, "/plates/not-a-uuid", jsonBody(t, map[string]any{}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", "not-a-uuid")
	w := httptest.NewRecorder()

	h.Update(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Archive ──────────────────────────────────────────────────────────────────

func TestPlateHandler_Archive_Success(t *testing.T) {
	svc := &mockPlateService{
		archiveFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodDelete, "/plates/"+plateID.String(), nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.Archive(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPlateHandler_Archive_Forbidden(t *testing.T) {
	svc := &mockPlateService{
		archiveFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
			return plateservice.ErrForbidden
		},
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodDelete, "/plates/"+plateID.String(), nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.Archive(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── RecordUse ────────────────────────────────────────────────────────────────

func TestPlateHandler_RecordUse_Success(t *testing.T) {
	svc := &mockPlateService{
		recordUseFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/use", nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.RecordUse(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPlateHandler_RecordUse_Unauthenticated(t *testing.T) {
	h := newPlateHandler(&mockPlateService{})

	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/use", nil)
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.RecordUse(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── ReplaceTags ─────────────────────────────────────────────────────────────

func TestPlateHandler_ReplaceTags_Success(t *testing.T) {
	svc := &mockPlateService{
		replaceTagsFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ []string) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPut, "/plates/"+plateID.String()+"/tags", jsonBody(t, map[string]any{"tags": []string{"go", "docker"}}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.ReplaceTags(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ─── Approve / Reject ─────────────────────────────────────────────────────────

func TestPlateHandler_Approve_Success(t *testing.T) {
	svc := &mockPlateService{
		approveFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/approve", nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.Approve(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPlateHandler_Reject_Success(t *testing.T) {
	svc := &mockPlateService{
		rejectFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/reject", nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.Reject(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ─── GrantBadge / RevokeBadge ─────────────────────────────────────────────────

func TestPlateHandler_GrantBadge_Success(t *testing.T) {
	svc := &mockPlateService{
		grantBadgeFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, _ *string) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/badges", jsonBody(t, map[string]any{"badge_slug": "featured"}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.GrantBadge(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPlateHandler_GrantBadge_NotFound(t *testing.T) {
	svc := &mockPlateService{
		grantBadgeFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, _ *string) error {
			return plateservice.ErrNotFound
		},
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodPost, "/plates/"+plateID.String()+"/badges", jsonBody(t, map[string]any{"badge_slug": "nope"}))
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	w := httptest.NewRecorder()

	h.GrantBadge(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlateHandler_RevokeBadge_Success(t *testing.T) {
	svc := &mockPlateService{
		revokeBadgeFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string) error { return nil },
	}

	h := newPlateHandler(svc)
	plateID := uuid.New()
	r := httptest.NewRequest(http.MethodDelete, "/plates/"+plateID.String()+"/badges/featured", nil)
	r = withAccountID(r, uuid.New())
	r = withChiParam(r, "id", plateID.String())
	r = withChiParam(r, "slug", "featured")
	w := httptest.NewRecorder()

	h.RevokeBadge(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
