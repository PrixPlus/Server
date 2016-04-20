package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/prixplus/server/settings"
)

// DB Singleton
var db *sql.DB

func connect() (*sql.DB, error) {

	sets, err := settings.Get()

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		sets.DB.User, sets.DB.Password, sets.DB.Host, sets.DB.Name, sets.DB.SSLMode)

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	// Testing DB connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// If we are in development
	// so we will create temporary tables
	// and insert some tests data
	if !sets.IsProduction() {
		err = createDevData(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func Get() (*sql.DB, error) {
	if db == nil {
		return connect()
	}
	return db, nil
}
