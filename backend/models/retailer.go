package models

import (
	"gorm.io/gorm"
)

type Retailer struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	ContactInfo string `json:"contact_info"`
	Location    string `json:"location"`
	Password    string `json:"password"`
}

// CreateRetailer inserts a new retailer into the database
func (r *Retailer) CreateRetailer(db *gorm.DB) error {
	return db.Create(r).Error
}

// GetAllRetailers retrieves all retailers from the database
func GetAllRetailers(db *gorm.DB) ([]Retailer, error) {
	var retailers []Retailer
	err := db.Find(&retailers).Error
	return retailers, err
}