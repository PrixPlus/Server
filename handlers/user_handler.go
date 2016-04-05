package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Just for testing
func Hello(c *gin.Context) {
	userName, ok := c.Get("userName")
	if !ok {
		c.JSON(505, gin.H{
			"error": "User not find in the header token",
		})
	}
	c.JSON(200, gin.H{
		"text": "Hello " + userName.(string) + "!",
	})
}

func InserUserHandler(c *gin.Context) {
	userName, ok := c.Get("userName")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not logged",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"text": "Hello " + userName.(string) + "!",
	})
}
