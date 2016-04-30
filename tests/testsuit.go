package tests

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prixplus/server/database"
	"github.com/prixplus/server/errs"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	router *gin.Engine

	suite.Suite
}

// If this method receive an error not nil
// it logs the errors stack and forces the test fail
func (t *TestSuite) NoError(err error) {
	if err != nil {
		errs.LogError(err)
		t.FailNow(err.Error())
	}
}

// SetUp the test environment
func (t *TestSuite) SetupSuite() {

	fmt.Println("### Initializing new Test Suite")

	// Force to use test.json configs
	os.Setenv("GO_ENV", "test")

	// Load singleton settings
	_, err := settings.Get()
	t.NoError(err)

	// Init DB singleton connection
	_, err = database.Get()
	t.NoError(err)

	// Routing the API
	t.router, err = routers.Init()
	t.NoError(err)

	// Creating temporary schemas
	// It will not insert nothing
	err = CreateTempTables()
	t.NoError(err)
}

func (t *TestSuite) SetupTest() {
	fmt.Println("### Initializing new Test")
	err := InsertTestEntities()
	t.NoError(err)
	fmt.Println("### Starting Test")
}

func (t *TestSuite) TearDownTest() {
	fmt.Println("### Finalizing Test")
	err := TruncateTempTables()
	t.NoError(err)
}

func (t *TestSuite) TearDownSuite() {
	fmt.Println("### Finalizing Test Suite")
	err := database.Close()
	t.NoError(err)
}
