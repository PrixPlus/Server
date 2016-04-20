package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/prixplus/server/errs"
)

type User struct {
	Id       int64  `json:"id"`
	Password string `json:"-"` // Not send or receive
	Email    string `json:"email"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %v>", u.Id, u.Email)
}

func (u User) Delete(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM users WHERE id=$1")
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

func (u *User) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users(password, email) VALUES($1,$2) RETURNING id")
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

func (u User) Update(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE users SET password=$1, email=$2 WHERE id=$3")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(u.Password, u.Email, u.Id)
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
func (u *User) Get(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT id, password, email FROM users WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR password=$2) AND " +
		"($3='' OR email=$3)")
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
func (u *User) GetAll(db *sql.DB) ([]User, error) {

	users := []User{}

	stmt, err := db.Prepare("SELECT id, password, email FROM users WHERE " +
		"($1=0 OR id=$1) AND " +
		"($2='' OR password=$2) AND " +
		"($3='' OR email=$3)")
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
