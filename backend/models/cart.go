package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	RetailerID uint       `json:"retailer_id"`
	Items      []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

// CreateCart inserts a new cart into the database
func (c *Cart) CreateCart(db *gorm.DB) error {
	return db.Create(c).Error
}

// GetCart retrieves a cart by retailer ID from the database
func GetCart(db *gorm.DB, retailerID uint) (Cart, error) {
	var cart Cart
	err := db.Preload("Items").Where("retailer_id = ?", retailerID).First(&cart).Error
	return cart, err
}

// AddCartItem inserts a new cart item into the database
func (ci *CartItem) AddCartItem(db *gorm.DB) error {
	return db.Create(ci).Error
}

// DeleteCartItem deletes a cart item from the database
func DeleteCartItem(db *gorm.DB, itemID uint) error {
	return db.Delete(&CartItem{}, itemID).Error
}
