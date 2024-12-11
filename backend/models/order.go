package models

import (
	"gorm.io/gorm"
)

type Order struct {
	OrderID    uint   `json:"order_id" gorm:"primaryKey"`
	RetailerID uint   `json:"retailer_id" gorm:"index"`
	ProductID  uint   `json:"product_id" gorm:"index"`
	Quantity   int    `json:"quantity"`
	Status     string `json:"status"` // e.g., pending, completed, canceled
}

// CreateOrder inserts a new order into the database
func (o *Order) CreateOrder(db *gorm.DB) error {
	return db.Create(o).Error
}

// GetOrdersByRetailer retrieves orders for a specific retailer from the database
func GetOrdersByRetailer(db *gorm.DB, retailerID uint) ([]Order, error) {
	var orders []Order
	err := db.Where("retailer_id = ?", retailerID).Find(&orders).Error
	return orders, err
}
