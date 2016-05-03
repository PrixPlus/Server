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

type UserSuite struct {
	TestSuite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

// Tests Get Me method using testToken
func (t *UserSuite) TestGetMe() {

	req, err := http.NewRequest("GET", "/api/me", nil)
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	users, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(users, 1, "not returned just 1 user")
}

// Tests [POST] User using a brand new user
func (t *UserSuite) TestCreateUser() {

	// Creating a new user using this email and pass
	login := &models.Login{Email: "brandnewuser@email.com", Password: "123456"}

	body, err := json.Marshal(login)
	t.NoError(err)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	t.NoError(err)

	resp := httptest.NewRecorder()

	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusCreated, resp.Code, "response code should be Created (201). Body: %s", string(resp.Body.Bytes()))

	// Location should point to the created content: /api/users/2 (if User.Id=2)
	t.Require().Regexp(`\/api\/users\/\d+`, resp.Header().Get("Location"), "location header should return the adress to retrieve the inserted user")

	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	users, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(users, 1, " not returned just 1 user")
}

// Tests [POST] User using testLogin from the already created testUser
// this method should return Conflict status with a message
func (t *UserSuite) TestCreateConflict() {
	// Using testLogin from testUser already in use
	body, err := json.Marshal(testLogin)
	t.NoError(err)

	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	t.NoError(err)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusConflict, resp.Code, "response code should be Conflict (409). Body: %s")

	var data map[string][]string
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	messages, ok := data["messages"]
	t.Require().Equal(ok, true, "messages not found in response")
	t.Require().Len(messages, 1, "not returned just 1 message")
}

// Tests [PUT] User using testToken and modifying the testUser
func (t *UserSuite) TestModifyUser() {

	modifiedUser := &models.User{}

	*modifiedUser = *testUser

	modifiedUser.Email = "modifiedEmail@test.com"
	modifiedUser.Password = "modifiedPass"

	body, err := json.Marshal(modifiedUser)
	t.NoError(err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/api/users/%d", modifiedUser.Id), bytes.NewReader(body))
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	users, ok := data["results"]
	t.Require().Equal(ok, true, "results not found in response")
	t.Require().Len(users, 1, " not returned just 1 user")
	// Cleanning password before the evaluation
	// Because user's password isn't returned
	modifiedUser.Password = ""
	t.Require().Equal(modifiedUser, users[0], "user modfied should be equals to the user sent")

}
