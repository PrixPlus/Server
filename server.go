package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/prixplus/server/errs"

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
		errs.LogError(errors.Wrap(err, "Error getting settings"))
		return
	}

	// Init DB singleton connection
	db, err := database.Get()
	if err != nil {
		errs.LogError(errors.Wrap(err, "Error getting database"))
		return
	}

	// Close DB when main() returns
	defer db.Close()

	// Creating temporary schemas and insert some tests datas
	// do it whenever server isn't in production
	if !sets.IsProduction() {
		err = tests.CreateTempTables()
		if err != nil {
			errs.LogError(errors.Wrap(err, "Error creating temporary schemas"))
			return
		}
		err = tests.InsertTestEntities()
		if err != nil {
			errs.LogError(errors.Wrap(err, "Error creating tests entities"))
			return
		}
	}

	// Logging the mode server is starting
	fmt.Printf("Server starting in %s mode at address: %s", sets.Env, ":8080\n\n")

	// Init router
	router, err := routers.Init()
	if err != nil {
		errs.LogError(errors.Wrap(err, "Error initializing router"))
		return
	}

	// Shut the server down gracefully
	processStopedBySignal(db)

	// Manners allows you to shut your Go webserver down gracefully, without dropping any requests
	err = manners.ListenAndServe(":8080", router)
	if err != nil {
		errs.LogError(errors.Wrap(err, "Error starting server"))
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
