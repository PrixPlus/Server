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

// Get the user from the current session
func getMe(token *models.Token, c *C) *models.User {

	// Test Refresh Token!
	req, err := http.NewRequest("GET", "/api/me", nil)
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving User from response
	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	users, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(users, HasLen, 1) // Test if return just one user in list

	return users[0]
}

// Updates user states
func putUser(user *models.User, token *models.Token, c *C) *models.User {

	body, err := json.Marshal(user)
	c.Assert(err, IsNil)

	// Test Refresh Token!
	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/users/%d", user.Id), bytes.NewReader(body))
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving User from response
	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	users, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(users, HasLen, 1) // Test if return just one user in list

	return users[0]
}

// Creates a new user
func postUser(login *models.Login, c *C) *models.User {

	body, err := json.Marshal(login)
	c.Assert(err, IsNil)

	// Trying to create a new User
	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusCreated)
	// Location should point to the created content: /api/users/2 (if User.Id=2)
	c.Assert(resp.Header().Get("Location"), Matches, `\/api\/users\/\d+`, Commentf("Locatioon doesn't matches: %s", resp.Header().Get("Location")))

	// Retrieving User from response
	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	users, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(users, HasLen, 1) // Test if return just one user in list

	return users[0]
}

// Creates a new user
func postUserMustConflict(login *models.Login, c *C) {
	// Trying to create a new User with same email
	// server should return StatusConflict
	body, err := json.Marshal(login)
	c.Assert(err, IsNil)
	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	c.Assert(err, IsNil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusConflict)
}
