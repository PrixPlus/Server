package middlewares

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/prixplus/server/errs"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	// Do some initialization logic here
	// Foo()
	return func(c *gin.Context) {

		// We will always respond the request
		// even if it panics or is aborted
		defer func() {
			r := recover()

			// Tests if has not had any error in this request
			// it means that request wasn't aborted and stack didn't go in panic
			if r == nil && !c.IsAborted() {
				return
			}

			if r != nil {
				c.AbortWithError(http.StatusInternalServerError, errors.Errorf("Panic! %s", r))
			}

			// Status code -1 == do not override the current status code
			c.JSON(-1, gin.H{
				"errors": c.Errors.Errors(),
			})

		}()

		// Process all next handlers
		c.Next()

		// Logs errors presents in context, if there is any
		errs.LogContextErrors(c)
	}
}
