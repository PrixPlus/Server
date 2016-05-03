package models

import (
	"github.com/jmoiron/sqlx"
)

type Inserter interface {
	Insert(tx *sqlx.Tx) error
}
