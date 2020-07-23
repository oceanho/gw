package api

import (
	"github.com/gin-gonic/gin"
	gw2 "github.com/oceanho/gw"
)

func GetEnv(c *gw2.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func CreateEnv(c *gw2.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func ModifyEnv(c *gw2.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}

func DestroyEnv(c *gw2.ApiContext) {
	//user := c.User
	c.JSON(200, gin.H{
		"status": "succ",
		"payload": gin.H{
			"token": "",
		},
	})
}
