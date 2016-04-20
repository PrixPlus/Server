package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/handler"
	"github.com/prixplus/server/middleware"
	"github.com/prixplus/server/settings"
)

// HTTP methods and status code follow REST convention
// http://www.restapitutorial.com/lessons/httpmethods.html
func Init() *gin.Engine {

	sets, err := settings.Get()
	if err != nil {
		log.Fatal("Error getting Settings: ", err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Auth())

	r.LoadHTMLGlob(sets.Dir + "templates/*")
	r.StaticFile("/favicon.ico", sets.Dir+"assets/favicon.ico")
	r.Static("/assets", sets.Dir+"assets")

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

		api.POST("/login", handler.Login())
		api.GET("/refresh_token", handler.Refresh())

		// Get user himself
		api.GET("/me", handler.GetMe())
		// Insert a new user
		api.POST("/users", handler.PostUser())
		// Update an new user
		api.PUT("/users/:id", handler.PutUser())

	}

	return r
}
