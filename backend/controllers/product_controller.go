package controllers

import (
	"net/http"
	"path/filepath"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv" // Import strconv for string to uint conversion

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	c.HTML(http.StatusOK, "products_list.html", gin.H{"products": products})
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	// Parse form data
	product.Name = c.PostForm("name")
	product.Description = c.PostForm("description")
	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}
	product.Price = price

	quantity, err := strconv.Atoi(c.PostForm("quantity"))
	if err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}
		product.ImageURL = filepath
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := product.CreateProduct(db, product); err != nil {
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
