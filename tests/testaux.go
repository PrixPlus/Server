// Create temporary schemas in DB and insert some tests entities
package tests

import "github.com/prixplus/server/database"

// Creating temporary test schemas
func CreateTempTables() error {
	for _, sql := range testTables {
		_, err := database.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func TruncateTempTables() error {
	for table, _ := range testTables {
		_, err := database.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			return err
		}
	}
	return nil
}

// Inserting all test entities
func InsertTestEntities() error {
	for _, e := range testEntities {
		err := e.Insert(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
