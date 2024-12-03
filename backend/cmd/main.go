package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to PointaFam"})
	})
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
