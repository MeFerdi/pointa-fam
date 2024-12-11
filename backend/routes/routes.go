package routes

import (
	"pointafam/backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Farmer Routes
	r.GET("/api/farmers", controllers.GetFarmers)
	r.POST("/api/farmers", controllers.CreateFarmer)

	// Product Routes
	r.GET("/api/products", controllers.GetProducts)
	r.POST("/api/products", controllers.CreateProduct)
	r.PUT("/api/products/:id", controllers.UpdateProduct)
	r.DELETE("/api/products/:id", controllers.DeleteProduct)

	// Retailer Routes
	r.GET("/api/retailers", controllers.GetRetailers)
	r.POST("/api/retailers", controllers.CreateRetailer)

	// Order Routes
	r.POST("/api/orders", controllers.CreateOrder)
	r.GET("/api/orders/:retailer_id", controllers.GetOrders)
}
