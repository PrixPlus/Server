package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prixplus/server/models"
	"github.com/stretchr/testify/suite"
)

type AuthSuite struct {
	TestSuite
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

// Tests Login method using testLogin
func (t *AuthSuite) TestGetToken() {

	body, err := json.Marshal(testLogin)
	t.NoError(err)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewReader(body))
	t.NoError(err)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "response code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string]*models.Token
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	token, ok := data["token"]
	t.Require().Equal(ok, true, "token not found in response")
	t.Require().NotNil(token, "token should not be nil")
	t.Require().NotEmpty(token.Raw, "token raw should not be empty")
}

// Tests Refresh Token method using testToken
func (t *AuthSuite) TestRefreshToken() {

	req, err := http.NewRequest("GET", "/api/refresh_token", nil)
	t.NoError(err)
	req.Header.Add("Authorization", "Bearer "+testToken.Raw)
	resp := httptest.NewRecorder()
	t.router.ServeHTTP(resp, req)
	t.Require().Equal(http.StatusOK, resp.Code, "code should be OK (200). Body: %s", string(resp.Body.Bytes()))

	var data map[string]*models.Token
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	t.NoError(err)

	token, ok := data["token"]
	t.Require().Equal(ok, true, "token not found in response")
	t.Require().NotNil(token, "token should not be nil")
	t.Require().NotEmpty(token.Raw, "token raw should not be empty")
}
