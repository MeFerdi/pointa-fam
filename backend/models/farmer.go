package models

import (
	"gorm.io/gorm"
)

type Farmer struct {
	gorm.Model

	ID uint `gorm:"primaryKey"`

	Name        string `json:"name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
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
