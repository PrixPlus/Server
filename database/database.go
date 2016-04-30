package database

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
	"github.com/prixplus/server/settings"
)

// DB Singleton
var db *sql.DB

func connect() (*sql.DB, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Erro getting settings")
	}

	fmt.Println("Connecting to database", sets.DB.Name)

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		sets.DB.User, sets.DB.Password, sets.DB.Host, sets.DB.Name, sets.DB.SSLMode)

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening database connection")
	}

	// Testing DB connection
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "Erro pinging database")
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
		return errors.Wrap(err, "Erro closing database")
	}
	db = nil
	return nil
}

// Implementing some DB methods for fast using
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
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
		return nil, errors.Wrap(err, "Erro getting database instance")
	}

	return db.Prepare(query)
}
