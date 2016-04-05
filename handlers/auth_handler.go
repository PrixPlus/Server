package handlers

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/models"
	"net/http"
	"time"
)

// Duplicate in Middleware and Handler
// It should be in model Auth
const (
	relm             = "Prix"
	signingAlgorithm = "HS256"
	timeout          = time.Hour
)

var secretKey = []byte("7hE Pr!x V3ry 53CRE7 K3Y 7h47 nO0N3 kN0w5!")

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"email": "EMAIL@EMAIL", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func Login(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var login model.User

		// Testing Login Error Mesages!
		c.Error(errors.New("Testing these cute error messages"))
		c.Error(errors.New("I love it!"))

		err := c.BindJSON(&login)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))
			return
		}

		// One transation just to get the user? Lol...
		// It's because the method User.Get() requires an transation
		tx, err := db.Begin()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user := model.User{Email: login.Email, Password: login.Password}
		err = user.Get(tx)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		err = tx.Commit()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Create the token
		token := jwt.New(jwt.GetSigningMethod(signingAlgorithm))

		// To add more data in token
		// not used because UserName is alread fine at this momment
		/*
			if mw.PayloadFunc != nil {
				for key, value := range PayloadFunc(user.UserName) {
					token.Claims[key] = value
				}
			}
		*/

		expire := time.Now().Add(timeout)
		token.Claims["username"] = user.UserName
		token.Claims["exp"] = expire.Unix()

		// I could use some key id to identify whay secret key are we using
		// but it is optional and isn't utilized in this package
		// token.Header["kid"] = "Identify My Secret Key Some Way"

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the AuthMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func Refresh(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName, ok := c.Get("userName")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not logged"))
			return
		}

		// Create the token
		newToken := jwt.New(jwt.GetSigningMethod(signingAlgorithm))

		expire := time.Now().Add(timeout)
		newToken.Claims["username"] = userName
		newToken.Claims["exp"] = expire.Unix()

		tokenString, err := newToken.SignedString(secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Save Token refresh in DB

		c.JSON(http.StatusOK, gin.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
