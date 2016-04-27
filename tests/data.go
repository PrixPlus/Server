// Create temporary schemas in DB and insert some tests entities
package tests

import (
	"github.com/prixplus/server/auth"

	"github.com/prixplus/server/database"
	"github.com/prixplus/server/models"
)

// Temporary table schemas
var schemas = []string{
	`CREATE TEMP TABLE users (id serial, password text, email text)`,
	`CREATE TEMP TABLE products (id serial, gtin text, description text, thumbnail text, price real, priceavg real, pricemax real, pricemin real)`,
}

// Temporary test entities
var (
	LoginTest = &models.Login{Email: "test@test.com", Password: "123456"}
	// Attention: User.Id will be fulfill when inserted
	UserTest = &models.User{Email: "test@test.com", Password: "$2a$10$tisC/yatxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"}
	// Initializes in InitTestDB
	TokenTest = &models.Token{}
	// The first test product
	ProductTest1 = &models.Product{Gtin: "7894900700046", Description: "REFRIGERANTE COCA COLA LATA ZERO 350ML", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/7894900700220/zhksxcua", Price: 4.53, PriceMax: 6.30, PriceMin: 3.90}
	ProductTest2 = &models.Product{Gtin: "0789840233945", Description: "FANDANGOS SALGADINHO DE MILHO SABOR PRESUNTO (HAM FLAVOR SNACK) BAG 63G -BRAZIL", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/789840233945", Price: 21.10, PriceMax: 22.10, PriceMin: 20.10}
	ProductTest3 = &models.Product{Gtin: "7892840211240", Description: "SALGADINHO DORITOS QUEIJO NACHO 110G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/doritos-queijo-nacho-elma-chips-110-g_600x600-PU6a203_1.jpg", Price: 4.53, PriceMax: 6.30, PriceMin: 3.90}
)

// Creating temporary test schemas
func CreateTempTables() error {

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
func InsertTestEntityies() error {

	// Adding UserTest
	err := UserTest.Insert(nil)
	if err != nil {
		return err
	}

	TokenTest, err = auth.NewToken(*UserTest)
	if err != nil {
		return err
	}

	// Adding ProductTest
	err = ProductTest1.Insert(nil)
	if err != nil {
		return err
	}
	// Adding ProductTest
	err = ProductTest2.Insert(nil)
	if err != nil {
		return err
	}
	// Adding ProductTest
	err = ProductTest3.Insert(nil)
	if err != nil {
		return err
	}

	return nil
}
