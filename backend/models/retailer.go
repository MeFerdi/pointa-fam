package models

import (
    "gorm.io/gorm"
)

type Retailer struct {
    gorm.Model
    ID          uint   `json:"id" gorm:"primaryKey"`
    Name        string `json:"name"`
    PhoneNumber string `json:"phoneNumber"`
    Location    string `json:"location"`
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

// GetRetailerByID retrieves a retailer by their ID
func GetRetailerByID(db *gorm.DB, id uint) (*Retailer, error) {
    var retailer Retailer
    err := db.First(&retailer, id).Error
    return &retailer, err
}

// UpdateRetailer updates the retailer's profile in the database
func (r *Retailer) UpdateRetailer(db *gorm.DB) error {
    return db.Save(r).Error
}

// DeleteRetailer deletes a retailer from the database
func DeleteRetailer(db *gorm.DB, id uint) error {
    return db.Delete(&Retailer{}, id).Error
}