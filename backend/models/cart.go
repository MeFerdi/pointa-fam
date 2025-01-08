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
	ID        uint    `json:"id" gorm:"primaryKey"`
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product"`
	UserID    uint    `json:"user_id"`
}

// CreateCart inserts a new cart into the database
func (c *Cart) CreateCart(db *gorm.DB) error {
	return db.Create(c).Error
}

// GetCart retrieves a cart by retailer ID from the database
func GetCart(db *gorm.DB, retailerID uint) (Cart, error) {
	var cart Cart
	err := db.Preload("Items.Product").Where("retailer_id = ?", retailerID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = Cart{RetailerID: retailerID}
		if err := cart.CreateCart(db); err != nil {
			return cart, err
		}
		return cart, nil
	}
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

// GetCartItemsByCartID retrieves all cart items for a specific cart
func GetCartItemsByCartID(db *gorm.DB, cartID uint) ([]CartItem, error) {
	var cartItems []CartItem
	if err := db.Preload("Product").Where("cart_id = ?", cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}
