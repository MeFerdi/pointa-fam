package services

import (
	"errors"
	"pointafam/backend/models"

	"gorm.io/gorm"
)

type FarmerService struct {
	DB *gorm.DB // Database connection instance
}

// NewFarmerService initializes a new FarmerService
func NewFarmerService(db *gorm.DB) *FarmerService {
	return &FarmerService{DB: db}
}

// CreateFarmer creates a new farmer in the database
func (s *FarmerService) CreateFarmer(farmer *models.Farmer) error {
	// Validate farmer data (e.g., check for empty fields)
	if farmer.Name == "" {
		return errors.New("farmer name and contact info are required")
	}

	return s.DB.Create(farmer).Error
}

// GetAllFarmers retrieves all farmers from the database
func (s *FarmerService) GetAllFarmers() ([]models.Farmer, error) {
	var farmers []models.Farmer
	err := s.DB.Find(&farmers).Error
	return farmers, err
}

// UpdateFarmer updates an existing farmer's details
func (s *FarmerService) UpdateFarmer(id uint, updatedFarmer *models.Farmer) error {
	// Check if the farmer exists before updating
	var farmer models.Farmer
	if err := s.DB.First(&farmer, id).Error; err != nil {
		return errors.New("farmer not found")
	}

	return s.DB.Model(&farmer).Updates(updatedFarmer).Error
}

// DeleteFarmer removes a farmer from the database by ID
func (s *FarmerService) DeleteFarmer(id uint) error {
	// Check if the farmer exists before deleting
	var farmer models.Farmer
	if err := s.DB.First(&farmer, id).Error; err != nil {
		return errors.New("farmer not found")
	}

	return s.DB.Delete(&farmer).Error
}

// ManageProducts allows farmers to manage their products (add/update/remove)
func (s *FarmerService) ManageProducts(farmerID uint, products []models.Product) error {
	// Check if the farmer exists
	var farmer models.Farmer
	if err := s.DB.First(&farmer, farmerID).Error; err != nil {
		return errors.New("farmer not found")
	}

	// Begin a transaction
	tx := s.DB.Begin()

	// Remove existing products for the farmer
	if err := tx.Where("farmer_id = ?", farmerID).Delete(&models.Product{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new products for the farmer
	for _, product := range products {
		product.FarmerID = farmerID
		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}
