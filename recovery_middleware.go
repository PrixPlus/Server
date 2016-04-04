package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RecoveryMiddleware(db *sql.DB) gin.HandlerFunc {
	// Do some initialization logic here
	// Foo()
	return func(c *gin.Context) {

		log.Println("Always responds requests")

		defer func() {

			if r := recover(); r != nil || c.IsAborted() {

				log.Println("Abort found")

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
			}

		}()

		c.Next()
	}
}
