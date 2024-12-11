package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv" // Import strconv for string to uint conversion

	"github.com/gin-gonic/gin"
)

// ProductService instance
var productService *services.ProductService

// SetProductService is called to set the ProductService instance (usually done in main)
func SetProductService(service *services.ProductService) {
	productService = service
}

// GetProducts retrieves all products from the database
func GetProducts(c *gin.Context) {
	products, err := productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// CreateProduct adds a new product for a farmer
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := productService.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product by ID
func UpdateProduct(c *gin.Context) {
	var product models.Product
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := productService.UpdateProduct(uint(id), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct removes a product by ID from the database.
func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := productService.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete product"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
