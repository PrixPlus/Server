package tests

import (
	"errors"
	"github.com/prixplus/server/auth"

	"github.com/prixplus/server/database"
	"github.com/prixplus/server/models"
)

// Temporary table schemas
var schemas = []string{
	`CREATE TEMP TABLE users (id serial, password text, email text)`,
}

// Temporary test entities
var (
	LoginTest = &models.Login{Email: "test@test.com", Password: "123456"}
	// Attention: User.Id will be fulfill when inserted
	UserTest = &models.User{Email: "test@test.com", Password: "$2a$10$tisC/yatxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"}
	// Initializes in InitTestDB
	TokenTest = &models.Token{}
)

// Create temporary schemas in DB and insert some tests entities
func InitData() error {

	// Creates TEMP schema
	err := createTestSchema()
	if err != nil {
		return errors.New("Error creating schema: " + err.Error())
	}

	// Populating our debugging server
	err = insertTestEntityies()
	if err != nil {
		return errors.New("Error populating schema: " + err.Error())
	}

	return nil
}

// Creating temporary test schemas
func createTestSchema() error {

	db, err := database.Get()
	if err != nil {
		return err
	}

	for _, q := range schemas {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

// Inserting temporary test entities
func insertTestEntityies() error {

	// Adding UserTest
	err := UserTest.Insert(nil)
	if err != nil {
		return err
	}

	TokenTest, err = auth.NewToken(*UserTest)
	if err != nil {
		return err
	}

	return nil
}
