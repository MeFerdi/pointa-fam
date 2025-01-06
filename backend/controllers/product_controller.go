package controllers

import (
	"log"
	"net/http"
	"path/filepath"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductService instance
var productService *services.ProductService

// SetProductService is called to set the ProductService instance (usually done in main)
func SetProductService(service *services.ProductService) {
	productService = service
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

// CreateProduct handles the creation of a new product.
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Parse form data
	product.Name = c.PostForm("name")
	product.Description = c.PostForm("description")
	priceStr := c.PostForm("price")
	if priceStr == "" {
		log.Printf("Price is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price is required"})
		return
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("Error parsing price: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}
	product.Price = price

	quantityStr := c.PostForm("quantity")
	if quantityStr == "" {
		log.Printf("Quantity is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity is required"})
		return
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		log.Printf("Error parsing quantity: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}
	product.Quantity = quantity

	product.Category = c.PostForm("category")

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		// Save the file to the server
		filename := filepath.Base(file.Filename)
		filepath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			log.Printf("Error saving file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}
		product.ImageURL = filepath
	} else {
		log.Printf("Error getting file: %v", err)
	}
	userIDStr := c.MustGet("userID").(string)
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		log.Printf("Error converting userID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to convert userID"})
		return
	}
	product.UserID = uint(userID)

	db := c.MustGet("db").(*gorm.DB)
	if err := product.CreateProduct(db); err != nil {
		log.Printf("Error creating product: %v", err)
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
