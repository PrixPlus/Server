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

// Interface for errors used in github.com/pkg/errors
type ErrorLocation interface {
	Location() (string, int)
	error // It still an error
}

func LogContextErrors(c *gin.Context) {
	// Do nothing if there is no errors in this context
	if len(c.Errors) == 0 {
		return
	}

	// If there is errors in context, lets log them all
	fmt.Println("---------- ERROR ----------")
	defer fmt.Println("---------------------------")

	f, err := getLogErrorFile()
	if err != nil {
		errors.Fprint(os.Stdout, errors.Wrap(err, "openning error log file"))
		return
	}
	defer f.Close()

	t := time.Now()
	l := fmt.Sprintf("%s - [%s] %s - IP: %s\n", t.Format(time.Stamp), c.Request.Method, c.Request.URL.String(), c.ClientIP())
	f.WriteString(l)
	fmt.Println(l)

	for _, e := range c.Errors {
		errors.Fprint(f, e.Err)
		errors.Fprint(os.Stdout, e.Err)
	}

}

func LogError(e error) {

	fmt.Println("---------- ERROR ----------")
	defer fmt.Println("---------------------------")

	f, err := getLogErrorFile()
	if err != nil {
		errors.Fprint(os.Stdout, errors.Wrap(err, "openning error log file"))
		return
	}
	defer f.Close()

	t := time.Now()
	l := fmt.Sprintf("%s - System error\n", t.Format(time.Stamp))
	f.WriteString(l)
	fmt.Printf(l)

	errors.Fprint(f, e)
	errors.Fprint(os.Stdout, e)
}

func getLogErrorFile() (*os.File, error) {

	sets, err := settings.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting settings")
	}

	filename := sets.Gopath + "/" + sets.LogFile

	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
