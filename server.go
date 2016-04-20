package main

import (
	"log"

	"github.com/braintree/manners"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/router"
	"github.com/prixplus/server/settings"
)

func main() {

	log.Println()

	// Load singleton settings
	sets, err := settings.Get()
	if err != nil {
		log.Fatal("Error loading settings: ", err)
	}

	// Init DB singleton connection
	db, err := database.Get()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
	}

	// Close DB when main() returns
	defer db.Close()

	// Logging the mode server is starting
	log.Printf("Server starting in %s mode at address: %s", sets.Env, ":8080\n\n")

	// Init routes
	handler := router.Init()

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	manners.ListenAndServe(":8080", handler)
}
