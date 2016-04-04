package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
