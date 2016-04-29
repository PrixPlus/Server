package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/prixplus/server/tests"

	"github.com/braintree/manners"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"
)

func main() {

	fmt.Println()

	// Load singleton settings
	sets, err := settings.Get()
	if err != nil {
		fmt.Println("Error loading settings: ", err)
		return
	}

	// Init DB singleton connection
	db, err := database.Get()
	if err != nil {
		fmt.Println("Error initializing DB: ", err)
		return
	}

	// Close DB when main() returns
	defer db.Close()

	// Creating temporary schemas and insert some tests datas
	// do it whenever server isn't in production
	if !sets.IsProduction() {
		err = tests.CreateTempTables()
		if err != nil {
			fmt.Println("Error creating temporary schemas: ", err)
			return
		}
		err = tests.InsertTestEntities()
		if err != nil {
			fmt.Println("Error creating tests entities: ", err)
			return
		}
	}

	// Logging the mode server is starting
	fmt.Printf("Server starting in %s mode at address: %s", sets.Env, ":8080\n\n")

	// Init router
	router, err := routers.Init()
	if err != nil {
		fmt.Println("Error initializing router: ", err)
		return
	}

	// Shut the server down gracefully
	processStopedBySignal(db)

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	err = manners.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
	defer manners.Close()

}

// Shut the server down gracefully if receive a interrupt signal
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
		fmt.Println("THIS PROCESS IS WAITING SIGNAL TO STOP GRACEFULLY")
		for sig := range c {
			fmt.Println("\n\nSTOPED BY SIGNAL:", sig.String())
			fmt.Println("SHUTTING DOWN GRACEFULLY!")
			fmt.Println("\nGod bye!")
			manners.Close()
			db.Close()
			os.Exit(1)
		}
	}()
}
