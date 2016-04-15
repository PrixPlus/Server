package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/handler"
	"github.com/prixplus/server/middleware"
	"os"
)

// HTTP methods and status code follow REST convention
// http://www.restapitutorial.com/lessons/httpmethods.html
func Init(db *sql.DB) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.Recovery(db))
	r.Use(middleware.Auth(db))

	directory := os.Getenv("GOPATH") + "/src/github.com/prixplus/server/"
	r.LoadHTMLGlob(directory + "templates/*")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.Static("/assets", "./assets")

	// Main route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "Prix!",
		})
	})

	// API VERSION 0
	// because will be like api1, api2, api3...
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		api.POST("/login", handler.Login(db))
		api.GET("/refresh_token", handler.Refresh(db))

		// Get user himself
		api.GET("/me", handler.GetMe(db))
		// Insert a new user
		api.POST("/users", handler.PostUser(db))
		// Update an new user
		api.PUT("/users/:id", handler.PutUser(db))

	}

	return r
}
