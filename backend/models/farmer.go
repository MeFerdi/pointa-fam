package models

import (
	"gorm.io/gorm"
)

type Farmer struct {
	gorm.Model

	ID          uint      `gorm:"primaryKey"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	Location    string    `json:"location"`
	Products    []Product `json:"products" gorm:"foreignKey:FarmerID"` // One-to-many relationship
}

// CreateFarmer inserts a new farmer into the database
func (f *Farmer) CreateFarmer(db *gorm.DB) error {
	return db.Create(f).Error
}

// GetAllFarmers retrieves all farmers from the database
func GetAllFarmers(db *gorm.DB) ([]Farmer, error) {
	var farmers []Farmer
	err := db.Preload("Products").Find(&farmers).Error // Preload Products to include them in the response
	return farmers, err
}

// GetFarmerByID retrieves a farmer by their ID
func GetFarmerByID(db *gorm.DB, id uint) (*Farmer, error) {
	var farmer Farmer
	err := db.Preload("Products").First(&farmer, id).Error // Preload Products to include them in the response
	return &farmer, err
}

// UpdateFarmer updates the farmer's profile in the database
func (f *Farmer) UpdateFarmer(db *gorm.DB) error {
	return db.Save(f).Error
}
