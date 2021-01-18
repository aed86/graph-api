package main

import (
	"github.com/gin-gonic/gin"

	"github.com/aed86/graph-api/config"
)

func main() {

	_ = config.Connect()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong23",
		})
	})
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})
	r.Run(":3001")
}