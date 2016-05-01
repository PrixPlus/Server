// Create temporary schemas in DB and insert some tests entities
package tests

import (
	"time"

	"github.com/prixplus/server/models"
)

// Temporary table schemas
var testTables = map[string]string{
	"users":    `CREATE TEMP TABLE users (id serial, password text, email text)`,
	"products": `CREATE TEMP TABLE products (id serial, gtin text, description text, thumbnail text, priceavg real, pricemax real, pricemin real)`,
}

// Login for the testUser
var testLogin = &models.Login{Email: "test@test.com", Password: "123456"}

// Expire date for the generated testToken
var testTokenExpire, _ = time.Parse(time.RFC3339Nano, "2114-11-24T03:10:03.951466305-03:00")

// Token for the testUser valid for 100 years
var testToken = &models.Token{Raw: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ1NzI0ODI1MzgsImlkIjoxfQ.VWRWQZJRA_oH8bYCcxTHYwXs6PEvBkNcxSPVrK7Be4uUztzW8IShOWP-wDwScL4PW4fyeGl3bujyRWoPANQFnA", Expire: testTokenExpire}

// User Email: "test@test.com", Password: "123456" (in hash)
var testUser = &models.User{Email: "test@test.com", Password: "$2a$10$tisC/y*atxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"}

// Test product
var testProduct = &models.Product{Gtin: "7894900700046", Description: "REFRIGERANTE COCA COLA LATA ZERO 350ML", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/7894900700220/zhksxcua", PriceAvg: 4.53, PriceMax: 6.30, PriceMin: 3.90}

// Show the number of products inserted
var testProductsLen = 3

// Entities list to be saved
var testEntities = []models.Inserter{
	testUser,
	testProduct,
	// More products
	&models.Product{Gtin: "7892840222949", Description: "SALGADINHO FANDANGOS QUEIJO 175G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/fandangos-presunto-elma-chips-200-g_300x300-PU6a3eb_1.jpg", PriceAvg: 21.10, PriceMax: 22.10, PriceMin: 20.10},
	&models.Product{Gtin: "7891095100934", Description: "PIPOCA PARA MICRO-ONDAS YOKI MANTEIGA 100G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/doritos-queijo-nacho-elma-chips-110-g_600x600-PU6a203_1.jpg", PriceAvg: 1.89, PriceMax: 1.50, PriceMin: 2.10},
}
