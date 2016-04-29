package tests

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	router *gin.Engine

	suite.Suite
}

// SetUp the test environment
func (t *TestSuite) SetupSuite() {

	fmt.Println("### Initializing new Test Suite")

	// Force to use test.json configs
	os.Setenv("GO_ENV", "test")

	// Load singleton settings
	_, err := settings.Get()
	t.Require().Nil(err, "Err initializing settings")

	// Init DB singleton connection
	_, err = database.Get()
	t.Require().Nil(err, "Err initializing database")

	// Routing the API
	t.router, err = routers.Init()
	t.Require().Nil(err, "Err initializing router")

	// Creating temporary schemas
	// It will not insert nothing
	err = CreateTempTables()
	t.Require().Nil(err, "Err creting temporary tables")
}

func (t *TestSuite) SetupTest() {
	fmt.Println("### Initializing new Test")
	err := InsertTestEntities()
	t.Require().Nil(err, "Err inserting test entities")
	fmt.Println("### Starting Test")
}

func (t *TestSuite) TearDownTest() {
	fmt.Println("### Finalizing Test")
	err := TruncateTempTables()
	t.Require().Nil(err, "Err truncating tables")
}

func (t *TestSuite) TearDownSuite() {
	fmt.Println("### Finalizing Test Suite")
	err := database.Close()
	t.Require().Nil(err, "Err closing database")
}
