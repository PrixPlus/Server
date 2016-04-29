// Create temporary schemas in DB and insert some tests entities
package tests

import "github.com/prixplus/server/models"

// Temporary table schemas
var testTables = map[string]string{
	"users":    `CREATE TEMP TABLE users (id serial, password text, email text)`,
	"products": `CREATE TEMP TABLE products (id serial, gtin text, description text, thumbnail text, priceavg real, pricemax real, pricemin real)`,
}

var testEntities = []models.Inserter{
	//
	// User
	//
	// Email: "test@test.com", Password: "123456" (in hash)
	&models.User{Email: "test@test.com", Password: "$2a$10$tisC/y*atxRhEIPNPAgH.yexTuPpGQ4BRAqsVrGViteXPsPDpe1Mx2"},
	//
	// Products
	//
	&models.Product{Gtin: "7894900700046", Description: "REFRIGERANTE COCA COLA LATA ZERO 350ML", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/7894900700220/zhksxcua", PriceAvg: 4.53, PriceMax: 6.30, PriceMin: 3.90},
	&models.Product{Gtin: "7892840222949", Description: "SALGADINHO FANDANGOS QUEIJO 175G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/fandangos-presunto-elma-chips-200-g_300x300-PU6a3eb_1.jpg", PriceAvg: 21.10, PriceMax: 22.10, PriceMin: 20.10},
	&models.Product{Gtin: "7891095100934", Description: "PIPOCA PARA MICRO-ONDAS YOKI MANTEIGA 100G", Thumbnail: "https://s3.amazonaws.com/bluesoft-cosmos/products/doritos-queijo-nacho-elma-chips-110-g_600x600-PU6a203_1.jpg", PriceAvg: 1.89, PriceMax: 1.50, PriceMin: 2.10},
}
