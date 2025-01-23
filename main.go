package main

import (
	"log"
	"net/http"
	"os"

	"pointafam/backend/config"
	"pointafam/backend/controllers"
	"pointafam/backend/middleware"
	"pointafam/backend/migrations"
	"pointafam/backend/models"
	"pointafam/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Subscriber model
type Subscriber struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique;not null"`
}

func main() {
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
	db.AutoMigrate(&Subscriber{}, &models.User{}, &models.Cart{}, &models.CartItem{})
	migrations.Migrate(db)
	controllers.SetDB(db)

	productService := services.NewProductService(db)
	controllers.SetProductService(productService)

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Set trusted proxies
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Serve static files
	r.Static("/public", "./public")

	// Serve the HTML file for homepage
	r.GET("/", func(c *gin.Context) {
		c.File("./index.html")
	})

	// Serve About Us page
	r.GET("/about", func(c *gin.Context) {
		c.File("./public/static/about.html")
	})

	// Serve Products by Category page
	r.GET("/products", func(c *gin.Context) {
		c.File("./public/static/products_by_category.html")
	})

	// Serve Contact Us page
	r.GET("/contact", func(c *gin.Context) {
		c.File("./public/static/contact.html")
	})

	// Serve role fields dynamically
	r.GET("/role_fields.html", func(c *gin.Context) {
		role := c.Query("role")
		c.HTML(http.StatusOK, "./public/role_fields.html", gin.H{"role": role})
	})

	// Serve auth.html and login.html through routes
	r.GET("/register", func(c *gin.Context) {
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
	r.GET("/cart", func(c *gin.Context) {
		c.File("./public/static/cart.html")
	})

	// Authentication routes
	r.POST("/api/register", controllers.SignUp)
	r.POST("/api/login", controllers.Login)

	// Use the DB middleware
	r.Use(middleware.DBMiddleware(db))

	// Endpoint to handle subscription
	r.POST("/subscribe", func(c *gin.Context) {
		var request struct {
			Email string `json:"email"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
			return
		}

		if request.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email address"})
			return
		}

		// Save the email to the database
		subscriber := Subscriber{Email: request.Email}
		if err := db.Create(&subscriber).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Subscription failed. Please try again."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Thank you for subscribing!"})
	})

	// Public API endpoint to get products by category
	r.GET("/api/products", controllers.GetProductsByCategory)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/user/:id", controllers.GetUserProfile)
		api.PUT("/user/:id", controllers.UpdateUserProfile)
		// api.DELETE("/user/:id", controllers.DeleteUser)

		// api.GET("/products/:id", controllers.GetProductByID)
		api.GET("/api/products/category", controllers.GetProductsByCategory)
		// api.GET("/user/:id/products", controllers.GetProductsByUser)
		api.POST("/api/user/:id/profile-picture", controllers.UploadProfilePicture)
		// api.GET("/user/:id/profile-picture", controllers.GetProfilePicture)
	}
	api.POST("/products", controllers.CreateProduct)
	api.PUT("/products/:id", controllers.UpdateProduct)
	api.DELETE("/products/:id", controllers.DeleteProduct)

	api.POST("/cart", controllers.AddToCart)
	api.GET("/cart/:retailer_id", controllers.ViewCart)
	api.DELETE("/cart/:item_id", controllers.RemoveFromCart)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
