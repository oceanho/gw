package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	g := r.Group("/api/v1")
	g.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"payload": "welcome.",
		})
	})
	r.Run(":9000")
}