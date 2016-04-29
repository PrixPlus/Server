package models

import "database/sql"

type Inserter interface {
	Insert(tx *sql.Tx) error
}
