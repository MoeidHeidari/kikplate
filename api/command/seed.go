package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"github.com/spf13/cobra"
)

type SeedCommand struct{}

func (c *SeedCommand) Short() string {
	return "Seed the database with initial data"
}

func (c *SeedCommand) Setup(_ *cobra.Command) {}

func (c *SeedCommand) Run() lib.CommandRunner {
	return func(
		logger lib.Logger,
		badgeRepo repository.BadgeRepository,
	) {
		ctx := context.Background()
		logger.Info("Seeding badges...")

		badges := []model.Badge{
			{
				ID:          uuid.New(),
				Slug:        "featured",
				Name:        "Featured",
				Description: "Hand-picked by the Kikplate team",
				Icon:        "star",
				Tier:        model.BadgeTierOfficial,
			},
			{
				ID:          uuid.New(),
				Slug:        "verified",
				Name:        "Verified",
				Description: "Ownership verified by Kikplate",
				Icon:        "shield",
				Tier:        model.BadgeTierVerified,
			},
			{
				ID:          uuid.New(),
				Slug:        "popular",
				Name:        "Popular",
				Description: "Used by many developers",
				Icon:        "fire",
				Tier:        model.BadgeTierCommunity,
			},
			{
				ID:          uuid.New(),
				Slug:        "official",
				Name:        "Official",
				Description: "Maintained by a trusted organization",
				Icon:        "badge-check",
				Tier:        model.BadgeTierOfficial,
			},
			{
				ID:          uuid.New(),
				Slug:        "sponsored",
				Name:        "Sponsored",
				Description: "Sponsored by a Kikplate partner",
				Icon:        "handshake",
				Tier:        model.BadgeTierSponsored,
			},
		}

		for _, badge := range badges {
			existing, err := badgeRepo.GetBySlug(ctx, badge.Slug)
			if err != nil {
				logger.Errorf("  ✗ %s: failed to check existence: %v", badge.Slug, err)
				continue
			}
			if existing != nil {
				logger.Infof("  ~ %s: already exists, skipping", badge.Slug)
				continue
			}

			b := badge
			if err := badgeRepo.Create(ctx, &b); err != nil {
				logger.Errorf("  ✗ %s: %v", badge.Slug, err)
				continue
			}
			logger.Infof("  ✓ %s", badge.Slug)
		}

		logger.Info("Seeding complete.")
	}
}

func NewSeedCommand() *SeedCommand {
	return &SeedCommand{}
}
