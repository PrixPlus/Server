package middlewares

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// Apply for private routes
func Auth(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := parseToken(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Your token is broken: "+err.Error()))
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Ivalid token provided or it's expired"))
			return
		}

		// Check in DB if token user still logged
		// maybe user logout

		userName := token.Claims["username"].(string)
		// Why use it? Dunno...
		// Maybe we will never use it is removed
		// c.Set("Claims", token.Claims)
		c.Set("userName", userName)

		// User will have full acess !?
		// If we need to filter some areas to some user profiles
		// Here we check if this user has this permission
		if false {
			c.AbortWithError(http.StatusForbidden, errors.New("You don't have permission to access"))
			return
		}

		c.Next()
	}
}

func parseToken(c *gin.Context) (*jwt.Token, error) {
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		return nil, errors.New("Auth header empty")
	}

	// The first word of Authorization header should be Bearer
	if len(auth) < 6 || auth[0:7] != "Bearer " {
		return nil, errors.New("Invalid auth header")
	}

	token, err := jwt.Parse(auth[7:], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(signingAlgorithm) != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}

		// I could use some key id to identify whay secret key are we using
		// but it is optional and isn't utilized in this package
		// secretKey, err := myLookupForSecretKey(token.Header["kid"])

		return secretKey, nil
	})

	return token, err
}
