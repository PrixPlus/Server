package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "os"
	)

func main() {

    prixpath := os.Getenv("GOPATH")+"/src/github.com/prixplus/server/"

    r := gin.Default()
    r.LoadHTMLGlob(prixpath+"templates/*")
    r.Static("/assets", "./assets")
    //r.StaticFS("/more_static", http.Dir("my_file_system"))
    r.StaticFile("/favicon.ico", "./resources/favicon.ico")

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{
            "title": "Prix!",
        })
    })
    r.Run(":8080")
}