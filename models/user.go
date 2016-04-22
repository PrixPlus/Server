package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/prixplus/server/database"

	"github.com/prixplus/server/errs"
)

type User struct {
	Id       int64  `json:"id,string"`          // Send as a string
	Password string `json:"password,omitempty"` // Omitted if empty
	Email    string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %v>", u.Id, u.Email)
}

func (u User) Delete(tx *sql.Tx) error {
	query := "DELETE FROM users WHERE id=$1"
	stmt, err := database.Prepare(query, tx)
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
	query := "INSERT INTO users(password, email) VALUES($1,$2) RETURNING id"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	err = stmt.QueryRow(u.Password, u.Email).Scan(&u.Id)
	if err != nil {
		return err
	}

	log.Printf("Inserted User %s\n", u)

	return nil
}

// Update user in database
func (u User) Update(tx *sql.Tx) error {
	query := "UPDATE users SET email=$1, password=$2 WHERE id=$3"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Email, u.Password, u.Id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect != 1 {
		return errors.New(fmt.Sprintf("%d rows affected in UPDATE to User.Id %d", affect, u.Id))
	}

	log.Printf("Updated User %s\n", u)

	return nil
}

// This method should return just one Elem or an error
// You can get any combination of the fields
func (u *User) Get(tx *sql.Tx) error {
	query := "SELECT id, password, email FROM users WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR password=$2) AND " +
		"($3='' OR email=$3)"
	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(u.Id, u.Password, u.Email)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&u.Id, &u.Password, &u.Email)
		if err != nil {
			return err
		}
	} else {
		// User not found, clear the reference
		*u = User{}
		return errs.ElementNotFound
	}

	// Check if this Elem returned is not unique
	if rows.Next() {
		*u = User{}
		return errors.New("Element not unique")
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	log.Printf("Geted User %s\n", u)

	return nil
}

// This method should return all Elements in db
// equals to the Elem given
func (u *User) GetAll(tx *sql.Tx) ([]User, error) {
	query := "SELECT id, password, email FROM users WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR password=$2) AND " +
		"($3='' OR email=$3)"

	users := []User{}

	stmt, err := database.Prepare(query, tx)
	if err != nil {
		return users, err
	}

	rows, err := stmt.Query(u.Id, u.Password, u.Email)
	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Id, &u.Password, &u.Email)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		return users, err
	}

	log.Printf("Geted %d users like %s\n", len(users), u)

	return users, nil
}
