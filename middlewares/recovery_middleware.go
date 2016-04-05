package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Recovery(db *sql.DB) gin.HandlerFunc {
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

			// Ops, something goes wrong
			// maybe bad request, or claims unauthorized...

			log.Println("### Recovery Middleware ###")

			if r != nil {
				c.Error(errors.New(fmt.Sprintf("Panic! %s", r)))
				c.Status(http.StatusInternalServerError)
			}

			status := c.Writer.Status()

			c.JSON(status, gin.H{
				"status":      status,
				"status_text": http.StatusText(status),
				"messages":    c.Errors.Errors(),
			})

		}()

		c.Next()
	}
}
