package main

import (
	"log"
	"net/http"
	"os"
	"pointafam/backend/config"
	"pointafam/backend/controllers"
	"pointafam/backend/middleware"
	"pointafam/backend/migrations"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
		api.GET("/farmers", controllers.GetFarmers)
		api.POST("/farmers", controllers.CreateFarmer)
		api.PUT("/farmers/:id", controllers.UpdateFarmer)
		api.DELETE("/farmers/:id", controllers.DeleteFarmer)

		api.GET("/products", controllers.GetProducts)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)

		api.GET("/retailers", controllers.GetRetailers)
		// api.POST("/retailers", controllers.CreateRetailer)
		api.PUT("/retailers/:id", controllers.UpdateRetailer)
		api.DELETE("/retailers/:id", controllers.DeleteRetailer)

		api.POST("/orders", controllers.CreateOrder)
		api.GET("/orders/:retailer_id", controllers.GetOrders)

		api.POST("/cart", controllers.AddToCart)
		api.GET("/cart/:retailer_id", controllers.ViewCart)
		api.DELETE("/cart/:item_id", controllers.DeleteFromCart)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
