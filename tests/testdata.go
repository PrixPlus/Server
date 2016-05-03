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
var testTokenExpire, _ = time.Parse(time.RFC3339Nano, "2500-01-01T00:00:00.000000000-03:00")

// Token for the testUser valid for 500 years
var testToken = &models.Token{Raw: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI1MjM2NDAwLCJ1aWQiOjF9.0oJlPNyM5LXeRUK89lLBFjAzbehqpzXFKu6bsM5GXGBIw7bTnkcOGcq530T3sWXnrKAp7qV983Cvu5syyfpCCQ", Expire: testTokenExpire}

// User Email: "test@test.com", Password: "123456" (in hash)
var testUser = &models.User{Email: "test@test.com", Password: "$2a$10$tisC/y*atxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"}

var testUsers = []models.Inserter{
	testUser,
	&models.User{Email: "foo@bar.com", Password: "$2a$10$tisC/y*atxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"},
	&models.User{Email: "prix@plus.com", Password: "$2a$10$tisC/y*atxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"},
}

var testProduct = &models.Product{Gtin: "7894900700046", Description: "REFRIGERANTE COCA COLA LATA ZERO 350ML", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/7894900700220/zhksxcua", PriceAvg: 4.53, PriceMax: 6.30, PriceMin: 3.90}

var testProducts = []models.Inserter{
	testProduct,
	&models.Product{Gtin: "7892840222949", Description: "SALGADINHO FANDANGOS QUEIJO 175G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/fandangos-presunto-elma-chips-200-g_300x300-PU6a3eb_1.jpg", PriceAvg: 21.10, PriceMax: 22.10, PriceMin: 20.10},
	&models.Product{Gtin: "7891095100934", Description: "PIPOCA PARA MICRO-ONDAS YOKI MANTEIGA 100G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/doritos-queijo-nacho-elma-chips-110-g_600x600-PU6a203_1.jpg", PriceAvg: 1.89, PriceMax: 1.50, PriceMin: 2.10},
}

// Test entities lists to be saved in development or test
var testEntities = [][]models.Inserter{
	testUsers,
	testProducts,
}
