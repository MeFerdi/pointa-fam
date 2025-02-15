package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductService instance
var productService *services.ProductService

// SetProductService is called to set the ProductService instance (usually done in main)
func SetProductService(service *services.ProductService) {
	productService = service
}

func GetProducts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductsByCategory retrieves products by category from the database
func GetProductsByCategory(c *gin.Context) {
	category := c.Query("category")
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("category = ?", category).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductsByUser retrieves products added by a specific user
func GetProductsByUser(c *gin.Context) {
	userID := c.Param("id")
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("user_id = ?", userID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
