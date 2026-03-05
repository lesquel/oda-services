package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/lesquel/oda-shared/domain"
)

// Connect opens a GORM connection to the PostgreSQL database.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)
	return db, nil
}

// RunMigrations runs GORM AutoMigrate for all domain models.
// Only the write-api connects to the primary database and runs migrations.
func RunMigrations(db *gorm.DB) error {
	// Enable uuid-ossp extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Printf("uuid-ossp extension: %v", err)
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.Poem{},
		&domain.Like{},
		&domain.EmotionTag{},
		&domain.Bookmark{},
		&domain.EmotionCatalog{},
		&domain.ModerationLog{},
	); err != nil {
		return err
	}

	// Ensure UUID defaults via raw SQL (pgx driver strips uuid_generate_v4() from AutoMigrate)
	tables := []string{"users", "refresh_tokens", "poems", "likes", "emotion_tags", "bookmarks", "emotion_catalog", "moderation_logs"}
	for _, t := range tables {
		db.Exec(`ALTER TABLE ` + t + ` ALTER COLUMN id SET DEFAULT uuid_generate_v4()`)
	}

	// Foreign key constraints (idempotent – GORM adds ALTER TABLE … ADD CONSTRAINT)
	fks := []string{
		`ALTER TABLE refresh_tokens DROP CONSTRAINT IF EXISTS fk_refresh_tokens_user;
 ALTER TABLE refresh_tokens ADD CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`,
		`ALTER TABLE poems DROP CONSTRAINT IF EXISTS fk_poems_author;
 ALTER TABLE poems ADD CONSTRAINT fk_poems_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE`,
		`ALTER TABLE likes DROP CONSTRAINT IF EXISTS fk_likes_user;
 ALTER TABLE likes ADD CONSTRAINT fk_likes_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`,
		`ALTER TABLE likes DROP CONSTRAINT IF EXISTS fk_likes_poem;
 ALTER TABLE likes ADD CONSTRAINT fk_likes_poem FOREIGN KEY (poem_id) REFERENCES poems(id) ON DELETE CASCADE`,
		`ALTER TABLE emotion_tags DROP CONSTRAINT IF EXISTS fk_emotion_tags_user;
 ALTER TABLE emotion_tags ADD CONSTRAINT fk_emotion_tags_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`,
		`ALTER TABLE emotion_tags DROP CONSTRAINT IF EXISTS fk_emotion_tags_poem;
 ALTER TABLE emotion_tags ADD CONSTRAINT fk_emotion_tags_poem FOREIGN KEY (poem_id) REFERENCES poems(id) ON DELETE CASCADE`,
		`ALTER TABLE emotion_tags DROP CONSTRAINT IF EXISTS fk_emotion_tags_catalog;
 ALTER TABLE emotion_tags ADD CONSTRAINT fk_emotion_tags_catalog FOREIGN KEY (emotion_id) REFERENCES emotion_catalog(id) ON DELETE CASCADE`,
		`ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS fk_bookmarks_user;
 ALTER TABLE bookmarks ADD CONSTRAINT fk_bookmarks_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`,
		`ALTER TABLE bookmarks DROP CONSTRAINT IF EXISTS fk_bookmarks_poem;
 ALTER TABLE bookmarks ADD CONSTRAINT fk_bookmarks_poem FOREIGN KEY (poem_id) REFERENCES poems(id) ON DELETE CASCADE`,
		`ALTER TABLE moderation_logs DROP CONSTRAINT IF EXISTS fk_moderation_logs_poem;
 ALTER TABLE moderation_logs ADD CONSTRAINT fk_moderation_logs_poem FOREIGN KEY (poem_id) REFERENCES poems(id) ON DELETE CASCADE`,
	}
	for _, fk := range fks {
		if err := db.Exec(fk).Error; err != nil {
			log.Printf("FK constraint warning: %v", err)
		}
	}

	return nil
}
