package plate

import (
	"github.com/kickplate/api/events"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository"
	"github.com/kickplate/api/service/githubapp"
	"gorm.io/gorm"
)

type RepositoryManifestYAML struct {
	Owner             string           `yaml:"owner"`
	Name              string           `yaml:"name"`
	Description       string           `yaml:"description"`
	Category          string           `yaml:"category"`
	VerificationToken string           `yaml:"verification_token"`
	Tags              []string         `yaml:"tags"`
	Variables         []map[string]any `yaml:"variables"`
	Dependencies      []map[string]any `yaml:"dependencies"`
}

type KickplateYAML = RepositoryManifestYAML

type plateService struct {
	db           *gorm.DB
	env          lib.Env
	plates       repository.PlateRepository
	orgs         repository.OrganizationRepository
	orgMembers   repository.OrganizationMemberRepository
	tags         repository.PlateTagRepository
	members      repository.PlateMemberRepository
	badges       repository.PlateBadgeRepository
	badgeCatalog repository.BadgeRepository
	accounts     repository.AccountRepository
	users        repository.UserRepository
	reviews      repository.PlateReviewRepository
	githubApp    githubapp.Service
	logger       lib.Logger
	emitter      *events.EventEmitter
}

func NewPlateService(
	db lib.Database,
	env lib.Env,
	plates repository.PlateRepository,
	orgs repository.OrganizationRepository,
	orgMembers repository.OrganizationMemberRepository,
	tags repository.PlateTagRepository,
	members repository.PlateMemberRepository,
	badges repository.PlateBadgeRepository,
	badgeCatalog repository.BadgeRepository,
	accounts repository.AccountRepository,
	users repository.UserRepository,
	reviews repository.PlateReviewRepository,
	githubApp githubapp.Service,
	logger lib.Logger,
	emitter *events.EventEmitter,
) PlateService {
	return &plateService{
		db:           db.DB,
		env:          env,
		plates:       plates,
		orgs:         orgs,
		orgMembers:   orgMembers,
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		reviews:      reviews,
		githubApp:    githubApp,
		logger:       logger,
		emitter:      emitter,
	}
}

func NewPlateServiceForTest(
	plates repository.PlateRepository,
	orgs repository.OrganizationRepository,
	orgMembers repository.OrganizationMemberRepository,
	tags repository.PlateTagRepository,
	members repository.PlateMemberRepository,
	badges repository.PlateBadgeRepository,
	badgeCatalog repository.BadgeRepository,
	accounts repository.AccountRepository,
	users repository.UserRepository,
	reviews repository.PlateReviewRepository,
	logger lib.Logger,
) PlateService {
	return &plateService{
		db:           nil,
		env:          lib.Env{},
		plates:       plates,
		orgs:         orgs,
		orgMembers:   orgMembers,
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		reviews:      reviews,
		logger:       logger,
		emitter:      events.NewEventEmitter(),
	}
}
