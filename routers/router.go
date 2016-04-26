package routers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/handlers"
	"github.com/prixplus/server/middlewares"
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
	r.Use(middlewares.Recovery())
	r.Use(middlewares.Auth())

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

		api.POST("/login", handlers.Login())
		api.GET("/refresh_token", handlers.Refresh())

		// Get user himself
		api.GET("/me", handlers.GetMe())
		// Insert a new user
		api.POST("/users", handlers.PostUser())
		// Update an user
		api.PUT("/users/:id", handlers.PutUser())

		// Get products like product given
		api.GET("/products", handlers.GetProductList())
		// Get products by id
		api.GET("/products/:id", handlers.GetProduct())
		// Insert a new product
		api.POST("/products", handlers.PostProduct())
		// Update a product
		api.PUT("/products/:id", handlers.PutProduct())

	}

	return r
}
