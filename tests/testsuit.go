package tests

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prixplus/server/db"
	"github.com/prixplus/server/errs"
	"github.com/prixplus/server/routers"
	"github.com/prixplus/server/settings"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	router  *gin.Engine
	notFail bool

	suite.Suite
}

// SetUp the test environment
func (t *TestSuite) SetupSuite() {
	// fmt.Println("### Initializing new Test Suite")

	// Force to use test.json configs
	os.Setenv("GO_ENV", "test")

	// Load singleton settings
	_, err := settings.Get()
	t.NoError(errors.Wrap(err, "getting settings"))

	// Init DB singleton connection
	_, err = db.Get()
	t.NoError(errors.Wrap(err, "getting db"))

	// Routing the API
	t.router, err = routers.Init()
	t.NoError(errors.Wrap(err, "initializing router"))
}

func (t *TestSuite) SetupTest() {
	// fmt.Println("### Initializing new Test")

	err := DropTempTablesIfExist()
	t.NoError(errors.Wrap(err, "dropping temporary tables if it does exist"))

	err = CreateTempTables()
	t.NoError(errors.Wrap(err, "creating temporary tables"))

	err = InsertTestEntities()
	t.NoError(errors.Wrap(err, "inserting test entities"))

	// fmt.Println("### Starting Test")
}

func (t *TestSuite) TearDownTest() {
	// fmt.Println("### Finalizing Test")

}

func (t *TestSuite) TearDownSuite() {
	// fmt.Println("### Finalizing Test Suite")

	err := db.Close()
	t.NoError(errors.Wrap(err, "closing db connection"))
}

// If this method receive an error not nil
// it logs the errors stack and forces the test fail
// It is important to handle errors using err.LogError
// because it will print all the stack trace returned
func (t *TestSuite) NoError(err error) {
	if err != nil {
		if e, ok := err.(errs.ErrorLocation); ok {
			errs.LogError(e)
		}
		//errs.LogError(err)
		t.FailNow(err.Error())
	}
}
