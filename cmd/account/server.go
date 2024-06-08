package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router
	r := gin.Default()

	// Define a route and handler function
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Run the Gin server
	r.Run() // Default listens on :8080
}
