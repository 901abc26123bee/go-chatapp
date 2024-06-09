package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Create a Gin router
	r := gin.Default()

	// Define a route and handler function
	r.GET("/test", func(c *gin.Context) {
		// fix cors for local run, run with docker com-pose fixed in nginx
		c.Header("Access-Control-Allow-Origin", "*") // Allow requests from any origin
		c.Header("Access-Control-Allow-Methods", "GET") // Allow only GET method

		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})

		log.Println("ack")
	})

	// Run the Gin server
	r.Run(":8080")
}
