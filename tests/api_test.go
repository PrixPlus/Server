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

type LoginSuite struct {
	TestSuite
}

func TestLoginSuite(t *testing.T) {
	fmt.Println("### Running Login Suit")
	suite.Run(t, new(LoginSuite))
}

// Testing Refresh Token
func (t *LoginSuite) TestCreateUser() {

	// Creating a new user
	login := &models.Login{Email: "TestLoginAndRefreshToken@email.com", Password: "123456"}

	body, err := json.Marshal(login)
	t.Nil(err, "Err encoding login")

	// Trying to create a new User
	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(body))

	t.Nil(err, "Err requesting api")

	resp := httptest.NewRecorder()

	t.router.ServeHTTP(resp, req)
	//t.Require().Equal(http.StatusCreated, resp.Code, "Response code error: %s", string(resp.Body.Bytes()))

	// Location should point to the created content: /api/users/2 (if User.Id=2)
	t.Require().Regexp(`\/api\/users\/\d+`, resp.Header().Get("Location"), "Location header doesn't match")

	// Retrieving User from response
	var data map[string][]*models.User
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.Nil(err, "Err decoding users from response")

	users, ok := data["results"]
	t.Equal(ok, true, "Err results not found in response")
	t.Len(users, 1, "Err not returned just 1 user")

}
