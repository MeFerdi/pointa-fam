package main

import (
	"log"
	"pointafam/backend/config"
	"pointafam/backend/controllers"
	"pointafam/backend/migrations"
	"pointafam/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite" // Change to postgres if using PostgreSQL
	"gorm.io/gorm"
)

func main() {
	// Load configuration from .env file
	cfg := config.LoadConfig()

	
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Run migrations to set up the database schema
	migrations.Migrate(db)

	// Initialize services and set them in controllers
	farmerService := services.NewFarmerService(db)
	controllers.SetFarmerService(farmerService)

	productService := services.NewProductService(db)
	controllers.SetProductService(productService)

	retailerService := services.NewRetailerService(db)
	controllers.SetRetailerService(retailerService)

	orderService := services.NewOrderService(db)
	controllers.SetOrderService(orderService)

	// Initialize Gin router
	r := gin.Default()

	// // Serve static files from the frontend directory
	// r.Static("/frontend", "./frontend")

	// // Serve index.html at root URL
	// r.GET("/", func(c *gin.Context) {
	// 	c.File("./frontend/index.html")
	// })

	// Define routes for farmers
	r.GET("/api/farmers", controllers.GetFarmers)
	r.POST("/api/farmers", controllers.CreateFarmer)
	r.PUT("/api/farmers/:id", controllers.UpdateFarmer)
	r.DELETE("/api/farmers/:id", controllers.DeleteFarmer)

	// Define routes for products
	r.GET("/api/products", controllers.GetProducts)
	r.POST("/api/products", controllers.CreateProduct)
	r.PUT("/api/products/:id", controllers.UpdateProduct)
	r.DELETE("/api/products/:id", controllers.DeleteProduct)

	// Define routes for retailers
	r.GET("/api/retailers", controllers.GetRetailers)
	r.POST("/api/retailers", controllers.CreateRetailer)
	r.PUT("/api/retailers/:id", controllers.UpdateRetailer)
	r.DELETE("/api/retailers/:id", controllers.DeleteRetailer)

	// Define routes for orders
	r.POST("/api/orders", controllers.CreateOrder)
	r.GET("/api/orders/:retailer_id", controllers.GetOrders)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
