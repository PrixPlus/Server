package middlewares

import (
	"net/http"

	"github.com/pkg/errors"

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
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid auth header"))
			return
		}

		token, err := auth.ParseToken(header[7:])
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "parsing token"))
			return
		}

		// Check in DB if token user still logged
		// maybe user logout

		uid, ok := token.Claims["uid"] // uid comming by float64 !?
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.Errorf("uid not found int token provided"))
			return
		}

		// This step is just necessary
		// because for some reason uid is stored as float
		idFloat64, ok := uid.(float64)
		if !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("casting uid to float64"))
			return
		}

		c.Set("uid", int64(idFloat64))

		c.Next()
	}
}
