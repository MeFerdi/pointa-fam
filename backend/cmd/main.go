package main

import (
	"log"
	"pointafam/backend/config"
	"pointafam/backend/controllers"
	"pointafam/backend/migrations"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	migrations.Migrate(db)
	controllers.SetDB(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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

	// Define routes for cart
	r.POST("/api/cart", controllers.AddToCart)
	r.GET("/api/cart/:retailer_id", controllers.ViewCart)
	r.DELETE("/api/cart/:item_id", controllers.DeleteFromCart)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
