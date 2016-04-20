package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/errs"
	"github.com/prixplus/server/model"
	"github.com/prixplus/server/settings"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"email": "EMAIL@EMAIL", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		sets, err := settings.Get()
		if err != nil {
			log.Fatal("Error getting Settings: ", err)
		}

		db, err := database.Get()
		if err != nil {
			log.Fatal("Error getting DB: ", err)
		}

		var login model.Login

		err = c.BindJSON(&login)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.New("Error parsing JSON: "+err.Error()))
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))
			return
		}

		u := model.User{Email: login.Email}
		err = u.Get(db)
		if err == errs.ElementNotFound {
			c.AbortWithError(errs.Status[err], errors.New("User not found!"))
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error getting user: "+err.Error()))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(u.Password))
		if err == bcrypt.ErrHashTooShort {
			c.Error(errors.New("PASS:" + login.Password + ". PASS2:" + u.Password))
			c.AbortWithError(http.StatusBadRequest, errors.New("This password is too short: "+err.Error()))
			return
		} else if err == bcrypt.ErrMismatchedHashAndPassword {
			c.Error(errors.New("PASS:" + login.Password + ". PASS2:" + u.Password))
			c.AbortWithError(http.StatusBadRequest, errors.New("Password does not match: "+err.Error()))
			return
		}

		// Create the token
		token := jwt.New(jwt.GetSigningMethod(sets.JWT.Algorithm))

		log.Printf("### USER ID: %v\n", u.Id)

		timeout := time.Hour * sets.JWT.Expiration

		expire := time.Now().Add(timeout)
		token.Claims["id"] = u.Id
		token.Claims["exp"] = expire.Unix()

		// I could use some key id to identify what secret key are we using
		// but it is optional and isn't utilized in this package
		// token.Header["kid"] = "Id of the Secret Key used to encrypt this token"

		tokenString, err := token.SignedString([]byte(sets.JWT.SecretKey))
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
func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {

		sets, err := settings.Get()
		if err != nil {
			log.Fatal("Error getting Settings: ", err)
		}

		id, ok := c.Get("id")
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, errors.New("User not logged"))
			return
		}

		// Need to verify if this token still valid
		// because an user could close this session intentionaly

		// Create the token
		newToken := jwt.New(jwt.GetSigningMethod(sets.JWT.Algorithm))

		timeout := time.Hour * sets.JWT.Expiration

		expire := time.Now().Add(timeout)
		newToken.Claims["id"] = id
		newToken.Claims["exp"] = expire.Unix()

		tokenString, err := newToken.SignedString(sets.JWT.SecretKey)
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
