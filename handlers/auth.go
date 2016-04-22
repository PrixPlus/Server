package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/auth"
	"github.com/prixplus/server/errs"
	"github.com/prixplus/server/models"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"email": "EMAIL@EMAIL", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var login models.Login

		err := c.BindJSON(&login)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))
			return
		}

		user := models.User{Email: login.Email}
		err = user.Get(nil)
		if err == errs.ElementNotFound {
			c.AbortWithError(errs.Status[err], errors.New("User not found!"))
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error getting user: "+err.Error()))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err == bcrypt.ErrHashTooShort {
			c.AbortWithError(http.StatusBadRequest, errors.New("This password is too short: "+err.Error()))
			return
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			c.Error(errors.New("Password received:" + login.Password))
			c.AbortWithError(http.StatusBadRequest, errors.New("Password does not match: "+err.Error()))
			return
		}

		token, err := auth.NewToken(user)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error creating new token: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the Authmiddlewares.
// Reply will be of the form {"token": "TOKEN"}.
func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, ok := c.Get("id")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not logged"))
			return
		}

		idFloat64, ok := id.(int64) // float64?
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error casting Claims"))
			return
		}
		// We just need User.Id to create a new token
		user := models.User{Id: idFloat64}

		token, err := auth.NewToken(user)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Error creating new token: "+err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
