package handler

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/model"
	"log"
	"net/http"
	"time"
)

// Duplicate in Middleware and Handler
// It should be in model Auth ?
const (
	relm             = "Prix"
	signingAlgorithm = "HS256"
	timeout          = time.Hour * 24 * 30 // Stay logged a month
)

var secretKey = []byte("7hE Pr!x V3ry 53CRE7 K3Y 7h47 nO0N3 kN0w5!")

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"email": "EMAIL@EMAIL", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func Login(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var user model.User

		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		if len(user.Email) == 0 || len(user.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))
			return
		}

		u := model.User{Email: user.Email, Password: user.Password}
		err = u.Get(db)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error getting user: "+err.Error()))
			return
		}

		// Create the token
		token := jwt.New(jwt.GetSigningMethod(signingAlgorithm))

		log.Printf("### USER ID: %v\n", u.Id)

		expire := time.Now().Add(timeout)
		token.Claims["id"] = u.Id
		token.Claims["exp"] = expire.Unix()

		// I could use some key id to identify what secret key are we using
		// but it is optional and isn't utilized in this package
		// token.Header["kid"] = "Id of the Secret Key used to encrypt this token"

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error creating new token: "+err.Error()))
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
		id, ok := c.Get("id")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not logged"))
			return
		}

		// Need to verify if this token still valid
		// because an user could close this session intentionaly

		// Create the token
		newToken := jwt.New(jwt.GetSigningMethod(signingAlgorithm))

		expire := time.Now().Add(timeout)
		newToken.Claims["id"] = id
		newToken.Claims["exp"] = expire.Unix()

		tokenString, err := newToken.SignedString(secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error creating new refresh token: "+err.Error()))
			return
		}

		// Save Token refresh in DB...

		c.JSON(http.StatusOK, gin.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}
