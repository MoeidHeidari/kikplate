package plate_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	plateservice "github.com/kickplate/api/service/plate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── mocks ────────────────────────────────────────────────────────────────────

type mockPlateRepo struct {
	plates map[uuid.UUID]*model.Plate
}

func newMockPlateRepo() *mockPlateRepo {
	return &mockPlateRepo{plates: make(map[uuid.UUID]*model.Plate)}
}

func (m *mockPlateRepo) Create(_ context.Context, p *model.Plate) error {
	m.plates[p.ID] = p
	return nil
}

func (m *mockPlateRepo) GetByID(_ context.Context, id uuid.UUID) (*model.Plate, error) {
	p, ok := m.plates[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockPlateRepo) GetBySlug(_ context.Context, slug string) (*model.Plate, error) {
	for _, p := range m.plates {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockPlateRepo) List(_ context.Context, filter repository.PlateFilter) ([]*model.Plate, int, error) {
	var result []*model.Plate
	for _, p := range m.plates {
		if filter.Visibility != nil && p.Visibility != *filter.Visibility {
			continue
		}
		if filter.Category != "" && p.Category != filter.Category {
			continue
		}
		result = append(result, p)
	}
	return result, len(result), nil
}

func (m *mockPlateRepo) Update(_ context.Context, p *model.Plate) error {
	m.plates[p.ID] = p
	return nil
}

func (m *mockPlateRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.plates, id)
	return nil
}

func (m *mockPlateRepo) IncrementUseCount(_ context.Context, id uuid.UUID) error {
	if p, ok := m.plates[id]; ok {
		p.UseCount++
	}
	return nil
}

func (m *mockPlateRepo) UpdateSyncState(_ context.Context, id uuid.UUID, state repository.PlateSyncState) error {
	return nil
}

func (m *mockPlateRepo) ListDueForSync(_ context.Context, limit int) ([]*model.Plate, error) {
	return nil, nil
}

// ─────────────────────────────────────────────────────────────────────────────

type mockTagRepo struct {
	tags map[uuid.UUID][]string
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{tags: make(map[uuid.UUID][]string)}
}

func (m *mockTagRepo) CreateMany(_ context.Context, plateID uuid.UUID, tags []string) error {
	m.tags[plateID] = append(m.tags[plateID], tags...)
	return nil
}

func (m *mockTagRepo) ListByPlate(_ context.Context, plateID uuid.UUID) ([]*model.PlateTag, error) {
	var result []*model.PlateTag
	for _, tag := range m.tags[plateID] {
		result = append(result, &model.PlateTag{PlateID: plateID, Tag: tag})
	}
	return result, nil
}

func (m *mockTagRepo) DeleteByPlate(_ context.Context, plateID uuid.UUID) error {
	m.tags[plateID] = nil
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────

type mockMemberRepo struct {
	members map[string]*model.PlateMember
}

func newMockMemberRepo() *mockMemberRepo {
	return &mockMemberRepo{members: make(map[string]*model.PlateMember)}
}

func (m *mockMemberRepo) key(plateID, accountID uuid.UUID) string {
	return plateID.String() + ":" + accountID.String()
}

func (m *mockMemberRepo) Create(_ context.Context, member *model.PlateMember) error {
	m.members[m.key(member.PlateID, member.AccountID)] = member
	return nil
}

func (m *mockMemberRepo) GetByPlateAndAccount(_ context.Context, plateID, accountID uuid.UUID) (*model.PlateMember, error) {
	return m.members[m.key(plateID, accountID)], nil
}

func (m *mockMemberRepo) ListByPlate(_ context.Context, plateID uuid.UUID) ([]*model.PlateMember, error) {
	var result []*model.PlateMember
	for _, mem := range m.members {
		if mem.PlateID == plateID {
			result = append(result, mem)
		}
	}
	return result, nil
}

func (m *mockMemberRepo) ListByAccount(_ context.Context, accountID uuid.UUID) ([]*model.PlateMember, error) {
	var result []*model.PlateMember
	for _, mem := range m.members {
		if mem.AccountID == accountID {
			result = append(result, mem)
		}
	}
	return result, nil
}

func (m *mockMemberRepo) UpdateLastUsedAt(_ context.Context, plateID, accountID uuid.UUID, t time.Time) error {
	k := m.key(plateID, accountID)
	if mem, ok := m.members[k]; ok {
		mem.LastUsedAt = &t
	}
	return nil
}

func (m *mockMemberRepo) Delete(_ context.Context, plateID, accountID uuid.UUID) error {
	delete(m.members, m.key(plateID, accountID))
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────

type mockBadgeRepo struct {
	grants map[string]*model.PlateBadge
}

func newMockBadgeRepo() *mockBadgeRepo {
	return &mockBadgeRepo{grants: make(map[string]*model.PlateBadge)}
}

func (m *mockBadgeRepo) key(plateID, badgeID uuid.UUID) string {
	return plateID.String() + ":" + badgeID.String()
}

func (m *mockBadgeRepo) Grant(_ context.Context, pb *model.PlateBadge) error {
	m.grants[m.key(pb.PlateID, pb.BadgeID)] = pb
	return nil
}

func (m *mockBadgeRepo) ListByPlate(_ context.Context, plateID uuid.UUID) ([]*model.PlateBadge, error) {
	var result []*model.PlateBadge
	for _, pb := range m.grants {
		if pb.PlateID == plateID {
			result = append(result, pb)
		}
	}
	return result, nil
}

func (m *mockBadgeRepo) Revoke(_ context.Context, plateID, badgeID uuid.UUID) error {
	delete(m.grants, m.key(plateID, badgeID))
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────

type mockBadgeCatalog struct {
	badges map[string]*model.Badge
}

func newMockBadgeCatalog() *mockBadgeCatalog {
	return &mockBadgeCatalog{badges: make(map[string]*model.Badge)}
}

func (m *mockBadgeCatalog) Create(_ context.Context, b *model.Badge) error {
	m.badges[b.Slug] = b
	return nil
}

func (m *mockBadgeCatalog) GetByID(_ context.Context, id uuid.UUID) (*model.Badge, error) {
	for _, b := range m.badges {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, nil
}

func (m *mockBadgeCatalog) GetBySlug(_ context.Context, slug string) (*model.Badge, error) {
	return m.badges[slug], nil
}

func (m *mockBadgeCatalog) List(_ context.Context) ([]*model.Badge, error) {
	var result []*model.Badge
	for _, b := range m.badges {
		result = append(result, b)
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	users map[uuid.UUID]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uuid.UUID]*model.User)}
}

func (m *mockUserRepo) Create(_ context.Context, u *model.User) error {
	m.users[u.ID] = u
	return nil
}
func (m *mockUserRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	return m.users[id], nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) GetByUsername(_ context.Context, username string) (*model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) Update(_ context.Context, u *model.User) error {
	m.users[u.ID] = u
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
type mockAccountRepo struct {
	accounts map[uuid.UUID]*model.Account
}

func newMockAccountRepo() *mockAccountRepo {
	return &mockAccountRepo{accounts: make(map[uuid.UUID]*model.Account)}
}

func (m *mockAccountRepo) Create(_ context.Context, a *model.Account) error {
	m.accounts[a.ID] = a
	return nil
}

func (m *mockAccountRepo) GetByID(_ context.Context, id uuid.UUID) (*model.Account, error) {
	return m.accounts[id], nil
}

func (m *mockAccountRepo) GetByProvider(_ context.Context, provider, providerUserID string) (*model.Account, error) {
	return nil, nil
}

func (m *mockAccountRepo) GetByUserID(_ context.Context, userID uuid.UUID) (*model.Account, error) {
	return nil, nil
}

func (m *mockAccountRepo) Update(_ context.Context, a *model.Account) error {
	m.accounts[a.ID] = a
	return nil
}

func (m *mockAccountRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.accounts, id)
	return nil
}

// ─── helpers ──────────────────────────────────────────────────────────────────

type testDeps struct {
	plates       *mockPlateRepo
	tags         *mockTagRepo
	members      *mockMemberRepo
	badges       *mockBadgeRepo
	badgeCatalog *mockBadgeCatalog
	accounts     *mockAccountRepo
	users        *mockUserRepo
}

func newDeps() *testDeps {
	return &testDeps{
		plates:       newMockPlateRepo(),
		tags:         newMockTagRepo(),
		members:      newMockMemberRepo(),
		badges:       newMockBadgeRepo(),
		badgeCatalog: newMockBadgeCatalog(),
		accounts:     newMockAccountRepo(),
		users:        newMockUserRepo(),
	}
}

func newServiceFromDeps(d *testDeps) plateservice.PlateService {
	return plateservice.NewPlateServiceForTest(
		d.plates,
		d.tags,
		d.members,
		d.badges,
		d.badgeCatalog,
		d.accounts,
		d.users,
		lib.GetLogger(),
	)
}

func seedAccount(d *testDeps) *model.Account {
	a := &model.Account{
		ID:             uuid.New(),
		Provider:       "local",
		ProviderUserID: uuid.New().String(),
	}
	_ = d.accounts.Create(context.Background(), a)
	return a
}

func seedFilePlate(d *testDeps, ownerID uuid.UUID) *model.Plate {
	content := "FROM golang:1.22-alpine"
	filename := "Dockerfile"
	now := time.Now()
	p := &model.Plate{
		ID:         uuid.New(),
		OwnerID:    ownerID,
		Type:       model.PlateTypeFile,
		Slug:       "go-dockerfile",
		Name:       "Go Dockerfile",
		Category:   "devops",
		Status:     model.PlateStatusPending,
		Visibility: model.PlateVisibilityPublic,
		Content:    &content,
		Filename:   &filename,
		IsVerified: true,
		VerifiedAt: &now,
	}
	_ = d.plates.Create(context.Background(), p)
	_ = d.members.Create(context.Background(), &model.PlateMember{
		ID:        uuid.New(),
		PlateID:   p.ID,
		AccountID: ownerID,
		Role:      model.PlateMemberRoleOwner,
	})
	return p
}

// ─── SubmitFile ───────────────────────────────────────────────────────────────

func TestSubmitFile_Success(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	account := seedAccount(d)

	p, err := svc.SubmitFile(context.Background(), account.ID, plateservice.SubmitFileInput{
		Name:       "Go Dockerfile",
		Category:   "devops",
		Visibility: model.PlateVisibilityPublic,
		Filename:   "Dockerfile",
		Content:    "FROM golang:1.22-alpine",
		Tags:       []string{"go", "docker"},
	})

	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, model.PlateTypeFile, p.Type)
	assert.Equal(t, model.PlateStatusPending, p.Status)
	assert.True(t, p.IsVerified)
	assert.Equal(t, "go-dockerfile", p.Slug)
}

func TestSubmitFile_CreatesOwnerMemberRow(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	account := seedAccount(d)

	p, err := svc.SubmitFile(context.Background(), account.ID, plateservice.SubmitFileInput{
		Name:       "Go Dockerfile",
		Category:   "devops",
		Visibility: model.PlateVisibilityPublic,
		Filename:   "Dockerfile",
		Content:    "FROM golang:1.22-alpine",
	})

	require.NoError(t, err)

	member, err := d.members.GetByPlateAndAccount(context.Background(), p.ID, account.ID)
	require.NoError(t, err)
	require.NotNil(t, member)
	assert.Equal(t, model.PlateMemberRoleOwner, member.Role)
}

func TestSubmitFile_StoresTags(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	account := seedAccount(d)

	p, err := svc.SubmitFile(context.Background(), account.ID, plateservice.SubmitFileInput{
		Name:       "Go Dockerfile",
		Category:   "devops",
		Visibility: model.PlateVisibilityPublic,
		Filename:   "Dockerfile",
		Content:    "FROM golang:1.22-alpine",
		Tags:       []string{"go", "docker"},
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"go", "docker"}, d.tags.tags[p.ID])
}

func TestSubmitFile_NoTags_DoesNotError(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	account := seedAccount(d)

	_, err := svc.SubmitFile(context.Background(), account.ID, plateservice.SubmitFileInput{
		Name:       "Go Dockerfile",
		Category:   "devops",
		Visibility: model.PlateVisibilityPublic,
		Filename:   "Dockerfile",
		Content:    "FROM golang:1.22-alpine",
	})

	require.NoError(t, err)
}

// ─── GetBySlug ────────────────────────────────────────────────────────────────

func TestGetBySlug_PublicPlate_AnyCallerCanRead(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	stranger := uuid.New()
	result, err := svc.GetBySlug(context.Background(), p.Slug, stranger)

	require.NoError(t, err)
	assert.Equal(t, p.ID, result.ID)
}

func TestGetBySlug_NotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)

	result, err := svc.GetBySlug(context.Background(), "does-not-exist", uuid.New())

	assert.Nil(t, result)
	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

func TestGetBySlug_PrivatePlate_NonMemberGetsForbidden(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	p.Visibility = model.PlateVisibilityPrivate
	_ = d.plates.Update(context.Background(), p)

	stranger := uuid.New()
	result, err := svc.GetBySlug(context.Background(), p.Slug, stranger)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

func TestGetBySlug_PrivatePlate_MemberCanRead(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	p.Visibility = model.PlateVisibilityPrivate
	_ = d.plates.Update(context.Background(), p)

	result, err := svc.GetBySlug(context.Background(), p.Slug, owner.ID)

	require.NoError(t, err)
	assert.Equal(t, p.ID, result.ID)
}

// ─── List ─────────────────────────────────────────────────────────────────────

func TestList_DefaultsToPublicVisibility(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	seedFilePlate(d, owner.ID)

	plates, total, err := svc.List(context.Background(), repository.PlateFilter{}, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, plates, 1)
}

func TestList_PrivatePlateNotReturnedByDefault(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	p.Visibility = model.PlateVisibilityPrivate
	_ = d.plates.Update(context.Background(), p)

	plates, total, err := svc.List(context.Background(), repository.PlateFilter{}, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, plates)
}

// ─── Update ───────────────────────────────────────────────────────────────────

func TestUpdate_Owner_CanUpdate(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	newDesc := "Updated description"
	result, err := svc.Update(context.Background(), p.ID, owner.ID, plateservice.UpdatePlateInput{
		Description: &newDesc,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated description", *result.Description)
}

func TestUpdate_Stranger_Forbidden(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	newDesc := "Should not work"
	result, err := svc.Update(context.Background(), p.ID, uuid.New(), plateservice.UpdatePlateInput{
		Description: &newDesc,
	})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, plateservice.ErrForbidden)
}

func TestUpdate_PlateNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)

	newDesc := "desc"
	result, err := svc.Update(context.Background(), uuid.New(), uuid.New(), plateservice.UpdatePlateInput{
		Description: &newDesc,
	})

	assert.Nil(t, result)
	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

// ─── Archive ──────────────────────────────────────────────────────────────────

func TestArchive_Owner_SetsStatusArchived(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.Archive(context.Background(), p.ID, owner.ID)

	require.NoError(t, err)
	stored := d.plates.plates[p.ID]
	assert.Equal(t, model.PlateStatusArchived, stored.Status)
}

func TestArchive_Stranger_Forbidden(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.Archive(context.Background(), p.ID, uuid.New())

	assert.ErrorIs(t, err, plateservice.ErrForbidden)
}

func TestArchive_PlateNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)

	err := svc.Archive(context.Background(), uuid.New(), uuid.New())

	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

// ─── RecordUse ────────────────────────────────────────────────────────────────

func TestRecordUse_NewUser_CreatesMemberRow(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	stranger := uuid.New()

	err := svc.RecordUse(context.Background(), p.ID, stranger)

	require.NoError(t, err)
	member, err := d.members.GetByPlateAndAccount(context.Background(), p.ID, stranger)
	require.NoError(t, err)
	require.NotNil(t, member)
	assert.Equal(t, model.PlateMemberRoleMember, member.Role)
	assert.NotNil(t, member.LastUsedAt)
}

func TestRecordUse_ExistingUser_UpdatesLastUsedAt(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	user := seedAccount(d)

	before := time.Now().Add(-1 * time.Minute)
	_ = d.members.Create(context.Background(), &model.PlateMember{
		ID:         uuid.New(),
		PlateID:    p.ID,
		AccountID:  user.ID,
		Role:       model.PlateMemberRoleMember,
		LastUsedAt: &before,
	})

	err := svc.RecordUse(context.Background(), p.ID, user.ID)

	require.NoError(t, err)
	member, _ := d.members.GetByPlateAndAccount(context.Background(), p.ID, user.ID)
	assert.True(t, member.LastUsedAt.After(before))
}

func TestRecordUse_IncrementsUseCount(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.RecordUse(context.Background(), p.ID, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 1, d.plates.plates[p.ID].UseCount)
}

// ─── ReplaceTags ─────────────────────────────────────────────────────────────

func TestReplaceTags_Owner_ReplacesTags(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	d.tags.tags[p.ID] = []string{"go", "docker"}

	err := svc.ReplaceTags(context.Background(), p.ID, owner.ID, []string{"rust", "wasm"})

	require.NoError(t, err)
	assert.Equal(t, []string{"rust", "wasm"}, d.tags.tags[p.ID])
}

func TestReplaceTags_EmptySlice_ClearsTags(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	d.tags.tags[p.ID] = []string{"go", "docker"}

	err := svc.ReplaceTags(context.Background(), p.ID, owner.ID, []string{})

	require.NoError(t, err)
	assert.Empty(t, d.tags.tags[p.ID])
}

func TestReplaceTags_Stranger_Forbidden(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.ReplaceTags(context.Background(), p.ID, uuid.New(), []string{"rust"})

	assert.ErrorIs(t, err, plateservice.ErrForbidden)
}

// ─── Approve / Reject ─────────────────────────────────────────────────────────

func TestApprove_SetsStatusApprovedAndPublishedAt(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	admin := uuid.New()

	err := svc.Approve(context.Background(), p.ID, admin)

	require.NoError(t, err)
	stored := d.plates.plates[p.ID]
	assert.Equal(t, model.PlateStatusApproved, stored.Status)
	assert.NotNil(t, stored.PublishedAt)
}

func TestApprove_PlateNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)

	err := svc.Approve(context.Background(), uuid.New(), uuid.New())

	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

