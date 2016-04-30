package middlewares

import (
	"fmt"
	"net/http"

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

			// Test if everything is fine
			// it means request wasn't aborted and didn't go into panic
			if r == nil && !c.IsAborted() {
				return
			}

			if r != nil {
				c.Error(fmt.Errorf("Panic! %s", r))
				c.Status(http.StatusInternalServerError)
			}

			// -1 == not override the current error code
			//better than: status := c.Writer.Status()
			c.JSON(-1, gin.H{
				"errors": c.Errors.Errors(),
			})

		}()

		c.Next()

		// Logs the errors present in context
		errs.LogContextErrors(c)
	}
}
