package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pointafam/backend/config"
	"pointafam/backend/controllers"
	"pointafam/backend/middleware"
	"pointafam/backend/migrations"
	"pointafam/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) > 1 {
		// Handle user deletion
		userID := os.Args[1]
		db, err := gorm.Open(sqlite.Open("./data/pointafam.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Could not connect to database: %v", err)
		}

		id, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			log.Fatalf("Invalid user ID: %v", err)
		}

		if err := models.DeleteUser(db, uint(id)); err != nil {
			log.Fatalf("Could not delete user: %v", err)
		}

		fmt.Printf("User with ID %d deleted successfully\n", id)
		return
	}

	cfg := config.LoadConfig()

	// Ensure the data directory exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		if err := os.Mkdir("./data", os.ModePerm); err != nil {
			log.Fatalf("Could not create data directory: %v", err)
		}
	}

	// Check if the database file exists
	if _, err := os.Stat(cfg.DBPath); os.IsNotExist(err) {
		log.Printf("Database file does not exist: %s", cfg.DBPath)
	} else {
		log.Printf("Database file exists: %s", cfg.DBPath)
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Run migrations
	migrations.Migrate(db)
	controllers.SetDB(db)

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Set trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Serve static files
	r.Static("/static", "./public/static")

	// Serve the HTML file
	r.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})

	// Serve role fields dynamically
	r.GET("/role_fields.html", func(c *gin.Context) {
		role := c.Query("role")
		c.HTML(http.StatusOK, "./public/role_fields.html", gin.H{"role": role})
	})

	// Serve auth.html and login.html through routes
	r.GET("/auth", func(c *gin.Context) {
		c.File("./public/static/auth.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.File("./public/static/login.html")
	})

	// Serve farmer and retailer dashboards
	r.GET("/farmer/dashboard", func(c *gin.Context) {
		c.File("./public/static/farmer_dashboard.html")
	})

	r.GET("/retailer/dashboard", func(c *gin.Context) {
		c.File("./public/static/retailer_dashboard.html")
	})

	// Authentication routes
	r.POST("/api/register", controllers.SignUp)
	r.POST("/api/login", controllers.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/user/:id", controllers.GetUserProfile)
		api.PUT("/user/:id", controllers.UpdateUserProfile)
		api.DELETE("/user/:id", controllers.DeleteUser) // Add this line

		api.GET("/products", controllers.GetProducts)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)

		api.POST("/cart", controllers.AddToCart)
		api.GET("/cart/:retailer_id", controllers.ViewCart)
		api.DELETE("/cart/:item_id", controllers.DeleteFromCart)

		// api.GET("/orders/current", controllers.GetCurrentOrders)
		// api.GET("/orders/history", controllers.GetOrderHistory)

		// api.GET("/suppliers", controllers.GetSuppliers)
		// api.POST("/suppliers", controllers.CreateSupplier)

		// api.GET("/inventory", controllers.GetInventory)

		// api.GET("/sales-analytics", controllers.GetSalesAnalytics)

		// api.GET("/messages", controllers.GetMessages)
		// api.GET("/feedback", controllers.GetFeedback)

		// api.GET("/resources", controllers.GetResources)
		// api.GET("/faqs", controllers.GetFAQs)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
