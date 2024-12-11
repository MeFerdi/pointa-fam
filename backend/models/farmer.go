package models

import (
	"gorm.io/gorm"
)

type Farmer struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	ContactInfo string `json:"contact_info"`
	Password    string `json:"password"` // Add password field
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
