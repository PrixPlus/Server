package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/prixplus/server/model"
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
	// and insert some dub data
	if sets.IsProduction() {
		return db, nil
	}

	// Creates TEMP schema
	err = createSchema(db)
	if err != nil {
		return nil, errors.New("Error creating schema: " + err.Error())
	}

	// Populating our debugging server
	err = populatingSchema(db)
	if err != nil {
		return nil, errors.New("Error populating schema: " + err.Error())
	}

	return db, nil
}

func Get() (*sql.DB, error) {
	if db == nil {
		return connect()
	}
	return db, nil
}

func createSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TEMP TABLE users (id serial, password text, email text)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func populatingSchema(db *sql.DB) error {

	// Adding Prixty user
	u := model.User{Password: "pass", Email: "user@prix.plus"}
	err := u.Insert(db)
	if err != nil {
		return err
	}

	return nil
}

// NOT USED
// JUST FOR EXAMPLE
// Testing some selects
/*
	a.Email = "dub@wars"

	err = a.Update(tx)
	if err != nil {
		return err
	}

	a2 := model.User{Email: "admin@admin"}

	err = a2.Get(tx)
	if err != nil {
		return err
	}

	err = a2.Delete(tx)
	if err != nil {
		return err
	}
*/
