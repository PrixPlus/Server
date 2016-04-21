package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/auth"
)

// Apply for private routes
func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

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

		token, err := auth.ParseToken(header[7:])
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("Error in token: "+err.Error()))
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
