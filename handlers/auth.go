package handlers

import (
	"net/http"

	"github.com/pkg/errors"

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
			c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "parsing JSON"))
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "email or password can not be empty"))
			return
		}

		user := models.User{Email: login.Email}
		err = user.Get(nil)
		if err == errs.ElemNotFound {
			c.AbortWithError(errs.Status[err], errors.Wrapf(err, "user not found with email %s", login.Email))
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrapf(err, "getting user with email %s", login.Email))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err == bcrypt.ErrHashTooShort {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "password is too short"))
			return
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			c.Error(errors.New("Password received:" + login.Password))
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "password does not match"))
			return
		}

		token, err := auth.NewToken(user)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrap(err, "creating new token"))
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

		userId, err := auth.GetUserIdFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrapf(err, "user not logged %d"))
			return
		}

		// We just need User.Id to create a new token
		user := models.User{Id: userId}

		token, err := auth.NewToken(user)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.Wrap(err, "creating new token"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
