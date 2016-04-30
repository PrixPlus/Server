// Create temporary schemas in DB and insert some tests entities
package tests

import (
	"github.com/pkg/errors"

	"github.com/prixplus/server/database"
)

// Creating temporary test schemas
func CreateTempTables() error {
	for _, sql := range testTables {
		_, err := database.Exec(sql)
		if err != nil {
			return errors.Wrap(err, "Error creating temporary tables")
		}
	}
	return nil
}

func TruncateTempTables() error {
	for table, _ := range testTables {
		_, err := database.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			return errors.Wrap(err, "Error truncating temporary tables")
		}
	}
	return nil
}

// Inserting all test entities
func InsertTestEntities() error {
	for _, e := range testEntities {
		err := e.Insert(nil)
		if err != nil {
			return errors.Wrapf(err, "Error inserting temporary entity: %#v", e)
		}
	}
	return nil
}
