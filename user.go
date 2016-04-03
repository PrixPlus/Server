package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type User struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.UserName, u.Email)
}

func (u User) Delete(tx *sql.Tx) error {
	stmt, err := tx.Prepare("DELETE FROM users WHERE id=$1")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect != 1 {
		return errors.New(fmt.Sprintf("%d rows affected in DELETE to User.Id %s", affect, u.Id))
	}

	log.Printf("Deleted User %s\n", u)

	return nil
}

func (u *User) Insert(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO users(username, password, email) VALUES($1,$2,$3) RETURNING id")
	if err != nil {
		return err
	}

	err = stmt.QueryRow(u.UserName, u.Password, u.Email).Scan(&u.Id)
	if err != nil {
		return err
	}

	log.Printf("Inserted User %s\n", u)

	return nil
}

func (u User) Update(tx *sql.Tx) error {
	stmt, err := tx.Prepare("UPDATE users SET username=$1, password=$2, email=$3 WHERE id=$4")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.UserName, u.Password, u.Email, u.Id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect != 1 {
		return errors.New(fmt.Sprintf("%d rows affected in UPDATE to User.Id %s", affect, u.Id))
	}

	log.Printf("Updated User %s\n", u)

	return nil
}

// This method should return just one Elem or an error
// You can get any combination of the fields
func (u *User) Get(tx *sql.Tx) error {
	stmt, err := tx.Prepare("SELECT id, username, password, email FROM users WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR username=$2) AND " +
		"($3='' OR password=$3) AND " +
		"($4='' OR email=$4)")
	if err != nil {
		return err
	}

	rows, err := stmt.Query(u.Id, u.UserName, u.Password, u.Email)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&u.Id, &u.UserName, &u.Password, &u.Email)
		if err != nil {
			return err
		}
	} else {
		// User not found, clear the reference
		*u = User{}
		return ErrElemNotFound
	}

	// Check if this Elem returned is not unique
	if rows.Next() {
		*u = User{}
		return ErrElemNotUnique
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	log.Printf("Geted User %s\n", u)

	return nil
}
