package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prixplus/server/model"
)

const (
	user     = "postgres"
	password = "pass"
	host     = "prix.plus"
	dbname   = "prix"
	sslmode  = "disable"
)

func InitDB() (*sql.DB, error) {

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		user, password, host, dbname, sslmode)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	// Testing DB connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// If we are in debuggin
	// so we will create temporary tables
	// and insert some dub data
	if !gin.IsDebugging() {
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
