package plate

import (
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type KickplateYAML struct {
	Owner             string           `yaml:"owner"`
	Name              string           `yaml:"name"`
	Description       string           `yaml:"description"`
	Category          string           `yaml:"category"`
	VerificationToken string           `yaml:"verification_token"`
	Tags              []string         `yaml:"tags"`
	Variables         []map[string]any `yaml:"variables"`
	Dependencies      []map[string]any `yaml:"dependencies"`
}

type plateService struct {
	db           *gorm.DB
	env          lib.Env
	plates       repository.PlateRepository
	orgs         repository.OrganizationRepository
	tags         repository.PlateTagRepository
	members      repository.PlateMemberRepository
	badges       repository.PlateBadgeRepository
	badgeCatalog repository.BadgeRepository
	accounts     repository.AccountRepository
	users        repository.UserRepository
	reviews      repository.PlateReviewRepository
	logger       lib.Logger
}

func NewPlateService(
	db lib.Database,
	env lib.Env,
	plates repository.PlateRepository,
	orgs repository.OrganizationRepository,
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
		db:           db.DB,
		env:          env,
		plates:       plates,
		orgs:         orgs,
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		reviews:      reviews,
		logger:       logger,
	}
}

func NewPlateServiceForTest(
	plates repository.PlateRepository,
	orgs repository.OrganizationRepository,
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
		tags:         tags,
		members:      members,
		badges:       badges,
		badgeCatalog: badgeCatalog,
		accounts:     accounts,
		users:        users,
		reviews:      reviews,
		logger:       logger,
	}
}
