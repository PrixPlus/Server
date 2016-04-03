package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	user     = "postgres"
	password = "pass"
	host     = "prix.plus"
	dbname   = "prix"
	sslmode  = "disable"
)

var (
	ErrElemNotFound  = errors.New("Element not found")
	ErrElemNotUnique = errors.New("Element not unique")
)

func createSchema(db *sql.DB) error {
	queries := []string{
		`CREATE TEMP TABLE users (id serial, username text, password text, email text)`,
	}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

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

	// Creates TEMP schema
	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	// Testing insert with some users
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	z := User{UserName: "admin", Password: "pass", Email: "admin@admin"}
	err = z.Insert(tx)
	if err != nil {
		return nil, err
	}

	u := User{UserName: "Testing", Password: "pass", Email: "test@test"}
	err = u.Insert(tx)
	if err != nil {
		return nil, err
	}

	u.Email = "dub@wars"

	err = u.Update(tx)
	if err != nil {
		return nil, err
	}

	a := User{Email: "admin@admin"}

	err = a.Get(tx)
	if err != nil {
		return nil, err
	}

	err = u.Delete(tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return db, nil
}
