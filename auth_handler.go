package main

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	relm             = "Prix"
	signingAlgorithm = "HS256"
	timeout          = time.Hour
)

var secretKey = []byte("7hE Pr!x V3ry 53CRE7 K3Y 7h47 nO0N3 kN0w5!")

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"email": "EMAIL@EMAIL", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func LoginHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var login User
		c.Error(errors.New("TESTING"))
		c.Error(errors.New("Fucking shit"))

		err := c.BindJSON(&login)
		if err != nil {
			//unauthorized(c, http.StatusBadRequest, "Something goes wrong: "+err.Error())

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if len(login.Email) == 0 || len(login.Password) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Email or Password can not be empty"))

			return
		}

		panic("FUCK U")

		// One transation just to get the user? Lol...
		// It's because the method User.Get() requires an transation
		tx, err := db.Begin()
		if err != nil {
			unauthorized(c, http.StatusBadRequest, "Something goes wrong: "+err.Error())
			return
		}

		user := User{Email: login.Email, Password: login.Password}
		err = user.Get(tx)
		if err != nil {
			unauthorized(c, http.StatusBadRequest, "Something goes wrong: "+err.Error())
			return
		}

		err = tx.Commit()
		if err != nil {
			unauthorized(c, http.StatusBadRequest, "Something goes wrong: "+err.Error())
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
			unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild: "+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":  tokenString,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

// Apply for private routes
func AuthMiddleware(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		token, err := parseToken(c)
		if err != nil {
			unauthorized(c, http.StatusUnauthorized, "Your token is broken: "+err.Error())
			return
		}

		if !token.Valid {
			unauthorized(c, http.StatusBadRequest, "Ivalid token provided")
			return
		}

		userName := token.Claims["username"].(string)
		// Why use it? Dunno...
		// Maybe we will never use it is removed
		// c.Set("Claims", token.Claims)
		c.Set("userName", userName)

		// User will have full acess !?
		// If we need to filter some areas to some user profiles
		// Here we check if this user has this permission
		if false {
			unauthorized(c, http.StatusForbidden, "You don't have permission to access")
			return
		}

		c.Next()
	}
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the AuthMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func RefreshHandler(c *gin.Context) {

	token, err := parseToken(c)
	if err != nil {
		unauthorized(c, http.StatusUnauthorized, "Your token is broken: "+err.Error())
		return
	}

	if !token.Valid {
		unauthorized(c, http.StatusBadRequest, "Ivalid token provided")
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(signingAlgorithm))

	expire := time.Now().Add(timeout)
	newToken.Claims["username"] = token.Claims["username"]
	newToken.Claims["exp"] = expire.Unix()

	tokenString, err := newToken.SignedString(secretKey)
	if err != nil {
		unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
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

func unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "Bearer realm="+relm)

	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()

	return
}
