package services

import (
	"pointafam/backend/models"

	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB // Database connection instance
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return order.CreateOrder(s.DB)
}

func (s *OrderService) GetOrdersByRetailer(retailerID uint) ([]models.Order, error) {
	return models.GetOrdersByRetailer(s.DB, retailerID)
}
