package router

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/handlers"
	"github.com/prixplus/server/middlewares"
	"os"
)

func Init(db *sql.DB) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	// We have our own Recovery
	// r.Use(gin.Recovery())
	r.Use(middlewares.Recovery(db))

	directory := os.Getenv("GOPATH") + "/src/github.com/prixplus/server/"
	r.LoadHTMLGlob(directory + "templates/*")
	r.Static("/assets", "./assets")
	//r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Main route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "Prix!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", handlers.Login(db))

	auth := r.Group("/auth")
	auth.Use(middlewares.Auth(db))
	{
		auth.GET("/refresh_token", handlers.Refresh(db))
		auth.GET("/hello", handlers.Hello)

	}

	return r
}
