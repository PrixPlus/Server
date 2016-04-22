package main

import (
	"database/sql"
	"github.com/prixplus/server/tests"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Creating temporary schemas and insert some tests datas
	// do it whenever server isn't in production
	if !sets.IsProduction() {
		tests.InitData()
	}

	// Logging the mode server is starting
	log.Printf("Server starting in %s mode at address: %s", sets.Env, ":8080\n\n")

	// Init routes
	routes := routers.Init()

	// Shut the server down gracefully
	processStopedBySignal(db)

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	err = manners.ListenAndServe(":8080", routes)
	if err != nil {
		log.Fatal("Error starting server: ", err)
		return
	}
	defer manners.Close()

}

// Shut the server down gracefully if receive a interrupt s
func processStopedBySignal(db *sql.DB) {
	// Stop server if someone kills the process
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	//signal.Notify(c, syscall.SIGSTOP)
	signal.Notify(c, syscall.SIGABRT) // ABORT
	signal.Notify(c, syscall.SIGKILL) // KILL
	signal.Notify(c, syscall.SIGTERM) // TERMINATION
	signal.Notify(c, syscall.SIGINT)  // TERMINAL INTERRUPT (Ctrl+C)
	signal.Notify(c, syscall.SIGSTOP) // STOP
	signal.Notify(c, syscall.SIGTSTP) // TERMINAL STOP (Ctrl+Z)
	signal.Notify(c, syscall.SIGQUIT) // QUIT (Ctrl+\)
	go func() {
		log.Println("THIS PROCESS IS WAITING SIGNAL TO STOP GRACEFULLY")
		for sig := range c {
			log.Println("\n\nSTOPED BY SIGNAL:", sig.String())
			log.Println("SHUTTING DOWN GRACEFULLY!")
			log.Println("\nGod bye!")
			manners.Close()
			db.Close()
			os.Exit(1)
		}
	}()
}
