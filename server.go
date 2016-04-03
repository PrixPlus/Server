package main

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func HelloHandler(c *gin.Context) {
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

func main() {

	prixpath := os.Getenv("GOPATH") + "/src/github.com/prixplus/server/"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Init DB connection
	db, err := InitDB()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
	}

	// Close DB when main returns
	defer db.Close()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.LoadHTMLGlob(prixpath + "templates/*")
	r.Static("/assets", "./assets")
	//r.StaticFS("/more_static", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/close", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "finishing server",
		})
		manners.Close()
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Prix!",
		})
	})

	r.POST("/login", LoginHandler(db))
	r.GET("/refresh_token", RefreshHandler)

	auth := r.Group("/auth")
	auth.Use(AuthMiddleware(db))
	{
		auth.GET("/hello", HelloHandler)

	}

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	manners.ListenAndServe(":"+port, r)
}
