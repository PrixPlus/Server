package api_tests

import (
	"github.com/prixplus/server/tests"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var router *gin.Engine

type TestSuite struct{}

var _ = Suite(&TestSuite{})

// SetUp the test environment
func (s *TestSuite) SetUpSuite(c *C) {

	// Force to use test.json configs
	os.Setenv("GO_ENV", "test")

	// Load singleton settings
	_, err := settings.Get()
	c.Assert(err, IsNil)

	// Init DB singleton connection
	_, err = database.Get()
	c.Assert(err, IsNil)

	// Routing the API
	router = routers.Init()

	// Creating temporary schemas and insert some tests entities
	tests.InitData()
}

// When all finishes
func (s *TestSuite) TearDownSuite(c *C) {
	// Closing DB singleton connection
	err := database.Close()
	c.Assert(err, IsNil)
}
