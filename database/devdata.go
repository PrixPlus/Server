package database

import (
	"database/sql"
	"errors"

	"github.com/prixplus/server/models"
)

// Create Dev Data
func createDevData(db *sql.DB) error {

	// Creates TEMP schema
	err := createSchema(db)
	if err != nil {
		return errors.New("Error creating schema: " + err.Error())
	}

	// Populating our debugging server
	err = populatingSchema(db)
	if err != nil {
		return errors.New("Error populating schema: " + err.Error())
	}

	return nil
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
	u := models.User{Password: "123456", Email: "user@prix.plus"}
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

	a2 := models.User{Email: "admin@admin"}

	err = a2.Get(tx)
	if err != nil {
		return err
	}

	err = a2.Delete(tx)
	if err != nil {
		return err
	}
*/