func TestReject_SetsStatusRejected(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	admin := uuid.New()

	err := svc.Reject(context.Background(), p.ID, admin)

	require.NoError(t, err)
	stored := d.plates.plates[p.ID]
	assert.Equal(t, model.PlateStatusRejected, stored.Status)
}

func TestReject_PlateNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)

	err := svc.Reject(context.Background(), uuid.New(), uuid.New())

	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

// ─── GrantBadge / RevokeBadge ─────────────────────────────────────────────────

func TestGrantBadge_KnownSlug_Grants(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	admin := uuid.New()

	badgeID := uuid.New()
	_ = d.badgeCatalog.Create(context.Background(), &model.Badge{
		ID:   badgeID,
		Slug: "featured",
		Name: "Featured",
		Tier: model.BadgeTierOfficial,
	})

	err := svc.GrantBadge(context.Background(), p.ID, admin, "featured", nil)

	require.NoError(t, err)
	grants, _ := d.badges.ListByPlate(context.Background(), p.ID)
	require.Len(t, grants, 1)
	assert.Equal(t, badgeID, grants[0].BadgeID)
	assert.Equal(t, admin.String(), grants[0].GrantedBy)
}

func TestGrantBadge_UnknownSlug_ReturnsNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.GrantBadge(context.Background(), p.ID, uuid.New(), "nonexistent", nil)

	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}

func TestRevokeBadge_RemovesGrant(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)
	admin := uuid.New()

	badgeID := uuid.New()
	_ = d.badgeCatalog.Create(context.Background(), &model.Badge{
		ID:   badgeID,
		Slug: "featured",
		Name: "Featured",
		Tier: model.BadgeTierOfficial,
	})

	require.NoError(t, svc.GrantBadge(context.Background(), p.ID, admin, "featured", nil))
	require.NoError(t, svc.RevokeBadge(context.Background(), p.ID, admin, "featured"))

	grants, _ := d.badges.ListByPlate(context.Background(), p.ID)
	assert.Empty(t, grants)
}

func TestRevokeBadge_UnknownSlug_ReturnsNotFound(t *testing.T) {
	d := newDeps()
	svc := newServiceFromDeps(d)
	owner := seedAccount(d)
	p := seedFilePlate(d, owner.ID)

	err := svc.RevokeBadge(context.Background(), p.ID, uuid.New(), "nonexistent")

	assert.ErrorIs(t, err, plateservice.ErrNotFound)
}
