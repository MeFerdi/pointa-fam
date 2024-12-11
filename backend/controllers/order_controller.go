package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv" // Import strconv for string to uint conversion

	"github.com/gin-gonic/gin"
)

// OrderService instance
var orderService *services.OrderService

// SetOrderService is called to set the OrderService instance (usually done in main)
func SetOrderService(service *services.OrderService) {
	orderService = service
}

// CreateOrder places a new order for a retailer
func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := orderService.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrders retrieves all orders for a specific retailer
func GetOrders(c *gin.Context) {
	retailerIDStr := c.Param("retailer_id")
	retailerID, err := strconv.ParseUint(retailerIDStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retailer ID"})
		return
	}

	orders, err := orderService.GetOrdersByRetailer(uint(retailerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
