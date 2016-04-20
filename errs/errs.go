// Stores API's errors
package errs

import (
	"errors"
)

var ElementNotFound = errors.New("Element not found")

var Status = map[error]int{ElementNotFound: 471}
