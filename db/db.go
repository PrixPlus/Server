package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
	"github.com/prixplus/server/settings"
)

type DB struct {
	*sqlx.DB
}

// DB Singleton
var db *sqlx.DB

func connect() (*sqlx.DB, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting settings")
	}

	// fmt.Println("Connecting to database", sets.DB.Name)

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		sets.DB.User, sets.DB.Password, sets.DB.Host, sets.DB.Name, sets.DB.SSLMode)

	db, err = sqlx.Connect("postgres", dbinfo)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to database")
	}

	return db, nil
}

func Get() (*sqlx.DB, error) {
	if db == nil {
		return connect()
	}
	return db, nil
}

func Close() error {
	err := db.Close()
	if err != nil {
		return errors.Wrap(err, "closing database")
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
func Prepare(query string, tx *sqlx.Tx) (*sqlx.Stmt, error) {

	// Use the transaction if provided
	if tx != nil {
		return tx.Preparex(query)
	}

	db, err := Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting database instance")
	}

	return db.Preparex(query)
}

// Return a statement for the given query
// using the transaction if provided
func PrepareNamed(query string, tx *sqlx.Tx) (*sqlx.NamedStmt, error) {

	// Use the transaction if provided
	if tx != nil {
		return tx.PrepareNamed(query)
	}

	db, err := Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting database instance")
	}

	return db.PrepareNamed(query)
}
