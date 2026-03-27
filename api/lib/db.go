package lib

import (
	"fmt"
	"time"

	"github.com/kickplate/api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(env Env, logger Logger) Database {
	logger.Info("Initializing database connection")
	username := env.DBUsername
	password := env.DBPassword
	host := env.DBHost
	port := env.DBPort
	dbname := env.DBName

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, username, password, dbname, port)

	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.GetGormLogger(),
		})
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			logger.Infof("Database connection attempt %d/%d failed, retrying in 2 seconds...", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		logger.Info("DSN: ", dsn)
		logger.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	logger.Info("Database connection established")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Account{},
		&model.Organization{},
		&model.EmailVerification{},
		&model.Plate{},
		&model.PlateTag{},
		&model.PlateMember{},
		&model.PlateReview{},
		&model.Badge{},
		&model.PlateBadge{},
	); err != nil {
		logger.Panicf("AutoMigrate failed: %v", err)
	}

	logger.Info("Database migrations applied")

	if err := db.Exec(`DROP TABLE IF EXISTS sync_log`).Error; err != nil {
		logger.Panicf("Failed to drop legacy sync_log table: %v", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		logger.Warnf("pg_trgm extension: %v", err)
	}

	_ = db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS search_vector`)
	if err := db.Exec(`
    ALTER TABLE plate
    ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(category, '')), 'C')
    ) STORED
`).Error; err != nil {
		logger.Warnf("search_vector: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_search
    ON plate USING GIN(search_vector)
`).Error; err != nil {
		logger.Warnf("idx_plate_search: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_name_trgm
    ON plate USING GIN(name gin_trgm_ops)
`).Error; err != nil {
		logger.Warnf("idx_plate_name_trgm: %v", err)
	}
	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_desc_trgm
    ON plate USING GIN(description gin_trgm_ops)
`).Error; err != nil {
		logger.Warnf("idx_plate_desc_trgm: %v", err)
	}

	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_status_visibility_usecount
    ON plate (status, visibility, bookmark_count DESC)
`).Error; err != nil {
		logger.Warnf("idx_plate_status_visibility_usecount: %v", err)
	}
	if err := db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_plate_tag_tag
    ON plate_tag (tag)
`).Error; err != nil {
		logger.Warnf("idx_plate_tag_tag: %v", err)
	}

	if err := db.Exec(`DELETE FROM plate WHERE type = 'file'`).Error; err != nil {
		logger.Warnf("cleanup file plates: %v", err)
	}
	if err := db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS content`).Error; err != nil {
		logger.Warnf("drop plate.content: %v", err)
	}
	if err := db.Exec(`ALTER TABLE plate DROP COLUMN IF EXISTS filename`).Error; err != nil {
		logger.Warnf("drop plate.filename: %v", err)
	}

	logger.Info("Extended migrations applied")

	return Database{
		DB: db,
	}
}

func (d Database) Close() error {
	if d.DB == nil {
		return nil
	}
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
