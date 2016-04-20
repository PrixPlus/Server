package main

import (
	"log"

	"github.com/braintree/manners"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"
)

func main() {

	log.Println()

	// Load singleton settings
	sets, err := settings.Get()
	if err != nil {
		log.Fatal("Error loading settings: ", err)
		return
	}

	// Init DB singleton connection
	db, err := database.Get()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
		return
	}

	// Close DB when main() returns
	defer db.Close()

	// Logging the mode server is starting
	log.Printf("Server starting in %s mode at address: %s", sets.Env, ":8080\n\n")

	// Init routes
	routes := routers.Init()

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	err = manners.ListenAndServe(":8080", routes)
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
	defer manners.Close()

	defer log.Println("God bye!")
}
