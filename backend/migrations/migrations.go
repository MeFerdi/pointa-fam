package migrations

import (
	"log"
	"pointafam/backend/models"

	"gorm.io/gorm"
)

// Migrate function to run migrations
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Farmer{},
		// &models.Farm{},
		&models.Product{},
		&models.Retailer{},
		&models.Order{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrated successfully!")
}
