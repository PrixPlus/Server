package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prixplus/server/models"
	"github.com/stretchr/testify/suite"
)

type ProductSuite struct {
	TestSuite
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}

// Tests [GET] /api/products/:ID method using testProduct
func (t *ProductSuite) xTestGetProduct() {

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/products/%d", testProduct.Id), nil)
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	products, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(products, 1, "not returned just 1 product")
}

// Tests [GET] /api/products
func (t *ProductSuite) xTestGetProductList() {

	req, err := http.NewRequest("GET", "/api/products", nil)
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	products, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(products, testProductsLen, "not returned just %d products", testProductsLen)
}

// Tests [POST] /api/products using a brand new product
func (t *ProductSuite) xTestCreateProduct() {

	// Creating a new user using this email and pass
	product := &models.Product{Gtin: "0000123456789", Description: "BRAND NEW PRODUCT", Thumbnail: "https://s3.amazonaws.com/pictures/products/123456789/kspzwgow", PriceAvg: 1.52, PriceMax: 2.31, PriceMin: 1.10}

	body, err := json.Marshal(product)
	t.NoError(err)

	req, err := http.NewRequest("POST", "/api/products", bytes.NewReader(body))
	t.NoError(err)

	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusCreated, resp.Code, "response code should be Created (201). Body: %s", string(resp.Body.Bytes()))

	// Location should point to the created content: /api/users/4 (if Product.Id=4)
	t.Require().Regexp(`\/api\/products\/\d+`, resp.Header().Get("Location"), "location header should return the adress to retrieve the new content")

	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	products, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(products, 1, "not returned just 1 product")
}

// Tests [PUT] /api/products/:id using testToken and modifying the testProduct
func (t *ProductSuite) TestModifyProduct() {

	modifiedProduct := &models.Product{}

	*modifiedProduct = *testProduct

	modifiedProduct.Description = "Modified description"

	body, err := json.Marshal(modifiedProduct)
	t.NoError(err)

	// Test Refresh Token!
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/products/%d", modifiedProduct.Id), bytes.NewReader(body))
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	products, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(products, 1, "not returned just 1 product")
	t.Require().Equal(modifiedProduct, products[0], "product modfied should be equals to the product sent")
}

/*
// TODO !!!

// Tests [DELET] /api/products/:id using testToken and testProduct
func (t *ProductSuite) TestDeletProduct() {

	// Test Refresh Token!
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/products/%d", testProduct.Id), nil)
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.Product
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	products, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(products, 1, "not returned just 1 product")
	t.Require().Equal(modifiedProduct, products[0], "product modfied should be equals to the product sent")
}
*/
