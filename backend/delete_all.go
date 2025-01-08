package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("../data/pointafam.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Delete all users
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("Could not delete users: %v", err)
	}
	log.Println("All users deleted successfully")

	// Delete all farmers
	if err := db.Exec("DELETE FROM farmers").Error; err != nil {
		log.Fatalf("Could not delete farmers: %v", err)
	}
	log.Println("All farmers deleted successfully")

	// Delete all retailers
	if err := db.Exec("DELETE FROM retailers").Error; err != nil {
		log.Fatalf("Could not delete retailers: %v", err)
	}
	log.Println("All retailers deleted successfully")

	// Delete all products
	if err := db.Exec("DELETE FROM products").Error; err != nil {
		log.Fatalf("Could not delete products: %v", err)
	}
	log.Println("All products deleted successfully")
}
