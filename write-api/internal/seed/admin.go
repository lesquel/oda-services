package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/lesquel/oda-shared/domain"
	"github.com/lesquel/oda-shared/hasher"
	"gorm.io/gorm"
)

// SeedAdmin ensures an admin user exists when adminEmail and adminPassword are set.
func SeedAdmin(db *gorm.DB, adminEmail, adminPassword string) error {
	if adminEmail == "" || adminPassword == "" {
		log.Println("ℹ️  ADMIN_EMAIL/ADMIN_PASSWORD not set — skipping admin seed")
		return nil
	}

	if len(adminPassword) < 16 {
		log.Fatal("FATAL: ADMIN_PASSWORD must be at least 16 characters")
	}

	hashedPwd, err := hasher.HashPassword(adminPassword)
	if err != nil {
		log.Fatalf("FATAL: failed to hash admin password: %v", err)
	}

	var user domain.User
	result := db.Unscoped().Where("email = ?", adminEmail).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		admin := &domain.User{
			ID:           uuid.NewString(),
			Email:        adminEmail,
			Username:     "admin",
			PasswordHash: hashedPwd,
			Role:         "admin",
			IsActive:     true,
		}
		if err := db.Create(admin).Error; err != nil {
			admin.Username = "admin_" + adminEmail[:4]
			if err := db.Create(admin).Error; err != nil {
				log.Fatalf("FATAL: admin seed failed: %v", err)
			}
		}
		log.Printf("✅ Admin user created: %s", adminEmail)
	} else if result.Error == nil {
		db.Model(&user).Updates(map[string]interface{}{
			"password_hash": hashedPwd,
			"role":          "admin",
			"deleted_at":    nil,
		})
		log.Printf("✅ Admin user updated: %s (role=admin)", adminEmail)
	} else {
		log.Printf("⚠️  Admin seed DB error: %v", result.Error)
	}

	return nil
}
