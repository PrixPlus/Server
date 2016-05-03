package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Settings struct {
	Production bool
	Debug      bool

	DB struct {
		User     string
		Password string
		Host     string
		Name     string
		SSLMode  string
	}

	JWT struct {
		Relm       string
		Algorithm  string
		Expiration time.Duration
		SecretKey  string
	}
	LogFile string

	Env    string `json:-`
	Gopath string `json:-`
	Dir    string `json:-`
}

func (s *Settings) IsProduction() bool {
	return s.Production
}

// Singleton settings
var sets *Settings

func load() (*Settings, error) {

	gopath := os.Getenv("GOPATH")

	dir := gopath + "/src/github.com/prixplus/server/"

	env := os.Getenv("GO_ENV")
	if env == "" {
		fmt.Println("Warning: Setting devlopment environment due to lack of GO_ENV value")
		env = "dev"
	}

	// Reading the JSON file in $GOPATH folder
	filePath := gopath + "/" + env + ".json"
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "reading config file %s", filePath)
	}

	sets = &Settings{}
	err = json.Unmarshal(jsonFile, &sets)
	if err != nil {
		return nil, errors.Wrap(err, "while parsing config file")
	}

	sets.Env = env
	sets.Gopath = gopath
	sets.Dir = dir

	// Ensure that will ever exist the LogFile
	if len(sets.LogFile) == 0 {
		sets.LogFile = "errs.log"
	}

	// If we are in Development
	// so we will set GIN to DebugMode
	if sets.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return sets, nil
}

func Get() (*Settings, error) {
	if sets == nil {
		return load()
	}
	return sets, nil
}
