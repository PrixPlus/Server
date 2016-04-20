package models

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/routers"
	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	gin.SetMode(gin.DebugMode)

	// Init DB connection
	db, err := InitDB()
	if err != nil {
		log.Fatal("Error initializing DB: ", err)
	}

	// Close DB when main() returns
	defer db.Close()

	req, _ := http.NewRequest("GET", "/foo", nil)
	resp := httptest.NewRecorder()
	routers.Init(db).ServeHTTP(resp, req)
	assert.Equal(t, resp.Body.String(), "bar")
}
