package services

import (
	"errors"
	"pointafam/backend/models"

	"gorm.io/gorm"
)

type RetailerService struct {
	DB *gorm.DB // Database connection instance
}

// NewRetailerService initializes a new RetailerService with a GORM DB instance
func NewRetailerService(db *gorm.DB) *RetailerService {
	return &RetailerService{DB: db}
}

// CreateRetailer adds a new retailer to the database
func (s *RetailerService) CreateRetailer(retailer *models.Retailer) error {
	// Validate retailer data (e.g., check for empty fields)
	if retailer.Name == "" || retailer.ContactInfo == "" {
		return errors.New("retailer name and contact info are required")
	}

	return s.DB.Create(retailer).Error
}

// GetAllRetailers retrieves all retailers from the database
func (s *RetailerService) GetAllRetailers() ([]models.Retailer, error) {
	var retailers []models.Retailer
	err := s.DB.Find(&retailers).Error
	return retailers, err
}

// UpdateRetailer updates an existing retailer's details by ID
func (s *RetailerService) UpdateRetailer(id string, updatedRetailer *models.Retailer) error {
	var retailer models.Retailer
	// Find the retailer by ID
	if err := s.DB.First(&retailer, id).Error; err != nil {
		return errors.New("retailer not found")
	}

	// Update retailer details
	return s.DB.Model(&retailer).Updates(updatedRetailer).Error
}

// DeleteRetailer removes a retailer from the database by ID
func (s *RetailerService) DeleteRetailer(id string) error {
	var retailer models.Retailer
	// Find the retailer by ID
	if err := s.DB.First(&retailer, id).Error; err != nil {
		return errors.New("retailer not found")
	}

	// Delete the retailer record from the database
	return s.DB.Delete(&retailer).Error
}
