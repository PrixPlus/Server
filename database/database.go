package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/prixplus/server/settings"
)

// DB Singleton
var db *sql.DB

func connect() (*sql.DB, error) {

	sets, err := settings.Get()

	log.Println("Connecting to database", sets.DB.Name)

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		sets.DB.User, sets.DB.Password, sets.DB.Host, sets.DB.Name, sets.DB.SSLMode)

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	// Testing DB connection
	err = db.Ping()
	if err != nil {
		log.Println("Error when pinging DB:", err)
		return nil, err
	}

	return db, nil
}

func Get() (*sql.DB, error) {
	if db == nil {
		return connect()
	}
	return db, nil
}

func Close() error {
	err := db.Close()
	if err != nil {
		return err
	}
	db = nil
	return nil
}

// Return a statement for the given query
// using the transaction if provided
func Prepare(query string, tx *sql.Tx) (*sql.Stmt, error) {

	// Transaction provided
	if tx != nil {
		return tx.Prepare(query)
	}

	db, err := Get()
	if err != nil {
		return nil, err
	}

	return db.Prepare(query)
}
