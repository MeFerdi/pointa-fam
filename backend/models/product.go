package models

import (
	"gorm.io/gorm"
)

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	ImageURL    string  `json:"image_url"`
	Category    string  `json:"category" binding:"required"`
	FarmerID    uint    `json:"farm_id"`
	UserID      uint    `json:"user_id"`
}

// CreateProduct inserts a new product into the database
func (p *Product) CreateProduct(db *gorm.DB) error {
	return db.Create(p).Error
}

// GetAllProducts retrieves all products from the database
func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Find(&products).Error
	return products, err
}

// UpdateProduct updates an existing product in the database
func UpdateProduct(db *gorm.DB, id uint, product *Product) error {
	return db.Model(&Product{}).Where("id = ?", id).Updates(product).Error
}

// DeleteProduct removes a product from the database by ID
func DeleteProduct(db *gorm.DB, id uint) error {
	return db.Delete(&Product{}, id).Error
}
