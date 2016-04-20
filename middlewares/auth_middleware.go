package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/settings"
)

// Apply for private routes
func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		sets, err := settings.Get()
		if err != nil {
			log.Fatal("Error getting Settings: ", err)
		}

		header := c.Request.Header.Get("Authorization")

		// If user isn't trying to identify himself
		// So call the handler and handler will decide
		// if request can or can't be precessed
		if header == "" {
			c.Next()
			return
		}

		// The first word of Authorization header should be Bearer
		if len(header) < 6 || header[0:7] != "Bearer " {
			c.AbortWithError(http.StatusBadRequest, errors.New("Invalid auth header"))
			return
		}

		token, err := jwt.Parse(header[7:], func(token *jwt.Token) (interface{}, error) {
			// Check if encryption algorithm in token is the same
			if jwt.GetSigningMethod(sets.JWT.Algorithm) != token.Method {
				return nil, errors.New("Invalid signing algorithm")
			}

			// I could use some key id to identify whay secret key are we using
			// but it is optional and isn't utilized in this package
			// secretKey, err := myLookupForSecretKey(token.Header["kid"])

			return []byte(sets.JWT.SecretKey), nil
		})

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error parsing token: "+err.Error()))
			return
		}

		if !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Ivalid token provided or it's expired"))
			return
		}

		// Check in DB if token user still logged
		// maybe user logout

		id, ok := token.Claims["id"] // Id comming by float64 !?
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("Ivalid token provided, id not found: %v, Claims: %#v", id, token.Claims)))
			return
		}

		idFloat64, ok := id.(float64)
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error casting Claims"))
			return
		}

		c.Set("id", int64(idFloat64))

		c.Next()
	}
}
