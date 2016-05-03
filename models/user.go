package models

import (
	"encoding/json"

	"github.com/prixplus/server/db"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
	"github.com/prixplus/server/errs"
)

type User struct {
	Id       int64  `json:"id,string"`          // Send as a string
	Password string `json:"password,omitempty"` // Omitted if empty
	Email    string `json:"email"`
}

func (u User) String() string {
	s, err := json.Marshal(u)
	if err != nil { // Just log the error
		errs.LogError(errors.Wrap(err, "encoding json"))
	}

	return string(s)
}

func (u User) Delete(tx *sqlx.Tx) error {
	query := "DELETE FROM users WHERE id=:id"
	stmtx, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	_, err = stmtx.Exec(u)
	if err != nil {
		return errors.Wrap(err, "executing named query")
	}

	// fmt.Printf("User deleted %s\n", u)
	return nil
}

func (u *User) Insert(tx *sqlx.Tx) error {
	query := "INSERT INTO " +
		"users(password, email) " +
		"VALUES(:password, :email) " +
		"RETURNING id"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	err = stmt.QueryRow(u).StructScan(u)
	if err != nil {
		return errors.Wrap(err, "scanning named query")
	}

	// fmt.Printf("User inserted %s\n", u)
	return nil
}

// Update user in database
func (u User) Update(tx *sqlx.Tx) error {
	query := "UPDATE users SET email=:email, password=:password WHERE id=:id"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	_, err = stmt.Exec(u)
	if err != nil {
		return errors.Wrap(err, "executing named query")
	}

	// fmt.Printf("User updated %s\n", u)
	return nil
}

// This method should return just one Elem or an error
// You can get any combination of the fields
func (u *User) Get(tx *sqlx.Tx) error {
	query := "SELECT * FROM users WHERE " +
		"(:id=0 OR id=:id) AND " +
		"(:password='' OR password=:password) AND " +
		"(:email='' OR email=:email)"
	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return errors.Wrap(err, "preparing named query")
	}

	err = stmt.Get(u, u)
	if err != nil {
		return errors.Wrap(err, "getting the user")
	}

	// fmt.Printf("User geted %s\n", u)
	return nil
}

// This method should return all Elements in db
// equals to the Elem given
func (u *User) GetAll(tx *sqlx.Tx) ([]User, error) {
	query := "SELECT * FROM users WHERE " +
		"(:id=0 OR id=:id) AND " +
		"(:password='' OR password=:password) AND " +
		"(:email='' OR email=:email)"

	users := []User{}

	stmt, err := db.PrepareNamed(query, tx)
	if err != nil {
		return users, errors.Wrap(err, "preparing named query")
	}

	err = stmt.Select(&users, u)
	if err != nil {
		return users, errors.Wrap(err, "selecting users")
	}
	// fmt.Printf("%d Users geted like %s\n", len(users), u)
	return users, nil
}
