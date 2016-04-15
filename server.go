package main

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/router"
	"log"
	"os"
	"strings"
)

func main() {
	// If var $MODE is set to RELEASE,
	// than starts server in release mode
	mode := strings.ToLower(os.Getenv("MODE"))
	if mode == "release" && gin.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init DB connection
	db, err := InitDB()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
	}

	// Close DB when main() returns
	defer db.Close()

	// Logging the mode server is starting
	log.Printf("Server starting in %s mode at address: %s", gin.Mode(), ":8080")

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	manners.ListenAndServe(":8080", router.Init(db))
}
