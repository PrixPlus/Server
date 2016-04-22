package api_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/prixplus/server/models"

	. "gopkg.in/check.v1"
)

func getToken(login *models.Login, c *C) *models.Token {

	body, err := json.Marshal(login)
	c.Assert(err, IsNil)

	// Trying to login with our LoginTest from UserTest
	req, err := http.NewRequest("POST", "/api/login", bytes.NewReader(body))
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving token from response
	var data map[string]*models.Token
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	token, ok := data["token"]
	c.Assert(ok, Equals, true)
	c.Assert(token, NotNil)

	return token
}

func refreshToken(token *models.Token, c *C) *models.Token {

	// Trying to get a new Token
	req, err := http.NewRequest("GET", "/api/refresh_token", nil)
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving token from response
	var data map[string]*models.Token
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	token, ok := data["token"]
	c.Assert(ok, Equals, true)
	c.Assert(token, NotNil)

	return token
}
