package migrations

import (
	"log"
	"pointafam/backend/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Farmer{},
		&models.Product{},
		&models.Retailer{},
		&models.Order{},
		&models.Cart{},
		&models.CartItem{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrated successfully!")
}
