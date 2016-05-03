// Create temporary schemas in DB and insert some tests entities
package tests

import (
	"github.com/pkg/errors"

	"github.com/prixplus/server/db"
)

// Creating temporary test schemas
func CreateTempTables() error {
	for _, sql := range testTables {
		_, err := db.Exec(sql)
		if err != nil {
			return errors.Wrap(err, "creating temporary tables")
		}
	}
	return nil
}

func DropTempTablesIfExist() error {
	for table, _ := range testTables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE")
		if err != nil {
			return errors.Wrap(err, "removing temporary tables")
		}
	}
	return nil
}

// Inserting all test entities
func InsertTestEntities() error {
	for _, e := range testEntities {
		err := e.Insert(nil)
		if err != nil {
			return errors.Wrapf(err, "inserting temporary entity: %#v", e)
		}
	}
	return nil
}
