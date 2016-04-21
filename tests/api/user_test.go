package api_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/prixplus/server/models"

	"github.com/prixplus/server/database"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"
	"github.com/prixplus/server/tests"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type UserTestSuite struct{}

var _ = Suite(&UserTestSuite{})

var router *gin.Engine

// SetUp the test environment
func (s *UserTestSuite) SetUpSuite(c *C) {

	os.Setenv("GO_ENV", "test")

	// Load singleton settings
	_, err := settings.Get()
	if err != nil {
		log.Fatal("Error loading settings: ", err)
		return
	}

	// Init DB singleton connection
	_, err = database.Get()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
		return
	}

	// Routing the API
	router = routers.Init()

	// Creating temporary schemas and insert some tests entities
	tests.InitTest()
}

// When all finishes
func (s *UserTestSuite) TearDownSuite(c *C) {
	// Closing DB singleton connection
	err := database.Close()
	if err != nil {
		log.Fatal("Error closing DB: ", err)
		return
	}
}

// Testing Login, Refresh Token
func (s *UserTestSuite) TestLogin(c *C) {
	getToken(tests.LoginTest, c)
}

// Testing get user from current session
func (s *UserTestSuite) TestGetMe(c *C) {
	token := getToken(tests.LoginTest, c)
	getMe(token, c)
}

// Testing Refresh Token
// this method uses the TokenTest
func (s *UserTestSuite) TestRefreshToken(c *C) {

	token := getToken(tests.LoginTest, c)

	// Test Refresh Token!
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

}

// Testing User Creation
func (s *UserTestSuite) TestPostUser(c *C) {

	login := models.Login{Email: "newuser@email.com", Password: "123456"}

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

	// Trying to create a new User with same email (Should block)
	req, err = http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	c.Assert(err, IsNil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusConflict)

}

// Testing modifying the user
func (s *UserTestSuite) TestPutUser(c *C) {

	token := getToken(tests.LoginTest, c)
	user := getMe(token, c) // Coppying values

	user.Email = "newemail@email.com"

	userModified := putUser(user, token, c)

	c.Assert(userModified, DeepEquals, user)

	// TESTING IF USER HAS REALLY CHANGED IN DB
	me := getMe(token, c)

	// Test if our changes has saved
	c.Assert(userModified, DeepEquals, me)

	// Now backing user to initial state
	user.Email = tests.UserTest.Email

	userModified = putUser(user, token, c)

	c.Assert(userModified, DeepEquals, user)
}

func getMe(token *models.Token, c *C) *models.User {

	// Test Refresh Token!
	req, err := http.NewRequest("GET", "/api/me", nil)
	c.Assert(err, IsNil)
	req.Header.Add("Authorization", "Bearer "+token.Raw)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	c.Assert(resp.Code, Equals, http.StatusOK)

	// Retrieving User from response
	var data map[string][]models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	users, ok := data["results"]
	c.Assert(ok, Equals, true)
	c.Assert(users, HasLen, 1) // Test if return just one user in list

	return &users[0]
}

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
	var data map[string]models.Token
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	c.Assert(err, IsNil)

	token, ok := data["token"]
	c.Assert(ok, Equals, true)
	c.Assert(token, NotNil)

	return &token
}

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
