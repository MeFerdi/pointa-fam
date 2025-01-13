package models

import (
	"gorm.io/gorm"
)

type Farmer struct {
	gorm.Model

	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`
}

// CreateFarmer inserts a new farmer into the database
func (f *Farmer) CreateFarmer(db *gorm.DB) error {
	return db.Create(f).Error
}

// GetAllFarmers retrieves all farmers from the database
func GetAllFarmers(db *gorm.DB) ([]Farmer, error) {
	var farmers []Farmer
	err := db.Find(&farmers).Error
	return farmers, err
}

// GetFarmerByID retrieves a farmer by their ID
func GetFarmerByID(db *gorm.DB, id uint) (*Farmer, error) {
	var farmer Farmer
	err := db.First(&farmer, id).Error
	return &farmer, err
}

// UpdateFarmer updates the farmer's profile in the database
func (f *Farmer) UpdateFarmer(db *gorm.DB) error {
	return db.Save(f).Error
}
