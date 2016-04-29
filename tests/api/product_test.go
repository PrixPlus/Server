package api_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/prixplus/server/models"
	. "gopkg.in/check.v1"
)

// Get products from a given example
func getProduct(product *models.Product, token *models.Token, c *C) *models.Product {

	// Test Refresh Token!
	req, err := http.NewRequest("GET", fmt.Sprintf("/api/products/%d", product.Id), nil)
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving products from response
	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	products, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(products, HasLen, 1) // Test if return just one product in list

	return products[0]
}

// Get products from a given example
func getProductList(q string, token *models.Token, c *C) []*models.Product {

	// Test Refresh Token!
	req, err := http.NewRequest("GET", "/api/products", nil)
	req.URL.Query().Add("q", q)
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)

	fmt.Printf("REQ: %#v\n", req.URL.RawQuery)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving products from response
	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	products, ok := data["results"]
	c.Assert(ok, Equals, true)
	//c.Assert(products, HasLen, 1) // Test if return just one product in list

	return products
}

// Updates product states
func putProduct(product *models.Product, token *models.Token, c *C) *models.Product {

	body, err := json.Marshal(product)
	c.Assert(err, IsNil)

	// Test Refresh Token!
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/products/%d", product.Id), bytes.NewReader(body))
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK, Commentf("Response: %s", resp.Body.String()))

	// Retrieving product from response
	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	products, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(products, HasLen, 1) // Test if return just one product in list

	return products[0]
}

// Creates a new product
func postProduct(product *models.Product, token *models.Token, c *C) *models.Product {

	body, err := json.Marshal(product)
	c.Assert(err, IsNil)

	// Trying to create a new Product
	req, err := http.NewRequest("POST", "/api/products", bytes.NewReader(body))
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusCreated)
	// Location should point to the created content: /api/products/2 (if Product.Id=2)
	c.Assert(resp.Header().Get("Location"), Matches, `\/api\/products\/\d+`, Commentf("Locatioon doesn't matches: %s", resp.Header().Get("Location")))

	// Retrieving Product from response
	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	products, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(products, HasLen, 1) // Test if return just one product in list

	return products[0]
}

// Creates an existing product
func postProductMustConflict(product *models.Product, c *C) {
	// Trying to create a new Product with Gtin
	// server should return StatusConflict
	body, err := json.Marshal(product)
	c.Assert(err, IsNil)
	req, err := http.NewRequest("POST", "/api/products", bytes.NewReader(body))
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusConflict)
}
