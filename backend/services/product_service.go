package services

import (
	"pointafam/backend/models"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return product.CreateProduct(s.DB)
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return models.GetAllProducts(s.DB)
}

func (s *ProductService) UpdateProduct(id uint, product *models.Product) error {
	return models.UpdateProduct(s.DB, id, product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return models.DeleteProduct(s.DB, id)
}
func (s *ProductService) GetProductsByCategory(category string) ([]models.Product, error) {
	var products []models.Product
	if err := s.DB.Where("category = ?", category).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
