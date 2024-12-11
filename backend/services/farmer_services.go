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
	if farmer.Name == "" || farmer.ContactInfo == "" {
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
func (s *FarmerService) ManageProducts(farmerID uint, product *models.Product) error {
	// Logic to manage products associated with this farmer.
	// This could include creating, updating, or deleting products.

	// Example: Adding a product to the database.
	product.FarmID = farmerID // Associate product with the specific farm/farmer
	return s.DB.Create(product).Error
}
