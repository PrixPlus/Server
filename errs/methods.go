// Stores API's errors
package errs

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prixplus/server/settings"
)

func LogContextErrors(c *gin.Context) {

	fmt.Println("---------- ERROR ----------")

	f, err := getLogErrorFile()
	if err != nil {
		fmt.Printf("Error openning error log file: %s", err)
		return
	}
	defer f.Close()

	t := time.Now()
	l := fmt.Sprintf("%s - [%s] %s - IP: %s\n", t.Format(time.Stamp), c.Request.Method, c.Request.URL.String(), c.ClientIP())
	f.WriteString(l)
	fmt.Printf(l)

	for _, e := range c.Errors {
		errors.Fprint(f, e.Err)
		errors.Fprint(os.Stdout, e.Err)
	}

	fmt.Println("---------------------------")
}

func LogError(e error) {

	fmt.Println("---------- ERROR ----------")

	f, err := getLogErrorFile()
	if err != nil {
		fmt.Printf("Error openning error log file: %s\n", err)
		return
	}
	defer f.Close()

	t := time.Now()
	l := fmt.Sprintf("%s - System error\n", t.Format(time.Stamp))
	f.WriteString(l)
	fmt.Printf(l)

	errors.Fprint(f, e)
	errors.Fprint(os.Stdout, e)

	fmt.Println("---------------------------")
}

func getLogErrorFile() (*os.File, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting settings")
	}

	filename := sets.Gopath + "/errs.log"

	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
